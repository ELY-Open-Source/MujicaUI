# 182 GreeksTable

## 用途

希腊值表。

## 最小示例

```go
rows := []finance.GreeksRow{{Label: "AAPL 100 C", Greeks: finance.OptionGreeks{Delta: 0.5, Gamma: 0.03, Theta: -0.1, Vega: 0.2, Rho: 0.05}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.GreeksTable(c, rows, finance.GreeksTableOptions{Notes: true})
}
ui.Render(view, 960, 440, 1)
```

## 参数

小数位范围 0—12；0 表示默认 4 位。

`rows` 为合约与五个数值；NaN 显示破折号。`Decimals` 默认 4；`Notes` 默认 false。

## 状态

默认与数据边界。

## 事件

展示组件；数据由调用方更新。

## 键盘操作

展示组件；表头提供说明和悬停提示。

## 限制

通过虚拟列表构建可见行；精度为显示规则，数值计算由调用方提供。
