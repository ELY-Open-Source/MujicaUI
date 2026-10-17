package feedback

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// UndoToastOptions configure an UndoToast.
type UndoToastOptions struct {
	// OnUndo is required; it runs when the user undoes.
	OnUndo func()
	// Duration is the countdown; zero is ToastDuration.
	Duration time.Duration
	// Action labels the button; empty uses "Undo".
	Action string
}

// UndoToastResult reports how an UndoToast ended this frame, if it did.
type UndoToastResult struct {
	Element ui.Element
	Undone  bool
	// Expired: the countdown ran out, or the user closed it; the action stands.
	Expired bool
}

// undoToastClock counts the time shown, stopped while hovered or focused.
type undoToastClock struct {
	message string
	spent   time.Duration
	runFrom time.Time
}

// UndoToast shows message at the bottom of the window with a countdown while *open; it never takes the focus.
func UndoToast(c *ui.Context, open *bool, message string, opts UndoToastOptions) UndoToastResult {
	if open == nil || message == "" || opts.OnUndo == nil {
		panic("mujica: UndoToast needs open, a message and OnUndo")
	}
	if opts.Duration < 0 {
		panic("mujica: negative UndoToast Duration")
	}
	clock := ui.Local(c.Root(), "undo_toast", func() *undoToastClock { return &undoToastClock{} })
	s := *clock
	var res UndoToastResult
	if !*open {
		*s = undoToastClock{}
		return res
	}
	now := c.Now()
	life := cmp.Or(opts.Duration, ToastDuration)
	if s.message != message {
		*s = undoToastClock{message: message, runFrom: now}
		c.Announce(message)
	}
	spent := s.spent
	if !s.runFrom.IsZero() {
		spent += now.Sub(s.runFrom)
	}
	if spent >= life {
		*open, res.Expired = false, true
		*s = undoToastClock{}
		slog.Debug("mujica: undo_toast expired")
		return res
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	ww, _ := c.Size()
	left := life - spent
	ui.Overlay(c, func() {
		dock := ui.Row(c).Absolute().Left(0).Right(0).Bottom(u * 4).Justify(ui.Center).PassThrough()
		dock.Children(func() {
			card := ui.Column(c).Width(min(420, ww-8*u)).Radius(theme.PopupRadius).Background(k.Surface).
				Border(theme.BorderWidth, k.ControlBorder).Role(ui.RoleStatus).Label(message).Clip()
			card.Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, 0.14))
			card.Transition(ui.ElementTransition{Duration: theme.PopupDuration, Enter: &ui.Motion{Y: 8}, Exit: &ui.Motion{}})
			res.Element = card
			card.Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(u*3).Padding(u*2, u*2, u*2, u*4).Children(func() {
					ui.Text(c, message).TextColor(k.Text).Grow(1).Shrink(1)
					secs := int(math.Ceil(left.Seconds()))
					label := fmt.Sprintf("%s (%d)", cmp.Or(opts.Action, core.Msg(c, "undo_toast.undo")), secs)
					if input.Button(c, label, input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("refresh-cw")}).Clicked() {
						res.Undone = true
					}
					if core.IconAction(c, icons.Must("x"), core.Msg(c, "toast.close")).Clicked() {
						res.Expired = true
					}
				})
				still := c.Preferences().ReduceMotion
				ui.Box(c).FillWidth().Height(2).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
					if still {
						return
					}
					rest := left
					if !s.runFrom.IsZero() {
						rest = life - s.spent - p.Now().Sub(s.runFrom)
						p.AnimationFrame()
					}
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * max(0, float32(rest)/float32(life)), H: r.H}, k.AccentText, 0)
				})
			})
			if card.Shortcut(0, ui.KeyEscape) {
				res.Expired = true
			}
			// Hover or focus holds the countdown.
			switch held := card.Hovered() || card.FocusWithin(); {
			case held && !s.runFrom.IsZero():
				s.spent, s.runFrom = spent, time.Time{}
			case !held && s.runFrom.IsZero():
				s.runFrom = now
			}
		})
	})
	switch {
	case res.Undone:
		*open, res.Expired = false, false
		*s = undoToastClock{}
		slog.Debug("mujica: undo_toast undo")
		opts.OnUndo()
	case res.Expired:
		*open = false
		*s = undoToastClock{}
		slog.Debug("mujica: undo_toast closed")
	case !s.runFrom.IsZero():
		// Wake for the next whole second of the count, and for the end.
		c.After(min(left, left-left.Truncate(time.Second)+time.Millisecond))
	}
	return res
}
