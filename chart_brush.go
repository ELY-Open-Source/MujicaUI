package mujica

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChartBrushOptions configure a ChartBrush.
type ChartBrushOptions struct {
	// Label names the brush for assistive technology; required.
	Label string
	// Min and Max span the track; Min < Max.
	Min, Max float64
	// Step is the grid both ends snap to; 0 means a hundredth of the span.
	Step float64
	// Values draws a preview line across the track.
	Values []float64
	// Height is the track height; 0 means 48 DIPs.
	Height   float32
	Format   func(v float64) string
	Disabled bool
}

type brushDrag struct {
	part  int
	grab  float64
	start float64
}

// ChartBrush selects the range *value on a track: drag the window or its ends, or use the keys.
func ChartBrush(c *ui.Context, value *[2]float64, opts ChartBrushOptions) ChoiceResult {
	if opts.Label == "" || !(opts.Max > opts.Min) || opts.Step < 0 || opts.Height < 0 {
		panic("mujica: ChartBrush needs a Label, Min < Max, Step ≥ 0 and Height ≥ 0")
	}
	step := opts.Step
	if step == 0 {
		step = (opts.Max - opts.Min) / 100
	}
	spec := numberInputResolve(NumberInputOptions{Min: opts.Min, Max: opts.Max, Step: step})
	if value[0] < opts.Min || value[1] > opts.Max || !(value[1] > value[0]) {
		panic("mujica: ChartBrush value must be an ascending range inside Min..Max")
	}
	h := opts.Height
	if h == 0 {
		h = 48
	}
	k := Tokens(c)
	u := Density(c).Unit()
	val := func(v float64) string { return chartValue(c, v, opts.Format) }
	snap := func(v float64) float64 {
		return max(spec.lo, min(spec.hi, spec.round(spec.lo+math.Round((v-spec.lo)/spec.step)*spec.step)))
	}
	before := *value
	var res ChoiceResult
	res.Element = ui.Column(c).Gap(u).MinWidth(0).Disabled(opts.Disabled)
	res.Element.Children(func() {
		track := ui.Box(c).Height(h).FillWidth().MinWidth(0).Radius(theme.ControlRadius).Background(k.SurfaceHover)
		b := track.Bounds()
		internal.Remeasure(c, track, b.W, b.H)
		at := func(v float64) float32 { return float32((v - opts.Min) / (opts.Max - opts.Min)) * b.W }
		from := func(px float32) float64 { return opts.Min + float64(px/max(b.W, 1))*(opts.Max-opts.Min) }
		drag := ui.Local(track, "mujica.brush.drag", func() brushDrag { return brushDrag{part: -1} })
		set := func(lo, hi float64) {
			lo, hi = snap(lo), snap(hi)
			if hi-lo >= spec.step-1e-9 {
				*value = [2]float64{lo, hi}
			}
		}
		track.Draw(func(p *ui.Painter, r ui.Rect) { brushPaint(p, r, *value, opts, k) })
		track.Children(func() {
			x0, x1 := at(value[0]), at(value[1])
			window := ui.Box(c).Absolute().Left(x0).Top(0).Width(max(x1-x0, 1)).Height(h).Focusable().
				Role(ui.RoleGroup).Label(fmt.Sprintf(msg(c, "chart_brush.window"), opts.Label, val(value[0]), val(value[1])))
			internal.Ring(window, k.Focus, theme.ControlRadius)
			ends := [2]*ui.Element{}
			for i, x := range []float32{x0, x1} {
				end := ui.Box(c).Absolute().Left(x - theme.MinHitSize/2).Top(0).Size(theme.MinHitSize, h).Focusable().
					Role(ui.RoleSlider).Label(fmt.Sprintf(msg(c, []string{"chart_brush.start", "chart_brush.end"}[i]), val(value[i])))
				internal.Ring(end, k.Focus, theme.ControlRadius)
				end.Draw(func(p *ui.Painter, r ui.Rect) {
					cx := internal.Snap(p, r.X+r.W/2)
					p.Fill(ui.Rect{X: cx - 3, Y: r.Y + r.H/2 - 10, W: 6, H: 20}, k.Text, 2)
				})
				ends[i] = end
			}
			if res.Element.IsDisabled() {
				return
			}
			for i, end := range ends {
				v, ok := brushKey(end, value[i], spec)
				switch {
				case !ok:
				case i == 0:
					set(min(v, value[1]-spec.step), value[1])
				default:
					set(value[0], max(v, value[0]+spec.step))
				}
			}
			w := value[1] - value[0]
			switch {
			case window.Shortcut(0, ui.KeyLeft):
				lo := max(spec.lo, spec.stepBy(value[0], -1))
				set(lo, lo+w)
			case window.Shortcut(0, ui.KeyRight):
				hi := min(spec.hi, spec.stepBy(value[1], 1))
				set(hi-w, hi)
			case window.Shortcut(0, ui.KeyHome):
				set(spec.lo, spec.lo+w)
			case window.Shortcut(0, ui.KeyEnd):
				set(spec.hi-w, spec.hi)
			}
			brushPointer(track, window, ends, drag, value, from, set)
		})
		ui.Text(c, val(value[0])+" – "+val(value[1])).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Role(ui.RoleStatus)
	})
	res.changed = *value != before
	if res.changed {
		slog.Debug("mujica: chart_brush range", "from", value[0], "to", value[1])
	}
	return res
}

