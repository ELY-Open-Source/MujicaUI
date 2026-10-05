package chart

import (
	"fmt"
	"math"
)

// Indicator helpers are pure; NaN marks insufficient data.

// indicatorCheck panics on a period below one.
func indicatorCheck(name string, n int) {
	if n < 1 {
		panic(fmt.Sprintf("mujica: %s period %d must be at least 1", name, n))
	}
}

func indicatorNaN(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = math.NaN()
	}
	return out
}

// SMA computes the n-period simple moving average.
func SMA(vals []float64, n int) []float64 {
	indicatorCheck("SMA", n)
	out := indicatorNaN(len(vals))
	var sum float64
	for i, v := range vals {
		sum += v
		if i >= n {
			sum -= vals[i-n]
		}
		if i >= n-1 {
			out[i] = sum / float64(n)
		}
	}
	return out
}

// EMA uses 2/(n+1), seeded by the first n-period SMA.
func EMA(vals []float64, n int) []float64 {
	indicatorCheck("EMA", n)
	out := indicatorNaN(len(vals))
	k := 2 / float64(n+1)
	var sum, prev float64
	for i, v := range vals {
		switch {
		case i < n-1:
			sum += v
		case i == n-1:
			prev = (sum + v) / float64(n)
			out[i] = prev
		default:
			prev = v*k + prev*(1-k)
			out[i] = prev
		}
	}
	return out
}

// Bollinger returns SMA and population-deviation bands.
func Bollinger(vals []float64, n int, width float64) (mid, upper, lower []float64) {
	indicatorCheck("Bollinger", n)
	if !(width > 0) {
		panic("mujica: Bollinger width must be positive")
	}
	mid = SMA(vals, n)
	upper, lower = indicatorNaN(len(vals)), indicatorNaN(len(vals))
	for i := n - 1; i < len(vals); i++ {
		var ss float64
		for _, v := range vals[i-n+1 : i+1] {
			ss += (v - mid[i]) * (v - mid[i])
		}
		d := width * math.Sqrt(ss/float64(n))
		upper[i], lower[i] = mid[i]+d, mid[i]-d
	}
	return mid, upper, lower
}

// VWAP averages typical prices weighted by cumulative volume.
func VWAP(candles []Candle) []float64 {
	out := indicatorNaN(len(candles))
	var pv, vol float64
	for i, c := range candles {
		pv += (c.High + c.Low + c.Close) / 3 * c.Volume
		vol += c.Volume
		if vol > 0 {
			out[i] = pv / vol
		}
	}
	return out
}

// MACD returns DIF, signal EMA, and their difference.
func MACD(vals []float64, fast, slow, signal int) (dif, dea, hist []float64) {
	if fast >= slow {
		panic("mujica: MACD fast period must be shorter than slow")
	}
	f, s := EMA(vals, fast), EMA(vals, slow)
	dif = indicatorNaN(len(vals))
	for i := range vals {
		dif[i] = f[i] - s[i]
	}
	dea = indicatorNaN(len(vals))
	if len(vals) >= slow {
		tail := EMA(dif[slow-1:], signal)
		copy(dea[slow-1:], tail)
	}
	hist = indicatorNaN(len(vals))
	for i := range vals {
		hist[i] = dif[i] - dea[i]
	}
	return dif, dea, hist
}

// RSI computes Wilder's relative strength index.
func RSI(vals []float64, n int) []float64 {
	indicatorCheck("RSI", n)
	out := indicatorNaN(len(vals))
	var gain, loss float64
	for i := 1; i < len(vals); i++ {
		d := vals[i] - vals[i-1]
		g, l := max(d, 0), max(-d, 0)
		if i <= n {
			gain += g
			loss += l
			if i < n {
				continue
			}
			gain, loss = gain/float64(n), loss/float64(n)
		} else {
			gain = (gain*float64(n-1) + g) / float64(n)
			loss = (loss*float64(n-1) + l) / float64(n)
		}
		if loss == 0 {
			out[i] = 100
		} else {
			out[i] = 100 - 100/(1+gain/loss)
		}
	}
	return out
}

// KDJ computes stochastic K, D, and J from RSV.
func KDJ(candles []Candle, n, m1, m2 int) (k, d, j []float64) {
	indicatorCheck("KDJ", n)
	indicatorCheck("KDJ m1", m1)
	indicatorCheck("KDJ m2", m2)
	k, d, j = indicatorNaN(len(candles)), indicatorNaN(len(candles)), indicatorNaN(len(candles))
	pk, pd := 50.0, 50.0
	for i := n - 1; i < len(candles); i++ {
		lo, hi := math.Inf(1), math.Inf(-1)
		for _, c := range candles[i-n+1 : i+1] {
			lo, hi = min(lo, c.Low), max(hi, c.High)
		}
		rsv := 50.0
		if hi > lo {
			rsv = (candles[i].Close - lo) / (hi - lo) * 100
		}
		pk = (float64(m1-1)*pk + rsv) / float64(m1)
		pd = (float64(m2-1)*pd + pk) / float64(m2)
		k[i], d[i], j[i] = pk, pd, 3*pk-2*pd
	}
	return k, d, j
}

// OBV accumulates signed volume according to close direction.
func OBV(candles []Candle) []float64 {
	out := make([]float64, len(candles))
	for i := 1; i < len(candles); i++ {
		out[i] = out[i-1]
		switch {
		case candles[i].Close > candles[i-1].Close:
			out[i] += candles[i].Volume
		case candles[i].Close < candles[i-1].Close:
			out[i] -= candles[i].Volume
		}
	}
	return out
}

// HeikinAshi turns candles into Heikin-Ashi candles, keeping times and volumes.
func HeikinAshi(candles []Candle) []Candle {
	out := make([]Candle, len(candles))
	for i, c := range candles {
		h := Candle{Time: c.Time, Volume: c.Volume, Close: (c.Open + c.High + c.Low + c.Close) / 4}
		if i == 0 {
			h.Open = (c.Open + c.Close) / 2
		} else {
			h.Open = (out[i-1].Open + out[i-1].Close) / 2
		}
		h.High, h.Low = max(c.High, h.Open, h.Close), min(c.Low, h.Open, h.Close)
		out[i] = h
	}
	return out
}
