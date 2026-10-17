package input

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// KnobOptions configure a Knob. Min and Max both zero mean 0 to 100.
type KnobOptions struct {
	Min, Max float64
	// Step snaps the value; zero is continuous, moving 1% per arrow press.
	Step float64
	// Ticks draws that many evenly spaced marks around the dial; zero draws none.
	Ticks int
	// Size is the dial diameter; zero is 56.
	Size float32
	// Format writes the value under the dial; nil prints as many decimals as Min and Step need.
	Format   func(float64) string
	Disabled bool
	// Label names the knob and shows above it.
	Label string
}

// knobSweep is the dial's travel, centered on twelve o'clock, in turns.
const knobSweep = 0.75

// knobDragPixels is how far a vertical drag travels the whole range.
const knobDragPixels float32 = 160

// Knob sets *value by dragging up and down, the arrows, Home and End; the dial shows the value and ticks.
func Knob(c *ui.Context, value *float64, opts KnobOptions) core.ChoiceResult {
	if value == nil {
		panic("mujica: Knob needs a value")
	}
	if opts.Ticks < 0 || !(opts.Size >= 0) || math.IsInf(float64(opts.Size), 1) {
		panic(fmt.Sprintf("mujica: Knob Ticks %d or Size %v invalid", opts.Ticks, opts.Size))
	}
	g, _ := SliderSetup(opts.Min, opts.Max, opts.Step, nil, false, c)
	format := sliderFormat(opts.Format, g)
	size := cmp.Or(opts.Size, 56)
	k := core.Tokens(c)
	old := g.snap(*value)
	*value = old
	col := ui.Column(c).AlignItems(ui.Center).Gap(core.Density(c).Unit()).Disabled(opts.Disabled)
	disabled := col.IsDisabled()
	var s ui.Element
	col.Children(func() {
		if opts.Label != "" {
			ui.Text(c, opts.Label).TextColor(k.Text).FontSize(core.FontSize(c, theme.CaptionSize))
		}
		var next float64
		var dragged bool
		s, next, dragged = SliderControl(c, old, g, g.lo, g.hi)
		want := ui.Local(s, "knob.want", func() float64 { return math.NaN() })
		if dragged {
			if math.IsNaN(*want) {
				*want = old
			}
			_, dy, _ := s.Dragged()
			*want = max(g.lo, min(g.hi, *want-float64(dy)*(g.hi-g.lo)/float64(knobDragPixels)))
			next = g.snap(*want)
		} else {
			*want = math.NaN()
		}
		if !disabled {
			*value = next
		}
		s.Size(size, size).FocusRing(false).Cursor(ui.CursorResizeNS).Label(cmp.Or(opts.Label, core.Msg(c, "knob.label")))
		hot := s.Hovered() || dragged
		v := *value
		s.Draw(func(p *ui.Painter, r ui.Rect) { knobPaint(p, r, k, g, v, opts.Ticks, hot, s.FocusVisible()) })
		ui.Text(c, format(*value)).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum")
	})
	changed := *value != old
	if changed {
		slog.Debug("mujica: knob changed", "label", opts.Label, "value", *value)
	}
	return core.NewChoiceResult(col, changed)
}

// knobPaint draws the track arc, the value arc, ticks and the diamond pointer.
func knobPaint(p *ui.Painter, r ui.Rect, k theme.Tokens, g SliderGeom, v float64, ticks int, hot, focus bool) {
	cx, cy := internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2)
	rad := r.W/2 - 6
	start := float32(-knobSweep / 2)
	frac := float32((v - g.lo) / (g.hi - g.lo))
	p.StrokePath(core.ProgressArc(p, cx, cy, rad, start, start+knobSweep), 3, k.Border)
	if frac > 0 {
		p.StrokePath(core.ProgressArc(p, cx, cy, rad, start, start+knobSweep*frac), 3, k.AccentText)
	}
	for i := range ticks {
		t := float64(start) + knobSweep*float64(i)/float64(max(ticks-1, 1))
		sin, cos := float32(math.Sin(2*math.Pi*t)), float32(math.Cos(2*math.Pi*t))
		p.Line(cx+sin*(rad+3), cy-cos*(rad+3), cx+sin*(rad+6), cy-cos*(rad+6), 1, k.TextMuted)
	}
	face := k.Surface
	if hot {
		face = k.SurfaceHover
	}
	var dial ui.Path
	dial.Circle(cx, cy, rad-6)
	p.FillPath(&dial, face)
	p.StrokePath(&dial, 1, k.ControlBorder)
	t := float64(start + knobSweep*frac)
	px := internal.Snap(p, cx+float32(math.Sin(2*math.Pi*t))*(rad-12))
	py := internal.Snap(p, cy-float32(math.Cos(2*math.Pi*t))*(rad-12))
	p.FillPath(internal.Diamond(px, py, 7), k.AccentText)
	if focus {
		internal.FocusRing(p, r, r.W/2, k.Focus)
	}
}
