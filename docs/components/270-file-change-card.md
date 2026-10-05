# 270 FileChangeCard 文件变更卡片

## 用途

预览一个文件的变更：增删行数、首处修改附近的差异行，提供接受、拒绝和打开完整差异的入口。

## 最小示例

```go
decision := mujica.FilePending
change := mujica.FileChange{Path: "ledger.md", Old: "Guests: 40\n", New: "Guests: 42\n"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.FileChangeCard(c, &decision, change, mujica.FileChangeCardOptions{}).OpenDiff {
		fmt.Println("open the full diff")
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func FileChangeCard(c *ui.Context, decision *FileDecision, change FileChange, opts FileChangeCardOptions) FileChangeResult`。`decision` 为 nil、`Path` 为空或 `Preview` 为负时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Preview` | `int` | 8 | 从首处修改前两行起显示的差异行数。 |
| `Disabled` | `bool` | `false` | 禁止决定与打开差异。 |

## 状态

- `FileChange` 字段：`Path`、`Old`、`New`；`Old` 为空表示新文件，`New` 为空表示删除。
- 差异由 `internal.LineDiff`（Myers 行差异）计算，文本不变时只算一次。
- 新增行带“+”、删除行带“−”和行号，颜色只是辅助。
- 决定状态 `FilePending`、`FileAccepted`、`FileRejected` 以图标和文字显示；已决定后显示“撤销”。

## 事件

- `Decided`：`*decision` 被用户改变。
- `OpenDiff`：请求调用方打开完整差异。

## 键盘操作

- Tab 依次到“打开差异”“拒绝”“接受”（或“撤销”），Enter / Space 触发。

## 限制

- 只做行级差异；词级高亮由 Git 组扩展。
