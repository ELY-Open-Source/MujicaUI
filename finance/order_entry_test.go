package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

type orderEntryProbe struct {
	s       OrderEntryState
	opts    OrderEntryOptions
	orders  []OrderDraft
	checks  []uint64
	checked []OrderDraft
}

func (p *orderEntryProbe) tester() *ui.Tester {
	return tradingTester(core.GreenUp, 520, 700, func(c *ui.Context) {
		if o, ok := OrderEntry(c, &p.s, p.opts).Submitted(); ok {
			p.orders = append(p.orders, o)
		}
	})
}

func TestOrderEntryLimitSubmitOnce(t *testing.T) {
	p := &orderEntryProbe{opts: OrderEntryOptions{Symbol: "MUJC", Last: 100}}
	tt := p.tester()
	_ = tt.Click("Quantity")
	tt.Type("10")
	_ = tt.Click("Limit price")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("99.5")
	if !tt.HasText("≈ 995.00") {
		t.Fatalf("estimate: %q", tt.Texts())
	}
	if err := tt.Click("Buy MUJC"); err != nil {
		t.Fatal(err)
	}
	_ = tt.Click("Buy MUJC")
	if len(p.orders) != 1 || p.orders[0] != (OrderDraft{Side: SideBuy, Type: OrderLimit, Quantity: 10, Price: 99.5}) {
		t.Fatalf("orders %+v", p.orders)
	}
	p.s.EndSubmit("")
	tt.Frame()
	_ = tt.Click("Buy MUJC")
	if len(p.orders) != 2 {
		t.Fatalf("after EndSubmit: %d orders", len(p.orders))
	}
}

func TestOrderEntryValidatesByType(t *testing.T) {
	p := &orderEntryProbe{opts: OrderEntryOptions{Symbol: "MUJC", Last: 100, MaxQuantity: 50}}
	tt := p.tester()
	_ = tt.Click("Buy MUJC")
	if len(p.orders) != 0 || !tt.HasText("Enter a quantity above zero.") || !tt.HasText("Required") {
		t.Fatalf("empty order: %q", tt.Texts())
	}
	if !tt.Focused("Quantity") {
		t.Fatal("first error not focused")
	}
	tt.Type("80")
	tt.Key(0, ui.KeyTab)
	if !tt.HasText("At most 50.") {
		t.Fatalf("max: %q", tt.Texts())
	}
	p.s.Draft.Type, p.s.Draft.Side, p.s.Draft.Quantity, p.s.Draft.Trigger = OrderStop, SideSell, 5, 101
	tt.Frame()
	_ = tt.Click("Sell MUJC")
	if len(p.orders) != 0 || !tt.HasText("A sell stop triggers below the last price.") || tt.HasText("Limit price") {
		t.Fatalf("stop: %q", tt.Texts())
	}
	p.s.Draft.Trigger = 98
	tt.Frame()
	_ = tt.Click("Sell MUJC")
	if len(p.orders) != 1 || p.orders[0].Type != OrderStop || p.orders[0].Trigger != 98 {
		t.Fatalf("orders %+v", p.orders)
	}
}

// An old draft's check cannot clear a newer draft.
func TestOrderEntryAsyncCheckIsVersioned(t *testing.T) {
	p := &orderEntryProbe{}
	p.opts = OrderEntryOptions{Symbol: "MUJC", Last: 100, Check: func(o OrderDraft, v uint64) {
		p.checks, p.checked = append(p.checks, v), append(p.checked, o)
	}}
	p.s.started, p.s.Draft = true, OrderDraft{Side: SideBuy, Type: OrderMarket, Quantity: 10}
	tt := p.tester()
	_ = tt.Click("Buy MUJC")
	if len(p.checks) != 1 || len(p.orders) != 0 {
		t.Fatalf("checks %v orders %d", p.checks, len(p.orders))
	}
	p.s.Draft.Quantity = 20
	tt.Frame()
	if p.s.Resolve(p.checks[0], "") {
		t.Fatal("stale check accepted")
	}
	tt.Frame()
	if len(p.checks) != 2 || p.checked[1].Quantity != 20 || len(p.orders) != 0 {
		t.Fatalf("checks %v orders %d", p.checks, len(p.orders))
	}
	p.s.Resolve(p.checks[1], "Not enough buying power.")
	tt.Frame()
	tt.Frame()
	if len(p.orders) != 0 || !tt.HasText("Not enough buying power.") {
		t.Fatalf("failed check submitted: %q", tt.Texts())
	}
}

func TestOrderEntryDisabledAndReset(t *testing.T) {
	p := &orderEntryProbe{opts: OrderEntryOptions{Symbol: "MUJC", Disabled: true, Initial: OrderDraft{Quantity: 1, Price: 50}}}
	tt := p.tester()
	_ = tt.Click("Quantity")
	tt.Type("9")
	_ = tt.Click("Buy MUJC")
	if p.s.Draft.Quantity != 1 || len(p.orders) != 0 {
		t.Fatalf("disabled entry changed: %+v", p.s.Draft)
	}
	p.opts.Disabled = false
	p.s.Draft.Quantity, p.s.Draft.Side = 7, SideSell
	tt.Frame()
	_ = tt.Click("Reset")
	if p.s.Draft != (OrderDraft{Side: SideBuy, Quantity: 1, Price: 50}) {
		t.Fatalf("reset: %+v", p.s.Draft)
	}
}

func TestOrderEntrySnapshot(t *testing.T) {
	snapshot(t, "order_entry", 420, 520, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			s := ui.Local(c.Root(), "snap", func() OrderEntryState { return OrderEntryState{} })
			if !s.started {
				s.started, s.Draft = true, OrderDraft{Side: SideBuy, Type: OrderStopLimit, Quantity: 25, Price: 101.5, Trigger: 101}
			}
			OrderEntry(c, s, OrderEntryOptions{Symbol: "MUJC", Last: 100.2})
		})
	})
}
