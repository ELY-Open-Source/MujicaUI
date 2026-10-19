package overlay

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// SplitButtonOptions configure a SplitButton.
type SplitButtonOptions struct {
	// Variant is Primary (the zero value) or Secondary.
	Variant input.ButtonVariant
	Icon    *ui.SVG
	// MenuLabel names the menu half; empty uses "More options".
	MenuLabel string
	Disabled  bool
}

// SplitButton runs its main action from the left half and opens build's menu from the right; each half takes focus.
func SplitButton(c *ui.Context, label string, opts SplitButtonOptions, build func(m *PopupMenu)) core.ChoiceResult {
	if label == "" || build == nil {
		panic("mujica: SplitButton needs a label and a menu")
	}
	if opts.Variant != input.Primary && opts.Variant != input.Secondary {
		panic(fmt.Sprintf("mujica: SplitButton Variant %d must be Primary or Secondary", opts.Variant))
	}
	k := core.Tokens(c)
	d := core.Density(c)
	menuName := opts.MenuLabel
	if menuName == "" {
		menuName = core.Msg(c, "split_button.more")
	}
	group := ui.Row(c).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	disabled := group.IsDisabled()
	open := ui.Local(group, "split_button.open", func() bool { return false })
	if disabled {
		*open = false
	}
	var clicked bool
	group.Children(func() {
		main := ui.ButtonBase(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*4).Gap(d.Unit()*2).
			Radius(theme.ControlRadius, 0, 0, theme.ControlRadius)
		splitButtonLook(c, k, main, opts.Variant, false)
		main.Children(func() {
			if opts.Icon != nil {
				ui.Icon(c, opts.Icon).Size(16, 16)
			}
			ui.Text(c, label).SingleLine()
		})
		clicked = main.Clicked() && !disabled
		if opts.Variant == input.Primary {
			ui.Box(c).Width(theme.BorderWidth).Background(k.OnAccent.Alpha(0.45)).Role(ui.RoleNone)
		}
		more := ui.ButtonBase(c).Role(ui.RoleMenuButton).Label(menuName).MinHeight(d.ControlHeight()).Width(theme.MinHitSize+d.Unit()).
			Radius(0, theme.ControlRadius, theme.ControlRadius, 0)
		splitButtonLook(c, k, more, opts.Variant, *open)
		if opts.Variant == input.Secondary {
			more.BorderWidth(theme.BorderWidth, theme.BorderWidth, theme.BorderWidth, 0)
		}
		more.Children(func() { ui.Icon(c, icons.Must("chevron-down")).Size(14, 14) })
		if !disabled && (more.Clicked() || more.Shortcut(0, ui.KeyDown)) {
			*open = true
			slog.Debug("mujica: split_button open")
		}
		// A tip coming and going ahead of the open menu would re-key its items.
		if !*open {
			more.Tooltip(menuName)
		}
		PopupMenuOpen(c, group, open, build)
	})
	if clicked {
		slog.Debug("mujica: split_button action", "label", label)
	}
	return core.NewChoiceResult(group, clicked)
}

// splitButtonLook colors one half as a Button of variant v would be.
func splitButtonLook(c *ui.Context, k theme.Tokens, b ui.Element, v input.ButtonVariant, open bool) {
	internal.Ring(b, k.Focus, theme.ControlRadius)
	bg, fg, border := input.ButtonColors(k, v, b.Hovered() && !b.IsDisabled(), b.Pressed() || open)
	b.Background(bg).TextColor(fg).Transition(core.ColorFade)
	if border.A > 0 {
		b.Border(theme.BorderWidth, border)
	}
}
