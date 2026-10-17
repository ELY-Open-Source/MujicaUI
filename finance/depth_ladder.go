package finance

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DepthLadderOptions configure a DepthLadder. TickSize is required.
type DepthLadderOptions struct {
	TickSize float64
	// Decimals of prices and sizes; zero means 2 and 0.
	Decimals, SizeDecimals int
	// Rows shown; zero means 15.
	Rows     int
	Disabled bool
}

// DepthLadderState is the ladder's scroll from last price, in ticks.
type DepthLadderState struct {
	Offset int
	cache  depthLadderCache
}

type depthLadderCache struct {
	key        bookKey
	bids, asks map[int64]float64
}

// DepthLadderResult is a built DepthLadder and its clicked cell.
type DepthLadderResult struct {
	Element ui.Element
	price   float64
	side    OrderSide
}

// Clicked returns the clicked cell's price and side.
func (r DepthLadderResult) Clicked() (price float64, side OrderSide, ok bool) {
	return r.price, r.side, r.side != 0
}

// DepthLadder is a price ladder around last; cell clicks order.
func DepthLadder(c *ui.Context, s *DepthLadderState, last float64, book OrderBookData, opts DepthLadderOptions) DepthLadderResult {
	dec := core.PriceDecimals("DepthLadder", opts.Decimals)
	if opts.SizeDecimals < 0 || opts.Rows < 0 {
		panic("mujica: DepthLadder SizeDecimals and Rows cannot be negative")
	}
	scale := math.Pow(10, float64(dec))
	tick := int64(math.Round(opts.TickSize * scale))
	if tick < 1 || math.Abs(opts.TickSize*scale-float64(tick)) > 1e-6 {
		panic(fmt.Sprintf("mujica: DepthLadder TickSize %v is not a whole number of %d-decimal units", opts.TickSize, dec))
	}
	rows := opts.Rows
	if rows == 0 {
		rows = 15
	}
	depthLadderRefresh(s, book, scale, tick)
	center := int64(math.Round(last*scale/float64(tick))) * tick
	k := core.Tokens(c)
	d := core.Density(c)
	var r DepthLadderResult
	r.Element = ui.Column(c).Gap(d.Unit()).Disabled(opts.Disabled).Role(ui.RoleTable).Label(core.Msg(c, "depth_ladder.label"))
	disabled := r.Element.IsDisabled()
	move := func(by int) {
		if !disabled && by != 0 {
			s.Offset += by
			slog.Debug("mujica: depth_ladder scrolled", "offset", s.Offset)
		}
	}
	if r.Element.Shortcut(0, ui.KeyPageUp) {
		move(rows / 2)
	}
	if r.Element.Shortcut(0, ui.KeyPageDown) {
		move(-rows / 2)
	}
	r.Element.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, core.Msg(c, "depth_ladder.bid")).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).Grow(1).Basis(0).TextAlign(ui.End)
			ui.Text(c, core.Msg(c, "depth_ladder.price")).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).Grow(1).Basis(0).TextAlign(ui.Center)
			ui.Text(c, core.Msg(c, "depth_ladder.ask")).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).Grow(1).Basis(0)
		})
		grid := ui.Column(c).FocusGroup(ui.Vertical | ui.Horizontal)
		grid.HandleInput(func(ev ui.InputEvent) bool {
			if ev.Kind != ui.InputScroll || disabled || ev.DY == 0 {
				return false
			}
			move(int(math.Copysign(1, float64(-ev.DY))))
			return true
		})
		grid.Children(func() {
			top := center + int64(s.Offset+rows/2)*tick
			for i := range rows {
				t := top - int64(i)*tick
				depthLadderRow(c, &r, s, t, t == center, scale, dec, opts.SizeDecimals)
			}
		})
		if s.Offset != 0 {
			if input.Button(c, core.Msg(c, "depth_ladder.recenter"), input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("refresh-cw"), Disabled: disabled}).Clicked() {
				s.Offset = 0
			}
		}
	})
	return r
}

// depthLadderRow is a buy cell, the price, a sell cell.
func depthLadderRow(c *ui.Context, r *DepthLadderResult, s *DepthLadderState, t int64, current bool, scale float64, dec, sizeDec int) {
	k := core.Tokens(c)
	d := core.Density(c)
	price := float64(t) / scale
	text := data.StatisticFormat(price, dec)
	ui.Row(c.Key(t)).MinHeight(d.ControlHeight()).AlignItems(ui.Stretch).BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).Children(func() {
		cell := func(side OrderSide, size float64) {
			label := fmt.Sprintf(core.Msg(c, fmt.Sprintf("depth_ladder.order%d", side)), text)
			b := ui.ButtonBase(c).Grow(1).Basis(0).Padding(0, d.Unit()*2).Label(label)
			if side == SideBuy {
				b.Justify(ui.End)
			} else {
				b.Justify(ui.Start)
			}
			internal.Ring(b, k.Focus, 0)
			if b.Hovered() {
				b.Background(core.TrendColor(c, side.sign()).Alpha(0.14))
			}
			b.Children(func() {
				if size > 0 {
					ui.Text(c, data.StatisticFormat(size, sizeDec)).TextColor(core.TrendColor(c, side.sign())).FontFeatures("tnum").FontWeight(600).SingleLine()
				}
			})
			if b.Clicked() && !b.IsDisabled() {
				r.price, r.side = price, side
				slog.Debug("mujica: depth_ladder order", "side", side)
			}
		}
		cell(SideBuy, s.cache.bids[t])
		mid := ui.Row(c).Grow(1).Basis(0).Justify(ui.Center).AlignItems(ui.Center).Gap(d.Unit()).Background(k.Background)
		if current {
			mid.Background(k.Selection).Label(core.Msg(c, "depth_ladder.last") + " " + text)
			mid.DrawOver(func(p *ui.Painter, rect ui.Rect) {
				internal.Frame(p, rect, 2, 4, k.AccentText)
			})
		}
		mid.Children(func() {
			if current {
				ui.Box(c).Size(8, 8).Draw(func(p *ui.Painter, rect ui.Rect) {
					p.FillPath(internal.Diamond(rect.X+4, rect.Y+4, 8), k.AccentText)
				})
			}
			t := ui.Text(c, text).FontFeatures("tnum").SingleLine()
			if current {
				t.FontWeight(700)
			}
		})
		cell(SideSell, s.cache.asks[t])
	})
}

// depthLadderRefresh indexes sizes by tick when the book changed.
func depthLadderRefresh(s *DepthLadderState, book OrderBookData, scale float64, tick int64) {
	key := bookKeyOf(book, float64(tick)/scale)
	if s.cache.bids != nil && s.cache.key == key {
		return
	}
	index := func(ls []BookLevel) map[int64]float64 {
		m := make(map[int64]float64, len(ls))
		for _, l := range ls {
			m[int64(math.Round(l.Price*scale/float64(tick)))*tick] += l.Size
		}
		return m
	}
	s.cache = depthLadderCache{key: key, bids: index(book.Bids), asks: index(book.Asks)}
}
