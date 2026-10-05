# 233 ParameterPanel 参数面板

## 用途

温度、Top-p、最大长度等参数；范围与默认值；重置。

## 最小示例

```go
temperature := 1.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	params := []mujica.ModelParam{{Label: "Temperature", Value: &temperature, Min: 0, Max: 2, Step: 0.1, Default: 1}}
	mujica.ParameterPanel(c, params)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`params` 每项 `Label`、`Value` 必填；`Min < Max`、`Step > 0`；`Default` 与当前值须为有限范围内数值，且位于 Min 起算的 Step 网格；范围与默认值显示在滑轨下。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

ChoiceResult.Changed() 表示参数或默认值重置改变了值。

## 键盘操作

Tab 到滑块与重置；方向键按 Step 调整，Home/End 到范围端点；Enter/Space 重置。

## 限制

参数校验失败 panic，数字步进复用 Slider 与 NumberInput 的步长网格。
