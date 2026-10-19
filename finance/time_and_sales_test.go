package finance

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func timeAndSalesTrades(n int) []TradePrint {
	at := time.Date(2026, 10, 5, 9, 30, 0, 0, time.UTC)
	out := make([]TradePrint, n)
	for i := range out {
		side := SideBuy
		if i%2 == 1 {
			side = SideSell
		}
		out[i] = TradePrint{ID: fmt.Sprint(i), Time: at.Add(time.Duration(i) * time.Second), Price: 100 + float64(i%50)/100, Size: float64(i%7 + 1), Side: side}
	}
	return out
}

func TestTimeAndSalesNewestFirstAndLargeMark(t *testing.T) {
	var s TimeAndSalesState
	trades := timeAndSalesTrades(3)
	trades[1].Size = 5000
	tt := tradingTester(core.GreenUp, 500, 400, func(c *ui.Context) {
		TimeAndSales(c, &s, trades, TimeAndSalesOptions{LargeSize: 1000})
	})
	newest, _ := tt.Find("09:30:02")
	oldest, _ := tt.Find("09:30:00")
	if newest.Y >= oldest.Y {
		t.Fatalf("newest %v not above oldest %v", newest, oldest)
	}
	if !tt.HasText("5,000") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if _, ok := tt.Find("09:30:01 Ask 100.01 5,000, large trade"); !ok {
		t.Fatal("large trade not spoken")
	}
}

func TestTimeAndSalesPauseHoldsView(t *testing.T) {
	var s TimeAndSalesState
	trades := timeAndSalesTrades(2)
	n := 0
	tt := tradingTester(core.GreenUp, 500, 400, func(c *ui.Context) {
		if TimeAndSales(c, &s, trades, TimeAndSalesOptions{}).Changed() {
			n++
		}
	})
	if err := tt.Click("Pause"); err != nil {
		t.Fatal(err)
	}
	trades = timeAndSalesTrades(5)
	tt.Frame()
	if tt.HasText("09:30:04") || !tt.HasText("Resume (3 new)") || !s.Paused || n != 1 {
		t.Fatalf("paused view moved: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("09:30:04") || s.Paused || n != 2 {
		t.Fatalf("resume: %q", tt.Texts())
	}
}

func TestTimeAndSalesBuildsOnlyVisibleRows(t *testing.T) {
	var s TimeAndSalesState
	trades := timeAndSalesTrades(100_000)
	tt := tradingTester(core.GreenUp, 500, 400, func(c *ui.Context) {
		TimeAndSales(c, &s, trades, TimeAndSalesOptions{})
	})
	if !tt.HasText("13:16:39") {
		t.Fatalf("newest trade not on top: %q", tt.Texts()[:8])
	}
	if got := len(tt.Texts()); got > 120 {
		t.Fatalf("built %d texts for 100,000 trades", got)
	}
}

func TestTimeAndSalesSnapshot(t *testing.T) {
	trades := timeAndSalesTrades(8)
	trades[5].Size = 4000
	snapshot(t, "time_and_sales", 420, 360, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			s := ui.Local(c.Root(), "snap", func() TimeAndSalesState { return TimeAndSalesState{} })
			TimeAndSales(c, s, trades, TimeAndSalesOptions{LargeSize: 1000, Height: 320})
		})
	})
}
