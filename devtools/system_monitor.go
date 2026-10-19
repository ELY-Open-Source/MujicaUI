package devtools

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/agent"
	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/theme"
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
	Status data.DataStatus
	Error  string
}

// SystemMonitor combines resource charts sampled by the caller.
func SystemMonitor(c *ui.Context, times []time.Time, metrics []MonitorMetric, opts SystemMonitorOptions) core.ChoiceResult {
	if opts.Height < 0 || opts.Status != data.DataReady && opts.Status != data.DataLoading && opts.Status != data.DataFailed {
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
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = ui.Column(c).Shrink(0).Gap(u * 3)
	r.Element.Children(func() {
		if len(metrics) == 0 {
			ui.Column(c).Height(h).Children(func() {
				r = core.NewChoiceResult(r.Element, data.DataPlaceholder(c, opts.Status, opts.Error, core.Msg(c, "system_monitor.empty"), ""))
			})
			return
		}
		for _, m := range metrics {
			panel := ui.Column(c.Key(m.Name)).Gap(u)
			series := ui.Local(panel, "mujica.monitor.series", func() []chart.ChartSeries { return []chart.ChartSeries{{Name: m.Name, Values: m.Samples}} })
			(*series)[0].Name, (*series)[0].Values = m.Name, m.Samples
			panel.Children(func() {
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					ui.Text(c, m.Name).FontWeight(600).Grow(1)
					last := math.NaN()
					if n := len(m.Samples); n > 0 {
						last = m.Samples[n-1]
					}
					if !math.IsNaN(last) && opts.Status == data.DataReady {
						ui.Text(c, fmt.Sprintf("%s %s", agent.SandboxNum(last), m.Unit)).FontFeatures("tnum").
							FontSize(core.FontSize(c, theme.CaptionSize+1)).TextColor(k.Text).Role(ui.RoleStatus)
					}
				})
				unit := m.Unit
				res := chart.LineChart(c, *series, chart.LineChartOptions{
					Frame: chart.ChartFrame{Label: m.Name, Height: h, Status: opts.Status, Error: opts.Error},
					Times: times, Format: func(v float64) string { return agent.SandboxNum(v) + " " + unit },
				})
				r = core.NewChoiceResult(r.Element, r.Changed() || res.Retried())
			})
		}
	})
	if r.Changed() {
		slog.Debug("mujica: system_monitor retry")
	}
	return r
}
