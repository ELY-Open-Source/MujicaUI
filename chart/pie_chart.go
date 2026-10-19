package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ChartSlice is one share of a PieChart or FunnelChart.
type ChartSlice struct {
	Label string
	Value float64
}

// PieChartOptions configure a PieChart.
type PieChartOptions struct {
	Frame ChartFrame
	// Donut cuts out the middle and writes the total there.
	Donut bool
	// MinAngle merges two or more narrower slices into Other.
	MinAngle float64
	Format   func(v float64) string
}

// PieChart shows shares as slices labeled with percents.
func PieChart(c *ui.Context, slices []ChartSlice, opts PieChartOptions) ChartResult {
	if opts.MinAngle < 0 || opts.MinAngle >= 360 {
		panic("mujica: PieChart MinAngle must be in [0, 360)")
	}
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	pal := internal.Palette(k)
	size := core.FontSize(c, theme.CaptionSize)
	parts := ChartMemo(ch, "parts", [2]any{ChartID(slices), opts.MinAngle}, func() []ChartSlice {
		return PieMerge(slices, opts.MinAngle, core.Msg(c, "pie_chart.other"))
	})
	total := 0.0
	for _, s := range parts {
		total += s.Value
	}
	n := len(parts)
	if total == 0 {
		n = 0
	}
	keys := make([]ChartLegendItem, len(parts))
	for i, s := range parts {
		keys[i] = ChartLegendItem{Name: s.Label, Color: pal[i%len(pal)], Shape: i}
	}
	// angle returns slice i's start and end, clockwise.
	angle := func(i int) (float64, float64) {
		a := 0.0
		for j := range i {
			a += parts[j].Value
		}
		return 2 * math.Pi * a / total, 2 * math.Pi * (a + parts[i].Value) / total
	}
	geom := func(r ui.Rect) (cx, cy, r0, r1 float32) {
		r1 = min(r.W, r.H)/2 - 8
		if opts.Donut {
			r0 = r1 * 0.58
		}
		return r.X + r.W/2, r.Y + r.H/2, r0, r1
	}
	share := func(i int) string { return fmt.Sprintf("%.1f%%", 100*parts[i].Value/total) }
	spec := ChartSpec{Marks: n, Legend: keys}
	spec.Readout = func(i int) string {
		return parts[i].Label + " · " + ChartValue(c, parts[i].Value, opts.Format) + " · " + share(i)
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		cx, cy, r0, r1 := geom(r)
		d := float32(math.Hypot(float64(x-cx), float64(y-cy)))
		if d < r0 || d > r1 {
			return -1
		}
		a := math.Atan2(float64(x-cx), float64(cy-y))
		if a < 0 {
			a += 2 * math.Pi
		}
		for i := range parts {
			if a0, a1 := angle(i); a >= a0 && a < a1 {
				return i
			}
		}
		return -1
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		cx, cy, r0, r1 := geom(r)
		for i := range parts {
			a0, a1 := angle(i)
			if a1 <= a0 {
				continue
			}
			grow := float32(0)
			if i == active {
				grow = 6
			}
			p.FillPath(internal.Sector(cx, cy, r0, r1+grow, a0, a1), keys[i].Color)
			p.StrokePath(internal.Sector(cx, cy, r0, r1+grow, a0, a1), 1.5, k.Surface)
			// Wide slices carry their shape and percent.
			if a1-a0 > 0.38 {
				mx, my := internal.Polar(cx, cy, (r0+r1)/2+float32(chartBool(r0 == 0))*r1*0.12, (a0+a1)/2)
				ink := ChartInk(k, keys[i].Color)
				internal.Marker(p, i, mx, my-8, 9, ink)
				// The percent shows where the ring is wide enough to hold it.
				if w, _ := p.MeasureText(0, ui.Span{Text: share(i), Size: size}); r1-r0 >= w+6 {
					ChartText(p, share(i), mx, my+2, size, ink, 0.5, 0)
				}
			}
		}
		if opts.Donut {
			ChartText(p, ChartValue(c, total, opts.Format), cx, cy, core.FontSize(c, theme.H3Size), k.Text, 0.5, 0.75)
			ChartText(p, core.Msg(c, "pie_chart.total"), cx, cy+4, size, k.TextMuted, 0.5, 0)
		}
	}
	res, _ := ch.Build(spec)
	return res
}

// PieMerge folds two or more narrow slices into other.
func PieMerge(slices []ChartSlice, minAngle float64, other string) []ChartSlice {
	total := 0.0
	for _, s := range slices {
		if s.Value < 0 || math.IsNaN(s.Value) || math.IsInf(s.Value, 0) || s.Label == "" {
			panic("mujica: PieChart slices need a label and a finite value ≥ 0")
		}
		total += s.Value
	}
	var keep []ChartSlice
	small, rest := 0, 0.0
	for _, s := range slices {
		if total > 0 && s.Value/total*360 < minAngle {
			small++
			rest += s.Value
		}
	}
	for _, s := range slices {
		if small < 2 || total == 0 || s.Value/total*360 >= minAngle {
			keep = append(keep, s)
		}
	}
	if small >= 2 {
		keep = append(keep, ChartSlice{Label: other, Value: rest})
	}
	return keep
}
