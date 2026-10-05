package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
	"time"
)

func marketsMarketStatusDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "191.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)

	phase := ui.Local(c.Root(), "191.phase", func() mujica.MarketPhase { return mujica.MarketOpen })
	if mujica.RadioGroup(c, phase, []mujica.RadioItem[mujica.MarketPhase]{{Value: mujica.MarketClosed, Label: "Closed"}, {Value: mujica.MarketPreMarket, Label: "Pre-market"}, {Value: mujica.MarketOpen, Label: "Open"}, {Value: mujica.MarketAfterHours, Label: "After hours"}, {Value: mujica.MarketHoliday, Label: "Holiday"}}, mujica.RadioGroupOptions{Horizontal: true}).Changed() {
		log(fmt.Sprint(*phase))
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		mujica.MarketStatus(c, *phase, mujica.MarketStatusOptions{Market: "NYSE", Now: now, Next: now.Add(5 * time.Hour), NextPhase: mujica.MarketAfterHours})
	})
}

func marketsTradingSessionClockDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "192.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
	markets := *ui.Local(c.Root(), "markets.TradingSessionClock.markets.pages_markets_3", func() []mujica.SessionMarket {
		return []mujica.SessionMarket{{Name: "NYSE", Location: time.FixedZone("EDT", -4*3600), Sessions: []mujica.TradingSession{{Phase: mujica.MarketOpen, Start: 570, End: 960}}}}
	})

	holiday := ui.Local(c.Root(), "192.holiday", func() bool { return false })
	if mujica.Checkbox(c, holiday, "Holiday", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint(*holiday))
	}
	markets[0].Holiday = *holiday
	minute := ui.Local(c.Root(), "192.minute", func() float64 { return 870 })
	if mujica.Slider(c, minute, mujica.SliderOptions{Min: 0, Max: 1439, Step: 1, Label: "UTC minute", ShowValue: true}).Changed() {
		log(fmt.Sprint(*minute))
	}
	now = now.Truncate(24 * time.Hour).Add(time.Duration(*minute) * time.Minute)
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		mujica.TradingSessionClock(c, markets, mujica.TradingSessionClockOptions{Now: now})
	})
}

func marketsEconomicCalendarDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "193.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "193.state", func() mujica.EconomicCalendarState { return mujica.EconomicCalendarState{} })
	rows := *ui.Local(c.Root(), "markets.EconomicCalendar.rows.pages_markets_3", func() []mujica.EconomicEvent {
		return []mujica.EconomicEvent{{ID: "jobs", Country: "US", Title: "Payrolls", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Importance: 3, Forecast: "150K", Previous: "142K", Actual: "254K"}}
	})

	status, empty, err := marketsDataState(c, "193", log)
	if empty || status != mujica.DataReady {
		rows = nil
	}
	beforeFilter := fmt.Sprint(state.Countries, state.MinImportance)
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.EconomicCalendar(c, state, rows, mujica.EconomicCalendarOptions{Status: status, Error: err})
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
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "194.state", func() mujica.DataTableState[string] { return mujica.DataTableState[string]{} })
	rows := *ui.Local(c.Root(), "markets.EarningsCalendar.rows.pages_markets_3", func() []mujica.EarningsEvent {
		return []mujica.EarningsEvent{{ID: "a", Symbol: "AAPL", Name: "Apple", Date: mujica.Date{Year: 2026, Month: 10, Day: 30}, Timing: mujica.EarningsAfterClose, EPSEstimate: 1.6, EPSActual: 1.64}}
	})

	status, empty, err := marketsDataState(c, "194", log)
	if empty || status != mujica.DataReady {
		rows = nil
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.EarningsCalendar(c, state, rows, mujica.EarningsCalendarOptions{Status: status, Error: err})
		if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
			log(fmt.Sprint(state.Selected()))
		}
	})
}

func marketsNewsFeedDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "195.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "195.state", func() mujica.NewsFeedState { return mujica.NewsFeedState{} })
	rows := *ui.Local(c.Root(), "markets.NewsFeed.rows.pages_markets_3", func() []mujica.NewsItem {
		return []mujica.NewsItem{{ID: "a", Title: "Apple unveils new chips", Source: "Wire", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Symbols: []string{"AAPL"}}}
	})

	status, empty, err := marketsDataState(c, "195", log)
	if empty || status != mujica.DataReady {
		rows = nil
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.NewsFeed(c, state, rows, mujica.NewsFeedOptions{Status: status, Error: err})
		if r.Opened != "" || r.Symbol != "" || r.Retry || r.ReadChanged {
			log(fmt.Sprint(r.Opened, r.Symbol))
		}
	})
}
