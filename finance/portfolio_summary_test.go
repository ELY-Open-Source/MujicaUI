package finance

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestPortfolioSummaryAccountsAndTotals(t *testing.T) {
	account := ""
	disabled := false
	changed := false
	rows := []PortfolioAccount{{ID: "a", Name: "Brokerage", Total: 120000, DayPnL: 1250, TotalReturn: 20000, Cost: 100000}, {ID: "b", Name: "Retirement", Total: 80000, DayPnL: -300, TotalReturn: 5000, Cost: 75000}}
	tt := chatTester(960, 300, func(c *ui.Context) {
		r := PortfolioSummary(c, &account, rows, PortfolioSummaryOptions{Disabled: disabled})
		changed = changed || r.Changed()
	})
	if !tt.HasText("$200,000.00") || !tt.HasText("+$950.00") || !tt.HasText("+$25,000.00") {
		t.Fatalf("totals %q", tt.Texts())
	}
	_ = tt.Click("Retirement")
	tt.Frame()
	if account != "b" || !changed || !tt.HasText("$80,000.00") {
		t.Fatalf("account %s texts %q", account, tt.Texts())
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Brokerage")
	if account != "b" {
		t.Fatal("disabled changed")
	}
	if !math.IsNaN(portfolioPercent(2, 0)) {
		t.Fatal("zero base")
	}
}
func TestPortfolioSummarySnapshot(t *testing.T) {
	acct := ""
	snapshot(t, "portfolio_summary", 760, 260, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			PortfolioSummary(c, &acct, []PortfolioAccount{{ID: "a", Name: "Brokerage", Total: 120000, DayPnL: 1250, TotalReturn: 20000, Cost: 100000},
				{ID: "b", Name: "Retirement", Total: 80000, DayPnL: -300, TotalReturn: 5000, Cost: 75000}}, PortfolioSummaryOptions{})

		})
	})
}
