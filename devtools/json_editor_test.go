package devtools

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func jsonEditorTester(text *string, opts JsonEditorOptions, last *JsonEditorResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
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

func TestJsonEditorRejectsTrailingDelimiter(t *testing.T) {
	for _, text := range []string{"{} ]", "1 }", "true ]"} {
		var r JsonEditorResult
		jsonEditorTester(&text, JsonEditorOptions{}, &r)
		if r.Valid {
			t.Fatalf("accepted trailing delimiter: %q", text)
		}
	}
}

func TestJsonEditorHeightBounds(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("reversed heights accepted")
		}
	}()
	text := "{}"
	var r JsonEditorResult
	jsonEditorTester(&text, JsonEditorOptions{MinHeight: 300, MaxHeight: 100}, &r)
}

func TestJsonEditorAncestorDisabledFormat(t *testing.T) {
	text := `{"a":1}`
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() { JsonEditor(c, &text, JsonEditorOptions{}) })
	}, 600, 460)
	tt.Key(0, ui.KeyTab)
	disabled = true
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("JSON")
	tt.Type("x")
	if text != `{"a":1}` {
		t.Fatal("disabled editor changed")
	}
}

func TestJsonEditorTrailingValueLocation(t *testing.T) {
	text := "{}\n{}"
	var result JsonEditorResult
	jsonEditorTester(&text, JsonEditorOptions{}, &result)
	if result.Valid || result.Line != 2 || result.Column != 1 {
		t.Fatalf("trailing value at %d:%d", result.Line, result.Column)
	}
}

func TestJsonEditorPreservesLargeNumbers(t *testing.T) {
	text := `{"value":1e309}`
	var result JsonEditorResult
	tt := jsonEditorTester(&text, JsonEditorOptions{}, &result)
	if !result.Valid {
		t.Fatal("valid large JSON number rejected")
	}
	_ = tt.Click("Format")
	if text != "{\n  \"value\": 1e309\n}" {
		t.Fatal("large number changed by formatting")
	}
}
