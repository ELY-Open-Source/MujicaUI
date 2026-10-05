# 435 WeekPicker · 周选择器

## 用途

选择整周：悬停时预览整周，值为该周第一天。

## 最小示例

```go
var week core.Date
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if datetime.WeekPicker(c, &week, datetime.WeekPickerOptions{WeekStart: datetime.WeekStartMonday, Clearable: true}).Changed() {
		fmt.Println("week of", week)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func WeekPicker(c *ui.Context, value *Date, opts WeekPickerOptions) ChoiceResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 周起始日：跟随语言或周日、周一、周六。 |
| `Min / Max` | `Date` | 零值 | 所选周必须落在范围内。 |
| `Today` | `Date` | 零值（系统时钟） | 标记为今天的日期。 |
| `Clearable` | `bool` | false | 显示清空按钮。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

关闭、打开（整周预览带、已选周）、禁用。

## 事件

`Changed()` 报告选中周变化。

## 键盘操作

- Enter、Space 或 ↓ 打开；方向键移动；Enter 选择所在周；Escape 关闭。

## 限制

- 值不是周起始日时 panic；更换周起始日后由调用方重设值。
