# 078 Statistic · 统计数值

## 用途

突出显示一个数值，附带标题、单位、前后缀、精度与变化值说明。组件只格式化调用方给出的数值，不计算业务指标。

## 最小示例

```go
change := 12.5
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Statistic(c, 12840, data.StatisticOptions{
		Title: "Pilgrims", Unit: "souls", Precision: 0,
		Change: &change, ChangeUnit: "%", ChangeNote: "vs last feast",
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Statistic(c *ui.Context, value float64, opts StatisticOptions) ui.Element`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 空 | 数值上方的标题。 |
| `Precision` | `int` | 0 | 小数位数；负数 panic。 |
| `Prefix`, `Suffix` | `string` | 空 | 紧贴数字的前后缀，如 `£`、`k`。 |
| `Unit` | `string` | 空 | 数字后的较小单位文字。 |
| `Format` | `func(float64) string` | 千分位分组 | 自定义数字格式。 |
| `Change` | `*float64` | nil | 变化值：正数显示上箭头与 Success 色，负数下箭头与 Danger 色，零为持平；同时带正负号。 |
| `ChangeUnit` | `string` | 空 | 变化值单位，如 `%`。 |
| `ChangeNote` | `string` | 空 | 变化值说明，如“较上次”。 |

## 状态

仅展示：上升、下降、持平（图标、符号与颜色共同表示）。

## 事件

无。

## 键盘操作

不获取焦点；分组名称把标题、数值、单位与变化一起读出。

## 限制

- 默认格式使用逗号分组与点号小数；其他地区格式通过 `Format` 提供。
