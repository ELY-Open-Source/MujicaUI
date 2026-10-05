package mujica

import (
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// PnL is profit and loss: realized, unrealized and today's.
type PnL struct {
	Realized, Unrealized, Today float64
}

// PnLDisplayOptions configure a PnLDisplay.
type PnLDisplayOptions struct {
	// Decimals of the amounts; zero means 2.
	Decimals int
	// Currency prefixes each amount, as "$".
	Currency string
}

// PnLDisplay shows three profits signed; *hidden masks them.
func PnLDisplay(c *ui.Context, pnl PnL, hidden *bool, opts PnLDisplayOptions) ChoiceResult {
	dec := priceDecimals("PnLDisplay", opts.Decimals)
	k := Tokens(c)
	u := Density(c).Unit()
	var r ChoiceResult
	r.Element = ui.Row(c).Gap(u * 6).AlignItems(ui.Start).Wrap().Role(ui.RoleGroup).Label(msg(c, "pnl_display.label"))
	r.Element.Children(func() {
		for _, item := range []struct {
			key string
			v   float64
		}{{"pnl_display.realized", pnl.Realized}, {"pnl_display.unrealized", pnl.Unrealized}, {"pnl_display.today", pnl.Today}} {
			title := msg(c, item.key)
			sign := trendSign(item.v)
			text := statisticSigned(item.v, dec)
			if opts.Currency != "" {
				text = pnlDisplaySign(sign) + opts.Currency + statisticFormat(math.Abs(item.v), dec)
			}
			col := trendColor(c, sign)
			spoken := title + ", " + trendWord(c, sign) + " " + text
			if *hidden {
				text, col, spoken = "••••", k.TextMuted, title+", "+msg(c, "pnl_display.masked")
			}
			ui.Column(c).Gap(u).Label(spoken).Role(ui.RoleText).Children(func() {
				ui.Text(c, title).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Role(ui.RoleNone)
				ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(col).Role(ui.RoleNone).Children(func() {
					if !*hidden {
						ui.Icon(c, trendIcon(sign)).Size(16, 16)
					}
					ui.Text(c, text).FontSize(FontSize(c, theme.H3Size)).FontFeatures("tnum").FontWeight(600).SingleLine()
				})
			})
		}
		name := msg(c, "pnl_display.hide")
		icon := icons.Must("eye-off")
		if *hidden {
			name, icon = msg(c, "pnl_display.show"), icons.Must("eye")
		}
		b := iconAction(c, icon, name)
		if b.Clicked() && !b.IsDisabled() {
			*hidden = !*hidden
			r.changed = true
			slog.Debug("mujica: pnl_display mask", "hidden", *hidden)
		}
	})
	return r
}

func pnlDisplaySign(sign int) string {
	switch {
	case sign > 0:
		return "+"
	case sign < 0:
		return "−"
	}
	return ""
}
