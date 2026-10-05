# 215 MessageEditor 消息编辑器

## 用途

编辑已发送的消息：多行文本框加取消与保存。保存时调用方据草稿产生新分支（配合 BranchNavigator）。

## 最小示例

```go
if chat.MessageEditor(c, &draft, chat.MessageEditorOptions{Original: "Hello"}).Saved() {
	fmt.Println("new branch from", draft)
}
```

签名：`func MessageEditor(c *ui.Context, draft *string, opts MessageEditorOptions) MessageEditorResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `draft` | `*string` | 必填 | 正在编辑的文本。 |
| `Original` | `string` | "" | 原文；草稿与原文相同或只有空白时不能保存。 |

## 状态

保存按钮在内容无效时禁用；文本框保留输入法、撤销与剪贴板行为。

## 事件

`Saved()`：点击保存或 Cmd/Ctrl+Enter。`Canceled()`：点击取消或在文本框中按 Escape。

## 键盘操作

Cmd/Ctrl+Enter 保存，Escape 取消；Tab 依次到文本框、取消、保存。

## 限制

- 不保存历史：分支与原文由调用方保存。
