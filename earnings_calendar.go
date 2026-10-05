package mujica

import (
	"cmp"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// EarningsTiming is when a report comes. The zero value is EarningsUnknown.
type EarningsTiming int

const (
	EarningsUnknown EarningsTiming = iota
	EarningsBeforeOpen
	EarningsAfterClose
)

// EarningsEvent is a scheduled report; NaN actuals are not yet reported.
type EarningsEvent struct {
	ID, Symbol, Name string
	Date             Date
	Timing           EarningsTiming
	EPSEstimate      float64
	EPSActual        float64
}

// EarningsCalendarOptions configure an EarningsCalendar.
type EarningsCalendarOptions struct {
	Version  int
	Status   DataStatus
	Error    string
	Disabled bool
}

// EarningsCalendar lists reports with before-open or after-close timing, estimate, actual and surprise.
func EarningsCalendar(c *ui.Context, s *DataTableState[string], events []EarningsEvent, opts EarningsCalendarOptions) DataTableView {
	if s == nil {
		panic("mujica: EarningsCalendar needs a state")
	}
	cols := []DataColumn[EarningsEvent]{
		{ID: "date", Title: msg(c, "earnings_calendar.date"), Width: 110, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return a.Date.Compare(b.Date) },
			Text:    func(e EarningsEvent) string { return e.Date.String() }},
		{ID: "timing", Title: msg(c, "earnings_calendar.timing"), Width: 130, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return cmp.Compare(a.Timing, b.Timing) },
			Text:    func(e EarningsEvent) string { return msg(c, earningsTimings[earningsTimingCheck(e.Timing)].key) },
			Cell:    func(e EarningsEvent) { earningsTiming(c, e.Timing) }},
		{ID: "symbol", Title: msg(c, "earnings_calendar.symbol"), Width: 90, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return cmp.Compare(a.Symbol, b.Symbol) },
			Text:    func(e EarningsEvent) string { return e.Symbol }},
		{ID: "name", Title: msg(c, "earnings_calendar.name"), Text: func(e EarningsEvent) string { return e.Name }},
		{ID: "estimate", Title: msg(c, "earnings_calendar.estimate"), Width: 90, Align: ui.End, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return cmp.Compare(a.EPSEstimate, b.EPSEstimate) },
			Text:    func(e EarningsEvent) string { return earningsValue(e.EPSEstimate) }},
		{ID: "actual", Title: msg(c, "earnings_calendar.actual"), Width: 90, Align: ui.End, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return cmp.Compare(a.EPSActual, b.EPSActual) },
			Text:    func(e EarningsEvent) string { return earningsValue(e.EPSActual) }},
		{ID: "surprise", Title: msg(c, "earnings_calendar.surprise"), Width: 110, Align: ui.End,
			Text: func(e EarningsEvent) string { return earningsSurpriseText(e) },
			Cell: func(e EarningsEvent) {
				if d := e.EPSActual - e.EPSEstimate; !math.IsNaN(d) {
					tradingSigned(c, d, earningsSurpriseText(e))
				} else {
					ui.Text(c, "—").TextColor(Tokens(c).TextMuted).FillWidth().TextAlign(ui.End)
				}
			}},
	}
	var v DataTableView
	col := ui.Column(c).MinHeight(0).Grow(1).Disabled(opts.Disabled)
	col.Children(func() {
		v = DataTable(c, s, events, DataTableOptions[EarningsEvent, string]{Columns: cols, Key: func(e EarningsEvent) string { return e.ID },
			Version: opts.Version, Status: opts.Status, Error: opts.Error, Empty: msg(c, "earnings_calendar.empty")})
		v.Element.Grow(1)
	})
	v.Element = col
	return v
}

var earningsTimings = [3]struct{ key, icon string }{
	{"earnings_calendar.unknown", "clock"}, {"earnings_calendar.before", "sun"}, {"earnings_calendar.after", "moon"},
}

func earningsTimingCheck(t EarningsTiming) EarningsTiming {
	if t < EarningsUnknown || t > EarningsAfterClose {
		panic("mujica: EarningsTiming out of range")
	}
	return t
}

// earningsTiming shows the timing as an icon and words.
func earningsTiming(c *ui.Context, t EarningsTiming) {
	k := Tokens(c)
	tm := earningsTimings[earningsTimingCheck(t)]
	ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
		ui.Icon(c, icons.Must(tm.icon)).Size(14, 14).TextColor(k.TextMuted)
		ui.Text(c, msg(c, tm.key)).FontSize(FontSize(c, theme.CaptionSize)).SingleLine()
	})
}

func earningsValue(v float64) string {
	if math.IsNaN(v) {
		return "—"
	}
	return statisticFormat(v, 2)
}

// earningsSurpriseText is actual minus estimate, with the percent of the estimate.
func earningsSurpriseText(e EarningsEvent) string {
	d := e.EPSActual - e.EPSEstimate
	if math.IsNaN(d) {
		return "—"
	}
	if e.EPSEstimate == 0 {
		return statisticSigned(d, 2)
	}
	return statisticSigned(d, 2) + " (" + statisticSigned(d/math.Abs(e.EPSEstimate)*100, 1) + "%)"
}
