# 257 ConversationExport · 会话导出

## 用途

选择格式（Markdown、JSON、纯文本）与范围，预览后生成导出文字；写文件由调用方完成。

## 最小示例

```go
var state chat.ConversationExportState
msgs := []chat.SessionMessage{{Role: chat.SessionUser, Text: "Hello"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.ConversationExport(c, &state, msgs, chat.ConversationExportOptions{Title: "Greeting"}); v.Exported {
		log.Println(len(v.Content), "bytes to save")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ConversationExport(c *ui.Context, s *ConversationExportState, messages []SessionMessage, opts ConversationExportOptions) ConversationExportView`。另有纯函数 `ExportConversation(title, messages, format) string`。

`ConversationExportState`：`Format`（`ExportMarkdown`、`ExportJSON`、`ExportText`）、`Range`（`ExportAll`、`ExportLast`、`ExportSelected`）、`Last`（0 → 10）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 会话标题，写入导出内容。 |
| `Selected` | `[]int` | nil | `ExportSelected` 使用的消息下标；越界 panic。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

格式单选、范围单选（无选中消息时“已选”不可选）、最近 N 条数字框、等宽字体预览、导出按钮、禁用。

## 事件

`Exported`：导出请求；`Content` 为生成的文字，`Format` 为格式。

## 键盘操作

- Tab 依次经过格式组、范围组、数字框、预览、导出按钮；单选组内方向键选择。

## 限制

- 导出不解析消息文字：Markdown 只加标题与作者行，原文原样写入。
- 预览只在格式、范围或消息变化时重新生成。
