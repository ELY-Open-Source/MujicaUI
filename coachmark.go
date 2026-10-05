package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// CoachmarkStep uses an anchor built in the current frame.
type CoachmarkStep struct {
	Anchor      *ui.Element
	Title, Body string
}

// CoachmarkResult reports tour completion or dismissal.
type CoachmarkResult struct {
	Element             *ui.Element
	finished, dismissed bool
}

// Finished reports Done on the last step.
func (r CoachmarkResult) Finished() bool { return r.finished }

// Dismissed reports Escape, a click outside or Skip.
func (r CoachmarkResult) Dismissed() bool { return r.dismissed }

// Coachmark provides anchored tips and keyboard-accessible step navigation.
func Coachmark(c *ui.Context, open *bool, step *int, steps []CoachmarkStep) CoachmarkResult {
	if len(steps) == 0 || *step < 0 || *step >= len(steps) {
		panic(fmt.Sprintf("mujica: Coachmark step %d outside %d steps", *step, len(steps)))
	}
	var r CoachmarkResult
	if !*open {
		return r
	}
	cur := steps[*step]
	if cur.Anchor == nil || cur.Title == "" {
		panic("mujica: CoachmarkStep needs an Anchor and a Title")
	}
	if cur.Anchor.IsDisabled() {
		*open = false
		return r
	}
	k := Tokens(c)
	d := Density(c)
	if b := cur.Anchor.Bounds(); b.W > 0 {
		const pad = 4
		ui.Overlay(c, func() {
			ui.Box(c).Absolute().Left(b.X-pad).Top(b.Y-pad).Size(b.W+2*pad, b.H+2*pad).PassThrough().Role(ui.RoleNone).
				Border(2, k.AccentText).Radius(theme.CardRadius)
		})
	}
	from := *step
	r.Element = Popover(c, cur.Anchor, open, PopoverOptions{Title: cur.Title, Width: 300}, func() {
		ui.Column(c).Gap(d.Unit() * 3).Children(func() {
			if cur.Body != "" {
				ui.Text(c, cur.Body)
			}
			ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprintf(msg(c, "coachmark.count"), *step+1, len(steps))).TextColor(k.TextMuted).
					FontSize(FontSize(c, theme.CaptionSize)).Grow(1)
				if *step < len(steps)-1 && Button(c, msg(c, "coachmark.skip"), ButtonOptions{Variant: Ghost}).Clicked() {
					*open, r.dismissed = false, true
				}
				if *step > 0 && Button(c, msg(c, "coachmark.back"), ButtonOptions{Variant: Secondary}).Clicked() {
					*step--
				}
				last := *step == len(steps)-1
				next := msg(c, "coachmark.next")
				if last {
					next = msg(c, "coachmark.done")
				}
				if Button(c, next, ButtonOptions{}).AutoFocus().Clicked() {
					if last {
						*open, r.finished = false, true
					} else {
						*step++
					}
				}
			})
		})
	})
	if !*open && !r.finished {
		r.dismissed = true
	}
	if *step != from || !*open {
		slog.Debug("mujica: coachmark step", "step", *step, "open", *open)
	}
	return r
}
