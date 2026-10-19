package finance

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

type OptionLegKind int

const (
	LegCall OptionLegKind = iota
	LegPut
	// LegStock is the underlying, bought or sold at Premium.
	LegStock
)

type OptionLeg struct {
	Kind     OptionLegKind
	Strike   float64
	Premium  float64
	Quantity float64
}

func PayoffAt(legs []OptionLeg, price float64) float64 {
	sum := 0.0
	for _, l := range legs {
		switch l.Kind {
		case LegCall:
			sum += l.Quantity * (max(price-l.Strike, 0) - l.Premium)
		case LegPut:
			sum += l.Quantity * (max(l.Strike-price, 0) - l.Premium)
		case LegStock:
			sum += l.Quantity * (price - l.Premium)
		default:
			panic("mujica: OptionLeg Kind out of range")
		}
	}
	return sum
}

func Breakevens(legs []OptionLeg, lo, hi float64) []float64 {
	kinks := []float64{lo, hi}
	for _, l := range legs {
		if l.Kind != LegStock && l.Strike > lo && l.Strike < hi {
			kinks = append(kinks, l.Strike)
		}
	}
	slices.Sort(kinks)
	var out []float64
	for i := 1; i < len(kinks); i++ {
		a, b := kinks[i-1], kinks[i]
		fa, fb := PayoffAt(legs, a), PayoffAt(legs, b)
		if fa == 0 && (len(out) == 0 || out[len(out)-1] != a) {
			out = append(out, a)
		}
		if fa*fb < 0 {
			out = append(out, a+(b-a)*fa/(fa-fb))
		}
	}
	if PayoffAt(legs, hi) == 0 && (len(out) == 0 || out[len(out)-1] != hi) {
		out = append(out, hi)
	}
	return out
}

// PayoffDiagramOptions configure a PayoffDiagram.
type PayoffDiagramOptions struct {
	Frame     chart.ChartFrame
	Spot      float64
	Low, High float64
	Points    int
	Format    func(v float64) string
}

