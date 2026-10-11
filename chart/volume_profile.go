package chart

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ProfileBin is the volume traded between two prices.
type ProfileBin struct {
	Low, High, Volume float64
}

// VolumeProfileOf spreads candle volumes across equal price bins.
func VolumeProfileOf(candles []Candle, n int) []ProfileBin {
	if n < 1 {
		panic("mujica: VolumeProfileOf needs at least one bin")
	}
	if len(candles) == 0 {
		return nil
	}
	marketCheckCandles(candles)
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, c := range candles {
		lo, hi = min(lo, c.Low), max(hi, c.High)
	}
	if hi == lo {
		hi = lo + 1
	}
	w := (hi - lo) / float64(n)
	bins := make([]ProfileBin, n)
	for i := range bins {
		bins[i] = ProfileBin{Low: lo + float64(i)*w, High: lo + float64(i+1)*w}
	}
	for _, c := range candles {
		span := c.High - c.Low
		for i := range bins {
			if span == 0 {
				if c.Low >= bins[i].Low && (c.Low < bins[i].High || i == n-1) {
					bins[i].Volume += c.Volume
					break
				}
				continue
			}
			overlap := min(c.High, bins[i].High) - max(c.Low, bins[i].Low)
			if overlap > 0 {
				bins[i].Volume += c.Volume * overlap / span
			}
		}
	}
	return bins
}

// ValueArea finds POC and its contiguous volume fraction.
func ValueArea(bins []ProfileBin, frac float64) (poc, lo, hi int) {
	if len(bins) == 0 || !(frac > 0 && frac <= 1) {
		panic("mujica: ValueArea needs bins and a fraction in (0, 1]")
	}
	var total float64
	for i, b := range bins {
		if b.Volume < 0 || math.IsNaN(b.Volume) || math.IsInf(b.Volume, 0) {
			panic("mujica: profile volumes must be finite and non-negative")
		}
		total += b.Volume
		if b.Volume > bins[poc].Volume {
			poc = i
		}
	}
	lo, hi = poc, poc
	in := bins[poc].Volume
	for in < frac*total {
		// Grow toward the fuller neighbor, as the market-profile rule does.
		down, up := -1.0, -1.0
		if lo > 0 {
			down = bins[lo-1].Volume
		}
		if hi < len(bins)-1 {
			up = bins[hi+1].Volume
		}
		if down < 0 && up < 0 {
			break
		}
		if up >= down {
			hi++
			in += up
		} else {
			lo--
			in += down
		}
	}
	return poc, lo, hi
}

// VolumeProfileOptions configure a VolumeProfile.
type VolumeProfileOptions struct {
	Frame ChartFrame
	// ValueArea defaults to 0.7.
	ValueArea float64
	Decimals  int
}

