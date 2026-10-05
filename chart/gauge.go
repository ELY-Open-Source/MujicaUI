package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// GaugeTone colors a Gauge segment; its label still names it.
type GaugeTone int

const (
	GaugeNeutral GaugeTone = iota
	GaugeGood
	GaugeWarning
	GaugeDanger
)

// GaugeSegment is a named dial range up to To.
type GaugeSegment struct {
	To    float64
	Label string
	Tone  GaugeTone
}

// GaugeOptions configure a Gauge; Frame.Height 0 means 200 DIPs.
type GaugeOptions struct {
	Frame    ChartFrame
	Min, Max float64
	// Segments split Min..Max ascending; the last ends at Max.
	Segments []GaugeSegment
	Unit     string
	// Full draws a whole circle instead of a half.
	Full   bool
	Format func(v float64) string
}

// Gauge points a needle at value on named segments.
func Gauge(c *ui.Context, value float64, opts GaugeOptions) ChartResult {
	if !(opts.Max > opts.Min) || math.IsNaN(value) || value < opts.Min || value > opts.Max {
		panic("mujica: Gauge needs Min < Max and a value between them")
	}
	prev := opts.Min
	for i, s := range opts.Segments {
		if s.To <= prev || s.Label == "" || i == len(opts.Segments)-1 && s.To != opts.Max {
			panic("mujica: Gauge segments need labels and ascending ends finishing at Max")
		}
		prev = s.To
	}
	if opts.Frame.Height == 0 {
		opts.Frame.Height = 200
	}
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	size := core.FontSize(c, theme.CaptionSize)
	tone := func(t GaugeTone) ui.Color {
		return map[GaugeTone]ui.Color{GaugeNeutral: internal.Palette(k)[5], GaugeGood: k.Success, GaugeWarning: k.Warning, GaugeDanger: k.Danger}[t]
	}
	val := func(v float64) string { return ChartValue(c, v, opts.Format) }
	seg := -1
	for i, s := range opts.Segments {
		if value <= s.To && seg < 0 {
			seg = i
		}
	}
	a0, a1 := -math.Pi/2, math.Pi/2
	if opts.Full {
		a0, a1 = -math.Pi, math.Pi
	}
	angle := func(v float64) float64 { return a0 + (a1-a0)*(v-opts.Min)/(opts.Max-opts.Min) }
	geom := func(r ui.Rect) (cx, cy, rad float32) {
		if opts.Full {
			rad = min(r.W, r.H)/2 - 12
			return r.X + r.W/2, r.Y + r.H/2, rad
		}
		rad = min(r.W/2-12, r.H-36)
		return r.X + r.W/2, r.Y + 8 + rad, rad
	}
	reading := val(value) + " " + opts.Unit
	if seg >= 0 {
		reading += " · " + opts.Segments[seg].Label
	}
	spec := ChartSpec{Marks: 1 + len(opts.Segments)}
	spec.Readout = func(i int) string {
		if i == 0 {
			return reading
		}
		s, from := opts.Segments[i-1], opts.Min
		if i > 1 {
			from = opts.Segments[i-2].To
		}
		return fmt.Sprintf("%s: %s–%s", s.Label, val(from), val(s.To))
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		cx, cy, rad := geom(r)
		d := float32(math.Hypot(float64(x-cx), float64(y-cy)))
		if d < rad*0.45 {
			return 0
		}
		if d > rad+6 {
			return -1
		}
		a := math.Atan2(float64(x-cx), float64(cy-y))
		from := a0
		for i, s := range opts.Segments {
			if a >= from && a <= angle(s.To) {
				return i + 1
			}
			from = angle(s.To)
		}
		return -1
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		cx, cy, rad := geom(r)
		from := a0
		if len(opts.Segments) == 0 {
			p.FillPath(internal.Sector(cx, cy, rad*0.78, rad, a0, a1), k.SurfaceHover)
		}
		for i, s := range opts.Segments {
			to := angle(s.To)
			grow := float32(0)
			if active == i+1 {
				grow = 4
			}
			p.FillPath(internal.Sector(cx, cy, rad*0.78-grow, rad+grow, from, to), tone(s.Tone))
			p.StrokePath(internal.Sector(cx, cy, rad*0.78-grow, rad+grow, from, to), 1.5, k.Surface)
			lx, ly := internal.Polar(cx, cy, rad+10, to)
			if i < len(opts.Segments)-1 {
				ChartText(p, val(s.To), lx, ly, size, k.TextMuted, 0.5, 1)
			}
			from = to
		}
		nx, ny := internal.Polar(cx, cy, rad*0.9, angle(value))
		var needle ui.Path
		needle.MoveTo(cx, cy).LineTo(internal.Snap(p, nx), internal.Snap(p, ny))
		w := float32(3)
		if active == 0 {
			w = 4.5
		}
		p.StrokePath(&needle, w, k.Text)
		p.FillPath(internal.Diamond(cx, cy, 12), k.Text)
		ty := cy + 34
		ChartText(p, val(value)+" "+opts.Unit, cx, ty, core.FontSize(c, theme.H3Size), k.Text, 0.5, 1)
		if seg >= 0 {
			ChartText(p, opts.Segments[seg].Label, cx, ty+2, size, k.TextMuted, 0.5, 0)
		}
	}
	res, _ := ch.Build(spec)
	return res
}
