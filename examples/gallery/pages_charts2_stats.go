package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
)

func c2LegendDemo(c *ui.Context, log func(string)) {
	type state struct {
		hidden             map[string]bool
		vertical, disabled bool
		max                float64
	}
	s := ui.Local(c.Root(), "121", func() state { return state{hidden: map[string]bool{}, max: 3} })
	c1Section(c, "ChartLegend", func() {
		c2Check(c, &s.vertical, "Vertical", log)
		c2Check(c, &s.disabled, "Disabled", log)
		c2Number(c, &s.max, "Max items", 0, 6, 1, log)
	}, func() {
		pal := internal.Palette(core.Tokens(c))
		names := []string{"Matins", "Lauds", "Prime", "Terce", "Sext", "None"}
		items := make([]chart.ChartLegendItem, len(names))
		for i, n := range names {
			items[i] = chart.ChartLegendItem{Name: n, Color: pal[i], Shape: i % 5}
		}
		ui.Column(c).Disabled(s.disabled).Children(func() {
			r := chart.ChartLegend(c, &s.hidden, items, chart.ChartLegendOptions{Vertical: s.vertical, MaxItems: int(s.max)})
			if expanded, ok := r.Expanded(); ok {
				log(fmt.Sprintf("ChartLegend: expanded %v", expanded))
			}
			if n, ok := r.Toggled(); ok {
				log("ChartLegend: toggled " + n)
			}
		})
		ui.Text(c, fmt.Sprintf("Hidden: %v", s.hidden))
	})
}

func c2BrushDemo(c *ui.Context, log func(string)) {
	type state struct {
		value    [2]float64
		disabled bool
		step     float64
	}
	s := ui.Local(c.Root(), "126", func() state { return state{value: [2]float64{20, 60}, step: 1} })
	c1Section(c, "ChartBrush", func() {
		c2Check(c, &s.disabled, "Disabled", log)
		c2Number(c, &s.step, "Step", 1, 10, 1, log)
	}, func() {
		vals := []float64{1, 4, 2, 7, 3, 8, 4, 9}
		r := chart.ChartBrush(c, &s.value, chart.ChartBrushOptions{Label: "Months", Max: 100, Step: s.step, Values: vals, Disabled: s.disabled})
		if r.Changed() {
			log(fmt.Sprintf("ChartBrush: range %g–%g", s.value[0], s.value[1]))
		}
		x := chart.ChartAxis(c, chart.ChartAxisOptions{Min: s.value[0], Max: s.value[1]})
		y := chart.ChartAxis(c, chart.ChartAxisOptions{Side: chart.SideLeft, Max: 10})
		xs := []float64{0, 15, 30, 45, 60, 75, 90, 100}
		var visible []int
		for i, v := range xs {
			if v >= s.value[0] && v <= s.value[1] {
				visible = append(visible, i)
			}
		}
		plotOf := func(r ui.Rect) ui.Rect { return ui.Rect{X: r.X + 40, Y: r.Y + 8, W: max(1, r.W-60), H: max(1, r.H-40)} }
		at := func(r ui.Rect, i int) (float32, float32) {
			plot := plotOf(r)
			return x.Map(plot, xs[visible[i]]), y.Map(plot, vals[visible[i]])
		}
		c1Log(log, "Zoom preview", chart.Chart(c, chart.ChartFrame{Label: "Zoom preview", Height: 140}, chart.ChartSpec{
			Marks: len(visible), Anchor: at,
			Pick: func(r ui.Rect, px, py float32) int {
				return internal.Nearest(len(visible), func(i int) (float32, float32) { return at(r, i) }, px, py, 24, false)
			},
			Readout: func(i int) string { return fmt.Sprintf("x %g · y %g", xs[visible[i]], vals[visible[i]]) },
			Draw: func(p *ui.Painter, r ui.Rect, active int) {
				plot := plotOf(r)
				x.Draw(p, plot)
				y.Draw(p, plot)
				for i := range visible {
					px, py := at(r, i)
					internal.Marker(p, 0, px, py, 8, core.Tokens(c).AccentText)
				}
			},
		}))
	})
}

