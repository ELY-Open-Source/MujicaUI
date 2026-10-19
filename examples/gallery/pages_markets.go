package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

var marketsPages = []page{
	{ID: "181", Name: "OptionChain", ZhName: "期权链", Category: "Markets", Demo: marketsOptionChainDemo, Code: `var state finance.OptionChainState
rows := []finance.OptionStrike{{Strike: 100, Call: finance.OptionQuote{Bid: 3, Ask: 3.2, Last: 3.1}, Put: finance.OptionQuote{Bid: 2, Ask: 2.2, Last: 2.1}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.OptionChain(c, &state, []string{"Oct 17", "Nov 21"}, rows, finance.OptionChainOptions{Spot: 100.4})
	if r.Changed || r.Submitted || r.ExpiryChanged {
		fmt.Println(state.Expiry, state.Strike, state.Side)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 切换到期日或列表；上下、Home/End 选择行；左右切换认购认沽；Enter 确认。"},
	{ID: "182", Name: "GreeksTable", ZhName: "希腊值表", Category: "Markets", Demo: marketsGreeksTableDemo, Code: `rows := []finance.GreeksRow{{Label: "AAPL 100 C", Greeks: finance.OptionGreeks{Delta: 0.5, Gamma: 0.03, Theta: -0.1, Vega: 0.2, Rho: 0.05}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.GreeksTable(c, rows, finance.GreeksTableOptions{Notes: true})
}
ui.Render(view, 960, 440, 1)
`, Keys: "展示组件；表头提供说明和悬停提示。"},
	{ID: "183", Name: "PayoffDiagram", ZhName: "到期损益图", Category: "Markets", Demo: marketsPayoffDiagramDemo, Code: `legs := []finance.OptionLeg{{Kind: finance.LegCall, Strike: 100, Premium: 3, Quantity: 1}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.PayoffDiagram(c, legs, finance.PayoffDiagramOptions{Frame: chart.ChartFrame{Label: "Call payoff"}, Spot: 103})
	if r.Changed() || r.Retried() {
		fmt.Println(r.Point)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦；方向键、Home/End 逐点读数；Escape 清除。"},
	{ID: "184", Name: "Watchlist", ZhName: "自选列表", Category: "Markets", Demo: marketsWatchlistDemo, Code: `var state finance.WatchlistState
groups := []finance.WatchGroup{{Name: "Technology", Items: []finance.WatchItem{{Symbol: "AAPL", Name: "Apple", Price: 227.5, Change: 2.7, ChangePercent: 1.2}}}, {Name: "Energy", Items: []finance.WatchItem{{Symbol: "XOM", Price: 117.4}}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.Watchlist(c, &state, groups, finance.WatchlistOptions{})
	if r.Changed || r.Submitted || r.Moved {
		fmt.Println(state.Selected, r.Rows, r.To)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 切换分组或列菜单；上下、Home/End 选择；Enter 确认；菜单 Space 开关列。"},
	{ID: "185", Name: "SymbolSearch", ZhName: "代码搜索", Category: "Markets", Demo: marketsSymbolSearchDemo, Code: `var state finance.SymbolSearchState
rows := []finance.MarketSymbol{{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.SymbolSearch(c, &state, rows, finance.SymbolSearchOptions{Markets: []string{"NASDAQ", "NYSE"}, Recent: rows})
	if r.Changed || r.Picked || r.Retry {
		fmt.Println(state.Query, state.Market, r.Chosen.Code)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦输入与市场；Down 进入结果；Tab 到达结果按钮，Enter 选择。"},
	{ID: "186", Name: "SymbolBadge", ZhName: "代码徽标", Category: "Markets", Demo: marketsSymbolBadgeDemo, Code: `symbol := finance.MarketSymbol{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.SymbolBadge(c, symbol, finance.SymbolBadgeOptions{Clickable: true})
	if r.Changed() {
		fmt.Println(symbol.Code)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Clickable 开启时 Tab 聚焦，Enter/Space 激活。"},
	{ID: "187", Name: "MarketHeatmap", ZhName: "市场热力图", Category: "Markets", Demo: marketsMarketHeatmapDemo, Code: `var path []string
root := chart.HeatNode{Name: "Market", Children: []chart.HeatNode{{Name: "Technology", Children: []chart.HeatNode{{Name: "AAPL", Cap: 3000, Change: 1.2}, {Name: "MSFT", Cap: 2800, Change: -0.6}}}, {Name: "Energy", Cap: 450, Change: -2.1}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := chart.MarketHeatmap(c, &path, root, chart.MarketHeatmapOptions{Frame: chart.ChartFrame{Label: "Market"}})
	if r.Changed() || r.ValueChanged() || r.Retried() {
		fmt.Println(path, r.Point)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦图；方向键逐块读数；Enter 下钻；Backspace 返回；路径按钮可用 Tab。"},
	{ID: "188", Name: "Screener", ZhName: "选股器", Category: "Markets", Demo: marketsScreenerDemo, Code: `var state finance.ScreenerState
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
`, Keys: "Tab 到达字段、比较、数值和保存；Select 方向键选择；NumberInput 方向键步进；表格方向键选择。"},
	{ID: "189", Name: "MarketOverview", ZhName: "市场概览", Category: "Markets", Demo: marketsMarketOverviewDemo, Code: `indices := []finance.IndexQuote{{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.MarketOverview(c, indices, finance.MarketOverviewOptions{Breadth: finance.MarketBreadth{Advancing: 2100, Declining: 1400, Unchanged: 300}})
	if r.Chosen >= 0 {
		fmt.Println(indices[r.Chosen].Name)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 到达指数卡片；Enter/Space 选择。"},
	{ID: "190", Name: "IndexCard", ZhName: "指数卡片", Category: "Markets", Demo: marketsIndexCardDemo, Code: `quote := finance.IndexQuote{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.IndexCard(c, quote, finance.IndexCardOptions{Clickable: true})
	if r.Changed() {
		fmt.Println(quote.Name)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦卡片或迷你图；Enter/Space 激活卡片，图表方向键读数。"},
	{ID: "191", Name: "MarketStatus", ZhName: "市场状态", Category: "Markets", Demo: marketsMarketStatusDemo, Code: `now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.MarketStatus(c, finance.MarketOpen, finance.MarketStatusOptions{Market: "NYSE", Now: now, Next: now.Add(5 * time.Hour), NextPhase: finance.MarketAfterHours})
}
ui.Render(view, 960, 440, 1)
`, Keys: "展示组件。"},
	{ID: "192", Name: "TradingSessionClock", ZhName: "交易时段时钟", Category: "Markets", Demo: marketsTradingSessionClockDemo, Code: `now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
markets := []finance.SessionMarket{{Name: "NYSE", Location: time.FixedZone("EDT", -4*3600), Sessions: []finance.TradingSession{{Phase: finance.MarketOpen, Start: 570, End: 960}}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.TradingSessionClock(c, markets, finance.TradingSessionClockOptions{Now: now})
}
ui.Render(view, 960, 440, 1)
`, Keys: "展示组件。"},
	{ID: "193", Name: "EconomicCalendar", ZhName: "经济日历", Category: "Markets", Demo: marketsEconomicCalendarDemo, Code: `var state finance.EconomicCalendarState
rows := []finance.EconomicEvent{{ID: "jobs", Country: "US", Title: "Payrolls", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Importance: 3, Forecast: "150K", Previous: "142K", Actual: "254K"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.EconomicCalendar(c, &state, rows, finance.EconomicCalendarOptions{})
	if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
		fmt.Println(state.Table.Selected())
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 到达国家过滤、重要性和表格；方向键选择；Enter 提交行。"},
	{ID: "194", Name: "EarningsCalendar", ZhName: "财报日历", Category: "Markets", Demo: marketsEarningsCalendarDemo, Code: `var state data.DataTableState[string]
rows := []finance.EarningsEvent{{ID: "a", Symbol: "AAPL", Name: "Apple", Date: core.Date{Year: 2026, Month: 10, Day: 30}, Timing: finance.EarningsAfterClose, EPSEstimate: 1.6, EPSActual: 1.64}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.EarningsCalendar(c, &state, rows, finance.EarningsCalendarOptions{})
	if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
		fmt.Println(state.Selected())
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦表格；方向键及 Home/End 选择；Enter 提交；表头支持键盘排序。"},
	{ID: "195", Name: "NewsFeed", ZhName: "新闻流", Category: "Markets", Demo: marketsNewsFeedDemo, Code: `var state finance.NewsFeedState
rows := []finance.NewsItem{{ID: "a", Title: "Apple unveils new chips", Source: "Wire", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Symbols: []string{"AAPL"}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.NewsFeed(c, &state, rows, finance.NewsFeedOptions{})
	if r.Opened != "" || r.Symbol != "" || r.Retry {
		fmt.Println(r.Opened, r.Symbol)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦列表或相关代码；上下选择；Enter 打开并标记已读。"},
	{ID: "196", Name: "SentimentGauge", ZhName: "情绪仪表", Category: "Markets", Demo: marketsSentimentGaugeDemo, Code: `value := 72.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.SentimentGauge(c, value, finance.SentimentGaugeOptions{Label: "Fear & Greed"})
}
ui.Render(view, 960, 440, 1)
`, Keys: "展示组件，数值与文字同时呈现。"},
	{ID: "197", Name: "PortfolioSummary", ZhName: "资产总览", Category: "Markets", Demo: marketsPortfolioSummaryDemo, Code: `account := ""
accounts := []finance.PortfolioAccount{{ID: "a", Name: "Brokerage", Total: 120000, DayPnL: 1250, TotalReturn: 20000, Cost: 100000}, {ID: "b", Name: "Retirement", Total: 80000, DayPnL: -300, TotalReturn: 5000, Cost: 75000}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.PortfolioSummary(c, &account, accounts, finance.PortfolioSummaryOptions{})
	if r.Changed() {
		fmt.Println(account)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦账户组；方向键切换账户。"},
	{ID: "198", Name: "AssetAllocationChart", ZhName: "资产配置图", Category: "Markets", Demo: marketsAssetAllocationChartDemo, Code: `items := []finance.AllocationItem{{Name: "Stocks", Value: 60, Target: 50}, {Name: "Bonds", Value: 30, Target: 40}, {Name: "Cash", Value: 10, Target: 10}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.AssetAllocationChart(c, items, finance.AssetAllocationChartOptions{Frame: chart.ChartFrame{Label: "Allocation"}})
	if r.Changed() || r.Retried() {
		fmt.Println(r.Point)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦环形图；方向键逐项读数；Home/End 到首尾。"},
	{ID: "199", Name: "PerformanceChart", ZhName: "收益曲线", Category: "Markets", Demo: marketsPerformanceChartDemo, Code: `values := []float64{100, 110, 95, 120}
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
`, Keys: "Tab 聚焦图和图例；方向键读数；Space 切换序列。"},
	{ID: "200", Name: "CurrencyConverter", ZhName: "货币换算", Category: "Markets", Demo: marketsCurrencyConverterDemo, Code: `state := finance.CurrencyConverterState{Amount: 100, From: "USD", To: "JPY"}
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
`, Keys: "Tab 到达金额、币种和互换；方向键步进、选择；Enter/Space 互换。"},
}

func marketsDataState(c *ui.Context, id string, log func(string)) (data.DataStatus, bool, string) {
	state := ui.Local(c.Root(), id+".status", func() int { return 0 })
	if input.RadioGroup(c, state, []input.RadioItem[int]{{Value: 0, Label: "Ready"}, {Value: 1, Label: "Loading"}, {Value: 2, Label: "Empty"}, {Value: 3, Label: "Error"}}, input.RadioGroupOptions{Horizontal: true}).Changed() {
		log(fmt.Sprint("status ", *state))
	}
	switch *state {
	case 1:
		return data.DataLoading, false, ""
	case 2:
		return data.DataReady, true, ""
	case 3:
		return data.DataFailed, false, "Sample failed"
	}
	return data.DataReady, false, ""
}
