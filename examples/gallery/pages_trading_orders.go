package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/finance"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type tradingSpreadState struct {
	ask, warn float64
}

func tradingSpreadDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "169", func() tradingSpreadState { return tradingSpreadState{ask: 100.05, warn: 0.5} })
	if input.Slider(c, &d.ask, input.SliderOptions{Min: 99.5, Max: 101, Step: 0.01, Label: "Ask (bid 99.95)", ShowValue: true}).Changed() {
		log(fmt.Sprintf("ask %.2f", d.ask))
	}
	if input.Slider(c, &d.warn, input.SliderOptions{Min: 0, Max: 2, Step: 0.1, Label: "Warn at percent", ShowValue: true}).Changed() {
		log(fmt.Sprintf("warn %.1f%%", d.warn))
	}
	finance.SpreadIndicator(c, 99.95, d.ask, finance.SpreadIndicatorOptions{WarnPercent: d.warn})
	formStates(c, func() {
		finance.SpreadIndicator(c, 99.95, 100.05, finance.SpreadIndicatorOptions{})
		finance.SpreadIndicator(c, 99, 101, finance.SpreadIndicatorOptions{WarnPercent: 0.5})
		finance.SpreadIndicator(c, 100.1, 100, finance.SpreadIndicatorOptions{})
	})
}

type tradingBarState struct {
	redUp    bool
	bid, ask float64
}

func tradingBarDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "170", func() tradingBarState { return tradingBarState{bid: 6234, ask: 3820} })
	tradingTrend(c, &d.redUp, log)
	for _, v := range []struct {
		name string
		p    *float64
	}{{"Bid volume", &d.bid}, {"Ask volume", &d.ask}} {
		if input.Slider(c, v.p, input.SliderOptions{Min: 0, Max: 10000, Step: 1, Label: v.name, ShowValue: true}).Changed() {
			log(fmt.Sprintf("%s %.0f", v.name, *v.p))
		}
	}
	ui.Box(c).Width(360).Children(func() { finance.BidAskBar(c, d.bid, d.ask, finance.BidAskBarOptions{}) })
}

type tradingEntryState struct {
	redUp, disabled, reject bool
	entry                   finance.OrderEntryState
}

func tradingEntryDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "171", func() tradingEntryState { return tradingEntryState{} })
	checks := *ui.Local(c.Root(), "171.checks", func() chan [2]any { return make(chan [2]any, 8) })
	placed := *ui.Local(c.Root(), "171.placed", func() chan string { return make(chan string, 1) })
	for drained := false; !drained; {
		select {
		case r := <-checks:
			v, msg := r[0].(uint64), r[1].(string)
			log(fmt.Sprintf("buying power check v%d: %q, kept %v", v, msg, d.entry.Resolve(v, msg)))
		case e := <-placed:
			d.entry.EndSubmit(e)
			log(fmt.Sprintf("order placed, error %q", e))
		default:
			drained = true
		}
	}
	tradingTrend(c, &d.redUp, log)
	formToggles(c, []string{"Disabled", "Broker rejects"}, &d.disabled, &d.reject)
	r := finance.OrderEntry(c, &d.entry, finance.OrderEntryOptions{Symbol: "MUJC", Last: 100.2, MaxQuantity: 5000, Disabled: d.disabled,
		Check: func(o finance.OrderDraft, v uint64) {
			log(fmt.Sprintf("check started v%d", v))
			go func() {
				time.Sleep(600 * time.Millisecond)
				msg := ""
				if o.Quantity*100 > 250000 {
					msg = "Not enough buying power."
				}
				checks <- [2]any{v, msg}
				c.Invalidate()
			}()
		}})
	if o, ok := r.Submitted(); ok {
		log(fmt.Sprintf("submitted %+v", o))
		fail := ""
		if d.reject {
			fail = "The broker rejected the order. Your entry is kept."
		}
		go func() {
			time.Sleep(800 * time.Millisecond)
			placed <- fail
			c.Invalidate()
		}()
	}
}

type tradingQuickState struct {
	redUp, confirm, disabled bool
	qty                      float64
}

func tradingQuickDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "172", func() tradingQuickState { return tradingQuickState{confirm: true, qty: 10} })
	tradingTrend(c, &d.redUp, log)
	formToggles(c, []string{"Confirm first", "Disabled"}, &d.confirm, &d.disabled)
	ui.Box(c).Width(400).Children(func() {
		if side, q, ok := finance.QuickTradeButtons(c, &d.qty, finance.QuickTradeButtonsOptions{Symbol: "MUJC", Confirm: d.confirm, Disabled: d.disabled}).Traded(); ok {
			log(fmt.Sprintf("trade %v %v", side, q))
		}
	})
}

type tradingPositionsState struct {
	redUp, disabled bool
	rows            []finance.PositionRow
	table           data.DataTableState[string]
}

func tradingPositionsDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "173", func() tradingPositionsState {
		return tradingPositionsState{rows: []finance.PositionRow{
			{ID: "p1", Symbol: "MUJC", Quantity: 100, AvgCost: 400, Price: 412.08},
			{ID: "p2", Symbol: "AVE", Quantity: 50, AvgCost: 90, Price: 88.2},
			{ID: "p3", Symbol: "OBLV", Quantity: 8, AvgCost: 1180, Price: 1204.5},
		}}
	})
	tradingTrend(c, &d.redUp, log)
	formToggles(c, []string{"Disabled"}, &d.disabled)
	ui.Box(c).Width(1000).Height(260).Children(func() {
		r := finance.PositionTable(c, &d.table, d.rows, finance.PositionTableOptions{Disabled: d.disabled})
		if id, ok := r.Acted(); ok {
			for i, p := range d.rows {
				if p.ID == id {
					d.rows = append(d.rows[:i:i], d.rows[i+1:]...)
					break
				}
			}
			log("close " + id)
		}
		if r.Table.Changed() {
			id, _ := d.table.Selected()
			log("select " + id)
		}
	})
}

type tradingOrdersState struct {
	disabled bool
	rows     []finance.OrderRow
	table    finance.OrderTableState
	version  int
}

func tradingOrdersDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "174", func() tradingOrdersState {
		return tradingOrdersState{rows: []finance.OrderRow{
			{ID: "o1", Symbol: "MUJC", Side: finance.SideBuy, Type: finance.OrderLimit, Price: 99.5, Quantity: 100, Filled: 40, Status: finance.OrderPartial},
			{ID: "o2", Symbol: "AVE", Side: finance.SideSell, Type: finance.OrderMarket, Quantity: 10, Filled: 10, Status: finance.OrderFilled},
			{ID: "o3", Symbol: "OBLV", Side: finance.SideBuy, Type: finance.OrderStop, Price: 12, Quantity: 5, Status: finance.OrderCanceled},
			{ID: "o4", Symbol: "DOLO", Side: finance.SideSell, Type: finance.OrderLimit, Price: 7.25, Quantity: 30, Status: finance.OrderOpen},
		}}
	})
	formToggles(c, []string{"Disabled"}, &d.disabled)
	ui.Box(c).Width(1000).Height(300).Children(func() {
		if id, ok := finance.OrderTable(c, &d.table, d.rows, finance.OrderTableOptions{Disabled: d.disabled, Version: d.version}).Acted(); ok {
			for i := range d.rows {
				if d.rows[i].ID == id {
					d.rows[i].Status = finance.OrderCanceled
				}
			}
			d.version++
			log("cancel " + id)
		}
	})
}

type tradingHistoryState struct {
	disabled bool
	rows     []finance.FillRow
	table    finance.TradeHistoryTableState
}

func tradingHistoryDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "175", func() tradingHistoryState {
		var rows []finance.FillRow
		for i := range 30 {
			side := finance.SideBuy
			if i%3 == 0 {
				side = finance.SideSell
			}
			rows = append(rows, finance.FillRow{ID: fmt.Sprint(i), Time: time.Date(2026, 9, 10+i/2, 9+i%6, 0, 0, 0, time.Local),
				Symbol: []string{"MUJC", "AVE", "OBLV"}[i%3], Side: side, Price: 100 + float64(i), Quantity: float64(i%5+1) * 10, Fee: 0.5})
		}
		return tradingHistoryState{rows: rows}
	})
	formToggles(c, []string{"Disabled"}, &d.disabled)
	ui.Box(c).Width(1000).Height(360).Children(func() {
		if fills, ok := finance.TradeHistoryTable(c, &d.table, d.rows, finance.TradeHistoryTableOptions{Disabled: d.disabled}).Exported(); ok {
			log(fmt.Sprintf("export %d fills", len(fills)))
		}
	})
}

type tradingPnLState struct {
	redUp, hidden bool
}

func tradingPnLDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "176", func() tradingPnLState { return tradingPnLState{} })
	tradingTrend(c, &d.redUp, log)
	if finance.PnLDisplay(c, finance.PnL{Realized: 1200.5, Unrealized: -340.2, Today: 88}, &d.hidden, finance.PnLDisplayOptions{Currency: "$"}).Changed() {
		log(fmt.Sprintf("hidden %v", d.hidden))
	}
}

type tradingLeverageState struct {
	disabled bool
	value    float64
}

func tradingLeverageDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "177", func() tradingLeverageState { return tradingLeverageState{value: 5} })
	formToggles(c, []string{"Disabled"}, &d.disabled)
	ui.Box(c).Width(380).Children(func() {
		if finance.LeverageSlider(c, &d.value, finance.LeverageSliderOptions{Warning: 10, Danger: 50, Disabled: d.disabled}).Changed() {
			log(fmt.Sprintf("leverage applied %v×", d.value))
		}
	})
}

type tradingMarginState struct{ ratio float64 }

func tradingMarginDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "178", func() tradingMarginState { return tradingMarginState{ratio: 0.55} })
	if input.Slider(c, &d.ratio, input.SliderOptions{Min: 0, Max: 1.4, Step: 0.01, Label: "Margin ratio", ShowValue: true}).Changed() {
		log(fmt.Sprintf("ratio %.2f", d.ratio))
	}
	ui.Box(c).Width(360).Children(func() { finance.MarginIndicator(c, d.ratio, finance.MarginIndicatorOptions{}) })
	formStates(c, func() {
		for _, r := range []float64{0.3, 0.85, 1.1} {
			ui.Box(c).Width(260).Children(func() { finance.MarginIndicator(c, r, finance.MarginIndicatorOptions{}) })
		}
	})
}

type tradingRiskState struct{ level float64 }

func tradingRiskDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "179", func() tradingRiskState { return tradingRiskState{level: 3} })
	if input.Slider(c, &d.level, input.SliderOptions{Min: 1, Max: 5, Step: 1, Label: "Level", ShowValue: true}).Changed() {
		log(fmt.Sprintf("level %.0f", d.level))
	}
	ui.Box(c).Width(360).Children(func() { finance.RiskMeter(c, int(d.level), finance.RiskMeterOptions{}) })
}

type tradingConfirmState struct {
	open, acknowledge bool
}

func tradingConfirmDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "180", func() tradingConfirmState { return tradingConfirmState{acknowledge: true} })
	formToggles(c, []string{"Require acknowledgement"}, &d.acknowledge)
	if tradingAction(c, "Review order") {
		d.open = true
		log("open")
	}
	ack := ""
	if d.acknowledge {
		ack = "I understand leveraged products can lose more than the deposit"
	}
	r := finance.OrderConfirmDialog(c, &d.open, finance.OrderConfirmDialogOptions{
		Lines: []finance.OrderSummaryLine{{Label: "Symbol", Value: "MUJC"}, {Label: "Side", Value: "Buy"}, {Label: "Type", Value: "Limit 100.20"},
			{Label: "Quantity", Value: "100"}, {Label: "Estimated amount", Value: "10,020.00"}},
		Risk: "5× leverage: a 20% move against you wipes out the margin.", Acknowledge: ack,
	})
	if r.Confirmed {
		log("confirmed")
	}
	if r.Canceled {
		log("canceled")
	}
}
