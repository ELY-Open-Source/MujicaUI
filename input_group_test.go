package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

func TestInputGroupTypingAndButtons(t *testing.T) {
	value, clicked, submitted := "", -1, 0
	opts := InputGroupOptions{
		Prefix: "https://", Suffix: ".com", Label: "Site",
		Buttons: []InputGroupButton{{Label: "Clear", Icon: icons.Must("x"), IconOnly: true}, {Label: "Go"}},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		p := InputGroup(c, &value, opts)
		if p.Clicked >= 0 {
			clicked = p.Clicked
			if p.Clicked == 0 {
				value = ""
			}
		}
		if p.Submitted() {
			submitted++
		}
	}, 500, 120)
	if err := tt.Click("Site"); err != nil {
		t.Fatal(err)
	}
	tt.Type("mujica")
	tt.Key(0, ui.KeyEnter)
	if value != "mujica" || submitted != 1 {
		t.Fatalf("value %q submitted %d", value, submitted)
	}
	if err := tt.Click("Go"); err != nil {
		t.Fatal(err)
	}
	if clicked != 1 {
		t.Fatalf("clicked %d, want 1", clicked)
	}
	if err := tt.Click("Clear"); err != nil {
		t.Fatal(err)
	}
	if clicked != 0 || value != "" {
		t.Fatalf("clear: clicked %d value %q", clicked, value)
	}
	if !tt.HasText("https://") || !tt.HasText(".com") {
		t.Fatal("addons missing")
	}
}

func TestInputGroupTabOrder(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		InputGroup(c, &value, InputGroupOptions{Label: "Query", Buttons: []InputGroupButton{{Label: "Back", Leading: true}, {Label: "Search"}}})
	}, 500, 120)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Back") {
		t.Fatal("leading button is not first")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Query") {
		t.Fatal("input is not second")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Search") {
		t.Fatal("trailing button is not third")
	}
}

func TestInputGroupBlockedStates(t *testing.T) {
	for _, opts := range []InputGroupOptions{{Label: "Code", Disabled: true}, {Label: "Code", ReadOnly: true}} {
		value, clicked := "abc", -1
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			opts.Buttons = []InputGroupButton{{Label: "Copy"}}
			if p := InputGroup(c, &value, opts); p.Clicked >= 0 {
				clicked = p.Clicked
			}
		}, 400, 100)
		_ = tt.Click("Code")
		tt.Type("zz")
		_ = tt.Click("Copy")
		if value != "abc" {
			t.Fatalf("%+v: value changed to %q", opts, value)
		}
		if opts.Disabled && clicked != -1 {
			t.Fatal("disabled button clicked")
		}
		if opts.ReadOnly && clicked != 0 {
			t.Fatal("read-only group blocked its button")
		}
	}
}

func TestInputGroupSnapshot(t *testing.T) {
	a, b, e := "mujica", "", "bad value"
	snapshot(t, "input_group", 520, 200, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).Children(func() {
			InputGroup(c, &a, InputGroupOptions{Prefix: "https://", Suffix: ".com", Label: "Site"})
			InputGroup(c, &b, InputGroupOptions{Placeholder: "Search the archive", Label: "Search",
				Buttons: []InputGroupButton{{Label: "Search", Icon: icons.Must("search"), IconOnly: true, Leading: true}, {Label: "Go"}}})
			InputGroup(c, &e, InputGroupOptions{Invalid: true, Suffix: "kg", Label: "Weight"})
		})
	})
}

// Input queued while focused must not edit a group disabled since.
func TestInputGroupDisabledAfterFocus(t *testing.T) {
	value, disabled, changes := "abc", false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			if InputGroup(c, &value, InputGroupOptions{Label: "Code"}).Changed() {
				changes++
			}
		})
	}, 400, 100)
	_ = tt.Click("Code")
	disabled = true
	tt.Type("zz")
	tt.Key(0, ui.KeyBackspace)
	tt.Frame()
	if value != "abc" || changes != 0 {
		t.Fatalf("value %q changes %d", value, changes)
	}
}

func TestInputGroupHitAreasCompact(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Density: theme.Compact})
		InputGroup(c, &value, InputGroupOptions{Label: "Query", Buttons: []InputGroupButton{{Label: "Clear", Icon: icons.Must("x"), IconOnly: true}}})
	}, 400, 100)
	if r, ok := tt.Find("Clear"); !ok || r.W < 28 || r.H < 28 {
		t.Fatalf("button %v", r)
	}
}

// Cut and paste after disabling a field with a selection change neither value nor clipboard.
func TestInputGroupDisabledEditCommands(t *testing.T) {
	value, disabled, changes := "secret", false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if InputGroup(c, &value, InputGroupOptions{Label: "Code", Disabled: disabled}).Changed() {
			changes++
		}
	}, 400, 100)
	_ = tt.Click("Code")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.SetClipboard("clip")
	disabled = true
	tt.Frame()
	tt.Key(ui.Cmd, ui.KeyX)
	tt.Command("cut")
	tt.Command("paste")
	tt.Command("undo")
	tt.Frame()
	if value != "secret" || tt.Clipboard() != "clip" || changes != 0 {
		t.Fatalf("value %q clipboard %q changes %d", value, tt.Clipboard(), changes)
	}
	disabled = false
	tt.Frame()
	tt.Frame()
	if value != "secret" || tt.Clipboard() != "clip" || changes != 0 {
		t.Fatalf("after enabling: value %q clipboard %q", value, tt.Clipboard())
	}
}
