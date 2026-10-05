package main

const c2ChartLegendCode = `hidden := map[string]bool{}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	k := mujica.Tokens(c)
	items := []mujica.ChartLegendItem{{Name: "Pews", Color: k.Info}, {Name: "Candles", Color: k.Success, Shape: 1}}
	mujica.ChartLegend(c, &hidden, items, mujica.ChartLegendOptions{MaxItems: 1})
}
ui.Render(view, 400, 200, 1)`

const c2ChartTooltipCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	k := mujica.Tokens(c)
	mujica.ChartTooltip(c, 40, 40, []mujica.ChartTooltipLine{{Text: "Lauds"}, {Text: "Pews 5", Color: k.Info}, {Text: "Candles 2", Color: k.Success, Shape: 1}})
}
ui.Render(view, 400, 200, 1)`

const c2ChartCrosshairCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	x := mujica.ChartAxis(c, mujica.ChartAxisOptions{Max: 10})
	y := mujica.ChartAxis(c, mujica.ChartAxisOptions{Side: mujica.SideLeft, Max: 100})
	hair := mujica.ChartCrosshair(c, mujica.ChartCrosshairOptions{X: &x, Y: &y})
	ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		plot := ui.Rect{X: r.X + 60, Y: r.Y + 20, W: r.W - 80, H: r.H - 60}
		x.Draw(p, plot)
		y.Draw(p, plot)
		hair.Draw(p, plot, x.Map(plot, 5), y.Map(plot, 70))
	})
}
ui.Render(view, 400, 300, 1)`

const c2ChartAxisCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	axis := mujica.ChartAxis(c, mujica.ChartAxisOptions{Side: mujica.SideLeft, Kind: mujica.AxisLog, Min: 1, Max: 1000})
	ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		axis.Draw(p, ui.Rect{X: r.X + 60, Y: r.Y + 20, W: r.W - 80, H: r.H - 60})
	})
}
ui.Render(view, 400, 300, 1)`

const c2ChartGridCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	x := mujica.ChartAxis(c, mujica.ChartAxisOptions{Max: 10})
	y := mujica.ChartAxis(c, mujica.ChartAxisOptions{Side: mujica.SideLeft, Max: 100})
	grid := mujica.ChartGrid(c, mujica.ChartGridOptions{X: &x, Y: &y, Minor: 2})
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
	mujica.Use(c, mujica.Settings{})
	mujica.ChartBrush(c, &value, mujica.ChartBrushOptions{Label: "Months", Max: 100, Step: 1, Values: []float64{2, 4, 3, 8, 5}})
}
ui.Render(view, 400, 160, 1)`

const c2ChartAnnotationCode = `notes := []mujica.ChartNote{{Kind: mujica.NoteHLine, Value: 50, Text: "Goal", Draggable: true}, {Kind: mujica.NoteBand, Value: 6, End: 8, Text: "Lent"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	x := mujica.ChartAxis(c, mujica.ChartAxisOptions{Max: 10})
	y := mujica.ChartAxis(c, mujica.ChartAxisOptions{Side: mujica.SideLeft, Max: 100})
	annotation := mujica.ChartAnnotation(c, notes, mujica.ChartAnnotationOptions{X: &x, Y: &y})
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
	mujica.Use(c, mujica.Settings{})
	data, err := mujica.ChartExport(c, func(c *ui.Context) {
		mujica.SparkBar(c, []float64{3, -2, 5}, mujica.SparkBarOptions{Frame: mujica.ChartFrame{Label: "Alms"}})
	}, mujica.ChartExportOptions{Width: 300, Height: 100, Scale: 2})
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
	mujica.Use(c, mujica.Settings{})
	mujica.SparkBar(c, values, mujica.SparkBarOptions{Frame: mujica.ChartFrame{Label: "Alms"}, BarWidth: 10})
}
ui.Render(view, 400, 100, 1)`

const c2TrendIndicatorCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TrendIndicator(c, 12.5, mujica.TrendIndicatorOptions{Precision: 1, Unit: "%", Note: "vs Lent"})
}
ui.Render(view, 400, 100, 1)`

const c2ChartEmptyStateCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ChartEmptyState(c, mujica.DataFailed, mujica.ChartEmptyStateOptions{Height: 180, Error: "The archive did not answer"})
}
ui.Render(view, 400, 240, 1)`

const c2ViolinPlotCode = `groups := []mujica.BoxGroup{{Name: "Choir", Values: []float64{1, 2, 3, 4, 7, 9}}, {Name: "Nave", Values: []float64{2, 4, 5, 6, 8}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ViolinPlot(c, groups, mujica.ViolinPlotOptions{Frame: mujica.ChartFrame{Label: "Hymns"}, Bandwidth: 1})
}
ui.Render(view, 400, 300, 1)`

const c2ParallelCoordinatesCode = `filters := map[string][2]float64{}
axes := []string{"Age", "Hymns", "Vigils"}
rows := [][]float64{{20, 3, 9}, {35, 8, 2}, {50, 5, 5}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ParallelCoordinates(c, &filters, axes, rows, mujica.ParallelCoordinatesOptions{Frame: mujica.ChartFrame{Label: "Monks"}, Names: []string{"Agnes", "Clare", "Hild"}})
}
ui.Render(view, 400, 300, 1)`

const c2ChordDiagramCode = `names := []string{"Nave", "Choir", "Crypt"}
matrix := [][]float64{{0, 5, 2}, {4, 0, 1}, {1, 2, 0}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ChordDiagram(c, names, matrix, mujica.ChordDiagramOptions{Frame: mujica.ChartFrame{Label: "Paths"}})
}
ui.Render(view, 400, 300, 1)`

const c2DensityPlotCode = `groups := []mujica.BoxGroup{{Name: "Matins", Values: []float64{1, 2, 2, 3, 7}}, {Name: "Vespers", Values: []float64{3, 4, 5, 6, 6}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.DensityPlot(c, groups, mujica.DensityPlotOptions{Frame: mujica.ChartFrame{Label: "Ages"}, Bandwidth: 1, Points: 128})
}
ui.Render(view, 400, 300, 1)`
