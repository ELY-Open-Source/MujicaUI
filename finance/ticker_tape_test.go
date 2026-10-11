package finance

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func tickerTapeItems(n int) []TickerItem {
	var items []TickerItem
	for i := range n {
		items = append(items, TickerItem{Symbol: fmt.Sprintf("SYM%d", i), Price: 100 + float64(i), ChangePercent: float64(i%3 - 1)})
	}
	return items
}

func tickerTapeX(t *testing.T, tt *ui.Tester) float32 {
	t.Helper()
	r, ok := tt.Find("SYM0")
	if !ok {
		t.Fatalf("no SYM0 in %q", tt.Texts())
	}
	return r.X
}

func TestTickerTapeScrollsAndPausesOnHover(t *testing.T) {
	tt := tradingTester(core.GreenUp, 400, 80, func(c *ui.Context) {
		ui.Box(c).Width(380).Children(func() { TickerTape(c, tickerTapeItems(6), TickerTapeOptions{Speed: 200}) })
	})
	x0 := tickerTapeX(t, tt)
	time.Sleep(60 * time.Millisecond)
	tt.Frame()
	x1 := tickerTapeX(t, tt)
	if x1 >= x0 {
		t.Fatalf("did not scroll: %v -> %v", x0, x1)
	}
	tt.Move(100, 25)
	tt.Frame()
	x2 := tickerTapeX(t, tt)
	time.Sleep(60 * time.Millisecond)
	tt.Frame()
	if x3 := tickerTapeX(t, tt); x3 != x2 {
		t.Fatalf("moved under the pointer: %v -> %v", x2, x3)
	}
}

func TestTickerTapePausesOnFocus(t *testing.T) {
	tt := tradingTester(core.GreenUp, 400, 80, func(c *ui.Context) {
		ui.Box(c).Width(380).Children(func() { TickerTape(c, tickerTapeItems(6), TickerTapeOptions{Speed: 200}) })
	})
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Market ticker") {
		t.Fatal("tape not focusable")
	}
	x0 := tickerTapeX(t, tt)
	time.Sleep(60 * time.Millisecond)
	tt.Frame()
	if x1 := tickerTapeX(t, tt); x1 != x0 {
		t.Fatalf("moved while focused: %v -> %v", x0, x1)
	}
}

func TestTickerTapePagesUnderReducedMotion(t *testing.T) {
	tt := tradingTester(core.GreenUp, 500, 80, func(c *ui.Context) {
		TickerTape(c, tickerTapeItems(5), TickerTapeOptions{PerPage: 2})
	})
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if !tt.HasText("SYM0") || tt.HasText("SYM2") || !tt.HasText("1 / 3") || !tt.HasText("−1.00%") {
		t.Fatalf("page 1: %q", tt.Texts())
	}
	if err := tt.Click("Next page"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("SYM2") || tt.HasText("SYM0") || !tt.HasText("+1.00%") {
		t.Fatalf("page 2: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	if !tt.HasText("SYM4") || !tt.HasText("3 / 3") {
		t.Fatalf("wrap back: %q", tt.Texts())
	}
}

func TestTickerTapeSnapshot(t *testing.T) {
	snapshot(t, "ticker_tape", 520, 60, func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() {
			TickerTape(c, []TickerItem{{"MUJC", 412.08, -1.35}, {"AVE", 88.2, 2.1}, {"OBLV", 1204.5, 0}}, TickerTapeOptions{})
		})
	})
}
