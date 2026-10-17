package code

import (
	"fmt"
	"log/slog"
	"path"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Breakpoint stops at File:Line when Condition holds.
type Breakpoint struct {
	ID        string
	File      string
	Line      int
	Enabled   bool
	Condition string
}

// BreakpointListOptions configure a BreakpointList.
type BreakpointListOptions struct {
	Disabled bool
}

// BreakpointListResult fields hold breakpoint IDs.
type BreakpointListResult struct {
	Element                 ui.Element
	Toggled, Edited, Opened string
}

type breakpointEdit struct{ id, draft string }

// BreakpointList lists breakpoints with toggles and conditions.
func BreakpointList(c *ui.Context, bps *[]Breakpoint, opts BreakpointListOptions) BreakpointListResult {
	if bps == nil {
		panic("mujica: BreakpointList needs breakpoints")
	}
	seen := map[string]bool{}
	for _, b := range *bps {
		if b.ID == "" || seen[b.ID] {
			panic("mujica: Breakpoint IDs must be set and unique")
		}
		seen[b.ID] = true
	}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	var res BreakpointListResult
	res.Element = CodePanel(c).Disabled(opts.Disabled).Label(core.Msg(c, "breakpoint_list.label"))
	edit := ui.Local(res.Element, "breakpoint_list.edit", func() breakpointEdit { return breakpointEdit{} })
	disabled := res.Element.IsDisabled()
	if disabled {
		edit.id = ""
	}
	res.Element.Children(func() {
		if len(*bps) == 0 {
			CodeCaption(c, core.Msg(c, "breakpoint_list.empty")).Padding(u * 4)
			return
		}
		layout.ScrollArea(c, nil, layout.ScrollAreaOptions{}, func() {
			for i := range *bps {
				b := &(*bps)[i]
				ui.Column(c.Key(b.ID)).Padding(u, u*2).Gap(u).BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).Children(func() {
					ui.Row(c).MinHeight(d.ControlHeight()).Gap(u * 2).AlignItems(ui.Center).Children(func() {
						on := b.Enabled
						input.Checkbox(c, &on, "", input.CheckboxOptions{Disabled: opts.Disabled}).Label(fmt.Sprintf(core.Msg(c, "breakpoint_list.enable"), path.Base(b.File), b.Line))
						if on != b.Enabled && !disabled {
							b.Enabled, res.Toggled = on, b.ID
							slog.Debug("mujica: breakpoint_list toggle", "id", b.ID, "enabled", on)
						}
						if breakpointLocation(c, *b).Clicked() && !disabled {
							res.Opened = b.ID
							slog.Debug("mujica: breakpoint_list open", "id", b.ID)
						}
						if core.IconAction(c, icons.Must("funnel"), fmt.Sprintf(core.Msg(c, "breakpoint_list.edit"), path.Base(b.File), b.Line)).Clicked() && !disabled {
							*edit = breakpointEdit{b.ID, b.Condition}
						}
					})
					switch {
					case edit.id == b.ID:
						in := input.TextInput(c, &edit.draft, input.TextInputOptions{Placeholder: core.Msg(c, "breakpoint_list.placeholder"), Label: core.Msg(c, "breakpoint_list.condition")})
						core.AutoFocus(c, in.Input)
						if in.Input.Shortcut(0, ui.KeyEscape) {
							edit.id = ""
						} else if in.Submitted() && !disabled {
							b.Condition, res.Edited, edit.id = edit.draft, b.ID, ""
							slog.Debug("mujica: breakpoint_list condition", "id", b.ID)
						}
					case b.Condition != "":
						CodeMono(c, fmt.Sprintf(core.Msg(c, "breakpoint_list.when"), b.Condition)).TextColor(k.Warning).Padding(0, 0, 0, 28)
					}
				})
			}
		}).Grow(1)
	})
	return res
}

// breakpointLocation is a button opening the location.
func breakpointLocation(c *ui.Context, b Breakpoint) ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	e := ui.ButtonBase(c).Grow(1).MinWidth(0).MinHeight(theme.MinHitSize).Justify(ui.Start).Gap(d.Unit()*2).Padding(0, d.Unit()).
		Radius(theme.ControlRadius).Label(fmt.Sprintf("%s:%d", b.File, b.Line))
	internal.Ring(e, k.Focus, theme.ControlRadius)
	if e.Hovered() {
		e.Background(k.SurfaceHover)
	}
	fg := k.Text
	if !b.Enabled {
		fg = k.TextMuted
	}
	return e.Children(func() {
		// A filled dot is enabled, a ring disabled.
		dot := ui.Box(c).Size(10, 10).Radius(5).Shrink(0).Role(ui.RoleNone).Border(theme.BorderWidth, k.Danger)
		if b.Enabled {
			dot.Background(k.Danger)
		}
		ui.Text(c, path.Base(b.File)).TextColor(fg).SingleLine()
		CodeCaption(c, fmt.Sprintf("%s:%d", path.Dir(b.File), b.Line)).MinWidth(0)
	})
}
