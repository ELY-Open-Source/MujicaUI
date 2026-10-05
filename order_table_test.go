package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func orderTableSample() []OrderRow {
	return []OrderRow{
		{ID: "o1", Symbol: "MUJC", Side: SideBuy, Type: OrderLimit, Price: 99.5, Quantity: 100, Filled: 40, Status: OrderPartial},
		{ID: "o2", Symbol: "AVE", Side: SideSell, Type: OrderMarket, Quantity: 10, Filled: 10, Status: OrderFilled},
		{ID: "o3", Symbol: "OBLV", Side: SideBuy, Type: OrderStop, Price: 12, Quantity: 5, Status: OrderCanceled},
		{ID: "o4", Symbol: "DOLO", Side: SideSell, Type: OrderLimit, Price: 7.25, Quantity: 30, Status: OrderOpen},
	}
}

func orderTableTester(s *OrderTableState, rows *[]OrderRow, opts *OrderTableOptions, acted *string) *ui.Tester {
	return tradingTester(GreenUp, 1000, 400, func(c *ui.Context) {
		ui.Box(c).Width(960).Height(320).Children(func() {
			if id, ok := OrderTable(c, s, *rows, *opts).Acted(); ok {
				*acted = id
			}
		})
	})
}

func TestOrderTableFilterProgressAndCancel(t *testing.T) {
	var s OrderTableState
	rows, opts, acted := orderTableSample(), OrderTableOptions{}, ""
	tt := orderTableTester(&s, &rows, &opts, &acted)
	if !tt.HasText("40 / 100") || !tt.HasText("Partly filled") || !tt.HasText("Market") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if err := tt.Click("Working"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("AVE") || tt.HasText("OBLV") || !tt.HasText("DOLO") {
		t.Fatalf("working filter: %q", tt.Texts())
	}
	if _, ok := tt.Find("Cancel AVE"); ok {
		t.Fatal("filled order offers cancel")
	}
	if err := tt.Click("Cancel DOLO"); err != nil {
		t.Fatal(err)
	}
	if acted != "o4" {
		t.Fatalf("acted %q", acted)
	}
	_ = tt.Click("Closed")
	if !tt.HasText("OBLV") || tt.HasText("MUJC") {
		t.Fatalf("closed filter: %q", tt.Texts())
	}
}

func TestOrderTableFiltersOnlyOnChange(t *testing.T) {
	var s OrderTableState
	rows, opts, acted := orderTableSample(), OrderTableOptions{}, ""
	tt := orderTableTester(&s, &rows, &opts, &acted)
	view := s.cache.visible
	tt.Frame()
	if &s.cache.visible[0] != &view[0] {
		t.Fatal("refiltered without a change")
	}
	rows[0].Status = OrderFilled
	opts.Version++
	_ = tt.Click("Working")
	if len(s.cache.visible) != 1 {
		t.Fatalf("version change ignored: %d rows", len(s.cache.visible))
	}
}

func TestOrderTableDisabledAndBadFill(t *testing.T) {
	var s OrderTableState
	rows, opts, acted := orderTableSample(), OrderTableOptions{Disabled: true}, ""
	tt := orderTableTester(&s, &rows, &opts, &acted)
	_ = tt.Click("Cancel DOLO")
	_ = tt.Click("Working")
	if acted != "" || s.Filter != OrdersAll {
		t.Fatal("disabled table acted")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("overfilled order accepted")
		}
	}()
	var s2 OrderTableState
	bad := []OrderRow{{ID: "x", Quantity: 1, Filled: 2}}
	orderTableTester(&s2, &bad, &OrderTableOptions{}, &acted)
}

func TestOrderTableSnapshot(t *testing.T) {
	snapshot(t, "order_table", 960, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Height(280).Width(940).Children(func() {
			s := ui.Local(c.Root(), "snap", func() OrderTableState { return OrderTableState{} })
			OrderTable(c, s, orderTableSample(), OrderTableOptions{})
		})
	})
}

// 100,000 orders filter once and build only visible rows.
func TestOrderTableLargeData(t *testing.T) {
	var s OrderTableState
	rows := make([]OrderRow, 100_000)
	for i := range rows {
		rows[i] = OrderRow{ID: fmt.Sprint(i), Symbol: fmt.Sprintf("S%05d", i), Side: SideBuy, Price: 1, Quantity: 10, Status: OrderStatus(i % 5)}
	}
	opts, acted := OrderTableOptions{}, ""
	tt := orderTableTester(&s, &rows, &opts, &acted)
	_ = tt.Click("Working")
	if len(s.cache.visible) != 40_000 {
		t.Fatalf("working rows %d", len(s.cache.visible))
	}
	if n := len(tt.Texts()); n > 400 {
		t.Fatalf("built %d texts", n)
	}
}
