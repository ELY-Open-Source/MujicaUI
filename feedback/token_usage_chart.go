package feedback

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

// TokenUsage is the tokens one session or step used.
type TokenUsage struct {
	Label         string
	Input, Output int
}

// TokenUsageChartOptions configure a TokenUsageChart.
type TokenUsageChartOptions struct {
	// Frame names, sizes and states the chart.
	Frame chart.ChartFrame
	// Line draws lines instead of stacked bars.
	Line bool
	// Version changes when usage is edited in place.
	Version int
}

type tokenUsageSeries struct {
	src     *TokenUsage
	n       int
	version int
	labels  []string
	series  []chart.ChartSeries
}

// TokenUsageChart shows input and output tokens per step.
func TokenUsageChart(c *ui.Context, usage []TokenUsage, opts TokenUsageChartOptions) chart.ChartResult {
	if opts.Frame.Label == "" {
		opts.Frame.Label = core.Msg(c, "token_usage_chart.label")
	}
	holder := ui.Column(c).AlignItems(ui.Stretch)
	cache := ui.Local(holder, "series", func() *tokenUsageSeries { return &tokenUsageSeries{} })
	s := *cache
	var first *TokenUsage
	if len(usage) > 0 {
		first = &usage[0]
	}
	if s.series == nil || s.src != first || s.n != len(usage) || s.version != opts.Version {
		step := max(1, (len(usage)+255)/256)
		n := (len(usage) + step - 1) / step
		in, out := make([]float64, n), make([]float64, n)
		s.labels = make([]string, n)
		for i, u := range usage {
			if u.Input < 0 || u.Output < 0 {
				panic("mujica: TokenUsage counts must not be negative")
			}
			at := i / step
			in[at] += float64(u.Input)
			out[at] += float64(u.Output)
			if i%step == 0 {
				s.labels[at] = u.Label
			}
			if step > 1 && (i%step == step-1 || i == len(usage)-1) {
				s.labels[at] += " … " + u.Label
			}
		}
		s.series = []chart.ChartSeries{{Name: core.Msg(c, "token_usage_chart.input"), Values: in}, {Name: core.Msg(c, "token_usage_chart.output"), Values: out}}
		s.src, s.n, s.version = first, len(usage), opts.Version
	}
	s.series[0].Name = core.Msg(c, "token_usage_chart.input")
	s.series[1].Name = core.Msg(c, "token_usage_chart.output")
	format := func(v float64) string { return strconv.FormatFloat(v, 'f', 0, 64) }
	var res chart.ChartResult
	holder.Children(func() {
		if opts.Line {
			res = chart.LineChart(c, s.series, chart.LineChartOptions{Frame: opts.Frame, Format: format})
			return
		}
		res = chart.BarChart(c, s.labels, s.series, chart.BarChartOptions{Frame: opts.Frame, Stacked: true, Format: format})
	})
	res.Element = holder
	return res
}
