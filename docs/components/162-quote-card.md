# 162 QuoteCard

报价卡片。

## 用途

显示代码、名称、现价、涨跌额与涨跌幅，并提供迷你图插槽。

## 最小示例

```go
quote := mujica.Quote{Symbol: "MUJC", Name: "Ave Mujica Holdings", Price: 412.08, Change: -5.62, ChangePercent: -1.35}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.QuoteCard(c, quote, mujica.QuoteCardOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `q` | `Quote` | 必填 | Symbol 必填，空值 panic；Name、Price、Change、ChangePercent。 |
| `Decimals` | `int` | 2 | 价格与涨跌额小数位。 |
| `Chart` | `func()` | nil | 迷你图插槽，高 40 DIP，例如 Sparkline。 |

## 状态

默认；涨、跌、平三种方向。

## 事件

无。

## 键盘操作

卡片本身不获取焦点；插槽内控件自行处理键盘。

## 限制

- 不计算涨跌，数值由调用方提供。
- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
