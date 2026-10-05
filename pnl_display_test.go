package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestPnLDisplaySignsAndMask(t *testing.T) {
	hidden, n := false, 0
	tt := tradingTester(GreenUp, 520, 160, func(c *ui.Context) {
		if PnLDisplay(c, PnL{Realized: 1200.5, Unrealized: -340, Today: 0}, &hidden, PnLDisplayOptions{Currency: "$"}).Changed() {
			n++
		}
	})
	for _, want := range []string{"+$1,200.50", "−$340.00", "$0.00"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	if err := tt.Click("Hide amounts"); err != nil {
		t.Fatal(err)
	}
	if !hidden || n != 1 || tt.HasText("+$1,200.50") {
		t.Fatalf("mask: hidden %v changes %d texts %q", hidden, n, tt.Texts())
	}
	if _, ok := tt.Find("Realized, hidden"); !ok {
		t.Fatal("masked amount still spoken")
	}
	tt.Key(0, ui.KeySpace)
	if hidden || n != 2 {
		t.Fatalf("keyboard toggle: hidden %v", hidden)
	}
}

func TestPnLDisplaySnapshot(t *testing.T) {
	hidden := false
	snapshot(t, "pnl_display", 460, 100, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			PnLDisplay(c, PnL{Realized: 1200.5, Unrealized: -340.2, Today: 88}, &hidden, PnLDisplayOptions{Currency: "$"})
		})
	})
}
