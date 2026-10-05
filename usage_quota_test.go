package mujica

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var quotaItems = []QuotaItem{{Label: "API calls", Used: 4200, Limit: 10000}, {Label: "Storage", Used: 9, Limit: 10, Unit: "GB"}, {Label: "Seats", Used: 12, Limit: 10}}

func TestUsageQuotaLevelsAndReset(t *testing.T) {
	reset := time.Duration(5 * time.Hour)
	tt := accountTester(500, 400, func(c *ui.Context) {
		UsageQuota(c, quotaItems, UsageQuotaOptions{ResetAt: c.Now().Add(reset)})
	})
	for _, want := range []string{"4,200 of 10,000 ", "9 of 10 GB", "Near limit", "Limit reached", "Resets in 5 h"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	reset = 72 * time.Hour
	tt.Frame()
	if tt.HasText("Resets in 5 h") {
		t.Fatalf("far reset: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("zero limit accepted")
		}
	}()
	accountTester(100, 100, func(c *ui.Context) { UsageQuota(c, []QuotaItem{{Label: "x"}}, UsageQuotaOptions{}) })
}

func TestUsageQuotaSnapshot(t *testing.T) {
	snapshot(t, "usage_quota", 440, 230, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			UsageQuota(c, quotaItems, UsageQuotaOptions{ResetAt: c.Now().Add(5*time.Hour + time.Minute)})
		})
	})
}

func TestUsageQuotaRejectsNonFiniteValues(t *testing.T) {
	for _, tc := range []struct{ used, limit, warn float64 }{
		{1, 10, math.NaN()}, {math.Inf(1), 10, 0}, {1, math.Inf(1), 0}, {math.NaN(), 10, 0}, {1, math.NaN(), 0},
	} {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("non-finite quota accepted")
				}
			}()
			accountTester(500, 300, func(c *ui.Context) {
				UsageQuota(c, []QuotaItem{{Label: "Calls", Used: tc.used, Limit: tc.limit}}, UsageQuotaOptions{Warn: tc.warn})
			})
		})
	}
}
