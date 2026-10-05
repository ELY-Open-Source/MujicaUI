package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// RegenerateKind is what a RegenerateMenu choice asks for.
type RegenerateKind int

const (
	// RegenerateAgain retries with the same model.
	RegenerateAgain RegenerateKind = iota
	// RegenerateModel retries with Choice.Model.
	RegenerateModel
	// RegenerateVariant retries asking for Choice.Variant, such as longer or shorter.
	RegenerateVariant
)

// RegenerateChoice is the retry the user picked.
type RegenerateChoice struct {
	Kind    RegenerateKind
	Model   string
	Variant string
}

// RegenerateMenuOptions configure a RegenerateMenu.
type RegenerateMenuOptions struct {
	// Models lists models to retry with; Current is marked and left out.
	Models  []string
	Current string
	// Variants replaces the built-in Longer, Shorter, Simpler and More detailed.
	Variants []string
	Disabled bool
}

// RegenerateMenuView is a built menu and the retry chosen this frame.
type RegenerateMenuView struct {
	Element *ui.Element
	Chosen  bool
	Choice  RegenerateChoice
}

// RegenerateMenu offers regenerate, retry with another model, and variants such as longer or shorter.
func RegenerateMenu(c *ui.Context, opts RegenerateMenuOptions) RegenerateMenuView {
	variants := opts.Variants
	if variants == nil {
		for _, key := range []string{"regenerate_menu.longer", "regenerate_menu.shorter", "regenerate_menu.simpler", "regenerate_menu.detailed"} {
			variants = append(variants, msg(c, key))
		}
	}
	var v RegenerateMenuView
	pick := func(ch RegenerateChoice) { v.Chosen, v.Choice = true, ch }
	v.Element = DropdownMenu(c, msg(c, "regenerate_menu.regenerate"), DropdownMenuOptions{Icon: icons.Must("refresh-cw"), Disabled: opts.Disabled}, func(m *PopupMenu) {
		if m.Item(msg(c, "regenerate_menu.again"), PopupMenuItemOptions{}) {
			pick(RegenerateChoice{Kind: RegenerateAgain})
		}
		others := 0
		for _, model := range opts.Models {
			if model != opts.Current {
				others++
			}
		}
		m.Submenu(msg(c, "regenerate_menu.model"), PopupMenuItemOptions{Disabled: others == 0}, func(m *PopupMenu) {
			for _, model := range opts.Models {
				if model != opts.Current && m.Item(model, PopupMenuItemOptions{}) {
					pick(RegenerateChoice{Kind: RegenerateModel, Model: model})
				}
			}
		})
		m.Separator()
		m.Group(msg(c, "regenerate_menu.variants"), func() {
			for _, variant := range variants {
				if m.Item(variant, PopupMenuItemOptions{}) {
					pick(RegenerateChoice{Kind: RegenerateVariant, Variant: variant})
				}
			}
		})
	})
	if v.Element.IsDisabled() {
		v.Chosen = false
	}
	if v.Chosen {
		slog.Debug("mujica: regenerate_menu choose", "kind", v.Choice.Kind, "model", v.Choice.Model, "variant", v.Choice.Variant)
	}
	return v
}
