# 271 MultiFileDiffReview 多文件差异审阅

## 用途

审阅多个文件的变更：左侧文件列表，右侧所选文件的完整差异，逐文件接受或拒绝，全部接受，并显示审阅进度。

## 最小示例

```go
files := []mujica.FileChange{
	{Path: "ledger.md", Old: "Guests: 40\n", New: "Guests: 42\n"},
	{Path: "menu.md", Old: "soup\n", New: "soup\ntart\n"},
}
decisions := make([]mujica.FileDecision, len(files))
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.MultiFileDiffReview(c, &decisions, files, mujica.MultiFileDiffReviewOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func MultiFileDiffReview(c *ui.Context, decisions *[]FileDecision, files []FileChange, opts MultiFileDiffReviewOptions) ChoiceResult`。`decisions` 与 `files` 数量不同、`Path` 为空或重复、`Height` 为负时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Height` | `float32` | 420 | 组件高度。 |
| `Disabled` | `bool` | `false` | 禁止选择与决定。 |

## 状态

- 进度为“已审阅 N / 总数”，以文字和细进度条显示。
- 所选文件用选中底色和左侧竖条标出。
- 文件超过 200 个或差异超过 200 行时虚拟化。
- “全部接受”只接受尚未决定的文件，全部决定后禁用。

## 事件

- `Changed()`：用户改变了任一文件的决定。

## 键盘操作

- Tab 进入文件列表；↑ ↓ 移动，Enter 选择；Tab 到右侧按钮后 Enter 接受或拒绝。

## 限制

- 文件很多时方向键只在已构建的可见行间移动，滚动后再继续。
