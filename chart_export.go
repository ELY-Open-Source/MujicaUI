package mujica

import (
	"bytes"
	"image/png"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"
)

// ChartExportOptions configure a ChartExport.
type ChartExportOptions struct {
	// Width and Height are the image size in DIPs; required.
	Width, Height int
	// Scale is device pixels per DIP; 0 means 2.
	Scale float32
}

// ChartExport returns a PNG rendered in the current theme.
func ChartExport(c *ui.Context, chart func(c *ui.Context), opts ChartExportOptions) ([]byte, error) {
	if chart == nil || opts.Width <= 0 || opts.Height <= 0 || opts.Scale < 0 || math.IsNaN(float64(opts.Scale)) || math.IsInf(float64(opts.Scale), 0) {
		panic("mujica: ChartExport needs a chart, a positive size and Scale ≥ 0")
	}
	scale := opts.Scale
	if scale == 0 {
		scale = 2
	}
	w := current(c)
	settings := Settings{Mode: Light, Density: w.density, Locale: w.locale, Messages: w.messages}
	if w.tokens.Dark {
		settings.Mode = Dark
	}
	tokens := w.tokens
	ready := false
	tester := ui.NewTester(func(c *ui.Context) {
		if !ready {
			return
		}
		if settings.Mode == Dark {
			settings.Dark = &tokens
		} else {
			settings.Light = &tokens
		}
		Use(c, settings)
		ui.Column(c).Fill().Background(tokens.Surface).Padding(Density(c).Unit() * 4).Children(func() { chart(c) })
	}, opts.Width, opts.Height)
	tester.SetScale(scale)
	tester.SetPreferences(c.Preferences())
	ready = true
	tester.Frame()
	img := tester.Image()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	slog.Debug("mujica: chart_export png", "bytes", buf.Len(), "width", opts.Width, "height", opts.Height)
	return buf.Bytes(), nil
}
