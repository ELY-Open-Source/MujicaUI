package mujica

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// MonitorMetric holds samples aligned with the shared timestamps.
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

// SystemMonitor combines resource charts sampled by the caller.
func SystemMonitor(c *ui.Context, times []time.Time, metrics []MonitorMetric, opts SystemMonitorOptions) ChoiceResult {
	if opts.Height < 0 || opts.Status != DataReady && opts.Status != DataLoading && opts.Status != DataFailed {
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
	var r ChoiceResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 3)
	r.Element.Children(func() {
		if len(metrics) == 0 {
			ui.Column(c).Height(h).Children(func() { r.changed = dataPlaceholder(c, opts.Status, opts.Error, msg(c, "system_monitor.empty"), "") })
			return
		}
		for _, m := range metrics {
			panel := ui.Column(c).Key(m.Name).Gap(u)
			series := ui.Local(panel, "mujica.monitor.series", func() []ChartSeries { return []ChartSeries{{Name: m.Name, Values: m.Samples}} })
			(*series)[0].Name, (*series)[0].Values = m.Name, m.Samples
			panel.Children(func() {
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
				res := LineChart(c, *series, LineChartOptions{
					Frame: ChartFrame{Label: m.Name, Height: h, Status: opts.Status, Error: opts.Error},
					Times: times, Format: func(v float64) string { return sandboxNum(v) + " " + unit },
				})
				r.changed = r.changed || res.Retried()
			})
		}
	})
	if r.changed {
		slog.Debug("mujica: system_monitor retry")
	}
	return r
}
