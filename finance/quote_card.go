package finance

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Quote is an instrument's price and change since the close.
type Quote struct {
	Symbol, Name  string
	Price         float64
	Change        float64
	ChangePercent float64
}

// QuoteCardOptions configure a QuoteCard.
type QuoteCardOptions struct {
	// Decimals of the price and change; zero means 2.
	Decimals int
	// Chart builds the mini chart slot; nil leaves it out.
	Chart func()
}

// QuoteCard shows a symbol, name, price, change and mini chart.
func QuoteCard(c *ui.Context, q Quote, opts QuoteCardOptions) ui.Element {
	if q.Symbol == "" {
		panic("mujica: QuoteCard needs a Symbol")
	}
	dec := core.PriceDecimals("QuoteCard", opts.Decimals)
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	e := ui.Column(c).Gap(u*2).Padding(u*4).MinWidth(200).Radius(theme.CardRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.Border).Role(ui.RoleGroup).Label(q.Symbol)
	e.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.End).Children(func() {
			ui.Text(c, q.Symbol).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600).SingleLine()
			if q.Name != "" {
				ui.Text(c, q.Name).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine().
					Ellipsis("…").MinWidth(0).Shrink(1)
			}
		})
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Wrap().Children(func() {
			PriceText(c, q.Price, PriceTextOptions{Decimals: dec, Size: theme.H2Size})
			PriceChangeBadge(c, q.Change, q.ChangePercent, PriceChangeBadgeOptions{Decimals: dec})
		})
		if opts.Chart != nil {
			ui.Box(c).Height(40).FillWidth().Children(opts.Chart)
		}
	})
	return e
}
