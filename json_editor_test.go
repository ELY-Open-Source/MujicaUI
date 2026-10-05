package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func jsonEditorTester(text *string, opts JsonEditorOptions, last *JsonEditorResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		*last = JsonEditor(c, text, opts)
	}, 600, 460)
}

func TestJsonEditorFormatAndValid(t *testing.T) {
	text := `{"hall":"east","guests":[1,2]}`
	var r JsonEditorResult
	tt := jsonEditorTester(&text, JsonEditorOptions{}, &r)
	if !r.Valid || !tt.HasText("Valid JSON") {
		t.Fatalf("valid: %+v", r)
	}
	_ = tt.Click("Format")
	if text != "{\n  \"hall\": \"east\",\n  \"guests\": [\n    1,\n    2\n  ]\n}" {
		t.Fatalf("formatted %q", text)
	}
}

func TestJsonEditorErrorLocation(t *testing.T) {
	text := "{\n  \"hall\": \"east\",\n  \"guests\": [1, 2,]\n}"
	var r JsonEditorResult
	tt := jsonEditorTester(&text, JsonEditorOptions{}, &r)
	if r.Valid || r.Line != 3 || r.Column != 19 {
		t.Fatalf("location %+v", r)
	}
	if !tt.HasText("   3 | ") || !tt.HasText("                         ^") {
		t.Fatalf("excerpt: %q", tt.Texts())
	}
	_ = tt.Click("Format") // disabled while invalid
	if r.Changed {
		t.Fatal("formatted invalid JSON")
	}
	text = "{} {}"
	tt.Frame()
	if r.Valid || r.Line != 1 {
		t.Fatalf("trailing value: %+v", r)
	}
}

func TestJsonEditorTypingAndDisabled(t *testing.T) {
	text := ""
	var r JsonEditorResult
	tt := jsonEditorTester(&text, JsonEditorOptions{}, &r)
	if r.Valid {
		t.Fatal("empty text is valid")
	}
	tt.Key(0, ui.KeyTab) // Format is disabled while invalid
	tt.Type(`[true]`)
	if !r.Valid || text != "[true]" {
		t.Fatalf("typed %q valid %v", text, r.Valid)
	}
	tt = jsonEditorTester(&text, JsonEditorOptions{Disabled: true}, &r)
	tt.Key(0, ui.KeyTab)
	tt.Type("x")
	if text != "[true]" {
		t.Fatalf("disabled edited: %q", text)
	}
}

func TestJsonEditorSnapshot(t *testing.T) {
	snapshot(t, "json_editor", 560, 360, func(c *ui.Context) {
		text := "{\n  \"hall\": \"east wing\",\n  \"guests\": [1, 2,]\n}"
		ui.Column(c).Padding(16).Children(func() { JsonEditor(c, &text, JsonEditorOptions{MinHeight: 120}) })
	})
}
