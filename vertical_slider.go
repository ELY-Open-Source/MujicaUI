package mujica

import (
	"cmp"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// VerticalSliderOptions configure a VerticalSlider. Min and Max both zero mean 0 to 100.
type VerticalSliderOptions struct {
	Min, Max float64
	// Step snaps the value; zero is continuous.
	Step  float64
	Marks []SliderMark
	// Height is the track length; zero is 160.
	Height float32
	// ShowValue shows the value beside the thumb while hovered, dragged or focused.
	ShowValue bool
	Format    func(float64) string
	Disabled  bool
	Label     string
}

// VerticalSlider is a Slider standing up: max at the top, Up and Right raise it.
func VerticalSlider(c *ui.Context, value *float64, opts VerticalSliderOptions) ChoiceResult {
	if value == nil {
		panic("mujica: VerticalSlider needs a value")
	}
	if !(opts.Height >= 0) || math.IsInf(float64(opts.Height), 1) {
		panic("mujica: VerticalSlider Height must be finite and >= 0")
	}
	g, _ := sliderSetup(opts.Min, opts.Max, opts.Step, opts.Marks, false, c)
	format := sliderFormat(opts.Format, g)
	k := Tokens(c)
	size := FontSize(c, theme.CaptionSize)
	h := cmp.Or(opts.Height, 160)
	old := g.snap(*value)
	*value = old
	var s *ui.Element
	var next float64
	var held bool
	wrap := ui.Column(c).Disabled(opts.Disabled)
	wrap.Children(func() { s, next, held = sliderControl(c, old, g, g.lo, g.hi) })
	// The base maps the pointer across; here the height decides.
	y := func(r ui.Rect, v float64) float32 {
		return r.Y + r.H - sliderThumb/2 - (r.H-sliderThumb)*float32((v-g.lo)/(g.hi-g.lo))
	}
	if held {
		if _, py, _ := s.PointerPosition(); s.Bounds().H > sliderThumb {
			f := 1 - float64((py-sliderThumb/2)/(s.Bounds().H-sliderThumb))
			next = g.snap(g.lo + max(0, min(1, f))*(g.hi-g.lo))
		}
	}
	w := Density(c).ControlHeight()
	for _, m := range opts.Marks {
		if m.Label != "" {
			w += 40
			break
		}
	}
	s.Size(w, h).PaddingY(sliderThumb / 2).FocusRing(false).Cursor(ui.CursorPointer)
	if opts.Label != "" {
		s.Label(opts.Label)
	}
	if !wrap.IsDisabled() {
		*value = next
	}
	changed := *value != old
	if changed {
		slog.Debug("mujica: vertical_slider changed", "label", opts.Label, "value", *value)
	}
	hot := s.Hovered() || held
	v := *value
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		x := r.X + Density(c).ControlHeight()/2
		top, bot := y(r, g.hi), y(r, g.lo)
		p.Fill(ui.Rect{X: x - sliderTrack/2, Y: top, W: sliderTrack, H: bot - top}, k.Border, 0)
		p.Stroke(ui.Rect{X: x - sliderTrack/2, Y: top, W: sliderTrack, H: bot - top}, k.ControlBorder, 0, 1)
		p.Fill(ui.Rect{X: x - sliderTrack/2, Y: y(r, v), W: sliderTrack, H: bot - y(r, v)}, k.AccentText, 0)
		for _, m := range opts.Marks {
			my := y(r, m.Value)
			p.Fill(ui.Rect{X: x + sliderTrack/2 + 2, Y: my - 0.5, W: 5, H: 1}, k.TextMuted, 0)
			if m.Label != "" {
				_, th := p.MeasureText(1000, ui.Span{Text: m.Label, Size: size})
				p.Text(x+sliderTrack/2+10, my-th/2, m.Label, size, k.TextMuted)
			}
		}
		sliderThumbPaint(p, k, x, y(r, v), hot, s.FocusVisible())
		if opts.ShowValue && (hot || s.Focused()) {
			sliderBubblePaint(p, k, x, y(r, v), format(v), size)
		}
	})
	return ChoiceResult{wrap, changed}
}
