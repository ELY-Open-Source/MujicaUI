package main

const c2ChartLegendCode = `hidden := map[string]bool{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	k := core.Tokens(c)
	items := []chart.ChartLegendItem{{Name: "Pews", Color: k.Info}, {Name: "Candles", Color: k.Success, Shape: 1}}
	chart.ChartLegend(c, &hidden, items, chart.ChartLegendOptions{MaxItems: 1})
}
ui.Render(view, 400, 200, 1)`

const c2ChartTooltipCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	k := core.Tokens(c)
	chart.ChartTooltip(c, 40, 40, []chart.ChartTooltipLine{{Text: "Lauds"}, {Text: "Pews 5", Color: k.Info}, {Text: "Candles 2", Color: k.Success, Shape: 1}})
}
ui.Render(view, 400, 200, 1)`

const c2ChartCrosshairCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	x := chart.ChartAxis(c, chart.ChartAxisOptions{Max: 10})
	y := chart.ChartAxis(c, chart.ChartAxisOptions{Side: chart.SideLeft, Max: 100})
	hair := chart.ChartCrosshair(c, chart.ChartCrosshairOptions{X: &x, Y: &y})
	ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		plot := ui.Rect{X: r.X + 60, Y: r.Y + 20, W: r.W - 80, H: r.H - 60}
		x.Draw(p, plot)
		y.Draw(p, plot)
		hair.Draw(p, plot, x.Map(plot, 5), y.Map(plot, 70))
	})
}
ui.Render(view, 400, 300, 1)`

const c2ChartAxisCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	axis := chart.ChartAxis(c, chart.ChartAxisOptions{Side: chart.SideLeft, Kind: chart.AxisLog, Min: 1, Max: 1000})
	ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		axis.Draw(p, ui.Rect{X: r.X + 60, Y: r.Y + 20, W: r.W - 80, H: r.H - 60})
	})
}
ui.Render(view, 400, 300, 1)`

const c2ChartGridCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	x := chart.ChartAxis(c, chart.ChartAxisOptions{Max: 10})
	y := chart.ChartAxis(c, chart.ChartAxisOptions{Side: chart.SideLeft, Max: 100})
	grid := chart.ChartGrid(c, chart.ChartGridOptions{X: &x, Y: &y, Minor: 2})
	ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		plot := ui.Rect{X: r.X + 60, Y: r.Y + 20, W: r.W - 80, H: r.H - 60}
		grid.Draw(p, plot)
		x.Draw(p, plot)
		y.Draw(p, plot)
	})
}
ui.Render(view, 400, 300, 1)`

const c2ChartBrushCode = `value := [2]float64{20, 60}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.ChartBrush(c, &value, chart.ChartBrushOptions{Label: "Months", Max: 100, Step: 1, Values: []float64{2, 4, 3, 8, 5}})
}
ui.Render(view, 400, 160, 1)`

const c2ChartAnnotationCode = `notes := []chart.ChartNote{{Kind: chart.NoteHLine, Value: 50, Text: "Goal", Draggable: true}, {Kind: chart.NoteBand, Value: 6, End: 8, Text: "Lent"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	x := chart.ChartAxis(c, chart.ChartAxisOptions{Max: 10})
	y := chart.ChartAxis(c, chart.ChartAxisOptions{Side: chart.SideLeft, Max: 100})
	annotation := chart.ChartAnnotation(c, notes, chart.ChartAnnotationOptions{X: &x, Y: &y})
	box := ui.Box(c).Fill().Focusable().Label("Annotations")
	b := box.Bounds()
	annotation.Drag(box, ui.Rect{X: 60, Y: 20, W: b.W - 80, H: b.H - 60})
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		plot := ui.Rect{X: r.X + 60, Y: r.Y + 20, W: r.W - 80, H: r.H - 60}
		x.Draw(p, plot)
		y.Draw(p, plot)
		annotation.Draw(p, plot)
	})
}
ui.Render(view, 400, 300, 1)`

const c2ChartExportCode = `view := func(c *ui.Context) {
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
ui.Render(view, 400, 160, 1)`

const c2SparkBarCode = `values := []float64{3, -2, 4, 1, -3}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.SparkBar(c, values, chart.SparkBarOptions{Frame: chart.ChartFrame{Label: "Alms"}, BarWidth: 10})
}
ui.Render(view, 400, 100, 1)`

const c2TrendIndicatorCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.TrendIndicator(c, 12.5, chart.TrendIndicatorOptions{Precision: 1, Unit: "%", Note: "vs Lent"})
}
ui.Render(view, 400, 100, 1)`

const c2ChartEmptyStateCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.ChartEmptyState(c, data.DataFailed, chart.ChartEmptyStateOptions{Height: 180, Error: "The archive did not answer"})
}
ui.Render(view, 400, 240, 1)`

const c2ViolinPlotCode = `groups := []chart.BoxGroup{{Name: "Choir", Values: []float64{1, 2, 3, 4, 7, 9}}, {Name: "Nave", Values: []float64{2, 4, 5, 6, 8}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.ViolinPlot(c, groups, chart.ViolinPlotOptions{Frame: chart.ChartFrame{Label: "Hymns"}, Bandwidth: 1})
}
ui.Render(view, 400, 300, 1)`

const c2ParallelCoordinatesCode = `filters := map[string][2]float64{}
axes := []string{"Age", "Hymns", "Vigils"}
rows := [][]float64{{20, 3, 9}, {35, 8, 2}, {50, 5, 5}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.ParallelCoordinates(c, &filters, axes, rows, chart.ParallelCoordinatesOptions{Frame: chart.ChartFrame{Label: "Monks"}, Names: []string{"Agnes", "Clare", "Hild"}})
}
ui.Render(view, 400, 300, 1)`

const c2ChordDiagramCode = `names := []string{"Nave", "Choir", "Crypt"}
matrix := [][]float64{{0, 5, 2}, {4, 0, 1}, {1, 2, 0}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.ChordDiagram(c, names, matrix, chart.ChordDiagramOptions{Frame: chart.ChartFrame{Label: "Paths"}})
}
ui.Render(view, 400, 300, 1)`

const c2DensityPlotCode = `groups := []chart.BoxGroup{{Name: "Matins", Values: []float64{1, 2, 2, 3, 7}}, {Name: "Vespers", Values: []float64{3, 4, 5, 6, 6}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.DensityPlot(c, groups, chart.DensityPlotOptions{Frame: chart.ChartFrame{Label: "Ages"}, Bandwidth: 1, Points: 128})
}
ui.Render(view, 400, 300, 1)`
