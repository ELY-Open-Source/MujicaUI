package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/finance"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

func marketsOptionChainDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "181.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "181.state", func() finance.OptionChainState { return finance.OptionChainState{} })
	rows := *ui.Local(c.Root(), "markets.OptionChain.rows.pages_markets_1", func() []finance.OptionStrike {
		return []finance.OptionStrike{{Strike: 100, Call: finance.OptionQuote{Bid: 3, Ask: 3.2, Last: 3.1}, Put: finance.OptionQuote{Bid: 2, Ask: 2.2, Last: 2.1}}}
	})

	precision := ui.Local(c.Root(), "markets.OptionChain.precision", func() float64 { return 2 })
	if input.Slider(c, precision, input.SliderOptions{Min: 1, Max: 8, Step: 1, Label: "Decimals", ShowValue: true}).Changed() {
		log(fmt.Sprint("decimals ", *precision))
	}
	spot := ui.Local(c.Root(), "181.spot", func() float64 { return 100.4 })
	if input.Slider(c, spot, input.SliderOptions{Min: 95, Max: 105, Step: 0.1, Label: "Spot", ShowValue: true}).Changed() {
		log(fmt.Sprint("spot ", *spot))
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.OptionChain(c, state, []string{"Oct 17", "Nov 21"}, rows, finance.OptionChainOptions{Decimals: int(*precision), Spot: *spot})
		if r.Changed || r.Submitted || r.ExpiryChanged {
			log(fmt.Sprint(state.Expiry, state.Strike, state.Side))
		}
	})
}

func marketsGreeksTableDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "182.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	rows := *ui.Local(c.Root(), "markets.GreeksTable.rows.pages_markets_1", func() []finance.GreeksRow {
		return []finance.GreeksRow{{Label: "AAPL 100 C", Greeks: finance.OptionGreeks{Delta: 0.5, Gamma: 0.03, Theta: -0.1, Vega: 0.2, Rho: 0.05}}}
	})

	precision := ui.Local(c.Root(), "markets.GreeksTable.precision", func() float64 { return 2 })
	if input.Slider(c, precision, input.SliderOptions{Min: 1, Max: 8, Step: 1, Label: "Decimals", ShowValue: true}).Changed() {
		log(fmt.Sprint("decimals ", *precision))
	}
	notes := ui.Local(c.Root(), "182.notes", func() bool { return true })
	if input.Checkbox(c, notes, "Notes", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("notes ", *notes))
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		finance.GreeksTable(c, rows, finance.GreeksTableOptions{Decimals: int(*precision), Notes: *notes})
	})
}

func marketsPayoffDiagramDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "183.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	legs := *ui.Local(c.Root(), "markets.PayoffDiagram.legs.pages_markets_1", func() []finance.OptionLeg {
		return []finance.OptionLeg{{Kind: finance.LegCall, Strike: 100, Premium: 3, Quantity: 1}}
	})

	status := ui.Local(c.Root(), "183.status", func() int { return 0 })
	if input.RadioGroup(c, status, []input.RadioItem[int]{{Value: 0, Label: "Ready"}, {Value: 1, Label: "Loading"}, {Value: 2, Label: "Empty"}, {Value: 3, Label: "Error"}}, input.RadioGroupOptions{Horizontal: true}).Changed() {
		log(fmt.Sprint("status ", *status))
	}
	frame := chart.ChartFrame{Label: "PayoffDiagram"}
	if *status == 1 {
		frame.Status = data.DataLoading
	}
	if *status == 3 {
		frame.Status = data.DataFailed
		frame.Error = "Sample failed"
	}
	if *status == 2 {
		legs = nil
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.PayoffDiagram(c, legs, finance.PayoffDiagramOptions{Frame: frame, Spot: 103})
		if r.Changed() || r.Retried() {
			log(fmt.Sprint(r.Point))
		}
	})
}

func marketsWatchlistDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "184.disabled", func() bool { return false })
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "184.state", func() finance.WatchlistState { return finance.WatchlistState{} })
	groups := *ui.Local(c.Root(), "markets.Watchlist.groups.pages_markets_1", func() []finance.WatchGroup {
		return []finance.WatchGroup{{Name: "Technology", Items: []finance.WatchItem{{Symbol: "AAPL", Name: "Apple", Price: 227.5, Change: 2.7, ChangePercent: 1.2}, {Symbol: "MSFT", Name: "Microsoft", Price: 418.2, Change: -2.5, ChangePercent: -0.6}}}, {Name: "Energy", Items: []finance.WatchItem{{Symbol: "XOM", Price: 117.4}}}}
	})

	precision := ui.Local(c.Root(), "markets.Watchlist.precision", func() float64 { return 2 })
	if input.Slider(c, precision, input.SliderOptions{Min: 1, Max: 8, Step: 1, Label: "Decimals", ShowValue: true}).Changed() {
		log(fmt.Sprint("decimals ", *precision))
	}
	if input.Button(c, "Refresh quotes", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
		groups[0].Items = slices.Clone(groups[0].Items)
		groups[0].Items[0].Price += 0.1
		log("quotes refreshed")
	}
	beforeGroup, beforeColumns := state.Group, fmt.Sprint(state.Hidden)
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.Watchlist(c, state, groups, finance.WatchlistOptions{Decimals: int(*precision)})
		if r.Changed || r.Submitted || r.Moved {
			log(fmt.Sprint(state.Selected, r.Rows, r.To))
		}
		if r.Moved {
			items := groups[state.Group].Items
			var moving []finance.WatchItem
			var rest []finance.WatchItem
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
	if input.Checkbox(c, disabled, "Disabled", input.CheckboxOptions{}).Changed() {
		log(fmt.Sprint("disabled ", *disabled))
	}
	state := ui.Local(c.Root(), "185.state", func() finance.SymbolSearchState { return finance.SymbolSearchState{} })
	rows := *ui.Local(c.Root(), "markets.SymbolSearch.rows.pages_markets_1", func() []finance.MarketSymbol {
		return []finance.MarketSymbol{{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}}
	})

	status, empty, err := marketsDataState(c, "185", log)
	if empty || status != data.DataReady {
		rows = nil
	}
	results := []finance.MarketSymbol{}
	for _, row := range rows {
		if (state.Market == "" || state.Market == row.Market) && strings.Contains(strings.ToUpper(row.Code+" "+row.Name), strings.ToUpper(state.Query)) {
			results = append(results, row)
		}
	}
	ui.Column(c).Disabled(*disabled).MinHeight(0).Height(400).Children(func() {
		r := finance.SymbolSearch(c, state, results, finance.SymbolSearchOptions{Loading: status == data.DataLoading, Error: err, Markets: []string{"NASDAQ", "NYSE"}, Recent: rows})
		if r.Changed || r.Picked || r.Retry {
			log(fmt.Sprint(state.Query, state.Market, r.Chosen.Code))
		}
	})
}
