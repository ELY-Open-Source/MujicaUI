package mujica

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// AllocationItem is an asset class's value and target weight in percent.
type AllocationItem struct {
	Name   string
	Value  float64
	Target float64
}

// AssetAllocationChartOptions configure an AssetAllocationChart.
type AssetAllocationChartOptions struct {
	Frame ChartFrame
	// Tolerance is the drift in points that warns; zero means 5.
	Tolerance float64
	Format    func(v float64) string
}

// AssetAllocationChart shows holdings as a donut and lists each weight against its target.
func AssetAllocationChart(c *ui.Context, items []AllocationItem, opts AssetAllocationChartOptions) ChartResult {
	if opts.Tolerance < 0 {
		panic("mujica: AssetAllocationChart Tolerance is negative")
	}
	tol := opts.Tolerance
	if tol == 0 {
		tol = 5
	}
	total, targets := 0.0, 0.0
	for _, it := range items {
		if it.Name == "" || it.Value < 0 || it.Target < 0 || math.IsNaN(it.Value) || math.IsNaN(it.Target) {
			panic("mujica: AllocationItem needs a Name and values ≥ 0")
		}
		total += it.Value
		targets += it.Target
	}
	if len(items) > 0 && math.Abs(targets-100) > 0.01 && targets != 0 {
		panic(fmt.Sprintf("mujica: AllocationItem targets sum to %v, not 100", targets))
	}
	k := Tokens(c)
	u := Density(c).Unit()
	pal := internal.Palette(k)
	var res ChartResult
	row := ui.Row(c).Gap(u * 6).AlignItems(ui.Start).Wrap()
	row.Children(func() {
		slices := make([]ChartSlice, len(items))
		for i, it := range items {
			slices[i] = ChartSlice{Label: it.Name, Value: it.Value}
		}
		frame := opts.Frame
		if frame.Height == 0 {
			frame.Height = 220
		}
		ui.Column(c).Width(frame.Height + 40).Children(func() {
			res = PieChart(c, slices, PieChartOptions{Frame: frame, Donut: true, Format: opts.Format})
		})
		if total == 0 {
			return
		}
		size := FontSize(c, theme.CaptionSize)
		g := ui.Grid(c).ColumnTracks(ui.FitContent(), ui.FitContent(), ui.FitContent(), ui.FitContent()).GapX(u * 4).GapY(u * 2).
			Role(ui.RoleTable).Label(frame.Label)
		g.Children(func() {
			for _, h := range []string{"asset_allocation_chart.asset", "asset_allocation_chart.weight", "asset_allocation_chart.target", "asset_allocation_chart.drift"} {
				ui.Text(c, msg(c, h)).FontSize(size).TextColor(k.TextMuted).Role(ui.RoleColumnHeader)
			}
			for i, it := range items {
				w := it.Value / total * 100
				drift := w - it.Target
				ui.Row(c).Key(it.Name).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(12, 12).Draw(func(p *ui.Painter, r ui.Rect) { internal.Marker(p, i, r.X+6, r.Y+6, 10, pal[i%len(pal)]) })
					ui.Text(c, it.Name).SingleLine()
				})
				ui.Text(c, fmt.Sprintf("%.1f%%", w)).FontFeatures("tnum").JustifySelf(ui.End)
				ui.Text(c, fmt.Sprintf("%.1f%%", it.Target)).FontFeatures("tnum").TextColor(k.TextMuted).JustifySelf(ui.End)
				allocationDrift(c, drift, tol)
			}
		})
	})
	res.Element = row
	return res
}

// allocationDrift shows the drift in points; beyond tolerance it warns with an icon and a word.
func allocationDrift(c *ui.Context, drift, tol float64) {
	k := Tokens(c)
	text := statisticSigned(drift, 1) + " pt"
	over := math.Abs(drift) > tol
	ink := k.TextMuted
	if over {
		ink = k.Warning
	}
	ui.Row(c).Gap(4).AlignItems(ui.Center).JustifySelf(ui.End).TextColor(ink).Children(func() {
		if over {
			ui.Icon(c, icons.Must("triangle-alert")).Size(12, 12)
			word := "asset_allocation_chart.over"
			if drift < 0 {
				word = "asset_allocation_chart.under"
			}
			ui.Text(c, msg(c, word)).FontSize(FontSize(c, theme.CaptionSize))
		}
		ui.Text(c, text).FontFeatures("tnum")
	})
}
