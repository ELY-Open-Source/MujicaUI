package main

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
)

func c1Series(gap bool) []chart.ChartSeries {
	a := []float64{3, 5, 4, 6, 8, 7, 9, 11}
	if gap {
		a = []float64{3, 5, math.NaN(), 6, 8, 7, 9, 11}
	}
	return []chart.ChartSeries{
		{Name: "Matins", Values: a},
		{Name: "Vespers", Values: []float64{2, 3, 5, 4, 6, 9, 8, 7}},
		{Name: "Compline", Values: []float64{1, 2, 2, 3, 2, 4, 5, 4}},
	}
}

func c1LineDemo(c *ui.Context, log func(string)) {
	type state struct {
		status    string
		gap, time bool
		big       bool
	}
	s := ui.Local(c.Root(), "101", func() state { return state{status: "Ready", gap: true} })
	c1Section(c, "LineChart", func() {
		ui.Checkbox(c, &s.gap, "Missing value")
		ui.Checkbox(c, &s.time, "Time axis")
		ui.Checkbox(c, &s.big, "100,000 points")
	}, func() {
		frame, show := c1Frame(c, "Attendance", &s.status)
		series, opts := c1Series(s.gap), chart.LineChartOptions{Frame: frame}
		if s.big {
			series = []chart.ChartSeries{{Name: "Wave", Values: c1Wave()}}
		}
		if s.time {
			t0 := time.Date(2026, 10, 5, 6, 0, 0, 0, time.UTC)
			opts.Times = make([]time.Time, len(series[0].Values))
			for i := range opts.Times {
				opts.Times[i] = t0.Add(time.Duration(i) * time.Hour)
			}
		}
		if !show {
			series = nil
			opts.Times = nil
		}
		c1Log(log, "LineChart", chart.LineChart(c, series, opts))
	})
}

var c1WaveData []float64

// c1Wave is 100,000 points, made once.
func c1Wave() []float64 {
	if c1WaveData == nil {
		c1WaveData = make([]float64, 100_000)
		for i := range c1WaveData {
			c1WaveData[i] = math.Sin(float64(i)/900) + 0.3*math.Sin(float64(i)/37)
		}
	}
	return c1WaveData
}

func c1AreaDemo(c *ui.Context, log func(string)) {
	type state struct {
		status  string
		stacked bool
		base    float64
	}
	s := ui.Local(c.Root(), "102", func() state { return state{status: "Ready", stacked: true} })
	c1Section(c, "AreaChart", func() {
		ui.Checkbox(c, &s.stacked, "Stacked")
		ui.Text(c, fmt.Sprintf("Baseline %.0f", s.base))
		ui.StepSlider(c, &s.base, -4, 4, 1).Width(140)
	}, func() {
		frame, show := c1Frame(c, "Attendance", &s.status)
		series := c1Series(false)
		if !show {
			series = nil
		}
		c1Log(log, "AreaChart", chart.AreaChart(c, series, chart.AreaChartOptions{Frame: frame, Stacked: s.stacked, Baseline: s.base}))
	})
}

func c1BarDemo(c *ui.Context, log func(string)) {
	type state struct {
		status                      string
		horizontal, stacked, labels bool
	}
	s := ui.Local(c.Root(), "103", func() state { return state{status: "Ready", labels: true} })
	c1Section(c, "BarChart", func() {
		ui.Checkbox(c, &s.horizontal, "Horizontal")
		ui.Checkbox(c, &s.stacked, "Stacked")
		ui.Checkbox(c, &s.labels, "Labels")
	}, func() {
		frame, show := c1Frame(c, "Tithes", &s.status)
		cats := []string{"Nave", "Choir", "Crypt", "Bells"}
		series := []chart.ChartSeries{{Name: "Gold", Values: []float64{4, -2, 3, 5}}, {Name: "Silver", Values: []float64{2, 5, 1, -1}}}
		if !show {
			series = nil
		}
		c1Log(log, "BarChart", chart.BarChart(c, cats, series, chart.BarChartOptions{Frame: frame, Horizontal: s.horizontal, Stacked: s.stacked, Labels: s.labels}))
	})
}

