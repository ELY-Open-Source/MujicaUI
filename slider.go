package mujica

import (
	"log/slog"
	"math"
	"strconv"
	"strings"

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
	// Format writes the value; nil prints as many decimals as Min and Step need.
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

// sliderKey is the last key released on a slider, and a move waiting for it.
type sliderKey struct {
	up      ui.Key
	pending int
}

type sliderKeyID struct{}

// sliderControl builds a SliderBase on a copy of v and returns where the pointer or a key moves v.
func sliderControl(c *ui.Context, v float64, g sliderGeom, from, to float64) (s *ui.Element, next float64, dragged bool) {
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
	up, dir := ks.up, ks.pending
	ks.up, ks.pending = 0, 0
	if _, _, dragged = s.Dragged(); dragged {
		return s, max(from, min(to, g.snap(raw))), true
	}
	unsure := dir != 0
	switch {
	case raw < v:
		dir, unsure = -1, false
	case raw > v:
		dir, unsure = 1, false
	}
	if dir == 0 {
		return s, v, false
	}
	end, limit := g.lo, from
	if dir > 0 {
		end, limit = g.hi, to
	}
	jump := raw == end
	// An arrow within 1% of an end lands on it too: the released key tells which it was.
	if unsure || jump && math.Abs(v-end) <= (g.hi-g.lo)/100*(1+1e-9) {
		switch up {
		case ui.KeyHome, ui.KeyEnd:
			jump = true
		case ui.KeyLeft, ui.KeyRight, ui.KeyUp, ui.KeyDown:
			jump = false
		default:
			ks.pending = dir
			return s, v, false
		}
	}
	if jump {
		return s, limit, false
	}
	step := g.step
	if step <= 0 {
		step = (g.hi - g.lo) / 100
	}
	return s, max(from, min(to, g.snap(v+float64(dir)*step))), false
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

// sliderDecimals counts the decimals the shortest form of each value needs.
func sliderDecimals(vs ...float64) int {
	n := 0
	for _, v := range vs {
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if i := strings.IndexByte(s, '.'); i >= 0 {
			n = max(n, len(s)-i-1)
		}
	}
	return n
}

func sliderFormat(format func(float64) string, g sliderGeom) func(float64) string {
	if format != nil {
		return format
	}
	prec := sliderDecimals(g.lo, g.step)
	if g.step == 0 {
		prec = max(prec, 1)
	}
	return func(v float64) string { return strconv.FormatFloat(v, 'f', prec, 64) }
}

// Slider sets *value by dragging anywhere on its box, the arrows (one Step), Home and End.
func Slider(c *ui.Context, value *float64, opts SliderOptions) ChoiceResult {
	g, h := sliderSetup(opts.Min, opts.Max, opts.Step, opts.Marks, opts.ShowValue, c)
	format := sliderFormat(opts.Format, g)
	size := FontSize(c, theme.CaptionSize)
	k := Tokens(c)
	old := g.snap(*value)
	*value = old
	var s *ui.Element
	var next float64
	var held bool
	// Disabled on a wrapper blocks the base before it takes input.
	ui.Column(c).Disabled(opts.Disabled).Children(func() {
		s, next, held = sliderControl(c, old, g, g.lo, g.hi)
	})
	s.Height(h).MinWidth(120).PaddingX(sliderThumb / 2).FocusRing(false).Cursor(ui.CursorPointer)
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
		sliderTrackPaint(p, r, g, k, opts.Marks, size, g.x(r, g.lo), g.x(r, v))
		x := g.x(r, v)
		sliderThumbPaint(p, k, x, r.Y+g.trackY, hover || held, s.FocusVisible())
		if opts.ShowValue && (hover || held || s.Focused()) {
			sliderBubblePaint(p, k, x, r.Y+g.trackY, format(v), size)
		}
	})
	return ChoiceResult{s, changed}
}

// sliderTrackPaint draws the rail, the filled part from x0 to x1, and the marks.
func sliderTrackPaint(p *ui.Painter, r ui.Rect, g sliderGeom, k theme.Tokens, marks []SliderMark, size, x0, x1 float32) {
	y := r.Y + g.trackY
	left, right := g.x(r, g.lo), g.x(r, g.hi)
	p.Fill(ui.Rect{X: left, Y: y - sliderTrack/2, W: right - left, H: sliderTrack}, k.Border, 0)
	p.Stroke(ui.Rect{X: left, Y: y - sliderTrack/2, W: right - left, H: sliderTrack}, k.ControlBorder, 0, 1)
	p.Fill(ui.Rect{X: x0, Y: y - sliderTrack/2, W: x1 - x0, H: sliderTrack}, k.AccentText, 0)
	for _, m := range marks {
		x := g.x(r, m.Value)
		p.Fill(ui.Rect{X: x - 0.5, Y: y + sliderTrack/2 + 2, W: 1, H: 5}, k.TextMuted, 0)
		if m.Label != "" {
			w, _ := p.MeasureText(1000, ui.Span{Text: m.Label, Size: size})
			p.Text(x-w/2, y+sliderTrack/2+8, m.Label, size, k.TextMuted)
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
func sliderBubblePaint(p *ui.Painter, k theme.Tokens, x, y float32, s string, size float32) {
	w, h := p.MeasureText(1000, ui.Span{Text: s, Size: size})
	b := ui.Rect{X: x - w/2 - 6, Y: y - sliderThumb/2 - h - 10, W: w + 12, H: h + 4}
	p.Fill(b, k.Surface, theme.ControlRadius)
	p.Stroke(b, k.ControlBorder, theme.ControlRadius, theme.BorderWidth)
	p.Text(b.X+6, b.Y+2, s, size, k.Text)
}
