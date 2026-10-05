package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func orderBookSample() OrderBookData {
	return OrderBookData{
		Bids: []BookLevel{{100.00, 5}, {99.99, 3}, {99.94, 2}, {99.90, 4}},
		Asks: []BookLevel{{100.01, 1}, {100.02, 2}, {100.07, 6}},
	}
}

func TestOrderBookClickReportsPrice(t *testing.T) {
	var s OrderBookState
	var price float64
	var side OrderSide
	tt := tradingTester(core.GreenUp, 400, 500, func(c *ui.Context) {
		if p, sd, ok := OrderBook(c, &s, orderBookSample(), OrderBookOptions{}).Clicked(); ok {
			price, side = p, sd
		}
	})
	if err := tt.Click("Ask 100.02, size 2"); err != nil {
		t.Fatal(err)
	}
	if price != 100.02 || side != SideSell {
		t.Fatalf("clicked %v %v", price, side)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if price != 100.01 || side != SideSell {
		t.Fatalf("keyboard %v %v", price, side)
	}
	if !tt.HasText("0.01 (0.010%)") {
		t.Fatalf("spread missing: %q", tt.Texts())
	}
}

func TestOrderBookGroupsLevelsByTicks(t *testing.T) {
	s := OrderBookState{Group: 0.05}
	tt := tradingTester(core.GreenUp, 400, 500, func(c *ui.Context) {
		OrderBook(c, &s, orderBookSample(), OrderBookOptions{Groups: []float64{0.01, 0.05, 0.1}})
	})
	// Bids round down, asks round up.
	want := []BookLevel{{100.00, 5}, {99.95, 3}, {99.90, 6}}
	for i, l := range want {
		if s.cache.bids[i] != l {
			t.Fatalf("bids %v", s.cache.bids)
		}
	}
	if s.cache.asks[0] != (BookLevel{100.05, 3}) || s.cache.asks[1] != (BookLevel{100.10, 6}) {
		t.Fatalf("asks %v", s.cache.asks)
	}
	if !tt.HasText("99.95") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestOrderBookAggregatesOnlyOnChange(t *testing.T) {
	var s OrderBookState
	book := orderBookSample()
	tt := tradingTester(core.GreenUp, 400, 500, func(c *ui.Context) { OrderBook(c, &s, book, OrderBookOptions{}) })
	first := &s.cache.bids[0]
	tt.Frame()
	if &s.cache.bids[0] != first {
		t.Fatal("regrouped without a change")
	}
	book.Bids[0].Size = 9
	book.Version++
	tt.Frame()
	if s.cache.bids[0].Size != 9 {
		t.Fatal("version change ignored")
	}
}

func TestOrderBookDisabledAndBadGroup(t *testing.T) {
	var s OrderBookState
	clicked := false
	tt := tradingTester(core.GreenUp, 400, 500, func(c *ui.Context) {
		_, _, clicked = OrderBook(c, &s, orderBookSample(), OrderBookOptions{Disabled: true}).Clicked()
	})
	_ = tt.Click("Bid 100.00, size 5")
	if clicked {
		t.Fatal("disabled book clicked")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("off-tick group accepted")
		}
	}()
	bad := OrderBookState{Group: 0.015}
	tradingTester(core.GreenUp, 100, 100, func(c *ui.Context) {
		OrderBook(c, &bad, orderBookSample(), OrderBookOptions{Groups: []float64{0.015}})
	})
}

func TestOrderBookSnapshot(t *testing.T) {
	snapshot(t, "order_book", 320, 420, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			s := ui.Local(c.Root(), "snap", func() OrderBookState { return OrderBookState{} })
			OrderBook(c, s, orderBookSample(), OrderBookOptions{Depth: 4, Groups: []float64{0.01, 0.05}})
		})
	})
}

// A 100,000-level book builds only the shown depth.
func TestOrderBookLargeBook(t *testing.T) {
	var s OrderBookState
	book := OrderBookData{}
	for i := range 100_000 {
		book.Bids = append(book.Bids, BookLevel{Price: 1000 - float64(i)/100, Size: 1})
		book.Asks = append(book.Asks, BookLevel{Price: 1000.01 + float64(i)/100, Size: 1})
	}
	tt := tradingTester(core.GreenUp, 400, 900, func(c *ui.Context) { OrderBook(c, &s, book, OrderBookOptions{}) })
	if n := len(tt.Texts()); n > 80 {
		t.Fatalf("built %d texts", n)
	}
	if len(s.cache.bids) != 100_000 || s.cache.bids[0].Price != 1000 {
		t.Fatalf("bids %d", len(s.cache.bids))
	}
}
