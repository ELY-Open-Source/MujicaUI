package mujica

import (
	"strconv"

	"github.com/egoist/mygo/ui"
)

// TokenUsage is the tokens one session or step used.
type TokenUsage struct {
	Label         string
	Input, Output int
}

// TokenUsageChartOptions configure a TokenUsageChart.
type TokenUsageChartOptions struct {
	// Frame names, sizes and states the chart; an empty Label uses the built-in name.
	Frame ChartFrame
	// Line draws input and output as lines instead of stacked bars.
	Line bool
}

type tokenUsageSeries struct {
	src    *TokenUsage
	n      int
	labels []string
	series []ChartSeries
}

// TokenUsageChart shows input and output tokens per session or step, as stacked bars or lines.
func TokenUsageChart(c *ui.Context, usage []TokenUsage, opts TokenUsageChartOptions) ChartResult {
	if opts.Frame.Label == "" {
		opts.Frame.Label = msg(c, "token_usage_chart.label")
	}
	holder := ui.Column(c).AlignItems(ui.Stretch)
	cache := ui.Local(holder, "series", func() *tokenUsageSeries { return &tokenUsageSeries{} })
	s := *cache
	var first *TokenUsage
	if len(usage) > 0 {
		first = &usage[0]
	}
	if s.series == nil || s.src != first || s.n != len(usage) {
		in, out := make([]float64, len(usage)), make([]float64, len(usage))
		s.labels = make([]string, len(usage))
		for i, u := range usage {
			if u.Input < 0 || u.Output < 0 {
				panic("mujica: TokenUsage counts must not be negative")
			}
			in[i], out[i], s.labels[i] = float64(u.Input), float64(u.Output), u.Label
		}
		s.series = []ChartSeries{{Name: msg(c, "token_usage_chart.input"), Values: in}, {Name: msg(c, "token_usage_chart.output"), Values: out}}
		s.src, s.n = first, len(usage)
	}
	format := func(v float64) string { return strconv.FormatFloat(v, 'f', 0, 64) }
	var res ChartResult
	holder.Children(func() {
		if opts.Line {
			res = LineChart(c, s.series, LineChartOptions{Frame: opts.Frame, Format: format})
			return
		}
		res = BarChart(c, s.labels, s.series, BarChartOptions{Frame: opts.Frame, Stacked: true, Format: format})
	})
	res.Element = holder
	return res
}
