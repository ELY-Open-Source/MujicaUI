package mujica

import (
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChartAxisKind is how an axis maps values.
type ChartAxisKind int

const (
	AxisLinear ChartAxisKind = iota
	// AxisLog maps powers of ten evenly; Min must be above zero.
	AxisLog
	// AxisTime maps From..To; values are Unix seconds.
	AxisTime
	// AxisCategory gives each category a band; values are indices.
	AxisCategory
)

// ChartSide is where an axis runs along the plot.
type ChartSide int

const (
	SideBottom ChartSide = iota
	SideLeft
	SideRight
	SideTop
)

// ChartAxisOptions configure a ChartAxis.
type ChartAxisOptions struct {
	Kind ChartAxisKind
	Side ChartSide
	// Min and Max span a linear or log axis.
	Min, Max float64
	// From and To span a time axis.
	From, To   time.Time
	Categories []string
	// Ticks fixes the tick values; nil picks them.
	Ticks []float64
	// Format writes tick and readout labels; nil picks one per kind.
	Format func(v float64) string
}

// ChartAxisPart is an axis ready to map values and draw itself in a plot.
type ChartAxisPart struct {
	opts   ChartAxisOptions
	lo, hi float64
	ticks  []float64
	step   float64
	layout string
	k      theme.Tokens
	size   float32
}

// ChartAxis resolves an axis: its domain, ticks and labels.
func ChartAxis(c *ui.Context, opts ChartAxisOptions) ChartAxisPart {
	a := ChartAxisPart{opts: opts, k: Tokens(c), size: FontSize(c, theme.CaptionSize)}
	switch opts.Kind {
	case AxisLinear, AxisLog:
		a.lo, a.hi = opts.Min, opts.Max
		if !(a.hi > a.lo) || opts.Kind == AxisLog && a.lo <= 0 {
			panic("mujica: ChartAxis needs Min < Max, and Min > 0 for a log axis")
		}
	case AxisTime:
		if !opts.To.After(opts.From) {
			panic("mujica: ChartAxis time axis needs From before To")
		}
		a.lo, a.hi = chartSeconds(opts.From), chartSeconds(opts.To)
	case AxisCategory:
		if len(opts.Categories) == 0 {
			panic("mujica: ChartAxis category axis needs Categories")
		}
		a.lo, a.hi = -0.5, float64(len(opts.Categories))-0.5
	default:
		panic("mujica: unknown ChartAxisKind")
	}
	if opts.Side > SideTop {
		panic("mujica: unknown ChartSide")
	}
	a.ticks = opts.Ticks
	if a.ticks == nil {
		a.ticks = a.autoTicks()
	}
	return a
}

func chartSeconds(t time.Time) float64 { return float64(t.UnixNano()) / 1e9 }

// autoTicks picks round ticks inside the domain.
func (a *ChartAxisPart) autoTicks() []float64 {
	var out []float64
	switch a.opts.Kind {
	case AxisLog:
		// Powers of ten, with 2 and 5 between them on short spans.
		mants := []float64{1}
		if a.hi/a.lo < 1000 {
			mants = []float64{1, 2, 5}
		}
		for e := math.Floor(math.Log10(a.lo)); math.Pow(10, e) <= a.hi*(1+1e-9); e++ {
			for _, m := range mants {
				if v := m * math.Pow(10, e); v >= a.lo*(1-1e-9) && v <= a.hi*(1+1e-9) {
					out = append(out, v)
				}
			}
		}
	case AxisTime:
		ts, layout := internal.TimeTicks(a.opts.From, a.opts.To, 5)
		a.layout = layout
		for _, t := range ts {
			out = append(out, chartSeconds(t))
		}
	case AxisCategory:
		for i := range a.opts.Categories {
			out = append(out, float64(i))
		}
	default:
		lo, hi, step := internal.Nice(a.lo, a.hi, 5)
		a.step = step
		for _, t := range internal.Ticks(lo, hi, step) {
			if t >= a.lo-step*1e-9 && t <= a.hi+step*1e-9 {
				out = append(out, t)
			}
		}
	}
	return out
}

// horizontal reports whether the axis runs across the plot.
func (a ChartAxisPart) horizontal() bool {
	return a.opts.Side == SideBottom || a.opts.Side == SideTop
}

// fraction is how far along the axis v lies, 0 to 1.
func (a ChartAxisPart) fraction(v float64) float64 {
	if a.opts.Kind == AxisLog {
		return (math.Log10(v) - math.Log10(a.lo)) / (math.Log10(a.hi) - math.Log10(a.lo))
	}
	return (v - a.lo) / (a.hi - a.lo)
}

// Map returns the pixel of v along the axis inside plot.
func (a ChartAxisPart) Map(plot ui.Rect, v float64) float32 {
	f := float32(a.fraction(v))
	if a.horizontal() {
		return plot.X + f*plot.W
	}
	return plot.Y + plot.H - f*plot.H
}

// Invert returns the value at pixel px along the axis inside plot.
func (a ChartAxisPart) Invert(plot ui.Rect, px float32) float64 {
	f := float64((px - plot.X) / plot.W)
	if !a.horizontal() {
		f = float64((plot.Y + plot.H - px) / plot.H)
	}
	if a.opts.Kind == AxisLog {
		return math.Pow(10, math.Log10(a.lo)+f*(math.Log10(a.hi)-math.Log10(a.lo)))
	}
	return a.lo + f*(a.hi-a.lo)
}

// Ticks returns the tick values.
func (a ChartAxisPart) Ticks() []float64 { return a.ticks }

// Label writes v as the axis does.
func (a ChartAxisPart) Label(v float64) string {
	switch {
	case a.opts.Format != nil:
		return a.opts.Format(v)
	case a.opts.Kind == AxisCategory:
		return a.opts.Categories[max(0, min(len(a.opts.Categories)-1, int(math.Round(v))))]
	case a.opts.Kind == AxisTime:
		layout := a.layout
		if layout == "" {
			layout = "2006-01-02 15:04"
		}
		return chartTime(v).Format(layout)
	case a.opts.Kind == AxisLog:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case a.step > 0:
		return internal.FormatTick(v, a.step)
	}
	return strconv.FormatFloat(v, 'g', 6, 64)
}

// readout writes v one decimal finer than ticks, for crosshair labels.
func (a ChartAxisPart) readout(v float64) string {
	switch {
	case a.opts.Format != nil, a.opts.Kind == AxisCategory, a.opts.Kind == AxisTime:
		return a.Label(v)
	case a.opts.Kind == AxisLog:
		return strconv.FormatFloat(v, 'g', 3, 64)
	case a.step > 0:
		return internal.FormatTick(v, a.step/10)
	}
	return strconv.FormatFloat(v, 'g', 6, 64)
}

// Draw paints the axis line and tick labels on its side of plot.
func (a ChartAxisPart) Draw(p *ui.Painter, plot ui.Rect) {
	pos := make([]float32, len(a.ticks))
	labels := make([]string, len(a.ticks))
	for i, t := range a.ticks {
		pos[i], labels[i] = a.Map(plot, t), a.Label(t)
	}
	k := a.k
	switch a.opts.Side {
	case SideBottom:
		internal.Rule(p, plot.X, plot.X+plot.W, plot.Y+plot.H, k.ControlBorder)
		internal.LabelsX(p, plot, plot.Y+plot.H+6, pos, labels, false, a.size, k)
	case SideTop:
		internal.Rule(p, plot.X, plot.X+plot.W, plot.Y, k.ControlBorder)
		internal.LabelsX(p, plot, plot.Y-6, pos, labels, true, a.size, k)
	case SideLeft:
		internal.RuleV(p, plot.X, plot.Y, plot.Y+plot.H, k.ControlBorder)
		internal.LabelsY(p, plot.X-6, pos, labels, false, a.size, k)
	case SideRight:
		internal.RuleV(p, plot.X+plot.W, plot.Y, plot.Y+plot.H, k.ControlBorder)
		internal.LabelsY(p, plot.X+plot.W+6, pos, labels, true, a.size, k)
	}
}
