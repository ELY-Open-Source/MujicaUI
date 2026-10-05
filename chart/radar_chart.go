package chart

import (
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// RadarChartOptions configure a RadarChart.
type RadarChartOptions struct {
	Frame ChartFrame
	// Max is the outer ring's value; 0 means automatic.
	Max    float64
	Format func(v float64) string
}

// RadarChart draws series as polygons over three or more axes.
func RadarChart(c *ui.Context, axes []string, series []ChartSeries, opts RadarChartOptions) ChartResult {
	if n := chartCheckSeries(series); len(axes) < 3 || n != len(axes) && len(series) > 0 || opts.Max < 0 {
		panic("mujica: RadarChart needs three or more axes, one value per axis and Max ≥ 0")
	}
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	pal := internal.Palette(k)
	size := core.FontSize(c, theme.CaptionSize)
	keys := make([]ChartLegendItem, len(series))
	top := opts.Max
	for i, s := range series {
		keys[i] = ChartLegendItem{Name: s.Name, Color: pal[i%len(pal)], Shape: i}
		for _, v := range s.Values {
			if v < 0 {
				panic("mujica: RadarChart values must be ≥ 0")
			}
			if !ch.Hidden(s.Name) && opts.Max == 0 && !math.IsNaN(v) {
				top = max(top, v)
			}
		}
	}
	n := len(axes)
	if top == 0 {
		n = 0
	}
	var step float64
	if n > 0 && opts.Max == 0 {
		_, top, step = internal.Nice(0, top, 4)
	} else if n > 0 {
		_, _, step = internal.Nice(0, top, 4)
	}
	geom := func(r ui.Rect) (cx, cy, rad float32) {
		return r.X + r.W/2, r.Y + r.H/2 + 4, min(r.W/2-70, r.H/2-24)
	}
	angle := func(i int) float64 { return 2 * math.Pi * float64(i) / float64(len(axes)) }
	spec := ChartSpec{Marks: n, Legend: keys, Hidden: ch.Toggles(len(series) > 1)}
	spec.Readout = func(i int) string {
		parts := []string{axes[i]}
		for _, s := range series {
			if !ch.Hidden(s.Name) {
				parts = append(parts, s.Name+" "+ChartValue(c, s.Values[i], opts.Format))
			}
		}
		return strings.Join(parts, " · ")
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		cx, cy, rad := geom(r)
		if math.Hypot(float64(x-cx), float64(y-cy)) > float64(rad)+16 {
			return -1
		}
		return internal.Nearest(len(axes), func(i int) (float32, float32) { return internal.Polar(cx, cy, rad, angle(i)) }, x, y, 2*rad+32, false)
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		cx, cy, rad := geom(r)
		for _, t := range internal.Ticks(step, top, step) {
			var ring ui.Path
			for i := range axes {
				x, y := internal.Polar(cx, cy, rad*float32(t/top), angle(i))
				if i == 0 {
					ring.MoveTo(x, y)
				} else {
					ring.LineTo(x, y)
				}
			}
			ring.Close()
			p.StrokePath(&ring, 1, k.Border)
			ChartText(p, internal.FormatTick(t, step), cx-4, cy-rad*float32(t/top), size, k.TextMuted, 1, 0.5)
		}
		for i, a := range axes {
			x, y := internal.Polar(cx, cy, rad, angle(i))
			w := float32(1)
			if i == active {
				w = 2
			}
			var spoke ui.Path
			spoke.MoveTo(cx, cy).LineTo(x, y)
			p.StrokePath(&spoke, w, k.ControlBorder)
			lx, ly := internal.Polar(cx, cy, rad+10, angle(i))
			ax := 0.5 + 0.5*float32(math.Sin(angle(i)))
			col := k.TextMuted
			if i == active {
				col = k.Text
			}
			ChartText(p, a, lx, ly, size, col, 1-ax, 0.5+0.5*float32(math.Cos(angle(i))))
		}
		for si, s := range series {
			if ch.Hidden(s.Name) {
				continue
			}
			var poly ui.Path
			pt := func(i int) (float32, float32) {
				v := s.Values[i]
				if math.IsNaN(v) {
					v = 0
				}
				return internal.Polar(cx, cy, rad*float32(v/top), angle(i))
			}
			for i := range axes {
				if i == 0 {
					poly.MoveTo(pt(i))
				} else {
					poly.LineTo(pt(i))
				}
			}
			poly.Close()
			p.FillPath(&poly, keys[si].Color.Alpha(0.16))
			p.StrokePath(&poly, 2, keys[si].Color)
			for i := range axes {
				x, y := pt(i)
				sz := float32(7)
				if i == active {
					sz = 11
				}
				internal.Marker(p, si, x, y, sz, keys[si].Color)
			}
		}
	}
	res, _ := ch.Build(spec)
	return res
}
