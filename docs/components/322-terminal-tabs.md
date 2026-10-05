# 322 TerminalTabs · 终端标签

## 用途

管理多个终端的标签：切换、重命名、关闭；仍有进程运行（`Busy`）的终端关闭前弹出确认。

## 最小示例

```go
tabs := []mujica.TerminalTab{{ID: "a", Title: "zsh"}, {ID: "b", Title: "build", Busy: true}}
active := "a"
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.TerminalTabs(c, &tabs, &active, mujica.TerminalTabsOptions{}, func(id string) {
		ui.Text(c, "terminal "+id)
	})
	if r.Closed != "" {
		fmt.Println("close pty", r.Closed)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func TerminalTabs(c *ui.Context, tabs *[]TerminalTab, active *string, opts TerminalTabsOptions, content func(id string)) TerminalTabsResult`。`*tabs` 与 `*active` 由调用方持有，组件直接修改：重命名改 `Title`，关闭删除元素并把 `active` 移到相邻标签。`active` 不是任何标签 ID 时 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `TerminalTab.ID` | `string` | 必填 | 稳定标识，也是事件值。 |
| `TerminalTab.Title` | `string` | 必填 | 标签文字。 |
| `TerminalTab.Busy` | `bool` | false | 进程仍在运行：显示小圆点并在关闭前确认。 |
| `Disabled` | `bool` | false | 不响应切换、重命名与关闭。 |
| `content` | `func(id string)` | 必填 | 构建当前标签的内容。 |

## 状态

当前标签（Selection 背景、粗体与底部 AccentText 线）、悬停、键盘焦点环、正在运行（圆点与无障碍描述“进程运行中”）、重命名输入框、关闭确认对话框、无标签时的空状态、禁用。

## 事件

`Changed`：当前标签变化；`Renamed`：被重命名的标签 ID；`Closed`：被关闭的标签 ID（调用方应关闭对应 PTY）。

## 键盘操作

- ←/→、Home/End：在标签间移动并切换。
- F2 或双击：重命名；Enter 确认，Escape 取消。
- Delete：关闭焦点所在标签；运行中的标签先确认，确认框默认聚焦“取消”。

## 限制

- 不提供新建按钮与拖动排序；新建由调用方向 `*tabs` 追加。
