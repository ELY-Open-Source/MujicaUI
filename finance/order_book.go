package finance

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// OrderSide is the buy or sell side.
type OrderSide int

const (
	SideBuy OrderSide = iota + 1
	SideSell
)

// sign is +1 for buying and -1 for selling.
func (s OrderSide) sign() int {
	switch s {
	case SideBuy:
		return 1
	case SideSell:
		return -1
	}
	panic(fmt.Sprintf("mujica: unknown OrderSide %d", s))
}

// BookLevel is a price and the size resting there.
type BookLevel struct {
	Price, Size float64
}

// OrderBookData is the book; bump Version on in-place changes.
type OrderBookData struct {
	Bids, Asks []BookLevel
	Version    int
}

// OrderBookOptions configure an OrderBook.
type OrderBookOptions struct {
	// Decimals of prices and sizes; zero means 2 and 0.
	Decimals, SizeDecimals int
	// Depth is the levels shown each side; zero means 8.
	Depth int
	// Groups are the offered aggregation steps; nil shows raw levels.
	Groups   []float64
	Disabled bool
}

// OrderBookState is an OrderBook's chosen grouping and its aggregated view.
type OrderBookState struct {
	// Group is the aggregation step; zero takes Groups[0].
	Group float64
	cache orderBookCache
}

type orderBookCache struct {
	key        bookKey
	valid      bool
	bids, asks []BookLevel
}

// bookKey tells when a book or its grouping changed.
type bookKey struct {
	bid, ask *BookLevel
	nb, na   int
	version  int
	group    float64
}

func bookKeyOf(book OrderBookData, group float64) bookKey {
	first := func(ls []BookLevel) *BookLevel {
		if len(ls) == 0 {
			return nil
		}
		return &ls[0]
	}
	return bookKey{first(book.Bids), first(book.Asks), len(book.Bids), len(book.Asks), book.Version, group}
}

// OrderBookResult is a built OrderBook and its clicked level.
type OrderBookResult struct {
	Element *ui.Element
	price   float64
	side    OrderSide
}

// Clicked returns the clicked level's price and side.
func (r OrderBookResult) Clicked() (price float64, side OrderSide, ok bool) {
	return r.price, r.side, r.side != 0
}

// OrderBook shows clickable asks over bids with depth bars.
func OrderBook(c *ui.Context, s *OrderBookState, book OrderBookData, opts OrderBookOptions) OrderBookResult {
	dec := core.PriceDecimals("OrderBook", opts.Decimals)
	if opts.SizeDecimals < 0 || opts.Depth < 0 {
		panic("mujica: OrderBook SizeDecimals and Depth cannot be negative")
	}
	depth := opts.Depth
	if depth == 0 {
		depth = 8
	}
	if len(opts.Groups) > 0 && s.Group == 0 {
		s.Group = opts.Groups[0]
	}
	if s.Group != 0 && !slices.Contains(opts.Groups, s.Group) {
		panic(fmt.Sprintf("mujica: OrderBook Group %v is not one of Groups", s.Group))
	}
	orderBookRefresh(s, book, dec)
	k := core.Tokens(c)
	d := core.Density(c)
	var r OrderBookResult
	r.Element = ui.Column(c).Gap(d.Unit()).MinWidth(240).Disabled(opts.Disabled).Role(ui.RoleTable).Label(core.Msg(c, "order_book.label"))
	r.Element.Children(func() {
		if len(opts.Groups) > 0 {
			var choices []input.SelectOption[float64]
			for _, g := range opts.Groups {
				choices = append(choices, input.SelectOption[float64]{Value: g, Label: data.StatisticFormat(g, dec)})
			}
			ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				ui.Text(c, core.Msg(c, "order_book.group")).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize))
				if input.Select(c, &s.Group, choices, input.SelectOptions{Label: core.Msg(c, "order_book.group"), Disabled: opts.Disabled}).Changed() {
					slog.Debug("mujica: order_book group", "step", s.Group)
				}
			})
		}
		asks := s.cache.asks[:min(depth, len(s.cache.asks))]
		bids := s.cache.bids[:min(depth, len(s.cache.bids))]
		total := max(orderBookSum(asks), orderBookSum(bids))
		orderBookHeader(c)
		list := ui.Column(c).FocusGroup(ui.Vertical)
		list.Children(func() {
			cum := orderBookSum(asks)
			for i := len(asks) - 1; i >= 0; i-- {
				orderBookRow(c, &r, asks[i], cum, total, SideSell, dec, opts.SizeDecimals)
				cum -= asks[i].Size
			}
			if len(asks) > 0 && len(bids) > 0 {
				ui.Row(c).Padding(d.Unit(), d.Unit()*2).BorderWidth(theme.BorderWidth, 0).BorderColor(k.Border).Children(func() {
					SpreadIndicator(c, bids[0].Price, asks[0].Price, SpreadIndicatorOptions{Decimals: dec})
				})
			}
			cum = 0
			for _, l := range bids {
				cum += l.Size
				orderBookRow(c, &r, l, cum, total, SideBuy, dec, opts.SizeDecimals)
			}
		})
	})
	if r.side != 0 {
		slog.Debug("mujica: order_book clicked", "side", r.side)
	}
	return r
}

