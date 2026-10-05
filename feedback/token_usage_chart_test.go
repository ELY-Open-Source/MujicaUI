package feedback

import (
	"strconv"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func tokenUsageSample() []TokenUsage {
	return []TokenUsage{{Label: "Plan", Input: 1200, Output: 300}, {Label: "Search", Input: 4800, Output: 650}, {Label: "Write", Input: 2100, Output: 2900}}
}

func TestTokenUsageChartBarsKeyboard(t *testing.T) {
	usage := tokenUsageSample()
	h := &chartHarness{build: func(c *ui.Context) chart.ChartResult { return TokenUsageChart(c, usage, TokenUsageChartOptions{}) }}
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
	h := &chartHarness{build: func(c *ui.Context) chart.ChartResult {
		return TokenUsageChart(c, usage, TokenUsageChartOptions{Line: true, Frame: chart.ChartFrame{Label: "Tokens per step"}})
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
		core.Use(c, core.Settings{})
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

func TestTokenUsageChartLocaleRefresh(t *testing.T) {
	locale := core.EnUS
	usage := tokenUsageSample()
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: locale})
		TokenUsageChart(c, usage, TokenUsageChartOptions{})
	}, 520, 420)
	locale = core.ZhCN
	tt.Frame()
	if !tt.HasText("输入") || !tt.HasText("输出") {
		t.Fatalf("stale legend: %q", tt.Texts())
	}
}

func TestTokenUsageChartLargeDataBounded(t *testing.T) {
	usage := make([]TokenUsage, 100000)
	for i := range usage {
		usage[i] = TokenUsage{Label: strconv.Itoa(i), Input: 1, Output: 2}
	}
	var cache *tokenUsageSeries
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := TokenUsageChart(c, usage, TokenUsageChartOptions{})
		cache = *ui.Local(r.Element, "series", func() *tokenUsageSeries { panic("missing series cache") })
	}, 520, 420)
	if len(cache.labels) > 256 {
		t.Fatalf("built %d categories", len(cache.labels))
	}
	in, out := 0.0, 0.0
	for i := range cache.labels {
		in += cache.series[0].Values[i]
		out += cache.series[1].Values[i]
	}
	if in != 100000 || out != 200000 {
		t.Fatalf("aggregated totals %v %v", in, out)
	}
	first := &cache.series[0].Values[0]
	tt.Frame()
	if &cache.series[0].Values[0] != first {
		t.Fatal("unchanged data reaggregated")
	}
}
