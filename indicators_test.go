package mujica

import (
	"math"
	"testing"
	"time"
)

func approx(t *testing.T, name string, got, want []float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: len %d, want %d", name, len(got), len(want))
	}
	for i := range got {
		if math.IsNaN(want[i]) != math.IsNaN(got[i]) || !math.IsNaN(want[i]) && math.Abs(got[i]-want[i]) > 1e-9 {
			t.Errorf("%s[%d] = %v, want %v", name, i, got[i], want[i])
		}
	}
}

var nan = math.NaN()

func TestSMAAndEMA(t *testing.T) {
	v := []float64{1, 2, 3, 4, 5, 6}
	approx(t, "SMA", SMA(v, 3), []float64{nan, nan, 2, 3, 4, 5})
	approx(t, "EMA", EMA(v, 3), []float64{nan, nan, 2, 3, 4, 5})
	approx(t, "EMA2", EMA([]float64{10, 20, 10}, 1), []float64{10, 20, 10})
	approx(t, "EMA3", EMA([]float64{2, 4, 8}, 2), []float64{nan, 3, 3 + (8-3)*2.0/3})
}

func TestBollinger(t *testing.T) {
	mid, up, lo := Bollinger([]float64{2, 4, 4, 4, 5, 5, 7, 9}, 8, 2)
	// Mean 5, population deviation 2.
	approx(t, "mid", mid[7:], []float64{5})
	approx(t, "upper", up[7:], []float64{9})
	approx(t, "lower", lo[7:], []float64{1})
	if !math.IsNaN(up[6]) {
		t.Fatal("band before the period")
	}
}

func candlesOf(rows ...[5]float64) []Candle {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	out := make([]Candle, len(rows))
	for i, r := range rows {
		out[i] = Candle{Time: base.Add(time.Duration(i) * time.Hour), Open: r[0], High: r[1], Low: r[2], Close: r[3], Volume: r[4]}
	}
	return out
}

func TestVWAPAndOBV(t *testing.T) {
	cs := candlesOf([5]float64{10, 12, 9, 12, 100}, [5]float64{12, 15, 12, 12, 300}, [5]float64{12, 12, 6, 9, 100})
	// Typical prices 11, 13, 9.
	approx(t, "VWAP", VWAP(cs), []float64{11, (1100 + 3900) / 400.0, (1100 + 3900 + 900) / 500.0})
	approx(t, "OBV", OBV(cs), []float64{0, 0, -100})
}

func TestMACD(t *testing.T) {
	v := []float64{1, 2, 3, 4, 5, 6, 7}
	dif, dea, hist := MACD(v, 2, 3, 2)
	f, s := EMA(v, 2), EMA(v, 3)
	for i := 2; i < len(v); i++ {
		if math.Abs(dif[i]-(f[i]-s[i])) > 1e-12 {
			t.Fatalf("dif[%d]", i)
		}
	}
	want := EMA(dif[2:], 2)
	approx(t, "dea", dea[2:], want)
	if !math.IsNaN(dea[1]) || math.Abs(hist[6]-(dif[6]-dea[6])) > 1e-12 {
		t.Fatal("hist or warm-up wrong")
	}
	approx(t, "dif line", dif[2:], []float64{0.5, 0.5, 0.5, 0.5, 0.5})
	approx(t, "hist line", hist[3:], []float64{0, 0, 0, 0})
}

func TestRSI(t *testing.T) {
	r := RSI([]float64{1, 2, 3, 2}, 2)
	approx(t, "RSI", r, []float64{nan, nan, 100, 50})
}

func TestKDJ(t *testing.T) {
	cs := candlesOf([5]float64{1, 2, 0, 1, 0}, [5]float64{1, 4, 0, 4, 0}, [5]float64{4, 4, 2, 3, 0})
	k, d, j := KDJ(cs, 2, 3, 3)
	k1 := 200.0 / 3
	d1 := (100 + k1) / 3
	// i=2: range 0..4, RSV 75.
	k2 := (2*k1 + 75) / 3
	d2 := (2*d1 + k2) / 3
	approx(t, "K", k, []float64{nan, k1, k2})
	approx(t, "D", d, []float64{nan, d1, d2})
	approx(t, "J", j, []float64{nan, 3*k1 - 2*d1, 3*k2 - 2*d2})
}

func TestHeikinAshi(t *testing.T) {
	cs := candlesOf([5]float64{10, 14, 8, 12, 5}, [5]float64{12, 13, 9, 10, 6})
	h := HeikinAshi(cs)
	if h[0].Open != 11 || h[0].Close != 11 || h[0].High != 14 || h[0].Low != 8 {
		t.Fatalf("first %+v", h[0])
	}
	if h[1].Open != 11 || h[1].Close != 11 || h[1].High != 13 || h[1].Low != 9 || h[1].Volume != 6 {
		t.Fatalf("second %+v", h[1])
	}
}

func TestIndicatorsPanicOnBadPeriods(t *testing.T) {
	for name, f := range map[string]func(){
		"SMA":   func() { SMA(nil, 0) },
		"MACD":  func() { MACD(nil, 26, 12, 9) },
		"BOLL":  func() { Bollinger(nil, 20, 0) },
		"KDJ m": func() { KDJ(nil, 9, 0, 3) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			f()
		}()
	}
}

func TestMarketWindowPanZoom(t *testing.T) {
	v := MarketView{}
	if f, to := marketWindow(v, 1000, 400); f != 880 || to != 1000 {
		t.Fatalf("default %d..%d", f, to)
	}
	marketPan(&v, 1000, 400, -100)
	if v.End != 900 || v.Count != 120 {
		t.Fatalf("pan back %+v", v)
	}
	marketZoom(&v, 1000, 400, 0.5)
	if f, to := marketWindow(v, 1000, 400); f != 840 || to != 900 {
		t.Fatalf("zoom in %d..%d", f, to)
	}
	marketZoom(&v, 1000, 100, 10)
	if v.Count != 100 {
		t.Fatalf("zoom out capped by room: %+v", v)
	}
	marketPan(&v, 1000, 100, 5000)
	if v.End != 0 {
		t.Fatalf("pan to the newest follows: %+v", v)
	}
	if f, to := marketWindow(MarketView{Count: 50, End: 10}, 1000, 400); f != 0 || to != 50 {
		t.Fatalf("start clamp %d..%d", f, to)
	}
}

func TestIntervalClose(t *testing.T) {
	loc := time.UTC
	at := time.Date(2026, 10, 7, 13, 47, 10, 0, loc) // a Wednesday
	for iv, want := range map[Interval]time.Time{
		"1m":  time.Date(2026, 10, 7, 13, 48, 0, 0, loc),
		"15m": time.Date(2026, 10, 7, 14, 0, 0, 0, loc),
		"4h":  time.Date(2026, 10, 7, 16, 0, 0, 0, loc),
		"1D":  time.Date(2026, 10, 8, 0, 0, 0, 0, loc),
		"1W":  time.Date(2026, 10, 12, 0, 0, 0, 0, loc),
		"1M":  time.Date(2026, 11, 1, 0, 0, 0, 0, loc),
	} {
		if got := iv.Close(at, loc); !got.Equal(want) {
			t.Errorf("%s: %v, want %v", iv, got, want)
		}
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown interval accepted")
		}
	}()
	Interval("7m").Close(at, loc)
}
