package chart

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ChartBrushOptions configure a ChartBrush.
type ChartBrushOptions struct {
	// Label names the brush for assistive technology; required.
	Label string
	// Min and Max span the track; Min < Max.
	Min, Max float64
	// Step defaults to one hundredth of the span.
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

// ChartBrush selects a range with dragging and keyboard controls.
func ChartBrush(c *ui.Context, value *[2]float64, opts ChartBrushOptions) core.ChoiceResult {
	if opts.Label == "" || !(opts.Max > opts.Min) || math.IsInf(opts.Min, 0) || math.IsInf(opts.Max, 0) || math.IsNaN(opts.Step) || math.IsInf(opts.Step, 0) || opts.Step < 0 || math.IsNaN(float64(opts.Height)) || math.IsInf(float64(opts.Height), 0) || opts.Height != 0 && opts.Height < theme.MinHitSize {
		panic("mujica: ChartBrush needs a Label, finite domain and step, and Height ≥ 28")
	}
	step := opts.Step
	if step == 0 {
		step = (opts.Max - opts.Min) / 100
	}
	spec := input.NumberInputResolve(input.NumberInputOptions{Min: opts.Min, Max: opts.Max, Step: step})
	if value[0] < opts.Min || value[1] > opts.Max || !(value[1] > value[0]) {
		panic("mujica: ChartBrush value must be an ascending range inside Min..Max")
	}
	h := opts.Height
	if h == 0 {
		h = 48
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	val := func(v float64) string { return ChartValue(c, v, opts.Format) }
	snap := func(v float64) float64 {
		return max(spec.Lo, min(spec.Hi, spec.Round(spec.Lo+math.Round((v-spec.Lo)/spec.Step)*spec.Step)))
	}
	before := *value
	var res core.ChoiceResult
	res.Element = ui.Column(c).Gap(u).MinWidth(0).Disabled(opts.Disabled)
	res.Element.Children(func() {
		track := ui.Box(c).Height(h).FillWidth().MinWidth(0).Radius(theme.ControlRadius).Background(k.SurfaceHover)
		b := track.Bounds()
		internal.Remeasure(c, track, b.W, b.H)
		at := func(v float64) float32 { return float32((v-opts.Min)/(opts.Max-opts.Min)) * b.W }
		from := func(px float32) float64 { return opts.Min + float64(px/max(b.W, 1))*(opts.Max-opts.Min) }
		drag := ui.Local(track, "mujica.brush.drag", func() brushDrag { return brushDrag{part: -1} })
		set := func(lo, hi float64) {
			lo, hi = snap(lo), snap(hi)
			if hi-lo >= spec.Step-1e-9 {
				*value = [2]float64{lo, hi}
			}
		}
		move := func(lo, hi float64) {
			w := hi - lo
			lo = max(spec.Lo, min(lo, spec.Hi-w))
			set(lo, lo+w)
		}
		preview := ui.Local(track, "mujica.brush.preview", func() brushPreview { return brushPreview{} })
		if id := ChartID(opts.Values); preview.id != id {
			*preview = brushPreview{id: id, lo: math.Inf(1), hi: math.Inf(-1)}
			for _, v := range opts.Values {
				if math.IsInf(v, 0) {
					panic("mujica: ChartBrush preview values must be finite or NaN")
				}
				if !math.IsNaN(v) {
					preview.lo, preview.hi = min(preview.lo, v), max(preview.hi, v)
				}
			}
		}
		track.Draw(func(p *ui.Painter, r ui.Rect) { brushPaint(p, r, *value, opts, k, preview) })
		track.Children(func() {
			x0, x1 := at(value[0]), at(value[1])
			end := func(i int, x float32) *ui.Element {
				e := ui.Box(c).Absolute().Left(x-theme.MinHitSize/2).Top(0).Size(theme.MinHitSize, h).Focusable().
					Role(ui.RoleSlider).Label(fmt.Sprintf(core.Msg(c, []string{"chart_brush.start", "chart_brush.end"}[i]), val(value[i])))
				internal.Ring(e, k.Focus, theme.ControlRadius)
				return e.Draw(func(p *ui.Painter, r ui.Rect) {
					cx := internal.Snap(p, r.X+r.W/2)
					p.Fill(ui.Rect{X: cx - 3, Y: r.Y + r.H/2 - 10, W: 6, H: 20}, k.Text, 2)
				})
			}
			// Tab order runs start, window, end.
			ends := [2]*ui.Element{end(0, x0)}
			window := ui.Box(c).Absolute().Left(x0 + theme.MinHitSize/2).Top(0).Width(max(x1-x0-theme.MinHitSize, theme.MinHitSize)).Height(h).Focusable().
				Role(ui.RoleGroup).Label(fmt.Sprintf(core.Msg(c, "chart_brush.window"), opts.Label, val(value[0]), val(value[1])))
			internal.Ring(window, k.Focus, theme.ControlRadius)
			ends[1] = end(1, x1)
			if res.Element.IsDisabled() {
				drag.part = -1
				return
			}
			for i, end := range ends {
				v, ok := brushKey(end, value[i], spec)
				switch {
				case !ok:
				case i == 0:
					set(min(v, value[1]-spec.Step), value[1])
				default:
					set(value[0], max(v, value[0]+spec.Step))
				}
			}
			w := value[1] - value[0]
			switch {
			case window.Shortcut(0, ui.KeyLeft):
				lo := max(spec.Lo, spec.StepBy(value[0], -1))
				set(lo, lo+w)
			case window.Shortcut(0, ui.KeyRight):
				hi := min(spec.Hi, spec.StepBy(value[1], 1))
				set(hi-w, hi)
			case window.Shortcut(0, ui.KeyHome):
				set(spec.Lo, spec.Lo+w)
			case window.Shortcut(0, ui.KeyEnd):
				set(spec.Hi-w, spec.Hi)
			}
			brushPointer(track, window, ends, drag, value, from, set, move)
		})
		ui.Text(c, val(value[0])+" – "+val(value[1])).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Role(ui.RoleStatus)
	})
	res = core.NewChoiceResult(res.Element, *value != before)
	if res.Changed() {
		c.Announce(fmt.Sprintf(core.Msg(c, "chart_brush.window"), opts.Label, val(value[0]), val(value[1])))
		slog.Debug("mujica: chart_brush range", "from", value[0], "to", value[1])
	}
	return res
}

// brushKey handles arrow, Page, Home and End keys.
func brushKey(e *ui.Element, v float64, spec input.NumberInputSpec) (float64, bool) {
	if e.IsDisabled() {
		return v, false
	}
	switch {
	case e.Shortcut(0, ui.KeyLeft), e.Shortcut(0, ui.KeyDown):
		return spec.StepBy(v, -1), true
	case e.Shortcut(0, ui.KeyRight), e.Shortcut(0, ui.KeyUp):
		return spec.StepBy(v, 1), true
	case e.Shortcut(0, ui.KeyPageDown):
		return spec.StepBy(v, -10), true
	case e.Shortcut(0, ui.KeyPageUp):
		return spec.StepBy(v, 10), true
	case e.Shortcut(0, ui.KeyHome):
		return spec.Lo, true
	case e.Shortcut(0, ui.KeyEnd):
		return spec.Hi, true
	}
	return v, false
}

// brushPointer drags handles, the window, or a new range.
func brushPointer(track, window *ui.Element, ends [2]*ui.Element, drag *brushDrag, value *[2]float64, from func(float32) float64, set, move func(lo, hi float64)) {
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
		lo := at - drag.grab
		// Shift both ends together at a boundary.
		move(lo, lo+w)
	default:
		set(min(drag.start, at), max(drag.start, at))
	}
}

// brushPaint draws the preview and frames the selected range.
type brushPreview struct {
	id     string
	lo, hi float64
	width  float32
	runs   [][]int
}

func brushPaint(p *ui.Painter, r ui.Rect, value [2]float64, opts ChartBrushOptions, k theme.Tokens, preview *brushPreview) {
	lo, hi := preview.lo, preview.hi
	if n := len(opts.Values); n > 1 && hi > lo {
		x := func(i int) float32 { return r.W * float32(i) / float32(n-1) }
		if preview.width != r.W || preview.runs == nil {
			preview.width, preview.runs = r.W, nil
			chartRuns(opts.Values, x, func(idx []int) { preview.runs = append(preview.runs, idx) })
		}
		for _, idx := range preview.runs {
			var line ui.Path
			for j, i := range idx {
				px := internal.Snap(p, r.X+x(i))
				py := internal.Snap(p, r.Y+4+(r.H-8)*float32(1-(opts.Values[i]-lo)/(hi-lo)))
				if j == 0 {
					line.MoveTo(px, py)
				} else {
					line.LineTo(px, py)
				}
			}
			p.StrokePath(&line, 1.25, k.TextMuted)
		}
	}
	x0 := r.X + r.W*float32((value[0]-opts.Min)/(opts.Max-opts.Min))
	x1 := r.X + r.W*float32((value[1]-opts.Min)/(opts.Max-opts.Min))
	p.Fill(ui.Rect{X: r.X, Y: r.Y, W: x0 - r.X, H: r.H}, k.Background.Alpha(0.6), 0)
	p.Fill(ui.Rect{X: x1, Y: r.Y, W: r.X + r.W - x1, H: r.H}, k.Background.Alpha(0.6), 0)
	p.Stroke(ui.Rect{X: x0, Y: r.Y, W: x1 - x0, H: r.H}, k.AccentText, 0, 1.5)
}
