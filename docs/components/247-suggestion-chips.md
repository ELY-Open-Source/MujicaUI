# 247 SuggestionChips · 建议标签

## 用途

起始提示或追问建议。点击可填入输入框，或配置为直接发送。

## 最小示例

```go
input := ""
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.SuggestionChips(c, &input, []string{"Plan a ball", "Write a sonnet"}, mujica.SuggestionChipsOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func SuggestionChips(c *ui.Context, input *string, suggestions []string, opts SuggestionChipsOptions) SuggestionChipsView`。建议不能为空字符串。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Send` | `bool` | false | 点击即请求发送；默认把建议填入 `*input`。 |
| `Label` | `string` | "" | 组名称，如“追问建议”。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认（右箭头）、发送模式（上箭头）、悬停、按下、聚焦、禁用。

## 事件

`Chosen`：点击的建议。`Sent` 为 true 时调用方应立即发送 `Chosen`；否则建议已写入 `*input`。

## 键盘操作

- Tab 进入组；←/→ 在建议间移动；Enter/Space 选择。

## 限制

- 建议换行排列；数量很多时请调用方截取。
