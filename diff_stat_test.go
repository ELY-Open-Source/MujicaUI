package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDiffStatBlocks(t *testing.T) {
	for _, tc := range []struct{ a, d, n, plus, minus int }{
		{0, 0, 5, 0, 0},
		{10, 0, 5, 5, 0},
		{0, 3, 5, 0, 5},
		{6, 4, 5, 3, 2},
		{1000, 1, 5, 4, 1},
		{1, 1000, 5, 1, 4},
		{1, 1, 1, 1, 0},
	} {
		p, m := diffStatBlocks(tc.a, tc.d, tc.n)
		if p != tc.plus || m != tc.minus {
			t.Errorf("blocks(%d,%d,%d) = %d,%d want %d,%d", tc.a, tc.d, tc.n, p, m, tc.plus, tc.minus)
		}
	}
}

func TestDiffStatTextAndName(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		DiffStat(c, 12, 3, DiffStatOptions{})
	}, 300, 60)
	if !tt.HasText("+12") || !tt.HasText("−3") {
		t.Fatalf("counts: %q", tt.Texts())
	}
	if _, ok := tt.Find("12 additions, 3 deletions"); !ok {
		t.Fatal("not named for assistive technology")
	}
	for _, tc := range []struct{ a, d, n int }{{-1, 0, 0}, {0, -1, 0}, {1, 1, -2}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%v did not panic", tc)
				}
			}()
			ui.Render(func(c *ui.Context) {
				Use(c, Settings{})
				DiffStat(c, tc.a, tc.d, DiffStatOptions{Blocks: tc.n})
			}, 100, 50, 1)
		}()
	}
}

func TestDiffStatSnapshot(t *testing.T) {
	snapshot(t, "diff_stat", 360, 120, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(10).Children(func() {
			DiffStat(c, 128, 32, DiffStatOptions{})
			DiffStat(c, 4, 0, DiffStatOptions{})
			DiffStat(c, 0, 0, DiffStatOptions{Blocks: 8})
		})
	})
}
