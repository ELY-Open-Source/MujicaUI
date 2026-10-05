package main

// Snippets are the bodies of the Examples in example_charts1_test.go.

const c1LineCode = `series := []mujica.ChartSeries{
	{Name: "Pews", Values: []float64{3, 5, math.NaN(), 4, 8}},
	{Name: "Candles", Values: []float64{5, 2, 6, 7, 1}},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.LineChart(c, series, mujica.LineChartOptions{Frame: mujica.ChartFrame{Label: "Chapel"}})
}
ui.Render(view, 400, 300, 1)`

const c1AreaCode = `series := []mujica.ChartSeries{
	{Name: "Matins", Values: []float64{3, 4, 6, 5}},
	{Name: "Vespers", Values: []float64{2, 3, 2, 4}},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.AreaChart(c, series, mujica.AreaChartOptions{Frame: mujica.ChartFrame{Label: "Attendance"}, Stacked: true})
}
ui.Render(view, 400, 300, 1)`

const c1BarCode = `categories := []string{"Nave", "Choir", "Crypt"}
series := []mujica.ChartSeries{{Name: "Gold", Values: []float64{4, -2, 3}}, {Name: "Silver", Values: []float64{2, 5, 1}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.BarChart(c, categories, series, mujica.BarChartOptions{Frame: mujica.ChartFrame{Label: "Tithes"}, Labels: true})
}
ui.Render(view, 400, 300, 1)`

const c1PieCode = `slices := []mujica.ChartSlice{{Label: "Wax", Value: 50}, {Label: "Oil", Value: 30}, {Label: "Salt", Value: 4}, {Label: "Ash", Value: 3}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PieChart(c, slices, mujica.PieChartOptions{Frame: mujica.ChartFrame{Label: "Stores"}, Donut: true, MinAngle: 20})
}
ui.Render(view, 400, 300, 1)`

const c1ScatterCode = `var selected []mujica.ScatterKey
groups := []mujica.ScatterGroup{{Name: "Monks", Points: []mujica.ScatterPoint{{X: 1, Y: 2, Size: 4}, {X: 3, Y: 5, Size: 9}}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ScatterChart(c, &selected, groups, mujica.ScatterChartOptions{Frame: mujica.ChartFrame{Label: "Choir"}, Bubbles: true})
}
ui.Render(view, 400, 300, 1)`

const c1RadarCode = `axes := []string{"Voice", "Latin", "Chant", "Organ"}
series := []mujica.ChartSeries{{Name: "Agnès", Values: []float64{8, 6, 9, 4}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.RadarChart(c, axes, series, mujica.RadarChartOptions{Frame: mujica.ChartFrame{Label: "Skills"}})
}
ui.Render(view, 400, 300, 1)`

const c1HeatCode = `rows, cols := []string{"Mon", "Tue"}, []string{"Matins", "Lauds", "Vespers"}
values := [][]float64{{1, 4, 9}, {3, math.NaN(), 7}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.HeatmapChart(c, rows, cols, values, mujica.HeatmapChartOptions{Frame: mujica.ChartFrame{Label: "Hours"}})
}
ui.Render(view, 400, 300, 1)`

const c1CalendarCode = `year := 2026
values := map[mujica.Date]float64{{Year: 2026, Month: time.March, Day: 4}: 5}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.CalendarHeatmap(c, &year, values, mujica.CalendarHeatmapOptions{Frame: mujica.ChartFrame{Label: "Prayers"}})
}
ui.Render(view, 400, 300, 1)`

const c1HistogramCode = `values := []float64{1, 2, 2, 3, 4, 5, 9, 10}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Histogram(c, values, mujica.HistogramOptions{Frame: mujica.ChartFrame{Label: "Ages"}, BinWidth: 5})
}
ui.Render(view, 400, 300, 1)`

const c1BoxCode = `groups := []mujica.BoxGroup{{Name: "Choir", Values: []float64{1, 2, 3, 4, 5, 40}}, {Name: "Nave", Values: []float64{4, 5, 6, 7}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.BoxPlot(c, groups, mujica.BoxPlotOptions{Frame: mujica.ChartFrame{Label: "Hymns"}})
}
ui.Render(view, 400, 300, 1)`

