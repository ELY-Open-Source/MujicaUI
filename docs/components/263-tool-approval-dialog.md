# 263 ToolApprovalDialog 工具审批对话框

## 用途

工具执行前请求批准：显示工具名、执行理由和参数预览，允许或拒绝。遵循 Dialog 与 AlertDialog 规则。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := agent.ToolApprovalDialog(c, &open, agent.ToolApprovalOptions{
		Tool: "delete_files", Args: "{\"glob\":\"*.tmp\"}", Reason: "Clean the build folder.",
	})
	if r.Allowed {
		fmt.Println("run the tool")
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func ToolApprovalDialog(c *ui.Context, open *bool, opts ToolApprovalOptions) ToolApprovalResult`。`Tool` 为空时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Tool` | `string` | 必填 | 工具名。 |
| `Args` | `string` | `""` | 参数，合法 JSON 会格式化高亮。 |
| `Reason` | `string` | `""` | Agent 给出的理由，作为对话框说明。 |

## 状态

- 打开时焦点在“拒绝”；点击遮罩不做任何事。
- 参数过长时在预览区内滚动，按钮始终可见。

## 事件

- `Allowed`：点击允许，或焦点在允许上按 Enter / Space。
- `Denied`：点击拒绝、焦点在拒绝上按 Enter、或按 Escape。
- `Element`：对话框面板，关闭时为 nil。

## 键盘操作

- 打开时焦点在拒绝，Enter 只触发聚焦按钮，必须 Tab 到“允许”才能批准。
- Escape 等同拒绝；Tab 在对话框内循环；关闭后焦点回到打开者。

## 限制

- 输入法组合期间按 Escape 仍会拒绝（MyGo 不公开组合状态，`ui/editor.go:73`）。
