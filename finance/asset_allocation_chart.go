package finance

import (
	"fmt"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

type AllocationItem struct {
	Name   string
	Value  float64
	Target float64
}

// AssetAllocationChartOptions configure an AssetAllocationChart.
type AssetAllocationChartOptions struct {
	Frame     chart.ChartFrame
	Version   int
	Tolerance float64
	Format    func(v float64) string
}

func AssetAllocationChart(c *ui.Context, items []AllocationItem, opts AssetAllocationChartOptions) chart.ChartResult {
	if opts.Tolerance < 0 {
		panic("mujica: AssetAllocationChart Tolerance is negative")
	}
	tol := opts.Tolerance
	if tol == 0 {
		tol = 5
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	pal := internal.Palette(k)
	var res chart.ChartResult
	row := ui.Row(c).Gap(u * 6).AlignItems(ui.Start).Wrap()
	row.Children(func() {
		memo := ui.Local(row, "mujica.allocation", func() allocationSlices { return allocationSlices{} })
		if id := chart.ChartID(items) + fmt.Sprint(opts.Version) + core.Msg(c, "pie_chart.other"); id != memo.id {
			total, targets := 0.0, 0.0
			for _, it := range items {
				if it.Name == "" || it.Value < 0 || it.Target < 0 || math.IsNaN(it.Value) || math.IsNaN(it.Target) || math.IsInf(it.Value, 0) || math.IsInf(it.Target, 0) {
					panic("mujica: AllocationItem needs a Name and values ≥ 0")
				}
				total += it.Value
				targets += it.Target
			}
			if len(items) > 0 && math.Abs(targets-100) > 0.01 && targets != 0 {
				panic(fmt.Sprintf("mujica: AllocationItem targets sum to %v, not 100", targets))
			}

			memo.total = total
			memo.id, memo.slices = id, make([]chart.ChartSlice, len(items))
			for i, it := range items {
				memo.slices[i] = chart.ChartSlice{Label: it.Name, Value: it.Value}
			}
			memo.parts = chart.PieMerge(memo.slices, 1, core.Msg(c, "pie_chart.other"))
			memo.markers = make([]int, len(items))
			for i, it := range items {
				at := slices.IndexFunc(memo.parts, func(part chart.ChartSlice) bool { return part.Label == it.Name })
				if at < 0 {
					at = len(memo.parts) - 1
				}
				memo.markers[i] = at
			}
		}
		parts := memo.parts
		total := memo.total
		frame := opts.Frame
		if frame.Height == 0 {
			frame.Height = 220
		}
		ui.Column(c).Width(max(frame.Height+40, 300)).Children(func() {
			res = chart.PieChart(c, parts, chart.PieChartOptions{Frame: frame, Donut: true, Format: opts.Format})
		})
		if total == 0 {
			return
		}
		size := core.FontSize(c, theme.CaptionSize)
		ui.Column(c).MinWidth(340).Grow(1).Children(func() {
			g := ui.Grid(c).ColumnTracks(ui.Fr(2), ui.Fr(1), ui.Fr(1), ui.Fr(2)).GapX(u * 2).GapY(u * 2).
				Role(ui.RoleTable).Label(frame.Label)
			g.Children(func() {
				for _, h := range []string{"asset_allocation_chart.asset", "asset_allocation_chart.weight", "asset_allocation_chart.target", "asset_allocation_chart.drift"} {
					ui.Text(c, core.Msg(c, h)).FontSize(size).TextColor(k.TextMuted).Role(ui.RoleColumnHeader)
				}
			})
			state := ui.Local(row, "mujica.allocation.rows", func() ui.ListState { return ui.ListState{} })
			ui.List(c, state, len(items), func(i int) {
				it := items[i]
				ui.Grid(c).ColumnTracks(ui.Fr(2), ui.Fr(1), ui.Fr(1), ui.Fr(2)).GapX(u * 2).MinHeight(core.Density(c).ControlHeight()).Children(func() {
					w := it.Value / total * 100
					drift := w - it.Target
					ui.Row(c.Key(it.Name)).Gap(u * 2).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(12, 12).Draw(func(p *ui.Painter, r ui.Rect) {
							internal.Marker(p, memo.markers[i], r.X+6, r.Y+6, 10, pal[memo.markers[i]%len(pal)])
						})
						ui.Text(c, it.Name).SingleLine()
					})
					ui.Text(c, fmt.Sprintf("%.1f%%", w)).FontFeatures("tnum").JustifySelf(ui.End)
					ui.Text(c, fmt.Sprintf("%.1f%%", it.Target)).FontFeatures("tnum").TextColor(k.TextMuted).JustifySelf(ui.End)
					allocationDrift(c, drift, tol)
				})
			}).Height(min(float32(len(items))*core.Density(c).ControlHeight(), 240))
		})
	})
	res.Element = row
	return res
}

type allocationSlices struct {
	id      string
	slices  []chart.ChartSlice
	total   float64
	parts   []chart.ChartSlice
	markers []int
}

func allocationDrift(c *ui.Context, drift, tol float64) {
	k := core.Tokens(c)
	text := data.StatisticSigned(drift, 1) + " pt"
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
			ui.Text(c, core.Msg(c, word)).FontSize(core.FontSize(c, theme.CaptionSize))
		}
		ui.Text(c, text).FontFeatures("tnum")
	})
}
