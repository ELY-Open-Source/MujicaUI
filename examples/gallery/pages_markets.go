package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
)

var marketsPages = []page{
	{ID: "181", Name: "OptionChain", ZhName: "期权链", Category: "Markets", Demo: marketsOptionChainDemo, Code: `var state mujica.OptionChainState
rows := []mujica.OptionStrike{{Strike: 100, Call: mujica.OptionQuote{Bid: 3, Ask: 3.2, Last: 3.1}, Put: mujica.OptionQuote{Bid: 2, Ask: 2.2, Last: 2.1}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.OptionChain(c, &state, []string{"Oct 17", "Nov 21"}, rows, mujica.OptionChainOptions{Spot: 100.4})
	if r.Changed || r.Submitted || r.ExpiryChanged {
		fmt.Println(state.Expiry, state.Strike, state.Side)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 切换到期日或列表；上下、Home/End 选择行；左右切换认购认沽；Enter 确认。"},
	{ID: "182", Name: "GreeksTable", ZhName: "希腊值表", Category: "Markets", Demo: marketsGreeksTableDemo, Code: `rows := []mujica.GreeksRow{{Label: "AAPL 100 C", Greeks: mujica.OptionGreeks{Delta: 0.5, Gamma: 0.03, Theta: -0.1, Vega: 0.2, Rho: 0.05}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.GreeksTable(c, rows, mujica.GreeksTableOptions{Notes: true})
}
ui.Render(view, 960, 440, 1)
`, Keys: "展示组件；表头提供说明和悬停提示。"},
	{ID: "183", Name: "PayoffDiagram", ZhName: "到期损益图", Category: "Markets", Demo: marketsPayoffDiagramDemo, Code: `legs := []mujica.OptionLeg{{Kind: mujica.LegCall, Strike: 100, Premium: 3, Quantity: 1}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.PayoffDiagram(c, legs, mujica.PayoffDiagramOptions{Frame: mujica.ChartFrame{Label: "Call payoff"}, Spot: 103})
	if r.Changed() || r.Retried() {
		fmt.Println(r.Point)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦；方向键、Home/End 逐点读数；Escape 清除。"},
	{ID: "184", Name: "Watchlist", ZhName: "自选列表", Category: "Markets", Demo: marketsWatchlistDemo, Code: `var state mujica.WatchlistState
groups := []mujica.WatchGroup{{Name: "Technology", Items: []mujica.WatchItem{{Symbol: "AAPL", Name: "Apple", Price: 227.5, Change: 2.7, ChangePercent: 1.2}}}, {Name: "Energy", Items: []mujica.WatchItem{{Symbol: "XOM", Price: 117.4}}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.Watchlist(c, &state, groups, mujica.WatchlistOptions{})
	if r.Changed || r.Submitted || r.Moved {
		fmt.Println(state.Selected, r.Rows, r.To)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 切换分组或列菜单；上下、Home/End 选择；Enter 确认；菜单 Space 开关列。"},
	{ID: "185", Name: "SymbolSearch", ZhName: "代码搜索", Category: "Markets", Demo: marketsSymbolSearchDemo, Code: `var state mujica.SymbolSearchState
rows := []mujica.MarketSymbol{{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.SymbolSearch(c, &state, rows, mujica.SymbolSearchOptions{Markets: []string{"NASDAQ", "NYSE"}, Recent: rows})
	if r.Changed || r.Picked || r.Retry {
		fmt.Println(state.Query, state.Market, r.Chosen.Code)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦输入与市场；Down 进入结果；Tab 到达结果按钮，Enter 选择。"},
	{ID: "186", Name: "SymbolBadge", ZhName: "代码徽标", Category: "Markets", Demo: marketsSymbolBadgeDemo, Code: `symbol := mujica.MarketSymbol{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.SymbolBadge(c, symbol, mujica.SymbolBadgeOptions{Clickable: true})
	if r.Changed() {
		fmt.Println(symbol.Code)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Clickable 开启时 Tab 聚焦，Enter/Space 激活。"},
	{ID: "187", Name: "MarketHeatmap", ZhName: "市场热力图", Category: "Markets", Demo: marketsMarketHeatmapDemo, Code: `var path []string
root := mujica.HeatNode{Name: "Market", Children: []mujica.HeatNode{{Name: "Technology", Children: []mujica.HeatNode{{Name: "AAPL", Cap: 3000, Change: 1.2}, {Name: "MSFT", Cap: 2800, Change: -0.6}}}, {Name: "Energy", Cap: 450, Change: -2.1}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.MarketHeatmap(c, &path, root, mujica.MarketHeatmapOptions{Frame: mujica.ChartFrame{Label: "Market"}})
	if r.Changed() || r.ValueChanged() || r.Retried() {
		fmt.Println(path, r.Point)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦图；方向键逐块读数；Enter 下钻；Backspace 返回；路径按钮可用 Tab。"},
	{ID: "188", Name: "Screener", ZhName: "选股器", Category: "Markets", Demo: marketsScreenerDemo, Code: `var state mujica.ScreenerState
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
`, Keys: "Tab 到达字段、比较、数值和保存；Select 方向键选择；NumberInput 方向键步进；表格方向键选择。"},
	{ID: "189", Name: "MarketOverview", ZhName: "市场概览", Category: "Markets", Demo: marketsMarketOverviewDemo, Code: `indices := []mujica.IndexQuote{{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.MarketOverview(c, indices, mujica.MarketOverviewOptions{Breadth: mujica.MarketBreadth{Advancing: 2100, Declining: 1400, Unchanged: 300}})
	if r.Chosen >= 0 {
		fmt.Println(indices[r.Chosen].Name)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 到达指数卡片；Enter/Space 选择。"},
	{ID: "190", Name: "IndexCard", ZhName: "指数卡片", Category: "Markets", Demo: marketsIndexCardDemo, Code: `quote := mujica.IndexQuote{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.IndexCard(c, quote, mujica.IndexCardOptions{Clickable: true})
	if r.Changed() {
		fmt.Println(quote.Name)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦卡片或迷你图；Enter/Space 激活卡片，图表方向键读数。"},
	{ID: "191", Name: "MarketStatus", ZhName: "市场状态", Category: "Markets", Demo: marketsMarketStatusDemo, Code: `now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.MarketStatus(c, mujica.MarketOpen, mujica.MarketStatusOptions{Market: "NYSE", Now: now, Next: now.Add(5 * time.Hour), NextPhase: mujica.MarketAfterHours})
}
ui.Render(view, 960, 440, 1)
`, Keys: "展示组件。"},
	{ID: "192", Name: "TradingSessionClock", ZhName: "交易时段时钟", Category: "Markets", Demo: marketsTradingSessionClockDemo, Code: `now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
markets := []mujica.SessionMarket{{Name: "NYSE", Location: time.FixedZone("EDT", -4*3600), Sessions: []mujica.TradingSession{{Phase: mujica.MarketOpen, Start: 570, End: 960}}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TradingSessionClock(c, markets, mujica.TradingSessionClockOptions{Now: now})
}
ui.Render(view, 960, 440, 1)
`, Keys: "展示组件。"},
	{ID: "193", Name: "EconomicCalendar", ZhName: "经济日历", Category: "Markets", Demo: marketsEconomicCalendarDemo, Code: `var state mujica.EconomicCalendarState
rows := []mujica.EconomicEvent{{ID: "jobs", Country: "US", Title: "Payrolls", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Importance: 3, Forecast: "150K", Previous: "142K", Actual: "254K"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.EconomicCalendar(c, &state, rows, mujica.EconomicCalendarOptions{})
	if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
		fmt.Println(state.Table.Selected())
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 到达国家过滤、重要性和表格；方向键选择；Enter 提交行。"},
	{ID: "194", Name: "EarningsCalendar", ZhName: "财报日历", Category: "Markets", Demo: marketsEarningsCalendarDemo, Code: `var state mujica.DataTableState[string]
rows := []mujica.EarningsEvent{{ID: "a", Symbol: "AAPL", Name: "Apple", Date: mujica.Date{Year: 2026, Month: 10, Day: 30}, Timing: mujica.EarningsAfterClose, EPSEstimate: 1.6, EPSActual: 1.64}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.EarningsCalendar(c, &state, rows, mujica.EarningsCalendarOptions{})
	if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
		fmt.Println(state.Selected())
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦表格；方向键及 Home/End 选择；Enter 提交；表头支持键盘排序。"},
	{ID: "195", Name: "NewsFeed", ZhName: "新闻流", Category: "Markets", Demo: marketsNewsFeedDemo, Code: `var state mujica.NewsFeedState
rows := []mujica.NewsItem{{ID: "a", Title: "Apple unveils new chips", Source: "Wire", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Symbols: []string{"AAPL"}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.NewsFeed(c, &state, rows, mujica.NewsFeedOptions{})
	if r.Opened != "" || r.Symbol != "" || r.Retry {
		fmt.Println(r.Opened, r.Symbol)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦列表或相关代码；上下选择；Enter 打开并标记已读。"},
	{ID: "196", Name: "SentimentGauge", ZhName: "情绪仪表", Category: "Markets", Demo: marketsSentimentGaugeDemo, Code: `value := 72.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.SentimentGauge(c, value, mujica.SentimentGaugeOptions{Label: "Fear & Greed"})
}
ui.Render(view, 960, 440, 1)
`, Keys: "展示组件，数值与文字同时呈现。"},
	{ID: "197", Name: "PortfolioSummary", ZhName: "资产总览", Category: "Markets", Demo: marketsPortfolioSummaryDemo, Code: `account := ""
accounts := []mujica.PortfolioAccount{{ID: "a", Name: "Brokerage", Total: 120000, DayPnL: 1250, TotalReturn: 20000, Cost: 100000}, {ID: "b", Name: "Retirement", Total: 80000, DayPnL: -300, TotalReturn: 5000, Cost: 75000}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.PortfolioSummary(c, &account, accounts, mujica.PortfolioSummaryOptions{})
	if r.Changed() {
		fmt.Println(account)
	}
}
ui.Render(view, 960, 440, 1)
`, Keys: "Tab 聚焦账户组；方向键切换账户。"},
	{ID: "198", Name: "AssetAllocationChart", ZhName: "资产配置图", Category: "Markets", Demo: marketsAssetAllocationChartDemo, Code: `items := []mujica.AllocationItem{{Name: "Stocks", Value: 60, Target: 50}, {Name: "Bonds", Value: 30, Target: 40}, {Name: "Cash", Value: 10, Target: 10}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.AssetAllocationChart(c, items, mujica.AssetAllocationChartOptions{Frame: mujica.ChartFrame{Label: "Allocation"}})
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
`, Keys: "Tab 聚焦图和图例；方向键读数；Space 切换序列。"},
	{ID: "200", Name: "CurrencyConverter", ZhName: "货币换算", Category: "Markets", Demo: marketsCurrencyConverterDemo, Code: `state := mujica.CurrencyConverterState{Amount: 100, From: "USD", To: "JPY"}
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
`, Keys: "Tab 到达金额、币种和互换；方向键步进、选择；Enter/Space 互换。"},
}

func marketsDataState(c *ui.Context, id string, log func(string)) (mujica.DataStatus, bool, string) {
	state := ui.Local(c.Root(), id+".status", func() int { return 0 })
	if mujica.RadioGroup(c, state, []mujica.RadioItem[int]{{Value: 0, Label: "Ready"}, {Value: 1, Label: "Loading"}, {Value: 2, Label: "Empty"}, {Value: 3, Label: "Error"}}, mujica.RadioGroupOptions{Horizontal: true}).Changed() {
		log(fmt.Sprint("status ", *state))
	}
	switch *state {
	case 1:
		return mujica.DataLoading, false, ""
	case 2:
		return mujica.DataReady, true, ""
	case 3:
		return mujica.DataFailed, false, "Sample failed"
	}
	return mujica.DataReady, false, ""
}
