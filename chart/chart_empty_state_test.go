package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

func TestChartEmptyStateKeepsSizeAndRetries(t *testing.T) {
	status := data.DataReady
	retries, disabled := 0, false
	var heights []float32
	var el *ui.Element
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(400).Disabled(disabled).Children(func() {
			r := ChartEmptyState(c, status, ChartEmptyStateOptions{Height: 180, Empty: "No hymns this week", Error: "The choir did not answer"})
			el = r.Element
			if r.Changed() {
				retries++
			}
		})
	}, 440, 300)
	heights = append(heights, el.Bounds().H)
	if !tt.HasText("No hymns this week") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	status = data.DataLoading
	tt.Frame()
	heights = append(heights, el.Bounds().H)
	if !tt.HasText("Loading") {
		t.Fatalf("loading: %q", tt.Texts())
	}
	status = data.DataFailed
	tt.Frame()
	heights = append(heights, el.Bounds().H)
	if err := tt.Click("Retry"); err != nil || retries != 1 || !tt.HasText("The choir did not answer") {
		t.Fatalf("retry: %v %d", err, retries)
	}
	if heights[0] != 180 || heights[1] != 180 || heights[2] != 180 {
		t.Fatalf("heights %v", heights)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Retry")
	if retries != 1 {
		t.Fatal("disabled retry fired")
	}
}

func TestChartEmptyStateSnapshot(t *testing.T) {
	snapshot(t, "chart_empty_state", 640, 220, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			for _, st := range []data.DataStatus{data.DataReady, data.DataLoading, data.DataFailed} {
				ui.Column(c).Width(192).Border(1, core.Tokens(c).Border).Children(func() {
					ChartEmptyState(c, st, ChartEmptyStateOptions{Height: 180})
				})
			}
		})
	})
}
