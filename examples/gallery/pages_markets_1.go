package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
	"slices"
	"strings"
)

func marketsOptionChainDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "181.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "181.state", func() mujica.OptionChainState { return mujica.OptionChainState{} })
	rows := *ui.Local(c.Root(), "markets.OptionChain.rows.pages_markets_1", func() []mujica.OptionStrike {
		return []mujica.OptionStrike{{Strike: 100, Call: mujica.OptionQuote{Bid: 3, Ask: 3.2, Last: 3.1}, Put: mujica.OptionQuote{Bid: 2, Ask: 2.2, Last: 2.1}}}
	})

	precision := ui.Local(c.Root(), "markets.OptionChain.precision", func() float64 { return 2 })
	if mujica.Slider(c, precision, mujica.SliderOptions{Min: 1, Max: 8, Step: 1, Label: "Decimals", ShowValue: true}).Changed() {
		log(fmt.Sprint("decimals ", *precision))
	}
	spot := ui.Local(c.Root(), "181.spot", func() float64 { return 100.4 })
	if mujica.Slider(c, spot, mujica.SliderOptions{Min: 95, Max: 105, Step: 0.1, Label: "Spot", ShowValue: true}).Changed() {
		log(fmt.Sprint("spot ", *spot))
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.OptionChain(c, state, []string{"Oct 17", "Nov 21"}, rows, mujica.OptionChainOptions{Decimals: int(*precision), Spot: *spot})
		if r.Changed || r.Submitted || r.ExpiryChanged {
			log(fmt.Sprint(state.Expiry, state.Strike, state.Side))
		}
	})
}

func marketsGreeksTableDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "182.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	rows := *ui.Local(c.Root(), "markets.GreeksTable.rows.pages_markets_1", func() []mujica.GreeksRow {
		return []mujica.GreeksRow{{Label: "AAPL 100 C", Greeks: mujica.OptionGreeks{Delta: 0.5, Gamma: 0.03, Theta: -0.1, Vega: 0.2, Rho: 0.05}}}
	})

	precision := ui.Local(c.Root(), "markets.GreeksTable.precision", func() float64 { return 2 })
	if mujica.Slider(c, precision, mujica.SliderOptions{Min: 1, Max: 8, Step: 1, Label: "Decimals", ShowValue: true}).Changed() {
		log(fmt.Sprint("decimals ", *precision))
	}
	notes := ui.Local(c.Root(), "182.notes", func() bool { return true })
	if mujica.Checkbox(c, notes, "Notes", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("notes ", *notes))
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		mujica.GreeksTable(c, rows, mujica.GreeksTableOptions{Decimals: int(*precision), Notes: *notes})
	})
}

func marketsPayoffDiagramDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "183.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	legs := *ui.Local(c.Root(), "markets.PayoffDiagram.legs.pages_markets_1", func() []mujica.OptionLeg {
		return []mujica.OptionLeg{{Kind: mujica.LegCall, Strike: 100, Premium: 3, Quantity: 1}}
	})

	status := ui.Local(c.Root(), "183.status", func() int { return 0 })
	if mujica.RadioGroup(c, status, []mujica.RadioItem[int]{{Value: 0, Label: "Ready"}, {Value: 1, Label: "Loading"}, {Value: 2, Label: "Empty"}, {Value: 3, Label: "Error"}}, mujica.RadioGroupOptions{Horizontal: true}).Changed() {
		log(fmt.Sprint("status ", *status))
	}
	frame := mujica.ChartFrame{Label: "PayoffDiagram"}
	if *status == 1 {
		frame.Status = mujica.DataLoading
	}
	if *status == 3 {
		frame.Status = mujica.DataFailed
		frame.Error = "Sample failed"
	}
	if *status == 2 {
		legs = nil
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.PayoffDiagram(c, legs, mujica.PayoffDiagramOptions{Frame: frame, Spot: 103})
		if r.Changed() || r.Retried() {
			log(fmt.Sprint(r.Point))
		}
	})
}

func marketsWatchlistDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "184.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "184.state", func() mujica.WatchlistState { return mujica.WatchlistState{} })
	groups := *ui.Local(c.Root(), "markets.Watchlist.groups.pages_markets_1", func() []mujica.WatchGroup {
		return []mujica.WatchGroup{{Name: "Technology", Items: []mujica.WatchItem{{Symbol: "AAPL", Name: "Apple", Price: 227.5, Change: 2.7, ChangePercent: 1.2}, {Symbol: "MSFT", Name: "Microsoft", Price: 418.2, Change: -2.5, ChangePercent: -0.6}}}, {Name: "Energy", Items: []mujica.WatchItem{{Symbol: "XOM", Price: 117.4}}}}
	})

	precision := ui.Local(c.Root(), "markets.Watchlist.precision", func() float64 { return 2 })
	if mujica.Slider(c, precision, mujica.SliderOptions{Min: 1, Max: 8, Step: 1, Label: "Decimals", ShowValue: true}).Changed() {
		log(fmt.Sprint("decimals ", *precision))
	}
	if mujica.Button(c, "Refresh quotes", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
		groups[0].Items = slices.Clone(groups[0].Items)
		groups[0].Items[0].Price += 0.1
		log("quotes refreshed")
	}
	beforeGroup, beforeColumns := state.Group, fmt.Sprint(state.Hidden)
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.Watchlist(c, state, groups, mujica.WatchlistOptions{Decimals: int(*precision)})
		if r.Changed || r.Submitted || r.Moved {
			log(fmt.Sprint(state.Selected, r.Rows, r.To))
		}
		if r.Moved {
			items := groups[state.Group].Items
			var moving []mujica.WatchItem
			var rest []mujica.WatchItem
			at := r.To
			for j, item := range items {
				if slices.Contains(r.Rows, j) {
					moving = append(moving, item)
					if j < r.To {
						at--
					}
				} else {
					rest = append(rest, item)
				}
			}
			groups[state.Group].Items = slices.Insert(rest, at, moving...)
		}
	})
	if state.Group != beforeGroup || fmt.Sprint(state.Hidden) != beforeColumns {
		log(fmt.Sprint("group/columns ", state.Group, state.Hidden))
	}
}

func marketsSymbolSearchDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "185.disabled", func() bool { return false })
	if mujica.Checkbox(c, disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "185.state", func() mujica.SymbolSearchState { return mujica.SymbolSearchState{} })
	rows := *ui.Local(c.Root(), "markets.SymbolSearch.rows.pages_markets_1", func() []mujica.MarketSymbol {
		return []mujica.MarketSymbol{{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}}
	})

	status, empty, err := marketsDataState(c, "185", log)
	if empty || status != mujica.DataReady {
		rows = nil
	}
	results := []mujica.MarketSymbol{}
	for _, row := range rows {
		if (state.Market == "" || state.Market == row.Market) && strings.Contains(strings.ToUpper(row.Code+" "+row.Name), strings.ToUpper(state.Query)) {
			results = append(results, row)
		}
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := mujica.SymbolSearch(c, state, results, mujica.SymbolSearchOptions{Loading: status == mujica.DataLoading, Error: err, Markets: []string{"NASDAQ", "NYSE"}, Recent: rows})
		if r.Changed || r.Picked || r.Retry {
			log(fmt.Sprint(state.Query, state.Market, r.Chosen.Code))
		}
	})
}
