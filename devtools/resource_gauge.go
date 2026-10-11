package devtools

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/agent"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

// ResourceReading is one resource's current use.
type ResourceReading struct {
	// Name is required, as "CPU".
	Name        string
	Used, Total float64
	// Unit follows the amounts, as "GB"; empty shows percent only.
	Unit string
}

// ResourceGaugeOptions configure a ResourceGauge; thresholds are fractions of Total.
type ResourceGaugeOptions struct {
	// Warn and Critical default to 0.7 and 0.9.
	Warn, Critical float64
	// Height of each dial; zero means 150.
	Height float32
	Status data.DataStatus
	Error  string
}

// ResourceGauge shows resource usage with named threshold bands.
func ResourceGauge(c *ui.Context, readings []ResourceReading, opts ResourceGaugeOptions) core.ChoiceResult {
	if opts.Status != data.DataReady && opts.Status != data.DataLoading && opts.Status != data.DataFailed {
		panic("mujica: ResourceGauge needs a known status")
	}
	warn, crit := opts.Warn, opts.Critical
	if warn == 0 {
		warn = 0.7
	}
	if crit == 0 {
		crit = 0.9
	}
	if !(0 < warn && warn < crit && crit < 1) || opts.Height < 0 {
		panic("mujica: ResourceGauge needs 0 < Warn < Critical < 1 and Height >= 0")
	}
	h := opts.Height
	if h == 0 {
		h = 150
	}
	u := core.Density(c).Unit()
	row := ui.Row(c).Wrap().Gap(u * 4).Shrink(0)
	result := core.ChoiceResult{Element: row}
	row.Children(func() {
		if len(readings) == 0 {
			ui.Column(c).Height(h).Children(func() {
				result = core.NewChoiceResult(result.Element, data.DataPlaceholder(c, opts.Status, opts.Error, core.Msg(c, "resource_gauge.empty"), ""))
			})
			return
		}
		for i, r := range readings {
			if r.Name == "" || !(r.Total > 0) || r.Used < 0 || r.Used > r.Total || math.IsNaN(r.Used) || math.IsInf(r.Used, 0) || math.IsInf(r.Total, 0) {
				panic(fmt.Sprintf("mujica: ResourceGauge reading %d needs a Name and 0 <= Used <= Total", i))
			}
			frac := r.Used / r.Total
			ui.Column(c).Key(r.Name).Width(200).Gap(u).AlignItems(ui.Center).Children(func() {
				ui.Text(c, r.Name).FontWeight(600)
				gauge := chart.Gauge(c, frac*100, chart.GaugeOptions{
					Frame: chart.ChartFrame{Label: r.Name, Height: h, Status: opts.Status, Error: opts.Error}, Min: 0, Max: 100, Unit: "%",
					Format: func(v float64) string { return fmt.Sprintf("%.0f", v) },
					Segments: []chart.GaugeSegment{
						{To: warn * 100, Label: core.Msg(c, "resource_gauge.normal"), Tone: chart.GaugeGood},
						{To: crit * 100, Label: core.Msg(c, "resource_gauge.high"), Tone: chart.GaugeWarning},
						{To: 100, Label: core.Msg(c, "resource_gauge.critical"), Tone: chart.GaugeDanger},
					},
				})
				result = core.NewChoiceResult(result.Element, result.Changed() || gauge.Retried())
				if r.Unit != "" {
					agent.AgentCaption(c, fmt.Sprintf("%s / %s %s", agent.SandboxNum(r.Used), agent.SandboxNum(r.Total), r.Unit))
				}
			})
		}
	})
	if result.Changed() {
		slog.Debug("mujica: resource_gauge retry")
	}
	return result
}
