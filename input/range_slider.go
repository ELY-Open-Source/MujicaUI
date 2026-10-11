package input

import (
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// RangeSliderOptions configure a RangeSlider. Min and Max both zero mean 0 to 100.
type RangeSliderOptions struct {
	Min, Max float64
	// Step snaps both ends to Min plus multiples of it; zero is continuous.
	Step      float64
	Marks     []SliderMark
	ShowValue bool
	Format    func(float64) string
	Disabled  bool
	// Label names the range; the thumbs read "<Label> minimum/maximum, value".
	Label string
}

// RangeSlider sets *low and *high, which never cross; a press moves the nearer thumb, each takes the keys.
func RangeSlider(c *ui.Context, low, high *float64, opts RangeSliderOptions) core.ChoiceResult {
	g, h := SliderSetup(opts.Min, opts.Max, opts.Step, opts.Marks, opts.ShowValue, c)
	format := sliderFormat(opts.Format, g)
	size := core.FontSize(c, theme.CaptionSize)
	k := core.Tokens(c)
	*low, *high = g.snap(math.Min(*low, *high)), g.snap(math.Max(*low, *high))
	oldLow, oldHigh := *low, *high
	e := ui.Box(c).Height(h).MinWidth(120).Shrink(0).Role(ui.RoleGroup).Disabled(opts.Disabled).Cursor(ui.CursorPointer)
	if opts.Label != "" {
		e.Label(opts.Label)
	}
	drag := ui.Local(e, rangeSliderDrag{}, func() int { return -1 })
	px, _, _ := e.PointerPosition()
	r := e.Bounds()
	at := g.snap(g.lo + float64((px-SliderThumb/2)/max(r.W-SliderThumb, 1))*(g.hi-g.lo))
	var knobs [2]*ui.Element
	pressed := false
	e.Children(func() {
		for i, v := range []*float64{low, high} {
			from, to := g.lo, *high
			if i == 1 {
				from, to = *low, g.hi
			}
			x := g.X(ui.Rect{W: r.W}, *v) - SliderThumb/2
			var kb *ui.Element
			var next float64
			var held bool
			// Each thumb is a SliderBase, which reports the range and value to assistive technology.
			ui.Box(c).Key(i).Absolute().Left(x).Top(g.TrackY-SliderThumb/2).Size(SliderThumb, SliderThumb).Children(func() {
				kb, next, held = SliderControl(c, *v, g, from, to)
			})
			kb.Size(SliderThumb, SliderThumb).FocusRing(false)
			name := core.Msg(c, []string{"range_slider.min", "range_slider.max"}[i])
			if opts.Label != "" {
				name = opts.Label + " " + name
			}
			kb.Label(name)
			if held {
				pressed = true
				*drag = i
			} else if !opts.Disabled {
				*v = next
			}
			knobs[i] = kb
		}
	})
	if _, _, ok := e.Dragged(); ok {
		pressed = true
	}
	if !pressed || opts.Disabled {
		*drag = -1
	} else {
		if *drag < 0 {
			*drag = 0
			if math.Abs(at-*high) < math.Abs(at-*low) || at > *high {
				*drag = 1
			}
		}
		if *drag == 0 {
			*low = min(at, *high)
		} else {
			*high = max(at, *low)
		}
		knobs[*drag].Focus()
	}
	changed := *low != oldLow || *high != oldHigh
	if changed {
		slog.Debug("mujica: range slider changed", "label", opts.Label, "low", *low, "high", *high)
	}
	lo, hi, which := *low, *high, *drag
	hover := e.Hovered()
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		x0, x1 := g.X(r, lo), g.X(r, hi)
		SliderTrackPaint(p, r, g, k, opts.Marks, size, x0, x1)
		for i, x := range []float32{x0, x1} {
			hot := which == i || knobs[i].Hovered()
			SliderThumbPaint(p, k, x, r.Y+g.TrackY, hot, knobs[i].FocusVisible())
			if opts.ShowValue && (hot || hover || knobs[i].Focused()) {
				rangeSliderBubble(p, k, x, r.Y+g.TrackY, format([]float64{lo, hi}[i]), i, size)
			}
		}
	})
	return core.NewChoiceResult(e, changed)
}

// rangeSliderDrag keys which thumb the pointer drags, -1 for none.
type rangeSliderDrag struct{}

// rangeSliderBubble keeps the two tips apart: the low one leans left, the high one right.
func rangeSliderBubble(p *ui.Painter, k theme.Tokens, x, y float32, s string, i int, size float32) {
	w, _ := p.MeasureText(1000, ui.Span{Text: s, Size: size})
	shift := (w+12)/2 - SliderThumb/2
	if i == 0 {
		shift = -shift
	}
	sliderBubblePaint(p, k, x+shift, y, s, size)
}
