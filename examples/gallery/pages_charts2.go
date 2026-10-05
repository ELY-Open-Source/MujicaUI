package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

const c2Category = "Charts II"

var charts2Pages = []page{
	{ID: "121", Name: "ChartLegend", ZhName: "图例", Category: c2Category, Demo: c2LegendDemo, Code: c2ChartLegendCode, Keys: "Tab enters; arrows move; Space toggles; Enter expands."},
	{ID: "122", Name: "ChartTooltip", ZhName: "读数浮层", Category: c2Category, Demo: c2TooltipDemo, Code: c2ChartTooltipCode, Keys: c1Keys},
	{ID: "123", Name: "ChartCrosshair", ZhName: "十字光标", Category: c2Category, Demo: c2CrosshairDemo, Code: c2ChartCrosshairCode, Keys: c1Keys},
	{ID: "124", Name: "ChartAxis", ZhName: "坐标轴", Category: c2Category, Demo: c2AxisDemo, Code: c2ChartAxisCode, Keys: c1Keys},
	{ID: "125", Name: "ChartGrid", ZhName: "图表网格", Category: c2Category, Demo: c2GridDemo, Code: c2ChartGridCode, Keys: c1Keys},
	{ID: "126", Name: "ChartBrush", ZhName: "区间刷选", Category: c2Category, Demo: c2BrushDemo, Code: c2ChartBrushCode, Keys: "Tab: start, window, end; arrows step; PageUp/Down step ten; Home/End reach bounds."},
	{ID: "127", Name: "ChartAnnotation", ZhName: "图表标注", Category: c2Category, Demo: c2AnnotationDemo, Code: c2ChartAnnotationCode, Keys: c1Keys + " Alt+PageUp/Down selects a movable note; Alt+arrows move it."},
	{ID: "128", Name: "ChartExport", ZhName: "图表导出", Category: c2Category, Demo: c2ExportDemo, Code: c2ChartExportCode, Keys: "Tab reaches Export PNG; Enter or Space renders bytes in the current theme."},
	{ID: "129", Name: "ViolinPlot", ZhName: "小提琴图", Category: c2Category, Demo: c2ViolinDemo, Code: c2ViolinPlotCode, Keys: c1Keys},
	{ID: "130", Name: "ParallelCoordinates", ZhName: "平行坐标", Category: c2Category, Demo: c2ParallelDemo, Code: c2ParallelCoordinatesCode, Keys: c1Keys + " Alt+PageUp/Down selects an axis; Alt+Up/Down adjusts its lower bound; Alt+Shift+Up/Down adjusts its upper bound; Alt+Delete clears."},
	{ID: "131", Name: "ChordDiagram", ZhName: "弦图", Category: c2Category, Demo: c2ChordDemo, Code: c2ChordDiagramCode, Keys: c1Keys},
	{ID: "132", Name: "SparkBar", ZhName: "迷你柱", Category: c2Category, Demo: c2SparkBarDemo, Code: c2SparkBarCode, Keys: "Tab focuses bars; Left/Right read points; Home/End jump; Escape restores the last value."},
	{ID: "133", Name: "TrendIndicator", ZhName: "趋势指示", Category: c2Category, Demo: c2TrendDemo, Code: c2TrendIndicatorCode, Keys: "Change and Invert controls use the keyboard; the indicator is read as text."},
	{ID: "138", Name: "DensityPlot", ZhName: "密度图", Category: c2Category, Demo: c2DensityDemo, Code: c2DensityPlotCode, Keys: c1Keys + " Space on legend entries hides or shows a group."},
	{ID: "140", Name: "ChartEmptyState", ZhName: "图表占位", Category: c2Category, Demo: c2EmptyDemo, Code: c2ChartEmptyStateCode, Keys: "Tab reaches Retry; Enter or Space retries."},
}

func c2Check(c *ui.Context, value *bool, label string, log func(string)) {
	if mujica.Checkbox(c, value, label, mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprintf("%s = %v", label, *value))
	}
}

func c2Number(c *ui.Context, value *float64, label string, lo, hi, step float64, log func(string)) {
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, label)
		r := mujica.NumberInput(c, value, mujica.NumberInputOptions{Label: label, Min: lo, Max: hi, Step: step})
		r.Element.Width(150)
		if r.Changed() {
			log(fmt.Sprintf("%s = %g", label, *value))
		}
	})
}

func c2Select(c *ui.Context, value *string, label string, choices []string, log func(string)) {
	items := make([]mujica.SelectOption[string], len(choices))
	for i, s := range choices {
		items[i] = mujica.SelectOption[string]{Value: s, Label: s}
	}
	r := mujica.Select(c, value, items, mujica.SelectOptions{Label: label})
	r.Element.Width(160)
	if r.Changed() {
		log(label + " = " + *value)
	}
}

func c2Frame(c *ui.Context, label string, status *string, log func(string)) (mujica.ChartFrame, bool) {
	c2Select(c, status, "Status", c1Statuses, log)
	f := mujica.ChartFrame{Label: label}
	if *status == "Loading" {
		f.Status = mujica.DataLoading
	}
	if *status == "Failed" {
		f.Status, f.Error = mujica.DataFailed, "The archive did not answer"
	}
	return f, *status != "Empty"
}