// VolumeProfile draws volume distribution, value area, and POC.
func VolumeProfile(c *ui.Context, bins []ProfileBin, opts VolumeProfileOptions) ChartResult {
	frac := opts.ValueArea
	if frac == 0 {
		frac = 0.7
	}
	if !(frac > 0 && frac <= 1) {
		panic("mujica: VolumeProfile ValueArea must be in (0, 1]")
	}
	dec := core.PriceDecimals("VolumeProfile", opts.Decimals)
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	va := ChartMemo(ch, "va", [2]any{ChartID(bins), frac}, func() [4]float64 {
		vmax := 0.0
		for i, b := range bins {
			if math.IsNaN(b.Volume) || math.IsInf(b.Volume, 0) || math.IsInf(b.High, 0) || math.IsInf(b.Low, 0) || !(b.High > b.Low) || b.Volume < 0 || i > 0 && b.Low < bins[i-1].High-1e-9 {
				panic(fmt.Sprintf("mujica: VolumeProfile bin %d is empty, negative or out of order", i))
			}
			vmax = max(vmax, b.Volume)
		}
		if len(bins) == 0 {
			return [4]float64{}
		}
		poc, lo, hi := ValueArea(bins, frac)
		return [4]float64{float64(poc), float64(lo), float64(hi), vmax}
	})
	poc, vlo, vhi, vmax := int(va[0]), int(va[1]), int(va[2]), va[3]
	inArea, outArea := k.AccentText, k.TextMuted
	keys := []ChartLegendItem{{Name: core.Msg(c, "volume_profile.value_area"), Color: inArea, Shape: 1}, {Name: core.Msg(c, "volume_profile.other"), Color: outArea, Shape: 4}}
	drawn := ui.Local(ch.Res.Element, "mujica.profile.drawn", func() int { return 0 })
	size := core.FontSize(c, theme.CaptionSize)
	geom := func(r ui.Rect) (ui.Rect, internal.Linear, internal.Linear) {
		plot := ChartInset(r, 60, 26)
		return plot, internal.Linear{D0: 0, D1: max(vmax, 1), R0: float64(plot.X), R1: float64(plot.X + plot.W)},
			internal.Linear{D0: bins[0].Low, D1: bins[len(bins)-1].High, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	spec := ChartSpec{Marks: len(bins), Legend: keys}
	spec.Move = func(i int, key ui.Key) int {
		if key == ui.KeyDown || key == ui.KeyLeft {
			return max(i-1, 0)
		}
		return min(i+1, len(bins)-1)
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, _, ys := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		v := ys.Invert(y)
		i := sort.Search(len(bins), func(i int) bool { return bins[i].High >= v })
		if i < len(bins) && v >= bins[i].Low {
			return i
		}
		return -1
	}
	spec.Readout = func(i int) string {
		b := bins[i]
		parts := []string{core.PriceUngrouped(b.Low, dec) + "–" + core.PriceUngrouped(b.High, dec), core.Msg(c, "candlestick_chart.volume") + " " + data.StatisticFormat(b.Volume, 0)}
		if i == poc {
			parts = append(parts, core.Msg(c, "volume_profile.poc"))
		}
		if i >= vlo && i <= vhi {
			parts = append(parts, core.Msg(c, "volume_profile.value_area"))
		}
		return strings.Join(parts, " · ")
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, xs, ys := geom(r)
		_, _, step := internal.Nice(ys.D0, ys.D1, 5)
		internal.AxisY(p, plot, ys, internal.Ticks(ys.D0, ys.D1, step), func(v float64) string { return internal.FormatTick(v, step) }, size, k)
		_, _, xstep := internal.Nice(0, xs.D1, 3)
		var tx []float32
		var tl []string
		for _, t := range internal.Ticks(0, xs.D1, xstep) {
			tx, tl = append(tx, xs.Map(t)), append(tl, volumeShort(t, xstep))
		}
		internal.AxisX(p, plot, tx, tl, size, k)
		indices := ChartMemo(ch, "profile.indices", [3]any{ChartID(bins), plot.Y, plot.H}, func() []int {
			return internal.Columns(len(bins), func(i int) float32 { return ys.Map(bins[i].Low) }, func(i int) float64 { return bins[i].Volume })
		})
		*drawn = len(indices)
		for _, i := range indices {
			b := bins[i]
			y0, y1 := ys.Map(b.High), ys.Map(b.Low)
			bar := ui.Rect{X: plot.X, Y: y0 + 1, W: xs.Map(b.Volume) - plot.X, H: max(1, y1-y0-2)}
			// Value-area bars are filled, the rest outlined.
			if i >= vlo && i <= vhi {
				p.Fill(bar, inArea, 0)
			} else {
				p.Stroke(bar, outArea, 0, 1)
			}
			if i == active {
				p.Stroke(ui.Rect{X: bar.X - 1, Y: bar.Y - 1, W: bar.W + 2, H: bar.H + 2}, k.Text, 0, 1.5)
			}
		}
		for _, line := range []struct {
			v   float64
			key string
		}{{(bins[poc].Low + bins[poc].High) / 2, "volume_profile.poc"}, {bins[vhi].High, "volume_profile.vah"}, {bins[vlo].Low, "volume_profile.val"}} {
			y := internal.Snap(p, ys.Map(line.v))
			p.Line(plot.X, y, plot.X+plot.W, y, 1, k.Warning)
			mktTag(p, plot.X+plot.W, y, core.Msg(c, line.key)+" "+core.PriceUngrouped(line.v, dec), size, k.Warning, k.Surface)
		}
	}
	res, _ := ch.Build(spec)
	return res
}
