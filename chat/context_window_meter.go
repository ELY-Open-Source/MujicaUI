package chat

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ContextSegment describes one part of the context window.
type ContextSegment struct {
	Label  string
	Tokens int
}

// ContextWindowMeterOptions configure a ContextWindowMeter.
type ContextWindowMeterOptions struct {
	// Warn sets the warning threshold; zero means 0.8.
	Warn float64
}

// ContextWindowMeter displays segmented usage and threshold warnings.
// Palette colors and shaped legend markers identify segments.
func ContextWindowMeter(c *ui.Context, segments []ContextSegment, limit int, opts ContextWindowMeterOptions) *ui.Element {
	if limit <= 0 || opts.Warn < 0 || opts.Warn > 1 {
		panic("mujica: ContextWindowMeter needs limit > 0 and 0 <= Warn <= 1")
	}
	used := 0
	for _, s := range segments {
		if s.Label == "" || s.Tokens < 0 {
			panic("mujica: ContextSegment needs a Label and non-negative Tokens")
		}
		used += s.Tokens
	}
	warn := opts.Warn
	if warn == 0 {
		warn = 0.8
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	pal := internal.Palette(k)
	share := float64(used) / float64(limit)
	pct := int(math.Round(share * 100))
	summary := fmt.Sprintf(core.Msg(c, "context_window_meter.used"), tokenThousands(used), tokenThousands(limit), pct)
	col := ui.Column(c).Gap(u * 2).MinWidth(200).Role(ui.RoleMeter).Label(core.Msg(c, "context_window_meter.label")).Description(summary)
	col.Children(func() {
		ui.Row(c).Gap(u * 2).Children(func() {
			ui.Text(c, core.Msg(c, "context_window_meter.label")).FontWeight(600).Grow(1).SingleLine()
			ui.Text(c, summary).FontSize(small).FontFeatures("tnum").TextColor(k.TextMuted).SingleLine()
		})
		ui.Box(c).Height(10).Radius(theme.ControlRadius).Background(k.Border).Clip().Draw(func(p *ui.Painter, r ui.Rect) {
			x := r.X
			for i, s := range segments {
				w := r.W * float32(min(float64(s.Tokens)/float64(limit), 1))
				w = min(w, r.X+r.W-x)
				p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, pal[i%len(pal)], 0)
				x += w
				// A surface gap parts neighbouring segments.
				if i < len(segments)-1 && w > 0 {
					p.Fill(ui.Rect{X: internal.Snap(p, x-1), Y: r.Y, W: 1, H: r.H}, k.Surface, 0)
				}
			}
			mark := r.X + r.W*float32(warn)
			p.Fill(ui.Rect{X: internal.Snap(p, mark), Y: r.Y, W: 1, H: r.H}, k.Text, 0)
		})
		ui.Row(c).Gap(u * 4).Wrap().Children(func() {
			for i, s := range segments {
				ui.Row(c).Gap(u).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(10, 10).Draw(func(p *ui.Painter, r ui.Rect) {
						internal.Marker(p, i, r.X+r.W/2, r.Y+r.H/2, 9, pal[i%len(pal)])
					})
					ui.Text(c, s.Label+" "+tokenThousands(s.Tokens)).FontSize(small).FontFeatures("tnum").SingleLine()
				})
			}
			ui.Text(c, core.Msg(c, "context_window_meter.free")+" "+tokenThousands(max(limit-used, 0))).FontSize(small).
				FontFeatures("tnum").TextColor(k.TextMuted).SingleLine()
		})
		if share >= warn {
			color := k.Warning
			if used > limit {
				color = k.Danger
			}
			ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(color).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must("triangle-alert")).Size(14, 14)
				ui.Text(c, fmt.Sprintf(core.Msg(c, "context_window_meter.warn"), pct)).FontSize(small)
			})
		}
	})
	return col
}
