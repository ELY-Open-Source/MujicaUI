package mujica

import (
	"fmt"
	"log/slog"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// IndicatorPaneSpec is one lower pane: MACD, RSI, KDJ or OBV.
type IndicatorPaneSpec struct {
	Kind IndicatorKind
	// Params nil uses the kind's defaults.
	Params []float64
}

// IndicatorPanes is the caller-held stack of panes and their height shares.
type IndicatorPanes struct {
	Panes []IndicatorPaneSpec
	// Sizes are the panes' height shares; they reset to equal when the pane count changes.
	Sizes []float32
}

// IndicatorPaneOptions configure an IndicatorPane.
type IndicatorPaneOptions struct {
	// Height is the whole stack's height; zero is 140 per pane.
	Height float32
	// Label prefixes each pane's accessible name, as the symbol.
	Label    string
	Disabled bool
}

// IndicatorPaneResult reports the stack's events.
type IndicatorPaneResult struct {
	Element *ui.Element
	// Removed is the index of the pane closed this frame, -1 for none.
	Removed int
	// Resized reports finished height changes, for saving.
	Resized bool
}

// IndicatorPane stacks indicator panes under a price chart, sharing its view; dividers between them resize.
func IndicatorPane(c *ui.Context, candles []Candle, view *MarketView, panes *IndicatorPanes, opts IndicatorPaneOptions) IndicatorPaneResult {
	if view == nil || panes == nil {
		panic("mujica: IndicatorPane needs a view and panes")
	}
	res := IndicatorPaneResult{Removed: -1}
	n := len(panes.Panes)
	for _, p := range panes.Panes {
		if _ = p.Kind.String(); indicatorInfo[p.Kind].overlay {
			panic(fmt.Sprintf("mujica: %v is an overlay, not a pane indicator", p.Kind))
		}
		indicatorParams(p.Kind, p.Params)
	}
	if n == 0 {
		res.Element = ui.Column(c)
		return res
	}
	if len(panes.Sizes) != n {
		panes.Sizes = nil
	}
	h := opts.Height
	if h == 0 {
		h = 140 * float32(n)
	}
	build := func(i int) func() {
		return func() {
			col := ui.Column(c).Fill().MinHeight(0)
			b := col.Bounds()
			internal.Remeasure(c, col, b.W, b.H)
			u := Density(c).Unit()
			// The plot takes what the title, legend and readout rows leave.
			plotH := max(40, b.H-2*theme.MinHitSize-u*12)
			col.Children(func() {
				if indicatorPaneOne(c, candles, view, panes.Panes[i], opts, plotH) {
					res.Removed = i
				}
			})
		}
	}
	box := ui.Column(c).Height(h).FillWidth().Disabled(opts.Disabled)
	box.Children(func() {
		if n == 1 {
			build(0)()
			return
		}
		list := make([]ResizablePanel, n)
		for i, p := range panes.Panes {
			list[i] = ResizablePanel{Label: indicatorName(p.Kind, indicatorParams(p.Kind, p.Params)), Min: 60, Build: build(i)}
		}
		r := ResizablePanelGroup(c, &panes.Sizes, list, ResizablePanelGroupOptions{Vertical: true})
		r.Element.Fill()
		res.Resized = r.Committed
	})
	if res.Removed >= 0 {
		slog.Debug("mujica: indicator_pane removed", "kind", panes.Panes[res.Removed].Kind)
		panes.Panes = append(panes.Panes[:res.Removed:res.Removed], panes.Panes[res.Removed+1:]...)
		panes.Sizes = nil
		c.Invalidate()
	}
	res.Element = box
	return res
}

// indicatorPaneOne draws one pane's lines over the view's window; it reports the close button.
func indicatorPaneOne(c *ui.Context, candles []Candle, view *MarketView, spec IndicatorPaneSpec, opts IndicatorPaneOptions, height float32) bool {
	k := Tokens(c)
	params := indicatorParams(spec.Kind, spec.Params)
	name := indicatorName(spec.Kind, params)
	closed := false
	ui.Row(c).AlignItems(ui.Center).Gap(Density(c).Unit()).Children(func() {
		ui.Text(c, name).FontSize(FontSize(c, theme.CaptionSize)).FontWeight(600).TextColor(k.Text).Grow(1)
		closed = iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "indicator_pane.remove"), name)).Clicked()
	})
	label := strings.TrimSpace(opts.Label + " " + name)
	ch := newChart(c, ChartFrame{Label: label, Height: height})
	series := chartMemo(ch, "series", fmt.Sprintf("%s|%v|%v", chartID(candles), spec.Kind, params), func() []ChartSeries {
		marketCheckCandles(candles)
		return indicatorSeries(spec.Kind, params, candles)
	})
	m := &mktCore{c: c, k: k, ch: ch, candles: candles, view: view}
	m.window()
	pal := internal.Palette(k)
	lines := make([]mktLine, len(series))
	keys := make([]chartKey, len(series))
	for i, s := range series {
		lines[i] = mktLine{name: s.Name, values: s.Values, color: pal[(i+1)%len(pal)], shape: i}
		keys[i] = chartKey{Name: s.Name, Color: lines[i].color, Shape: i}
	}
	yr := chartMemo(ch, "y", [3]any{chartID(candles), m.from, m.to}, func() [2]float64 {
		lo, hi := math.Inf(1), math.Inf(-1)
		for _, l := range lines {
			for i := m.from; i < m.to; i++ {
				if v := l.values[i]; !math.IsNaN(v) {
					lo, hi = min(lo, v), max(hi, v)
				}
			}
		}
		if spec.Kind == IndicatorRSI {
			lo, hi = 0, 100
		}
		return [2]float64{lo, hi}
	})
	lo, hi := yr[0], yr[1]
	spec2 := chartSpec{n: len(candles), keys: keys}
	if m.to <= m.from || math.IsInf(lo, 0) {
		spec2.n = 0
	}
	m.lo, m.hi = lo, hi
	spec2.pick = func(r ui.Rect, x, y float32) int {
		g := m.paneGeom(r)
		if !g.plot.Contains(x, y) {
			return -1
		}
		return g.index(m, x)
	}
	spec2.readout = func(i int) string {
		parts := []string{candles[i].Time.Format("2006-01-02 15:04")}
		for _, l := range lines {
			parts = append(parts, l.name+" "+chartValue(c, l.values[i], func(v float64) string { return fmt.Sprintf("%.2f", v) }))
		}
		return strings.Join(parts, " · ")
	}
	spec2.draw = func(p *ui.Painter, r ui.Rect, active int) {
		g := m.paneGeom(r)
		size := FontSize(c, theme.CaptionSize)
		internal.AxisY(p, g.plot, g.ys, g.ticks, func(v float64) string { return volumeShort(v, g.step) }, size, k)
		p.Clip(g.plot, 0, func() {
			for i, l := range lines {
				if spec.Kind == IndicatorMACD && i == 2 {
					indicatorPaneBars(p, g, m, l.values, c)
					continue
				}
				mktPaintLine(p, g, m, l.values, l.color, i)
			}
		})
		if active >= m.from && active < m.to {
			x := internal.Snap(p, g.x(m, active))
			p.Line(x, g.plot.Y, x, g.plot.Y+g.plot.H, 1, k.TextMuted)
		}
	}
	ch.build(spec2)
	return closed && !ch.res.Element.IsDisabled()
}

