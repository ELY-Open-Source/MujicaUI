# 185 SymbolSearch

## 用途

代码搜索。

## 最小示例

```go
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
```

## 参数

`s.Query/Market` 为调用方搜索条件；`results` 是调用方异步结果。`Markets` nil 隐藏过滤；`Recent` 为空查询时显示；`Loading/Error/Disabled` 默认零值。

## 状态

默认、聚焦、禁用；异步结果支持加载、失败和空。

## 事件

if r.Changed || r.Picked || r.Retry { 调用方回调(state.Query, state.Market, r.Chosen.Code) }

## 键盘操作

Tab 聚焦输入与市场；Down 进入结果；Tab 到达结果按钮，Enter 选择。

## 限制

调用方处理搜索、市场过滤和最近记录；组件提供异步加载、失败和重试界面。
