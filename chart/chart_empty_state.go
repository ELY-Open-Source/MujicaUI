package chart

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
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

// ChartEmptyState preserves plot size and reports retry clicks.
func ChartEmptyState(c *ui.Context, status data.DataStatus, opts ChartEmptyStateOptions) core.ChoiceResult {
	if opts.Height < 0 {
		panic("mujica: ChartEmptyState Height must be ≥ 0")
	}
	h := opts.Height
	if h == 0 {
		h = 240
	}
	var res core.ChoiceResult
	res.Element = ui.Column(c).Height(h).FillWidth().MinWidth(0)
	res.Element.Children(func() {
		res = core.NewChoiceResult(res.Element, data.DataPlaceholder(c, status, opts.Error, opts.Empty, core.Msg(c, "chart_empty_state.empty")) && !res.Element.IsDisabled())
	})
	return res
}