// paneGeom is geom without the time axis, for lower panes.
func (m *mktCore) paneGeom(r ui.Rect) mktGeom {
	plot := chartInset(r, 60, 4)
	lo, hi, step := internal.Nice(m.lo, m.hi, 2)
	return mktGeom{plot: plot, band: internal.Band{N: m.to - m.from, R0: plot.X, R1: plot.X + plot.W, Pad: 0.3},
		ys: internal.Linear{D0: lo, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}, ticks: internal.Ticks(lo, hi, step), step: step}
}

// indicatorPaneBars draws MACD's histogram: hollow bars above zero, filled below.
func indicatorPaneBars(p *ui.Painter, g mktGeom, m *mktCore, vals []float64, c *ui.Context) {
	k := Tokens(c)
	zero := g.ys.Map(0)
	w := max(1, g.band.Width())
	for i := m.from; i < m.to; i++ {
		v := vals[i]
		if math.IsNaN(v) {
			continue
		}
		y := g.ys.Map(v)
		x := internal.Snap(p, g.x(m, i))
		bar := ui.Rect{X: x - w/2, Y: min(y, zero), W: w, H: max(1, float32(math.Abs(float64(y-zero))))}
		if v >= 0 && w >= 3 {
			p.Fill(bar, k.Surface, 0)
			p.Stroke(bar, trendColor(c, 1), 0, 1)
		} else {
			p.Fill(bar, trendColor(c, trendSign(v)), 0)
		}
	}
}
