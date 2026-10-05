# 025 TextArea 多行文本框

## 用途

多行文本编辑，可设最小与最大高度，超出最大高度后在框内滚动；可显示字数；提交快捷键可配置。编辑核心是 MyGo `ui.TextAreaBase`，只排版可见段落，适合长文本。

## 最小示例

```go
r := mujica.TextArea(c, &notes, mujica.TextAreaOptions{Label: "Notes", MaxHeight: 160, MaxLength: 200})
if r.Submitted() {
	fmt.Println("send", notes)
}
```

签名：`func TextArea(c *ui.Context, value *string, opts TextAreaOptions) InputResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*string` | 必填 | 编辑的字符串。 |
| `Placeholder` | `string` | `""` | 占位文字。 |
| `Label` | `string` | `""` | 无障碍名称。 |
| `MinHeight` | `float32` | `0`（三行） | 编辑区最小高度，DIP。 |
| `MaxHeight` | `float32` | `0`（随内容增长） | 编辑区最大高度，超出后滚动；小于 `MinHeight` 时 panic。 |
| `MaxLength` | `int` | `0`（不显示） | 显示“当前 / 上限”字数（按 Unicode 字符计）；超出时变红并提示超出数量。 |
| `SubmitKey` / `SubmitMods` | `ui.Key` / `ui.Modifiers` | Cmd+Enter（非 macOS 为 Ctrl+Enter） | 提交快捷键；`SubmitKey` 为零时使用默认值。 |
| `Disabled` | `bool` | `false` | 禁用。 |
| `ReadOnly` | `bool` | `false` | 只读，可选择复制。 |
| `Error` | `string` | `""` | 错误文字；为空且超出 `MaxLength` 时自动显示超出提示。 |

## 状态

默认、悬停、聚焦、只读、错误、禁用，表现同 TextInput。

## 事件

`Changed()`：内容被用户改变。`Submitted()`：按下提交快捷键。

## 键盘操作

Enter 与 Shift+Enter 换行；提交快捷键提交；方向键、PageUp/PageDown 移动光标并滚动；撤销、剪贴板与输入法同 TextInput。

## 限制

- `MaxLength` 只提示，不截断输入，以免打断输入法组合与粘贴。
- 提交快捷键不能是不带修饰键的 Enter：MyGo 编辑器先接收 Enter 并插入换行（`editor.go:161`）。
- 只读实现方式同 TextInput。
