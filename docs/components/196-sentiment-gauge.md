# 196 SentimentGauge

## 用途

情绪仪表。

## 最小示例

```go
value := 72.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.SentimentGauge(c, value, finance.SentimentGaugeOptions{Label: "Fear & Greed"})
}
ui.Render(view, 960, 440, 1)
```

## 参数

`value` 范围 0—100；`Label` 默认空；`Size` 默认 240 DIP，指定值至少 48 DIP。`Status/Error/Empty` 默认就绪、空文本、false，保留同一尺寸的加载、失败与空占位。

## 状态

默认与数据边界。

## 事件

`Retry` 回调报告错误占位的重试请求；数值由调用方更新。

## 键盘操作

Tab 聚焦仪表；方向键和 Home/End 朗读当前数值与区间文字。

## 限制

区间边界为 20、40、60、80；值域错误触发 panic。
