package feedback

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// AlertOptions configure an Alert.
type AlertOptions struct {
	Severity core.Severity
	// Title goes above the message in bold; empty shows none.
	Title string
	// Action labels a button at the end; empty shows none.
	Action string
}

// Alert shows an inline message; a non-nil open adds a dismiss button that clears it.
func Alert(c *ui.Context, open *bool, message string, opts AlertOptions) core.AlertResult {
	if message == "" {
		panic("mujica: Alert needs a message")
	}
	if open != nil && !*open {
		return core.AlertResult{}
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	tone, icon := core.SeverityLook(k, opts.Severity)
	name := core.Msg(c, "alert."+opts.Severity.String())
	var res core.AlertResult
	box := ui.Row(c).FillWidth().AlignItems(ui.Start).Gap(u*3).Padding(u*3, u*3, u*3, u*4).
		Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Role(ui.RoleGroup).Label(name).Clip()
	// A bar in the severity color along the leading edge.
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: 3, H: r.H}, tone, 0)
	})
	// Warnings and errors are read out when they appear or change.
	said := ui.Local(box, "alert.said", func() string { return "" })
	if news := name + ": " + message; opts.Severity >= core.SeverityWarning && *said != news {
		*said = news
		c.Announce(news)
	}
	box.Children(func() {
		ui.Icon(c, icon).Size(20, 20).TextColor(tone).Margin(1, 0, 0, 0)
		ui.Column(c).Grow(1).Shrink(1).Gap(u).Children(func() {
			if opts.Title != "" {
				ui.Text(c, opts.Title).Bold().TextColor(k.Text)
			}
			ui.Text(c, message).TextColor(k.Text).LineHeight(theme.LineHeight)
			if opts.Action != "" {
				ui.Row(c).Margin(u, 0, 0, 0).Children(func() {
					if input.Button(c, opts.Action, input.ButtonOptions{Variant: input.Outline}).Clicked() {
						slog.Debug("mujica: alert action", "severity", opts.Severity, "action", opts.Action)
						res.Action = true
					}
				})
			}
		})
		if open != nil {
			dismiss := input.Button(c, "", input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("x"), Label: core.Msg(c, "alert.close")})
			dismiss.TextColor(k.TextMuted).Padding(0).MinHeight(theme.MinHitSize).Margin(-u, -u, -u, 0)
			if dismiss.Clicked() {
				slog.Debug("mujica: alert closed", "severity", opts.Severity)
				*open, res.Closed = false, true
			}
		}
	})
	res.Element = box
	return res
}
