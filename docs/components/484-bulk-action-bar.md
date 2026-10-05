# 484 BulkActionBar 批量操作条

## 用途

多选后出现的操作条：显示选中数、调用方提供的批量操作和“清除选择”。

## 最小示例

```go
selected := 3
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := input.BulkActionBar(c, selected, input.BulkActionBarOptions{}, func() {
		input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary})
	})
	if r.Changed() {
		selected = 0
	}
}
```

签名：`func BulkActionBar(c *ui.Context, selected int, opts BulkActionBarOptions, actions func()) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `selected` | `int` | 必填 | 选中数；0 时不显示（`Element` 为 nil），负数 panic。 |
| `actions` | `func()` | `nil` | 批量操作按钮插槽。 |
| `Disabled` | `bool` | `false` | 整条不可用。 |

## 状态

浮层样式，160 ms 淡入；选中数变化时播报。

## 事件

`Changed()`：用户清除选择（按钮或 Escape）。

## 键盘操作

Tab 进入操作按钮与“清除选择”；焦点在条内时 Escape 清除选择。

## 限制

位置由调用方放置，不自动悬浮在窗口底部。
