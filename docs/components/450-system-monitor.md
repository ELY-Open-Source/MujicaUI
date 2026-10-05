# 450 SystemMonitor · 系统监视器

## 用途

在同一时间轴上组合多资源折线，展示最新采样值。

## 最小示例

```go
times := []time.Time{time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC), time.Date(2026, 10, 5, 9, 0, 5, 0, time.UTC)}
metrics := []devtools.MonitorMetric{{Name: "CPU", Unit: "%", Samples: []float64{42, 51}}, {Name: "Memory", Unit: "GB", Samples: []float64{3.2, 3.4}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if devtools.SystemMonitor(c, times, metrics, devtools.SystemMonitorOptions{}).Changed() {
		fmt.Println("retry")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Height 为单图高度，默认 110 DIP；Status、Error 控制加载与失败；每个 MonitorMetric.Name 必填，Samples 与 times 等长。

## 状态

实时采样、空资源、加载、失败、父级禁用。

## 事件

Changed() 表示请求重试。

## 键盘操作

Tab 到图例及图；方向键读取采样点；Home、End 到边界。

## 限制

采样由调用方提供；更新采样提交新切片；NaN 表示缺测。
