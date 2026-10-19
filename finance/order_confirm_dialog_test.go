package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/overlay"
)

func orderConfirmTester(open *bool, opts OrderConfirmDialogOptions, res *overlay.AlertDialogResult) *ui.Tester {
	return tradingTester(core.GreenUp, 600, 500, func(c *ui.Context) {
		if r := OrderConfirmDialog(c, open, opts); r.Confirmed || r.Canceled {
			*res = r
		}
	})
}

var orderConfirmLines = []OrderSummaryLine{{"Symbol", "MUJC"}, {"Side", "Buy"}, {"Quantity", "100"}, {"Est. amount", "10,020.00"}}

func TestOrderConfirmDialogFocusesCancel(t *testing.T) {
	open := true
	var res overlay.AlertDialogResult
	tt := orderConfirmTester(&open, OrderConfirmDialogOptions{Lines: orderConfirmLines, Risk: "Leveraged products can lose more than the deposit."}, &res)
	if !tt.Focused("Cancel") || !tt.HasText("10,020.00") || !tt.HasText("Leveraged products can lose more than the deposit.") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if !res.Canceled || open {
		t.Fatalf("Enter on Cancel: %+v open %v", res, open)
	}
	open, res = true, overlay.AlertDialogResult{}
	tt.Frame()
	if err := tt.Click("Place order"); err != nil {
		t.Fatal(err)
	}
	if !res.Confirmed || open {
		t.Fatalf("confirm: %+v", res)
	}
}

func TestOrderConfirmDialogAcknowledgeAndEscape(t *testing.T) {
	open := true
	var res overlay.AlertDialogResult
	tt := orderConfirmTester(&open, OrderConfirmDialogOptions{Lines: orderConfirmLines, Acknowledge: "I understand the risk"}, &res)
	_ = tt.Click("Place order")
	if res.Confirmed || !open {
		t.Fatal("confirmed before acknowledging")
	}
	_ = tt.Click("I understand the risk")
	_ = tt.Click("Place order")
	if !res.Confirmed {
		t.Fatal("acknowledged confirm refused")
	}
	open, res = true, overlay.AlertDialogResult{}
	tt.Frame()
	_ = tt.Click("Place order")
	if res.Confirmed {
		t.Fatal("acknowledgement survived reopening")
	}
	tt.Key(0, ui.KeyEscape)
	if !res.Canceled || open {
		t.Fatalf("escape: %+v", res)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("empty summary accepted")
		}
	}()
	open = true
	orderConfirmTester(&open, OrderConfirmDialogOptions{}, &res)
}

func TestOrderConfirmDialogSnapshot(t *testing.T) {
	open := true
	snapshot(t, "order_confirm_dialog", 560, 480, func(c *ui.Context) {
		OrderConfirmDialog(c, &open, OrderConfirmDialogOptions{Lines: orderConfirmLines,
			Risk: "Leveraged products can lose more than the deposit.", Acknowledge: "I understand the risk"})
	})
}
