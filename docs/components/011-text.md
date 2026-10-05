# 011 Text 文本

## 用途

正文文本：自动换行、按行数省略、可选择复制。`EditableText` 是 Text 的可编辑模式，就地修改一段短文本（如文件名），并展示调用方的保存进度与保存结果。

## 最小示例

```go
name := "Coronation.md"
saveErr := ""
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Text(c, "Order ID: 4821-AX", mujica.TextOptions{Selectable: true})
	mujica.Text(c, "A long note that stops after two lines.", mujica.TextOptions{MaxLines: 2})
	if mujica.EditableText(c, &name, mujica.EditableTextOptions{Error: saveErr}).Changed() {
		saveErr = "" // save name here; set saveErr on failure
	}
}
```

## 参数

`func Text(c *ui.Context, s string, opts TextOptions) *ui.Element`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `MaxLines` | `int` | `0` | 超过该行数时以省略号截断；`0` 不截断；`1` 为单行。负数 panic。 |
| `Selectable` | `bool` | `false` | 允许用指针或键盘选择并复制文本。 |

`func EditableText(c *ui.Context, value *string, opts EditableTextOptions) *ui.Element`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Saving` | `bool` | `false` | 调用方正在保存：显示加载标记，并禁止再次编辑。 |
| `Error` | `string` | `""` | 调用方保存失败的说明，显示在文本下方并报告给辅助技术。 |
| `Label` | `string` | `""` | 辅助技术读出的名称。 |

## 状态

- 默认：可编辑文本带虚线下边，提示可以编辑。
- 悬停：背景变为 `SurfaceHover`，尺寸不变。
- 聚焦：2 DIP 焦点环。
- 编辑中：显示 MyGo 文本输入框，酒红边框，选中扩展名之前的部分。
- 保存中（`Saving`）：显示加载标记，阻止进入编辑。
- 错误（`Error`）：虚线变为 `Danger`，下方显示错误说明，保留已确认的值。

## 事件

- `EditableText(...).Changed()`：用户确认了一个新值，`*value` 已更新。调用方在此保存，并通过 `Saving` 与 `Error` 回传结果。
- 取消编辑不产生事件，`*value` 不变。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 聚焦可编辑文本 |
| Enter / 双击 | 开始编辑 |
| Enter | 确认 |
| Escape | 取消并恢复原值 |
| 焦点移走 | 视为确认 |
| Cmd/Ctrl+A、Cmd/Ctrl+C | 在 `Selectable` 文本上全选、复制 |

## 限制

- 编辑模式基于 `ui.EditableText`。MyGo 不报告“取消”事件，组件只能保证取消后值不变，不能通知调用方。
- 编辑模式只支持单行文本。
- 值为空字符串时文本区域没有可见文字；组件保留最小宽度和虚线，保证仍可点击。
- MyGo 无法检测文本是否真的被截断，所以截断时不会自动显示完整内容的提示。
