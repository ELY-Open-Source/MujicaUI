# 415 MailList 邮件列表

## 用途

展示虚拟化邮件列表，包含未读粗体、星标、附件计数、多选与打开事件。

## 最小示例

```go
var s mujica.ListState[string]
mails := []mujica.MailItem{{ID: "1", From: "Sakiko", Subject: "Second act", Unread: true, Attachments: 1}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if r := mujica.MailList(c, &s, mails, mujica.MailListOptions{}); r.Starred != "" {
		mails[0].Starred = !mails[0].Starred
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`ListState[string]` 保存滚动与选择；`MailItem` 包含必填 ID、`From`、`Subject`、`Snippet`、`Time`、`Unread`、`Starred`、`Attachments`。`Height` 默认 420 DIP，必须 ≥ 0；`Disabled` 默认 false。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

`Changed()` 表示选择变化；`Opened` 返回打开的邮件 ID；`Starred` 返回请求切换星标的 ID。调用方更新数据。

## 键盘操作

上下方向键、Home / End 移动选择；Enter 或双击打开；S 请求星标；Cmd/Ctrl+A 全选；Cmd/Ctrl 点击切换单项，Shift 点击选择范围。

## 限制

星标、已读更新、邮件下载与打开由调用方处理；列表只构建可见行。
