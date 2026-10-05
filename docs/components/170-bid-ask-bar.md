# 170 BidAskBar

买卖力量条。

## 用途

以分段条对比买卖成交量，两端标出占比与数量。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.BidAskBar(c, 6234, 3820, finance.BidAskBarOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `bid / ask` | `float64` | 必填 | 买量与卖量，负数或 NaN panic。 |
| `Decimals` | `int` | 0 | 数量小数位。 |
| `Label` | `string` | 内置“Bid and ask volume” | 无障碍名称。 |

## 状态

默认；两者都为 0 时显示中性灰条与 50%。

## 事件

无。

## 键盘操作

不获取焦点；以计量角色读出买卖占比。

## 限制

- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
