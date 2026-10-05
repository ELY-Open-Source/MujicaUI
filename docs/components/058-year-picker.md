# 058 YearPicker · 年份选择器

## 用途

按十年分页浏览并选择年份。

## 最小示例

```go
year := 2026
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	opts := datetime.YearPickerOptions{Label: "Vintage", Min: 1990, Max: 2030}
	if datetime.YearPicker(c, &year, opts).Changed() {
		fmt.Println("year", year)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func YearPicker(c *ui.Context, value *int, opts YearPickerOptions) DateTimeResult`，返回值的 `Element` 字段是组件元素。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Min` | `int` | 0 | 最早年份；0 不限制。 |
| `Max` | `int` | 0 | 最晚年份；0 不限制。 |
| `DisableYear` | `func(int) bool` | nil | 不可选年份。 |
| `Clearable` | `bool` | false | 有值时显示清空按钮。 |
| `Today` | `Date` | 零值（系统时钟） | 标记当前年份。 |
| `ReadOnly` | `bool` | false | 只读：可聚焦、可复制（Cmd/Ctrl+C），不能修改，不打开弹层。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | 内置“选择年份” | 辅助技术名称。 |

## 状态

默认、悬停、聚焦、选中（酒红填充 + 内嵌细框）、当前年（加粗 + 菱形）、不可选、相邻十年（弱化）、禁用。

## 事件

`Changed()`：选择或清空导致值变化时为 true。`*value` 为 0 表示未选择。

## 键盘操作

- Enter/Space/↓：打开。
- ←/→：前后一年；↑/↓：前后三年。
- PageUp/PageDown：上/下一个十年。
- Home/End：页首/页尾。
- Enter：选择并关闭；Escape：关闭。

## 限制

- 只支持公元 1 年及以后；值为负数时 panic。
- 日期与时间值不带时区：`Date`、`YearMonth`、`TimeOfDay` 只保存日历字段，组件内部计算不做 UTC 转换。
- `Today` 为零值时读取 `c.Now()`（系统本地时区）；需要固定“今天”（测试、跨时区应用）时显式传入。
- 无障碍：MyGo v0.2.10 没有公开设置弹出按钮展开状态（`expanded`）、数值范围（`hasRange`/`accRange`）与值（`accValue`）的接口，这些字段只在 MyGo 内部组件中赋值（ui/element.go:286-295，例如 ui/base.go:297、ui/indicators.go:61）。因此触发器只能声明为 `RolePopUpButton` 而无法报告展开状态；时间段以 `Description` 读出当前值与范围，错误通过 `Error()` 关联。
