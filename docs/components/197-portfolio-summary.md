# 197 PortfolioSummary

## 用途

资产总览。

## 最小示例

```go
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
```

## 参数

`account` 空值汇总全部账户；`accounts` 使用稳定 ID。`Currency` 默认 $；`Decimals` 默认 2；`Disabled` 默认 false。

## 状态

默认；交互区域支持聚焦与选中，Disabled 参数控制禁用。

## 事件

if r.Changed() { 调用方回调(account) }

## 键盘操作

Tab 聚焦账户组；方向键切换账户。

## 限制

收益、成本与资产数值由调用方计算；基数为 0 时隐藏比例。
