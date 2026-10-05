package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
)

func TestTextInputTypingAndSubmit(t *testing.T) {
	name, changes, submits := "", 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := TextInput(c, &name, TextInputOptions{Label: "Name", Placeholder: "Your name"})
		if r.Changed() {
			changes++
		}
		if r.Submitted() {
			submits++
		}
	}, 400, 200)
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Ada")
	tt.Key(0, ui.KeyEnter)
	if name != "Ada" || changes == 0 || submits != 1 {
		t.Fatalf("name %q changes %d submits %d", name, changes, submits)
	}
}

func TestTextInputUndoClipboardIME(t *testing.T) {
	s := ""
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TextInput(c, &s, TextInputOptions{Label: "Note"})
	}, 400, 200)
	_ = tt.Click("Note")
	tt.Compose("にほ", 2)
	if s != "" {
		t.Fatalf("composition leaked into value: %q", s)
	}
	tt.Type("日本")
	tt.Command("selectAll")
	tt.Command("copy")
	if tt.Clipboard() != "日本" {
		t.Fatalf("clipboard %q", tt.Clipboard())
	}
	tt.Type("x")
	tt.Command("undo")
	if s != "日本" {
		t.Fatalf("after undo %q", s)
	}
}

func TestTextInputClear(t *testing.T) {
	s, changes := "draft", 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if TextInput(c, &s, TextInputOptions{Label: "Draft", Clearable: true}).Changed() {
			changes++
		}
	}, 400, 200)
	if err := tt.Click("Clear"); err != nil {
		t.Fatal(err)
	}
	if s != "" || changes != 1 {
		t.Fatalf("after clear %q, changes %d", s, changes)
	}
	if _, ok := tt.Find("Clear"); ok {
		t.Fatal("clear button stays on an empty field")
	}
	tt.Type("again")
	tt.Key(0, ui.KeyEscape)
	if s != "" || changes != 3 {
		t.Fatalf("after Escape %q, changes %d", s, changes)
	}
}

func TestTextInputReadOnlyAndDisabled(t *testing.T) {
	ro, off := "fixed", "locked"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Gap(8).Children(func() {
			TextInput(c, &ro, TextInputOptions{Label: "RO", ReadOnly: true, Clearable: true})
			TextInput(c, &off, TextInputOptions{Label: "Off", Disabled: true})
		})
	}, 400, 200)
	_ = tt.Click("RO")
	tt.Command("selectAll")
	tt.Command("copy")
	tt.Type("zz")
	tt.Key(0, ui.KeyBackspace)
	if ro != "fixed" || tt.Clipboard() != "fixed" {
		t.Fatalf("read-only value %q clipboard %q", ro, tt.Clipboard())
	}
	if _, ok := tt.Find("Clear"); ok {
		t.Fatal("read-only field offers clear")
	}
	_ = tt.Click("Off")
	tt.Type("zz")
	if off != "locked" || tt.Focused("Off") {
		t.Fatalf("disabled value %q focused %v", off, tt.Focused("Off"))
	}
}

func TestTextInputSlotsAndError(t *testing.T) {
	s := "x"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TextInput(c, &s, TextInputOptions{
			Label:    "Price",
			Leading:  func() { ui.Text(c, "¥") },
			Trailing: func() { ui.Text(c, "CNY") },
			Error:    "Not a price",
		})
	}, 400, 200)
	for _, want := range []string{"¥", "CNY", "Not a price"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
}

func TestTextInputSnapshot(t *testing.T) {
	view := func(c *ui.Context) {
		a, b, e, ro, dis := "", "Elisabeth", "1789-13-01", "Read only", "Disabled"
		ui.Column(c).Padding(24).Gap(12).Width(360).Children(func() {
			TextInput(c, &a, TextInputOptions{Placeholder: "Placeholder", Label: "A"})
			TextInput(c, &b, TextInputOptions{Label: "B", Clearable: true, Leading: func() {
				ui.Icon(c, icons.Must("user")).Size(16, 16).TextColor(core.Tokens(c).TextMuted)
			}})
			TextInput(c, &e, TextInputOptions{Label: "E", Error: "No such date"})
			TextInput(c, &ro, TextInputOptions{Label: "R", ReadOnly: true})
			TextInput(c, &dis, TextInputOptions{Label: "D", Disabled: true})
		})
	}
	snapshot(t, "text_input", 420, 320, view)
	buttonDensitySnapshots(t, "text_input", 420, 360, view)
}
