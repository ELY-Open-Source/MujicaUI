package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/datetime"
)

func c1CalendarDemo(c *ui.Context, log func(string)) {
	type state struct {
		status, start string
		year          int
	}
	s := ui.Local(c.Root(), "108", func() state { return state{status: "Ready", start: "Monday", year: 2026} })
	c1Section(c, "CalendarHeatmap", func() { ui.Select(c, &s.start, []string{"Locale", "Sunday", "Monday", "Saturday"}) }, func() {
		frame, show := c1Frame(c, "Prayers", &s.status)
		frame.Height = 160
		values := map[core.Date]float64{}
		if show {
			d := core.Date{Year: s.year, Month: time.January, Day: 1}
			for i := range 366 {
				if (i*37)%5 != 0 {
					values[d.AddDays(i)] = float64((i * 13) % 9)
				}
			}
		}
		start := map[string]datetime.WeekStart{"Locale": datetime.WeekStartLocale, "Sunday": datetime.WeekStartSunday, "Monday": datetime.WeekStartMonday, "Saturday": datetime.WeekStartSaturday}[s.start]
		r := chart.CalendarHeatmap(c, &s.year, values, chart.CalendarHeatmapOptions{Frame: frame, WeekStart: start})
		c1Log(log, "CalendarHeatmap", r)
		if r.ValueChanged() {
			log("CalendarHeatmap: year " + strconv.Itoa(s.year))
		}
	})
}

func c1HistogramDemo(c *ui.Context, log func(string)) {
	type state struct {
		status     string
		bins       float64
		cumulative bool
		values     []float64
	}
	s := ui.Local(c.Root(), "109", func() state {
		v := make([]float64, 500)
		for i := range v {
			v[i] = 50 + 12*math.Sin(float64(i)*1.7)*math.Cos(float64(i)*0.31)
		}
		return state{status: "Ready", bins: 10, values: v}
	})
	c1Section(c, "Histogram", func() {
		ui.Text(c, fmt.Sprintf("%.0f bins", s.bins))
		ui.StepSlider(c, &s.bins, 3, 30, 1).Width(140)
		ui.Checkbox(c, &s.cumulative, "Cumulative")
	}, func() {
		frame, show := c1Frame(c, "Ages", &s.status)
		values := s.values
		if !show {
			values = nil
		}
		c1Log(log, "Histogram", chart.Histogram(c, values, chart.HistogramOptions{Frame: frame, Bins: int(s.bins), Cumulative: s.cumulative}))
	})
}

func c1BoxDemo(c *ui.Context, log func(string)) {
	type state struct{ status string }
	s := ui.Local(c.Root(), "110", func() state { return state{status: "Ready"} })
	c1Section(c, "BoxPlot", func() {}, func() {
		frame, show := c1Frame(c, "Hymns", &s.status)
		groups := []chart.BoxGroup{
			{Name: "Choir", Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 24}},
			{Name: "Nave", Values: []float64{4, 5, 5, 6, 7, 8}},
			{Name: "Crypt", Values: []float64{12, 14, 15, 18, 22, 25}},
		}
		if !show {
			groups = nil
		}
		c1Log(log, "BoxPlot", chart.BoxPlot(c, groups, chart.BoxPlotOptions{Frame: frame}))
	})
}

func c1WaterfallDemo(c *ui.Context, log func(string)) {
	type state struct{ status string }
	s := ui.Local(c.Root(), "111", func() state { return state{status: "Ready"} })
	c1Section(c, "WaterfallChart", func() {}, func() {
		frame, show := c1Frame(c, "Ledger", &s.status)
		steps := []chart.WaterfallStep{
			{Label: "Tithes", Value: 12}, {Label: "Repairs", Value: -5}, {Label: "Q1", Kind: chart.WaterfallSubtotal},
			{Label: "Alms", Value: 4}, {Label: "Bread", Value: -2}, {Label: "Year", Kind: chart.WaterfallTotal},
		}
		if !show {
			steps = nil
		}
		c1Log(log, "WaterfallChart", chart.WaterfallChart(c, steps, chart.WaterfallChartOptions{Frame: frame}))
	})
}

func c1FunnelDemo(c *ui.Context, log func(string)) {
	type state struct{ status string }
	s := ui.Local(c.Root(), "112", func() state { return state{status: "Ready"} })
	c1Section(c, "FunnelChart", func() {}, func() {
		frame, show := c1Frame(c, "Order", &s.status)
		stages := []chart.ChartSlice{{Label: "Visitors", Value: 1000}, {Label: "Novices", Value: 400}, {Label: "Vows", Value: 100}, {Label: "Priors", Value: 12}}
		if !show {
			stages = nil
		}
		c1Log(log, "FunnelChart", chart.FunnelChart(c, stages, chart.FunnelChartOptions{Frame: frame}))
	})
}

func c1BulletDemo(c *ui.Context, log func(string)) {
	type state struct {
		status        string
		value, target float64
		vertical      bool
	}
	s := ui.Local(c.Root(), "113", func() state { return state{status: "Ready", value: 270, target: 250} })
	c1Section(c, "BulletChart", func() {
		ui.Text(c, fmt.Sprintf("Value %.0f", s.value))
		ui.StepSlider(c, &s.value, 0, 300, 10).Width(140)
		ui.Checkbox(c, &s.vertical, "Vertical")
	}, func() {
		frame, _ := c1Frame(c, "Alms", &s.status)
		if s.vertical {
			frame.Height = 200
		}
		target := s.target
		c1Log(log, "BulletChart", chart.BulletChart(c, s.value, chart.BulletChartOptions{Frame: frame, Target: &target, Bands: []float64{150, 225, 300}, Vertical: s.vertical}))
	})
}

func c1SparkDemo(c *ui.Context, log func(string)) {
	type state struct{ status string }
	s := ui.Local(c.Root(), "114", func() state { return state{status: "Ready"} })
	c1Section(c, "Sparkline", func() {}, func() {
		frame, show := c1Frame(c, "Bells", &s.status)
		values := []float64{3, 5, 4, math.NaN(), 7, 6, 9, 8, 11}
		if !show {
			values = nil
		}
		ui.Row(c).Gap(24).Wrap().Children(func() {
			for _, k := range []chart.SparklineKind{chart.SparklineLine, chart.SparklineBar, chart.SparklineArea} {
				f := frame
				f.Label = fmt.Sprintf("Bells %d", k+1)
				c1Log(log, f.Label, chart.Sparkline(c, values, chart.SparklineOptions{Frame: f, Kind: k}))
			}
		})
	})
}

func c1GaugeDemo(c *ui.Context, log func(string)) {
	type state struct {
		status string
		value  float64
		full   bool
	}
	s := ui.Local(c.Root(), "115", func() state { return state{status: "Ready", value: 72} })
	c1Section(c, "Gauge", func() {
		ui.Text(c, fmt.Sprintf("Value %.0f", s.value))
		ui.StepSlider(c, &s.value, 0, 100, 1).Width(140)
		ui.Checkbox(c, &s.full, "Full circle")
	}, func() {
		frame, _ := c1Frame(c, "Nave", &s.status)
		segments := []chart.GaugeSegment{{To: 60, Label: "Calm", Tone: chart.GaugeGood}, {To: 85, Label: "Busy", Tone: chart.GaugeWarning}, {To: 100, Label: "Crowded", Tone: chart.GaugeDanger}}
		c1Log(log, "Gauge", chart.Gauge(c, s.value, chart.GaugeOptions{Frame: frame, Max: 100, Unit: "%", Segments: segments, Full: s.full}))
	})
}
