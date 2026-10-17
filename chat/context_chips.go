package chat

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ContextChipsOptions configure ContextChips.
type ContextChipsOptions struct {
	// Max sets visible chips before folding; zero means four.
	Max int
}

// RemoveResult carries the element and removed item ID.
type RemoveResult struct {
	Element ui.Element
	removed string
}

// Removed returns the ID whose remove button was pressed.
func (r RemoveResult) Removed() (string, bool) { return r.removed, r.removed != "" }

// ContextChips displays removable context with folded overflow.
func ContextChips(c *ui.Context, items []ContextItem, opts ContextChipsOptions) RemoveResult {
	if opts.Max < 0 {
		panic("mujica: ContextChips Max must not be negative")
	}
	limit := opts.Max
	if limit == 0 {
		limit = 4
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r RemoveResult
	r.Element = ui.Row(c).Gap(u * 2).Wrap().AlignItems(ui.Center).Role(ui.RoleList).Label(core.Msg(c, "context_chips.label"))
	expanded := ui.Local(r.Element, "expanded", func() bool { return false })
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		shown := len(items)
		if !*expanded {
			shown = min(shown, limit)
		}
		for _, it := range items[:shown] {
			ui.Row(c.Key(it.ID)).Gap(u).AlignItems(ui.Center).Padding(0, 0, 0, u*2).MinHeight(core.Density(c).ControlHeight()).
				Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).Role(ui.RoleListItem).
				Children(func() {
					ui.Icon(c, icons.Must(it.Kind.icon())).Size(14, 14).Shrink(0).TextColor(k.AccentText)
					ui.Text(c, it.Label).SingleLine().MaxWidth(180).FontSize(core.FontSize(c, theme.CaptionSize))
					if core.IconAction(c, icons.Must("x"), fmt.Sprintf(core.Msg(c, "context_chips.remove"), it.Label)).Clicked() && !disabled {
						r.removed = it.ID
						slog.Debug("mujica: context_chips remove", "id", it.ID)
					}
				})
		}
		if len(items) <= limit {
			return
		}
		label := fmt.Sprintf(core.Msg(c, "context_chips.more"), len(items)-limit)
		if *expanded {
			label = core.Msg(c, "context_chips.less")
		}
		b := ui.ButtonBase(c).MinHeight(core.Density(c).ControlHeight()).Padding(0, u*2).Radius(theme.ControlRadius).Label(label).TextColor(k.AccentText)
		internal.Ring(b, k.Focus, theme.ControlRadius)
		if b.Hovered() {
			b.Background(k.SurfaceHover)
		}
		b.Children(func() { ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine() })
		if b.Clicked() && !disabled {
			*expanded = !*expanded
			slog.Debug("mujica: context_chips fold", "expanded", *expanded)
		}
	})
	return r
}
