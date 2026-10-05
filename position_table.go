package mujica

import (
	"cmp"
	"log/slog"

	"github.com/egoist/mygo/ui"
)

// PositionRow is one holding; ID must be stable across updates.
type PositionRow struct {
	ID, Symbol string
	Quantity   float64
	AvgCost    float64
	Price      float64
}

// MarketValue is the quantity at the last price.
func (p PositionRow) MarketValue() float64 { return p.Quantity * p.Price }

// Unrealized is the floating profit against the average cost.
func (p PositionRow) Unrealized() float64 { return (p.Price - p.AvgCost) * p.Quantity }

// UnrealizedPercent is the floating profit as a percent of cost.
func (p PositionRow) UnrealizedPercent() float64 {
	if p.AvgCost == 0 {
		return 0
	}
	return (p.Price - p.AvgCost) / p.AvgCost * 100
}

// PositionTableOptions configure a PositionTable.
type PositionTableOptions struct {
	// Decimals of prices and amounts; zero means 2.
	Decimals int
	// Version must change when rows change in place.
	Version  int
	Status   DataStatus
	Error    string
	Disabled bool
}

// TableActionResult is a trading table and its acted row.
type TableActionResult struct {
	Element *ui.Element
	Table   DataTableView
	id      string
}

// Acted returns the ID of the row acted on.
func (r TableActionResult) Acted() (string, bool) { return r.id, r.id != "" }

// PositionTable lists holdings with floating profit and Close.
func PositionTable(c *ui.Context, s *DataTableState[string], rows []PositionRow, opts PositionTableOptions) TableActionResult {
	dec := priceDecimals("PositionTable", opts.Decimals)
	var r TableActionResult
	num := func(id, key string, v func(PositionRow) float64) DataColumn[PositionRow] {
		return DataColumn[PositionRow]{ID: id, Title: msg(c, key), Align: ui.End, Sortable: true,
			Compare: func(a, b PositionRow) int { return cmp.Compare(v(a), v(b)) },
			Text:    func(p PositionRow) string { return statisticFormat(v(p), dec) }}
	}
	pnl := num("pnl", "position_table.pnl", PositionRow.Unrealized)
	pnl.Cell = func(p PositionRow) { tradingSigned(c, p.Unrealized(), statisticSigned(p.Unrealized(), dec)) }
	pct := num("pct", "position_table.pnl_percent", PositionRow.UnrealizedPercent)
	pct.Cell = func(p PositionRow) {
		tradingSigned(c, p.UnrealizedPercent(), statisticSigned(p.UnrealizedPercent(), 2)+"%")
	}
	cols := []DataColumn[PositionRow]{
		{ID: "symbol", Title: msg(c, "position_table.symbol"), Sortable: true,
			Compare: func(a, b PositionRow) int { return cmp.Compare(a.Symbol, b.Symbol) },
			Text:    func(p PositionRow) string { return p.Symbol }},
		num("qty", "position_table.quantity", func(p PositionRow) float64 { return p.Quantity }),
		num("cost", "position_table.cost", func(p PositionRow) float64 { return p.AvgCost }),
		num("price", "position_table.price", func(p PositionRow) float64 { return p.Price }),
		num("value", "position_table.value", PositionRow.MarketValue),
		pnl, pct,
		{ID: "action", Title: msg(c, "position_table.action"), Width: 96,
			Text: func(PositionRow) string { return "" },
			Cell: func(p PositionRow) {
				if Button(c, msg(c, "position_table.close"), ButtonOptions{Variant: Secondary}).Label(msg(c, "position_table.close") + " " + p.Symbol).Clicked() {
					r.id = p.ID
				}
			}},
	}
	r.Element = ui.Column(c).Disabled(opts.Disabled).Grow(1).FillWidth().MinHeight(0)
	r.Element.Children(func() {
		r.Table = DataTable(c, s, rows, DataTableOptions[PositionRow, string]{Columns: cols, Key: func(p PositionRow) string { return p.ID },
			Version: opts.Version, Status: opts.Status, Error: opts.Error, Empty: msg(c, "position_table.empty")})
		r.Table.Element.Grow(1)
	})
	if r.Element.IsDisabled() {
		r.id = ""
	}
	if r.id != "" {
		slog.Debug("mujica: position_table close")
	}
	return r
}

// tradingSigned shows a signed amount with its arrow.
func tradingSigned(c *ui.Context, v float64, text string) {
	sign := trendSign(v)
	ui.Row(c).Gap(2).AlignItems(ui.Center).Justify(ui.End).FillWidth().TextColor(trendColor(c, sign)).Children(func() {
		ui.Icon(c, trendIcon(sign)).Size(12, 12)
		ui.Text(c, text).FontFeatures("tnum").FontWeight(600).SingleLine()
	})
}
