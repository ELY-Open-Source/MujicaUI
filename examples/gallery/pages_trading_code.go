package main

// tradingCode holds each page's copyable snippet, the same as its doc and Example.
var tradingCode = map[string]string{
	"TickerTape": `items := []finance.TickerItem{{Symbol: "MUJC", Price: 412.08, ChangePercent: -1.35}, {Symbol: "AVE", Price: 88.2, ChangePercent: 2.1}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.TickerTape(c, items, finance.TickerTapeOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"QuoteCard": `quote := finance.Quote{Symbol: "MUJC", Name: "Ave Mujica Holdings", Price: 412.08, Change: -5.62, ChangePercent: -1.35}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.QuoteCard(c, quote, finance.QuoteCardOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"PriceText": `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	core.UseTrend(c, core.RedUp)
	finance.PriceText(c, 18234.5, finance.PriceTextOptions{Compare: true, Reference: 18000})
}
ui.Render(view, 400, 300, 1)`,
	"PriceChangeBadge": `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.PriceChangeBadge(c, -5.62, -1.35, finance.PriceChangeBadgeOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"OrderBook": `var state finance.OrderBookState
book := finance.OrderBookData{
	Bids: []finance.BookLevel{{Price: 100.00, Size: 5}, {Price: 99.99, Size: 3}},
	Asks: []finance.BookLevel{{Price: 100.01, Size: 1}, {Price: 100.02, Size: 2}},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if price, side, ok := finance.OrderBook(c, &state, book, finance.OrderBookOptions{Groups: []float64{0.01, 0.05}}).Clicked(); ok {
		fmt.Println("order at", price, side)
	}
}
ui.Render(view, 400, 300, 1)`,
	"Level2Quotes": `data := finance.Level2Data{
	Bids: []finance.L2Level{{Price: 100.00, Size: 500, Orders: 12}},
	Asks: []finance.L2Level{{Price: 100.01, Size: 200, Orders: 3}},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.Level2Quotes(c, data, finance.Level2QuotesOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"TimeAndSales": `var state finance.TimeAndSalesState
trades := []finance.TradePrint{{ID: "1", Time: time.Now(), Price: 100.01, Size: 300, Side: finance.SideBuy}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.TimeAndSales(c, &state, trades, finance.TimeAndSalesOptions{LargeSize: 1000})
}
ui.Render(view, 400, 300, 1)`,
	"DepthLadder": `var state finance.DepthLadderState
book := finance.OrderBookData{Bids: []finance.BookLevel{{Price: 100, Size: 5}}, Asks: []finance.BookLevel{{Price: 100.01, Size: 2}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if price, side, ok := finance.DepthLadder(c, &state, 100, book, finance.DepthLadderOptions{TickSize: 0.01}).Clicked(); ok {
		fmt.Println("place", side, "at", price)
	}
}
ui.Render(view, 400, 300, 1)`,
	"SpreadIndicator": `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.SpreadIndicator(c, 99.95, 100.05, finance.SpreadIndicatorOptions{WarnPercent: 0.5})
}
ui.Render(view, 400, 300, 1)`,
	"BidAskBar": `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.BidAskBar(c, 6234, 3820, finance.BidAskBarOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"OrderEntry": `var state finance.OrderEntryState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if order, ok := finance.OrderEntry(c, &state, finance.OrderEntryOptions{Symbol: "MUJC", Last: 100.2}).Submitted(); ok {
		fmt.Println("place", order.Quantity)
		state.EndSubmit("")
	}
}
ui.Render(view, 400, 300, 1)`,
	"QuickTradeButtons": `quantity := 10.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if side, q, ok := finance.QuickTradeButtons(c, &quantity, finance.QuickTradeButtonsOptions{Symbol: "MUJC", Confirm: true}).Traded(); ok {
		fmt.Println(side, q)
	}
}
ui.Render(view, 400, 300, 1)`,
	"PositionTable": `var state data.DataTableState[string]
rows := []finance.PositionRow{{ID: "p1", Symbol: "MUJC", Quantity: 100, AvgCost: 400, Price: 412.08}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if id, ok := finance.PositionTable(c, &state, rows, finance.PositionTableOptions{}).Acted(); ok {
		fmt.Println("close", id)
	}
}
ui.Render(view, 400, 300, 1)`,
	"OrderTable": `var state finance.OrderTableState
rows := []finance.OrderRow{{ID: "o1", Symbol: "MUJC", Side: finance.SideBuy, Price: 99.5, Quantity: 100, Filled: 40, Status: finance.OrderPartial}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if id, ok := finance.OrderTable(c, &state, rows, finance.OrderTableOptions{}).Acted(); ok {
		fmt.Println("cancel", id)
	}
}
ui.Render(view, 400, 300, 1)`,
	"TradeHistoryTable": `var state finance.TradeHistoryTableState
rows := []finance.FillRow{{ID: "f1", Time: time.Now(), Symbol: "MUJC", Side: finance.SideBuy, Price: 400, Quantity: 10, Fee: 1.2}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if fills, ok := finance.TradeHistoryTable(c, &state, rows, finance.TradeHistoryTableOptions{}).Exported(); ok {
		fmt.Println("export", len(fills))
	}
}
ui.Render(view, 400, 300, 1)`,
	"PnLDisplay": `hidden := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.PnLDisplay(c, finance.PnL{Realized: 1200.5, Unrealized: -340.2, Today: 88}, &hidden, finance.PnLDisplayOptions{Currency: "$"})
}
ui.Render(view, 400, 300, 1)`,
	"LeverageSlider": `leverage := 5.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if finance.LeverageSlider(c, &leverage, finance.LeverageSliderOptions{Warning: 10, Danger: 50}).Changed() {
		fmt.Println("leverage", leverage)
	}
}
ui.Render(view, 400, 300, 1)`,
	"MarginIndicator": `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.MarginIndicator(c, 0.86, finance.MarginIndicatorOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"RiskMeter": `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.RiskMeter(c, 4, finance.RiskMeterOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"OrderConfirmDialog": `open := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.OrderConfirmDialog(c, &open, finance.OrderConfirmDialogOptions{
		Lines:       []finance.OrderSummaryLine{{Label: "Symbol", Value: "MUJC"}, {Label: "Quantity", Value: "100"}},
		Risk:        "Leveraged products can lose more than the deposit.",
		Acknowledge: "I understand the risk",
	})
	if r.Confirmed {
		fmt.Println("send order")
	}
}
ui.Render(view, 400, 300, 1)`,
}
