package chat

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
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
	// Models lists models to retry with; Current is left out.
	Models  []string
	Current string
	// Variants replaces the built-in Longer, Shorter, Simpler and More detailed.
	Variants []string
	Disabled bool
}

// RegenerateMenuView is a built menu and the chosen retry.
type RegenerateMenuView struct {
	Element *ui.Element
	Chosen  bool
	Choice  RegenerateChoice
}

// RegenerateMenu offers retry, another model, and rewrite variants.
func RegenerateMenu(c *ui.Context, opts RegenerateMenuOptions) RegenerateMenuView {
	variants := opts.Variants
	if variants == nil {
		for _, key := range []string{"regenerate_menu.longer", "regenerate_menu.shorter", "regenerate_menu.simpler", "regenerate_menu.detailed"} {
			variants = append(variants, core.Msg(c, key))
		}
	}
	var v RegenerateMenuView
	pick := func(ch RegenerateChoice) { v.Chosen, v.Choice = true, ch }
	v.Element = overlay.DropdownMenu(c, core.Msg(c, "regenerate_menu.regenerate"), overlay.DropdownMenuOptions{Icon: icons.Must("refresh-cw"), Disabled: opts.Disabled}, func(m *overlay.PopupMenu) {
		if m.Item(core.Msg(c, "regenerate_menu.again"), overlay.PopupMenuItemOptions{}) {
			pick(RegenerateChoice{Kind: RegenerateAgain})
		}
		others := 0
		for _, model := range opts.Models {
			if model != opts.Current {
				others++
			}
		}
		m.Submenu(core.Msg(c, "regenerate_menu.model"), overlay.PopupMenuItemOptions{Disabled: others == 0}, func(m *overlay.PopupMenu) {
			for _, model := range opts.Models {
				if model != opts.Current && m.Item(model, overlay.PopupMenuItemOptions{}) {
					pick(RegenerateChoice{Kind: RegenerateModel, Model: model})
				}
			}
		})
		m.Separator()
		m.Group(core.Msg(c, "regenerate_menu.variants"), func() {
			for _, variant := range variants {
				if m.Item(variant, overlay.PopupMenuItemOptions{}) {
					pick(RegenerateChoice{Kind: RegenerateVariant, Variant: variant})
				}
			}
		})
	})
	v.Element.AlignSelf(ui.Start)
	if v.Element.IsDisabled() {
		v.Chosen = false
	}
	if v.Chosen {
		slog.Debug("mujica: regenerate_menu choose", "kind", v.Choice.Kind, "model", v.Choice.Model, "variant", v.Choice.Variant)
	}
	return v
}
