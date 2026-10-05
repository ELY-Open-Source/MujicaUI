# 163 PriceText

价格文本。

## 用途

按固定小数位与千分位显示价格；可与参考价比较着色并加箭头；价格变化时闪烁一次。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	core.UseTrend(c, core.RedUp)
	finance.PriceText(c, 18234.5, finance.PriceTextOptions{Compare: true, Reference: 18000})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `float64` | 必填 | 价格。 |
| `Decimals` | `int` | 2 | 小数位，0 表示 2，负数 panic。 |
| `NoGrouping` | `bool` | false | 不显示千分位。 |
| `Compare / Reference` | `bool / float64` | false / 0 | 按相对 Reference 的方向着色并加箭头。 |
| `NoFlash` | `bool` | false | 变化时不闪烁。 |
| `Size` | `float32` | 正文字号 | 字号（DIP），随系统文字缩放。 |

## 状态

默认；上涨、下跌、持平；变化闪烁（360 ms 淡出）。

## 事件

无。

## 键盘操作

不获取焦点。

## 限制

- 闪烁由帧间比较值得出：同一元素的值改变才会闪烁。减少动画时不闪烁，结束后不再请求重绘。
- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
