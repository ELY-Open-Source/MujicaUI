# 449 ResourceGauge · 资源仪表

## 用途

展示 CPU、内存与磁盘使用量和阈值状态。

## 最小示例

```go
readings := []devtools.ResourceReading{{Name: "CPU", Used: 45, Total: 100}, {Name: "Memory", Used: 6.4, Total: 8, Unit: "GB"}, {Name: "Disk", Used: 470, Total: 500, Unit: "GB"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if devtools.ResourceGauge(c, readings, devtools.ResourceGaugeOptions{}).Changed() {
		fmt.Println("retry")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Warn 默认 0.7，Critical 默认 0.9，须 0 < Warn < Critical < 1；Height 默认 150 DIP；Status、Error 控制加载与失败；Name 必填，Total 为有限正数，Used 在 0 到 Total 内。

## 状态

正常、偏高、危险、空资源、加载、失败、父级禁用。

## 事件

Changed() 表示请求重试；调用方更新 readings。

## 键盘操作

Tab 聚焦仪表；方向键读取数值与阈值文字。

## 限制

资源采样与容量换算由调用方完成。
