package input

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
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
	// Format writes the value; nil prints as many decimals as Min and Step need.
	Format   func(float64) string
	Disabled bool
	// Label names the slider for assistive technology.
	Label string
}

const (
	SliderThumb  float32 = 16
	SliderTrack  float32 = 4
	sliderBubble float32 = 26
	sliderMarkH  float32 = 18
)

// SliderGeom maps values to x positions across a slider's box; grid holds the step points.
type SliderGeom struct {
	lo, hi, step float64
	grid         NumberInputSpec
	TrackY       float32
}

func (g SliderGeom) X(r ui.Rect, v float64) float32 {
	f := float32((v - g.lo) / (g.hi - g.lo))
	return r.X + SliderThumb/2 + (r.W-SliderThumb)*max(0, min(1, f))
}

// snap moves v to the nearest step point inside the range.
func (g SliderGeom) snap(v float64) float64 {
	v = max(g.lo, min(g.hi, v))
	if g.step <= 0 {
		return v
	}
	return g.grid.StepBy(v, 0)
}

// move goes n steps from v, 1% of the range each when continuous, staying in range.
func (g SliderGeom) move(v float64, n int) float64 {
	if g.step <= 0 {
		return max(g.lo, min(g.hi, v+float64(n)*(g.hi-g.lo)/100))
	}
	return g.grid.StepBy(v, n)
}

// sliderKey is the last key released on a slider and how many presses MyGo saw before it.
type sliderKey struct {
	up      ui.Key
	presses int
}

type sliderKeyID struct{}

// SliderControl builds a SliderBase on a copy of v and returns where the pointer or a key moves v.
func SliderControl(c *ui.Context, v float64, g SliderGeom, from, to float64) (s *ui.Element, next float64, dragged bool) {
	raw := v
	s = ui.SliderBase(c, &raw, g.lo, g.hi)
	ks := *ui.Local(s, sliderKeyID{}, func() *sliderKey { return &sliderKey{} })
	// MyGo gives key downs to the base alone (ui/base.go:118); releases still reach this handler.
	s.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputKeyUp {
			return false
		}
		ks.up = ev.Key
		return true
	})
	up := ks.up
	ks.up = 0
	if _, _, dragged = s.Dragged(); dragged {
		ks.presses = 0
		return s, max(from, min(to, g.snap(raw))), true
	}
	if raw != v {
		// A key press MyGo applied to the copy; its rounded result (ui/widgets.go:375) says nothing of the direction.
		ks.presses++
	}
	if up == 0 {
		return s, v, false
	}
	n := max(1, ks.presses)
	ks.presses = 0
	switch up {
	case ui.KeyLeft, ui.KeyDown:
		return s, max(from, min(to, g.move(v, -n))), false
	case ui.KeyRight, ui.KeyUp:
		return s, max(from, min(to, g.move(v, n))), false
	case ui.KeyHome:
		return s, g.snap(from), false
	case ui.KeyEnd:
		return s, g.snap(to), false
	}
	return s, v, false
}

func SliderSetup(min, max, step float64, marks []SliderMark, showValue bool, c *ui.Context) (SliderGeom, float32) {
	if min == 0 && max == 0 {
		max = 100
	}
	if max <= min || step < 0 {
		panic("mujica: slider needs Min < Max and Step >= 0")
	}
	g := SliderGeom{lo: min, hi: max, step: step}
	if step > 0 {
		// NumberInput's grid: exact bounds and integer step indexes at the decimals Min and Step need.
		g.grid = NumberInputResolve(NumberInputOptions{Min: min, Max: max, Step: step, Precision: SliderDecimals(min, step)})
	}
	h := core.Density(c).ControlHeight()
	g.TrackY = h / 2
	if showValue {
		g.TrackY += sliderBubble
		h += sliderBubble
	}
	for _, m := range marks {
		if m.Label != "" {
			return g, h + sliderMarkH
		}
	}
	return g, h
}

// SliderDecimals counts the decimals the shortest form of each value needs.
func SliderDecimals(vs ...float64) int {
	n := 0
	for _, v := range vs {
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if i := strings.IndexByte(s, '.'); i >= 0 {
			n = max(n, len(s)-i-1)
		}
	}
	return n
}

