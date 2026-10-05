package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTokenCounterLevels(t *testing.T) {
	for _, tc := range []struct {
		used  int
		text  string
		label string
	}{
		{1234, "≈ 1,234 / 8,000", "About 1,234 of 8,000 tokens"},
		{7500, "≈ 7,500 / 8,000", "About 7,500 of 8,000 tokens"},
		{9100, "≈ 9,100 / 8,000", "About 9,100 of 8,000 tokens, 1,100 over the limit"},
	} {
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			TokenCounter(c, tc.used, 8000)
		}, 300, 60)
		if !tt.HasText(tc.text) {
			t.Errorf("%d: texts %q", tc.used, tt.Texts())
		}
		if _, ok := tt.Find(tc.label); !ok {
			t.Errorf("%d: no %q", tc.used, tc.label)
		}
	}
	if tokenThousands(1234567) != "1,234,567" || tokenThousands(12) != "12" {
		t.Fatal("thousands separators")
	}
}

func TestTokenCounterBadLimit(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a zero limit")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		TokenCounter(c, 1, 0)
	}, 100, 100)
}

func TestTokenCounterSnapshot(t *testing.T) {
	snapshot(t, "token_counter", 240, 110, func(c *ui.Context) {
		ui.Column(c).Padding(12).Gap(8).Children(func() {
			TokenCounter(c, 1234, 8000)
			TokenCounter(c, 7500, 8000)
			TokenCounter(c, 9100, 8000)
		})
	})
}
