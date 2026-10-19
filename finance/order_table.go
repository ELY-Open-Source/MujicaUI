package finance

import (
	"cmp"
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/navigation"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// OrderStatus is where an order stands.
type OrderStatus int

const (
	OrderOpen OrderStatus = iota
	OrderPartial
	OrderFilled
	OrderCanceled
	OrderRejected
)

// working reports an order that can still fill.
func (s OrderStatus) working() bool { return s == OrderOpen || s == OrderPartial }

// OrderRow is one order; ID must be stable across updates.
type OrderRow struct {
	ID, Symbol string
	Side       OrderSide
	Type       OrderType
	Price      float64
	Quantity   float64
	Filled     float64
	Status     OrderStatus
}

// OrderFilter picks which orders an OrderTable shows.
type OrderFilter int

const (
	OrdersAll OrderFilter = iota
	OrdersWorking
	OrdersFilled
	OrdersClosed
)

func (f OrderFilter) keeps(s OrderStatus) bool {
	switch f {
	case OrdersWorking:
		return s.working()
	case OrdersFilled:
		return s == OrderFilled
	case OrdersClosed:
		return s == OrderCanceled || s == OrderRejected
	}
	return true
}

// OrderTableOptions configure an OrderTable.
type OrderTableOptions struct {
	// Decimals of prices and quantities; zero means 2 and 0.
	Decimals, QuantityDecimals int
	// Version must change when rows change in place.
	Version  int
	Status   data.DataStatus
	Error    string
	Disabled bool
}

// OrderTableState is the status filter and table state.
type OrderTableState struct {
	Filter OrderFilter
	Table  data.DataTableState[string]
	cache  tradingRows[OrderRow]
}

// OrderTable lists orders with filter, progress and Cancel.
func OrderTable(c *ui.Context, s *OrderTableState, rows []OrderRow, opts OrderTableOptions) TableActionResult {
	dec := core.PriceDecimals("OrderTable", opts.Decimals)
	if opts.QuantityDecimals < 0 {
		panic("mujica: OrderTable with negative QuantityDecimals")
	}
	qd := opts.QuantityDecimals
	shown := s.cache.view(rows, opts.Version, s.Filter, func(o OrderRow) bool {
		if o.Quantity <= 0 || o.Filled < 0 || o.Filled > o.Quantity {
			panic(fmt.Sprintf("mujica: OrderTable order %s fills %v of %v", o.ID, o.Filled, o.Quantity))
		}
		return s.Filter.keeps(o.Status)
	})
	var r TableActionResult
	text := func(id, key string, t func(OrderRow) string) data.DataColumn[OrderRow] {
		return data.DataColumn[OrderRow]{ID: id, Title: core.Msg(c, key), Sortable: true, Text: t,
			Compare: func(a, b OrderRow) int { return cmp.Compare(t(a), t(b)) }}
	}
	cols := []data.DataColumn[OrderRow]{
		text("symbol", "position_table.symbol", func(o OrderRow) string { return o.Symbol }),
		{ID: "side", Title: core.Msg(c, "order_table.side"), Text: func(o OrderRow) string { return core.Msg(c, fmt.Sprintf("order_entry.side%d", o.Side)) },
			Cell: func(o OrderRow) { tradingSide(c, o.Side) }},
		text("type", "order_entry.type", func(o OrderRow) string { return core.Msg(c, fmt.Sprintf("order_entry.type%d", o.Type)) }),
		{ID: "price", Title: core.Msg(c, "order_book.price"), Align: ui.End, Sortable: true,
			Compare: func(a, b OrderRow) int { return cmp.Compare(a.Price, b.Price) },
			Text: func(o OrderRow) string {
				if o.Type == OrderMarket {
					return core.Msg(c, "order_table.market")
				}
				return data.StatisticFormat(o.Price, dec)
			}},
		{ID: "filled", Title: core.Msg(c, "order_table.filled"), Width: 140,
			Text: func(o OrderRow) string {
				return data.StatisticFormat(o.Filled, qd) + " / " + data.StatisticFormat(o.Quantity, qd)
			},
			Cell: func(o OrderRow) { orderTableProgress(c, o, qd) }},
		text("status", "order_table.status", func(o OrderRow) string { return core.Msg(c, fmt.Sprintf("order_table.status%d", o.Status)) }),
		{ID: "action", Title: core.Msg(c, "position_table.action"), Width: 96, Text: func(OrderRow) string { return "" },
			Cell: func(o OrderRow) {
				if o.Status.working() && input.Button(c, core.Msg(c, "order_table.cancel"), input.ButtonOptions{Variant: input.Secondary}).
					Label(core.Msg(c, "order_table.cancel")+" "+o.Symbol).Clicked() {
					r.id = o.ID
				}
			}},
	}
	d := core.Density(c)
	r.Element = ui.Column(c).Gap(d.Unit() * 2).Disabled(opts.Disabled).Grow(1).FillWidth().MinHeight(0)
	r.Element.Children(func() {
		var segs []navigation.Segment
		for f := OrdersAll; f <= OrdersClosed; f++ {
			segs = append(segs, navigation.Segment{Label: core.Msg(c, fmt.Sprintf("order_table.filter%d", f))})
		}
		at := int(s.Filter)
		if navigation.SegmentedControl(c, &at, segs, navigation.SegmentedControlOptions{Label: core.Msg(c, "order_table.filter"), Disabled: opts.Disabled}).Changed() {
			s.Filter = OrderFilter(at)
			slog.Debug("mujica: order_table filter", "filter", s.Filter)
		}
		r.Table = data.DataTable(c, &s.Table, shown, data.DataTableOptions[OrderRow, string]{Columns: cols, Key: func(o OrderRow) string { return o.ID },
			Status: opts.Status, Error: opts.Error, Empty: core.Msg(c, "order_table.empty")})
		r.Table.Element.Grow(1)
	})
	if r.Element.IsDisabled() {
		r.id = ""
	}
	if r.id != "" {
		slog.Debug("mujica: order_table cancel")
	}
	return r
}

// orderTableProgress shows the filled share as bar and count.
func orderTableProgress(c *ui.Context, o OrderRow, qd int) {
	k := core.Tokens(c)
	share := float32(o.Filled / o.Quantity)
	ui.Row(c).Gap(core.Density(c).Unit() * 2).AlignItems(ui.Center).FillWidth().Children(func() {
		ui.Box(c).Size(48, 6).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, k.SurfaceHover, theme.ControlRadius)
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * share, H: r.H}, k.AccentText, theme.ControlRadius)
		})
		ui.Text(c, data.StatisticFormat(o.Filled, qd)+" / "+data.StatisticFormat(o.Quantity, qd)).FontFeatures("tnum").FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
	})
}

// tradingSide names a side with its arrow and trend color.
func tradingSide(c *ui.Context, side OrderSide) {
	ui.Row(c).Gap(2).AlignItems(ui.Center).TextColor(core.TrendColor(c, side.sign())).Children(func() {
		ui.Icon(c, core.TrendIcon(side.sign())).Size(12, 12)
		ui.Text(c, core.Msg(c, fmt.Sprintf("order_entry.side%d", side))).FontWeight(600).SingleLine()
	})
}

// tradingRows is a filtered view, rebuilt only on change.
type tradingRows[R any] struct {
	first   *R
	n, ver  int
	filter  any
	valid   bool
	visible []R
}

func (t *tradingRows[R]) view(rows []R, version int, filter any, keep func(R) bool) []R {
	var first *R
	if len(rows) > 0 {
		first = &rows[0]
	}
	if t.valid && t.first == first && t.n == len(rows) && t.ver == version && t.filter == filter {
		return t.visible
	}
	t.first, t.n, t.ver, t.filter, t.valid = first, len(rows), version, filter, true
	t.visible = make([]R, 0, len(rows))
	for _, r := range rows {
		if keep(r) {
			t.visible = append(t.visible, r)
		}
	}
	return t.visible
}
