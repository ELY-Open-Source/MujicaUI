# 089 Popconfirm 气泡确认

## 用途

在触发控件旁做轻量确认（归档、移出列表等可撤销操作）。说明与面板关联，确认与取消都可聚焦。高风险操作请用 AlertDialog。

## 最小示例

```go
open := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	b := input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary})
	if b.Clicked() {
		open = true
	}
	if overlay.Popconfirm(c, b, &open, overlay.PopconfirmOptions{Title: "Archive this letter?"}).Confirmed {
		fmt.Println("archived")
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func Popconfirm(c *ui.Context, anchor ui.Element, open *bool, opts PopconfirmOptions) PopconfirmResult`。`anchor` 为 nil 或 `Title` 为空时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 问题，也是面板的无障碍名称。 |
| `Description` | `string` | `""` | 补充说明，作为面板的无障碍描述。 |
| `ConfirmLabel` | `string` | 内置“确定 / Confirm” | 确认按钮文字。 |
| `CancelLabel` | `string` | 内置“取消 / Cancel” | 取消按钮文字。 |
| `Danger` | `bool` | `false` | 确认按钮使用危险变体，图标改为 Danger 色，打开时焦点改在取消。 |

`PopconfirmResult{Confirmed, Canceled bool}` 只在发生选择的那一帧为 true。

## 状态

与 Popover 相同的面板；下方放不下时翻到上方。

## 事件

- `Confirmed`：点击确认或在其上按 Enter / Space。
- `Canceled`：点击取消、Escape、外部点击。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| 打开时 | 焦点在确认按钮；`Danger` 时在取消按钮，连按 Enter 不会执行危险操作。 |
| Enter | 触发聚焦的按钮。 |
| Tab | 在取消与确认间移动，焦点顺序紧跟锚点。 |
| Escape | 取消并把焦点还给锚点。 |

## 限制

- 打开期间外部点击只用于关闭，不传给下方控件：`PopoverBase` 用覆盖整窗的背板接收点击（`ui/base.go:373`）。
- 位置只有下方与上方（`ui/widgets.go:689`）；宽度不超过窗口，说明过长时滚动，按钮始终可见。