const c1WaterfallCode = `steps := []mujica.WaterfallStep{
	{Label: "Tithes", Value: 12},
	{Label: "Repairs", Value: -5},
	{Label: "Year", Kind: mujica.WaterfallTotal},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.WaterfallChart(c, steps, mujica.WaterfallChartOptions{Frame: mujica.ChartFrame{Label: "Ledger"}})
}
ui.Render(view, 400, 300, 1)`

const c1FunnelCode = `stages := []mujica.ChartSlice{{Label: "Visitors", Value: 1000}, {Label: "Novices", Value: 400}, {Label: "Vows", Value: 100}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FunnelChart(c, stages, mujica.FunnelChartOptions{Frame: mujica.ChartFrame{Label: "Order"}})
}
ui.Render(view, 400, 300, 1)`

const c1BulletCode = `target := 250.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.BulletChart(c, 270, mujica.BulletChartOptions{Frame: mujica.ChartFrame{Label: "Alms"}, Target: &target, Bands: []float64{150, 225, 300}})
}
ui.Render(view, 400, 300, 1)`

const c1SparkCode = `values := []float64{3, 5, 4, 7, 6, 9}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Sparkline(c, values, mujica.SparklineOptions{Frame: mujica.ChartFrame{Label: "Bells"}, Kind: mujica.SparklineArea})
}
ui.Render(view, 400, 300, 1)`

const c1GaugeCode = `segments := []mujica.GaugeSegment{{To: 60, Label: "Calm", Tone: mujica.GaugeGood}, {To: 100, Label: "Crowded", Tone: mujica.GaugeDanger}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Gauge(c, 72, mujica.GaugeOptions{Frame: mujica.ChartFrame{Label: "Nave"}, Max: 100, Unit: "%", Segments: segments})
}
ui.Render(view, 400, 300, 1)`

const c1TreemapCode = `var path []string
root := mujica.ChartNode{Name: "Abbey", Children: []mujica.ChartNode{
	{Name: "Library", Children: []mujica.ChartNode{{Name: "Psalters", Value: 30}, {Name: "Maps", Value: 10}}},
	{Name: "Kitchen", Value: 25},
}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Treemap(c, &path, root, mujica.TreemapOptions{Frame: mujica.ChartFrame{Label: "Rooms"}})
}
ui.Render(view, 400, 300, 1)`

const c1SunburstCode = `var path []string
root := mujica.ChartNode{Name: "Abbey", Children: []mujica.ChartNode{
	{Name: "Library", Children: []mujica.ChartNode{{Name: "Psalters", Value: 30}, {Name: "Maps", Value: 10}}},
	{Name: "Kitchen", Value: 25},
}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Sunburst(c, &path, root, mujica.SunburstOptions{Frame: mujica.ChartFrame{Label: "Rooms"}})
}
ui.Render(view, 400, 300, 1)`

const c1SankeyCode = `nodes := []mujica.SankeyNode{{ID: "tithe", Label: "Tithes"}, {ID: "fund", Label: "Fund"}, {ID: "roof", Label: "Roof"}}
links := []mujica.SankeyLink{{From: "tithe", To: "fund", Value: 60}, {From: "fund", To: "roof", Value: 50}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.SankeyChart(c, nodes, links, mujica.SankeyChartOptions{Frame: mujica.ChartFrame{Label: "Money"}})
}
ui.Render(view, 400, 300, 1)`

const c1NetworkCode = `var state mujica.NetworkGraphState
nodes := []mujica.GraphNode{{ID: "a", Label: "Abbot"}, {ID: "b", Label: "Prior"}, {ID: "c", Label: "Novice"}}
edges := []mujica.GraphEdge{{From: "a", To: "b"}, {From: "b", To: "c"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.NetworkGraph(c, &state, nodes, edges, mujica.NetworkGraphOptions{Frame: mujica.ChartFrame{Label: "House"}, Force: true})
}
ui.Render(view, 400, 300, 1)`

const c1RealtimeCode = `buf := mujica.NewRealtimeBuffer(600, "Bell")
start := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
for i := range 60 {
	buf.Append(start.Add(time.Duration(i)*time.Second), float64(i%7))
}
var state mujica.RealtimeChartState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.RealtimeChart(c, &state, buf, mujica.RealtimeChartOptions{Frame: mujica.ChartFrame{Label: "Sound"}, Window: 30 * time.Second})
}
ui.Render(view, 400, 300, 1)`
