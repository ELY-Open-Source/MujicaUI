package mujica

import (
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// LayoutTransitionOptions configure a LayoutTransition.
type LayoutTransitionOptions struct {
	// Duration of a change; zero means StateDuration.
	Duration time.Duration
	// Row lays the content in a row.
	Row bool
	// Clip hides content outside the box while it resizes.
	Clip bool
}

// LayoutTransition animates its size and place when they change.
func LayoutTransition(c *ui.Context, opts LayoutTransitionOptions, content func()) *ui.Element {
	if opts.Duration < 0 {
		panic("mujica: LayoutTransition Duration is negative")
	}
	d := opts.Duration
	if d == 0 {
		d = theme.StateDuration
	}
	box := ui.Column(c).AlignItems(ui.Stretch)
	if opts.Row {
		box.Row()
	}
	if opts.Clip {
		box.Clip()
	}
	box.Transition(ui.ElementTransition{Duration: Motion(c, d), Position: true, Size: true})
	if content != nil {
		box.Children(content)
	}
	return box
}
