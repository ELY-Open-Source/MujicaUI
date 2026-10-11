package finance_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/finance"
)

func ExampleOptionChain() {
	var state finance.OptionChainState
	rows := []finance.OptionStrike{{Strike: 100, Call: finance.OptionQuote{Bid: 3, Ask: 3.2, Last: 3.1}, Put: finance.OptionQuote{Bid: 2, Ask: 2.2, Last: 2.1}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.OptionChain(c, &state, []string{"Oct 17", "Nov 21"}, rows, finance.OptionChainOptions{Spot: 100.4})
		if r.Changed || r.Submitted || r.ExpiryChanged {
			fmt.Println(state.Expiry, state.Strike, state.Side)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleGreeksTable() {
	rows := []finance.GreeksRow{{Label: "AAPL 100 C", Greeks: finance.OptionGreeks{Delta: 0.5, Gamma: 0.03, Theta: -0.1, Vega: 0.2, Rho: 0.05}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		finance.GreeksTable(c, rows, finance.GreeksTableOptions{Notes: true})
	}
	ui.Render(view, 960, 440, 1)
}

func ExamplePayoffDiagram() {
	legs := []finance.OptionLeg{{Kind: finance.LegCall, Strike: 100, Premium: 3, Quantity: 1}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.PayoffDiagram(c, legs, finance.PayoffDiagramOptions{Frame: chart.ChartFrame{Label: "Call payoff"}, Spot: 103})
		if r.Changed() || r.Retried() {
			fmt.Println(r.Point)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleWatchlist() {
	var state finance.WatchlistState
	groups := []finance.WatchGroup{{Name: "Technology", Items: []finance.WatchItem{{Symbol: "AAPL", Name: "Apple", Price: 227.5, Change: 2.7, ChangePercent: 1.2}}}, {Name: "Energy", Items: []finance.WatchItem{{Symbol: "XOM", Price: 117.4}}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.Watchlist(c, &state, groups, finance.WatchlistOptions{})
		if r.Changed || r.Submitted || r.Moved {
			fmt.Println(state.Selected, r.Rows, r.To)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleSymbolSearch() {
	var state finance.SymbolSearchState
	rows := []finance.MarketSymbol{{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.SymbolSearch(c, &state, rows, finance.SymbolSearchOptions{Markets: []string{"NASDAQ", "NYSE"}, Recent: rows})
		if r.Changed || r.Picked || r.Retry {
			fmt.Println(state.Query, state.Market, r.Chosen.Code)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleSymbolBadge() {
	symbol := finance.MarketSymbol{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.SymbolBadge(c, symbol, finance.SymbolBadgeOptions{Clickable: true})
		if r.Changed() {
			fmt.Println(symbol.Code)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleScreener() {
	var state finance.ScreenerState
	fields := []finance.ScreenerField{{ID: "pe", Label: "P/E", Decimals: 1}}
	rows := []finance.ScreenerRow{{Symbol: "AAPL", Values: map[string]float64{"pe": 34}}, {Symbol: "XOM", Values: map[string]float64{"pe": 13}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.Screener(c, &state, fields, rows, finance.ScreenerOptions{})
		if r.Changed || r.Save || r.Table.Changed() || r.Table.Submitted() || r.Table.SortChanged() || r.Table.Retried() {
			fmt.Println(state.Conditions)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleMarketOverview() {
	indices := []finance.IndexQuote{{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.MarketOverview(c, indices, finance.MarketOverviewOptions{Breadth: finance.MarketBreadth{Advancing: 2100, Declining: 1400, Unchanged: 300}})
		if r.Chosen >= 0 {
			fmt.Println(indices[r.Chosen].Name)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleIndexCard() {
	quote := finance.IndexQuote{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.IndexCard(c, quote, finance.IndexCardOptions{Clickable: true})
		if r.Changed() {
			fmt.Println(quote.Name)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleMarketStatus() {
	now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		finance.MarketStatus(c, finance.MarketOpen, finance.MarketStatusOptions{Market: "NYSE", Now: now, Next: now.Add(5 * time.Hour), NextPhase: finance.MarketAfterHours})
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleTradingSessionClock() {
	now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
	markets := []finance.SessionMarket{{Name: "NYSE", Location: time.FixedZone("EDT", -4*3600), Sessions: []finance.TradingSession{{Phase: finance.MarketOpen, Start: 570, End: 960}}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		finance.TradingSessionClock(c, markets, finance.TradingSessionClockOptions{Now: now})
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleEconomicCalendar() {
	var state finance.EconomicCalendarState
	rows := []finance.EconomicEvent{{ID: "jobs", Country: "US", Title: "Payrolls", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Importance: 3, Forecast: "150K", Previous: "142K", Actual: "254K"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.EconomicCalendar(c, &state, rows, finance.EconomicCalendarOptions{})
		if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
			fmt.Println(state.Table.Selected())
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleEarningsCalendar() {
	var state data.DataTableState[string]
	rows := []finance.EarningsEvent{{ID: "a", Symbol: "AAPL", Name: "Apple", Date: core.Date{Year: 2026, Month: 10, Day: 30}, Timing: finance.EarningsAfterClose, EPSEstimate: 1.6, EPSActual: 1.64}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.EarningsCalendar(c, &state, rows, finance.EarningsCalendarOptions{})
		if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
			fmt.Println(state.Selected())
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleNewsFeed() {
	var state finance.NewsFeedState
	rows := []finance.NewsItem{{ID: "a", Title: "Apple unveils new chips", Source: "Wire", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Symbols: []string{"AAPL"}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.NewsFeed(c, &state, rows, finance.NewsFeedOptions{})
		if r.Opened != "" || r.Symbol != "" || r.Retry {
			fmt.Println(r.Opened, r.Symbol)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleSentimentGauge() {
	value := 72.0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		finance.SentimentGauge(c, value, finance.SentimentGaugeOptions{Label: "Fear & Greed"})
	}
	ui.Render(view, 960, 440, 1)
}

func ExamplePortfolioSummary() {
	account := ""
	accounts := []finance.PortfolioAccount{{ID: "a", Name: "Brokerage", Total: 120000, DayPnL: 1250, TotalReturn: 20000, Cost: 100000}, {ID: "b", Name: "Retirement", Total: 80000, DayPnL: -300, TotalReturn: 5000, Cost: 75000}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.PortfolioSummary(c, &account, accounts, finance.PortfolioSummaryOptions{})
		if r.Changed() {
			fmt.Println(account)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleAssetAllocationChart() {
	items := []finance.AllocationItem{{Name: "Stocks", Value: 60, Target: 50}, {Name: "Bonds", Value: 30, Target: 40}, {Name: "Cash", Value: 10, Target: 10}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.AssetAllocationChart(c, items, finance.AssetAllocationChartOptions{Frame: chart.ChartFrame{Label: "Allocation"}})
		if r.Changed() || r.Retried() {
			fmt.Println(r.Point)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExamplePerformanceChart() {
	values := []float64{100, 110, 95, 120}
	benchmark := []float64{100, 105, 108, 112}
	times := []time.Time{time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.PerformanceChart(c, values, finance.PerformanceChartOptions{Frame: chart.ChartFrame{Label: "Performance"}, Times: times, Benchmark: benchmark})
		if r.Changed() || r.Retried() {
			fmt.Println(r.Point)
		}
		if name, ok := r.Toggled(); ok {
			fmt.Println(name)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleCurrencyConverter() {
	state := finance.CurrencyConverterState{Amount: 100, From: "USD", To: "JPY"}
	currencies := []finance.Currency{{Code: "USD", Name: "Dollar", Decimals: 2}, {Code: "JPY", Name: "Yen", Decimals: 0}}
	rate := 149.37
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := finance.CurrencyConverter(c, &state, currencies, finance.CurrencyConverterOptions{Rate: rate})
		if r.PairChanged {
			rate = 0
			fmt.Println(state.From, state.To)
		}
	}
	ui.Render(view, 960, 440, 1)
}
