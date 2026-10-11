package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestCurrencyFormatAndParse(t *testing.T) {
	for _, tc := range []struct {
		minor int64
		dec   int
		group bool
		want  string
	}{
		{123450, 2, true, "1,234.50"}, {123450, 2, false, "1234.50"}, {-5, 2, true, "-0.05"},
		{1234567, 0, true, "1,234,567"}, {0, 3, true, "0.000"}, {100000000, 2, true, "1,000,000.00"},
	} {
		if got := CurrencyFormat(tc.minor, tc.dec, tc.group); got != tc.want {
			t.Errorf("CurrencyFormat(%d, %d) = %q, want %q", tc.minor, tc.dec, got, tc.want)
		}
	}
	for text, want := range map[string]int64{"1,234.5": 123450, "0.07": 7, "-12": -1200, " 3 ": 300, "1.": 100} {
		if v, why := currencyInputParse(text, 2); v != want || why != "" {
			t.Errorf("parse %q = %d %q, want %d", text, v, why, want)
		}
	}
	for text, want := range map[string]string{"1.234": "currency_input.precision", "abc": "currency_input.invalid", "": "currency_input.invalid", "1.2.3": "currency_input.invalid", "99999999999999999999": "currency_input.invalid"} {
		if _, why := currencyInputParse(text, 2); why != want {
			t.Errorf("parse %q: %q, want %q", text, why, want)
		}
	}
}

// fieldText copies the focused field's text through the clipboard.
func fieldText(tt *ui.Tester) string {
	tt.Command("selectAll")
	tt.Command("copy")
	return tt.Clipboard()
}

func TestCurrencyInputTypingGroupingAndErrors(t *testing.T) {
	v := int64(123450)
	dis := false
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if CurrencyInput(c, &v, CurrencyInputOptions{Symbol: "$", Label: "Price", Disabled: dis}).Changed() {
			n++
		}
		Button(c, "Elsewhere", ButtonOptions{})
	}, 500, 200)
	if !tt.HasText("$") {
		t.Fatalf("no symbol: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	if got := fieldText(tt); got != "1234.50" {
		t.Fatalf("focused text %q, want ungrouped", got)
	}
	tt.Command("selectAll")
	tt.Type("2500.255")
	if v != 123450 || !tt.HasText("Use at most 2 decimal places") {
		t.Fatalf("too precise: %d %q", v, tt.Texts())
	}
	tt.Command("selectAll")
	tt.Type("2500.25")
	if v != 250025 || n != 1 || tt.HasText("decimal places") {
		t.Fatalf("typed: %d (%d)", v, n)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyTab)
	if got := fieldText(tt); got != "2500.25" {
		t.Fatalf("refocused text %q", got)
	}
	dis = true
	tt.Frame()
	_ = tt.Click("Price")
	tt.Type("9")
	if v != 250025 {
		t.Fatalf("disabled edit: %d", v)
	}
}

func TestCurrencyInputNoFractionAndRejects(t *testing.T) {
	v := int64(1500)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CurrencyInput(c, &v, CurrencyInputOptions{Symbol: "¥", NoFraction: true, Label: "Amount"})
	}, 300, 100)
	tt.Key(0, ui.KeyTab)
	if got := fieldText(tt); got != "1500" {
		t.Fatalf("yen %q", got)
	}
	tt.Command("selectAll")
	tt.Type("12.5")
	if v != 1500 || !tt.HasText("Use at most 0 decimal places") {
		t.Fatalf("fraction accepted: %d", v)
	}
	for _, o := range []CurrencyInputOptions{{Decimals: -1}, {Decimals: 7}, {Decimals: 2, NoFraction: true}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: no panic", o)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); CurrencyInput(c, &v, o) }, 200, 100)
		}()
	}
}

func TestCurrencyInputSnapshot(t *testing.T) {
	a, b, d := int64(1234550), int64(98000), int64(0)
	snapshot(t, "currency_input", 520, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(12).Children(func() {
			CurrencyInput(c, &a, CurrencyInputOptions{Symbol: "$"}).Element.Width(150)
			CurrencyInput(c, &b, CurrencyInputOptions{Symbol: "¥", NoFraction: true}).Element.Width(150)
			CurrencyInput(c, &d, CurrencyInputOptions{Symbol: "€", Disabled: true}).Element.Width(150)
		})
	})
}
