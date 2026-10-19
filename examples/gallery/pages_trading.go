package main

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/finance"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

const tradingCategory = "Trading"

// tradingTrend shows and applies the rise-color switch.
func tradingTrend(c *ui.Context, redUp *bool, log func(string)) {
	if input.Checkbox(c, redUp, "Red means up (China)", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprintf("red up %v", *redUp))
	}
	t := core.GreenUp
	if *redUp {
		t = core.RedUp
	}
	core.UseTrend(c, t)
}

// tradingAction is a demo button at its natural width.
func tradingAction(c *ui.Context, label string) bool {
	clicked := false
	ui.Row(c).Children(func() {
		clicked = input.Button(c, label, input.ButtonOptions{Variant: input.Secondary}).Clicked()
	})
	return clicked
}

// tradingBook is a deterministic book around mid, moved by step.
func tradingBook(mid float64, step int) finance.OrderBookData {
	var b finance.OrderBookData
	for i := range 12 {
		w := float64((i*7+step*3)%11 + 1)
		b.Bids = append(b.Bids, finance.BookLevel{Price: math.Round((mid-0.01*float64(i+1))*100) / 100, Size: w * 10})
		b.Asks = append(b.Asks, finance.BookLevel{Price: math.Round((mid+0.01*float64(i))*100) / 100, Size: float64((i*5+step)%9+1) * 10})
	}
	b.Version = step
	return b
}

var tradingPages = []page{
	{ID: "161", Name: "TickerTape", ZhName: "行情滚动条", Category: tradingCategory, Code: tradingCode["TickerTape"], Demo: tradingTickerDemo,
		Keys: "Tab focuses the tape, which pauses it. With reduced motion, Left and Right turn pages."},
	{ID: "162", Name: "QuoteCard", ZhName: "报价卡片", Category: tradingCategory, Code: tradingCode["QuoteCard"], Demo: tradingQuoteDemo,
		Keys: "The card takes no focus; controls in its chart slot do."},
	{ID: "163", Name: "PriceText", ZhName: "价格文本", Category: tradingCategory, Code: tradingCode["PriceText"], Demo: tradingPriceDemo,
		Keys: "Display only."},
	{ID: "164", Name: "PriceChangeBadge", ZhName: "涨跌徽标", Category: tradingCategory, Code: tradingCode["PriceChangeBadge"], Demo: tradingBadgeDemo,
		Keys: "Display only; the name reads the direction and values."},
	{ID: "165", Name: "OrderBook", ZhName: "盘口", Category: tradingCategory, Code: tradingCode["OrderBook"], Demo: tradingBookDemo,
		Keys: "Tab into the levels, Up and Down move, Enter picks a price."},
	{ID: "166", Name: "Level2Quotes", ZhName: "多档行情", Category: tradingCategory, Code: tradingCode["Level2Quotes"], Demo: tradingLevel2Demo,
		Keys: "Display only; each row reads its level, price, size and orders."},
	{ID: "167", Name: "TimeAndSales", ZhName: "逐笔成交", Category: tradingCategory, Code: tradingCode["TimeAndSales"], Demo: tradingTapeDemo,
		Keys: "Tab to Pause and the list; Up and Down move in the list."},
	{ID: "168", Name: "DepthLadder", ZhName: "价格阶梯", Category: tradingCategory, Code: tradingCode["DepthLadder"], Demo: tradingLadderDemo,
		Keys: "Arrows move between cells, Enter orders, Page Up and Page Down scroll, the wheel steps a tick."},
	{ID: "169", Name: "SpreadIndicator", ZhName: "价差指示", Category: tradingCategory, Code: tradingCode["SpreadIndicator"], Demo: tradingSpreadDemo,
		Keys: "Display only."},
	{ID: "170", Name: "BidAskBar", ZhName: "买卖力量条", Category: tradingCategory, Code: tradingCode["BidAskBar"], Demo: tradingBarDemo,
		Keys: "Display only; read as a meter."},
	{ID: "171", Name: "OrderEntry", ZhName: "下单面板", Category: tradingCategory, Code: tradingCode["OrderEntry"], Demo: tradingEntryDemo,
		Keys: "Tab through side, type and fields; Up and Down step numbers; a failed submit focuses the first error."},
	{ID: "172", Name: "QuickTradeButtons", ZhName: "快速交易按钮", Category: tradingCategory, Code: tradingCode["QuickTradeButtons"], Demo: tradingQuickDemo,
		Keys: "Enter or Space trades; with Confirm, Enter in the popup cancels and Escape closes it."},
	{ID: "173", Name: "PositionTable", ZhName: "持仓表", Category: tradingCategory, Code: tradingCode["PositionTable"], Demo: tradingPositionsDemo,
		Keys: "DataTable keys; Tab reaches each Close button."},
	{ID: "174", Name: "OrderTable", ZhName: "委托表", Category: tradingCategory, Code: tradingCode["OrderTable"], Demo: tradingOrdersDemo,
		Keys: "Arrows switch the status filter; DataTable keys; Tab reaches Cancel."},
	{ID: "175", Name: "TradeHistoryTable", ZhName: "成交记录表", Category: tradingCategory, Code: tradingCode["TradeHistoryTable"], Demo: tradingHistoryDemo,
		Keys: "Tab to the range picker, Export and the table."},
	{ID: "176", Name: "PnLDisplay", ZhName: "盈亏显示", Category: tradingCategory, Code: tradingCode["PnLDisplay"], Demo: tradingPnLDemo,
		Keys: "Tab to the eye button; Enter or Space masks the amounts."},
	{ID: "177", Name: "LeverageSlider", ZhName: "杠杆滑块", Category: tradingCategory, Code: tradingCode["LeverageSlider"], Demo: tradingLeverageDemo,
		Keys: "Arrows step levels, Home and End jump; Apply sets the value."},
	{ID: "178", Name: "MarginIndicator", ZhName: "保证金率指示", Category: tradingCategory, Code: tradingCode["MarginIndicator"], Demo: tradingMarginDemo,
		Keys: "Display only; read as a meter with its lines."},
	{ID: "179", Name: "RiskMeter", ZhName: "风险等级", Category: tradingCategory, Code: tradingCode["RiskMeter"], Demo: tradingRiskDemo,
		Keys: "Display only."},
	{ID: "180", Name: "OrderConfirmDialog", ZhName: "下单确认对话框", Category: tradingCategory, Code: tradingCode["OrderConfirmDialog"], Demo: tradingConfirmDemo,
		Keys: "Cancel has the focus; Escape cancels; the backdrop does not close it."},
}

