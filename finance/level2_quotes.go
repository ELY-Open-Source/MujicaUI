package finance

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// L2Level is a level's price, size and order count.
type L2Level struct {
	Price, Size float64
	Orders      int
}

// Level2Data holds the bid and ask levels, best first.
type Level2Data struct {
	Bids, Asks []L2Level
}

// Level2QuotesOptions configure a Level2Quotes.
type Level2QuotesOptions struct {
	// Decimals of prices and sizes; zero means 2 and 0.
	Decimals, SizeDecimals int
	// Levels shown each side; zero means 5.
	Levels int
}

// level2Mark is a level's last values and size move.
type level2Mark struct {
	price, size float64
	dir         int
	at          time.Time
	seen        bool
}

// Level2Quotes lists bid and ask levels; changes flash and mark.
func Level2Quotes(c *ui.Context, data Level2Data, opts Level2QuotesOptions) ui.Element {
	dec := core.PriceDecimals("Level2Quotes", opts.Decimals)
	if opts.SizeDecimals < 0 || opts.Levels < 0 {
		panic("mujica: Level2Quotes SizeDecimals and Levels cannot be negative")
	}
	n := opts.Levels
	if n == 0 {
		n = 5
	}
	u := core.Density(c).Unit()
	e := ui.Row(c).Gap(u * 3).AlignItems(ui.Start).Role(ui.RoleTable).Label(core.Msg(c, "level2_quotes.label"))
	marks := ui.Local(e, "marks", func() map[string]*level2Mark { return map[string]*level2Mark{} })
	e.Children(func() {
		for _, side := range []OrderSide{SideBuy, SideSell} {
			levels := data.Bids
			if side == SideSell {
				levels = data.Asks
			}
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
				for i := range n {
					var l *L2Level
					if i < len(levels) {
						l = &levels[i]
					}
					level2Row(c, *marks, side, i, l, dec, opts.SizeDecimals)
				}
			})
		}
	})
	return e
}

// level2Row is one level with its move mark.
func level2Row(c *ui.Context, marks map[string]*level2Mark, side OrderSide, i int, l *L2Level, dec, sizeDec int) {
	k := core.Tokens(c)
	d := core.Density(c)
	name := fmt.Sprintf(core.Msg(c, fmt.Sprintf("level2_quotes.side%d", side)), i+1)
	key := fmt.Sprint(side, i)
	m := marks[key]
	if m == nil {
		m = &level2Mark{}
		marks[key] = m
	}
	row := ui.Row(c).Gap(d.Unit()*2).AlignItems(ui.Center).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).
		Radius(theme.ControlRadius).BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border)
	if l == nil {
		row.Children(func() {
			ui.Text(c, name).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).Width(28)
			ui.Text(c, "—").TextColor(k.TextMuted)
		})
		return
	}
	if m.seen && (l.Price != m.price || l.Size != m.size) {
		// A new price is a new level, not a move.
		m.dir, m.at = 0, c.Now()
		if l.Price == m.price {
			m.dir = core.TrendSign(l.Size - m.size)
		}
		slog.Debug("mujica: level2_quotes level changed", "level", name)
	}
	m.price, m.size, m.seen = l.Price, l.Size, true
	if dur := core.Motion(c, priceTextFlash); dur > 0 && !m.at.IsZero() {
		if t := c.Now().Sub(m.at); t < dur {
			row.Background(k.Selection.Alpha(1 - float32(float64(t)/float64(dur))))
			c.AnimationFrame()
		}
	}
	col := core.TrendColor(c, side.sign())
	sizeText := data.StatisticFormat(l.Size, sizeDec)
	orders := fmt.Sprintf(core.Msg(c, "level2_quotes.orders"), l.Orders)
	row.Role(ui.RoleRow).Label(name + " " + data.StatisticFormat(l.Price, dec) + " " + sizeText + " " + orders)
	row.Children(func() {
		ui.Text(c, name).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).Width(28)
		ui.Text(c, data.StatisticFormat(l.Price, dec)).TextColor(col).FontFeatures("tnum").FontWeight(600).Grow(1).Basis(0).SingleLine()
		ui.Row(c).Gap(2).AlignItems(ui.Center).Justify(ui.End).Grow(1).Basis(0).Children(func() {
			if m.dir != 0 {
				ui.Icon(c, core.TrendIcon(m.dir)).Size(12, 12).TextColor(core.TrendColor(c, m.dir))
			}
			ui.Text(c, sizeText).FontFeatures("tnum").SingleLine()
		})
		ui.Text(c, fmt.Sprint(l.Orders)).TextColor(k.TextMuted).FontFeatures("tnum").TextAlign(ui.End).Width(36).SingleLine()
	})
}
