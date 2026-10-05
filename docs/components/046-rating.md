# 046 Rating

评分。

## 用途

以星形显示和设置评分，支持半星等步长、只读展示和数值说明。

## 最小示例

```go
stars := 3.5
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Rating(c, &stars, mujica.RatingOptions{Step: 0.5, ShowValue: true})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*float64` | 必填 | 评分，0 到 Max。 |
| `Max` | `int` | 5 | 星数；负数 panic。 |
| `Step` | `float64` | 1 | 最小变化量，须整除 1（如 0.5、0.25）；否则 panic。 |
| `ReadOnly` | `bool` | false | 只显示，可聚焦以便读出数值。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | 内置“Rating/评分” | 无障碍名称。 |
| `ShowValue` | `bool` | false | 在星形右侧显示“3.5 / 5”。 |

## 状态

默认、悬停（预览指针处的数值）、按下拖动、聚焦（焦点环）、只读、禁用。

## 事件

`ChoiceResult.Changed()` 用户设置了新值；同值点击不触发。

## 键盘操作

→/↑ 增加一个步长，←/↓ 减少，Home 清零，End 满分。

## 限制

- 屏幕阅读器以滑块角色读出名称，数值通过说明“3.5 of 5”提供；MyGo 未公开为自定义元素设置数值范围的接口。
