package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/finance"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func marketsSentimentGaugeDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "196.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	value := 72.0

	level := ui.Local(c.Root(), "196.value", func() float64 { return 72 })
	if input.Slider(c, level, input.SliderOptions{Min: 0, Max: 100, Step: 1, Label: "Sentiment", ShowValue: true}).Changed() {
		log(fmt.Sprint(*level))
	}
	value = *level
	status, empty, err := marketsDataState(c, "196", log)
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		finance.SentimentGauge(c, value, finance.SentimentGaugeOptions{Label: "Fear & Greed", Status: status, Empty: empty, Error: err, Retry: func() { log("retry requested") }})
	})
}

func marketsPortfolioSummaryDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "197.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	account := ui.Local(c.Root(), "197.account", func() string { return "" })
	accounts := *ui.Local(c.Root(), "markets.PortfolioSummary.accounts.pages_markets_4", func() []finance.PortfolioAccount {
		return []finance.PortfolioAccount{{ID: "a", Name: "Brokerage", Total: 120000, DayPnL: 1250, TotalReturn: 20000, Cost: 100000}, {ID: "b", Name: "Retirement", Total: 80000, DayPnL: -300, TotalReturn: 5000, Cost: 75000}}
	})

	assets := ui.Local(c.Root(), "197.assets", func() float64 { return 120000 })
	if input.NumberInput(c, assets, input.NumberInputOptions{Min: 0, Max: 1000000, Step: 100, Precision: 2, Label: "Brokerage assets"}).Changed() {
		log(fmt.Sprint(*assets))
	}
	accounts[0].Total = *assets
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.PortfolioSummary(c, account, accounts, finance.PortfolioSummaryOptions{})
		if r.Changed() {
			log(fmt.Sprint(account))
		}
	})
}

func marketsAssetAllocationChartDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "198.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	items := *ui.Local(c.Root(), "markets.AssetAllocationChart.items.pages_markets_4", func() []finance.AllocationItem {
		return []finance.AllocationItem{{Name: "Stocks", Value: 60, Target: 50}, {Name: "Bonds", Value: 30, Target: 40}, {Name: "Cash", Value: 10, Target: 10}}
	})

	status := ui.Local(c.Root(), "198.status", func() int { return 0 })
	if input.RadioGroup(c, status, []input.RadioItem[int]{{Value: 0, Label: "Ready"}, {Value: 1, Label: "Loading"}, {Value: 2, Label: "Empty"}, {Value: 3, Label: "Error"}}, input.RadioGroupOptions{Horizontal: true}).Changed() {
		log(fmt.Sprint("status ", *status))
	}
	frame := chart.ChartFrame{Label: "AssetAllocationChart"}
	if *status == 1 {
		frame.Status = data.DataLoading
	}
	if *status == 3 {
		frame.Status = data.DataFailed
		frame.Error = "Sample failed"
	}
	if *status == 2 {
		items = nil
	}
	tolerance := ui.Local(c.Root(), "198.tolerance", func() float64 { return 5 })
	if input.Slider(c, tolerance, input.SliderOptions{Min: 1, Max: 15, Step: 1, Label: "Drift tolerance", ShowValue: true}).Changed() {
		log(fmt.Sprint(*tolerance))
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.AssetAllocationChart(c, items, finance.AssetAllocationChartOptions{Frame: frame, Tolerance: *tolerance})
		if r.Changed() || r.Retried() {
			log(fmt.Sprint(r.Point))
		}
	})
}

func marketsPerformanceChartDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "199.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	values := *ui.Local(c.Root(), "markets.PerformanceChart.values.pages_markets_4", func() []float64 { return []float64{100, 110, 95, 120} })
	benchmark := *ui.Local(c.Root(), "markets.PerformanceChart.benchmark.pages_markets_4", func() []float64 { return []float64{100, 105, 108, 112} })
	times := *ui.Local(c.Root(), "markets.PerformanceChart.times.pages_markets_4", func() []time.Time {
		return []time.Time{time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)}
	})

	status := ui.Local(c.Root(), "199.status", func() int { return 0 })
	if input.RadioGroup(c, status, []input.RadioItem[int]{{Value: 0, Label: "Ready"}, {Value: 1, Label: "Loading"}, {Value: 2, Label: "Empty"}, {Value: 3, Label: "Error"}}, input.RadioGroupOptions{Horizontal: true}).Changed() {
		log(fmt.Sprint("status ", *status))
	}
	frame := chart.ChartFrame{Label: "PerformanceChart"}
	if *status == 1 {
		frame.Status = data.DataLoading
	}
	if *status == 3 {
		frame.Status = data.DataFailed
		frame.Error = "Sample failed"
	}
	if *status == 2 {
		values = nil
		benchmark = nil
		times = nil
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.PerformanceChart(c, values, finance.PerformanceChartOptions{Frame: frame, Times: times, Benchmark: benchmark})
		if r.Changed() || r.Retried() {
			log(fmt.Sprint(r.Point))
		}
		if name, ok := r.Toggled(); ok {
			log(fmt.Sprint(name))
		}
	})
}

func marketsCurrencyConverterDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "200.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "200.state", func() finance.CurrencyConverterState {
		return finance.CurrencyConverterState{Amount: 100, From: "USD", To: "JPY"}
	})
	currencies := *ui.Local(c.Root(), "markets.CurrencyConverter.currencies.pages_markets_4", func() []finance.Currency {
		return []finance.Currency{{Code: "USD", Name: "Dollar", Decimals: 2}, {Code: "JPY", Name: "Yen", Decimals: 0}}
	})
	rate := ui.Local(c.Root(), "200.rate", func() float64 { return 149.37 })

	if input.NumberInput(c, rate, input.NumberInputOptions{Min: 0, Max: 1000, Step: 0.0001, Precision: 4, Label: "Caller rate"}).Changed() {
		log(fmt.Sprint("rate ", *rate))
	}
	if input.Button(c, "Load pair rate", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
		*rate = 149.37
		if state.From == "JPY" {
			*rate = 1 / 149.37
		}
		log("pair rate loaded")
	}
	beforeAmount := state.Amount
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.CurrencyConverter(c, state, currencies, finance.CurrencyConverterOptions{Rate: *rate, RateFrom: state.From, RateTo: state.To})
		if r.PairChanged {
			*rate = 0
			log(fmt.Sprint(state.From, state.To))
		}
	})
	if state.Amount != beforeAmount {
		log(fmt.Sprint("amount ", state.Amount))
	}
}
