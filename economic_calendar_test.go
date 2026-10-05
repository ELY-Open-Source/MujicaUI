package mujica

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestEconomicCalendarCountryAndImportance(t *testing.T) {
	s := EconomicCalendarState{}
	rows := []EconomicEvent{{ID: "us", Country: "US", Title: "Payrolls", Time: mktNow, Importance: 3, Forecast: "150K", Previous: "142K", Actual: "254K"}, {ID: "eu", Country: "EU", Title: "CPI", Time: mktNow, Importance: 1}}
	tt := chatTester(960, 420, func(c *ui.Context) { EconomicCalendar(c, &s, rows, EconomicCalendarOptions{}).Element.Height(380) })
	if s.Table.Rows() != 2 || !tt.HasText("254K") || !tt.HasText("150K") || !tt.HasText("142K") {
		t.Fatalf("rows %d texts %q", s.Table.Rows(), tt.Texts())
	}
	_ = tt.Click("High")
	if s.MinImportance != 3 || s.Table.Rows() != 1 {
		t.Fatalf("importance %d rows %d", s.MinImportance, s.Table.Rows())
	}
	s.MinImportance = 0
	s.Countries = []string{"EU"}
	tt.Frame()
	if s.Table.Rows() != 1 || !tt.HasText("CPI") || tt.HasText("Payrolls") {
		t.Fatalf("country %q", tt.Texts())
	}
	marketsPanic(t, func(c *ui.Context) {
		EconomicCalendar(c, &s, []EconomicEvent{{ID: "x", Importance: 4}}, EconomicCalendarOptions{})
	})
}
func TestEconomicCalendarSnapshot(t *testing.T) {
	snapshot(t, "economic_calendar", 900, 360, func(c *ui.Context) {
		ec := ui.Local(c.Root(), "snapshot.economic_calendar", func() EconomicCalendarState { return EconomicCalendarState{} })
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			EconomicCalendar(c, ec, []EconomicEvent{{ID: "1", Country: "US", Title: "Nonfarm payrolls", Time: mktNow, Importance: 3, Forecast: "150K", Previous: "142K", Actual: "254K"},
				{ID: "2", Country: "EU", Title: "CPI", Time: mktNow, Importance: 2, Forecast: "1.8%", Previous: "2.2%"}}, EconomicCalendarOptions{}).Element.Height(300)

		})
	})
}
