package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ToolParam is a parameter in a tool's input schema.
type ToolParam struct {
	Name        string
	Type        string
	Required    bool
	Description string
}

// AgentTool is a tool an agent may call; the caller keeps Enabled.
type AgentTool struct {
	Name        string
	Description string
	Enabled     bool
	Params      []ToolParam
}

// ToolRegistryPanelOptions configure a ToolRegistryPanel.
type ToolRegistryPanelOptions struct {
	Disabled bool
}

// ToolRegistryResult is a built ToolRegistryPanel and its events.
type ToolRegistryResult struct {
	Element *ui.Element
	toggled string
	on      bool
}

// Toggled returns the tool switched this frame and its new state.
func (r ToolRegistryResult) Toggled() (name string, on, ok bool) {
	return r.toggled, r.on, r.toggled != ""
}

// ToolRegistryPanel lists tools with on/off switches and previews the chosen tool's parameters.
func ToolRegistryPanel(c *ui.Context, s *ListState[string], tools []AgentTool, opts ToolRegistryPanelOptions) ToolRegistryResult {
	seen := map[string]bool{}
	for _, t := range tools {
		if t.Name == "" || seen[t.Name] {
			panic("mujica: AgentTool needs a unique Name")
		}
		seen[t.Name] = true
	}
	k := Tokens(c)
	u := Density(c).Unit()
	row := ui.Row(c).AlignItems(ui.Stretch).Gap(u * 3).Disabled(opts.Disabled)
	res := ToolRegistryResult{Element: row}
	row.Children(func() {
		ui.Column(c).Width(300).Shrink(0).AlignItems(ui.Stretch).Children(func() {
			List(c, s, len(tools), ListOptions[string]{
				Key: func(i int) string { return tools[i].Name }, Label: func(i int) string { return tools[i].Name },
				Empty: msg(c, "tool_registry_panel.empty"),
			}, func(i int) {
				t := tools[i]
				ui.Column(c).Grow(1).MinWidth(0).Children(func() {
					ui.Text(c, t.Name).Font(theme.MonoFont).FontWeight(600).SingleLine()
					ui.Text(c, t.Description).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
				})
				on := t.Enabled
				Switch(c, &on, "", SwitchOptions{}).Label(fmt.Sprintf(msg(c, "tool_registry_panel.enable"), t.Name))
				if on != t.Enabled && !row.IsDisabled() {
					res.toggled, res.on = t.Name, on
					slog.Debug("mujica: tool registry toggle", "tool", t.Name, "on", on)
				}
			}).Element.Grow(1)
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(u*2).Padding(u*3, u*4).Background(k.Surface).Border(theme.BorderWidth, k.Border).
			Radius(theme.CardRadius).Children(func() {
			name, ok := s.Selected()
			var tool *AgentTool
			for i := range tools {
				if ok && tools[i].Name == name {
					tool = &tools[i]
				}
			}
			if tool == nil {
				ui.Text(c, msg(c, "tool_registry_panel.choose")).TextColor(k.TextMuted)
				return
			}
			toolSchema(c, *tool)
		})
	})
	return res
}

// toolSchema previews a tool's parameters as a typed list.
func toolSchema(c *ui.Context, t AgentTool) {
	k := Tokens(c)
	u := Density(c).Unit()
	small := FontSize(c, theme.CaptionSize)
	ui.Text(c, t.Name).Font(theme.SerifFont).FontWeight(600).FontSize(FontSize(c, theme.H3Size))
	if t.Description != "" {
		ui.Text(c, t.Description).TextColor(k.TextMuted)
	}
	ui.Box(c).Height(1).Background(k.Ornament).Role(ui.RoleNone)
	if len(t.Params) == 0 {
		ui.Text(c, msg(c, "tool_registry_panel.no_params")).TextColor(k.TextMuted)
		return
	}
	for _, p := range t.Params {
		ui.Column(c).Gap(2).Children(func() {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
				ui.Text(c, p.Name).Font(theme.MonoFont).FontWeight(600)
				ui.Text(c, p.Type).Font(theme.MonoFont).FontSize(small).TextColor(k.Info)
				if p.Required {
					ui.Text(c, msg(c, "tool_registry_panel.required")).FontSize(small).TextColor(k.AccentText).
						Padding(0, u).Border(theme.BorderWidth, k.AccentText).Radius(theme.ControlRadius)
				}
			})
			if p.Description != "" {
				ui.Text(c, p.Description).FontSize(small).TextColor(k.TextMuted)
			}
		})
	}
}
