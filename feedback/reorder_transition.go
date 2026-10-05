package feedback

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ReorderTransitionOptions configure a ReorderTransition.
type ReorderTransitionOptions struct {
	// Duration of a move; zero means StateDuration.
	Duration time.Duration
	// Gap between items in DIPs.
	Gap float32
}

// ReorderTransition slides keyed items to their new places.
func ReorderTransition[K comparable](c *ui.Context, keys []K, opts ReorderTransitionOptions, item func(key K)) *ui.Element {
	if opts.Duration < 0 || !(opts.Gap >= 0) || math.IsInf(float64(opts.Gap), 0) {
		panic("mujica: ReorderTransition Duration and Gap must not be negative")
	}
	d := opts.Duration
	if d == 0 {
		d = theme.StateDuration
	}
	seen := make(map[K]bool, len(keys))
	for _, k := range keys {
		if seen[k] {
			panic(fmt.Sprintf("mujica: ReorderTransition duplicate key %v", k))
		}
		seen[k] = true
	}
	move := ui.ElementTransition{Duration: core.Motion(c, d), Position: true, Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(opts.Gap)
	type state struct {
		keys    []K
		leaving map[K]time.Time
	}
	st := ui.Local(col, "reorder", func() state { return state{leaving: map[K]time.Time{}} })
	for _, k := range st.keys {
		if !seen[k] {
			st.leaving[k] = c.Now().Add(d)
		}
	}
	st.keys = append(st.keys[:0], keys...)
	col.Children(func() {
		for _, k := range keys {
			ui.Column(c).Key(k).AlignItems(ui.Stretch).Transition(move).Children(func() { item(k) })
		}
		for k, until := range st.leaving {
			if seen[k] || !c.Now().Before(until) {
				delete(st.leaving, k)
				continue
			}
			if c.Preferences().ReduceMotion {
				ui.Column(c).Key(k).Absolute().Invisible().Transition(ui.ElementTransition{Duration: d, Position: true})
				delete(st.leaving, k)
			}
		}
	})
	return col
}