// brushKey steps an end for the arrow, Page, Home and End keys.
func brushKey(e *ui.Element, v float64, spec numberInputSpec) (float64, bool) {
	switch {
	case e.Shortcut(0, ui.KeyLeft), e.Shortcut(0, ui.KeyDown):
		return spec.stepBy(v, -1), true
	case e.Shortcut(0, ui.KeyRight), e.Shortcut(0, ui.KeyUp):
		return spec.stepBy(v, 1), true
	case e.Shortcut(0, ui.KeyPageDown):
		return spec.stepBy(v, -10), true
	case e.Shortcut(0, ui.KeyPageUp):
		return spec.stepBy(v, 10), true
	case e.Shortcut(0, ui.KeyHome):
		return spec.lo, true
	case e.Shortcut(0, ui.KeyEnd):
		return spec.hi, true
	}
	return v, false
}

// brushPointer drags an end, the window, or a new range across the track.
func brushPointer(track, window *ui.Element, ends [2]*ui.Element, drag *brushDrag, value *[2]float64, from func(float32) float64, set func(lo, hi float64)) {
	px, _, _ := track.PointerPosition()
	at := from(px)
	parts := []*ui.Element{ends[0], ends[1], window, track}
	held := -1
	for i, e := range parts {
		if _, _, ok := e.Dragged(); ok && held < 0 {
			held = i
		}
	}
	if held < 0 {
		drag.part = -1
		return
	}
	if drag.part < 0 {
		*drag = brushDrag{part: held, grab: at - value[0], start: at}
	}
	w := value[1] - value[0]
	switch drag.part {
	case 0:
		set(min(at, value[1]), value[1])
	case 1:
		set(value[0], max(at, value[0]))
	case 2:
		set(at-drag.grab, at-drag.grab+w)
	default:
		set(min(drag.start, at), max(drag.start, at))
	}
}

// brushPaint draws the preview line, dims outside the range and frames the window.
func brushPaint(p *ui.Painter, r ui.Rect, value [2]float64, opts ChartBrushOptions, k theme.Tokens) {
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, v := range opts.Values {
		if !math.IsNaN(v) {
			lo, hi = min(lo, v), max(hi, v)
		}
	}
	if n := len(opts.Values); n > 1 && hi > lo {
		x := func(i int) float32 { return r.X + r.W*float32(i)/float32(n-1) }
		chartRuns(opts.Values, x, func(idx []int) {
			var line ui.Path
			for j, i := range idx {
				y := r.Y + 4 + (r.H-8)*float32(1-(opts.Values[i]-lo)/(hi-lo))
				if j == 0 {
					line.MoveTo(x(i), y)
				} else {
					line.LineTo(x(i), y)
				}
			}
			p.StrokePath(&line, 1.25, k.TextMuted)
		})
	}
	x0 := r.X + r.W*float32((value[0]-opts.Min)/(opts.Max-opts.Min))
	x1 := r.X + r.W*float32((value[1]-opts.Min)/(opts.Max-opts.Min))
	p.Fill(ui.Rect{X: r.X, Y: r.Y, W: x0 - r.X, H: r.H}, k.Background.Alpha(0.6), 0)
	p.Fill(ui.Rect{X: x1, Y: r.Y, W: r.X + r.W - x1, H: r.H}, k.Background.Alpha(0.6), 0)
	p.Stroke(ui.Rect{X: x0, Y: r.Y, W: x1 - x0, H: r.H}, k.AccentText, 0, 1.5)
}
