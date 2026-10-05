package mujica

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
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
}

// ResourceGauge shows CPU, memory and disk dials with named normal, high and critical bands.
func ResourceGauge(c *ui.Context, readings []ResourceReading, opts ResourceGaugeOptions) *ui.Element {
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
	u := Density(c).Unit()
	row := ui.Row(c).Wrap().Gap(u * 4).Shrink(0)
	row.Children(func() {
		for i, r := range readings {
			if r.Name == "" || !(r.Total > 0) || r.Used < 0 || r.Used > r.Total || math.IsNaN(r.Used) {
				panic(fmt.Sprintf("mujica: ResourceGauge reading %d needs a Name and 0 <= Used <= Total", i))
			}
			frac := r.Used / r.Total
			ui.Column(c).Key(r.Name).Width(200).Gap(u).AlignItems(ui.Center).Children(func() {
				ui.Text(c, r.Name).FontWeight(600)
				Gauge(c, frac*100, GaugeOptions{
					Frame: ChartFrame{Label: r.Name, Height: h}, Min: 0, Max: 100, Unit: "%",
					Format: func(v float64) string { return fmt.Sprintf("%.0f", v) },
					Segments: []GaugeSegment{
						{To: warn * 100, Label: msg(c, "resource_gauge.normal"), Tone: GaugeGood},
						{To: crit * 100, Label: msg(c, "resource_gauge.high"), Tone: GaugeWarning},
						{To: 100, Label: msg(c, "resource_gauge.critical"), Tone: GaugeDanger},
					},
				})
				if r.Unit != "" {
					agentCaption(c, fmt.Sprintf("%s / %s %s", sandboxNum(r.Used), sandboxNum(r.Total), r.Unit))
				}
			})
		}
	})
	return row
}
