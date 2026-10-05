package mujica

import (
	"fmt"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ProgressShape is how a progress indicator draws. The zero value is ProgressLinear.
type ProgressShape int

const (
	ProgressLinear ProgressShape = iota
	ProgressRing
)

// ProgressOptions configure a Progress.
type ProgressOptions struct {
	Shape ProgressShape
	// Indeterminate shows work of unknown length; the value is ignored.
	Indeterminate bool
	// Label names the indicator and shows beside it; empty reads "Loading".
	Label     string
	ShowValue bool
	// Size is the ring diameter (48 DIP when zero) or the bar width (fills when zero).
	Size float32
}

// Progress shows work done, value 0 to 1; indeterminate ones hold still under reduced motion.
func Progress(c *ui.Context, value float64, opts ProgressOptions) *ui.Element {
	if opts.Shape != ProgressLinear && opts.Shape != ProgressRing {
		panic(fmt.Sprintf("mujica: unknown ProgressShape %d", int(opts.Shape)))
	}
	if !opts.Indeterminate && (math.IsNaN(value) || value < 0 || value > 1) {
		panic(fmt.Sprintf("mujica: Progress value %v outside 0..1", value))
	}
	if !(opts.Size >= 0) || math.IsInf(float64(opts.Size), 1) {
		panic(fmt.Sprintf("mujica: Progress Size %v invalid", opts.Size))
	}
	k := Tokens(c)
	u := Density(c).Unit()
	still := c.Preferences().ReduceMotion
	name := opts.Label
	if name == "" {
		name = msg(c, "loading")
	}
	acc := value
	if opts.Indeterminate {
		acc = -1
	}
	caption := ""
	switch {
	case opts.Indeterminate && still:
		caption = msg(c, "progress.indeterminate")
	case !opts.Indeterminate && opts.ShowValue:
		caption = strconv.Itoa(int(math.Round(value*100))) + "%"
	}
	if opts.Shape == ProgressRing {
		return progressRing(c, k, u, acc, name, caption, opts, still)
	}
	col := ui.Column(c).Gap(u * 1.5)
	if opts.Size > 0 {
		col.Width(opts.Size)
	} else {
		col.FillWidth()
	}
	col.Children(func() {
		if opts.Label != "" || caption != "" {
			ui.Row(c).Gap(u * 2).Children(func() {
				ui.Text(c, opts.Label).TextColor(k.Text).Grow(1).Shrink(1)
				ui.Text(c, caption).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize)).FontFeatures("tnum").SingleLine()
			})
		}
		bar := ui.Progress(c, acc).Label(name).FillWidth().Height(u * 1.5).Radius(theme.ControlRadius).Background(k.Border)
		shown := float32(value)
		if !opts.Indeterminate {
			shown = progressAnimate(bar, shown, still)
		}
		bar.Draw(func(p *ui.Painter, r ui.Rect) {
			switch {
			case !opts.Indeterminate:
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * shown, H: r.H}, k.AccentText, theme.ControlRadius)
			case still:
				p.Fill(r, k.AccentText, theme.ControlRadius)
			default:
				p.AnimationFrame()
				phase := float32(p.Now().UnixMilli()%1400) / 1400
				w := r.W * 0.3
				x := r.X - w + (r.W+w)*phase
				p.Clip(r, theme.ControlRadius, func() {
					p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, k.AccentText, theme.ControlRadius)
				})
			}
		})
	})
	return col
}

func progressRing(c *ui.Context, k theme.Tokens, u float32, acc float64, name, caption string, opts ProgressOptions, still bool) *ui.Element {
	d := opts.Size
	if d == 0 {
		d = 48
	}
	row := ui.Row(c).AlignItems(ui.Center).Gap(u * 2)
	row.Children(func() {
		ring := ui.Progress(c, acc).Label(name).Size(d, d).Radius(0).Background(ui.Transparent)
		shown := float32(acc)
		if !opts.Indeterminate {
			shown = progressAnimate(ring, shown, still)
		}
		stroke := max(3, d/10)
		size := FontSize(c, theme.CaptionSize)
		tw, th := c.MeasureText(0, ui.Span{Text: caption, Size: size})
		inside := caption != "" && !opts.Indeterminate && tw+4 <= d-2*stroke && th <= d-2*stroke
		ring.Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy, rad := internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2), internal.Snap(p, (r.W-stroke)/2)
			var track ui.Path
			track.Circle(cx, cy, rad)
			p.StrokePath(&track, stroke, k.Border)
			switch {
			case !opts.Indeterminate:
				if shown > 0 {
					p.StrokePath(progressArc(p, cx, cy, rad, 0, shown), stroke, k.AccentText)
				}
			case still:
				p.StrokePath(&track, stroke, k.AccentText)
			default:
				p.AnimationFrame()
				start := float32(p.Now().UnixMilli()%1200) / 1200
				p.StrokePath(progressArc(p, cx, cy, rad, start, start+0.25), stroke, k.AccentText)
			}
		})
		if inside {
			ring.Children(func() {
				ui.Text(c, caption).FontSize(size).FontFeatures("tnum").TextColor(k.Text)
			})
			ring.Center()
		}
		if opts.Label != "" || (caption != "" && !inside) {
			ui.Column(c).Children(func() {
				if opts.Label != "" {
					ui.Text(c, opts.Label).TextColor(k.Text)
				}
				if caption != "" && !inside {
					ui.Text(c, caption).TextColor(k.TextMuted).FontSize(size).FontFeatures("tnum")
				}
			})
		}
	})
	return row
}

