package devtools

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
)

func TestStatCardValueChangeTrend(t *testing.T) {
	opts := StatCardOptions{Label: "Requests", Change: 0.124, HasChange: true, Trend: []float64{3, 5, 4, 8, 9}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		StatCard(c, "12,480", opts)
	}, 400, 200)
	for _, s := range []string{"Requests", "12,480", "+12.4%"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	opts.Change = -0.05
	tt.Frame()
	if !tt.HasText("-5.0%") {
		t.Fatalf("fall: %q", tt.Texts())
	}
	opts.HasChange = false
	tt.Frame()
	if tt.HasText("-5.0%") {
		t.Fatal("change shown without HasChange")
	}
}

func TestStatCardLoadingErrorSameSize(t *testing.T) {
	status, retries := data.DataReady, 0
	var h []float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			r := StatCard(c, "42 ms", StatCardOptions{Label: "Latency", Status: status, Error: "timeout"})
			h = append(h, r.Element.Bounds().H)
			if r.Changed() {
				retries++
			}
		})
	}, 400, 300)
	tt.Frame()
	ready := h[len(h)-1]
	status = data.DataLoading
	tt.Frame()
	tt.Frame()
	if !tt.HasText("Loading") || h[len(h)-1] != ready {
		t.Fatalf("loading: %q height %v vs %v", tt.Texts(), h[len(h)-1], ready)
	}
	status = data.DataFailed
	tt.Frame()
	tt.Frame()
	if !tt.HasText("timeout") || h[len(h)-1] != ready {
		t.Fatalf("failed: %q height %v vs %v", tt.Texts(), h[len(h)-1], ready)
	}
	_ = tt.Click("Retry")
	if retries != 1 {
		t.Fatalf("retries %d", retries)
	}
}

func TestStatCardSnapshot(t *testing.T) {
	snapshot(t, "stat_card", 720, 150, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(12).Children(func() {
			StatCard(c, "12,480", StatCardOptions{Label: "Requests", Change: 0.124, HasChange: true, Trend: []float64{3, 5, 4, 8, 6, 9}})
			StatCard(c, "182 ms", StatCardOptions{Label: "p95 latency", Change: 0.08, HasChange: true, LowerIsBetter: true, Trend: []float64{120, 140, 130, 182}})
			StatCard(c, "", StatCardOptions{Label: "Errors", Status: data.DataFailed, Error: "Source unreachable"})
		})
	})
}

func TestStatCardDisabledRetryAndKeyboardTrend(t *testing.T) {
	status := data.DataReady
	disabled := false
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			if StatCard(c, "42", StatCardOptions{Label: "CPU", Trend: []float64{1, 4, 2}, Status: status, Error: "offline"}).Changed() {
				retries++
			}
		})
	}, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	announced := false
	for _, value := range tt.Announcements() {
		announced = announced || value == "2"
	}
	if !announced || !tt.HasText("2") {
		t.Fatalf("trend boundary readout: %q", tt.Texts())
	}
	status, disabled = data.DataFailed, true
	tt.Frame()
	_ = tt.Click("Retry")
	tt.Key(0, ui.KeyEnter)
	if retries != 0 {
		t.Fatal("disabled retry")
	}
}

func TestStatCardTextScaleKeepsCardHeight(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); StatCard(c, "42", StatCardOptions{Label: "CPU"}) }, 600, 350)
	tt.SetPreferences(ui.Preferences{TextScale: 2, ReduceMotion: true})
	bounds, ok := tt.Find("CPU")
	if !ok || bounds.H < 232 {
		t.Fatalf("scaled card height %v", bounds.H)
	}
}
