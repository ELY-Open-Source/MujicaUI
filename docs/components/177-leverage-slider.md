# 177 LeverageSlider

杠杆滑块。

## 用途

在预设杠杆档位间选择草稿值，显示高风险与极高风险区间提示，按“确认”后才写入值。

## 最小示例

```go
leverage := 5.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if finance.LeverageSlider(c, &leverage, finance.LeverageSliderOptions{Warning: 10, Danger: 50}).Changed() {
		fmt.Println("leverage", leverage)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*float64` | 必填 | 当前杠杆，须为 Levels 之一。 |
| `Levels` | `[]float64` | 1–100 八档 | 递增档位，至少两档。 |
| `Warning / Danger` | `float64` | 0 | 高风险与极高风险起点；Danger 须大于 Warning。 |
| `Label` | `string` | 内置“Leverage” | 名称。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、聚焦、拖动、高风险、极高风险、草稿未确认、禁用。

## 事件

`ChoiceResult.Changed()` 在确认写入时为 true。

## 键盘操作

Tab 聚焦滑块：←/↓ 与 →/↑ 逐档、Home/End 到两端；Tab 到达确认按钮。

## 限制

- 滑块命中区高 44 DIP；风险区间以色带、图标和文字共同表达。
