package finance

import (
	"cmp"
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/datetime"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

// FillRow is one of the account's fills; ID is stable.
type FillRow struct {
	ID       string
	Time     time.Time
	Symbol   string
	Side     OrderSide
	Price    float64
	Quantity float64
	Fee      float64
}

// TradeHistoryTableOptions configure a TradeHistoryTable.
type TradeHistoryTableOptions struct {
	// Decimals of prices and quantities; zero means 2 and 0.
	Decimals, QuantityDecimals int
	// Version must change when rows change in place.
	Version  int
	Status   data.DataStatus
	Error    string
	Disabled bool
}

// TradeHistoryTableState is the date range and table state.
type TradeHistoryTableState struct {
	// Range keeps fills on its days; zero keeps all.
	Range datetime.DateRange
	Table data.DataTableState[string]
	cache tradingRows[FillRow]
}

// TradeHistoryResult is a built table and its export request.
type TradeHistoryResult struct {
	Element  ui.Element
	Table    data.DataTableView
	exported []FillRow
	export   bool
}

// Exported returns the range's fills when export was asked.
func (r TradeHistoryResult) Exported() ([]FillRow, bool) { return r.exported, r.export }

// TradeHistoryTable lists fills by date range, with export.
func TradeHistoryTable(c *ui.Context, s *TradeHistoryTableState, rows []FillRow, opts TradeHistoryTableOptions) TradeHistoryResult {
	dec := core.PriceDecimals("TradeHistoryTable", opts.Decimals)
	if opts.QuantityDecimals < 0 {
		panic("mujica: TradeHistoryTable with negative QuantityDecimals")
	}
	rng := s.Range
	shown := s.cache.view(rows, opts.Version, rng, func(f FillRow) bool {
		if rng.IsZero() {
			return true
		}
		d := core.DateOf(f.Time)
		return d.Compare(rng.Start) >= 0 && d.Compare(rng.End) <= 0
	})
	num := func(id, key string, dp int, v func(FillRow) float64) data.DataColumn[FillRow] {
		return data.DataColumn[FillRow]{ID: id, Title: core.Msg(c, key), Align: ui.End, Sortable: true,
			Compare: func(a, b FillRow) int { return cmp.Compare(v(a), v(b)) },
			Text:    func(f FillRow) string { return data.StatisticFormat(v(f), dp) }}
	}
	cols := []data.DataColumn[FillRow]{
		{ID: "time", Title: core.Msg(c, "trade_history_table.time"), Sortable: true, Width: 180,
			Compare: func(a, b FillRow) int { return a.Time.Compare(b.Time) },
			Text:    func(f FillRow) string { return f.Time.Format("2006-01-02 15:04:05") }},
		{ID: "symbol", Title: core.Msg(c, "position_table.symbol"), Sortable: true,
			Compare: func(a, b FillRow) int { return cmp.Compare(a.Symbol, b.Symbol) },
			Text:    func(f FillRow) string { return f.Symbol }},
		{ID: "side", Title: core.Msg(c, "order_table.side"), Text: func(f FillRow) string { return core.Msg(c, fmt.Sprintf("order_entry.side%d", f.Side)) },
			Cell: func(f FillRow) { tradingSide(c, f.Side) }},
		num("price", "order_book.price", dec, func(f FillRow) float64 { return f.Price }),
		num("qty", "order_entry.quantity", opts.QuantityDecimals, func(f FillRow) float64 { return f.Quantity }),
		num("fee", "trade_history_table.fee", dec, func(f FillRow) float64 { return f.Fee }),
	}
	d := core.Density(c)
	var r TradeHistoryResult
	r.Element = ui.Column(c).Gap(d.Unit() * 2).Disabled(opts.Disabled).Grow(1).FillWidth().MinHeight(0)
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			if datetime.DateRangePicker(c, &s.Range, datetime.DateRangePickerOptions{Label: core.Msg(c, "trade_history_table.range"), Disabled: disabled}).Changed() {
				slog.Debug("mujica: trade_history_table range", "days", s.Range.Days())
			}
			if !s.Range.IsZero() && input.Button(c, core.Msg(c, "trade_history_table.all"), input.ButtonOptions{Variant: input.Ghost, Disabled: disabled}).Clicked() && !disabled {
				s.Range = datetime.DateRange{}
			}
			ui.Box(c).Grow(1)
			if input.Button(c, core.Msg(c, "trade_history_table.export"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("upload"), Disabled: disabled || len(shown) == 0}).Clicked() && !disabled {
				r.exported, r.export = shown, true
				slog.Debug("mujica: trade_history_table export", "rows", len(shown))
			}
		})
		r.Table = data.DataTable(c, &s.Table, shown, data.DataTableOptions[FillRow, string]{Columns: cols, Key: func(f FillRow) string { return f.ID },
			Status: opts.Status, Error: opts.Error, Empty: core.Msg(c, "trade_history_table.empty")})
		r.Table.Element.Grow(1)
	})
	return r
}
