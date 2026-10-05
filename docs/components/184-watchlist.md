# 184 Watchlist

## 用途

自选列表。

## 最小示例

```go
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
```

## 参数

`s` 保存分组、代码和隐藏列；`groups` 至少一个；`Decimals` 默认 2；`Disabled` 默认 false。

## 状态

默认；交互区域支持聚焦与选中，Disabled 参数控制禁用。

## 事件

if r.Changed || r.Submitted || r.Moved { 调用方回调(state.Selected, r.Rows, r.To) }

## 键盘操作

Tab 切换分组或列菜单；上下、Home/End 选择；Enter 确认；菜单 Space 开关列。

## 限制

拖动仅返回 Moved；调用方应用 Rows/To 调整顺序；刷新采用稳定 Symbol。
