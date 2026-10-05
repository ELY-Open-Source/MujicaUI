package mujica

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// StaggerOptions configure a Stagger.
type StaggerOptions struct {
	// Interval between items; zero means 60 ms.
	Interval time.Duration
	// Duration of each item's entry; zero means PopupDuration.
	Duration time.Duration
	// Horizontal lays items in a row.
	Horizontal bool
}

type staggerState struct {
	born time.Time
	done bool
}

// Stagger builds n items entering one after another, once.
func Stagger(c *ui.Context, n int, opts StaggerOptions, item func(i int)) *ui.Element {
	if n < 0 || opts.Interval < 0 || opts.Duration < 0 {
		panic("mujica: Stagger needs non-negative n, Interval and Duration")
	}
	every, d := opts.Interval, opts.Duration
	if every == 0 {
		every = 60 * time.Millisecond
	}
	if d == 0 {
		d = theme.PopupDuration
	}
	box := ui.Column(c).AlignItems(ui.Stretch)
	if opts.Horizontal {
		box.Row().AlignItems(ui.Start)
	}
	st := ui.Local(box, "stagger", func() *staggerState { return &staggerState{born: c.Now()} })
	s := *st
	if n == 0 {
		return box
	}
	upto := n
	if !s.done && !c.Preferences().ReduceMotion {
		upto = min(n, 1+int(c.Now().Sub(s.born)/every))
		if upto < n {
			box.DrawOver(func(*ui.Painter, ui.Rect) { c.After(s.born.Add(time.Duration(upto) * every).Sub(c.Now())) })
		} else {
			s.done = true
			slog.Debug("mujica: stagger done", "items", n)
		}
	}
	if c.Preferences().ReduceMotion {
		s.done = true
	}
	enter := &ui.Motion{Y: 8}
	if opts.Horizontal {
		enter = &ui.Motion{X: 8}
	}
	box.Children(func() {
		for i := range upto {
			ui.Column(c).Key(i).AlignItems(ui.Stretch).
				Transition(ui.ElementTransition{Duration: Motion(c, d), Enter: enter}).Children(func() { item(i) })
		}
	})
	return box
}
