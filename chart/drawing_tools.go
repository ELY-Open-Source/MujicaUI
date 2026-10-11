package chart

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

// DrawingTool selects annotation creation or selection behavior.
type DrawingTool int

const (
	ToolSelect DrawingTool = iota
	ToolTrendLine
	ToolHorizontal
	ToolFibonacci
	ToolRectangle
	ToolText
)

// DrawingPoint anchors annotations to candle time and price.
type DrawingPoint struct {
	Time  time.Time
	Price float64
}

// Drawing is one annotation; one-point tools use A only.
type Drawing struct {
	Tool DrawingTool
	A, B DrawingPoint
	Text string
}

// DrawingTools holds annotations, the current tool, and undo history.
type DrawingTools struct {
	Items []Drawing
	Tool  DrawingTool
	// Text is what the next text annotation says.
	Text string
	// Locked stops selecting, moving, deleting and drawing.
	Locked bool

	sel     int // selected item plus one; zero is none
	pending *Drawing
	undo    [][]Drawing
	moving  bool
}

// fibLevels lists Fibonacci retracement ratios.
var fibLevels = []float64{0, 0.236, 0.382, 0.5, 0.618, 0.786, 1}

// Selected returns the selected drawing's index, or -1.
func (d *DrawingTools) Selected() int { return d.sel - 1 }

// CanUndo reports whether Undo has a step to take back.
func (d *DrawingTools) CanUndo() bool { return len(d.undo) > 0 }

func (d *DrawingTools) save() {
	d.undo = append(d.undo, slices.Clone(d.Items))
}

// Undo restores the preceding annotations unless locked.
func (d *DrawingTools) Undo() bool {
	n := len(d.undo)
	if d.Locked || n == 0 {
		return false
	}
	d.Items, d.undo = d.undo[n-1], d.undo[:n-1]
	d.sel, d.pending = 0, nil
	slog.Debug("mujica: drawing_tools undo", "items", len(d.Items))
	return true
}

// Delete removes the selected drawing, undoably.
func (d *DrawingTools) Delete() bool {
	if d.Locked || d.sel == 0 || d.sel > len(d.Items) {
		return false
	}
	d.save()
	d.Items = slices.Delete(slices.Clone(d.Items), d.sel-1, d.sel)
	d.sel = 0
	slog.Debug("mujica: drawing_tools delete")
	return true
}

// Clear removes every drawing, undoably.
func (d *DrawingTools) Clear() bool {
	if d.Locked || len(d.Items) == 0 {
		return false
	}
	d.save()
	d.Items, d.sel, d.pending = nil, 0, nil
	slog.Debug("mujica: drawing_tools clear")
	return true
}

// twoPoint reports whether t needs a second click.
func (t DrawingTool) twoPoint() bool {
	return t == ToolTrendLine || t == ToolFibonacci || t == ToolRectangle
}

// point maps plot coordinates onto the nearest candle.
func (d *DrawingTools) point(m *mktCore, g mktGeom, x, y float32) DrawingPoint {
	return DrawingPoint{Time: m.candles[g.index(m, x)].Time, Price: m.priceAt(g.ys.Invert(y))}
}

// priceAt resolves an axis value into a rounded price.
func (m *mktCore) priceAt(v float64) float64 {
	if m.base != 0 {
		v = m.base * (1 + v/100)
	}
	p := math.Pow(10, float64(m.dec))
	return math.Round(v*p) / p
}

// xy places a drawing point in the plot.
func (d *DrawingTools) xy(m *mktCore, g mktGeom, pt DrawingPoint) (float32, float32) {
	i := max(0, marketIndexAt(m.candles, pt.Time))
	return g.x(m, i), g.ys.Map(m.val(pt.Price))
}

