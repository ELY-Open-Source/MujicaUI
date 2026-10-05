package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
)

func marketsSymbolBadgeDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "186.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	symbol := *ui.Local(c.Root(), "markets.SymbolBadge.symbol.pages_markets_2", func() mujica.MarketSymbol {
		return mujica.MarketSymbol{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}
	})

	kind := ui.Local(c.Root(), "186.kind", func() string { return "Stock" })
	if mujica.RadioGroup(c, kind, []mujica.RadioItem[string]{{Value: "Stock", Label: "Stock"}, {Value: "ETF", Label: "ETF"}}, mujica.RadioGroupOptions{Horizontal: true}).Changed() {
		log(*kind)
	}
	symbol.Type = *kind
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.SymbolBadge(c, symbol, mujica.SymbolBadgeOptions{Clickable: true})
		if r.Changed() {
			log(fmt.Sprint(symbol.Code))
		}
	})
}

func marketsMarketHeatmapDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "187.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	path := ui.Local(c.Root(), "187.path", func() []string { return []string{} })
	root := *ui.Local(c.Root(), "markets.MarketHeatmap.root.pages_markets_2", func() mujica.HeatNode {
		return mujica.HeatNode{Name: "Market", Children: []mujica.HeatNode{{Name: "Technology", Children: []mujica.HeatNode{{Name: "AAPL", Cap: 3000, Change: 1.2}, {Name: "MSFT", Cap: 2800, Change: -0.6}}}, {Name: "Energy", Cap: 450, Change: -2.1}}}
	})

	status := ui.Local(c.Root(), "187.status", func() int { return 0 })
	if mujica.RadioGroup(c, status, []mujica.RadioItem[int]{{Value: 0, Label: "Ready"}, {Value: 1, Label: "Loading"}, {Value: 2, Label: "Empty"}, {Value: 3, Label: "Error"}}, mujica.RadioGroupOptions{Horizontal: true}).Changed() {
		log(fmt.Sprint("status ", *status))
	}
	frame := mujica.ChartFrame{Label: "MarketHeatmap"}
	if *status == 1 {
		frame.Status = mujica.DataLoading
	}
	if *status == 3 {
		frame.Status = mujica.DataFailed
		frame.Error = "Sample failed"
	}
	if *status == 2 {
		root.Children = nil
		if len(*path) > 0 {
			*path = nil
			log("path reset: empty data")
		}
	}
	scale := ui.Local(c.Root(), "187.scale", func() float64 { return 3 })
	if mujica.Slider(c, scale, mujica.SliderOptions{Min: 1, Max: 10, Step: 1, Label: "Color scale", ShowValue: true}).Changed() {
		log(fmt.Sprint(*scale))
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.MarketHeatmap(c, path, root, mujica.MarketHeatmapOptions{Frame: frame, Scale: *scale})
		if r.Changed() || r.ValueChanged() || r.Retried() {
			log(fmt.Sprint(path, r.Point))
		}
	})
}

func marketsScreenerDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "188.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "188.state", func() mujica.ScreenerState { return mujica.ScreenerState{} })
	fields := *ui.Local(c.Root(), "markets.Screener.fields.pages_markets_2", func() []mujica.ScreenerField { return []mujica.ScreenerField{{ID: "pe", Label: "P/E", Decimals: 1}} })
	rows := *ui.Local(c.Root(), "markets.Screener.rows.pages_markets_2", func() []mujica.ScreenerRow {
		return []mujica.ScreenerRow{{Symbol: "AAPL", Values: map[string]float64{"pe": 34}}, {Symbol: "XOM", Values: map[string]float64{"pe": 13}}}
	})

	status, empty, err := marketsDataState(c, "188", log)
	if empty || status != mujica.DataReady {
		rows = nil
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.Screener(c, state, fields, rows, mujica.ScreenerOptions{Status: status, Error: err})
		if r.Changed || r.Save || r.Table.Changed() || r.Table.Submitted() || r.Table.SortChanged() || r.Table.Retried() {
			log(fmt.Sprint(state.Conditions))
		}
	})
}

func marketsMarketOverviewDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "189.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	indices := *ui.Local(c.Root(), "markets.MarketOverview.indices.pages_markets_2", func() []mujica.IndexQuote {
		return []mujica.IndexQuote{{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}}
	})

	advancing := ui.Local(c.Root(), "189.advancing", func() float64 { return 2100 })
	if mujica.Slider(c, advancing, mujica.SliderOptions{Min: 0, Max: 4000, Step: 100, Label: "Advancing", ShowValue: true}).Changed() {
		log(fmt.Sprint(*advancing))
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.MarketOverview(c, indices, mujica.MarketOverviewOptions{Breadth: mujica.MarketBreadth{Advancing: int(*advancing), Declining: 1400, Unchanged: 300}})
		if r.Chosen >= 0 {
			log(fmt.Sprint(indices[r.Chosen].Name))
		}
	})
}

func marketsIndexCardDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "190.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	quote := *ui.Local(c.Root(), "markets.IndexCard.quote.pages_markets_2", func() mujica.IndexQuote {
		return mujica.IndexQuote{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}
	})

	precision := ui.Local(c.Root(), "markets.IndexCard.precision", func() float64 { return 2 })
	if mujica.Slider(c, precision, mujica.SliderOptions{Min: 1, Max: 8, Step: 1, Label: "Decimals", ShowValue: true}).Changed() {
		log(fmt.Sprint("decimals ", *precision))
	}
	price := ui.Local(c.Root(), "190.price", func() float64 { return 5751.07 })
	if mujica.NumberInput(c, price, mujica.NumberInputOptions{Min: 1, Max: 10000, Step: 0.01, Precision: 2, Label: "Price"}).Changed() {
		log(fmt.Sprint(*price))
	}
	quote.Price = *price
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.IndexCard(c, quote, mujica.IndexCardOptions{Decimals: int(*precision), Clickable: true})
		if r.Changed() {
			log(fmt.Sprint(quote.Name))
		}
	})
}