func sliderFormat(format func(float64) string, g SliderGeom) func(float64) string {
	if format != nil {
		return format
	}
	prec := SliderDecimals(g.lo, g.step)
	if g.step == 0 {
		prec = max(prec, 1)
	}
	return func(v float64) string { return strconv.FormatFloat(v, 'f', prec, 64) }
}

// Slider sets *value by dragging anywhere on its box, the arrows (one Step), Home and End.
func Slider(c *ui.Context, value *float64, opts SliderOptions) core.ChoiceResult {
	g, h := SliderSetup(opts.Min, opts.Max, opts.Step, opts.Marks, opts.ShowValue, c)
	format := sliderFormat(opts.Format, g)
	size := core.FontSize(c, theme.CaptionSize)
	k := core.Tokens(c)
	old := g.snap(*value)
	*value = old
	var s *ui.Element
	var next float64
	var held bool
	// Disabled on a wrapper blocks the base before it takes input.
	ui.Column(c).Disabled(opts.Disabled).Children(func() {
		s, next, held = SliderControl(c, old, g, g.lo, g.hi)
	})
	s.Height(h).MinWidth(120).PaddingX(SliderThumb / 2).FocusRing(false).Cursor(ui.CursorPointer)
	if opts.Label != "" {
		s.Label(opts.Label)
	}
	if !opts.Disabled {
		*value = next
	}
	changed := *value != old
	if changed {
		slog.Debug("mujica: slider changed", "label", opts.Label, "value", *value)
	}
	hover := s.Hovered()
	v := *value
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		SliderTrackPaint(p, r, g, k, opts.Marks, size, g.X(r, g.lo), g.X(r, v))
		x := g.X(r, v)
		SliderThumbPaint(p, k, x, r.Y+g.TrackY, hover || held, s.FocusVisible())
		if opts.ShowValue && (hover || held || s.Focused()) {
			sliderBubblePaint(p, k, x, r.Y+g.TrackY, format(v), size)
		}
	})
	return core.NewChoiceResult(s, changed)
}

// SliderTrackPaint draws the rail, the filled part from x0 to x1, and the marks.
func SliderTrackPaint(p *ui.Painter, r ui.Rect, g SliderGeom, k theme.Tokens, marks []SliderMark, size, x0, x1 float32) {
	y := r.Y + g.TrackY
	left, right := g.X(r, g.lo), g.X(r, g.hi)
	p.Fill(ui.Rect{X: left, Y: y - SliderTrack/2, W: right - left, H: SliderTrack}, k.Border, 0)
	p.Stroke(ui.Rect{X: left, Y: y - SliderTrack/2, W: right - left, H: SliderTrack}, k.ControlBorder, 0, 1)
	p.Fill(ui.Rect{X: x0, Y: y - SliderTrack/2, W: x1 - x0, H: SliderTrack}, k.AccentText, 0)
	for _, m := range marks {
		x := g.X(r, m.Value)
		p.Fill(ui.Rect{X: x - 0.5, Y: y + SliderTrack/2 + 2, W: 1, H: 5}, k.TextMuted, 0)
		if m.Label != "" {
			w, _ := p.MeasureText(1000, ui.Span{Text: m.Label, Size: size})
			p.Text(x-w/2, y+SliderTrack/2+8, m.Label, size, k.TextMuted)
		}
	}
}

// SliderThumbPaint draws the diamond thumb at (x, y), with its focus ring.
func SliderThumbPaint(p *ui.Painter, k theme.Tokens, x, y float32, hot, focus bool) {
	x, y = internal.Snap(p, x), internal.Snap(p, y)
	size := SliderThumb
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
func sliderBubblePaint(p *ui.Painter, k theme.Tokens, x, y float32, s string, size float32) {
	w, h := p.MeasureText(1000, ui.Span{Text: s, Size: size})
	b := ui.Rect{X: x - w/2 - 6, Y: y - SliderThumb/2 - h - 10, W: w + 12, H: h + 4}
	p.Fill(b, k.Surface, theme.ControlRadius)
	p.Stroke(b, k.ControlBorder, theme.ControlRadius, theme.BorderWidth)
	p.Text(b.X+6, b.Y+2, s, size, k.Text)
}