// input places, selects, moves, and deletes chart annotations.
func (d *DrawingTools) input(m *mktCore, plot *ui.Element, g mktGeom, x, y float32, over, held bool, dx, dy float32) bool {
	if d.Tool < ToolSelect || d.Tool > ToolText {
		panic(fmt.Sprintf("mujica: unknown DrawingTool %d", d.Tool))
	}
	if d.pending != nil && d.pending.Tool != d.Tool {
		d.pending, d.moving = nil, false
	}
	if plot.Shortcut(ui.Alt, ui.KeyEnter) && !d.Locked && len(d.Items) > 0 {
		d.sel = d.sel%len(d.Items) + 1
		m.st.alertFocus = -1
		m.c.Announce(fmt.Sprintf(core.Msg(m.c, "drawing_tools.selected"), d.sel))
	}
	if plot.Shortcut(ui.Cmd, ui.KeyZ) {
		d.Undo()
	}
	if d.Locked {
		d.pending, d.moving = nil, false
		return false
	}
	if plot.Shortcut(0, ui.KeyDelete) || plot.Shortcut(0, ui.KeyBackspace) {
		d.Delete()
	}
	if plot.Shortcut(0, ui.KeyEscape) {
		d.pending, d.sel = nil, 0
	}
	if s := d.sel - 1; s >= 0 && s < len(d.Items) {
		step := map[ui.Key][2]int{ui.KeyLeft: {-1, 0}, ui.KeyRight: {1, 0}, ui.KeyUp: {0, 1}, ui.KeyDown: {0, -1}}
		for key, v := range step {
			if plot.Shortcut(ui.Alt, key) {
				d.save()
				d.nudge(m, s, v[0], float64(v[1])*g.step/10)
			}
		}
	}
	if !over && !held {
		return false
	}
	if plot.Clicked() && d.Tool != ToolSelect {
		pt := d.point(m, g, x, y)
		d.place(m, pt)
		return true
	}
	if d.pending != nil {
		d.pending.B = d.point(m, g, x, y)
		return true
	}
	if d.Tool != ToolSelect {
		return held
	}
	switch {
	case held && !d.moving:
		hit := d.hit(m, g, x-dx, y-dy)
		if hit < 0 {
			d.sel = 0
			return false
		}
		d.sel, d.moving = hit+1, true
		d.save()
	case !held:
		d.moving = false
		if plot.Clicked() {
			d.sel = d.hit(m, g, x, y) + 1
		}
		return false
	}
	if dx != 0 || dy != 0 {
		steps := int(math.Round(float64(dx / g.band.Step())))
		d.nudge(m, d.sel-1, steps, m.priceAt(g.ys.Invert(y))-m.priceAt(g.ys.Invert(y-dy)))
	}
	return true
}

// nudge moves annotations by candles and price.
func (d *DrawingTools) nudge(m *mktCore, i, candles int, price float64) {
	it := &d.Items[i]
	move := func(pt *DrawingPoint) {
		j := max(0, min(len(m.candles)-1, marketIndexAt(m.candles, pt.Time)+candles))
		pt.Time, pt.Price = m.candles[j].Time, pt.Price+price
	}
	move(&it.A)
	move(&it.B)
}

// hit finds the newest annotation within 6 DIPs.
func (d *DrawingTools) hit(m *mktCore, g mktGeom, x, y float32) int {
	for i := len(d.Items) - 1; i >= 0; i-- {
		it := d.Items[i]
		ax, ay := d.xy(m, g, it.A)
		bx, by := d.xy(m, g, it.B)
		near := false
		switch it.Tool {
		case ToolHorizontal:
			near = math.Abs(float64(y-ay)) <= 6
		case ToolText:
			near = math.Abs(float64(x-ax)) <= 40 && math.Abs(float64(y-ay)) <= 10
		case ToolRectangle, ToolFibonacci:
			l, r, t, b := min(ax, bx)-6, max(ax, bx)+6, min(ay, by)-6, max(ay, by)+6
			near = x >= l && x <= r && y >= t && y <= b
		default:
			near = drawingSegmentDist(x, y, ax, ay, bx, by) <= 6
		}
		if near {
			return i
		}
	}
	return -1
}

// drawingSegmentDist measures distance to a line segment.
func drawingSegmentDist(x, y, ax, ay, bx, by float32) float64 {
	vx, vy := float64(bx-ax), float64(by-ay)
	t := 0.0
	if l := vx*vx + vy*vy; l > 0 {
		t = max(0, min(1, (float64(x-ax)*vx+float64(y-ay)*vy)/l))
	}
	return math.Hypot(float64(x-ax)-t*vx, float64(y-ay)-t*vy)
}

func (d *DrawingTools) place(m *mktCore, pt DrawingPoint) {
	if d.Tool.twoPoint() && d.pending == nil {
		d.pending = &Drawing{Tool: d.Tool, A: pt, B: pt}
	} else {
		item := Drawing{Tool: d.Tool, A: pt, B: pt, Text: d.Text}
		if d.pending != nil {
			item.A = d.pending.A
		}
		if d.Tool == ToolText && item.Text == "" {
			item.Text = core.Msg(m.c, "drawing_tools.note")
		}
		d.save()
		d.Items = append(slices.Clone(d.Items), item)
		d.sel, d.pending, d.Tool = len(d.Items), nil, ToolSelect
		slog.Debug("mujica: drawing_tools add", "tool", item.Tool)
	}
	m.c.Invalidate()
}
