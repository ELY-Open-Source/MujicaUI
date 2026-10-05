package chat

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/overlay"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ToolItem is one tool a ToolToggleMenu turns on or off.
type ToolItem struct {
	// ID and On are required.
	ID          string
	Label       string
	Description string
	On          *bool
	// NeedsAuth marks tools requiring caller-owned authorization.
	NeedsAuth, Authorized bool
}

// ToolToggleResult carries the trigger and requested action.
type ToolToggleResult struct {
	Element            *ui.Element
	changed, authorize string
}

// Changed returns the toggled tool ID.
func (r ToolToggleResult) Changed() (string, bool) { return r.changed, r.changed != "" }

// Authorize returns the tool ID requesting authorization.
func (r ToolToggleResult) Authorize() (string, bool) { return r.authorize, r.authorize != "" }

// ToolToggleMenu displays tool switches and authorization actions.
func ToolToggleMenu(c *ui.Context, tools []ToolItem) ToolToggleResult {
	on := 0
	for _, t := range tools {
		if t.ID == "" || t.On == nil {
			panic("mujica: ToolItem needs ID and On")
		}
		if *t.On {
			on++
		}
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := fmt.Sprintf(core.Msg(c, "tool_toggle_menu.count"), on)
	var r ToolToggleResult
	b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Padding(0, u*2).Gap(u).Radius(theme.ControlRadius).Label(label).TextColor(k.TextMuted)
	r.Element = b
	internal.Ring(b, k.Focus, theme.ControlRadius)
	open := ui.Local(b, "open", func() bool { return false })
	switch {
	case *open:
		b.Background(k.Selection).TextColor(k.Text)
	case b.Hovered():
		b.Background(k.SurfaceHover).TextColor(k.Text)
	}
	b.Children(func() {
		ui.Icon(c, icons.Must("sliders-horizontal")).Size(16, 16)
		if on > 0 {
			ui.Text(c, fmt.Sprint(on)).FontFeatures("tnum").FontWeight(600).TextColor(k.AccentText).SingleLine()
		}
	})
	disabled := b.IsDisabled()
	if b.Clicked() && !disabled {
		*open = !*open
		slog.Debug("mujica: tool_toggle_menu fold", "open", *open)
	}
	if disabled {
		*open = false
	}
	overlay.Popover(c, b, open, overlay.PopoverOptions{Title: core.Msg(c, "tool_toggle_menu.label"), Width: 320}, func() {
		for _, t := range tools {
			toolToggleRow(c, t, &r)
		}
	})
	if r.changed != "" || r.authorize != "" {
		slog.Debug("mujica: tool_toggle_menu", "changed", r.changed, "authorize", r.authorize)
	}
	return r
}

// toolToggleRow displays a switch or authorization action.
func toolToggleRow(c *ui.Context, t ToolItem, r *ToolToggleResult) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	locked := t.NeedsAuth && !t.Authorized
	ui.Row(c).Key(t.ID).Gap(u * 3).AlignItems(ui.Center).Children(func() {
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			ui.Row(c).Gap(u).AlignItems(ui.Center).MinHeight(core.Density(c).ControlHeight()).Children(func() {
				if locked {
					ui.Text(c, t.Label).FontWeight(600).SingleLine()
				} else if input.Switch(c, t.On, t.Label, input.SwitchOptions{}).Changed() {
					r.changed = t.ID
				}
				if t.NeedsAuth {
					icon := "key-round"
					if t.Authorized {
						icon = "shield-check"
					}
					ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(k.Warning).Label(core.Msg(c, "tool_toggle_menu.needs_auth"))
				}
			})
			if t.Description != "" {
				ui.Text(c, t.Description).FontSize(small).TextColor(k.TextMuted)
			}
			if locked {
				ui.Text(c, core.Msg(c, "tool_toggle_menu.needs_auth")).FontSize(small).TextColor(k.Warning)
			}
		})
		if locked {
			full := fmt.Sprintf(core.Msg(c, "tool_toggle_menu.authorize"), t.Label)
			if input.Button(c, core.Msg(c, "tool_toggle_menu.authorize_short"), input.ButtonOptions{Variant: input.Outline, Label: full}).Clicked() {
				r.authorize = t.ID
			}
		}
	})
}
