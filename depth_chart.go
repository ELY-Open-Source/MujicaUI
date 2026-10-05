package mujica

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// DepthLevel is one price level of an order book.
type DepthLevel struct {
	Price, Size float64
}

// DepthChartOptions configure a DepthChart.
type DepthChartOptions struct {
	Frame    ChartFrame
	Decimals int
}

// depthData is the book's cumulative sizes, bids then asks, from the mid outward.
type depthData struct {
	price, total []float64
	nBids        int
	mid, maxSize float64
}

// DepthChart draws cumulative bid and ask depth as stepped areas around the mid price; arrows walk the levels.
func DepthChart(c *ui.Context, bids, asks []DepthLevel, opts DepthChartOptions) ChartResult {
	dec := priceDecimals("DepthChart", opts.Decimals)
	ch := newChart(c, opts.Frame)
	k := Tokens(c)
	d := chartMemo(ch, "depth", [2]string{chartID(bids), chartID(asks)}, func() depthData { return depthBuild(bids, asks) })
	bidCol, askCol := trendColor(c, 1), trendColor(c, -1)
	keys := []chartKey{{Name: msg(c, "depth_chart.bids"), Color: bidCol, Shape: 1}, {Name: msg(c, "depth_chart.asks"), Color: askCol, Shape: 3}}
	size := FontSize(c, theme.CaptionSize)
	geom := func(r ui.Rect) (ui.Rect, internal.Linear, internal.Linear) {
		plot := chartInset(r, 56, 26)
		lo, hi := d.price[d.nBids-1], d.price[len(d.price)-1]
		return plot, internal.Linear{D0: lo, D1: hi, R0: float64(plot.X), R1: float64(plot.X + plot.W)},
			internal.Linear{D0: 0, D1: d.maxSize * 1.05, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	spec := chartSpec{n: len(d.price), keys: keys}
	if len(bids) == 0 || len(asks) == 0 {
		spec.n = 0
	}
	// Left and Right walk by price: lower bids to the left, higher asks to the right.
	order := func(i int) int {
		if i < d.nBids {
			return d.nBids - 1 - i
		}
		return i
	}
	spec.move = func(i int, key ui.Key) int {
		pos := order(i)
		if key == ui.KeyLeft || key == ui.KeyUp {
			pos = max(pos-1, 0)
		} else {
			pos = min(pos+1, len(d.price)-1)
		}
		return order(pos)
	}
	spec.pick = func(r ui.Rect, x, y float32) int {
		plot, xs, _ := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		return internal.Nearest(len(d.price), func(i int) (float32, float32) { return xs.Map(d.price[i]), 0 }, x, 0, plot.W, true)
	}
	spec.readout = func(i int) string {
		side := msg(c, "depth_chart.bid")
		if i >= d.nBids {
			side = msg(c, "depth_chart.ask")
		}
		return strings.Join([]string{side + " " + priceUngrouped(d.price[i], dec), fmt.Sprintf(msg(c, "depth_chart.total"), statisticFormat(d.total[i], 0)),
			msg(c, "depth_chart.mid") + " " + priceUngrouped(d.mid, dec)}, " · ")
	}
	spec.draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, xs, ys := geom(r)
		_, _, ystep := internal.Nice(0, d.maxSize*1.05, 4)
		internal.AxisY(p, plot, ys, internal.Ticks(0, d.maxSize*1.05, ystep), func(v float64) string { return volumeShort(v, ystep) }, size, k)
		_, _, xstep := internal.Nice(xs.D0, xs.D1, max(2, int(plot.W/90)))
		var tx []float32
		var tl []string
		for _, t := range internal.Ticks(xs.D0, xs.D1, xstep) {
			tx, tl = append(tx, xs.Map(t)), append(tl, internal.FormatTick(t, xstep))
		}
		internal.AxisX(p, plot, tx, tl, size, k)
		base := plot.Y + plot.H
		for side, idx := range [][2]int{{0, d.nBids}, {d.nBids, len(d.price)}} {
			col := []ui.Color{bidCol, askCol}[side]
			var area, edge ui.Path
			mx := internal.Snap(p, xs.Map(d.mid))
			area.MoveTo(mx, base)
			edge.MoveTo(mx, base)
			prevY := base
			for i := idx[0]; i < idx[1]; i++ {
				x, y := internal.Snap(p, xs.Map(d.price[i])), internal.Snap(p, ys.Map(d.total[i]))
				area.LineTo(x, prevY).LineTo(x, y)
				edge.LineTo(x, prevY).LineTo(x, y)
				prevY = y
			}
			end := plot.X
			if side == 1 {
				end = plot.X + plot.W
			}
			area.LineTo(end, prevY).LineTo(end, base).Close()
			edge.LineTo(end, prevY)
			p.FillPath(&area, col.Alpha(0.2))
			p.StrokePath(&edge, 2, col)
			label := []string{msg(c, "depth_chart.bids"), msg(c, "depth_chart.asks")}[side]
			chartText(p, label, []float32{plot.X + 6, plot.X + plot.W - 6}[side], plot.Y+4, size, col, []float32{0, 1}[side], 0)
		}
		mx := internal.Snap(p, xs.Map(d.mid))
		for y := plot.Y; y < base; y += 6 {
			p.Line(mx, y, mx, min(y+3, base), 1, k.TextMuted)
		}
		chartText(p, priceUngrouped(d.mid, dec), mx, plot.Y+4, size, k.Text, 0.5, 0)
		if active >= 0 {
			x, y := xs.Map(d.price[active]), ys.Map(d.total[active])
			p.Line(internal.Snap(p, x), plot.Y, internal.Snap(p, x), base, 1, k.TextMuted)
			shape := 1
			if active >= d.nBids {
				shape = 3
			}
			internal.Marker(p, shape, x, y, 10, keys[shape/3].Color)
		}
	}
	res, _ := ch.build(spec)
	return res
}

// depthBuild checks the book's order and sums sizes outward from the mid.
func depthBuild(bids, asks []DepthLevel) depthData {
	for i, l := range bids {
		if !(l.Size > 0) || i > 0 && !(l.Price < bids[i-1].Price) {
			panic(fmt.Sprintf("mujica: DepthChart bid %d needs a positive size and a price below the one before", i))
		}
	}
	for i, l := range asks {
		if !(l.Size > 0) || i > 0 && !(l.Price > asks[i-1].Price) {
			panic(fmt.Sprintf("mujica: DepthChart ask %d needs a positive size and a price above the one before", i))
		}
	}
	if len(bids) > 0 && len(asks) > 0 && !(bids[0].Price < asks[0].Price) {
		panic("mujica: DepthChart best bid must be below best ask")
	}
	var d depthData
	d.nBids = len(bids)
	var sum float64
	for _, l := range bids {
		sum += l.Size
		d.price, d.total = append(d.price, l.Price), append(d.total, sum)
	}
	d.maxSize = sum
	sum = 0
	for _, l := range asks {
		sum += l.Size
		d.price, d.total = append(d.price, l.Price), append(d.total, sum)
	}
	d.maxSize = max(d.maxSize, sum)
	if len(bids) > 0 && len(asks) > 0 {
		d.mid = (bids[0].Price + asks[0].Price) / 2
	}
	return d
}
