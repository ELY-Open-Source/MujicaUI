package mujica

import (
	"log/slog"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// SliderMark is a tick on a slider's track, with an optional caption.
type SliderMark struct {
	Value float64
	Label string
}

// SliderOptions configure a Slider. Min and Max both zero mean 0 to 100.
type SliderOptions struct {
	Min, Max float64
	// Step snaps the value to Min plus multiples of it; zero is continuous.
	Step  float64
	Marks []SliderMark
	// ShowValue shows the value above the thumb while hovered, dragged or focused.
	ShowValue bool
	// Format writes the value; nil prints it with the step's precision.
	Format   func(float64) string
	Disabled bool
	// Label names the slider for assistive technology.
	Label string
}

const (
	sliderThumb  float32 = 16
	sliderTrack  float32 = 4
	sliderBubble float32 = 26
	sliderMarkH  float32 = 18
)

// sliderGeom maps values to x positions across a slider's box.
type sliderGeom struct {
	lo, hi, step float64
	trackY       float32
}

func (g sliderGeom) x(r ui.Rect, v float64) float32 {
	f := float32((v - g.lo) / (g.hi - g.lo))
	return r.X + sliderThumb/2 + (r.W-sliderThumb)*max(0, min(1, f))
}

// snap rounds v to the nearest step inside the range.
func (g sliderGeom) snap(v float64) float64 {
	v = max(g.lo, min(g.hi, v))
	if g.step <= 0 {
		return v
	}
	v = g.lo + math.Round((v-g.lo)/g.step)*g.step
	if v > g.hi {
		v -= g.step
	}
	return math.Round(v*1e9) / 1e9
}

// keyed returns the value after a key moved old to raw, MyGo's 1% key step,
// as one step of the slider: Home and End land on the ends.
func (g sliderGeom) keyed(old, raw float64) float64 {
	switch {
	case g.step <= 0:
		return raw
	case raw == g.lo:
		return g.lo
	case raw == g.hi:
		return g.snap(g.hi)
	case raw > old:
		return g.snap(min(g.hi, old+g.step))
	}
	return g.snap(old - g.step)
}

func sliderSetup(min, max, step float64, marks []SliderMark, showValue bool, c *ui.Context) (sliderGeom, float32) {
	if min == 0 && max == 0 {
		max = 100
	}
	if max <= min || step < 0 {
		panic("mujica: slider needs Min < Max and Step >= 0")
	}
	g := sliderGeom{lo: min, hi: max, step: step}
	h := Density(c).ControlHeight()
	g.trackY = h / 2
	if showValue {
		g.trackY += sliderBubble
		h += sliderBubble
	}
	for _, m := range marks {
		if m.Label != "" {
			return g, h + sliderMarkH
		}
	}
	return g, h
}

func sliderFormat(format func(float64) string, step float64) func(float64) string {
	if format != nil {
		return format
	}
	prec := 0
	if step > 0 && step < 1 {
		prec = int(math.Ceil(-math.Log10(step) - 1e-9))
	} else if step == 0 {
		prec = 1
	}
	return func(v float64) string { return strconv.FormatFloat(v, 'f', prec, 64) }
}

// Slider sets *value between Min and Max by dragging anywhere on its box,
// the arrows (one Step), Home and End. The thumb is a diamond; the hit area
// is the whole rectangle.
func Slider(c *ui.Context, value *float64, opts SliderOptions) ChoiceResult {
	g, h := sliderSetup(opts.Min, opts.Max, opts.Step, opts.Marks, opts.ShowValue, c)
	format := sliderFormat(opts.Format, opts.Step)
	k := Tokens(c)
	old := g.snap(*value)
	*value = old
	target := value
	if opts.Disabled {
		// SliderBase takes the pointer as it is built: a disabled one moves a copy.
		scratch := old
		target = &scratch
	}
	s := ui.SliderBase(c, target, g.lo, g.hi).Height(h).MinWidth(120).PaddingX(sliderThumb / 2).
		FocusRing(false).Disabled(opts.Disabled).Cursor(ui.CursorPointer)
	if opts.Label != "" {
		s.Label(opts.Label)
	}
	if *value != old {
		if _, _, drag := s.Dragged(); drag {
			*value = g.snap(*value)
		} else {
			*value = g.keyed(old, *value)
		}
	}
	changed := *value != old
	if changed {
		slog.Debug("mujica: slider changed", "label", opts.Label, "value", *value)
	}
	_, _, held := s.Dragged()
	hover := s.Hovered()
	v := *value
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		sliderTrackPaint(p, r, g, k, opts.Marks, format, g.x(r, g.lo), g.x(r, v))
		x := g.x(r, v)
		sliderThumbPaint(p, k, x, r.Y+g.trackY, hover || held, s.FocusVisible())
		if opts.ShowValue && (hover || held || s.Focused()) {
			sliderBubblePaint(p, k, x, r.Y+g.trackY, format(v))
		}
	})
	return ChoiceResult{s, changed}
}

// sliderTrackPaint draws the rail, the filled part from x0 to x1, and the marks.
func sliderTrackPaint(p *ui.Painter, r ui.Rect, g sliderGeom, k theme.Tokens, marks []SliderMark, format func(float64) string, x0, x1 float32) {
	y := r.Y + g.trackY
	left, right := g.x(r, g.lo), g.x(r, g.hi)
	p.Fill(ui.Rect{X: left, Y: y - sliderTrack/2, W: right - left, H: sliderTrack}, k.Border, 0)
	p.Stroke(ui.Rect{X: left, Y: y - sliderTrack/2, W: right - left, H: sliderTrack}, k.ControlBorder, 0, 1)
	p.Fill(ui.Rect{X: x0, Y: y - sliderTrack/2, W: x1 - x0, H: sliderTrack}, k.AccentText, 0)
	for _, m := range marks {
		x := g.x(r, m.Value)
		p.Fill(ui.Rect{X: x - 0.5, Y: y + sliderTrack/2 + 2, W: 1, H: 5}, k.TextMuted, 0)
		if m.Label != "" {
			w, _ := p.MeasureText(1000, ui.Span{Text: m.Label, Size: theme.CaptionSize})
			p.Text(x-w/2, y+sliderTrack/2+8, m.Label, theme.CaptionSize, k.TextMuted)
		}
	}
}

// sliderThumbPaint draws the diamond thumb at (x, y), with its focus ring.
func sliderThumbPaint(p *ui.Painter, k theme.Tokens, x, y float32, hot, focus bool) {
	x, y = internal.Snap(p, x), internal.Snap(p, y)
	size := sliderThumb
	if hot {
		size += 2
	}
	p.FillPath(internal.Diamond(x, y, size+3), k.Surface)
	p.FillPath(internal.Diamond(x, y, size), k.AccentText)
	if focus {
		p.StrokePath(internal.Diamond(x, y, size+3+2*(theme.FocusGap+theme.FocusWidth)), theme.FocusWidth, k.Focus)
	}
}

// sliderBubblePaint draws the value tip above the thumb at (x, y).
func sliderBubblePaint(p *ui.Painter, k theme.Tokens, x, y float32, s string) {
	w, h := p.MeasureText(1000, ui.Span{Text: s, Size: theme.CaptionSize})
	b := ui.Rect{X: x - w/2 - 6, Y: y - sliderThumb/2 - h - 10, W: w + 12, H: h + 4}
	p.Fill(b, k.Surface, theme.ControlRadius)
	p.Stroke(b, k.ControlBorder, theme.ControlRadius, theme.BorderWidth)
	p.Text(b.X+6, b.Y+2, s, theme.CaptionSize, k.Text)
}
