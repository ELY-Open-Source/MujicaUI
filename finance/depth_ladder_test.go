package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func depthLadderBook() OrderBookData {
	return OrderBookData{Bids: []BookLevel{{100.00, 5}, {99.99, 3}}, Asks: []BookLevel{{100.01, 2}, {100.03, 7}}}
}

func TestDepthLadderCentersAndOrders(t *testing.T) {
	var s DepthLadderState
	var price float64
	var side OrderSide
	tt := tradingTester(core.GreenUp, 400, 600, func(c *ui.Context) {
		if p, sd, ok := DepthLadder(c, &s, 100.004, depthLadderBook(), DepthLadderOptions{TickSize: 0.01, Rows: 7}).Clicked(); ok {
			price, side = p, sd
		}
	})
	top, _ := tt.Find("100.03")
	mid, _ := tt.Find("Last price 100.00")
	bottom, _ := tt.Find("99.97")
	if !(top.Y < mid.Y && mid.Y < bottom.Y) || tt.HasText("100.04") {
		t.Fatalf("not centered: %v %v %v %q", top, mid, bottom, tt.Texts())
	}
	if err := tt.Click("Sell at 100.03"); err != nil {
		t.Fatal(err)
	}
	if price != 100.03 || side != SideSell {
		t.Fatalf("clicked %v %v", price, side)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyEnter)
	if price != 100.03 || side != SideBuy {
		t.Fatalf("keyboard %v %v", price, side)
	}
	if !tt.HasText("7") || !tt.HasText("5") {
		t.Fatalf("sizes: %q", tt.Texts())
	}
}

func TestDepthLadderScrollsAndRecenters(t *testing.T) {
	var s DepthLadderState
	tt := tradingTester(core.GreenUp, 400, 600, func(c *ui.Context) {
		DepthLadder(c, &s, 100, depthLadderBook(), DepthLadderOptions{TickSize: 0.01, Rows: 7})
	})
	_ = tt.Click("Buy at 100.00")
	tt.Key(0, ui.KeyPageUp)
	if s.Offset != 3 || !tt.HasText("100.06") {
		t.Fatalf("page up: offset %d %q", s.Offset, tt.Texts())
	}
	tt.Scroll(100, 200, 0, 40)
	if s.Offset != 2 {
		t.Fatalf("wheel: offset %d", s.Offset)
	}
	if err := tt.Click("Center on last price"); err != nil {
		t.Fatal(err)
	}
	if s.Offset != 0 || tt.HasText("Center on last price") {
		t.Fatalf("recenter: %d", s.Offset)
	}
}

// Only the visible rows are built, whatever the book's size.
func TestDepthLadderLargeBook(t *testing.T) {
	var s DepthLadderState
	book := OrderBookData{}
	for i := range 100_000 {
		book.Bids = append(book.Bids, BookLevel{Price: 100 - float64(i)/100, Size: 1})
	}
	tt := tradingTester(core.GreenUp, 400, 600, func(c *ui.Context) {
		DepthLadder(c, &s, 100, book, DepthLadderOptions{TickSize: 0.01, Rows: 9})
	})
	if n := len(tt.Texts()); n > 60 {
		t.Fatalf("built %d texts", n)
	}
	// A sentinel survives only if the index is kept.
	s.cache.bids[-1] = 42
	tt.Frame()
	if s.cache.bids[-1] != 42 {
		t.Fatal("index rebuilt without a change")
	}
}

func TestDepthLadderDisabledAndBadTick(t *testing.T) {
	var s DepthLadderState
	clicked := false
	tt := tradingTester(core.GreenUp, 400, 600, func(c *ui.Context) {
		_, _, clicked = DepthLadder(c, &s, 100, depthLadderBook(), DepthLadderOptions{TickSize: 0.01, Disabled: true}).Clicked()
	})
	_ = tt.Click("Buy at 100.00")
	tt.Scroll(100, 200, 0, 40)
	if clicked || s.Offset != 0 {
		t.Fatal("disabled ladder acted")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("zero tick accepted")
		}
	}()
	tradingTester(core.GreenUp, 100, 100, func(c *ui.Context) { DepthLadder(c, &s, 100, depthLadderBook(), DepthLadderOptions{}) })
}

func TestDepthLadderSnapshot(t *testing.T) {
	snapshot(t, "depth_ladder", 340, 340, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			s := ui.Local(c.Root(), "snap", func() DepthLadderState { return DepthLadderState{} })
			DepthLadder(c, s, 100, depthLadderBook(), DepthLadderOptions{TickSize: 0.01, Rows: 7})
		})
	})
}
