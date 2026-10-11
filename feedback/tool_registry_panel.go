package feedback

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ToolParam is a parameter in a tool's input schema.
type ToolParam struct {
	Name        string
	Type        string
	Required    bool
	Description string
}

// AgentTool is a callable tool; the caller keeps Enabled.
type AgentTool struct {
	Name        string
	Description string
	Enabled     bool
	Params      []ToolParam
}

// ToolRegistryPanelOptions configure a ToolRegistryPanel.
type ToolRegistryPanelOptions struct {
	Disabled bool
	// Version changes when tool names are edited in place.
	Version int
}

// ToolRegistryResult is a built ToolRegistryPanel and its events.
type ToolRegistryResult struct {
	Element *ui.Element
	List    data.ListView
	toggled string
	on      bool
}

// Toggled returns the tool switched this frame.
func (r ToolRegistryResult) Toggled() (name string, on, ok bool) {
	return r.toggled, r.on, r.toggled != ""
}

type toolRegistryIndex struct {
	src        *AgentTool
	n, version int
	byName     map[string]int
}

// ToolRegistryPanel lists tools with switches and a schema preview.
func ToolRegistryPanel(c *ui.Context, s *data.ListState[string], tools []AgentTool, opts ToolRegistryPanelOptions) ToolRegistryResult {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	row := ui.Row(c).AlignItems(ui.Stretch).Gap(u * 3).Disabled(opts.Disabled)
	res := ToolRegistryResult{Element: row}
	index := ui.Local(row, "index", func() toolRegistryIndex { return toolRegistryIndex{} })
	var first *AgentTool
	if len(tools) > 0 {
		first = &tools[0]
	}
	if index.byName == nil || index.src != first || index.n != len(tools) || index.version != opts.Version {
		byName := make(map[string]int, len(tools))
		for i, t := range tools {
			if _, exists := byName[t.Name]; t.Name == "" || exists {
				panic("mujica: AgentTool needs a unique Name")
			}
			byName[t.Name] = i
		}
		*index = toolRegistryIndex{src: first, n: len(tools), version: opts.Version, byName: byName}
	}
	row.Children(func() {
		ui.Column(c).Width(300).Shrink(0).AlignItems(ui.Stretch).Children(func() {
			res.List = data.List(c, s, len(tools), data.ListOptions[string]{
				Key: func(i int) string { return tools[i].Name }, Label: func(i int) string { return tools[i].Name },
				Empty: core.Msg(c, "tool_registry_panel.empty"),
			}, func(i int) {
				t := tools[i]
				ui.Column(c).Grow(1).MinWidth(0).Children(func() {
					ui.Text(c, t.Name).Font(theme.MonoFont).FontWeight(600).SingleLine()
					ui.Text(c, t.Description).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
				})
				on := t.Enabled
				input.Switch(c, &on, "", input.SwitchOptions{}).Label(fmt.Sprintf(core.Msg(c, "tool_registry_panel.enable"), t.Name))
				if on != t.Enabled && !row.IsDisabled() {
					res.toggled, res.on = t.Name, on
					slog.Debug("mujica: tool registry toggle", "tool", t.Name, "on", on)
				}
			})
			res.List.Element.Grow(1)
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(u*2).Padding(u*3, u*4).Background(k.Surface).Border(theme.BorderWidth, k.Border).
			Radius(theme.CardRadius).Children(func() {
			name, ok := s.Selected()
			i, known := index.byName[name]
			if !ok || !known {
				ui.Text(c, core.Msg(c, "tool_registry_panel.choose")).TextColor(k.TextMuted)
				return
			}
			toolSchema(c, tools[i])
		})
	})
	return res
}

// toolSchema previews a tool's parameters as a typed list.
func toolSchema(c *ui.Context, t AgentTool) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	ui.Text(c, t.Name).Font(theme.SerifFont).FontWeight(600).FontSize(core.FontSize(c, theme.H3Size))
	if t.Description != "" {
		ui.Text(c, t.Description).TextColor(k.TextMuted)
	}
	ui.Box(c).Height(1).Background(k.Ornament).Role(ui.RoleNone)
	if len(t.Params) == 0 {
		ui.Text(c, core.Msg(c, "tool_registry_panel.no_params")).TextColor(k.TextMuted)
		return
	}
	ui.Column(c).Key(t.Name).Grow(1).MinHeight(0).Children(func() {
		ui.List(c, nil, len(t.Params), func(i int) {
			p := t.Params[i]
			ui.Column(c).Gap(2).MinHeight(core.Density(c).ControlHeight()).Children(func() {
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					ui.Text(c, p.Name).Font(theme.MonoFont).FontWeight(600)
					ui.Text(c, p.Type).Font(theme.MonoFont).FontSize(small).TextColor(k.Info)
					if p.Required {
						ui.Text(c, core.Msg(c, "tool_registry_panel.required")).FontSize(small).TextColor(k.AccentText).
							Padding(0, u).Border(theme.BorderWidth, k.AccentText).Radius(theme.ControlRadius)
					}
				})
				if p.Description != "" {
					ui.Text(c, p.Description).FontSize(small).TextColor(k.TextMuted)
				}
			})
		}).Grow(1).MinHeight(0)
	})
}
