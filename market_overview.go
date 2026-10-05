package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// MarketBreadth counts the issues that rose, fell and stayed flat.
type MarketBreadth struct{ Advancing, Declining, Unchanged int }

// MarketOverviewOptions configure a MarketOverview.
type MarketOverviewOptions struct {
	Breadth  MarketBreadth
	AsOf     time.Time
	Decimals int
	Disabled bool
}

type MarketOverviewView struct {
	Element *ui.Element
	Chosen  int
}

func MarketOverview(c *ui.Context, indices []IndexQuote, opts MarketOverviewOptions) MarketOverviewView {
	b := opts.Breadth
	if b.Advancing < 0 || b.Declining < 0 || b.Unchanged < 0 {
		panic("mujica: MarketOverview breadth counts must not be negative")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	v := MarketOverviewView{Chosen: -1}
	v.Element = ui.Column(c).Gap(u * 3).Disabled(opts.Disabled)
	v.Element.Children(func() {
		ui.Row(c).Wrap().Gap(u * 3).Children(func() {
			for i, q := range indices {
				ui.Column(c).Key(q.Name).Children(func() {
					if IndexCard(c, q, IndexCardOptions{Decimals: opts.Decimals, Clickable: true, Disabled: opts.Disabled}).Changed() {
						v.Chosen = i
					}
				})
			}
		})
		if total := b.Advancing + b.Declining + b.Unchanged; total > 0 {
			marketBreadthBar(c, b, total)
		}
		if !opts.AsOf.IsZero() {
			ui.Text(c, fmt.Sprintf(msg(c, "market_overview.as_of"), opts.AsOf.Format("2006-01-02 15:04"))).
				FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		}
	})
	if v.Chosen >= 0 {
		slog.Debug("mujica: market_overview choose", "index", indices[v.Chosen].Name)
	}
	return v
}

func marketBreadthBar(c *ui.Context, b MarketBreadth, total int) {
	k := Tokens(c)
	u := Density(c).Unit()
	up, down := trendColor(c, 1), trendColor(c, -1)
	label := fmt.Sprintf(msg(c, "market_overview.breadth"), b.Advancing, b.Declining, b.Unchanged)
	ui.Column(c).Gap(u).MaxWidth(640).FillWidth().Role(ui.RoleMeter).Label(label).Children(func() {
		ui.Box(c).FillWidth().Height(10).Draw(func(p *ui.Painter, r ui.Rect) {
			x := r.X
			for _, part := range []struct {
				n   int
				col ui.Color
			}{{b.Advancing, up}, {b.Unchanged, k.Border}, {b.Declining, down}} {
				w := r.W * float32(part.n) / float32(total)
				p.Fill(ui.Rect{X: x, Y: r.Y, W: max(w-1, 0), H: r.H}, part.col, 1)
				x += w
			}
		})
		ui.Row(c).Gap(u * 4).Children(func() {
			for _, part := range []struct {
				n    int
				sign int
				key  string
			}{{b.Advancing, 1, "market_overview.up"}, {b.Unchanged, 0, "market_overview.flat"}, {b.Declining, -1, "market_overview.down"}} {
				ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(trendColor(c, part.sign)).Children(func() {
					ui.Icon(c, trendIcon(part.sign)).Size(12, 12)
					ui.Text(c, fmt.Sprintf(msg(c, part.key), part.n)).FontSize(FontSize(c, theme.CaptionSize)).FontFeatures("tnum")
				})
			}
		})
	})
}
