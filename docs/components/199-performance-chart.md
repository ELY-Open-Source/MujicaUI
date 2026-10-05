# 199 PerformanceChart

## 用途

收益曲线。

## 最小示例

```go
values := []float64{100, 110, 95, 120}
benchmark := []float64{100, 105, 108, 112}
times := []time.Time{time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.PerformanceChart(c, values, finance.PerformanceChartOptions{Frame: chart.ChartFrame{Label: "Performance"}, Times: times, Benchmark: benchmark})
	if r.Changed() || r.Retried() {
		fmt.Println(r.Point)
	}
	if name, ok := r.Toggled(); ok {
		fmt.Println(name)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`Version` 在原地修改净值或基准时递增。

`values` 为正净值；`Times` 必填且递增；`Benchmark` nil 隐藏或与净值等长。`Frame.Label` 必填；`Name/BenchmarkName` 默认本地化名称。

## 状态

默认与数据边界；图表的 Frame 支持加载、空、失败占位。

## 事件

if r.Changed() || r.Retried() { 调用方回调(r.Point) }; if name, ok := r.Toggled(); ok { 调用方回调(name) }

## 键盘操作

Tab 聚焦图和图例；方向键读数；Space 切换序列。

## 限制

净值先归一化为累计收益；最大回撤取峰值到谷值；日期与报价由调用方提供。