func orderBookSum(ls []BookLevel) float64 {
	var t float64
	for _, l := range ls {
		t += l.Size
	}
	return t
}

// orderBookHeader names the columns.
func orderBookHeader(c *ui.Context) {
	k := core.Tokens(c)
	d := core.Density(c)
	ui.Row(c).Padding(0, d.Unit()*2).Children(func() {
		for i, key := range []string{"order_book.price", "order_book.size", "order_book.total"} {
			t := ui.Text(c, core.Msg(c, key)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Grow(1).Basis(0)
			if i > 0 {
				t.TextAlign(ui.End)
			}
		}
	})
}

// orderBookRow is one clickable level with its depth bar.
func orderBookRow(c *ui.Context, r *OrderBookResult, l BookLevel, cum, total float64, side OrderSide, dec, sizeDec int) {
	k := core.Tokens(c)
	d := core.Density(c)
	col := core.TrendColor(c, side.sign())
	price := data.StatisticFormat(l.Price, dec)
	b := ui.ButtonBase(c).Key(fmt.Sprint(side, price)).MinHeight(d.ControlHeight()).
		Padding(0, d.Unit()*2).Radius(theme.ControlRadius).
		Label(fmt.Sprintf(core.Msg(c, "order_book.level"), core.Msg(c, fmt.Sprintf("order_book.side%d", side)), price, data.StatisticFormat(l.Size, sizeDec)))
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	share := float32(0)
	if total > 0 {
		share = float32(cum / total)
	}
	b.Draw(func(p *ui.Painter, rect ui.Rect) {
		w := rect.W * share
		p.Fill(ui.Rect{X: rect.X + rect.W - w, Y: rect.Y + 2, W: w, H: rect.H - 4}, col.Alpha(0.14), 0)
	})
	b.Children(func() {
		ui.Text(c, price).TextColor(col).FontFeatures("tnum").FontWeight(600).Grow(1).Basis(0).SingleLine()
		ui.Text(c, data.StatisticFormat(l.Size, sizeDec)).FontFeatures("tnum").TextAlign(ui.End).Grow(1).Basis(0).SingleLine()
		ui.Text(c, data.StatisticFormat(cum, sizeDec)).FontFeatures("tnum").TextColor(k.TextMuted).TextAlign(ui.End).Grow(1).Basis(0).SingleLine()
	})
	if b.Clicked() && !b.IsDisabled() {
		r.price, r.side = l.Price, side
	}
}

// orderBookRefresh regroups only when book or grouping changed.
func orderBookRefresh(s *OrderBookState, book OrderBookData, dec int) {
	key := bookKeyOf(book, s.Group)
	if s.cache.valid && s.cache.key == key {
		return
	}
	s.cache = orderBookCache{key: key, valid: true,
		bids: orderBookGroup(book.Bids, s.Group, dec, false), asks: orderBookGroup(book.Asks, s.Group, dec, true)}
}

// orderBookGroup merges levels into group steps, best first.
func orderBookGroup(ls []BookLevel, group float64, dec int, up bool) []BookLevel {
	scale := math.Pow(10, float64(dec))
	step := int64(1)
	if group != 0 {
		step = int64(math.Round(group * scale))
		if step < 1 || math.Abs(group*scale-float64(step)) > 1e-6 {
			panic(fmt.Sprintf("mujica: OrderBook group %v is not a whole number of %d-decimal ticks", group, dec))
		}
	}
	sizes := map[int64]float64{}
	for _, l := range ls {
		if l.Price < 0 || l.Size < 0 {
			panic("mujica: OrderBook levels need prices and sizes of zero or more")
		}
		t := int64(math.Round(l.Price * scale))
		// Bids round down, asks up: groups never beat real prices.
		b := t / step * step
		if up && b < t {
			b += step
		}
		sizes[b] += l.Size
	}
	out := make([]BookLevel, 0, len(sizes))
	for t, sz := range sizes {
		out = append(out, BookLevel{Price: float64(t) / scale, Size: sz})
	}
	slices.SortFunc(out, func(a, b BookLevel) int {
		if up {
			return cmp.Compare(a.Price, b.Price)
		}
		return cmp.Compare(b.Price, a.Price)
	})
	return out
}
