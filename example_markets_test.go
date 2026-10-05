package mujica_test

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
	"time"
)

func ExampleOptionChain() {
	var state mujica.OptionChainState
	rows := []mujica.OptionStrike{{Strike: 100, Call: mujica.OptionQuote{Bid: 3, Ask: 3.2, Last: 3.1}, Put: mujica.OptionQuote{Bid: 2, Ask: 2.2, Last: 2.1}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.OptionChain(c, &state, []string{"Oct 17", "Nov 21"}, rows, mujica.OptionChainOptions{Spot: 100.4})
		if r.Changed || r.Submitted || r.ExpiryChanged {
			fmt.Println(state.Expiry, state.Strike, state.Side)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleGreeksTable() {
	rows := []mujica.GreeksRow{{Label: "AAPL 100 C", Greeks: mujica.OptionGreeks{Delta: 0.5, Gamma: 0.03, Theta: -0.1, Vega: 0.2, Rho: 0.05}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.GreeksTable(c, rows, mujica.GreeksTableOptions{Notes: true})
	}
	ui.Render(view, 960, 440, 1)
}

func ExamplePayoffDiagram() {
	legs := []mujica.OptionLeg{{Kind: mujica.LegCall, Strike: 100, Premium: 3, Quantity: 1}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.PayoffDiagram(c, legs, mujica.PayoffDiagramOptions{Frame: mujica.ChartFrame{Label: "Call payoff"}, Spot: 103})
		if r.Changed() || r.Retried() {
			fmt.Println(r.Point)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleWatchlist() {
	var state mujica.WatchlistState
	groups := []mujica.WatchGroup{{Name: "Technology", Items: []mujica.WatchItem{{Symbol: "AAPL", Name: "Apple", Price: 227.5, Change: 2.7, ChangePercent: 1.2}}}, {Name: "Energy", Items: []mujica.WatchItem{{Symbol: "XOM", Price: 117.4}}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.Watchlist(c, &state, groups, mujica.WatchlistOptions{})
		if r.Changed || r.Submitted || r.Moved {
			fmt.Println(state.Selected, r.Rows, r.To)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleSymbolSearch() {
	var state mujica.SymbolSearchState
	rows := []mujica.MarketSymbol{{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.SymbolSearch(c, &state, rows, mujica.SymbolSearchOptions{Markets: []string{"NASDAQ", "NYSE"}, Recent: rows})
		if r.Changed || r.Picked || r.Retry {
			fmt.Println(state.Query, state.Market, r.Chosen.Code)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleSymbolBadge() {
	symbol := mujica.MarketSymbol{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.SymbolBadge(c, symbol, mujica.SymbolBadgeOptions{Clickable: true})
		if r.Changed() {
			fmt.Println(symbol.Code)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleMarketHeatmap() {
	var path []string
	root := mujica.HeatNode{Name: "Market", Children: []mujica.HeatNode{{Name: "Technology", Children: []mujica.HeatNode{{Name: "AAPL", Cap: 3000, Change: 1.2}, {Name: "MSFT", Cap: 2800, Change: -0.6}}}, {Name: "Energy", Cap: 450, Change: -2.1}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.MarketHeatmap(c, &path, root, mujica.MarketHeatmapOptions{Frame: mujica.ChartFrame{Label: "Market"}})
		if r.Changed() || r.ValueChanged() || r.Retried() {
			fmt.Println(path, r.Point)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleScreener() {
	var state mujica.ScreenerState
	fields := []mujica.ScreenerField{{ID: "pe", Label: "P/E", Decimals: 1}}
	rows := []mujica.ScreenerRow{{Symbol: "AAPL", Values: map[string]float64{"pe": 34}}, {Symbol: "XOM", Values: map[string]float64{"pe": 13}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.Screener(c, &state, fields, rows, mujica.ScreenerOptions{})
		if r.Changed || r.Save || r.Table.Changed() || r.Table.Submitted() || r.Table.SortChanged() || r.Table.Retried() {
			fmt.Println(state.Conditions)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleMarketOverview() {
	indices := []mujica.IndexQuote{{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.MarketOverview(c, indices, mujica.MarketOverviewOptions{Breadth: mujica.MarketBreadth{Advancing: 2100, Declining: 1400, Unchanged: 300}})
		if r.Chosen >= 0 {
			fmt.Println(indices[r.Chosen].Name)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleIndexCard() {
	quote := mujica.IndexQuote{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.IndexCard(c, quote, mujica.IndexCardOptions{Clickable: true})
		if r.Changed() {
			fmt.Println(quote.Name)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleMarketStatus() {
	now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MarketStatus(c, mujica.MarketOpen, mujica.MarketStatusOptions{Market: "NYSE", Now: now, Next: now.Add(5 * time.Hour), NextPhase: mujica.MarketAfterHours})
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleTradingSessionClock() {
	now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
	markets := []mujica.SessionMarket{{Name: "NYSE", Location: time.FixedZone("EDT", -4*3600), Sessions: []mujica.TradingSession{{Phase: mujica.MarketOpen, Start: 570, End: 960}}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TradingSessionClock(c, markets, mujica.TradingSessionClockOptions{Now: now})
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleEconomicCalendar() {
	var state mujica.EconomicCalendarState
	rows := []mujica.EconomicEvent{{ID: "jobs", Country: "US", Title: "Payrolls", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Importance: 3, Forecast: "150K", Previous: "142K", Actual: "254K"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.EconomicCalendar(c, &state, rows, mujica.EconomicCalendarOptions{})
		if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
			fmt.Println(state.Table.Selected())
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleEarningsCalendar() {
	var state mujica.DataTableState[string]
	rows := []mujica.EarningsEvent{{ID: "a", Symbol: "AAPL", Name: "Apple", Date: mujica.Date{Year: 2026, Month: 10, Day: 30}, Timing: mujica.EarningsAfterClose, EPSEstimate: 1.6, EPSActual: 1.64}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.EarningsCalendar(c, &state, rows, mujica.EarningsCalendarOptions{})
		if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
			fmt.Println(state.Selected())
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleNewsFeed() {
	var state mujica.NewsFeedState
	rows := []mujica.NewsItem{{ID: "a", Title: "Apple unveils new chips", Source: "Wire", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Symbols: []string{"AAPL"}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.NewsFeed(c, &state, rows, mujica.NewsFeedOptions{})
		if r.Opened != "" || r.Symbol != "" || r.Retry {
			fmt.Println(r.Opened, r.Symbol)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleSentimentGauge() {
	value := 72.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.SentimentGauge(c, value, mujica.SentimentGaugeOptions{Label: "Fear & Greed"})
	}
	ui.Render(view, 960, 440, 1)
}

func ExamplePortfolioSummary() {
	account := ""
	accounts := []mujica.PortfolioAccount{{ID: "a", Name: "Brokerage", Total: 120000, DayPnL: 1250, TotalReturn: 20000, Cost: 100000}, {ID: "b", Name: "Retirement", Total: 80000, DayPnL: -300, TotalReturn: 5000, Cost: 75000}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.PortfolioSummary(c, &account, accounts, mujica.PortfolioSummaryOptions{})
		if r.Changed() {
			fmt.Println(account)
		}
	}
	ui.Render(view, 960, 440, 1)
}

func ExampleAssetAllocationChart() {
	items := []mujica.AllocationItem{{Name: "Stocks", Value: 60, Target: 50}, {Name: "Bonds", Value: 30, Target: 40}, {Name: "Cash", Value: 10, Target: 10}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.AssetAllocationChart(c, items, mujica.AssetAllocationChartOptions{Frame: mujica.ChartFrame{Label: "Allocation"}})
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
		mujica.Use(c, mujica.Settings{})
		r := mujica.PerformanceChart(c, values, mujica.PerformanceChartOptions{Frame: mujica.ChartFrame{Label: "Performance"}, Times: times, Benchmark: benchmark})
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
	state := mujica.CurrencyConverterState{Amount: 100, From: "USD", To: "JPY"}
	currencies := []mujica.Currency{{Code: "USD", Name: "Dollar", Decimals: 2}, {Code: "JPY", Name: "Yen", Decimals: 0}}
	rate := 149.37
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.CurrencyConverter(c, &state, currencies, mujica.CurrencyConverterOptions{Rate: rate})
		if r.PairChanged {
			rate = 0
			fmt.Println(state.From, state.To)
		}
	}
	ui.Render(view, 960, 440, 1)
}