func PayoffDiagram(c *ui.Context, legs []OptionLeg, opts PayoffDiagramOptions) chart.ChartResult {
	lo, hi := payoffRange(legs, opts)
	n := opts.Points
	if n == 0 {
		n = 101
	}
	if n < 2 {
		panic("mujica: PayoffDiagram needs at least 2 Points")
	}
	ch := chart.NewChart(c, opts.Frame)
	k := core.Tokens(c)
	size := core.FontSize(c, theme.CaptionSize)
	type data struct {
		px, pl, be []float64
		gx, gp     []float64
		ylo, yhi   float64
	}
	d := chart.ChartMemo(ch, "payoff", fmt.Sprint(legs, lo, hi, n), func() data {
		var v data
		v.be = Breakevens(legs, lo, hi)
		v.ylo, v.yhi = 0, 0
		for i := range n {
			x := lo + (hi-lo)*float64(i)/float64(n-1)
			v.px, v.pl = append(v.px, x), append(v.pl, PayoffAt(legs, x))
			v.ylo, v.yhi = min(v.ylo, v.pl[i]), max(v.yhi, v.pl[i])
		}
		v.gx = slices.Clone(v.px)
		for _, l := range legs {
			if l.Kind != LegStock && l.Strike > lo && l.Strike < hi {
				v.gx = append(v.gx, l.Strike)
			}
		}
		v.gx = append(v.gx, v.be...)
		slices.Sort(v.gx)
		v.gx = slices.Compact(v.gx)
		for _, x := range v.gx {
			value := PayoffAt(legs, x)
			v.gp = append(v.gp, value)
			v.ylo, v.yhi = min(v.ylo, value), max(v.yhi, value)
		}
		return v
	})
	if len(legs) == 0 {
		n = 0
	}
	format := func(v float64) string { return chart.ChartValue(c, v, opts.Format) }
	ylo, yhi, step := internal.Nice(d.ylo, d.yhi, 4)
	ticks := internal.Ticks(ylo, yhi, step)
	geom := func(r ui.Rect) (ui.Rect, internal.Linear, internal.Linear) {
		plot := chart.ChartInset(r, 56, 26)
		return plot, internal.Linear{D0: lo, D1: hi, R0: float64(plot.X), R1: float64(plot.X + plot.W)},
			internal.Linear{D0: ylo, D1: yhi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	up, down := core.TrendColor(c, 1), core.TrendColor(c, -1)
	spec := chart.ChartSpec{Marks: n}
	spec.Readout = func(i int) string {
		word := [3]string{"payoff_diagram.loss", "payoff_diagram.even", "payoff_diagram.profit"}[core.TrendSign(d.pl[i])+1]
		return fmt.Sprintf(core.Msg(c, "payoff_diagram.readout"), format(d.px[i]), core.Msg(c, word), format(d.pl[i]))
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, xs, _ := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		return max(0, min(n-1, int(math.Round((xs.Invert(x)-lo)/(hi-lo)*float64(n-1)))))
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, xs, ys := geom(r)
		internal.AxisY(p, plot, ys, ticks, func(v float64) string { return internal.FormatTick(v, step) }, size, k)
		lx, ll, lstep := internal.Nice(lo, hi, max(2, int(plot.W/90)))
		var tx []float32
		var tl []string
		for _, t := range internal.Ticks(lx, ll, lstep) {
			if t >= lo && t <= hi {
				tx, tl = append(tx, xs.Map(t)), append(tl, internal.FormatTick(t, lstep))
			}
		}
		internal.AxisX(p, plot, tx, tl, size, k)
		zero := internal.Snap(p, ys.Map(0))
		pt := func(x, v float64) (float32, float32) { return internal.Snap(p, xs.Map(x)), internal.Snap(p, ys.Map(v)) }
		// Exact corners: the grid plus strikes and breakevens.
		indices := chart.ChartMemo(ch, "payoff.geometry", [2]any{chart.ChartID(d.gx), plot.W}, func() []int {
			return internal.Columns(len(d.gx), func(i int) float32 { return float32((d.gx[i]-lo)/(hi-lo)) * plot.W }, func(i int) float64 { return d.gp[i] })
		})

		for j := 1; j < len(indices); j++ {
			prev, i := indices[j-1], indices[j]
			a, b := d.gx[prev], d.gx[i]
			fa, fb := d.gp[prev], d.gp[i]
			col := up
			if fa+fb < 0 {
				col = down
			}
			var area ui.Path
			ax, ay := pt(a, fa)
			bx, by := pt(b, fb)
			area.MoveTo(ax, zero).LineTo(ax, ay).LineTo(bx, by).LineTo(bx, zero).Close()
			p.FillPath(&area, col.Alpha(0.16))
		}
		p.Line(plot.X, zero, plot.X+plot.W, zero, 1, k.ControlBorder)
		var line ui.Path
		for j, i := range indices {
			if j == 0 {
				line.MoveTo(pt(d.gx[i], d.gp[i]))
			} else {
				line.LineTo(pt(d.gx[i], d.gp[i]))
			}
		}
		p.StrokePath(&line, 2, k.AccentText)
		for _, b := range d.be {
			bx, _ := pt(b, 0)
			p.FillPath(internal.Diamond(bx, zero, 9), k.Text)
			chart.ChartText(p, format(b), bx, zero-6, size, k.Text, 0.5, 1)
		}
		if opts.Spot > lo && opts.Spot < hi {
			sx, _ := pt(opts.Spot, 0)
			for y := plot.Y; y < plot.Y+plot.H; y += 6 {
				p.Line(sx, y, sx, min(y+3, plot.Y+plot.H), 1, k.Info)
			}
			chart.ChartText(p, core.Msg(c, "payoff_diagram.spot")+" "+format(opts.Spot), sx+4, plot.Y, size, k.Info, 0, 0)
		}
		if active >= 0 {
			x, y := pt(d.px[active], d.pl[active])
			p.Line(x, plot.Y, x, plot.Y+plot.H, 1, k.TextMuted)
			internal.Marker(p, 0, x, y, 13, k.Surface)
			internal.Marker(p, 0, x, y, 9, k.AccentText)
		}
	}
	res, _ := ch.Build(spec)
	if len(d.be) > 0 && n > 0 {
		parts := make([]string, len(d.be))
		for i, b := range d.be {
			parts[i] = format(b)
		}
		res.Element.Children(func() {
			ui.Text(c, core.Msg(c, "payoff_diagram.breakevens")+" "+strings.Join(parts, ", ")).FontSize(size).TextColor(k.TextMuted)
		})
	}
	return res
}

func payoffRange(legs []OptionLeg, opts PayoffDiagramOptions) (lo, hi float64) {
	if opts.Low != 0 || opts.High != 0 {
		if !(opts.High > opts.Low) || opts.Low < 0 {
			panic("mujica: PayoffDiagram needs 0 ≤ Low < High")
		}
		return opts.Low, opts.High
	}
	lo, hi = math.Inf(1), math.Inf(-1)
	for _, l := range legs {
		ref := l.Strike
		if l.Kind == LegStock {
			ref = l.Premium
		}
		lo, hi = min(lo, ref), max(hi, ref)
	}
	if opts.Spot > 0 {
		lo, hi = min(lo, opts.Spot), max(hi, opts.Spot)
	}
	if math.IsInf(lo, 1) || hi <= 0 {
		return 0, 100
	}
	return lo * 0.7, hi * 1.3
}
