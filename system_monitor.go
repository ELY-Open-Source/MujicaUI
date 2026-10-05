package mujica

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// MonitorMetric is one sampled resource; Samples line up with the shared times.
type MonitorMetric struct {
	Name    string
	Unit    string
	Samples []float64
}

// SystemMonitorOptions configure a SystemMonitor.
type SystemMonitorOptions struct {
	// Height of each chart; zero means 110.
	Height float32
	Status DataStatus
	Error  string
}

// SystemMonitorResult reports this frame's events.
type SystemMonitorResult struct {
	Element *ui.Element
	Retry   bool
}

// SystemMonitor stacks a live line chart per metric over one time axis, sampled by the caller.
func SystemMonitor(c *ui.Context, times []time.Time, metrics []MonitorMetric, opts SystemMonitorOptions) SystemMonitorResult {
	if opts.Height < 0 {
		panic("mujica: SystemMonitor Height is negative")
	}
	for i, m := range metrics {
		if m.Name == "" || len(m.Samples) != len(times) {
			panic(fmt.Sprintf("mujica: SystemMonitor metric %d needs a Name and one sample per time", i))
		}
	}
	h := opts.Height
	if h == 0 {
		h = 110
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r SystemMonitorResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 3)
	r.Element.Children(func() {
		for _, m := range metrics {
			ui.Column(c).Key(m.Name).Gap(u).Children(func() {
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					ui.Text(c, m.Name).FontWeight(600).Grow(1)
					last := math.NaN()
					if n := len(m.Samples); n > 0 {
						last = m.Samples[n-1]
					}
					if !math.IsNaN(last) && opts.Status == DataReady {
						ui.Text(c, fmt.Sprintf("%s %s", sandboxNum(last), m.Unit)).FontFeatures("tnum").
							FontSize(FontSize(c, theme.CaptionSize+1)).TextColor(k.Text).Role(ui.RoleStatus)
					}
				})
				unit := m.Unit
				res := LineChart(c, []ChartSeries{{Name: m.Name, Values: m.Samples}}, LineChartOptions{
					Frame: ChartFrame{Label: m.Name, Height: h, Status: opts.Status, Error: opts.Error},
					Times: times, Format: func(v float64) string { return sandboxNum(v) + " " + unit },
				})
				r.Retry = r.Retry || res.Retried()
			})
		}
	})
	return r
}
