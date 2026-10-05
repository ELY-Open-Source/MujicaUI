package main

// Snippets are the bodies of the Examples in example_charts1_test.go.

const c1LineCode = `series := []chart.ChartSeries{
	{Name: "Pews", Values: []float64{3, 5, math.NaN(), 4, 8}},
	{Name: "Candles", Values: []float64{5, 2, 6, 7, 1}},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.LineChart(c, series, chart.LineChartOptions{Frame: chart.ChartFrame{Label: "Chapel"}})
}
ui.Render(view, 400, 300, 1)`

const c1AreaCode = `series := []chart.ChartSeries{
	{Name: "Matins", Values: []float64{3, 4, 6, 5}},
	{Name: "Vespers", Values: []float64{2, 3, 2, 4}},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.AreaChart(c, series, chart.AreaChartOptions{Frame: chart.ChartFrame{Label: "Attendance"}, Stacked: true})
}
ui.Render(view, 400, 300, 1)`

const c1BarCode = `categories := []string{"Nave", "Choir", "Crypt"}
series := []chart.ChartSeries{{Name: "Gold", Values: []float64{4, -2, 3}}, {Name: "Silver", Values: []float64{2, 5, 1}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.BarChart(c, categories, series, chart.BarChartOptions{Frame: chart.ChartFrame{Label: "Tithes"}, Labels: true})
}
ui.Render(view, 400, 300, 1)`

const c1PieCode = `slices := []chart.ChartSlice{{Label: "Wax", Value: 50}, {Label: "Oil", Value: 30}, {Label: "Salt", Value: 4}, {Label: "Ash", Value: 3}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.PieChart(c, slices, chart.PieChartOptions{Frame: chart.ChartFrame{Label: "Stores"}, Donut: true, MinAngle: 20})
}
ui.Render(view, 400, 300, 1)`

const c1ScatterCode = `var selected []chart.ScatterKey
groups := []chart.ScatterGroup{{Name: "Monks", Points: []chart.ScatterPoint{{X: 1, Y: 2, Size: 4}, {X: 3, Y: 5, Size: 9}}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.ScatterChart(c, &selected, groups, chart.ScatterChartOptions{Frame: chart.ChartFrame{Label: "Choir"}, Bubbles: true})
}
ui.Render(view, 400, 300, 1)`

const c1RadarCode = `axes := []string{"Voice", "Latin", "Chant", "Organ"}
series := []chart.ChartSeries{{Name: "Agnès", Values: []float64{8, 6, 9, 4}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.RadarChart(c, axes, series, chart.RadarChartOptions{Frame: chart.ChartFrame{Label: "Skills"}})
}
ui.Render(view, 400, 300, 1)`

const c1HeatCode = `rows, cols := []string{"Mon", "Tue"}, []string{"Matins", "Lauds", "Vespers"}
values := [][]float64{{1, 4, 9}, {3, math.NaN(), 7}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.HeatmapChart(c, rows, cols, values, chart.HeatmapChartOptions{Frame: chart.ChartFrame{Label: "Hours"}})
}
ui.Render(view, 400, 300, 1)`

const c1CalendarCode = `year := 2026
values := map[core.Date]float64{{Year: 2026, Month: time.March, Day: 4}: 5}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.CalendarHeatmap(c, &year, values, chart.CalendarHeatmapOptions{Frame: chart.ChartFrame{Label: "Prayers"}})
}
ui.Render(view, 400, 300, 1)`

const c1HistogramCode = `values := []float64{1, 2, 2, 3, 4, 5, 9, 10}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.Histogram(c, values, chart.HistogramOptions{Frame: chart.ChartFrame{Label: "Ages"}, BinWidth: 5})
}
ui.Render(view, 400, 300, 1)`

const c1BoxCode = `groups := []chart.BoxGroup{{Name: "Choir", Values: []float64{1, 2, 3, 4, 5, 40}}, {Name: "Nave", Values: []float64{4, 5, 6, 7}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.BoxPlot(c, groups, chart.BoxPlotOptions{Frame: chart.ChartFrame{Label: "Hymns"}})
}
ui.Render(view, 400, 300, 1)`

