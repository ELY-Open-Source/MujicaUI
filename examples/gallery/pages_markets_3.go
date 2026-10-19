package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/finance"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func marketsMarketStatusDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "191.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)

	phase := ui.Local(c.Root(), "191.phase", func() finance.MarketPhase { return finance.MarketOpen })
	if input.RadioGroup(c, phase, []input.RadioItem[finance.MarketPhase]{{Value: finance.MarketClosed, Label: "Closed"}, {Value: finance.MarketPreMarket, Label: "Pre-market"}, {Value: finance.MarketOpen, Label: "Open"}, {Value: finance.MarketAfterHours, Label: "After hours"}, {Value: finance.MarketHoliday, Label: "Holiday"}}, input.RadioGroupOptions{Horizontal: true}).Changed() {
		log(fmt.Sprint(*phase))
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		finance.MarketStatus(c, *phase, finance.MarketStatusOptions{Market: "NYSE", Now: now, Next: now.Add(5 * time.Hour), NextPhase: finance.MarketAfterHours})
	})
}

func marketsTradingSessionClockDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "192.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
	markets := *ui.Local(c.Root(), "markets.TradingSessionClock.markets.pages_markets_3", func() []finance.SessionMarket {
		return []finance.SessionMarket{{Name: "NYSE", Location: time.FixedZone("EDT", -4*3600), Sessions: []finance.TradingSession{{Phase: finance.MarketOpen, Start: 570, End: 960}}}}
	})

	holiday := ui.Local(c.Root(), "192.holiday", func() bool { return false })
	if input.Checkbox(c, holiday, "Holiday", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint(*holiday))
	}
	markets[0].Holiday = *holiday
	minute := ui.Local(c.Root(), "192.minute", func() float64 { return 870 })
	if input.Slider(c, minute, input.SliderOptions{Min: 0, Max: 1439, Step: 1, Label: "UTC minute", ShowValue: true}).Changed() {
		log(fmt.Sprint(*minute))
	}
	now = now.Truncate(24 * time.Hour).Add(time.Duration(*minute) * time.Minute)
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		finance.TradingSessionClock(c, markets, finance.TradingSessionClockOptions{Now: now})
	})
}

func marketsEconomicCalendarDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "193.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "193.state", func() finance.EconomicCalendarState { return finance.EconomicCalendarState{} })
	rows := *ui.Local(c.Root(), "markets.EconomicCalendar.rows.pages_markets_3", func() []finance.EconomicEvent {
		return []finance.EconomicEvent{{ID: "jobs", Country: "US", Title: "Payrolls", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Importance: 3, Forecast: "150K", Previous: "142K", Actual: "254K"}}
	})

	status, empty, err := marketsDataState(c, "193", log)
	if empty || status != data.DataReady {
		rows = nil
	}
	beforeFilter := fmt.Sprint(state.Countries, state.MinImportance)
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.EconomicCalendar(c, state, rows, finance.EconomicCalendarOptions{Status: status, Error: err})
		if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
			log(fmt.Sprint(state.Table.Selected()))
		}
	})
	if fmt.Sprint(state.Countries, state.MinImportance) != beforeFilter {
		log(fmt.Sprint("filter ", state.Countries, state.MinImportance))
	}
}

func marketsEarningsCalendarDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "194.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "194.state", func() data.DataTableState[string] { return data.DataTableState[string]{} })
	rows := *ui.Local(c.Root(), "markets.EarningsCalendar.rows.pages_markets_3", func() []finance.EarningsEvent {
		return []finance.EarningsEvent{{ID: "a", Symbol: "AAPL", Name: "Apple", Date: core.Date{Year: 2026, Month: 10, Day: 30}, Timing: finance.EarningsAfterClose, EPSEstimate: 1.6, EPSActual: 1.64}}
	})

	status, empty, err := marketsDataState(c, "194", log)
	if empty || status != data.DataReady {
		rows = nil
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.EarningsCalendar(c, state, rows, finance.EarningsCalendarOptions{Status: status, Error: err})
		if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
			log(fmt.Sprint(state.Selected()))
		}
	})
}

func marketsNewsFeedDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "195.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "195.state", func() finance.NewsFeedState { return finance.NewsFeedState{} })
	rows := *ui.Local(c.Root(), "markets.NewsFeed.rows.pages_markets_3", func() []finance.NewsItem {
		return []finance.NewsItem{{ID: "a", Title: "Apple unveils new chips", Source: "Wire", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Symbols: []string{"AAPL"}}}
	})

	status, empty, err := marketsDataState(c, "195", log)
	if empty || status != data.DataReady {
		rows = nil
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.NewsFeed(c, state, rows, finance.NewsFeedOptions{Status: status, Error: err})
		if r.Opened != "" || r.Symbol != "" || r.Retry || r.ReadChanged {
			log(fmt.Sprint(r.Opened, r.Symbol))
		}
	})
}
