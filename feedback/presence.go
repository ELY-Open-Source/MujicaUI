package feedback

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// PresenceKind is how content comes and goes.
type PresenceKind int

const (
	// PresenceFade fades in and out.
	PresenceFade PresenceKind = iota
	// PresenceSlide fades while moving up into place and down out.
	PresenceSlide
	// PresenceCollapse grows open and shrinks shut along the parent's direction.
	PresenceCollapse
)

// PresenceOptions configure a Presence.
type PresenceOptions struct {
	Kind PresenceKind
	// Duration is PopupDuration when zero.
	Duration time.Duration
}

// presenceMotion is where content comes from and goes.
func presenceMotion(kind PresenceKind) *ui.Motion {
	switch kind {
	case PresenceFade:
		return &ui.Motion{}
	case PresenceSlide:
		return &ui.Motion{Y: 8}
	case PresenceCollapse:
		return &ui.Motion{Collapse: true, Opacity: 1}
	}
	panic("mujica: unknown PresenceKind")
}

type presenceState struct{ shown bool }

// Presence builds content while show, entering and exiting with motion.
func Presence(c *ui.Context, show bool, opts PresenceOptions, content func()) ui.Element {
	m := presenceMotion(opts.Kind)
	if opts.Duration < 0 {
		panic("mujica: Presence Duration is negative")
	}
	d := opts.Duration
	if d == 0 {
		d = theme.PopupDuration
	}
	holder := ui.Column(c).AlignItems(ui.Stretch).Shrink(0)
	st := ui.Local(holder, "presence", func() presenceState { return presenceState{shown: show} })
	if st.shown != show {
		st.shown = show
		slog.Debug("mujica: presence", "show", show)
	}
	holder.Children(func() {
		if !show {
			if c.Preferences().ReduceMotion {
				// Rebuilding the empty slot cancels MyGo's active exit.
				ui.Column(c.Key("presence")).Invisible().Transition(ui.ElementTransition{Duration: d, Position: true})
			}
			return
		}
		box := ui.Column(c.Key("presence")).AlignItems(ui.Stretch)
		box.Transition(ui.ElementTransition{Duration: core.Motion(c, d), Enter: m, Exit: m})
		box.Children(content)
	})
	return holder
}
