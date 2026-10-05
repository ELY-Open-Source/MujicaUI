# 128 ChartExport 图表导出

## 用途

以当前窗口主题离屏绘制图表，返回 PNG 字节，供调用方保存或传输。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data, err := chart.ChartExport(c, func(c *ui.Context) {
		chart.SparkBar(c, []float64{3, -2, 5}, chart.SparkBarOptions{Frame: chart.ChartFrame{Label: "Alms"}})
	}, chart.ChartExportOptions{Width: 300, Height: 100, Scale: 2})
	if err != nil {
		panic(err)
	}
	bitmap, err := ui.DecodeBitmap(data)
	if err != nil {
		panic(err)
	}
	ui.Image(c, bitmap).Size(300, 100)
}
ui.Render(view, 400, 160, 1)
```

## 参数

`chart func(*ui.Context)` 必填；`Width, Height int` 为正 DIP 尺寸；有限 `Scale float32` 默认 2，负值 panic。

## 状态

沿用当前主题、密度、本地化、自定义 token 和桌面偏好。导出过程为同步渲染。

## 事件

返回 `([]byte, error)`；PNG 编码错误直接返回调用方。

## 键盘操作

导出触发按钮由调用方提供，支持 Enter 和 Space。

## 限制

图表回调应独立渲染，复用调用方数据；调用方在用户触发导出时调用以控制开销。组件只返回字节。
