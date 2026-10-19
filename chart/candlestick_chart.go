package chart

import (
	"fmt"
	"log/slog"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ChartType selects the price chart shape.
type ChartType int

const (
	ChartCandles ChartType = iota
	ChartOHLC
	ChartLine
	ChartArea
	ChartHeikinAshi
)

// CandlestickOptions configure a CandlestickChart and its layers.
type CandlestickOptions struct {
	Frame ChartFrame
	Type  ChartType
	// Decimals of prices; zero means 2.
	Decimals int
	// Overlays draw indicators over the prices.
	Overlays []IndicatorOverlay
	// Compare switches the chart to normalized percent prices.
	Compare *[]CompareSymbol
	// Alerts provides draggable lines and price-crossing callbacks.
	Alerts  *[]PriceAlertLine
	OnAlert func(a PriceAlertLine, last Candle)
	// Drawings, when set, takes clicks for the current tool.
	Drawings *DrawingTools
	// Sync links the window and crosshair under SyncID.
	Sync   *ChartSync
	SyncID string
}

// mktState is a price chart's own frame-to-frame state.
type mktState struct {
	drawn      int     // candles painted last frame
	dragRest   float32 // pan drag left over, under a candle
	lastClose  float64
	alertSeen  bool
	grab       int // alert index dragged, -1 none
	focus      int // keyboard mark last frame, to follow only when it moves
	alertFocus int
	plot       ui.Rect
	ys         internal.Linear // the price scale as last painted
}

// mktCore is one price chart being built.
type mktCore struct {
	c        *ui.Context
	k        theme.Tokens
	ch       *chart
	candles  []Candle
	view     *MarketView
	from, to int
	act      int // keyboard or click activation, mark plus one; zero is none
	lo, hi   float64
	base     float64 // percent mode divides by this close; 0 is prices
	dec      int
	opts     CandlestickOptions
	st       *mktState
	overlays []mktLine
}

// mktLine is one drawn indicator or compare line.
type mktLine struct {
	name   string
	values []float64
	color  ui.Color
	shape  int
}

// CandlestickChart paints visible candles with hollow rises and filled falls.
func CandlestickChart(c *ui.Context, candles []Candle, view *MarketView, opts CandlestickOptions) ChartResult {
	if view == nil {
		panic("mujica: CandlestickChart needs a view")
	}
	if opts.Type < ChartCandles || opts.Type > ChartHeikinAshi {
		panic(fmt.Sprintf("mujica: unknown ChartType %d", opts.Type))
	}
	if opts.Sync != nil && opts.SyncID == "" {
		panic("mujica: CandlestickChart Sync needs a SyncID")
	}
	ch := NewChart(c, opts.Frame)
	shown := ChartMemo(ch, "candles", [2]any{ChartID(candles), opts.Type}, func() []Candle {
		marketCheckCandles(candles)
		if opts.Type == ChartHeikinAshi {
			return HeikinAshi(candles)
		}
		return candles
	})
	m := &mktCore{c: c, k: core.Tokens(c), ch: ch, candles: shown, view: view, opts: opts, dec: core.PriceDecimals("CandlestickChart", opts.Decimals)}
	m.st = *ui.Local(ch.Res.Element, "mujica.mkt", func() *mktState { return &mktState{grab: -1, focus: -1, alertFocus: -1} })
	if ch.Res.Element.IsDisabled() {
		copy := *view
		view = &copy
		m.view = view
		m.opts.Sync = nil
		opts.Sync = nil
	}
	if s := opts.Sync; s != nil {
		s.load(opts.SyncID, shown, view)
	}
	m.window()
	m.overlays = mktOverlayLines(ch, m, candles)
	m.compareLines()
	m.yDomain()
	m.alertsCheck()
	res, plot := ch.Build(m.spec())
	if plot.Valid() {
		m.input(plot)
	}
	if s := opts.Sync; s != nil {
		s.store(opts.SyncID, shown, view, res.Point)
	}
	return res
}

// OHLCChart is CandlestickChart drawn as open-high-low-close bars.
func OHLCChart(c *ui.Context, candles []Candle, view *MarketView, opts CandlestickOptions) ChartResult {
	opts.Type = ChartOHLC
	return CandlestickChart(c, candles, view, opts)
}

// HeikinAshiChart switches between smoothed candles and raw OHLC.
func HeikinAshiChart(c *ui.Context, candles []Candle, view *MarketView, raw *bool, opts CandlestickOptions) ChartResult {
	if raw == nil {
		panic("mujica: HeikinAshiChart needs raw")
	}
	opts.Type = ChartHeikinAshi
	if *raw {
		opts.Type = ChartCandles
	}
	var res ChartResult
	col := ui.Column(c).Gap(core.Density(c).Unit() * 2).MinWidth(0)
	col.Children(func() {
		ui.Row(c).Children(func() {
			v := []bool{*raw}
			if input.ToggleGroup(c, &v, []input.ToggleItem[bool]{{Value: false, Label: core.Msg(c, "heikin_ashi_chart.smoothed")}, {Value: true, Label: core.Msg(c, "heikin_ashi_chart.raw")}},
				input.ToggleGroupOptions{Label: core.Msg(c, "heikin_ashi_chart.switch")}).Changed() {
				*raw = v[0]
				slog.Debug("mujica: heikin_ashi switch", "raw", *raw)
			}
		})
		res = CandlestickChart(c, candles, view, opts)
	})
	res.Element = col
	return res
}

// window resolves visible candles using the last measured width.
func (m *mktCore) window() {
	m.from, m.to = marketWindow(*m.view, len(m.candles), m.room())
}

// room is how many candles fit, three DIPs each.
func (m *mktCore) room() int {
	w := m.ch.Res.Element.Bounds().W - 64
	if w <= 0 {
		return 200
	}
	return max(10, int(w/3))
}

// yDomain caches visible prices, indicators, and alerts.
func (m *mktCore) yDomain() {
	if m.to <= m.from {
		return
	}
	var alerts []PriceAlertLine
	if m.opts.Alerts != nil {
		alerts = *m.opts.Alerts
	}
	ids := fmt.Sprintf("%s|%d|%d|%v|%d|%v", ChartID(m.candles), m.from, m.to, m.base, len(m.overlays), alerts)
	for _, l := range m.overlays {
		ids += ChartID(l.values)
	}
	d := ChartMemo(m.ch, "ydomain", ids, func() [2]float64 {
		lo, hi := math.Inf(1), math.Inf(-1)
		for i := m.from; i < m.to; i++ {
			lo, hi = min(lo, m.val(m.candles[i].Low)), max(hi, m.val(m.candles[i].High))
			for _, l := range m.overlays {
				if v := l.values[i]; !math.IsNaN(v) {
					lo, hi = min(lo, v), max(hi, v)
				}
			}
		}
		if m.base == 0 {
			for _, a := range alerts {
				lo, hi = min(lo, a.Price), max(hi, a.Price)
			}
		}
		return [2]float64{lo, hi}
	})
	m.lo, m.hi = d[0], d[1]
}

// val converts price to the current axis unit.
func (m *mktCore) val(p float64) float64 {
	if m.base == 0 {
		return p
	}
	return (p/m.base - 1) * 100
}

// mktGeom places candles and prices in a plot.
type mktGeom struct {
	plot  ui.Rect
	band  internal.Band
	ys    internal.Linear
	ticks []float64
	step  float64
}

func (m *mktCore) geom(r ui.Rect) mktGeom {
	plot := ChartInset(r, 60, 26)
	lo, hi, step := internal.Nice(m.lo, m.hi, 4)
	return mktGeom{
		plot:  plot,
		band:  internal.Band{N: m.to - m.from, R0: plot.X, R1: plot.X + plot.W, Pad: 0.3},
		ys:    internal.Linear{D0: lo, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)},
		ticks: internal.Ticks(lo, hi, step),
		step:  step,
	}
}

