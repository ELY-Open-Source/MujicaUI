package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSymbolSearchResultsAndAsync(t *testing.T) {
	s := SymbolSearchState{}
	rows := []MarketSymbol{{Code: "AAPL", Name: "Apple", Market: "NASDAQ"}}
	var got SymbolSearchView
	opts := SymbolSearchOptions{Markets: []string{"NASDAQ", "NYSE"}, Recent: rows}
	tt := chatTester(640, 450, func(c *ui.Context) {
		v := SymbolSearch(c, &s, rows, opts)
		if v.Picked || v.Changed || v.Retry {
			got = v
		}
	})
	if !tt.HasText("Recent") {
		t.Fatal("recent missing")
	}
	_ = tt.Click("AAPL")
	if !got.Picked || got.Chosen.Code != "AAPL" {
		t.Fatalf("pick %+v", got)
	}
	tt.Key(0, ui.KeyTab)
	s.Query = "apple"
	opts.Loading = true
	rows = nil
	tt.Frame()
	if !tt.HasText("Searching") {
		t.Fatal("loading missing")
	}
	opts.Loading = false
	opts.Error = "Caller failed"
	tt.Frame()
	_ = tt.Click("Retry")
	if !got.Retry {
		t.Fatal("retry missing")
	}
	opts.Error = ""
	tt.Frame()
	if !tt.HasText("No symbols match") {
		t.Fatal("empty missing")
	}
	opts.Disabled = true
	rows = []MarketSymbol{{Code: "AAPL"}}
	got = SymbolSearchView{}
	tt.Frame()
	_ = tt.Click("AAPL")
	if got.Picked {
		t.Fatal("disabled pick")
	}
}
func TestSymbolSearchSnapshot(t *testing.T) {
	var ss SymbolSearchState
	snapshot(t, "symbol_search", 520, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			SymbolSearch(c, &ss, nil, SymbolSearchOptions{Markets: []string{"NASDAQ", "NYSE"}, Recent: []MarketSymbol{{Code: "AAPL", Name: "Apple", Market: "NASDAQ"}}})

		})
	})
}

func TestSymbolSearchArrowNavigation(t *testing.T) {
	s := SymbolSearchState{Query: "tech"}
	rows := []MarketSymbol{{Code: "AAPL", Name: "Apple"}, {Code: "MSFT", Name: "Microsoft"}}
	var chosen string
	tt := chatTester(640, 450, func(c *ui.Context) {
		r := SymbolSearch(c, &s, rows, SymbolSearchOptions{})
		if r.Picked {
			chosen = r.Chosen.Code
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if chosen != "MSFT" {
		t.Fatalf("chosen %q", chosen)
	}
}
