# 121 ChartLegend 图例

## 用途

展示序列名称与形状；点击或键盘切换序列显隐，支持横纵布局和折叠。

## 最小示例

```go
hidden := map[string]bool{}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	k := mujica.Tokens(c)
	items := []mujica.ChartLegendItem{{Name: "Pews", Color: k.Info}, {Name: "Candles", Color: k.Success, Shape: 1}}
	mujica.ChartLegend(c, &hidden, items, mujica.ChartLegendOptions{MaxItems: 1})
}
ui.Render(view, 400, 200, 1)
```

## 参数

`hidden *map[string]bool`：名称映射到隐藏状态，nil 为静态图例。`items []ChartLegendItem`：唯一非空 `Name`、`Color`、`Shape`（0–4）。`Vertical` 默认 false；`MaxItems` 默认 0，显示全部，负值 panic。

## 状态

可见项显示形状，隐藏项显示空框与删除线；支持祖先禁用。折叠后显示展开按钮。

## 事件

`Toggled()` 返回本帧切换的序列名称和是否发生事件。

## 键盘操作

Tab 进入图例；方向键在项之间移动；Space 切换显隐；展开按钮支持 Enter 和 Space。

## 限制

调用方依据 hidden 控制数据序列。图例颜色来自调用方，应使用主题色；形状提供颜色之外的识别。
