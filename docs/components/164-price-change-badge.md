# 164 PriceChangeBadge

涨跌徽标。

## 用途

显示涨跌额与涨跌幅，带箭头、正负号与语义底色。

## 最小示例

```go
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PriceChangeBadge(c, -5.62, -1.35, mujica.PriceChangeBadgeOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `change` | `float64` | 必填 | 涨跌额。 |
| `percent` | `float64` | 必填 | 涨跌幅（%），固定 2 位小数。 |
| `Decimals` | `int` | 2 | 涨跌额小数位。 |
| `HidePercent` | `bool` | false | 只显示涨跌额。 |

## 状态

上涨、下跌、持平。

## 事件

无。

## 键盘操作

不获取焦点；无障碍名称含方向词与数值。

## 限制

- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
