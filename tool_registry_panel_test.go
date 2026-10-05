package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func toolSample() []AgentTool {
	return []AgentTool{
		{Name: "web_search", Description: "Search the web", Enabled: true, Params: []ToolParam{
			{Name: "query", Type: "string", Required: true, Description: "What to look for"},
			{Name: "limit", Type: "integer"},
		}},
		{Name: "read_file", Description: "Read a local file"},
	}
}

func TestToolRegistryPanelToggleAndPreview(t *testing.T) {
	var s ListState[string]
	tools := toolSample()
	disabled := false
	var toggled []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Height(400).Children(func() {
			r := ToolRegistryPanel(c, &s, tools, ToolRegistryPanelOptions{Disabled: disabled})
			r.Element.Grow(1)
			if name, on, ok := r.Toggled(); ok {
				toggled = append(toggled, name)
				for i := range tools {
					if tools[i].Name == name {
						tools[i].Enabled = on
					}
				}
			}
		})
	}, 800, 420)
	if !tt.HasText("Choose a tool to see its parameters") {
		t.Fatalf("empty preview: %q", tt.Texts())
	}
	if err := tt.Click("Search the web"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("query") || !tt.HasText("required") || !tt.HasText("integer") {
		t.Fatalf("schema: %q", tt.Texts())
	}
	if err := tt.Click("read_file"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("No parameters") {
		t.Fatalf("no params: %q", tt.Texts())
	}
	if err := tt.Click("Enable read_file"); err != nil {
		t.Fatal(err)
	}
	if len(toggled) != 1 || toggled[0] != "read_file" || !tools[1].Enabled {
		t.Fatalf("toggle %q enabled %v", toggled, tools[1].Enabled)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Enable read_file")
	if len(toggled) != 1 {
		t.Fatal("disabled panel toggled")
	}
}

func TestToolRegistryPanelInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		var s ListState[string]
		ToolRegistryPanel(c, &s, []AgentTool{{Name: "a"}, {Name: "a"}}, ToolRegistryPanelOptions{})
	}, 300, 200, 1)
}

func TestToolRegistryPanelSnapshot(t *testing.T) {
	tools := toolSample()
	snapshotAfter(t, "tool_registry_panel", 760, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() *ListState[string] { return &ListState[string]{} })
		ui.Column(c).Padding(16).Height(300).Children(func() {
			ToolRegistryPanel(c, *s, tools, ToolRegistryPanelOptions{}).Element.Grow(1)
		})
	}, func(tt *ui.Tester) { _ = tt.Click("Search the web") })
}
