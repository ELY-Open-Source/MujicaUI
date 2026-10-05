# 082 AlertDialog 警示对话框

## 用途

确认危险操作。默认聚焦“取消”，点击遮罩无效，Escape 等同取消；可要求先输入目标名称，未输入前确认按钮禁用；输入框中的 Enter 从不确认。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.Button(c, "Delete archive", input.ButtonOptions{Variant: input.Danger}).Clicked() {
		open = true
	}
	r := overlay.AlertDialog(c, &open, overlay.AlertDialogOptions{
		Title:       "Delete the archive?",
		Description: "Twelve letters will be destroyed.",
		ConfirmText: "archive",
	})
	if r.Confirmed {
		fmt.Println("deleted")
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func AlertDialog(c *ui.Context, open *bool, opts AlertDialogOptions) AlertDialogResult`。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 标题与无障碍名称；为空时 panic。 |
| `Description` | `string` | `""` | 后果说明，同时作为无障碍描述。 |
| `ConfirmLabel` | `string` | 内置“删除 / Delete” | 危险按钮文字。 |
| `CancelLabel` | `string` | 内置“取消 / Cancel” | 取消按钮文字。 |
| `ConfirmText` | `string` | `""` | 非空时显示输入框，输入内容与之完全相同后确认按钮才可用。 |

`AlertDialogResult{Confirmed, Canceled bool}` 只在发生选择的那一帧为 true，之后对话框关闭。

## 状态

- 确认按钮使用 Danger 变体；`ConfirmText` 未匹配时禁用（MyGo 统一变暗）。
- 输入框有边框与焦点环；内容保留在对话框打开期间，关闭后清空。

## 事件

- `Confirmed`：点击确认，或焦点在已解锁的确认按钮上按 Enter / Space。
- `Canceled`：点击取消、在取消上按 Enter / Space、或按 Escape。
- 点击遮罩不产生事件。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| 打开时 | 焦点在“取消”。 |
| Enter | 触发聚焦的按钮；在输入框中不确认，Tab 到确认按钮后再按。 |
| Escape | 取消并关闭。 |
| Tab | 在对话框内循环。 |

## 限制

- 输入框中的 Enter 从不确认，必须用 Tab 移到确认按钮：MyGo 在 Enter 时先提交输入法组合再处理按键（`ui/editor.go:688`），而组合状态不公开，无法区分“提交组合”的 Enter。
- 输入法组合期间按 Escape 仍会关闭浮层：MyGo 没有公开组合状态（`ui/editor.go:73` 的 `compose` 未导出，组合事件不交给带编辑器元素的 `HandleInput`，`ui/input.go:689`），编辑器也不接收 Escape（`ui/editor.go:163`），macOS 上按键先于输入法送达应用（`internal/darwin/surface.go:704`）。
- 关闭后焦点优先由 MyGo 还给打开者；打开者已消失、焦点落空时，焦点停在组件调用处的不可见占位元素上，Tab 从这里继续。不能回到打开者的父级：MyGo 不公开当前焦点或打开者（`ui/scope.go:123`、`ui/scope.go:137`）。
- 取消按钮的焦点由 `AutoFocus` 给出，打开时不显示焦点环（`ui/input.go:868`）；按 Tab 后显示。
