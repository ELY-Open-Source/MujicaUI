package finance

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

func TestEarningsCalendarTimingActualAndKeys(t *testing.T) {
	var s data.DataTableState[string]
	rows := []EarningsEvent{{ID: "a", Symbol: "AAPL", Date: core.Date{Year: 2026, Month: 10, Day: 30}, Timing: EarningsAfterClose, EPSEstimate: 1.6, EPSActual: 1.64}, {ID: "b", Symbol: "JPM", Date: core.Date{Year: 2026, Month: 10, Day: 11}, Timing: EarningsBeforeOpen, EPSEstimate: 4, EPSActual: math.NaN()}}
	tt := chatTester(960, 420, func(c *ui.Context) { EarningsCalendar(c, &s, rows, EarningsCalendarOptions{}).Element.Height(380) })
	if !tt.HasText("Before open") || !tt.HasText("After close") || !tt.HasText("+0.04 (+2.5%)") || !tt.HasText("—") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("AAPL")
	tt.Key(0, ui.KeyDown)
	id, ok := s.Selected()
	if !ok || id != "b" {
		t.Fatalf("selected %s %v", id, ok)
	}
	if earningsSurpriseText(EarningsEvent{EPSEstimate: 0, EPSActual: 2}) != "+2.00" {
		t.Fatal("zero estimate")
	}
	marketsPanic(t, func(c *ui.Context) {
		EarningsCalendar(c, &s, []EarningsEvent{{ID: "x", Timing: 99}}, EarningsCalendarOptions{})
	})
}
func TestEarningsCalendarSnapshot(t *testing.T) {
	snapshot(t, "earnings_calendar", 900, 300, func(c *ui.Context) {
		es := ui.Local(c.Root(), "snapshot.earnings_calendar", func() data.DataTableState[string] { return data.DataTableState[string]{} })
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			EarningsCalendar(c, es, []EarningsEvent{{ID: "1", Symbol: "AAPL", Name: "Apple", Date: core.Date{Year: 2026, Month: 10, Day: 30}, Timing: EarningsAfterClose, EPSEstimate: 1.6, EPSActual: 1.64},
				{ID: "2", Symbol: "JPM", Name: "JPMorgan", Date: core.Date{Year: 2026, Month: 10, Day: 11}, Timing: EarningsBeforeOpen, EPSEstimate: 4.0, EPSActual: math.NaN()}}, EarningsCalendarOptions{}).Element.Height(240)

		})
	})
}
