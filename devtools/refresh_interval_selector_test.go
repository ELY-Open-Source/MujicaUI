package devtools

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func refreshTester(interval *time.Duration, opts *RefreshIntervalSelectorOptions, refreshes, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := RefreshIntervalSelector(c, interval, *opts)
		if r.Refresh {
			*refreshes++
		}
		if r.Changed {
			*changes++
		}
	}, 500, 300)
}

func TestRefreshIntervalSelectorManualAndAuto(t *testing.T) {
	interval := time.Duration(0)
	opts := RefreshIntervalSelectorOptions{Intervals: []time.Duration{time.Second, time.Minute}}
	refreshes, changes := 0, 0
	tt := refreshTester(&interval, &opts, &refreshes, &changes)
	_ = tt.Click("Refresh")
	if refreshes != 1 {
		t.Fatalf("manual refreshes %d", refreshes)
	}
	tt.Advance(1100 * time.Millisecond)
	if refreshes != 1 {
		t.Fatal("refreshed while off")
	}
	_ = tt.Click("Refresh interval")
	_ = tt.Click("Every 1 s")
	if interval != time.Second || changes != 1 {
		t.Fatalf("interval %v changes %d", interval, changes)
	}
	start := tt.Now()
	if !overlayUntil(tt, func() bool { return refreshes == 2 }) {
		t.Fatal("no automatic refresh")
	}
	if tt.Now().Sub(start) < 900*time.Millisecond {
		t.Fatal("refreshed before the interval")
	}
}

func TestRefreshIntervalSelectorRefreshingAndDisabled(t *testing.T) {
	interval := time.Second
	opts := RefreshIntervalSelectorOptions{Intervals: []time.Duration{time.Second}, Refreshing: true}
	refreshes, changes := 0, 0
	tt := refreshTester(&interval, &opts, &refreshes, &changes)
	_ = tt.Click("Refresh")
	tt.Advance(1100 * time.Millisecond)
	if refreshes != 0 || !tt.HasText("Refresh") {
		t.Fatalf("refreshing: %d", refreshes)
	}
	opts = RefreshIntervalSelectorOptions{Intervals: []time.Duration{time.Second}, Disabled: true}
	tt.Frame()
	tt.Advance(1100 * time.Millisecond)
	_ = tt.Click("Refresh")
	if refreshes != 0 {
		t.Fatalf("disabled refreshed %d", refreshes)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unoffered interval did not panic")
		}
	}()
	bad := 7 * time.Second
	refreshTester(&bad, &RefreshIntervalSelectorOptions{}, &refreshes, &changes)
}

func TestRefreshIntervalSelectorSnapshot(t *testing.T) {
	snapshot(t, "refresh_interval_selector", 420, 80, func(c *ui.Context) {
		interval := 30 * time.Second
		ui.Column(c).Padding(16).Children(func() {
			RefreshIntervalSelector(c, &interval, RefreshIntervalSelectorOptions{})
		})
	})
}

func TestRefreshIntervalSelectorKeyboardAncestorDisable(t *testing.T) {
	interval := time.Duration(0)
	disabled := false
	changes, refreshes := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			r := RefreshIntervalSelector(c, &interval, RefreshIntervalSelectorOptions{})
			if r.Changed {
				changes++
			}
			if r.Refresh {
				refreshes++
			}
		})
	}, 500, 400)
	_ = tt.Click("Refresh interval")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if interval != 5*time.Second || changes != 1 {
		t.Fatal("keyboard interval")
	}
	_ = tt.Click("Refresh interval")
	disabled = true
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if interval != 5*time.Second || changes != 1 || refreshes != 0 || tt.HasText("Every 10 s") {
		t.Fatal("disabled popup remained active")
	}
}

func TestRefreshIntervalSelectorRejectsDuplicateIntervals(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate intervals accepted")
		}
	}()
	interval := time.Duration(0)
	refreshes, changes := 0, 0
	refreshTester(&interval, &RefreshIntervalSelectorOptions{Intervals: []time.Duration{time.Second, time.Second}}, &refreshes, &changes)
}
