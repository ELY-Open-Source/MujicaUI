# 069 Card · 卡片

## 用途

把一组相关内容放在带描边的表面上：标题、内容与页脚三个插槽。默认不增加点击行为；可选 `Decorated` 绘制宫廷风格的双线边框与四角菱形角饰。Card 是视觉回归基准之一，快照覆盖浅深主题与三种密度。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Card(c, data.CardOptions{Title: "Archive", Decorated: true, Footer: func() {
		input.Button(c, "Open", input.ButtonOptions{})
	}}, func() {
		ui.Text(c, "Letters of the court, sealed in wax.")
	}).Width(280)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Card(c *ui.Context, opts CardOptions, body func()) ui.Element`，`body` 构建内容插槽。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 空 | 衬线标题；为空时不显示标题区。同时作为分组的无障碍名称。 |
| `Footer` | `func()` | nil | 页脚插槽，位于分隔线下方并右对齐，常放操作按钮。 |
| `Decorated` | `bool` | false | 外框与内框双线、四角菱形与标题下菱形分隔线，均使用 Ornament 色。 |

## 状态

默认（Surface 表面 + Border 描边，4 DIP 圆角）、装饰（Ornament 双线与角饰）。卡片本身无悬停、按下、选中或禁用状态；内部控件保留各自状态。

## 事件

无。卡片不产生事件，交互来自内部控件。

## 键盘操作

卡片不获取焦点；Tab 按顺序进入内部控件。装饰只做绘制，不占焦点、不遮挡内容、不缩小点击区域。

## 限制

- 卡片不提供整卡点击；需要时由调用方在内容中放置按钮或链接。
- 角饰为矢量绘制，位于 4 DIP 内框的四角；装饰模式会额外增加 6 DIP 内边距。
