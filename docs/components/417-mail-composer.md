# 417 MailComposer 邮件编辑器

## 用途

编辑收件人、抄送、密送、主题、正文和附件，将草稿交给调用方保存。

## 最小示例

```go
draft := mujica.MailDraft{To: []string{"uika@example.com"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.MailComposer(c, &draft, mujica.MailComposerOptions{
		OnSaveDraft: func(d mujica.MailDraft) {}, // store d
	})
	if r.Sent {
		draft = mujica.MailDraft{}
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`MailDraft` 保存 `To`、`Cc`、`Bcc`、`Subject`、`Body`、`Attachments`。`Contacts` 提供联系人。`OnSaveDraft func(MailDraft)` 必填，收到独立的切片副本；`OnAttach func()` 可选；`Disabled` 默认 false。

## 状态

载入含抄送、密送地址的草稿时自动显示对应字段；点击地址字段标签会聚焦编辑器。

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

`Sent` 表示发送请求，`Discarded` 表示舍弃请求。编辑后静置两秒调用 OnSaveDraft；手动保存立即回调。

## 键盘操作

Tab 访问字段与操作；Cmd/Ctrl+S 保存草稿；组件内 Cmd/Ctrl+Enter 发送；正文 Enter 换行。

## 限制

邮箱语法由 net/mail.ParseAddress 校验；实际发送、附件选择与存储由调用方处理；保存回调在界面线程运行。
