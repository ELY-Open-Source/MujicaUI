package input

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// HoldToConfirmOptions configure a HoldToConfirm.
type HoldToConfirmOptions struct {
	// Duration is how long to hold; zero is 1.5 s.
	Duration time.Duration
	Disabled bool
}

// holdToConfirmState tracks one hold: by key or pointer, since start.
type holdToConfirmState struct {
	key, done bool
	start     time.Time
}

// step advances the hold to now; it reports completion once per hold.
func (s *holdToConfirmState) step(held bool, now time.Time, d time.Duration) (progress float64, confirmed bool) {
	switch {
	case !held:
		s.start, s.done = time.Time{}, false
		return 0, false
	case s.done:
		return 1, false
	case s.start.IsZero():
		s.start = now
	}
	progress = min(1, float64(now.Sub(s.start))/float64(d))
	if progress >= 1 {
		s.done = true
		return 1, true
	}
	return progress, false
}

// HoldToConfirm runs a dangerous action only after the pointer or Space is held for Duration.
func HoldToConfirm(c *ui.Context, label string, opts HoldToConfirmOptions) core.ChoiceResult {
	if label == "" {
		panic("mujica: HoldToConfirm needs a label")
	}
	if opts.Duration < 0 {
		panic("mujica: negative HoldToConfirm Duration")
	}
	d := opts.Duration
	if d == 0 {
		d = 1500 * time.Millisecond
	}
	k := core.Tokens(c)
	den := core.Density(c)
	b := ui.ButtonBase(c).MinHeight(den.ControlHeight()).MinWidth(theme.MinHitSize).Padding(0, den.Unit()*4).
		Radius(theme.ControlRadius).Disabled(opts.Disabled).Description(core.Msg(c, "hold_to_confirm.hint")).Clip()
	internal.Ring(b, k.Focus, theme.ControlRadius)
	st := ui.Local(b, "hold_to_confirm", func() *holdToConfirmState { return &holdToConfirmState{} })
	s := *st
	disabled := b.IsDisabled()
	// Space holds; Enter and Space never click, so no single key press confirms.
	b.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Key != ui.KeySpace && ev.Key != ui.KeyEnter || disabled {
			return false
		}
		if ev.Key == ui.KeySpace && (ev.Kind == ui.InputKeyDown || ev.Kind == ui.InputKeyUp) {
			s.key = ev.Kind == ui.InputKeyDown
		}
		return true
	})
	_ = b.Clicked() // a click alone does nothing; it keeps the base pressable
	held := !disabled && (s.key || b.Pressed())
	if !b.Focused() {
		s.key = false
		held = !disabled && b.Pressed()
	}
	progress, confirmed := s.step(held, c.Now(), d)
	still := c.Preferences().ReduceMotion
	if held && !confirmed && progress < 1 {
		if still {
			c.After(100 * time.Millisecond)
		} else {
			c.AnimationFrame()
		}
	}
	bg, fg, _ := ButtonColors(k, Danger, b.Hovered() && !disabled, false)
	b.Background(bg).TextColor(fg)
	text := label
	if held && still && progress < 1 {
		left := d - time.Duration(progress*float64(d))
		text = fmt.Sprintf(core.Msg(c, "hold_to_confirm.left"), left.Seconds())
	}
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		if progress > 0 && !still {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * float32(progress), H: r.H}, k.OnDanger.Alpha(0.28), 0)
		}
	})
	// Room for the countdown, so the width holds.
	w1, _ := c.MeasureText(0, ui.Span{Text: label})
	w2, _ := c.MeasureText(0, ui.Span{Text: fmt.Sprintf(core.Msg(c, "hold_to_confirm.left"), d.Seconds())})
	b.MinWidth(max(w1, w2) + den.Unit()*8)
	b.Children(func() { ui.Text(c, text).SingleLine() })
	if confirmed {
		c.Announce(core.Msg(c, "hold_to_confirm.done"))
		slog.Debug("mujica: hold_to_confirm confirmed", "label", label)
	}
	return core.NewChoiceResult(b, confirmed)
}
