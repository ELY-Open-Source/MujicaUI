package chart

import (
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// VolumeChartOptions configure a VolumeChart.
type VolumeChartOptions struct {
	Frame ChartFrame
	// Sync links the crosshair to other charts under SyncID.
	Sync   *ChartSync
	SyncID string
}

// VolumeChart paints visible volume bars using the shared MarketView.
func VolumeChart(c *ui.Context, candles []Candle, view *MarketView, opts VolumeChartOptions) ChartResult {
	if view == nil {
		panic("mujica: VolumeChart needs a view")
	}
	if opts.Sync != nil && opts.SyncID == "" {
		panic("mujica: VolumeChart Sync needs a SyncID")
	}
	if opts.Frame.Height == 0 {
		opts.Frame.Height = 120
	}
	ch := NewChart(c, opts.Frame)
	ChartMemo(ch, "check", ChartID(candles), func() bool { marketCheckCandles(candles); return true })
	if ch.Res.Element.IsDisabled() {
		copy := *view
		view = &copy
		opts.Sync = nil
	}
	m := &mktCore{c: c, k: core.Tokens(c), ch: ch, candles: candles, view: view}
	if s := opts.Sync; s != nil {
		s.load(opts.SyncID, candles, view)
	}
	m.window()
	hi := ChartMemo(ch, "vmax", [3]any{ChartID(candles), m.from, m.to}, func() float64 {
		v := 0.0
		for i := m.from; i < m.to; i++ {
			v = max(v, candles[i].Volume)
		}
		return v
	})
	k := m.k
	size := core.FontSize(c, theme.CaptionSize)
	geom := func(r ui.Rect) (ui.Rect, internal.Band, internal.Linear, []float64, float64) {
		plot := ChartInset(r, 60, 4)
		lo, top, step := internal.Nice(0, max(hi, 1), 2)
		return plot, internal.Band{N: m.to - m.from, R0: plot.X, R1: plot.X + plot.W, Pad: 0.3},
			internal.Linear{D0: lo, D1: top, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}, internal.Ticks(lo, top, step), step
	}
	spec := ChartSpec{Marks: len(candles)}
	if m.to <= m.from {
		spec.Marks = 0
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, band, _, _, _ := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		return max(m.from, min(m.to-1, m.from+int((x-band.R0)/band.Step())))
	}
	spec.Readout = func(i int) string {
		cd := candles[i]
		return strings.Join([]string{cd.Time.Format("2006-01-02 15:04"), core.Msg(c, "candlestick_chart.volume") + " " + data.StatisticFormat(cd.Volume, 0),
			core.TrendWord(c, core.TrendSign(cd.Close-cd.Open))}, " · ")
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, band, ys, ticks, step := geom(r)
		internal.AxisY(p, plot, ys, ticks, func(v float64) string { return volumeShort(v, step) }, size, k)
		w := max(1, band.Width())
		at := active
		if at < 0 && opts.Sync != nil {
			at = opts.Sync.cursorIndex(opts.SyncID, candles)
		}
		for i := m.from; i < m.to; i++ {
			cd := candles[i]
			x := internal.Snap(p, band.Center(i-m.from))
			y := ys.Map(cd.Volume)
			bar := ui.Rect{X: x - w/2, Y: y, W: w, H: plot.Y + plot.H - y}
			up := cd.Close >= cd.Open
			col := core.TrendColor(c, 1)
			if !up {
				col = core.TrendColor(c, -1)
			}
			if up && w >= 3 {
				p.Fill(bar, k.Surface, 0)
				p.Stroke(bar, col, 0, 1)
			} else {
				p.Fill(bar, col, 0)
			}
		}
		m.st0().drawn = m.to - m.from
		if at >= m.from && at < m.to {
			x := internal.Snap(p, band.Center(at-m.from))
			p.Line(x, plot.Y, x, plot.Y+plot.H, 1, k.TextMuted)
		}
	}
	res, plot := ch.Build(spec)
	if plot != nil && !plot.IsDisabled() {
		room := m.room()
		switch {
		case plot.Shortcut(ui.Shift, ui.KeyLeft):
			marketPan(view, len(candles), room, -max(1, (m.to-m.from)/10))
		case plot.Shortcut(ui.Shift, ui.KeyRight):
			marketPan(view, len(candles), room, max(1, (m.to-m.from)/10))
		}
		if f := ch.st.focus; f >= 0 && (f < m.from || f >= m.to) {
			ch.st.focus = max(m.from, min(m.to-1, f))
		}
	}
	if s := opts.Sync; s != nil {
		s.store(opts.SyncID, candles, view, res.Point)
	}
	return res
}

// st0 is the chart's market state, made on first use.
func (m *mktCore) st0() *mktState {
	if m.st == nil {
		m.st = *ui.Local(m.ch.Res.Element, "mujica.mkt", func() *mktState { return &mktState{grab: -1, focus: -1} })
	}
	return m.st
}

// volumeShort writes a volume tick as 1.2K, 3M or 4B.
func volumeShort(v, step float64) string {
	for _, u := range []struct {
		div  float64
		unit string
	}{{1e9, "B"}, {1e6, "M"}, {1e3, "K"}} {
		if math.Abs(v) >= u.div {
			return internal.FormatTick(v/u.div, step/u.div) + u.unit
		}
	}
	return internal.FormatTick(v, step)
}
