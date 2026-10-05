package mujica

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// LayoutTransitionOptions configure a LayoutTransition.
type LayoutTransitionOptions struct {
	// Duration of a change; zero means 200 ms.
	Duration time.Duration
	// Row lays the content in a row instead of a column.
	Row bool
	// Clip hides content outside the box while it resizes.
	Clip bool
}

// LayoutTransition holds content whose size and place animate from where they were (FLIP).
// Siblings that the change moves need a LayoutTransition of their own to move in step.
func LayoutTransition(c *ui.Context, opts LayoutTransitionOptions, content func()) *ui.Element {
	if opts.Duration < 0 {
		panic("mujica: LayoutTransition Duration is negative")
	}
	d := opts.Duration
	if d == 0 {
		d = 200 * time.Millisecond
	}
	box := ui.Column(c).AlignItems(ui.Stretch)
	if opts.Row {
		box.Row()
	}
	if opts.Clip {
		box.Clip()
	}
	box.Transition(ui.ElementTransition{Duration: d, Position: true, Size: true})
	if content != nil {
		box.Children(content)
	}
	return box
}
