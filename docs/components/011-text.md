# 011 Text 文本

## 用途

正文文本：自动换行、按行数省略、可选择复制。`EditableText` 是 Text 的可编辑模式，就地修改一段短文本（如文件名），并展示调用方的保存进度与保存结果。

## 最小示例

```go
name := "Coronation.md"
saveErr := ""
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	display.Text(c, "Order ID: 4821-AX", display.TextOptions{Selectable: true})
	display.Text(c, "A long note that stops after two lines.", display.TextOptions{MaxLines: 2})
	if display.EditableText(c, &name, display.EditableTextOptions{Error: saveErr}).Changed() {
		saveErr = "" // save name here; set saveErr on failure
	}
}
```

## 参数

`func Text(c *ui.Context, s string, opts TextOptions) ui.Element`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `MaxLines` | `int` | `0` | 超过该行数时以省略号截断；`0` 不截断；`1` 为单行。负数 panic。 |
| `Selectable` | `bool` | `false` | 允许用指针或键盘选择并复制文本。 |

`func EditableText(c *ui.Context, value *string, opts EditableTextOptions) EditableTextView`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Saving` | `bool` | `false` | 调用方正在保存：显示加载标记，并禁止再次编辑。 |
| `Error` | `string` | `""` | 调用方保存失败的说明，显示在文本下方并报告给辅助技术。 |
| `Label` | `string` | `""` | 辅助技术读出的名称；同时用在显示状态和编辑输入框上。 |

返回值 `EditableTextView`：`Element`；`Changed()` 表示本帧确认了新值；`Canceled` 表示本帧按 Escape 放弃了编辑。

## 状态

- 默认：可编辑文本带虚线下边，提示可以编辑。
- 悬停：背景变为 `SurfaceHover`，尺寸不变。
- 聚焦：2 DIP 焦点环（可编辑文本与 `Selectable` 文本）。
- 编辑中：显示文本输入框，酒红边框；名称与错误说明附在输入框上。
- 保存中（`Saving`）：显示加载标记，阻止进入编辑。
- 错误（`Error`）：虚线变为 `Danger`，下方显示错误说明，保留已确认的值。

## 事件

- `EditableTextView.Changed()`：用户确认了新值，`*value` 已更新。调用方在此保存，并通过 `Saving` 与 `Error` 回传结果。
- `EditableTextView.Canceled`：用户按 Escape 放弃编辑，`*value` 不变。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 聚焦可编辑文本；`Selectable` 文本也在 Tab 顺序中 |
| Enter / 双击 | 开始编辑 |
| Enter | 确认 |
| Escape | 取消并恢复原值 |
| 焦点移走 | 视为确认 |
| Cmd/Ctrl+A、Cmd/Ctrl+C | 在 `Selectable` 文本上全选、复制 |

## 限制

- 编辑模式基于 `ui.TextInputBase` 构建，而不是 `ui.EditableText`，因为后者内部的输入框无法通过公开接口设置名称和错误说明。代价是开始编辑时不会像 Finder 那样自动选中扩展名之前的部分，光标停在文本末尾（可用 Cmd/Ctrl+A 全选）。
- 编辑模式只支持单行文本。
- 值为空字符串时文本区域没有可见文字；组件保留最小宽度和虚线，保证仍可点击。
- MyGo 无法检测文本是否真的被截断，所以截断时不会自动显示完整内容的提示。
