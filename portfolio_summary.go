package mujica

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// PortfolioAccount is one account's value and profit.
type PortfolioAccount struct {
	ID, Name string
	Total    float64
	// DayPnL is today's profit; TotalReturn the profit since inception.
	DayPnL, TotalReturn float64
	Cost                float64
}

// PortfolioSummaryOptions configure a PortfolioSummary.
type PortfolioSummaryOptions struct {
	// Currency prefixes amounts; empty is "$".
	Currency string
	// Decimals of amounts; zero means 2.
	Decimals int
	Disabled bool
}

// PortfolioSummaryView is a built summary; Changed reports another account chosen.

func PortfolioSummary(c *ui.Context, account *string, accounts []PortfolioAccount, opts PortfolioSummaryOptions) ChoiceResult {
	if account == nil {
		panic("mujica: PortfolioSummary needs an account")
	}
	dec := priceDecimals("PortfolioSummary", opts.Decimals)
	cur := cmp.Or(opts.Currency, "$")
	k := Tokens(c)
	u := Density(c).Unit()
	var sum PortfolioAccount
	found := *account == ""
	for _, a := range accounts {
		if a.ID == "" {
			panic("mujica: PortfolioAccount needs an ID")
		}
		if *account == "" || a.ID == *account {
			found = found || a.ID == *account
			sum.Total += a.Total
			sum.DayPnL += a.DayPnL
			sum.TotalReturn += a.TotalReturn
			sum.Cost += a.Cost
		}
	}
	if !found {
		panic("mujica: PortfolioSummary account " + *account + " is not listed")
	}
	before := *account
	v := ChoiceResult{}
	v.Element = ui.Column(c).Gap(u*3).Padding(u*4).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Radius(theme.CardRadius).Disabled(opts.Disabled)
	v.Element.Children(func() {
		if len(accounts) > 1 {
			items := []RadioItem[string]{{Value: "", Label: msg(c, "portfolio_summary.all")}}
			for _, a := range accounts {
				items = append(items, RadioItem[string]{Value: a.ID, Label: a.Name})
			}
			RadioGroup(c, account, items, RadioGroupOptions{Horizontal: true, Label: msg(c, "portfolio_summary.accounts"), Disabled: opts.Disabled})
		}
		ui.Box(c).FillWidth().Height(9).Draw(func(p *ui.Painter, r ui.Rect) { internal.RuleDiamond(p, r, 6, 5, k.Ornament) })
		ui.Row(c).Wrap().Gap(u * 8).Children(func() {
			portfolioFigure(c, msg(c, "portfolio_summary.total"), func() {
				ui.Text(c, cur+statisticFormat(sum.Total, dec)).Font(theme.SerifFont).FontSize(FontSize(c, theme.H1Size)).FontFeatures("tnum")
			})
			prev := sum.Total - sum.DayPnL
			portfolioFigure(c, msg(c, "portfolio_summary.day"), func() {
				portfolioSigned(c, sum.DayPnL, cur, dec, portfolioPercent(sum.DayPnL, prev))
			})
			portfolioFigure(c, msg(c, "portfolio_summary.return"), func() {
				portfolioSigned(c, sum.TotalReturn, cur, dec, portfolioPercent(sum.TotalReturn, sum.Cost))
			})
		})
	})
	v.changed = *account != before && !v.Element.IsDisabled()
	if v.Changed() {
		slog.Debug("mujica: portfolio_summary account", "id", *account)
	}
	return v
}

func portfolioFigure(c *ui.Context, title string, value func()) {
	ui.Column(c).Gap(2).Role(ui.RoleGroup).Label(title).Children(func() {
		ui.Text(c, title).FontSize(FontSize(c, theme.CaptionSize)).TextColor(Tokens(c).TextMuted)
		value()
	})
}

func portfolioSigned(c *ui.Context, v float64, cur string, dec int, pct float64) {
	sign := trendSign(v)
	text := cur + statisticFormat(math.Abs(v), dec)
	if sign > 0 {
		text = "+" + text
	} else if sign < 0 {
		text = "−" + text
	}
	ui.Row(c).Gap(4).AlignItems(ui.Center).TextColor(trendColor(c, sign)).Label(trendWord(c, sign) + " " + text).Children(func() {
		ui.Icon(c, trendIcon(sign)).Size(16, 16)
		ui.Text(c, text).FontSize(FontSize(c, theme.H3Size)).FontWeight(600).FontFeatures("tnum")
		if !math.IsNaN(pct) {
			ui.Text(c, fmt.Sprintf("(%s%%)", statisticSigned(pct, 2))).FontFeatures("tnum")
		}
	})
}

func portfolioPercent(part, base float64) float64 {
	if base == 0 {
		return math.NaN()
	}
	return part / base * 100
}
