package internal

import (
	"fmt"
	"math"
	"testing"
	"time"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestNiceAndTicks(t *testing.T) {
	for _, tc := range []struct {
		lo, hi        float64
		n             int
		nlo, nhi, stp float64
	}{
		{0, 97, 5, 0, 100, 20},
		{-3.2, 7.9, 4, -5, 10, 5},
		{0.012, 0.087, 5, 0, 0.1, 0.02},
		{5, 5, 4, 4, 6, 0.5},
	} {
		lo, hi, step := Nice(tc.lo, tc.hi, tc.n)
		if math.Abs(lo-tc.nlo) > 1e-9 || math.Abs(hi-tc.nhi) > 1e-9 || math.Abs(step-tc.stp) > 1e-9 {
			t.Errorf("Nice(%v, %v, %d) = %v %v %v", tc.lo, tc.hi, tc.n, lo, hi, step)
		}
	}
	if got := fmt.Sprint(Ticks(-5, 10, 5)); got != "[-5 0 5 10]" {
		t.Errorf("ticks %s", got)
	}
	if FormatTick(0.30000000000000004, 0.1) != "0.3" || FormatTick(1500, 500) != "1500" || FormatTick(-1e-17, 0.5) != "0.0" {
		t.Error("FormatTick")
	}
}

func TestNicePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an inverted domain")
		}
	}()
	Nice(2, 1, 4)
}

func TestTimeTicks(t *testing.T) {
	lo := time.Date(2026, 10, 5, 9, 7, 0, 0, time.UTC)
	ticks, layout := TimeTicks(lo, lo.Add(2*time.Hour), 4)
	if layout != "15:04" || len(ticks) != 4 || ticks[0].Format("15:04") != "09:30" || ticks[3].Format("15:04") != "11:00" {
		t.Fatalf("ticks %v layout %s", ticks, layout)
	}
}

func TestLinearAndBand(t *testing.T) {
	s := Linear{D0: 0, D1: 10, R0: 100, R1: 0}
	if s.Map(5) != 50 || s.Invert(25) != 7.5 {
		t.Fatal("Linear")
	}
	b := Band{N: 4, R0: 0, R1: 200, Pad: 0.5}
	if b.Step() != 50 || b.Width() != 25 || b.Start(1) != 62.5 || b.Center(3) != 175 {
		t.Fatal("Band")
	}
}

func TestPaletteContrast(t *testing.T) {
	for _, k := range []theme.Tokens{theme.Light(), theme.Dark()} {
		pal := Palette(k)
		seen := map[[3]uint8]bool{}
		for i, c := range pal {
			if r := theme.Contrast(c, k.Surface); r < 3 {
				t.Errorf("dark=%v color %d contrast %.2f", k.Dark, i, r)
			}
			seen[[3]uint8{c.R, c.G, c.B}] = true
		}
		if len(pal) != 6 || len(seen) != 6 {
			t.Errorf("dark=%v palette not six distinct colors", k.Dark)
		}
	}
}

func TestNearest(t *testing.T) {
	pts := [][2]float32{{0, 0}, {10, 10}, {20, 0}}
	at := func(i int) (float32, float32) { return pts[i][0], pts[i][1] }
	if Nearest(3, at, 9, 2, 30, false) != 1 || Nearest(3, at, 18, 50, 30, true) != 2 || Nearest(3, at, 100, 100, 5, false) != -1 {
		t.Fatal("Nearest")
	}
}

func TestColumnsBounded(t *testing.T) {
	const n = 100_000
	x := func(i int) float32 { return float32(i) * 400 / n }
	y := func(i int) float64 { return math.Sin(float64(i)) }
	got := Columns(n, x, y)
	if len(got) > 4*400+4 || got[0] != 0 || got[len(got)-1] != n-1 {
		t.Fatalf("%d points kept, first %d last %d", len(got), got[0], got[len(got)-1])
	}
	for i := 1; i < len(got); i++ {
		if got[i] <= got[i-1] {
			t.Fatal("indices out of order")
		}
	}
}
