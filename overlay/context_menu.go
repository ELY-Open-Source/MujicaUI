package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"
)

// ContextMenuOptions configure a ContextMenu.
type ContextMenuOptions struct {
	Disabled bool
}

// contextMenuState is where the menu opened, relative to the target.
type contextMenuState struct {
	open, fresh bool // fresh: opened by a press the menu has not yet seen
	x, y        float32
}

// ContextMenu gives the container target an in-window menu at the pointer or below the focus.
func ContextMenu(c *ui.Context, target ui.Element, opts ContextMenuOptions, build func(m *PopupMenu)) ui.Element {
	if !target.Valid() {
		panic("mujica: ContextMenu needs a target")
	}
	if build == nil {
		panic("mujica: ContextMenu needs a build function")
	}
	st := ui.Local(target, "mujica.contextmenu", func() contextMenuState { return contextMenuState{} })
	if opts.Disabled {
		st.open = false
		return target
	}
	target.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputPointerDown || ev.Button != 1 {
			return false
		}
		st.open, st.fresh, st.x, st.y = true, true, ev.X, ev.Y
		slog.Debug("mujica: context_menu open", "by", "pointer")
		return true
	})
	if target.Shortcut(0, ui.KeyContextMenu) || target.Shortcut(ui.Shift, ui.KeyF10) {
		st.open, st.x, st.y = true, 0, target.Bounds().H
		slog.Debug("mujica: context_menu open", "by", "keyboard")
	}
	if !st.open {
		return target
	}
	var anchor ui.Element
	target.Children(func() {
		anchor = ui.Box(c.Key("mujica.contextmenu.anchor")).Absolute().Left(st.x).Top(st.y).Size(1, 1).PassThrough()
	})
	panel := PopupMenuOpen(c, anchor, &st.open, build)
	// The opening press fell outside the new panel; it must not dismiss it.
	if st.fresh && panel.Valid() && panel.Dismissed() {
		st.open = true
	}
	st.fresh = false
	return target
}
