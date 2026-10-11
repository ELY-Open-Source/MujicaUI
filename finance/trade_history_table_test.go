package finance

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/datetime"
)

func tradeHistorySample() []FillRow {
	day := func(d, h int) time.Time { return time.Date(2026, 10, d, h, 0, 0, 0, time.UTC) }
	return []FillRow{
		{ID: "f1", Time: day(1, 10), Symbol: "MUJC", Side: SideBuy, Price: 400, Quantity: 10, Fee: 1.2},
		{ID: "f2", Time: day(3, 14), Symbol: "AVE", Side: SideSell, Price: 88.2, Quantity: 5, Fee: 0.4},
		{ID: "f3", Time: day(5, 9), Symbol: "OBLV", Side: SideBuy, Price: 12, Quantity: 100, Fee: 0.6},
	}
}

func TestTradeHistoryTableRangeAndExport(t *testing.T) {
	s := TradeHistoryTableState{Range: datetime.DateRange{Start: core.Date{Year: 2026, Month: 10, Day: 2}, End: core.Date{Year: 2026, Month: 10, Day: 4}}}
	var exported []FillRow
	tt := tradingTester(core.GreenUp, 1000, 600, func(c *ui.Context) {
		ui.Box(c).Width(960).Height(400).Children(func() {
			if rows, ok := TradeHistoryTable(c, &s, tradeHistorySample(), TradeHistoryTableOptions{}).Exported(); ok {
				exported = rows
			}
		})
	})
	if tt.HasText("MUJC") || !tt.HasText("AVE") || tt.HasText("OBLV") {
		t.Fatalf("range: %q", tt.Texts())
	}
	if err := tt.Click("Export"); err != nil {
		t.Fatal(err)
	}
	if len(exported) != 1 || exported[0].ID != "f2" {
		t.Fatalf("exported %v", exported)
	}
	if err := tt.Click("All dates"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("MUJC") || !tt.HasText("OBLV") || !s.Range.IsZero() {
		t.Fatalf("all dates: %q", tt.Texts())
	}
	_ = tt.Click("Export")
	if len(exported) != 3 {
		t.Fatalf("exported %d", len(exported))
	}
}

func TestTradeHistoryTableDisabled(t *testing.T) {
	var s TradeHistoryTableState
	exported := false
	tt := tradingTester(core.GreenUp, 1000, 600, func(c *ui.Context) {
		ui.Box(c).Width(960).Height(400).Children(func() {
			_, exported = TradeHistoryTable(c, &s, tradeHistorySample(), TradeHistoryTableOptions{Disabled: true}).Exported()
		})
	})
	_ = tt.Click("Export")
	if exported {
		t.Fatal("disabled table exported")
	}
}

func TestTradeHistoryTableSnapshot(t *testing.T) {
	snapshot(t, "trade_history_table", 900, 260, func(c *ui.Context) {
		ui.Column(c).Padding(16).Height(240).Width(880).Children(func() {
			s := ui.Local(c.Root(), "snap", func() TradeHistoryTableState { return TradeHistoryTableState{} })
			TradeHistoryTable(c, s, tradeHistorySample(), TradeHistoryTableOptions{})
		})
	})
}

// 100,000 fills build only visible rows.
func TestTradeHistoryTableLargeData(t *testing.T) {
	var s TradeHistoryTableState
	rows := make([]FillRow, 100_000)
	for i := range rows {
		rows[i] = FillRow{ID: fmt.Sprint(i), Time: time.Date(2026, 1, 1, 0, 0, i, 0, time.UTC), Symbol: "S", Side: SideSell, Price: 1, Quantity: 1}
	}
	tt := tradingTester(core.GreenUp, 1000, 600, func(c *ui.Context) {
		ui.Box(c).Width(960).Height(400).Children(func() { TradeHistoryTable(c, &s, rows, TradeHistoryTableOptions{}) })
	})
	if n := len(tt.Texts()); n > 400 {
		t.Fatalf("built %d texts", n)
	}
}
