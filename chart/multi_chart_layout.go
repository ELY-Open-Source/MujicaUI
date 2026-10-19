package chart

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ChartCell holds one chart's symbol, interval, and window.
type ChartCell struct {
	Symbol   string
	Interval Interval
	View     MarketView
}

// MultiChartState preserves six cells across layout changes.
type MultiChartState struct {
	// Count is 1, 2, 4 or 6; zero is 1.
	Count int
	// Cells preserve their settings while hidden.
	Cells  [6]ChartCell
	Active int
}

// MultiChartLayoutOptions configure a MultiChartLayout.
type MultiChartLayoutOptions struct {
	// Height of the grid; zero is 480.
	Height float32
}

// MultiChartLayout builds 1, 2, 4, or 6 independent charts.
func MultiChartLayout(c *ui.Context, state *MultiChartState, opts MultiChartLayoutOptions, cell func(i int, cell *ChartCell)) core.ChoiceResult {
	if state == nil || cell == nil {
		panic("mujica: MultiChartLayout needs a state and a cell builder")
	}
	root := ui.Column(c).Gap(core.Density(c).Unit() * 2).FillWidth()
	if root.IsDisabled() {
		copy := *state
		state = &copy
	}
	if opts.Height < 0 {
		panic("mujica: MultiChartLayout Height must be non-negative")
	}
	if state.Count == 0 {
		state.Count = 1
	}
	cols := map[int]int{1: 1, 2: 2, 4: 2, 6: 3}[state.Count]
	if cols == 0 {
		panic(fmt.Sprintf("mujica: MultiChartLayout Count %d must be 1, 2, 4 or 6", state.Count))
	}
	if state.Active < 0 || state.Active >= state.Count {
		panic(fmt.Sprintf("mujica: MultiChartLayout Active %d outside %d cells", state.Active, state.Count))
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	h := opts.Height
	if h == 0 {
		h = 480
	}
	changed := false
	root.Children(func() {
		sel := []int{state.Count}
		var items []input.ToggleItem[int]
		for _, n := range []int{1, 2, 4, 6} {
			items = append(items, input.ToggleItem[int]{Value: n, Label: fmt.Sprint(n), Name: fmt.Sprintf(core.Msg(c, "multi_chart_layout.count"), n)})
		}
		ui.Row(c).Children(func() {
			if input.ToggleGroup(c, &sel, items, input.ToggleGroupOptions{Label: core.Msg(c, "multi_chart_layout.label")}).Changed() {
				state.Count, changed = sel[0], true
				state.Active = min(state.Active, state.Count-1)
				slog.Debug("mujica: multi_chart_layout count", "count", state.Count)
			}
		})
		cols = map[int]int{1: 1, 2: 2, 4: 2, 6: 3}[state.Count]
		rows := (state.Count + cols - 1) / cols
		grid := ui.Box(c).Grid().Columns(cols).Gap(u * 2).Height(h).FillWidth()
		grid.Children(func() {
			for i := range state.Count {
				box := ui.Column(c.Key(i)).MinWidth(0).MinHeight(0).Height((h - u*2*float32(rows-1)) / float32(rows)).Padding(u * 2).
					Background(k.Surface).Radius(theme.CardRadius).Clip().Role(ui.RoleGroup).
					Label(fmt.Sprintf(core.Msg(c, "multi_chart_layout.cell"), i+1, state.Cells[i].Symbol))
				if i == state.Active && state.Count > 1 {
					box.Border(2, k.AccentText)
					box.DrawOver(func(p *ui.Painter, r ui.Rect) { internal.Frame(p, r, 3, 5, k.Ornament) })
				} else {
					box.Border(theme.BorderWidth, k.Border)
				}
				if !box.IsDisabled() && (box.FocusWithin() || box.Clicked()) {
					if state.Active != i {
						state.Active, changed = i, true
					}
				}
				box.Children(func() { cell(i, &state.Cells[i]) })
			}
		})
	})
	return core.NewChoiceResult(root, changed)
}
