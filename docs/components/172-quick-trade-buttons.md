# 172 QuickTradeButtons

快速交易按钮。

## 用途

以预设数量一键市价买入或卖出；可配置在发送前就近确认。

## 最小示例

```go
quantity := 10.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if side, q, ok := mujica.QuickTradeButtons(c, &quantity, mujica.QuickTradeButtonsOptions{Symbol: "MUJC", Confirm: true}).Traded(); ok {
		fmt.Println(side, q)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `quantity` | `*float64` | 必填 | 当前数量，须为 Presets 之一，0 取第一项。 |
| `Symbol` | `string` | 必填 | 代码。 |
| `Presets` | `[]float64` | 1, 10, 100 | 数量预设。 |
| `Decimals` | `int` | 0 | 数量小数位。 |
| `Confirm` | `bool` | false | 发送前弹出确认，默认聚焦取消。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、悬停、按下、聚焦、确认中、禁用。

## 事件

`QuickTradeResult.Traded()` 返回方向与数量。

## 键盘操作

Tab 到达数量分段与买卖按钮；Enter/Space 触发；确认浮层中 Enter 默认取消，Escape 关闭。

## 限制

- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