type tradingTickerState struct {
	redUp bool
	items []finance.TickerItem
}

func tradingTickerDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "161", func() tradingTickerState {
		return tradingTickerState{items: []finance.TickerItem{{Symbol: "MUJC", Price: 412.08, ChangePercent: -1.35}, {Symbol: "AVE", Price: 88.2, ChangePercent: 2.1}, {Symbol: "OBLV", Price: 1204.5, ChangePercent: 0.4},
			{Symbol: "DOLO", Price: 7.25, ChangePercent: -0.8}, {Symbol: "MORT", Price: 33.1, ChangePercent: 1.2}, {Symbol: "TIMR", Price: 61.9, ChangePercent: -2.3}}}
	})
	tradingTrend(c, &d.redUp, log)
	if tradingAction(c, "Move prices") {
		for i := range d.items {
			d.items[i].ChangePercent = -d.items[i].ChangePercent
		}
		log("prices moved")
	}
	ui.Box(c).Width(560).Children(func() { finance.TickerTape(c, d.items, finance.TickerTapeOptions{}) })
	formStates(c, func() {
		ui.Box(c).Width(280).Children(func() { finance.TickerTape(c, d.items[:2], finance.TickerTapeOptions{Speed: 15}) })
	})
}

type tradingQuoteState struct {
	redUp bool
	price float64
}

func tradingQuoteDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "162", func() tradingQuoteState { return tradingQuoteState{price: 412.08} })
	tradingTrend(c, &d.redUp, log)
	if input.Slider(c, &d.price, input.SliderOptions{Min: 380, Max: 440, Step: 0.01, Label: "Price", ShowValue: true}).Changed() {
		log(fmt.Sprintf("price %.2f", d.price))
	}
	change := d.price - 417.7
	ui.Row(c).Gap(16).Wrap().Children(func() {
		finance.QuoteCard(c, finance.Quote{Symbol: "MUJC", Name: "Ave Mujica Holdings", Price: d.price, Change: change, ChangePercent: change / 417.7 * 100},
			finance.QuoteCardOptions{Chart: func() {
				k := core.Tokens(c)
				ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
					p.Line(r.X, r.Y+r.H*0.7, r.X+r.W/2, r.Y+r.H*0.3, 1.5, k.Info)
					p.Line(r.X+r.W/2, r.Y+r.H*0.3, r.X+r.W, r.Y+r.H*(0.5-float32(change)/40), 1.5, k.Info)
				})
			}})
		finance.QuoteCard(c, finance.Quote{Symbol: "AVE", Price: 88.2}, finance.QuoteCardOptions{})
	})
}

type tradingPriceState struct {
	redUp, noFlash, grouping bool
	price                    float64
}

func tradingPriceDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "163", func() tradingPriceState { return tradingPriceState{price: 18234.5, grouping: true} })
	tradingTrend(c, &d.redUp, log)
	formToggles(c, []string{"Thousands separators", "No flash"}, &d.grouping, &d.noFlash)
	ui.Row(c).Gap(8).Children(func() {
		if input.Button(c, "Tick up", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			d.price += 12.5
			log(fmt.Sprintf("price %.2f", d.price))
		}
		if input.Button(c, "Tick down", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			d.price -= 12.5
			log(fmt.Sprintf("price %.2f", d.price))
		}
	})
	finance.PriceText(c, d.price, finance.PriceTextOptions{Size: 28, Compare: true, Reference: 18234.5, NoFlash: d.noFlash, NoGrouping: !d.grouping})
	formStates(c, func() {
		finance.PriceText(c, 0.004521, finance.PriceTextOptions{Decimals: 6})
		finance.PriceText(c, 100, finance.PriceTextOptions{Compare: true, Reference: 100})
	})
}

