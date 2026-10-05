# 178 MarginIndicator

保证金率指示。

## 用途

显示保证金率（维持保证金 / 权益）及其与警戒线、强平线的关系。

## 最小示例

```go
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.MarginIndicator(c, 0.86, mujica.MarginIndicatorOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `ratio` | `float64` | 必填 | 保证金率，0 以上。 |
| `Warning / Liquidation` | `float64` | 0.8 / 1 | 警戒与强平线；须 0 < Warning < Liquidation。 |
| `Label` | `string` | 内置“Margin ratio” | 名称。 |

## 状态

安全、追加保证金、强制平仓。

## 事件

无。

## 键盘操作

不获取焦点；以计量角色读出比率、状态与两条线。

## 限制

- 状态同时用图标、文字与颜色表达；两条线粗细不同，并在下方以文字标注。
