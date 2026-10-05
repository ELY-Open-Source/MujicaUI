package chart

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// chartBuilders builds every Charts I component with a small data set.
func chartBuilders() map[string]func(c *ui.Context) ChartResult {
	year := 2026
	var path, path2 []string
	var sel []ScatterKey
	var graph NetworkGraphState
	var rt RealtimeChartState
	buf := NewRealtimeBuffer(10, "a")
	buf.Append(rtStart, 1)
	buf.Append(rtStart.Add(time.Second), 2)
	target := 3.0
	f := ChartFrame{Label: "Under test"}
	return map[string]func(c *ui.Context) ChartResult{
		"line": func(c *ui.Context) ChartResult { return LineChart(c, lineSeries(), LineChartOptions{Frame: f}) },
		"area": func(c *ui.Context) ChartResult { return AreaChart(c, lineSeries(), AreaChartOptions{Frame: f}) },
		"bar":  func(c *ui.Context) ChartResult { return BarChart(c, barCats, barSeries(), BarChartOptions{Frame: f}) },
		"pie":  func(c *ui.Context) ChartResult { return PieChart(c, pieSlices, PieChartOptions{Frame: f}) },
		"scat": func(c *ui.Context) ChartResult {
			return ScatterChart(c, &sel, scatterGroups, ScatterChartOptions{Frame: f})
		},
		"radar": func(c *ui.Context) ChartResult {
			return RadarChart(c, radarAxes, radarSeries(), RadarChartOptions{Frame: f})
		},
		"heat": func(c *ui.Context) ChartResult {
			return HeatmapChart(c, heatRows, heatCols, heatValues(), HeatmapChartOptions{Frame: f})
		},
		"cal": func(c *ui.Context) ChartResult {
			return CalendarHeatmap(c, &year, calendarValues(), CalendarHeatmapOptions{Frame: f})
		},
		"hist": func(c *ui.Context) ChartResult { return Histogram(c, []float64{1, 2, 3}, HistogramOptions{Frame: f}) },
		"box":  func(c *ui.Context) ChartResult { return BoxPlot(c, boxGroups, BoxPlotOptions{Frame: f}) },
		"fall": func(c *ui.Context) ChartResult {
			return WaterfallChart(c, waterfallSteps, WaterfallChartOptions{Frame: f})
		},
		"funnel": func(c *ui.Context) ChartResult { return FunnelChart(c, funnelStages, FunnelChartOptions{Frame: f}) },
		"bullet": func(c *ui.Context) ChartResult {
			return BulletChart(c, 2, BulletChartOptions{Frame: f, Target: &target, Bands: []float64{1, 4}})
		},
		"spark": func(c *ui.Context) ChartResult { return Sparkline(c, sparkValues, SparklineOptions{Frame: f}) },
		"gauge": func(c *ui.Context) ChartResult {
			return Gauge(c, 72, GaugeOptions{Frame: f, Max: 100, Segments: gaugeSegments})
		},
		"tree": func(c *ui.Context) ChartResult { return Treemap(c, &path, chartTree(), TreemapOptions{Frame: f}) },
		"sun":  func(c *ui.Context) ChartResult { return Sunburst(c, &path2, chartTree(), SunburstOptions{Frame: f}) },
		"sankey": func(c *ui.Context) ChartResult {
			return SankeyChart(c, sankeyNodes, sankeyLinks, SankeyChartOptions{Frame: f})
		},
		"graph": func(c *ui.Context) ChartResult {
			return NetworkGraph(c, &graph, graphNodes, graphEdges, NetworkGraphOptions{Frame: f})
		},
		"rt": func(c *ui.Context) ChartResult {
			return RealtimeChart(c, &rt, buf, RealtimeChartOptions{Frame: f, Window: time.Minute})
		},
	}
}

func TestChartsIgnoreInputWhenDisabled(t *testing.T) {
	for name, build := range chartBuilders() {
		values := 0
		h := &chartHarness{disabled: true, build: func(c *ui.Context) ChartResult {
			r := build(c)
			if r.ValueChanged() || r.Retried() {
				values++
			}
			if _, ok := r.Toggled(); ok {
				values++
			}
			return r
		}}
		tt := ui.NewTester(h.view, 560, 480)
		for range 6 {
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyRight)
			tt.Key(0, ui.KeyEnter)
			tt.Key(0, ui.KeySpace)
		}
		b := h.res.Element.Bounds()
		tt.Click("Under test")
		tt.ClickAt(b.X+b.W/2, b.Y+b.H/2)
		if h.res.Point != -1 || values != 0 {
			t.Errorf("%s: disabled chart took input: point %d, events %d", name, h.res.Point, values)
		}
	}
}

func TestChartsEnterAtFirstMark(t *testing.T) {
	for name, build := range chartBuilders() {
		h := &chartHarness{build: build}
		tt := ui.NewTester(h.view, 560, 480)
		tabTo(t, tt, "Under test")
		tt.Key(0, ui.KeyRight)
		if h.res.Point != 0 || chartReadout(tt) == msg2(tt) {
			t.Errorf("%s: first arrow read point %d", name, h.res.Point)
		}
	}
}

// msg2 is the hint line, read back to tell it from a readout.
func msg2(*ui.Tester) string { return "Focus the chart and use the arrow keys to read values" }
