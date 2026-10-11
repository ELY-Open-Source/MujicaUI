package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// BannerOptions configure a Banner.
type BannerOptions struct {
	Severity core.Severity
	// Action labels a button after the message; empty shows none.
	Action string
}

// Banner is a full-width notice for the top of a page; a non-nil open adds a close button.
func Banner(c *ui.Context, open *bool, message string, opts BannerOptions) core.AlertResult {
	if message == "" {
		panic("mujica: Banner needs a message")
	}
	k := core.Tokens(c)
	tone, icon := core.SeverityLook(k, opts.Severity)
	if open != nil && !*open {
		return core.AlertResult{}
	}
	u := core.Density(c).Unit()
	name := core.Msg(c, "alert."+opts.Severity.String())
	var res core.AlertResult
	bar := ui.Row(c).FillWidth().MinHeight(core.Density(c).ControlHeight()+2*u).AlignItems(ui.Center).Gap(u*3).Padding(u*2, u*3, u*2, u*4).
		Background(k.Surface.Mix(tone, 0.08)).Role(ui.RoleGroup).Label(name)
	bar.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 2, W: r.W, H: 2}, tone, 0)
	})
	said := ui.Local(bar, "banner.said", func() string { return "" })
	if news := name + ": " + message; opts.Severity >= core.SeverityWarning && *said != news {
		*said = news
		c.Announce(news)
	}
	bar.Children(func() {
		ui.Icon(c, icon).Size(20, 20).TextColor(tone)
		ui.Text(c, message).TextColor(k.Text).Grow(1).Shrink(1).LineHeight(theme.LineHeight)
		if opts.Action != "" && Button(c, opts.Action, ButtonOptions{Variant: Outline}).Clicked() {
			slog.Debug("mujica: banner action", "severity", opts.Severity)
			res.Action = true
		}
		if open != nil && core.IconAction(c, icons.Must("x"), core.Msg(c, "alert.close")).Clicked() {
			slog.Debug("mujica: banner closed", "severity", opts.Severity)
			*open, res.Closed = false, true
		}
	})
	res.Element = bar
	return res
}
