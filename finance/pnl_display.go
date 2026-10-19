package finance

import (
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
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
func PnLDisplay(c *ui.Context, pnl PnL, hidden *bool, opts PnLDisplayOptions) core.ChoiceResult {
	dec := core.PriceDecimals("PnLDisplay", opts.Decimals)
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = ui.Row(c).Gap(u * 6).AlignItems(ui.Start).Wrap().Role(ui.RoleGroup).Label(core.Msg(c, "pnl_display.label"))
	r.Element.Children(func() {
		for _, item := range []struct {
			key string
			v   float64
		}{{"pnl_display.realized", pnl.Realized}, {"pnl_display.unrealized", pnl.Unrealized}, {"pnl_display.today", pnl.Today}} {
			title := core.Msg(c, item.key)
			sign := core.TrendSign(item.v)
			text := data.StatisticSigned(item.v, dec)
			if opts.Currency != "" {
				text = pnlDisplaySign(sign) + opts.Currency + data.StatisticFormat(math.Abs(item.v), dec)
			}
			col := core.TrendColor(c, sign)
			spoken := title + ", " + core.TrendWord(c, sign) + " " + text
			if *hidden {
				text, col, spoken = "••••", k.TextMuted, title+", "+core.Msg(c, "pnl_display.masked")
			}
			ui.Column(c).Gap(u).Label(spoken).Role(ui.RoleText).Children(func() {
				ui.Text(c, title).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Role(ui.RoleNone)
				ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(col).Role(ui.RoleNone).Children(func() {
					if !*hidden {
						ui.Icon(c, core.TrendIcon(sign)).Size(16, 16)
					}
					ui.Text(c, text).FontSize(core.FontSize(c, theme.H3Size)).FontFeatures("tnum").FontWeight(600).SingleLine()
				})
			})
		}
		name := core.Msg(c, "pnl_display.hide")
		icon := icons.Must("eye-off")
		if *hidden {
			name, icon = core.Msg(c, "pnl_display.show"), icons.Must("eye")
		}
		b := core.IconAction(c, icon, name)
		if b.Clicked() && !b.IsDisabled() {
			*hidden = !*hidden
			r = core.NewChoiceResult(r.Element, true)
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
