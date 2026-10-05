package chart

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// paint draws axes, visible candles, layers, and the crosshair.
func (m *mktCore) paint(p *ui.Painter, r ui.Rect, active int) {
	g := m.geom(r)
	m.st.plot, m.st.ys = g.plot, g.ys
	k := m.k
	size := core.FontSize(m.c, theme.CaptionSize)
	label := func(v float64) string {
		if m.base != 0 {
			return internal.FormatTick(v, g.step) + "%"
		}
		return internal.FormatTick(v, g.step)
	}
	internal.AxisY(p, g.plot, g.ys, g.ticks, label, size, k)
	m.paintTimeAxis(p, g, size)
	p.Clip(g.plot, 0, func() {
		m.st.drawn = m.paintCandles(p, g)
		for _, l := range m.overlays {
			mktPaintLine(p, g, m, l.values, l.color, l.shape)
		}
		m.paintAlerts(p, g, size)
		if d := m.opts.Drawings; d != nil {
			d.paint(p, g, m, size)
		}
	})
	at := active
	if at < 0 && m.opts.Sync != nil {
		at = m.opts.Sync.cursorIndex(m.opts.SyncID, m.candles)
	}
	if at >= m.from && at < m.to {
		x := internal.Snap(p, g.x(m, at))
		y := internal.Snap(p, g.ys.Map(m.val(m.candles[at].Close)))
		p.Line(x, g.plot.Y, x, g.plot.Y+g.plot.H, 1, k.TextMuted)
		p.Line(g.plot.X, y, g.plot.X+g.plot.W, y, 1, k.TextMuted.Alpha(0.6))
		mktTag(p, g.plot.X+g.plot.W, y, label(m.val(m.candles[at].Close)), size, k.Text, k.Surface)
	}
}

// paintTimeAxis labels candles about 90 DIPs apart.
func (m *mktCore) paintTimeAxis(p *ui.Painter, g mktGeom, size float32) {
	every := max(1, int(math.Ceil(float64(90/g.band.Step()))))
	layout := core.Msg(m.c, "candlestick_chart.daily_time")
	if m.to-m.from > 1 && m.candles[m.from+1].Time.Sub(m.candles[m.from].Time) < 24*time.Hour {
		layout = "15:04"
	}
	var xs []float32
	var labels []string
	for i := m.from; i < m.to; i += every {
		xs = append(xs, g.x(m, i))
		labels = append(labels, m.candles[i].Time.Format(layout))
	}
	internal.AxisX(p, g.plot, xs, labels, size, m.k)
}

// paintCandles returns the number of visible candles painted.
func (m *mktCore) paintCandles(p *ui.Painter, g mktGeom) int {
	k := m.k
	w := max(1, g.band.Width())
	if m.opts.Type == ChartLine || m.opts.Type == ChartArea {
		var line, fill ui.Path
		for i := m.from; i < m.to; i++ {
			x, y := internal.Snap(p, g.x(m, i)), internal.Snap(p, g.ys.Map(m.val(m.candles[i].Close)))
			if i == m.from {
				line.MoveTo(x, y)
				fill.MoveTo(x, g.plot.Y+g.plot.H).LineTo(x, y)
			} else {
				line.LineTo(x, y)
				fill.LineTo(x, y)
			}
		}
		if m.opts.Type == ChartArea {
			fill.LineTo(internal.Snap(p, g.x(m, m.to-1)), g.plot.Y+g.plot.H).Close()
			p.FillPath(&fill, k.AccentText.Alpha(0.18))
		}
		p.StrokePath(&line, 2, k.AccentText)
		return m.to - m.from
	}
	for i := m.from; i < m.to; i++ {
		c := m.candles[i]
		up := c.Close >= c.Open
		col := core.TrendColor(m.c, 1)
		if !up {
			col = core.TrendColor(m.c, -1)
		}
		x := internal.Snap(p, g.x(m, i))
		yo, yc := g.ys.Map(m.val(c.Open)), g.ys.Map(m.val(c.Close))
		yh, yl := g.ys.Map(m.val(c.High)), g.ys.Map(m.val(c.Low))
		if m.opts.Type == ChartOHLC {
			p.Line(x, yh, x, yl, 1.5, col)
			p.Line(x-w/2, yo, x, yo, 1.5, col)
			p.Line(x, yc, x+w/2, yc, 1.5, col)
			continue
		}
		top, bot := min(yo, yc), max(yo, yc)
		body := ui.Rect{X: x - w/2, Y: top, W: w, H: max(1, bot-top)}
		p.Line(x, yh, x, top, 1, col)
		p.Line(x, bot, x, yl, 1, col)
		// Rises are hollow; falls are filled.
		if up && w >= 3 {
			p.Fill(body, k.Surface, 0)
			p.Stroke(body, col, 0, 1)
		} else {
			p.Fill(body, col, 0)
		}
	}
	return m.to - m.from
}

// mktPaintLine draws visible values with distinct series markers.
func mktPaintLine(p *ui.Painter, g mktGeom, m *mktCore, values []float64, col ui.Color, shape int) {
	var path ui.Path
	open := false
	for i := m.from; i < m.to; i++ {
		v := values[i]
		if math.IsNaN(v) {
			open = false
			continue
		}
		x, y := internal.Snap(p, g.x(m, i)), internal.Snap(p, g.ys.Map(v))
		if open {
			path.LineTo(x, y)
		} else {
			path.MoveTo(x, y)
			open = true
		}
	}
	p.StrokePath(&path, 1.5, col)
	every := max(1, int(60/max(g.band.Step(), 0.1)))
	for i := m.from + shape%every; i < m.to; i += every {
		if v := values[i]; !math.IsNaN(v) {
			internal.Marker(p, shape, g.x(m, i), g.ys.Map(v), 6, col)
		}
	}
}

// mktTag centers a price label on the right edge.
func mktTag(p *ui.Painter, right, y float32, s string, size float32, fg, bg ui.Color) {
	w, h := p.MeasureText(0, ui.Span{Text: s, Size: size})
	b := ui.Rect{X: right - w - 8, Y: y - h/2 - 2, W: w + 8, H: h + 4}
	p.Fill(b, bg, theme.ControlRadius)
	p.Stroke(b, fg, theme.ControlRadius, 1)
	p.Text(b.X+4, b.Y+2, s, size, fg)
}
