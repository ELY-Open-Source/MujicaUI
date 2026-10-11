package finance

import (
	"cmp"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

type EarningsTiming int

const (
	EarningsUnknown EarningsTiming = iota
	EarningsBeforeOpen
	EarningsAfterClose
)

type EarningsEvent struct {
	ID, Symbol, Name string
	Date             core.Date
	Timing           EarningsTiming
	EPSEstimate      float64
	EPSActual        float64
}

// EarningsCalendarOptions configure an EarningsCalendar.
type EarningsCalendarOptions struct {
	Version  int
	Status   data.DataStatus
	Error    string
	Disabled bool
}

func EarningsCalendar(c *ui.Context, s *data.DataTableState[string], events []EarningsEvent, opts EarningsCalendarOptions) data.DataTableView {
	if s == nil {
		panic("mujica: EarningsCalendar needs a state")
	}
	cols := []data.DataColumn[EarningsEvent]{
		{ID: "date", Title: core.Msg(c, "earnings_calendar.date"), Width: 110, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return a.Date.Compare(b.Date) },
			Text:    func(e EarningsEvent) string { return e.Date.String() }},
		{ID: "timing", Title: core.Msg(c, "earnings_calendar.timing"), Width: 130, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return cmp.Compare(a.Timing, b.Timing) },
			Text:    func(e EarningsEvent) string { return core.Msg(c, earningsTimings[earningsTimingCheck(e.Timing)].key) },
			Cell:    func(e EarningsEvent) { earningsTiming(c, e.Timing) }},
		{ID: "symbol", Title: core.Msg(c, "earnings_calendar.symbol"), Width: 90, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return cmp.Compare(a.Symbol, b.Symbol) },
			Text:    func(e EarningsEvent) string { return e.Symbol }},
		{ID: "name", Title: core.Msg(c, "earnings_calendar.name"), Text: func(e EarningsEvent) string { return e.Name }},
		{ID: "estimate", Title: core.Msg(c, "earnings_calendar.estimate"), Width: 90, Align: ui.End, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return cmp.Compare(a.EPSEstimate, b.EPSEstimate) },
			Text:    func(e EarningsEvent) string { return earningsValue(e.EPSEstimate) }},
		{ID: "actual", Title: core.Msg(c, "earnings_calendar.actual"), Width: 90, Align: ui.End, Sortable: true,
			Compare: func(a, b EarningsEvent) int { return cmp.Compare(a.EPSActual, b.EPSActual) },
			Text:    func(e EarningsEvent) string { return earningsValue(e.EPSActual) }},
		{ID: "surprise", Title: core.Msg(c, "earnings_calendar.surprise"), Width: 150, Align: ui.End,
			Text: func(e EarningsEvent) string { return earningsSurpriseText(e) },
			Cell: func(e EarningsEvent) {
				if d := e.EPSActual - e.EPSEstimate; !math.IsNaN(d) {
					tradingSigned(c, d, earningsSurpriseText(e))
				} else {
					ui.Text(c, "—").TextColor(core.Tokens(c).TextMuted).FillWidth().TextAlign(ui.End)
				}
			}},
	}
	var v data.DataTableView
	col := ui.Column(c).MinHeight(0).Grow(1).Disabled(opts.Disabled)
	opts.Disabled = col.IsDisabled()
	if opts.Disabled {
		copy := *s
		s = &copy
	}
	col.Children(func() {
		v = data.DataTable(c, s, events, data.DataTableOptions[EarningsEvent, string]{Columns: cols, Key: func(e EarningsEvent) string { return e.ID },
			Version: opts.Version, Status: opts.Status, Error: opts.Error, Empty: core.Msg(c, "earnings_calendar.empty")})
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
	k := core.Tokens(c)
	tm := earningsTimings[earningsTimingCheck(t)]
	ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
		ui.Icon(c, icons.Must(tm.icon)).Size(14, 14).TextColor(k.TextMuted)
		ui.Text(c, core.Msg(c, tm.key)).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
	})
}

func earningsValue(v float64) string {
	if math.IsNaN(v) {
		return "—"
	}
	return data.StatisticFormat(v, 2)
}

func earningsSurpriseText(e EarningsEvent) string {
	d := e.EPSActual - e.EPSEstimate
	if math.IsNaN(d) {
		return "—"
	}
	if e.EPSEstimate == 0 {
		return data.StatisticSigned(d, 2)
	}
	return data.StatisticSigned(d, 2) + " (" + data.StatisticSigned(d/math.Abs(e.EPSEstimate)*100, 1) + "%)"
}
