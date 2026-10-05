package finance

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// EconomicEvent is a scheduled release; Actual stays empty until published.
type EconomicEvent struct {
	ID, Country, Title string
	Time               time.Time
	// Importance is 1 (low) to 3 (high).
	Importance                 int
	Forecast, Previous, Actual string
}

type EconomicCalendarState struct {
	Countries []string
	// MinImportance hides events below it; zero shows all.
	MinImportance int
	Table         data.DataTableState[string]
	cache         economicCache
}

type economicCache struct {
	key       string
	rows      []EconomicEvent
	data      string
	countries []input.SelectOption[string]
}

// EconomicCalendarOptions configure an EconomicCalendar.
type EconomicCalendarOptions struct {
	Version  int
	Status   data.DataStatus
	Error    string
	Disabled bool
}

func EconomicCalendar(c *ui.Context, s *EconomicCalendarState, events []EconomicEvent, opts EconomicCalendarOptions) data.DataTableView {
	if s == nil {
		panic("mujica: EconomicCalendar needs a state")
	}
	if s.MinImportance < 0 || s.MinImportance > 3 {
		panic("mujica: EconomicCalendar MinImportance must be 0..3")
	}
	u := core.Density(c).Unit()
	var v data.DataTableView
	col := ui.Column(c).Gap(u * 2).MinHeight(0).Grow(1).Disabled(opts.Disabled)
	opts.Disabled = col.IsDisabled()
	if opts.Disabled {
		copy := *s
		s = &copy
	}
	col.Children(func() {
		if data := chart.ChartID(events) + fmt.Sprint(opts.Version); data != s.cache.data {
			s.cache.data, s.cache.countries = data, nil
			seen := map[string]bool{}
			for _, e := range events {
				if e.ID == "" || seen[e.ID] || e.Importance < 1 || e.Importance > 3 {

					panic("mujica: EconomicEvent needs an ID and Importance 1..3")
				}
				seen[e.ID] = true
				if !slices.ContainsFunc(s.cache.countries, func(o input.SelectOption[string]) bool { return o.Value == e.Country }) {
					s.cache.countries = append(s.cache.countries, input.SelectOption[string]{Value: e.Country, Label: e.Country})
				}
			}
			slices.SortFunc(s.cache.countries, func(a, b input.SelectOption[string]) int { return cmp.Compare(a.Value, b.Value) })
		}
		countries := s.cache.countries
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Wrap().Children(func() {
			input.MultiSelect(c, &s.Countries, countries, input.MultiSelectOptions{Placeholder: core.Msg(c, "economic_calendar.all_countries"),
				Label: core.Msg(c, "economic_calendar.countries"), Disabled: opts.Disabled}).Element.MinWidth(220)
			imp := []input.RadioItem[int]{{Value: 0, Label: core.Msg(c, "economic_calendar.all")}, {Value: 2, Label: core.Msg(c, "economic_calendar.medium_up")},
				{Value: 3, Label: core.Msg(c, "economic_calendar.high")}}
			input.RadioGroup(c, &s.MinImportance, imp, input.RadioGroupOptions{Horizontal: true, Label: core.Msg(c, "economic_calendar.importance"), Disabled: opts.Disabled})
		})
		rows := economicFilter(s, events, opts.Version)
		cols := []data.DataColumn[EconomicEvent]{
			{ID: "time", Title: core.Msg(c, "economic_calendar.time"), Width: 130, Sortable: true,
				Compare: func(a, b EconomicEvent) int { return a.Time.Compare(b.Time) },
				Text:    func(e EconomicEvent) string { return e.Time.Format("01-02 15:04") }},
			{ID: "country", Title: core.Msg(c, "economic_calendar.country"), Width: 80, Sortable: true,
				Compare: func(a, b EconomicEvent) int { return cmp.Compare(a.Country, b.Country) },
				Text:    func(e EconomicEvent) string { return e.Country }},
			{ID: "importance", Title: core.Msg(c, "economic_calendar.importance"), Width: 110, Sortable: true,
				Compare: func(a, b EconomicEvent) int { return cmp.Compare(a.Importance, b.Importance) },
				Text: func(e EconomicEvent) string {
					return core.Msg(c, fmt.Sprintf("economic_calendar.level_%d", e.Importance))
				},
				Cell: func(e EconomicEvent) { economicImportance(c, e.Importance) }},
			{ID: "event", Title: core.Msg(c, "economic_calendar.event"), Text: func(e EconomicEvent) string { return e.Title }},
			{ID: "actual", Title: core.Msg(c, "economic_calendar.actual"), Width: 90, Align: ui.End,
				Text: func(e EconomicEvent) string { return cmp.Or(e.Actual, "—") },
				Cell: func(e EconomicEvent) {
					ui.Text(c, cmp.Or(e.Actual, "—")).FontWeight(700).FontFeatures("tnum").JustifySelf(ui.End).FillWidth().TextAlign(ui.End)
				}},
			{ID: "forecast", Title: core.Msg(c, "economic_calendar.forecast"), Width: 90, Align: ui.End, Text: func(e EconomicEvent) string { return cmp.Or(e.Forecast, "—") }},
			{ID: "previous", Title: core.Msg(c, "economic_calendar.previous"), Width: 90, Align: ui.End, Text: func(e EconomicEvent) string { return cmp.Or(e.Previous, "—") }},
		}
		v = data.DataTable(c, &s.Table, rows, data.DataTableOptions[EconomicEvent, string]{Columns: cols, Key: func(e EconomicEvent) string { return e.ID },
			Version: opts.Version, Status: opts.Status, Error: opts.Error, Empty: core.Msg(c, "economic_calendar.empty")})
		v.Element.Grow(1)
	})
	v.Element = col
	return v
}

// economicImportance draws one to three diamonds and names the level.
func economicImportance(c *ui.Context, level int) {
	k := core.Tokens(c)
	word := core.Msg(c, fmt.Sprintf("economic_calendar.level_%d", level))
	ui.Row(c).Gap(4).AlignItems(ui.Center).Label(word).Children(func() {
		ui.Box(c).Size(34, 10).Draw(func(p *ui.Painter, r ui.Rect) {
			for i := range 3 {
				d := internal.Diamond(r.X+5+float32(i)*12, r.Y+5, 9)
				if i < level {
					p.FillPath(d, k.AccentText)
				} else {
					p.StrokePath(d, 1, k.ControlBorder)
				}
			}
		})
		ui.Text(c, word).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
	})
}

func economicFilter(s *EconomicCalendarState, events []EconomicEvent, version int) []EconomicEvent {
	key := fmt.Sprint(chart.ChartID(events), version, s.MinImportance, strings.Join(s.Countries, "\x00"))
	if key == s.cache.key {
		return s.cache.rows
	}
	s.cache.key, s.cache.rows = key, nil
	for _, e := range events {
		if e.Importance >= s.MinImportance && (len(s.Countries) == 0 || slices.Contains(s.Countries, e.Country)) {
			s.cache.rows = append(s.cache.rows, e)
		}
	}
	slog.Debug("mujica: economic_calendar filter", "shown", len(s.cache.rows))
	return s.cache.rows
}