const c1WaterfallCode = `steps := []chart.WaterfallStep{
	{Label: "Tithes", Value: 12},
	{Label: "Repairs", Value: -5},
	{Label: "Year", Kind: chart.WaterfallTotal},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.WaterfallChart(c, steps, chart.WaterfallChartOptions{Frame: chart.ChartFrame{Label: "Ledger"}})
}
ui.Render(view, 400, 300, 1)`

const c1FunnelCode = `stages := []chart.ChartSlice{{Label: "Visitors", Value: 1000}, {Label: "Novices", Value: 400}, {Label: "Vows", Value: 100}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.FunnelChart(c, stages, chart.FunnelChartOptions{Frame: chart.ChartFrame{Label: "Order"}})
}
ui.Render(view, 400, 300, 1)`

const c1BulletCode = `target := 250.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.BulletChart(c, 270, chart.BulletChartOptions{Frame: chart.ChartFrame{Label: "Alms"}, Target: &target, Bands: []float64{150, 225, 300}})
}
ui.Render(view, 400, 300, 1)`

const c1SparkCode = `values := []float64{3, 5, 4, 7, 6, 9}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.Sparkline(c, values, chart.SparklineOptions{Frame: chart.ChartFrame{Label: "Bells"}, Kind: chart.SparklineArea})
}
ui.Render(view, 400, 300, 1)`

const c1GaugeCode = `segments := []chart.GaugeSegment{{To: 60, Label: "Calm", Tone: chart.GaugeGood}, {To: 100, Label: "Crowded", Tone: chart.GaugeDanger}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.Gauge(c, 72, chart.GaugeOptions{Frame: chart.ChartFrame{Label: "Nave"}, Max: 100, Unit: "%", Segments: segments})
}
ui.Render(view, 400, 300, 1)`

const c1TreemapCode = `var path []string
root := chart.ChartNode{Name: "Abbey", Children: []chart.ChartNode{
	{Name: "Library", Children: []chart.ChartNode{{Name: "Psalters", Value: 30}, {Name: "Maps", Value: 10}}},
	{Name: "Kitchen", Value: 25},
}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.Treemap(c, &path, root, chart.TreemapOptions{Frame: chart.ChartFrame{Label: "Rooms"}})
}
ui.Render(view, 400, 300, 1)`

const c1SunburstCode = `var path []string
root := chart.ChartNode{Name: "Abbey", Children: []chart.ChartNode{
	{Name: "Library", Children: []chart.ChartNode{{Name: "Psalters", Value: 30}, {Name: "Maps", Value: 10}}},
	{Name: "Kitchen", Value: 25},
}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.Sunburst(c, &path, root, chart.SunburstOptions{Frame: chart.ChartFrame{Label: "Rooms"}})
}
ui.Render(view, 400, 300, 1)`

const c1SankeyCode = `nodes := []chart.SankeyNode{{ID: "tithe", Label: "Tithes"}, {ID: "fund", Label: "Fund"}, {ID: "roof", Label: "Roof"}}
links := []chart.SankeyLink{{From: "tithe", To: "fund", Value: 60}, {From: "fund", To: "roof", Value: 50}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.SankeyChart(c, nodes, links, chart.SankeyChartOptions{Frame: chart.ChartFrame{Label: "Money"}})
}
ui.Render(view, 400, 300, 1)`

const c1NetworkCode = `var state chart.NetworkGraphState
nodes := []chart.GraphNode{{ID: "a", Label: "Abbot"}, {ID: "b", Label: "Prior"}, {ID: "c", Label: "Novice"}}
edges := []chart.GraphEdge{{From: "a", To: "b"}, {From: "b", To: "c"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.NetworkGraph(c, &state, nodes, edges, chart.NetworkGraphOptions{Frame: chart.ChartFrame{Label: "House"}, Force: true})
}
ui.Render(view, 400, 300, 1)`

const c1RealtimeCode = `buf := chart.NewRealtimeBuffer(600, "Bell")
start := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
for i := range 60 {
	buf.Append(start.Add(time.Duration(i)*time.Second), float64(i%7))
}
var state chart.RealtimeChartState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.RealtimeChart(c, &state, buf, chart.RealtimeChartOptions{Frame: chart.ChartFrame{Label: "Sound"}, Window: 30 * time.Second})
}
ui.Render(view, 400, 300, 1)`