func (g mktGeom) x(m *mktCore, i int) float32 { return g.band.Center(i - m.from) }

// index is the candle under x, clamped to the window.
func (g mktGeom) index(m *mktCore, x float32) int {
	i := m.from + int((x-g.band.R0)/g.band.Step())
	return max(m.from, min(m.to-1, i))
}

func (m *mktCore) spec() ChartSpec {
	n := len(m.candles)
	var keys []ChartLegendItem
	for _, l := range m.overlays {
		keys = append(keys, ChartLegendItem{Name: l.name, Color: l.color, Shape: l.shape})
	}
	if n == 0 && m.opts.Compare != nil {
		pal := internal.Palette(m.k)
		for i, s := range *m.opts.Compare {
			keys = append(keys, ChartLegendItem{Name: s.Name, Color: pal[(i+1)%len(pal)], Shape: i + 1})
		}
	}
	spec := ChartSpec{Marks: n, Legend: keys}
	if m.to <= m.from {
		spec.Marks = 0
		return spec
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		g := m.geom(r)
		if !g.plot.Contains(x, y) {
			return -1
		}
		return g.index(m, x)
	}
	spec.Move = func(i int, key ui.Key) int {
		if key == ui.KeyLeft || key == ui.KeyUp {
			return max(i-1, 0)
		}
		return min(i+1, n-1)
	}
	if d := m.opts.Drawings; d != nil {
		spec.Activate = func(i int) {
			m.act = i + 1
		}
	}
	spec.Readout = m.readout
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) { m.paint(p, r, active) }
	return spec
}

// readout describes time, OHLC, volume, direction, and indicators.
func (m *mktCore) readout(i int) string {
	c := m.candles[i]
	f := func(v float64) string { return core.PriceUngrouped(v, m.dec) }
	prev := c.Open
	if i > 0 {
		prev = m.candles[i-1].Close
	}
	sign := core.TrendSign(c.Close - prev)
	pct := 0.0
	if prev != 0 {
		pct = (c.Close/prev - 1) * 100
	}
	parts := []string{
		c.Time.Format("2006-01-02 15:04"),
		fmt.Sprintf(core.Msg(m.c, "candlestick_chart.ohlc"), f(c.Open), f(c.High), f(c.Low), f(c.Close)),
		fmt.Sprintf("%+.2f%% %s", pct, core.TrendWord(m.c, sign)),
		core.Msg(m.c, "candlestick_chart.volume") + " " + data.StatisticFormat(c.Volume, 0),
	}
	for _, l := range m.overlays {
		parts = append(parts, l.name+" "+ChartValue(m.c, l.values[i], f))
	}
	return strings.Join(parts, " · ")
}
