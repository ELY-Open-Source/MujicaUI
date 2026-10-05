package overlay

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
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
	open, dismissed, focusIn bool
	enterAt, leaveAt         time.Time
	panel                    ui.Rect
}

// HoverCard previews content by a hovered or focused anchor; nil while hidden.
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
		// The popover's backdrop covers the anchor, so test the pointer by geometry.
		x, y, in := c.Root().PointerPosition()
		pointer = pointer || in && (b.Contains(x, y) || st.panel.Contains(x, y))
	}
	hot := pointer || keyboard || st.focusIn
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
			st.open, st.focusIn, st.leaveAt = false, false, time.Time{}
			return nil
		}
	}
	width := opts.Width
	if width == 0 {
		width = 280
	}
	core.OverlayLayer(c, st)
	// A popover keeps the card in the anchor's focus order and modal scope, even inside a Dialog.
	shown := true
	panel := ui.PopoverBase(c, anchor, &shown, func(panel *ui.Element) {
		core.OverlayStyle(c, panel, OverlayAnchored(c, panel, b))
		panel.Margin(6, 0, 0, 0).Width(PopoverWide(c, width)).Padding(core.Density(c).Unit() * 4).Gap(core.Density(c).Unit() * 2).
			Label(opts.Label).MaxHeight(PopoverRoom(c, b))
		ui.Scroll(c).Shrink(1).Children(content)
	})
	if !shown {
		slog.Debug("mujica: hover_card dismiss")
		if keyboard || st.focusIn {
			anchor.Focus() // at once, so the anchor stays hot and the card stays dismissed
		}
		st.open, st.dismissed, st.focusIn = false, true, false
		return panel
	}
	st.panel, st.focusIn = panel.Bounds(), panel.FocusWithin()
	return panel
}
