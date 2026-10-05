# 169 SpreadIndicator

价差指示。

## 用途

显示卖一减买一的价差与占中间价的百分比；超过阈值或买卖倒挂时给出带图标的文字提示。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.SpreadIndicator(c, 99.95, 100.05, finance.SpreadIndicatorOptions{WarnPercent: 0.5})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `bid / ask` | `float64` | 必填 | 买一价与卖一价。 |
| `Decimals` | `int` | 2 | 价差小数位。 |
| `WarnPercent` | `float64` | 0 | 价差百分比达到即提示；0 不提示，负数 panic。 |

## 状态

正常、价差过大、倒挂。

## 事件

无。

## 键盘操作

不获取焦点；以状态角色读出。

## 限制

- 无。
