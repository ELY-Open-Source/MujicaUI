package mujica

import (
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// LoaderOptions configure the loaders: DotsLoader, BarsLoader, PulseLoader, OrbitLoader and WaveLoader.
type LoaderOptions struct {
	// Size is the mark's height in DIPs; zero means 16.
	Size float32
	// Label is read out, and shown with ShowLabel; empty means "Loading".
	Label     string
	ShowLabel bool
	// Color paints the mark; zero uses AccentText.
	Color ui.Color
}

// loaderDraw paints a loader at phase t in [0, 1); t is 0 and still is set under reduced motion.
type loaderDraw func(p *ui.Painter, r ui.Rect, t float64, col ui.Color, still bool)

// loaderMark builds a loader of aspect w:1 that repaints only while painted and moving.
func loaderMark(c *ui.Context, opts LoaderOptions, aspect float32, period int64, draw loaderDraw) *ui.Element {
	if !(opts.Size >= 0) || math.IsInf(float64(opts.Size), 1) {
		panic("mujica: loader Size must be finite and not negative")
	}
	k := Tokens(c)
	size := opts.Size
	if size == 0 {
		size = 16
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "loading")
	}
	col := opts.Color
	if col.A == 0 {
		col = k.AccentText
	}
	still := c.Preferences().ReduceMotion
	row := ui.Row(c).Gap(Density(c).Unit() * 2).AlignItems(ui.Center).Shrink(0).Role(ui.RoleProgress).Label(label)
	row.Children(func() {
		ui.Box(c).Key(still).Size(size*aspect, size).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			t := 0.0
			if !still {
				t = float64(p.Now().UnixMilli()%period) / float64(period)
				p.AnimationFrame()
			}
			draw(p, r, t, col, still)
		})
		// Reduced motion shows the words, as the mark stands still.
		if opts.ShowLabel || still {
			ui.Text(c, label).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize)).SingleLine()
		}
	})
	return row
}

// loaderWave is a smooth 0..1 rise and fall over phase t.
func loaderWave(t float64) float32 { return float32(0.5 - 0.5*math.Cos(2*math.Pi*t)) }

// DotsLoader shows three dots lighting up in turn.
func DotsLoader(c *ui.Context, opts LoaderOptions) *ui.Element {
	return loaderMark(c, opts, 2.5, 1200, func(p *ui.Painter, r ui.Rect, t float64, col ui.Color, still bool) {
		d := r.H * 0.36
		for i := range 3 {
			a := float32(1)
			if !still {
				a = 0.3 + 0.7*loaderWave(t-float64(i)/6)
			}
			var dot ui.Path
			cx := r.X + r.W*(float32(i)+0.5)/3
			dot.Circle(cx, r.Y+r.H/2, d/2)
			p.FillPath(&dot, col.Alpha(a))
		}
	})
}
