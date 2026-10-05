package mujica

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// ReorderTransitionOptions configure a ReorderTransition.
type ReorderTransitionOptions struct {
	// Duration of a move; zero means 200 ms.
	Duration time.Duration
	// Gap between items in DIPs.
	Gap float32
}

// ReorderTransition lays out keyed items in a column; when their order changes they slide to their new places.
// Items added fade in and items removed fade out.
func ReorderTransition[K comparable](c *ui.Context, keys []K, opts ReorderTransitionOptions, item func(key K)) *ui.Element {
	if opts.Duration < 0 || opts.Gap < 0 {
		panic("mujica: ReorderTransition Duration and Gap must not be negative")
	}
	d := opts.Duration
	if d == 0 {
		d = 200 * time.Millisecond
	}
	seen := make(map[K]bool, len(keys))
	for _, k := range keys {
		if seen[k] {
			panic(fmt.Sprintf("mujica: ReorderTransition duplicate key %v", k))
		}
		seen[k] = true
	}
	move := ui.ElementTransition{Duration: d, Position: true, Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(opts.Gap)
	col.Children(func() {
		for _, k := range keys {
			ui.Column(c).Key(k).AlignItems(ui.Stretch).Transition(move).Children(func() { item(k) })
		}
	})
	return col
}
