# 416 MailReader 邮件阅读器

## 用途

展示邮件头、纯文本或安全 Markdown 子集正文，以及附件操作。

## 最小示例

```go
m := messaging.Mail{From: "Sakiko", Subject: "Second act", Body: "The score is attached.\n> Can you send it?",
	Attachments: []messaging.MailAttachment{{Name: "score.pdf", Size: 2048}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if r := messaging.MailReader(c, m); r.Link != "" {
		c.OpenURL(r.Link)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`Mail.From`、`Subject` 必填；`To`、`Cc`、`Date`、`Body` 为邮件内容。`Format` 默认 `MailPlain`，也支持 `MailMarkdown`。`Attachments []MailAttachment` 提供名称、字节数和 MIME 类型。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

`Link` 返回安全链接；`Opened`、`Downloaded` 返回附件索引，无事件时为 -1。

## 键盘操作

Tab 访问链接、引用开关和附件；Enter / Space 激活；正文支持选择复制。

## 限制

MarkdownView 处理安全 Markdown 子集，HTML 展示为文字；链接和附件动作由调用方执行。
