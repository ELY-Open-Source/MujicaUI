package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"
)

// DropdownMenuOptions configure a DropdownMenu.
type DropdownMenuOptions struct {
	// Icon goes before the label.
	Icon     *ui.SVG
	Disabled bool
	// Label names an icon-only trigger for assistive technology.
	Label string
}

// DropdownMenu is a button opening an in-window menu that build fills; it returns the button.
func DropdownMenu(c *ui.Context, label string, opts DropdownMenuOptions, build func(m *PopupMenu)) ui.Element {
	if build == nil {
		panic("mujica: DropdownMenu needs a build function")
	}
	if label == "" && opts.Label == "" {
		panic("mujica: DropdownMenu needs a label or Label")
	}
	b := ui.ButtonBase(c).Role(ui.RoleMenuButton).Disabled(opts.Disabled)
	if opts.Label != "" {
		b.Label(opts.Label)
	}
	open := ui.Local(b, "mujica.dropdown", func() bool { return false })
	if opts.Disabled {
		*open = false
	}
	if b.Clicked() || b.Shortcut(0, ui.KeyDown) {
		*open = true
		slog.Debug("mujica: dropdown_menu open")
	}
	popupMenuTrigger(c, b, label, opts.Icon, *open)
	PopupMenuOpen(c, b, open, build)
	return b
}
