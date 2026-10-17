# 204 MessageHeader 消息头

## 用途

在消息上方显示模型名、时间与状态（发送中、生成中、已发送、失败），状态同时用图标和文字表达。

## 最小示例

```go
chat.MessageHeader(c, chat.MessageHeaderOptions{Model: "herald-2", Time: time.Now(), Status: chat.MessageSent})
```

签名：`func MessageHeader(c *ui.Context, opts MessageHeaderOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Model` | `string` | "" | 模型名，加粗；空则不显示。 |
| `Time` | `time.Time` | 零值 | 时间，按 `15:04` 显示；零值不显示。 |
| `Status` | `MessageStatus` | `MessageNoStatus` | `MessageSending`、`MessageStreaming`、`MessageSent`、`MessageFailed`；未知值 panic。 |

## 状态

生成中显示菱形加载标记（减少动画时静止）；失败用 Danger 色并带警示图标。

## 事件

无。

## 键盘操作

不获取焦点。

## 限制

- 时间格式固定为 24 小时制时分，日期由 MessageList 的日期分隔显示。
