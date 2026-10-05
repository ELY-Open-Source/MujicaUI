package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func tokenUsageSample() []TokenUsage {
	return []TokenUsage{{Label: "Plan", Input: 1200, Output: 300}, {Label: "Search", Input: 4800, Output: 650}, {Label: "Write", Input: 2100, Output: 2900}}
}

func TestTokenUsageChartBarsKeyboard(t *testing.T) {
	usage := tokenUsageSample()
	h := &chartHarness{build: func(c *ui.Context) ChartResult { return TokenUsageChart(c, usage, TokenUsageChartOptions{}) }}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Token usage")
	want := []string{"Plan · Input 1200", "Plan · Output 300", "Search · Input 4800"}
	for i, w := range want {
		tt.Key(0, ui.KeyRight)
		if got := chartReadout(tt); got != w {
			t.Fatalf("mark %d readout %q, want %q", i, got, w)
		}
	}
	if !tt.HasText("Input") || !tt.HasText("Output") {
		t.Fatal("legend missing")
	}
}

func TestTokenUsageChartLineEmptyInvalid(t *testing.T) {
	usage := tokenUsageSample()
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return TokenUsageChart(c, usage, TokenUsageChartOptions{Line: true, Frame: ChartFrame{Label: "Tokens per step"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Tokens per step")
	tt.Key(0, ui.KeyRight)
	if h.res.Point < 0 {
		t.Fatal("line chart has no keyboard path")
	}
	usage = nil
	tt.Frame()
	if !tt.HasText("No data") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("negative tokens did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		TokenUsageChart(c, []TokenUsage{{Label: "x", Input: -1}}, TokenUsageChartOptions{})
	}, 300, 200, 1)
}

func TestTokenUsageChartSnapshot(t *testing.T) {
	usage := tokenUsageSample()
	snapshot(t, "token_usage_chart", 520, 360, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			TokenUsageChart(c, usage, TokenUsageChartOptions{})
		})
	})
}
