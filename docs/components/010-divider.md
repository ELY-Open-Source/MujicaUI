# 010 Divider 分隔线

## 用途

横向或纵向的细分隔线；横向可在中间放文字；可选旧金色单个菱形装饰。永不获取焦点，装饰不进入无障碍树。

## 最小示例

```go
ui.Column(c).Padding(16).Gap(12).Children(func() {
	layout.Divider(c, layout.DividerOptions{Ornament: true})
	layout.Divider(c, layout.DividerOptions{Label: "Chapter II"})
})
```

签名：`func Divider(c *ui.Context, opts DividerOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Vertical` | `bool` | false | 纵向，放在横向行中，高度随行拉伸。 |
| `Label` | `string` | "" | 横向分隔线中间的文字（衬线 12 DIP，TextMuted）。纵向带文字时 panic。 |
| `Ornament` | `bool` | false | 线条改用 Ornament 色，并在中间画一个 7 DIP 菱形；有 `Label` 时文字取代菱形。 |

## 状态

- 默认：1 DIP Border 色细线。
- 装饰：Ornament 色线与菱形。
- 无悬停、焦点等交互状态。

## 事件

无。

## 键盘操作

不获取焦点，Tab 直接跳过。

## 限制

- 方向不会根据父级自动判断，需要显式设置 `Vertical`。
- 纵向不支持文字。
