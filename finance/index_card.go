package finance

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

type IndexQuote struct {
	Name          string
	Price         float64
	Change        float64
	ChangePercent float64
	// Series is recent levels; empty hides the chart.
	Series []float64
}

// IndexCardOptions configure an IndexCard.
type IndexCardOptions struct {
	// Decimals of the level and change; zero means 2.
	Decimals int
	// Clickable makes the card a button.
	Clickable bool
	Disabled  bool
}

func IndexCard(c *ui.Context, q IndexQuote, opts IndexCardOptions) core.ChoiceResult {
	if q.Name == "" {
		panic("mujica: IndexCard needs a Name")
	}
	dec := core.PriceDecimals("IndexCard", opts.Decimals)
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var e ui.Element
	ui.Row(c).Shrink(0).Disabled(opts.Disabled).Children(func() {
		if opts.Clickable {
			e = ui.ButtonBase(c).Justify(ui.Start)
			internal.Ring(e, k.Focus, theme.CardRadius)
		} else {
			e = ui.Row(c).Role(ui.RoleGroup)
		}
	})
	e.Column().AlignItems(ui.Stretch).Gap(u).Padding(u*3).Width(200).Radius(theme.CardRadius).Background(k.Surface).
		Label(q.Name).Border(theme.BorderWidth, k.Border)
	if opts.Clickable && e.Hovered() {
		e.Border(theme.BorderWidth, k.ControlBorder).Background(k.SurfaceHover)
	}
	e.Children(func() {
		ui.Text(c, q.Name).Font(theme.SerifFont).TextColor(k.TextMuted).SingleLine()
		PriceText(c, q.Price, PriceTextOptions{Decimals: dec, Size: theme.H3Size})
		ui.Row(c).Children(func() {
			PriceChangeBadge(c, q.Change, q.ChangePercent, PriceChangeBadgeOptions{Decimals: dec})
		})
		if len(q.Series) > 0 {
			chart.Sparkline(c, q.Series, chart.SparklineOptions{Frame: chart.ChartFrame{Label: q.Name}, Width: 120})
		}
	})
	clicked := opts.Clickable && e.Clicked() && !e.IsDisabled()
	if clicked {
		slog.Debug("mujica: index_card click", "name", q.Name)
	}
	return core.NewChoiceResult(e, clicked)
}
