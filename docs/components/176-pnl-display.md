# 176 PnLDisplay

盈亏显示。

## 用途

并列显示已实现、未实现与当日盈亏，带正负号、箭头与语义色；可切换隐私遮蔽。

## 最小示例

```go
hidden := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.PnLDisplay(c, finance.PnL{Realized: 1200.5, Unrealized: -340.2, Today: 88}, &hidden, finance.PnLDisplayOptions{Currency: "$"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `pnl` | `PnL` | 必填 | Realized、Unrealized、Today。 |
| `hidden` | `*bool` | 必填 | 为 true 时以“••••”遮蔽。 |
| `Decimals` | `int` | 2 | 小数位。 |
| `Currency` | `string` | "" | 金额前缀，例如 "$"。 |

## 状态

正、负、零；已遮蔽。

## 事件

`ChoiceResult.Changed()` 在遮蔽切换时为 true。

## 键盘操作

Tab 到达遮蔽按钮，Enter/Space 切换。

## 限制

- 遮蔽时无障碍名称只读出“已隐藏”，不读出金额。
- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