// progressAnimate eases toward target; under reduced motion it lands there and keeps the state in step.
func progressAnimate(e *ui.Element, target float32, still bool) float32 {
	if still {
		// A changed target lands at once under reduced motion; the detour stops a running transition.
		e.Animate("progress.value", target+1, theme.StateDuration)
	}
	return e.Animate("progress.value", target, theme.StateDuration)
}

// progressArc is the arc from turn a to turn b, clockwise from twelve o'clock.
func progressArc(pt *ui.Painter, cx, cy, r float32, a, b float32) *ui.Path {
	var p ui.Path
	n := max(2, int(math.Ceil(float64((b-a)*64))))
	for i := 0; i <= n; i++ {
		t := float64(a + (b-a)*float32(i)/float32(n))
		x := cx + r*float32(math.Sin(2*math.Pi*t))
		y := cy - r*float32(math.Cos(2*math.Pi*t))
		if i == 0 {
			p.MoveTo(internal.Snap(pt, x), internal.Snap(pt, y))
		} else {
			p.LineTo(internal.Snap(pt, x), internal.Snap(pt, y))
		}
	}
	return &p
}

// MeterOptions configure a Meter.
type MeterOptions struct {
	// Min and Max bound the value; both zero mean 0 to 1.
	Min, Max float64
	// Levels are the High and Critical thresholds; Critical below Warning flags low values.
	Levels    *ui.MeterLevels
	Label     string
	ShowValue bool
	// Width is the bar width; zero fills.
	Width float32
}

// meterLevel is 0 normal, 1 warning, 2 critical.
func meterLevel(v float64, l *ui.MeterLevels) int {
	if l == nil || l.Warning == l.Critical {
		return 0
	}
	low := l.Critical < l.Warning
	switch {
	case !low && v >= l.Critical, low && v <= l.Critical:
		return 2
	case !low && v >= l.Warning, low && v <= l.Warning:
		return 1
	}
	return 0
}

// Meter shows a level within a range with threshold ticks and a level word, not color alone.
func Meter(c *ui.Context, value float64, opts MeterOptions) *ui.Element {
	lo, hi := opts.Min, opts.Max
	if lo == 0 && hi == 0 {
		hi = 1
	}
	if !finite(lo) || !finite(hi) || !(hi > lo) {
		panic(fmt.Sprintf("mujica: Meter range %v..%v is empty", lo, hi))
	}
	if math.IsNaN(value) || value < lo || value > hi {
		panic(fmt.Sprintf("mujica: Meter value %v outside %v..%v", value, lo, hi))
	}
	if l := opts.Levels; l != nil && !(l.Warning >= lo && l.Warning <= hi && l.Critical >= lo && l.Critical <= hi) {
		panic(fmt.Sprintf("mujica: Meter levels %+v outside %v..%v", *l, lo, hi))
	}
	k := Tokens(c)
	u := Density(c).Unit()
	level := meterLevel(value, opts.Levels)
	tone := []ui.Color{k.Success, k.Warning, k.Danger}[level]
	var word string
	if opts.Levels != nil {
		low := opts.Levels.Critical < opts.Levels.Warning
		word = msg(c, "meter.normal")
		switch {
		case level == 2:
			word = msg(c, "meter.critical")
		case level == 1 && low:
			word = msg(c, "meter.low_warning")
		case level == 1:
			word = msg(c, "meter.warning")
		}
	}
	frac := func(v float64) float32 { return float32((v - lo) / (hi - lo)) }
	col := ui.Column(c).Gap(u * 1.5)
	if !(opts.Width >= 0) || math.IsInf(float64(opts.Width), 1) {
		panic(fmt.Sprintf("mujica: Meter Width %v invalid", opts.Width))
	}
	if opts.Width > 0 {
		col.Width(opts.Width)
	} else {
		col.FillWidth()
	}
	col.Children(func() {
		caption := word
		if opts.ShowValue {
			caption = strconv.FormatFloat(value, 'f', -1, 64) + " / " + strconv.FormatFloat(hi, 'f', -1, 64)
			if word != "" {
				caption += " · " + word
			}
		}
		if opts.Label != "" || caption != "" {
			ui.Row(c).Gap(u * 2).Children(func() {
				ui.Text(c, opts.Label).TextColor(k.Text).Grow(1).Shrink(1)
				ui.Text(c, caption).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize)).FontFeatures("tnum").SingleLine()
			})
		}
		name := opts.Label
		if word != "" {
			name += " (" + word + ")"
		}
		bar := ui.Meter(c, value, lo, hi, opts.Levels).Label(name).FillWidth().Height(u * 2).
			Radius(theme.ControlRadius).Background(k.Border)
		bar.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * frac(value), H: r.H}, tone, theme.ControlRadius)
		})
		if l := opts.Levels; l != nil {
			// Ticks at the thresholds, over the fill.
			bar.DrawOver(func(p *ui.Painter, r ui.Rect) {
				for _, v := range []float64{l.Warning, l.Critical} {
					x := float32(math.Round(float64(r.X + r.W*frac(v))))
					p.Fill(ui.Rect{X: x - 1, Y: r.Y, W: 2, H: r.H}, k.Text.Alpha(0.6), 0)
				}
			})
		}
	})
	return col
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