func c1PieDemo(c *ui.Context, log func(string)) {
	type state struct {
		status string
		donut  bool
		min    float64
	}
	s := ui.Local(c.Root(), "104", func() state { return state{status: "Ready", donut: true, min: 20} })
	c1Section(c, "PieChart", func() {
		ui.Checkbox(c, &s.donut, "Donut")
		ui.Text(c, fmt.Sprintf("Merge under %.0f°", s.min))
		ui.StepSlider(c, &s.min, 0, 40, 5).Width(140)
	}, func() {
		frame, show := c1Frame(c, "Stores", &s.status)
		slices := []chart.ChartSlice{{Label: "Wax", Value: 50}, {Label: "Oil", Value: 30}, {Label: "Myrrh", Value: 12}, {Label: "Salt", Value: 4}, {Label: "Ash", Value: 4}}
		if !show {
			slices = nil
		}
		c1Log(log, "PieChart", chart.PieChart(c, slices, chart.PieChartOptions{Frame: frame, Donut: s.donut, MinAngle: s.min}))
	})
}

func c1ScatterDemo(c *ui.Context, log func(string)) {
	type state struct {
		status   string
		bubbles  bool
		selected []chart.ScatterKey
	}
	s := ui.Local(c.Root(), "105", func() state { return state{status: "Ready", bubbles: true} })
	c1Section(c, "ScatterChart", func() {
		ui.Checkbox(c, &s.bubbles, "Bubbles")
		ui.Text(c, fmt.Sprintf("%d selected", len(s.selected)))
	}, func() {
		frame, show := c1Frame(c, "Choir", &s.status)
		groups := []chart.ScatterGroup{
			{Name: "Monks", Points: []chart.ScatterPoint{{X: 1, Y: 2, Size: 4}, {X: 3, Y: 5, Size: 9}, {X: 6, Y: 4, Size: 1}, {X: 7, Y: 8, Size: 6}}},
			{Name: "Nuns", Points: []chart.ScatterPoint{{X: 2, Y: 7, Size: 16}, {X: 8, Y: 9, Size: 2}, {X: 5, Y: 2, Size: 8}}},
		}
		if !show {
			groups = nil
		}
		r := chart.ScatterChart(c, &s.selected, groups, chart.ScatterChartOptions{Frame: frame, Bubbles: s.bubbles})
		c1Log(log, "ScatterChart", r)
		if r.ValueChanged() {
			log(fmt.Sprintf("ScatterChart: %d selected", len(s.selected)))
		}
	})
}

func c1RadarDemo(c *ui.Context, log func(string)) {
	type state struct{ status string }
	s := ui.Local(c.Root(), "106", func() state { return state{status: "Ready"} })
	c1Section(c, "RadarChart", func() {}, func() {
		frame, show := c1Frame(c, "Skills", &s.status)
		series := []chart.ChartSeries{{Name: "Agnès", Values: []float64{8, 6, 9, 4, 7}}, {Name: "Clare", Values: []float64{5, 9, 6, 8, 3}}}
		if !show {
			series = nil
		}
		c1Log(log, "RadarChart", chart.RadarChart(c, []string{"Voice", "Latin", "Chant", "Organ", "Bells"}, series, chart.RadarChartOptions{Frame: frame}))
	})
}

func c1HeatDemo(c *ui.Context, log func(string)) {
	type state struct {
		status string
		gaps   bool
	}
	s := ui.Local(c.Root(), "107", func() state { return state{status: "Ready", gaps: true} })
	c1Section(c, "HeatmapChart", func() { ui.Checkbox(c, &s.gaps, "Empty cells") }, func() {
		frame, show := c1Frame(c, "Hours", &s.status)
		rows, cols := []string{"Mon", "Tue", "Wed", "Thu"}, []string{"Matins", "Lauds", "Terce", "Vespers", "Compline"}
		values := make([][]float64, len(rows))
		for r := range rows {
			values[r] = make([]float64, len(cols))
			for col := range cols {
				values[r][col] = float64((r*7 + col*5) % 11)
				if s.gaps && (r+col)%6 == 1 {
					values[r][col] = math.NaN()
				}
			}
		}
		if !show {
			rows, values = nil, nil
		}
		c1Log(log, "HeatmapChart", chart.HeatmapChart(c, rows, cols, values, chart.HeatmapChartOptions{Frame: frame}))
	})
}