type tradingBadgeState struct {
	redUp, hidePercent bool
	change             float64
}

func tradingBadgeDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "164", func() tradingBadgeState { return tradingBadgeState{change: 1.25} })
	tradingTrend(c, &d.redUp, log)
	formToggles(c, []string{"Hide percent"}, &d.hidePercent)
	if input.Slider(c, &d.change, input.SliderOptions{Min: -5, Max: 5, Step: 0.05, Label: "Change", ShowValue: true}).Changed() {
		log(fmt.Sprintf("change %.2f", d.change))
	}
	finance.PriceChangeBadge(c, d.change, d.change/1.47, finance.PriceChangeBadgeOptions{HidePercent: d.hidePercent})
	formStates(c, func() {
		finance.PriceChangeBadge(c, 2.4, 1.2, finance.PriceChangeBadgeOptions{})
		finance.PriceChangeBadge(c, -2.4, -1.2, finance.PriceChangeBadgeOptions{})
		finance.PriceChangeBadge(c, 0, 0, finance.PriceChangeBadgeOptions{})
	})
}

type tradingBookState struct {
	redUp, disabled bool
	step            int
	book            finance.OrderBookState
}

func tradingBookDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "165", func() tradingBookState { return tradingBookState{} })
	tradingTrend(c, &d.redUp, log)
	formToggles(c, []string{"Disabled"}, &d.disabled)
	if tradingAction(c, "Next update") {
		d.step++
		log(fmt.Sprintf("book update %d", d.step))
	}
	ui.Box(c).Width(320).Children(func() {
		r := finance.OrderBook(c, &d.book, tradingBook(100, d.step), finance.OrderBookOptions{Depth: 6, Groups: []float64{0.01, 0.05, 0.1}, Disabled: d.disabled})
		if price, side, ok := r.Clicked(); ok {
			log(fmt.Sprintf("level %v at %.2f", side, price))
		}
	})
}

type tradingLevel2State struct {
	redUp bool
	step  int
}

func tradingLevel2Demo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "166", func() tradingLevel2State { return tradingLevel2State{} })
	tradingTrend(c, &d.redUp, log)
	if tradingAction(c, "Next update") {
		d.step++
		log(fmt.Sprintf("level 2 update %d", d.step))
	}
	book := tradingBook(100, d.step)
	var data finance.Level2Data
	for i := range 5 {
		data.Bids = append(data.Bids, finance.L2Level{Price: book.Bids[i].Price, Size: book.Bids[i].Size, Orders: int(book.Bids[i].Size / 10)})
		data.Asks = append(data.Asks, finance.L2Level{Price: book.Asks[i].Price, Size: book.Asks[i].Size, Orders: int(book.Asks[i].Size / 10)})
	}
	ui.Box(c).Width(620).Children(func() { finance.Level2Quotes(c, data, finance.Level2QuotesOptions{}) })
}

type tradingTapeState struct {
	redUp  bool
	trades []finance.TradePrint
	tape   finance.TimeAndSalesState
}

func tradingTapeDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "167", func() tradingTapeState { return tradingTapeState{} })
	tradingTrend(c, &d.redUp, log)
	if tradingAction(c, "Add 20 trades") {
		at := time.Date(2026, 10, 5, 9, 30, 0, 0, time.Local)
		for range 20 {
			n := len(d.trades)
			side := finance.SideBuy
			if n%3 == 1 {
				side = finance.SideSell
			}
			d.trades = append(d.trades, finance.TradePrint{ID: fmt.Sprint(n), Time: at.Add(time.Duration(n) * time.Second),
				Price: 100 + float64(n%17)/100, Size: float64((n*37)%1500 + 1), Side: side})
		}
		log(fmt.Sprintf("%d trades", len(d.trades)))
	}
	ui.Box(c).Width(460).Children(func() {
		if finance.TimeAndSales(c, &d.tape, d.trades, finance.TimeAndSalesOptions{LargeSize: 1000}).Changed() {
			log(fmt.Sprintf("paused %v", d.tape.Paused))
		}
	})
}

type tradingLadderState struct {
	redUp, disabled bool
	last            float64
	ladder          finance.DepthLadderState
}

func tradingLadderDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "168", func() tradingLadderState { return tradingLadderState{last: 100} })
	tradingTrend(c, &d.redUp, log)
	formToggles(c, []string{"Disabled"}, &d.disabled)
	if input.Slider(c, &d.last, input.SliderOptions{Min: 99.9, Max: 100.1, Step: 0.01, Label: "Last price", ShowValue: true}).Changed() {
		log(fmt.Sprintf("last %.2f", d.last))
	}
	ui.Box(c).Width(360).Children(func() {
		r := finance.DepthLadder(c, &d.ladder, d.last, tradingBook(100, 0), finance.DepthLadderOptions{TickSize: 0.01, Rows: 11, Disabled: d.disabled})
		if price, side, ok := r.Clicked(); ok {
			log(fmt.Sprintf("order %v at %.2f", side, price))
		}
	})
}
