package mujica

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// AudioSpectrumOptions configure an AudioSpectrum.
type AudioSpectrumOptions struct {
	// Height is the chart's height; zero means 120.
	Height float32
	// BandLabel names band i in the readout; nil numbers them.
	BandLabel func(i int) string
	// Status and Error show loading and failure at one size.
	Status DataStatus
	Error  string
	// Label names the chart for assistive technology.
	Label string
}

type spectrumKey struct{}

// AudioSpectrum draws caller bands as bars; it never animates itself.
func AudioSpectrum(c *ui.Context, bands []float32, opts AudioSpectrumOptions) *ui.Element {
	for _, b := range bands {
		if b < 0 || b > 1 || math.IsNaN(float64(b)) {
			panic("mujica: AudioSpectrum band outside 0..1")
		}
	}
	h := opts.Height
	if h < 0 {
		panic("mujica: AudioSpectrum Height must not be negative")
	}
	if h == 0 {
		h = 120
	}
	k := Tokens(c)
	label := opts.Label
	if label == "" {
		label = msg(c, "audio_spectrum.label")
	}
	e := ui.Box(c).Height(h).MinWidth(120).Radius(theme.ControlRadius).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).Focusable().Label(label)
	internal.Ring(e, k.Focus, theme.ControlRadius)
	if mediaChartState(c, e, opts.Status, opts.Error, len(bands) == 0, "audio_spectrum.empty") {
		return e
	}
	cur := ui.Local(e, spectrumKey{}, func() int { return -1 })
	*cur = min(*cur, len(bands)-1)
	name := func(i int) string {
		if opts.BandLabel != nil {
			return opts.BandLabel(i)
		}
		return fmt.Sprintf(msg(c, "audio_spectrum.band"), i+1)
	}
	if moved, ok := chartCursor(e, cur, len(bands)); ok && moved {
		c.Announce(spectrumReadout(name(*cur), bands[*cur]))
	}
	readout := ""
	if *cur >= 0 {
		readout = spectrumReadout(name(*cur), bands[*cur])
		e.Description(readout)
	}
	sel, size := *cur, FontSize(c, theme.CaptionSize)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		pad := float32(8)
		n := float32(len(bands))
		w := (r.W - 2*pad) / n
		base := r.Y + r.H - pad
		for i, b := range bands {
			x := internal.Snap(p, r.X+pad+float32(i)*w)
			bh := internal.Snap(p, (r.H-2*pad-16)*b)
			col := k.AccentText
			if i == sel {
				col = k.Text
			}
			p.Fill(ui.Rect{X: x + 1, Y: base - bh, W: max(1, w-2), H: bh}, col, 0)
		}
		p.Fill(ui.Rect{X: r.X + pad, Y: base, W: r.W - 2*pad, H: 1}, k.Border, 0)
		if sel >= 0 {
			x := r.X + pad + float32(sel)*w
			p.StrokeDashed(ui.Rect{X: x, Y: r.Y + pad, W: w, H: base - r.Y - pad}, k.TextMuted, 0, 1)
			p.Text(r.X+pad, r.Y+4, readout, size, k.Text)
		}
	})
	return e
}

func spectrumReadout(name string, v float32) string {
	return fmt.Sprintf("%s: %d%%", name, int(math.Round(float64(v*100))))
}

// chartCursor moves *cur by arrows, Home and End.
func chartCursor(e *ui.Element, cur *int, n int) (moved, ok bool) {
	if e.IsDisabled() || n == 0 {
		return false, false
	}
	old := *cur
	switch {
	case e.Shortcut(0, ui.KeyRight):
		*cur = min(n-1, *cur+1)
	case e.Shortcut(0, ui.KeyLeft):
		*cur = max(0, *cur-1)
	case e.Shortcut(0, ui.KeyHome):
		*cur = 0
	case e.Shortcut(0, ui.KeyEnd):
		*cur = n - 1
	default:
		return false, false
	}
	return *cur != old, true
}

// mediaChartState fills e with a loading, error or empty state.
func mediaChartState(c *ui.Context, e *ui.Element, status DataStatus, errText string, empty bool, emptyKey string) bool {
	k := Tokens(c)
	switch {
	case status == DataLoading:
		e.Center().Children(func() { Spinner(c, SpinnerOptions{ShowLabel: true}) })
	case status == DataFailed:
		if errText == "" {
			panic("mujica: a failed media chart needs Error")
		}
		e.Center().Description(errText).Children(func() {
			ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(16, 16).TextColor(k.Danger)
				ui.Text(c, errText).TextColor(k.Danger).MaxLines(2)
			})
		})
	case status != DataReady:
		panic(fmt.Sprintf("mujica: media chart Status %d is not Ready, Loading or Failed", status))
	case empty:
		e.Center().Description(msg(c, emptyKey)).Children(func() {
			ui.Text(c, msg(c, emptyKey)).TextColor(k.TextMuted)
		})
	default:
		return false
	}
	return true
}
