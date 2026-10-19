package finance

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
)

func positionSample() []PositionRow {
	return []PositionRow{
		{ID: "p1", Symbol: "MUJC", Quantity: 100, AvgCost: 400, Price: 412.08},
		{ID: "p2", Symbol: "AVE", Quantity: 50, AvgCost: 90, Price: 88.2},
	}
}

func TestPositionTableValuesAndClose(t *testing.T) {
	var s data.DataTableState[string]
	rows := positionSample()
	closed := ""
	tt := tradingTester(core.GreenUp, 1000, 400, func(c *ui.Context) {
		ui.Box(c).Width(960).Height(300).Children(func() {
			if id, ok := PositionTable(c, &s, rows, PositionTableOptions{}).Acted(); ok {
				closed = id
			}
		})
	})
	for _, want := range []string{"41,208.00", "+1,208.00", "+3.02%", "−90.00", "−2.00%"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	if err := tt.Click("Close AVE"); err != nil {
		t.Fatal(err)
	}
	if closed != "p2" {
		t.Fatalf("closed %q", closed)
	}
}

// Selection follows the row ID when rows reorder and change.
func TestPositionTableStableIDs(t *testing.T) {
	var s data.DataTableState[string]
	rows := positionSample()
	tt := tradingTester(core.GreenUp, 1000, 400, func(c *ui.Context) {
		ui.Box(c).Width(960).Height(300).Children(func() { PositionTable(c, &s, rows, PositionTableOptions{}) })
	})
	_ = tt.Click("AVE")
	rows = []PositionRow{{ID: "p3", Symbol: "OBLV", Quantity: 1, AvgCost: 1, Price: 1}, {ID: "p2", Symbol: "AVE", Quantity: 60, AvgCost: 90, Price: 91}, rows[0]}
	tt.Frame()
	if id, ok := s.Selected(); !ok || id != "p2" {
		t.Fatalf("selection %q %v", id, ok)
	}
}

func TestPositionTableDisabledAndLargeData(t *testing.T) {
	var s data.DataTableState[string]
	rows := make([]PositionRow, 100_000)
	for i := range rows {
		rows[i] = PositionRow{ID: fmt.Sprint(i), Symbol: fmt.Sprintf("S%05d", i), Quantity: 1, AvgCost: 1, Price: 2}
	}
	closed := false
	tt := tradingTester(core.GreenUp, 1000, 400, func(c *ui.Context) {
		ui.Box(c).Width(960).Height(300).Children(func() {
			_, closed = PositionTable(c, &s, rows, PositionTableOptions{Disabled: true}).Acted()
		})
	})
	_ = tt.Click("Close S00000")
	if closed {
		t.Fatal("disabled table closed a position")
	}
	if n := len(tt.Texts()); n > 400 {
		t.Fatalf("built %d texts for 100,000 rows", n)
	}
}

func TestPositionTableSnapshot(t *testing.T) {
	snapshot(t, "position_table", 1040, 220, func(c *ui.Context) {
		ui.Column(c).Padding(16).Height(200).Width(1010).Children(func() {
			s := ui.Local(c.Root(), "snap", func() data.DataTableState[string] { return data.DataTableState[string]{} })
			PositionTable(c, s, positionSample(), PositionTableOptions{})
		})
	})
}
