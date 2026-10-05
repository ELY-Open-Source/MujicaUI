# 489 ChoiceChips 筛选芯片

## 用途

单选或多选的筛选芯片，自动换行。选中芯片用选中底色、强调边框和对勾表示。

## 最小示例

```go
var status []string
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.ChoiceChips(c, &status, []input.ToggleItem[string]{
		{Value: "open", Label: "Open"}, {Value: "closed", Label: "Closed"},
	}, input.ChoiceChipsOptions{Multiple: true, Label: "Status"})
}
```

签名：`func ChoiceChips[T comparable](c *ui.Context, value *[]T, items []ToggleItem[T], opts ChoiceChipsOptions) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*[]T` | 必填 | 选中的值；单选时多于一个 panic。 |
| `items` | `[]ToggleItem[T]` | 必填 | 芯片（Value、Label、Icon、Name、Disabled）。 |
| `Multiple` | `bool` | `false` | 允许多选；单选时可再次点击取消。 |
| `Disabled` | `bool` | `false` | 整组禁用。 |
| `Label` | `string` | `""` | 组名称。 |

## 状态

默认、悬停、按下、焦点、选中、禁用；颜色过渡 120 ms。

## 事件

`Changed()`：选择变化。

## 键盘操作

整组是一个 Tab 停靠点；左右方向键移动并跳过禁用芯片；Space/Enter 切换。

## 限制

无。
