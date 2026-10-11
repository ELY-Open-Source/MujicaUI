package input

import (
	"cmp"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// QRCodeOptions configure a QRCode.
type QRCodeOptions struct {
	Level QRLevel
	// Size is the side in DIPs with the quiet zone; zero is 160.
	Size float32
	// Label names the code for assistive technology; empty uses "QR code".
	Label string
}

// qrCodeCache keeps the last encoding so it runs only when the text or level changes.
type qrCodeCache struct {
	text  string
	level QRLevel
	m     QRMatrix
	err   error
	done  bool
}

// QRCode draws text as a QR code, encoded offline; too long a text shows an error.
func QRCode(c *ui.Context, text string, opts QRCodeOptions) *ui.Element {
	_ = opts.Level.String()
	if !(opts.Size >= 0) || math.IsInf(float64(opts.Size), 1) {
		panic("mujica: QRCode Size must be finite and >= 0")
	}
	size := opts.Size
	if size == 0 {
		size = 160
	}
	k := core.Tokens(c)
	e := ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleImage).Label(cmp.Or(opts.Label, core.Msg(c, "qr_code.label")))
	cache := ui.Local(e, "qr_code", func() qrCodeCache { return qrCodeCache{} })
	if !cache.done || cache.text != text || cache.level != opts.Level {
		m, err := EncodeQR([]byte(text), opts.Level)
		*cache = qrCodeCache{text: text, level: opts.Level, m: m, err: err, done: true}
		slog.Debug("mujica: qr_code encoded", "version", m.Version, "level", opts.Level, "err", err)
	}
	if cache.err != nil {
		e.Center().Padding(core.Density(c).Unit()*2).Background(k.Surface).Border(theme.BorderWidth, k.Danger).Radius(theme.CardRadius).
			Description(cache.err.Error())
		e.Children(func() {
			ui.Text(c, core.Msg(c, "qr_code.too_long")).TextColor(k.Danger).FontSize(core.FontSize(c, theme.CaptionSize)).TextAlign(ui.Center)
		})
		return e
	}
	m := cache.m
	// Scanners need dark modules on a light field, in either theme.
	light, dark := theme.Light().Surface, theme.Light().Text
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		n := m.Size() + 8
		s := p.Scale()
		mod := float32(math.Floor(float64(r.W*s)/float64(n))) / s
		if mod <= 0 {
			return
		}
		ox := float32(math.Round(float64((r.X+(r.W-mod*float32(n))/2)*s))) / s
		oy := float32(math.Round(float64((r.Y+(r.H-mod*float32(n))/2)*s))) / s
		p.Fill(r, light, theme.ControlRadius)
		for y, row := range m.Modules {
			for x, on := range row {
				if on {
					p.Fill(ui.Rect{X: ox + float32(x+4)*mod, Y: oy + float32(y+4)*mod, W: mod, H: mod}, dark, 0)
				}
			}
		}
	})
	return e
}
