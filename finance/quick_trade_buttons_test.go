package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

type quickTradeProbe struct {
	qty    float64
	opts   QuickTradeButtonsOptions
	trades []OrderSide
	sizes  []float64
}

func (p *quickTradeProbe) tester() *ui.Tester {
	return tradingTester(core.GreenUp, 520, 400, func(c *ui.Context) {
		ui.Box(c).Width(400).Children(func() {
			if side, q, ok := QuickTradeButtons(c, &p.qty, p.opts).Traded(); ok {
				p.trades, p.sizes = append(p.trades, side), append(p.sizes, q)
			}
		})
	})
}

func TestQuickTradeButtonsPresetsAndTrade(t *testing.T) {
	p := &quickTradeProbe{opts: QuickTradeButtonsOptions{Symbol: "MUJC"}}
	tt := p.tester()
	if p.qty != 1 {
		t.Fatalf("default preset %v", p.qty)
	}
	if err := tt.Click("10"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Sell 10"); err != nil {
		t.Fatal(err)
	}
	if len(p.trades) != 1 || p.trades[0] != SideSell || p.sizes[0] != 10 {
		t.Fatalf("trades %v %v", p.trades, p.sizes)
	}
}

func TestQuickTradeButtonsConfirmFocusesCancel(t *testing.T) {
	p := &quickTradeProbe{qty: 100, opts: QuickTradeButtonsOptions{Symbol: "MUJC", Confirm: true}}
	tt := p.tester()
	_ = tt.Click("Buy 100")
	if len(p.trades) != 0 || !tt.HasText("Buy 100 MUJC?") || !tt.Focused("Cancel") {
		t.Fatalf("confirm: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if len(p.trades) != 0 || tt.HasText("Buy 100 MUJC?") {
		t.Fatal("Enter traded instead of canceling")
	}
	_ = tt.Click("Buy 100")
	if err := tt.Click("Buy"); err != nil {
		t.Fatal(err)
	}
	if len(p.trades) != 1 || p.trades[0] != SideBuy || p.sizes[0] != 100 {
		t.Fatalf("trades %v", p.trades)
	}
}

func TestQuickTradeButtonsDisabledAndBadQuantity(t *testing.T) {
	p := &quickTradeProbe{opts: QuickTradeButtonsOptions{Symbol: "MUJC", Disabled: true}}
	tt := p.tester()
	_ = tt.Click("Buy 1")
	_ = tt.Click("10")
	if len(p.trades) != 0 || p.qty != 1 {
		t.Fatal("disabled buttons acted")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("quantity outside presets accepted")
		}
	}()
	(&quickTradeProbe{qty: 7, opts: QuickTradeButtonsOptions{Symbol: "MUJC"}}).tester()
}

func TestQuickTradeButtonsSnapshot(t *testing.T) {
	qty := 10.0
	snapshot(t, "quick_trade_buttons", 380, 130, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			QuickTradeButtons(c, &qty, QuickTradeButtonsOptions{Symbol: "MUJC"})
		})
	})
}