func c2ExportDemo(c *ui.Context, log func(string)) {
	type state struct {
		width, height, scale float64
		disabled             bool
		bitmap               *ui.Bitmap
		bytes                int
		error                string
	}
	s := ui.Local(c.Root(), "128", func() state { return state{width: 400, height: 240, scale: 2} })
	c1Section(c, "ChartExport", func() {
		c2Number(c, &s.width, "Width", 200, 600, 50, log)
		c2Number(c, &s.height, "Height", 160, 400, 40, log)
		c2Number(c, &s.scale, "Scale", 1, 3, 0.5, log)
		c2Check(c, &s.disabled, "Disabled", log)
	}, func() {
		if input.Button(c, "Export PNG", input.ButtonOptions{Disabled: s.disabled}).Clicked() {
			data, err := chart.ChartExport(c, func(c *ui.Context) {
				chart.BarChart(c, []string{"Nave", "Choir", "Crypt"}, []chart.ChartSeries{{Name: "Pews", Values: []float64{3, 5, 2}}}, chart.BarChartOptions{Frame: chart.ChartFrame{Label: "Exported chart", Height: float32(s.height) - 80}})
			}, chart.ChartExportOptions{Width: int(s.width), Height: int(s.height), Scale: float32(s.scale)})
			if err == nil {
				s.bitmap, err = ui.DecodeBitmap(data)
			}
			if err != nil {
				s.error = err.Error()
				log("ChartExport: " + s.error)
			} else {
				s.error = ""
				s.bytes = len(data)
				log(fmt.Sprintf("ChartExport: %d bytes", s.bytes))
			}
		}
		if s.error != "" {
			ui.Text(c, s.error).TextColor(core.Tokens(c).Danger).Role(ui.RoleStatus)
		}
		if s.bitmap != nil {
			ui.Image(c, s.bitmap).Size(float32(s.width), float32(s.height))
			ui.Text(c, fmt.Sprintf("PNG: %d bytes", s.bytes))
		}
	})
}

func c2SparkBarDemo(c *ui.Context, log func(string)) {
	type state struct {
		status   string
		width    float64
		disabled bool
	}
	s := ui.Local(c.Root(), "132", func() state { return state{status: "Ready", width: 10} })
	c1Section(c, "SparkBar", func() {
		c2Number(c, &s.width, "Bar width", 2, 20, 1, log)
		c2Check(c, &s.disabled, "Disabled", log)
	}, func() {
		f, show := c2Frame(c, "Alms", &s.status, log)
		vals := []float64{3, -2, 4, 1, -3}
		if !show {
			vals = nil
		}
		ui.Column(c).Disabled(s.disabled).Children(func() {
			c1Log(log, "SparkBar", chart.SparkBar(c, vals, chart.SparkBarOptions{Frame: f, BarWidth: float32(s.width)}))
		})
	})
}

func c2TrendDemo(c *ui.Context, log func(string)) {
	type state struct {
		change, precision float64
		invert            bool
	}
	s := ui.Local(c.Root(), "133", func() state { return state{change: 12.5, precision: 1} })
	c1Section(c, "TrendIndicator", func() {
		c2Number(c, &s.change, "Change", -100, 100, 0.5, log)
		c2Number(c, &s.precision, "Precision", 0, 3, 1, log)
		c2Check(c, &s.invert, "Invert", log)
	}, func() {
		chart.TrendIndicator(c, s.change, chart.TrendIndicatorOptions{Precision: int(s.precision), Unit: "%", Note: "vs Lent", Invert: s.invert})
		ui.Row(c).Gap(24).Children(func() {
			chart.TrendIndicator(c, 3, chart.TrendIndicatorOptions{})
			chart.TrendIndicator(c, -3, chart.TrendIndicatorOptions{})
			chart.TrendIndicator(c, 0, chart.TrendIndicatorOptions{})
		})
	})
}

func c2EmptyDemo(c *ui.Context, log func(string)) {
	type state struct {
		status   string
		height   float64
		disabled bool
	}
	s := ui.Local(c.Root(), "140", func() state { return state{status: "Failed", height: 180} })
	c1Section(c, "ChartEmptyState", func() {
		c2Number(c, &s.height, "Height", 160, 320, 40, log)
		c2Check(c, &s.disabled, "Disabled", log)
	}, func() {
		f, _ := c2Frame(c, "ChartEmptyState", &s.status, log)
		ui.Column(c).Disabled(s.disabled).Border(1, core.Tokens(c).Border).Children(func() {
			if chart.ChartEmptyState(c, f.Status, chart.ChartEmptyStateOptions{Height: float32(s.height), Error: f.Error}).Changed() {
				log("ChartEmptyState: retry")
				s.status = "Loading"
			}
		})
	})
}
