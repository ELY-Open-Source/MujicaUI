# 200 CurrencyConverter

## 用途

货币换算。

## 最小示例

```go
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
```

## 参数

`RateFrom/RateTo` 指定汇率所属货币对，支持数值相同的刷新结果；货币小数位范围 0—12。

`s` 保存金额及货币对；`currencies` 至少两种，Decimals 为最小单位精度。`Rate` 默认 0 表示待提供；`RateTime` 零值隐藏；`Disabled` 默认 false。

## 状态

默认；交互区域支持聚焦与选中，Disabled 参数控制禁用。

## 事件

if r.PairChanged { rate = 0; 调用方回调(state.From, state.To) }

## 键盘操作

Tab 到达金额、币种和互换；方向键步进、选择；Enter/Space 互换。

## 限制

币种切换产生 PairChanged，调用方立即将 Rate 置 0 并获取新汇率；结果按目标货币精度取整。
