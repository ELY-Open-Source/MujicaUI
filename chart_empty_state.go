package mujica

import (
	"github.com/egoist/mygo/ui"
)

// ChartEmptyStateOptions configure a ChartEmptyState.
type ChartEmptyStateOptions struct {
	// Height keeps the chart's size; 0 means 240.
	Height float32
	// Error describes a DataFailed status; empty uses the built-in text.
	Error string
	// Empty replaces the built-in "No data" text.
	Empty string
}

// ChartEmptyState holds a chart's place while it loads, fails or has no data; Changed reports Retry.
func ChartEmptyState(c *ui.Context, status DataStatus, opts ChartEmptyStateOptions) ChoiceResult {
	if opts.Height < 0 {
		panic("mujica: ChartEmptyState Height must be ≥ 0")
	}
	h := opts.Height
	if h == 0 {
		h = 240
	}
	var res ChoiceResult
	res.Element = ui.Column(c).Height(h).FillWidth().MinWidth(0)
	res.Element.Children(func() {
		res.changed = dataPlaceholder(c, status, opts.Error, opts.Empty, msg(c, "chart.empty")) && !res.Element.IsDisabled()
	})
	return res
}
