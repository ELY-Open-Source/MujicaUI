package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func optionChainTester(s *OptionChainState, strikes []OptionStrike, opts OptionChainOptions, got *OptionChainView) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Height(420).Children(func() {
			v := OptionChain(c, s, []string{"Oct 17", "Nov 21"}, strikes, opts)
			if v.Changed || v.Submitted || v.ExpiryChanged {
				*got = v
			}
		})
	}, 780, 440)
	// Rows show once the list has been laid out.
	tt.Frame()
	return tt
}

func TestOptionChainChooseByClickAndKeys(t *testing.T) {
	var s OptionChainState
	var got OptionChainView
	tt := optionChainTester(&s, mktStrikes(), OptionChainOptions{Spot: 100.4}, &got)
	tt.Frame()
	optionChainClick(t, tt, "102.00", OptionPut)
	if !s.Chosen || s.Strike != 102 || s.Side != OptionPut || !got.Changed {
		t.Fatalf("state %+v", s)
	}
	tt.Key(0, ui.KeyLeft)
	if s.Side != OptionCall || s.Strike != 102 {
		t.Fatalf("Left kept side %v", s.Side)
	}
	tt.Key(0, ui.KeyDown)
	if s.Strike != 103 {
		t.Fatalf("Down chose %v", s.Strike)
	}
	tt.Key(0, ui.KeyEnter)
	if !got.Submitted {
		t.Fatal("Enter did not submit")
	}
}

func TestOptionChainMarksInTheMoneyAndCenters(t *testing.T) {
	var strikes []OptionStrike
	for i := range 300 {
		strikes = append(strikes, OptionStrike{Strike: float64(i) + 0.5})
	}
	var s OptionChainState
	var got OptionChainView
	tt := optionChainTester(&s, strikes, OptionChainOptions{Spot: 200}, &got)
	tt.Frame()
	if !tt.HasText("200.50") || tt.HasText("Call 5.50") {
		t.Fatalf("not centered on the spot: %q", tt.Texts())
	}
	if n := len(tt.Texts()); n > 400 {
		t.Fatalf("built %d texts for 300 strikes", n)
	}
	if !tt.HasText("Call 199.50: bid 0.00, ask 0.00, last 0.00, in the money") || tt.HasText("Put 199.50: bid 0.00, ask 0.00, last 0.00, in the money") {
		t.Fatal("in-the-money marks wrong")
	}
}

func TestOptionChainExpiryAndDisabled(t *testing.T) {
	var s OptionChainState
	var got OptionChainView
	tt := optionChainTester(&s, mktStrikes(), OptionChainOptions{Spot: 100}, &got)
	_ = tt.Click("Oct 17")
	_ = tt.Click("Nov 21")
	if s.Expiry != "Nov 21" || !got.ExpiryChanged {
		t.Fatalf("expiry %q", s.Expiry)
	}
	var d OptionChainState
	got = OptionChainView{}
	tt = optionChainTester(&d, mktStrikes(), OptionChainOptions{Spot: 100, Disabled: true}, &got)
	optionChainClick(t, tt, "102.00", OptionPut)
	if d.Chosen || got.Changed {
		t.Fatal("disabled chain chose")
	}
}

func TestOptionChainRejectsUnsortedStrikes(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	var s OptionChainState
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		OptionChain(c, &s, []string{"x"}, []OptionStrike{{Strike: 2}, {Strike: 1}}, OptionChainOptions{})
	}, 100, 100, 1)
}

func TestOptionChainSnapshot(t *testing.T) {
	snapshot(t, "option_chain", 760, 400, func(c *ui.Context) {
		s := &OptionChainState{Strike: 99, Side: OptionCall, Chosen: true}
		ui.Column(c).Padding(16).Height(400).Children(func() {
			OptionChain(c, s, []string{"Oct 17", "Nov 21"}, mktStrikes(), OptionChainOptions{Spot: 100.4})
		})
	})
}

// optionChainClick clicks the call or put half of the row for strike.
func optionChainClick(t *testing.T, tt *ui.Tester, strike string, side OptionSide) {
	t.Helper()
	r, ok := tt.Find(strike)
	if !ok {
		t.Fatalf("%q not shown", strike)
	}
	x := r.X + r.W*0.2
	if side == OptionPut {
		x = r.X + r.W*0.8
	}
	tt.ClickAt(x, r.Y+r.H/2)
}

func TestOptionChainSameStrikeOtherSide(t *testing.T) {
	var s OptionChainState
	var got OptionChainView
	tt := optionChainTester(&s, mktStrikes(), OptionChainOptions{Spot: 100}, &got)
	optionChainClick(t, tt, "102.00", OptionCall)
	optionChainClick(t, tt, "102.00", OptionPut)
	if s.Strike != 102 || s.Side != OptionPut {
		t.Fatalf("contract %+v", s)
	}
}
