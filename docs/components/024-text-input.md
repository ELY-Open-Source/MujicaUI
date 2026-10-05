# 024 TextInput 文本输入框

## 用途

单行文本编辑。支持占位、前后插槽、清空、只读与错误。编辑核心是 MyGo `ui.TextInputBase`，因此保留选择、撤销重做、剪贴板、右键菜单与输入法组合。

## 最小示例

```go
r := mujica.TextInput(c, &name, mujica.TextInputOptions{Label: "Name", Placeholder: "Your name", Clearable: true})
if r.Submitted() {
	fmt.Println("hello", name)
}
```

签名：`func TextInput(c *ui.Context, value *string, opts TextInputOptions) InputResult`

`InputResult` 的 `Element` 字段是外层元素，`Input` 字段是内部编辑元素，`Changed()`、`Submitted()` 报告事件。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*string` | 必填 | 编辑的字符串；nil 时 panic。 |
| `Placeholder` | `string` | `""` | 空值时的占位文字（TextMuted 色）。 |
| `Label` | `string` | `""` | 无障碍名称，也供测试查找。 |
| `Leading` / `Trailing` | `func()` | `nil` | 在框内左侧、右侧构建内容，如图标或单位。 |
| `Clearable` | `bool` | `false` | 有内容时显示清空按钮；Escape 也会清空。 |
| `Disabled` | `bool` | `false` | 禁用，不可聚焦。 |
| `ReadOnly` | `bool` | `false` | 可聚焦、选择、复制，不可修改。 |
| `Error` | `string` | `""` | 非空时边框变为 Danger，并在下方显示该文字；辅助技术同时朗读。 |

## 状态

默认（Surface 底、ControlBorder 边）、悬停（边框加深）、聚焦（边框 Focus 色；键盘聚焦另加焦点环）、只读（Background 底、Border 边）、错误、禁用。禁用来自自身或任一祖先容器；聚焦时被禁用会放弃焦点，已排队的按键不再修改值、不触发事件。颜色变化 120 毫秒过渡，减少动画时直接切换。

## 事件

- `Changed()`：输入、粘贴、撤销或清空改变了值。只读时不会报告。
- `Submitted()`：在框内按 Enter。

## 键盘操作

Tab 聚焦（全选内容）；Enter 提交；Escape 清空（仅 `Clearable` 且有内容时，否则交给外层）；Cmd/Ctrl+Z、Shift+Cmd+Z 撤销重做；Cmd/Ctrl+A/C/X/V；组合中的文本不会写入值。

## 限制

- 输入法组合期间无法阻止 Enter 提交：MyGo 不公开组合状态，组合文本存在未导出的 `editor.compose`（`ui/editor.go:73`），带编辑器的元素收不到 `InputCompose` 事件（`ui/input.go:689`）。公开接口也无法间接得知：`Element.TextCaret`（`ui/handler.go:82`）只是自绘控件声明输入位置的设置方法，`Tester.TextCaret` 返回的 `Active`（`ui/headless.go:370-373`）只表示可编辑元素持有键盘，不表示正在组合。各平台输入法在组合时通常自己消费 Enter。
- MyGo 没有公开的只读输入接口（只读模式仅用于 `Text.Selectable`，`editor.go:77`）。`ReadOnly` 通过让编辑器编辑一个副本并在同帧还原实现，撤销栈里可能出现被还原的步骤。
- 清空按钮不进入 Tab 序列（与 MyGo SearchField 相同），键盘用 Escape 清空。
