# 082 AlertDialog 警示对话框

## 用途

确认危险操作。默认聚焦“取消”，点击遮罩无效，Escape 等同取消；可要求先输入目标名称，未确认目标前 Enter 不会执行。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.Button(c, "Delete archive", mujica.ButtonOptions{Variant: mujica.Danger}).Clicked() {
		open = true
	}
	r := mujica.AlertDialog(c, &open, mujica.AlertDialogOptions{
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

`AlertDialogResult{Confirmed, Cancelled bool}` 只在发生选择的那一帧为 true，之后对话框关闭。

## 状态

- 确认按钮使用 Danger 变体；`ConfirmText` 未匹配时禁用（MyGo 统一变暗）。
- 输入框有边框与焦点环；内容保留在对话框打开期间，关闭后清空。

## 事件

- `Confirmed`：点击确认，或在输入框内按 Enter 且目标已匹配。
- `Cancelled`：点击取消、在取消上按 Enter / Space、或按 Escape。
- 点击遮罩不产生事件。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| 打开时 | 焦点在“取消”。 |
| Enter | 触发聚焦的按钮；在输入框中只有目标匹配时确认。 |
| Escape | 取消并关闭。 |
| Tab | 在对话框内循环。 |

## 限制

- 取消按钮的焦点由 `AutoFocus` 给出，打开时不显示焦点环（MyGo 只在键盘移动焦点后设置可见焦点，`ui/input.go:868`）；按 Tab 后显示。
- 不替代系统原生警示框；需要系统样式时用 MyGo 的桌面对话框。
