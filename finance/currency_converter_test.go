package finance

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestCurrencyConverterPrecisionAmountAndNoRate(t *testing.T) {
	s := CurrencyConverterState{Amount: 1.23, From: "USD", To: "JPY"}
	var last CurrencyConverterView
	rate := 149.37
	disabled := false
	currencies := []Currency{{"USD", "Dollar", 2}, {"JPY", "Yen", 0}}
	tt := chatTester(960, 300, func(c *ui.Context) {
		last = CurrencyConverter(c, &s, currencies, CurrencyConverterOptions{Rate: rate, Disabled: disabled})
	})
	if last.Result != 184 || !tt.HasText("184 JPY") {
		t.Fatalf("result %v text %q", last.Result, tt.Texts())
	}
	rate = 0
	tt.Frame()
	if !math.IsNaN(last.Result) || !tt.HasText("Rate unavailable") {
		t.Fatal("missing rate")
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Swap currencies")
	if s.From != "USD" {
		t.Fatal("disabled swapped")
	}
	marketsPanic(t, func(c *ui.Context) { CurrencyConverter(c, &s, currencies, CurrencyConverterOptions{Rate: -1}) })
}
func TestCurrencyConverterSnapshot(t *testing.T) {
	cs := CurrencyConverterState{Amount: 100, From: "USD", To: "JPY"}
	snapshot(t, "currency_converter", 760, 260, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			CurrencyConverter(c, &cs, []Currency{{"USD", "US Dollar", 2}, {"JPY", "Yen", 0}, {"EUR", "Euro", 2}}, CurrencyConverterOptions{Rate: 149.37, RateTime: mktNow})

		})
	})
}
