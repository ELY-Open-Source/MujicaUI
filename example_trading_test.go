package mujica_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleTickerTape() {
	items := []mujica.TickerItem{{Symbol: "MUJC", Price: 412.08, ChangePercent: -1.35}, {Symbol: "AVE", Price: 88.2, ChangePercent: 2.1}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TickerTape(c, items, mujica.TickerTapeOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleQuoteCard() {
	quote := mujica.Quote{Symbol: "MUJC", Name: "Ave Mujica Holdings", Price: 412.08, Change: -5.62, ChangePercent: -1.35}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.QuoteCard(c, quote, mujica.QuoteCardOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePriceText() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.UseTrend(c, mujica.RedUp)
		mujica.PriceText(c, 18234.5, mujica.PriceTextOptions{Compare: true, Reference: 18000})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePriceChangeBadge() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.PriceChangeBadge(c, -5.62, -1.35, mujica.PriceChangeBadgeOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOrderBook() {
	var state mujica.OrderBookState
	book := mujica.OrderBookData{
		Bids: []mujica.BookLevel{{Price: 100.00, Size: 5}, {Price: 99.99, Size: 3}},
		Asks: []mujica.BookLevel{{Price: 100.01, Size: 1}, {Price: 100.02, Size: 2}},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if price, side, ok := mujica.OrderBook(c, &state, book, mujica.OrderBookOptions{Groups: []float64{0.01, 0.05}}).Clicked(); ok {
			fmt.Println("order at", price, side)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLevel2Quotes() {
	data := mujica.Level2Data{
		Bids: []mujica.L2Level{{Price: 100.00, Size: 500, Orders: 12}},
		Asks: []mujica.L2Level{{Price: 100.01, Size: 200, Orders: 3}},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Level2Quotes(c, data, mujica.Level2QuotesOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTimeAndSales() {
	var state mujica.TimeAndSalesState
	trades := []mujica.TradePrint{{ID: "1", Time: time.Now(), Price: 100.01, Size: 300, Side: mujica.SideBuy}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TimeAndSales(c, &state, trades, mujica.TimeAndSalesOptions{LargeSize: 1000})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDepthLadder() {
	var state mujica.DepthLadderState
	book := mujica.OrderBookData{Bids: []mujica.BookLevel{{Price: 100, Size: 5}}, Asks: []mujica.BookLevel{{Price: 100.01, Size: 2}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if price, side, ok := mujica.DepthLadder(c, &state, 100, book, mujica.DepthLadderOptions{TickSize: 0.01}).Clicked(); ok {
			fmt.Println("place", side, "at", price)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSpreadIndicator() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.SpreadIndicator(c, 99.95, 100.05, mujica.SpreadIndicatorOptions{WarnPercent: 0.5})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBidAskBar() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.BidAskBar(c, 6234, 3820, mujica.BidAskBarOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOrderEntry() {
	var state mujica.OrderEntryState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if order, ok := mujica.OrderEntry(c, &state, mujica.OrderEntryOptions{Symbol: "MUJC", Last: 100.2}).Submitted(); ok {
			fmt.Println("place", order.Quantity)
			state.EndSubmit("")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleQuickTradeButtons() {
	quantity := 10.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if side, q, ok := mujica.QuickTradeButtons(c, &quantity, mujica.QuickTradeButtonsOptions{Symbol: "MUJC", Confirm: true}).Traded(); ok {
			fmt.Println(side, q)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePositionTable() {
	var state mujica.DataTableState[string]
	rows := []mujica.PositionRow{{ID: "p1", Symbol: "MUJC", Quantity: 100, AvgCost: 400, Price: 412.08}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if id, ok := mujica.PositionTable(c, &state, rows, mujica.PositionTableOptions{}).Acted(); ok {
			fmt.Println("close", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOrderTable() {
	var state mujica.OrderTableState
	rows := []mujica.OrderRow{{ID: "o1", Symbol: "MUJC", Side: mujica.SideBuy, Price: 99.5, Quantity: 100, Filled: 40, Status: mujica.OrderPartial}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if id, ok := mujica.OrderTable(c, &state, rows, mujica.OrderTableOptions{}).Acted(); ok {
			fmt.Println("cancel", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTradeHistoryTable() {
	var state mujica.TradeHistoryTableState
	rows := []mujica.FillRow{{ID: "f1", Time: time.Now(), Symbol: "MUJC", Side: mujica.SideBuy, Price: 400, Quantity: 10, Fee: 1.2}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if fills, ok := mujica.TradeHistoryTable(c, &state, rows, mujica.TradeHistoryTableOptions{}).Exported(); ok {
			fmt.Println("export", len(fills))
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePnLDisplay() {
	hidden := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.PnLDisplay(c, mujica.PnL{Realized: 1200.5, Unrealized: -340.2, Today: 88}, &hidden, mujica.PnLDisplayOptions{Currency: "$"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLeverageSlider() {
	leverage := 5.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.LeverageSlider(c, &leverage, mujica.LeverageSliderOptions{Warning: 10, Danger: 50}).Changed() {
			fmt.Println("leverage", leverage)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMarginIndicator() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MarginIndicator(c, 0.86, mujica.MarginIndicatorOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRiskMeter() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.RiskMeter(c, 4, mujica.RiskMeterOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOrderConfirmDialog() {
	open := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.OrderConfirmDialog(c, &open, mujica.OrderConfirmDialogOptions{
			Lines:       []mujica.OrderSummaryLine{{Label: "Symbol", Value: "MUJC"}, {Label: "Quantity", Value: "100"}},
			Risk:        "Leveraged products can lose more than the deposit.",
			Acknowledge: "I understand the risk",
		})
		if r.Confirmed {
			fmt.Println("send order")
		}
	}
	ui.Render(view, 400, 300, 1)
}
