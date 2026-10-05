# 030 RadioGroup 单选组

## 用途

在一组选项中选一个，横向或纵向排列；选项值是任意可比较类型，作为稳定标识。基于 `ui.RadioGroup` 与 `ui.RadioBase`，整组是一个 Tab 停靠点。

## 最小示例

```go
items := []mujica.RadioItem[string]{{Value: "s", Label: "Small"}, {Value: "m", Label: "Medium"}, {Value: "l", Label: "Large"}}
if mujica.RadioGroup(c, &size, items, mujica.RadioGroupOptions{Label: "Size"}).Changed() {
	fmt.Println("size", size)
}
```

签名：`func RadioGroup[T comparable](c *ui.Context, value *T, items []RadioItem[T], opts RadioGroupOptions) InputResult`

## 参数

`RadioItem[T]`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Value` | `T` | 必填 | 稳定选项值；重复时 panic。 |
| `Label` | `string` | `""` | 标签，点击标签也会选择。 |
| `Disabled` | `bool` | `false` | 单项禁用，方向键跳过。 |

`RadioGroupOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Horizontal` | `bool` | `false` | 横向排列。 |
| `Label` | `string` | `""` | 组的无障碍名称。 |
| `Disabled` | `bool` | `false` | 禁用整组。 |

`value` 为 nil 或 `items` 为空时 panic。`*value` 不在选项中时没有选中项。

## 状态

未选（空心圆）、悬停与按下（AccentText 边）、选中（酒红实心圆与象牙白圆点；悬停与按下时底色加深、边框变为正文色）、键盘聚焦、禁用。

## 事件

`Changed()`：用户选择了另一个值（点击、Space 或方向键）。

## 键盘操作

Tab 进入组并落在已选项；Up/Down/Left/Right 移动并直接选择，循环，跳过禁用项；Home、End 选择首末可用项。

## 限制

- 每个选项按 `Value` 相等性取得稳定键（首次出现的序号），选项重新排序后焦点仍跟随原来的值；组件记住见过的值，长期不断变化的选项集合会让这份记录增长。

- 方向键对横向与纵向组都同时响应上下左右（MyGo `RadioGroup` 固定使用 `Horizontal|Vertical`）。
