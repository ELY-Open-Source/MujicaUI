package feedback

import (
	"strconv"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
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
	var s data.ListState[string]
	tools := toolSample()
	disabled := false
	var toggled []string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
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
		core.Use(c, core.Settings{})
		var s data.ListState[string]
		ToolRegistryPanel(c, &s, []AgentTool{{Name: "a"}, {Name: "a"}}, ToolRegistryPanelOptions{})
	}, 300, 200, 1)
}

func TestToolRegistryPanelSnapshot(t *testing.T) {
	tools := toolSample()
	snapshotAfter(t, "tool_registry_panel", 760, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() *data.ListState[string] { return &data.ListState[string]{} })
		ui.Column(c).Padding(16).Height(300).Children(func() {
			ToolRegistryPanel(c, *s, tools, ToolRegistryPanelOptions{}).Element.Grow(1)
		})
	}, func(tt *ui.Tester) { _ = tt.Click("Search the web") })
}

func TestToolRegistryPanelLargeDataKeyboardAndAncestorDisabled(t *testing.T) {
	tools := make([]AgentTool, 100000)
	for i := range tools {
		tools[i] = AgentTool{Name: strconv.Itoa(i)}
	}
	var s data.ListState[string]
	disabled := false
	toggles := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			r := ToolRegistryPanel(c, &s, tools, ToolRegistryPanelOptions{})
			r.Element.Height(300)
			if _, _, ok := r.Toggled(); ok {
				toggles++
			}
		})
	}, 700, 400)
	a, b := s.Rows.Visible()
	if b-a > 30 || len(tt.Texts()) > 150 {
		t.Fatalf("built %d rows", b-a)
	}
	_ = tt.Click("0")
	tt.Key(0, ui.KeyDown)
	if id, _ := s.Selected(); id != "1" {
		t.Fatalf("keyboard selected %q", id)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if toggles != 1 {
		t.Fatalf("keyboard toggles %d", toggles)
	}
	disabled = true
	tt.Key(0, ui.KeySpace)
	if toggles != 1 {
		t.Fatal("ancestor disabled toggled")
	}
}

func TestToolRegistryPanelLargeParameterPreview(t *testing.T) {
	tools := []AgentTool{{Name: "Tool", Params: make([]ToolParam, 100000)}}
	for i := range tools[0].Params {
		tools[0].Params[i] = ToolParam{Name: "Param " + strconv.Itoa(i), Type: "string"}
	}
	var s data.ListState[string]
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ToolRegistryPanel(c, &s, tools, ToolRegistryPanelOptions{}).Element.Height(300)
	}, 700, 400)
	if err := tt.Click("Tool"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Param 0") || len(tt.Texts()) > 150 {
		t.Fatalf("parameter preview built %d texts", len(tt.Texts()))
	}
}
