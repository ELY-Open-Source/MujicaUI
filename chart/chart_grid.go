package chart

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
)

// ChartGridOptions configure a ChartGrid.
type ChartGridOptions struct {
	// X and Y provide grid ticks; either may be nil.
	X, Y *ChartAxisPart
	// Minor subdivides major intervals with fainter lines.
	Minor int
}

// ChartGridPart is a grid ready to draw in a plot.
type ChartGridPart struct {
	opts         ChartGridOptions
	major, minor ui.Color
}

// ChartGrid draws pixel-aligned major and minor lines.
func ChartGrid(c *ui.Context, opts ChartGridOptions) ChartGridPart {
	if opts.Minor < 0 || opts.X == nil && opts.Y == nil {
		panic("mujica: ChartGrid needs an axis and Minor ≥ 0")
	}
	for _, a := range []*ChartAxisPart{opts.X, opts.Y} {
		if a != nil && a.opts.Kind == AxisCategory && opts.Minor > 0 {
			panic("mujica: ChartGrid has no minor lines on a category axis")
		}
	}
	k := core.Tokens(c)
	return ChartGridPart{opts: opts, major: k.Border, minor: k.Border.Mix(k.Surface, 0.55)}
}

// Draw paints minor lines, then major ones, inside plot.
func (g ChartGridPart) Draw(p *ui.Painter, plot ui.Rect) {
	for _, a := range []*ChartAxisPart{g.opts.X, g.opts.Y} {
		if a == nil {
			continue
		}
		line := func(v float64, col ui.Color) {
			at := a.Map(plot, v)
			if a.horizontal() {
				internal.RuleV(p, at, plot.Y, plot.Y+plot.H, col)
			} else {
				internal.Rule(p, plot.X, plot.X+plot.W, at, col)
			}
		}
		ticks := a.drawTicks(plot, 4)
		for i := 1; i < len(ticks) && g.opts.Minor > 0; i++ {
			for j := 1; j < g.opts.Minor; j++ {
				f := float64(j) / float64(g.opts.Minor)
				v := ticks[i-1] + (ticks[i]-ticks[i-1])*f
				if a.opts.Kind == AxisLog {
					v = math.Pow(10, math.Log10(ticks[i-1])+(math.Log10(ticks[i])-math.Log10(ticks[i-1]))*f)
				}
				line(v, g.minor)
			}
		}
		for _, t := range ticks {
			line(t, g.major)
		}
	}
}
