# 488 MentionInput 提及输入

## 用途

输入 @ 后弹出候选，选中后插入整体提及；提及作为一个单元被光标跳过、被一次删除。文本与提及列表都由调用方持有。

## 最小示例

```go
text, mentions := "", []mujica.Mention(nil)
people := []mujica.MentionCandidate{{ID: "u1", Name: "Alice"}, {ID: "u2", Name: "Bob"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.MentionInput(c, &text, &mentions, mujica.MentionInputOptions{
		Label: "Comment",
		Candidates: func(q string) []mujica.MentionCandidate {
			var out []mujica.MentionCandidate
			for _, p := range people {
				if strings.HasPrefix(strings.ToLower(p.Name), strings.ToLower(q)) {
					out = append(out, p)
				}
			}
			return out
		},
	})
}
```

签名：`func MentionInput(c *ui.Context, value *string, mentions *[]Mention, opts MentionInputOptions) InputResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*string` | 必填 | 全文。 |
| `mentions` | `*[]Mention` | 必填 | 提及（ID、Name、Start、End，按 rune 计）；与文本不符 panic。 |
| `Candidates` | `func(string) []MentionCandidate` | 必填 | 按 @ 后文字返回候选，仅在该文字变化时调用。 |
| `Placeholder` | `string` | `""` | 占位文字。 |
| `Label` | `string` | `""` | 无障碍名称。 |
| `Disabled` | `bool` | `false` | 不可编辑。 |
| `Error` | `string` | `""` | 错误信息。 |

## 状态

提及以强调色加浅底显示；候选列表用浮层样式，当前项有选中底色与左侧标记。输入法组合期间不提交。

## 事件

`Changed()`：文本或提及变化；`Submitted()`：候选关闭时按 Enter。

## 键盘操作

@ 后 Up/Down 选择候选，Enter 或 Tab 插入，Escape 关闭；左右键跨过提及，Backspace/Delete 整体删除；Shift+方向键选择；支持复制、剪切、粘贴、全选、撤销、重做命令。

## 限制

MyGo 编辑器没有设置光标的公开接口（ui/editor.go:117 只在文本变化时截断光标），无法在插入提及后把光标放到提及之后，因此本组件自绘单行编辑器；不支持多行与双击选词。
