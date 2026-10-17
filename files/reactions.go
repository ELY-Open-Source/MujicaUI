package files

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Reaction is an emoji and who reacted with it.
type Reaction struct {
	Emoji string
	Count int
	// Mine marks the user's own reaction.
	Mine bool
	// People names who reacted, for the tooltip.
	People []string
}

// ReactionsOptions configure Reactions.
type ReactionsOptions struct {
	// Choices are the emoji the add picker offers; required.
	Choices  []string
	Disabled bool
}

// ReactionsResult exposes reactions and requested changes.
type ReactionsResult struct {
	Element ui.Element
	// Toggled identifies a reaction; Added identifies a picker choice.
	Toggled, Added string
}

// Reactions displays counts, toggles and an add picker.
func Reactions(c *ui.Context, reactions []Reaction, opts ReactionsOptions) ReactionsResult {
	if len(opts.Choices) == 0 {
		panic("mujica: Reactions needs Choices")
	}
	u := core.Density(c).Unit()
	var res ReactionsResult
	res.Element = ui.Row(c).Wrap().Gap(u * 2).AlignItems(ui.Center).Role(ui.RoleGroup).Label(core.Msg(c, "reactions.label")).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	open := ui.Local(res.Element, "reactions.open", func() bool { return false })
	*open = *open && !disabled
	res.Element.Children(func() {
		for _, r := range reactions {
			if r.Count < 1 || r.Emoji == "" {
				panic("mujica: Reaction needs an Emoji and a positive Count")
			}
			if reactionChip(c, r).Clicked() && !disabled {
				res.Toggled = r.Emoji
			}
		}
		add := core.IconAction(c, icons.Must("plus"), core.Msg(c, "reactions.add"))
		if add.Clicked() && !disabled {
			*open = !*open
		}
		overlay.Popover(c, add, open, overlay.PopoverOptions{Label: core.Msg(c, "reactions.pick")}, func() {
			ui.Grid(c).Columns(min(len(opts.Choices), 6)).Gap(u).FocusGroup(ui.Horizontal).Children(func() {
				for _, e := range opts.Choices {
					b := ui.ButtonBase(c.Key(e)).Size(36, 36).Radius(theme.ControlRadius).Label(e)
					internal.Ring(b, core.Tokens(c).Focus, theme.ControlRadius)
					if b.Hovered() {
						b.Background(core.Tokens(c).SurfaceHover)
					}
					b.Children(func() { ui.Text(c, e).FontSize(core.FontSize(c, 20)) })
					if b.Clicked() && !disabled {
						res.Added, *open = e, false
					}
				}
			})
		})
	})
	if res.Toggled != "" || res.Added != "" {
		slog.Debug("mujica: reactions change", "toggled", res.Toggled != "", "added", res.Added != "")
	}
	return res
}

// reactionChip frames and marks the current user's reaction.
func reactionChip(c *ui.Context, r Reaction) ui.Element {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := fmt.Sprintf(core.Msg(c, "reactions.count"), r.Emoji, r.Count)
	if r.Mine {
		label += ", " + core.Msg(c, "reactions.mine")
	}
	b := ui.ButtonBase(c.Key(r.Emoji)).MinHeight(theme.MinHitSize).MinWidth(theme.MinHitSize).Padding(0, u*2).Gap(u).
		Radius(theme.MinHitSize / 2).Role(ui.RoleToggleButton).Selected(r.Mine).Label(label).Transition(core.ColorFade)
	if len(r.People) > 0 {
		b.Tooltip(strings.Join(r.People, ", "))
	}
	internal.Ring(b, k.Focus, theme.MinHitSize/2)
	bg, border := k.Surface, k.Border
	if r.Mine {
		bg, border = k.Selection, k.AccentText
	}
	if b.Hovered() && !r.Mine {
		bg = k.SurfaceHover
	}
	b.Background(bg).Border(theme.BorderWidth, border)
	return b.Children(func() {
		ui.Text(c, r.Emoji)
		fg := k.TextMuted
		if r.Mine {
			fg = k.AccentText
		}
		ui.Text(c, fmt.Sprint(r.Count)).FontFeatures("tnum").FontWeight(600).TextColor(fg)
	})
}
