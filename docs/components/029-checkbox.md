# 029 Checkbox 复选框

## 用途

切换布尔值，支持选中、未选、半选；点击方框或标签都会切换。`CheckboxGroup` 提供组内全选（对应 MyGo CheckboxGroup）。基于 `ui.CheckboxBase`。

## 最小示例

```go
input.CheckboxGroup(c, "Notifications", []input.CheckboxItem{
	{Label: "Mail", Value: &mail},
	{Label: "Calendar", Value: &calendar},
}, input.CheckboxGroupOptions{})
```

签名：

```go
func Checkbox(c *ui.Context, value *bool, label string, opts CheckboxOptions) ui.Element
func CheckboxGroup(c *ui.Context, label string, items []CheckboxItem, opts CheckboxGroupOptions) ChoiceResult
```

## 参数

`CheckboxOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Indeterminate` | `bool` | `false` | 显示半选横线；此时点击总是变为选中。是否半选由调用方决定。 |
| `Disabled` | `bool` | `false` | 禁用。 |

`CheckboxItem`：`Label string`、`Value *bool`（必填）、`Disabled bool`。`CheckboxGroupOptions`：`Disabled bool`。`items` 为空时 panic。

## 状态

未选（Surface 底、ControlBorder 边）、悬停与按下（AccentText 边）、选中（酒红底、象牙白勾；悬停与按下时底色加深、边框变为正文色）、半选（酒红底、象牙白横线）、键盘聚焦（方框外焦点环）、禁用。行高不低于密度控件高度。

全选框：全部选中为选中，部分选中为半选，全不选为未选。点击时若可用项未全选则全部勾选，否则全部清除；禁用项保持不变。

## 事件

- `Checkbox(...).Changed()`：用户切换了值。
- `CheckboxGroup(...).Changed()`：全选或任一项被切换。

## 键盘操作

Tab 聚焦（每个复选框一个停靠点）；Space 切换。

## 限制

- MyGo 只有内部字段能把复选框报告为“mixed”（`access.go:307`，`checked` 未导出）。半选时组件在无障碍描述中加入“部分选中”，但角色状态仍报告为未选中。
