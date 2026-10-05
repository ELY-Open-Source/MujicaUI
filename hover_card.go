package mujica

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"
)

// HoverCardOptions configure a HoverCard.
type HoverCardOptions struct {
	// OpenDelay is 500 ms and CloseDelay 300 ms when zero.
	OpenDelay  time.Duration
	CloseDelay time.Duration
	// Width is 280 DIPs when zero.
	Width float32
	// Label names the card for assistive technology.
	Label string
}

// hoverCardState is what a card keeps between frames while its anchor lives.
type hoverCardState struct {
	open, dismissed bool
	// keyMode keeps the card in the focus order while the keyboard is on it.
	keyMode          bool
	enterAt, leaveAt time.Time
	panel            ui.Rect
}

// HoverCard previews content beside anchor while the pointer rests on it or
// it has the keyboard focus. The pointer may move into the card; with the
// keyboard, Tab from the anchor enters it. It closes after CloseDelay once
// neither holds it, or at once on Escape. It returns the card, or nil.
func HoverCard(c *ui.Context, anchor *ui.Element, opts HoverCardOptions, content func()) *ui.Element {
	if anchor == nil {
		panic("mujica: HoverCard needs an anchor")
	}
	if opts.OpenDelay < 0 || opts.CloseDelay < 0 {
		panic("mujica: HoverCard delays must not be negative")
	}
	openDelay, closeDelay := opts.OpenDelay, opts.CloseDelay
	if openDelay == 0 {
		openDelay = 500 * time.Millisecond
	}
	if closeDelay == 0 {
		closeDelay = 300 * time.Millisecond
	}
	st := ui.Local(anchor, "mujica.hovercard", func() hoverCardState { return hoverCardState{} })
	b := anchor.Bounds()
	keyboard := anchor.FocusVisible()
	pointer := anchor.Hovered()
	if st.open {
		// The card's layer covers the anchor, so test the pointer by geometry.
		x, y, in := c.Root().PointerPosition()
		pointer = pointer || in && (b.Contains(x, y) || st.panel.Contains(x, y))
	}
	hot := pointer || keyboard || st.keyMode
	if !hot {
		st.dismissed = false
	}
	if st.dismissed {
		return nil
	}
	now := c.Now()
	if hot {
		st.leaveAt = time.Time{}
		if !st.open {
			if st.enterAt.IsZero() {
				st.enterAt = now
			}
			if wait := openDelay - now.Sub(st.enterAt); wait > 0 {
				c.After(wait)
				return nil
			}
			slog.Debug("mujica: hover_card open", "keyboard", keyboard)
			st.open = true
		}
	} else {
		st.enterAt = time.Time{}
		if !st.open {
			return nil
		}
		if st.leaveAt.IsZero() {
			st.leaveAt = now
		}
		if wait := closeDelay - now.Sub(st.leaveAt); wait > 0 {
			c.After(wait)
		} else {
			slog.Debug("mujica: hover_card close")
			st.open, st.keyMode, st.leaveAt = false, false, time.Time{}
			return nil
		}
	}
	width := opts.Width
	if width == 0 {
		width = 280
	}
	style := func(panel *ui.Element, p float32) {
		overlayStyle(c, panel, p)
		panel.Width(width).Padding(Density(c).Unit() * 4).Gap(Density(c).Unit() * 2).Label(opts.Label).
			MaxHeight(popoverRoom(c, b))
		ui.Scroll(c).Shrink(1).Children(content)
	}
	var panel *ui.Element
	dismissed, keyMode := false, keyboard || st.keyMode
	if keyMode {
		shown := true
		panel = ui.PopoverBase(c, anchor, &shown, func(panel *ui.Element) {
			style(panel.Margin(6, 0, 0, 0), overlayAnchored(c, panel, b))
		})
		dismissed = !shown
	} else {
		ui.Overlay(c, func() {
			panel = ui.Column(c).Absolute().Role(ui.RolePopup)
			prev := panel.Bounds()
			x, y := overlayPlace(c, b, width, prev.H, 6)
			panel.Left(x).Top(y)
			p := overlayEnter(c, panel)
			if prev.H == 0 {
				p = 0 // not measured yet
				c.AnimationFrame()
			}
			panel.Children(func() { style(panel, p) })
		})
		dismissed = c.Shortcut(0, ui.KeyEscape)
	}
	if dismissed {
		slog.Debug("mujica: hover_card dismiss")
		st.open, st.dismissed, st.keyMode = false, true, false
		if keyMode {
			anchor.Focus() // at once, so the anchor stays hot and the card stays dismissed
		}
		return panel
	}
	st.panel, st.keyMode = panel.Bounds(), keyboard || keyMode && panel.FocusWithin()
	return panel
}
