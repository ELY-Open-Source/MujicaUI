package input

import (
	"bytes"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

func TestButtonLoadingKeepsWidth(t *testing.T) {
	for _, icon := range []*ui.SVG{nil, icons.Must("check")} {
		loading := false
		var w float32
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			ui.Row(c).Children(func() {
				w = Button(c, "Save", ButtonOptions{Icon: icon, Loading: loading}).Bounds().W
			})
		}, 300, 200)
		tt.Frame()
		before := w
		loading = true
		tt.Frame()
		tt.Frame()
		if w != before || !tt.HasText("Save") {
			t.Fatalf("icon %v: width %v -> %v, texts %q", icon != nil, before, w, tt.Texts())
		}
	}
}

func inputReviewPanics(t *testing.T, name string, view func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: no panic", name)
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		view(c)
	}, 200, 100)
}

func TestIconOnlyButtonsNeedNames(t *testing.T) {
	inputReviewPanics(t, "Button", func(c *ui.Context) {
		Button(c, "", ButtonOptions{Icon: icons.Must("x")})
	})
	inputReviewPanics(t, "ButtonGroup", func(c *ui.Context) {
		ButtonGroup(c, []ButtonGroupItem{{Icon: icons.Must("bold")}}, ButtonGroupOptions{})
	})
}

func TestTextInputDisabledWhileFocused(t *testing.T) {
	s, off, focused, changes := "", false, false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := TextInput(c, &s, TextInputOptions{Label: "Name", Disabled: off})
		focused = r.Input.Focused()
		if r.Changed() || r.Submitted() {
			changes++
		}
	}, 400, 200)
	_ = tt.Click("Name")
	tt.Type("a")
	off = true
	tt.Type("b")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if s != "a" || changes != 1 || focused {
		t.Fatalf("value %q events %d focused %v", s, changes, focused)
	}
}

func TestTextAreaDisabledWhileFocused(t *testing.T) {
	s, off := "", false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TextArea(c, &s, TextAreaOptions{Label: "Notes", Disabled: off})
	}, 400, 300)
	_ = tt.Click("Notes")
	off = true
	tt.Type("late")
	if s != "" {
		t.Fatalf("disabled area took %q", s)
	}
}

func TestCheckboxDisabledWhileFocused(t *testing.T) {
	on, off, focused := false, false, false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		focused = Checkbox(c, &on, "Sync", CheckboxOptions{Disabled: off}).Focused()
	}, 300, 200)
	tt.Key(0, ui.KeyTab)
	off = true
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if on || focused {
		t.Fatalf("checked %v focused %v", on, focused)
	}
}

func TestRadioGroupDisabledWhileFocused(t *testing.T) {
	v, off := "a", false
	items := []RadioItem[string]{{Value: "a", Label: "A"}, {Value: "b", Label: "B"}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		RadioGroup(c, &v, items, RadioGroupOptions{Disabled: off})
	}, 300, 200)
	tt.Key(0, ui.KeyTab)
	off = true
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	if v != "a" {
		t.Fatalf("disabled group chose %q", v)
	}
}

func TestTextAreaPlainEnterSubmits(t *testing.T) {
	s, submits := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if TextArea(c, &s, TextAreaOptions{Label: "Chat", SubmitKey: ui.KeyEnter}).Submitted() {
			submits++
		}
	}, 400, 300)
	_ = tt.Click("Chat")
	tt.Type("hi")
	tt.Key(0, ui.KeyEnter)
	tt.Key(ui.Shift, ui.KeyEnter)
	if submits != 1 || s != "hi\n" {
		t.Fatalf("submits %d value %q", submits, s)
	}
}

func TestRadioGroupFocusFollowsValue(t *testing.T) {
	v := ""
	items := []RadioItem[string]{{Value: "a", Label: "A"}, {Value: "b", Label: "B"}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		RadioGroup(c, &v, items, RadioGroupOptions{})
	}, 300, 200)
	tt.Key(0, ui.KeyTab)
	items[0], items[1] = items[1], items[0]
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	if v != "a" {
		t.Fatalf("Space after reordering chose %q, want a", v)
	}
}

// inputReviewStates renders view plain, hovered and pressed at the named element.
func inputReviewStates(t *testing.T, name string, view func(c *ui.Context)) (plain, hover, pressed []byte) {
	t.Helper()
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Dark})
		view(c)
	}, 300, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("no %q", name)
	}
	plain = encode(t, tt.Image())
	// Colors still fade under reduced motion; let each fade finish on the Tester clock.
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(theme.StateDuration)
	hover = encode(t, tt.Image())
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(theme.StateDuration)
	pressed = encode(t, tt.Image())
	return
}

func TestCheckedMarksShowHoverAndPress(t *testing.T) {
	on, v := true, "a"
	for name, view := range map[string]func(c *ui.Context){
		"Sync": func(c *ui.Context) { Checkbox(c, &on, "Sync", CheckboxOptions{}) },
		"A": func(c *ui.Context) {
			RadioGroup(c, &v, []RadioItem[string]{{Value: "a", Label: "A"}}, RadioGroupOptions{})
		},
	} {
		plain, hover, pressed := inputReviewStates(t, name, view)
		if bytes.Equal(plain, hover) || bytes.Equal(hover, pressed) {
			t.Errorf("%s: checked state hides hover (%v) or press (%v)", name, bytes.Equal(plain, hover), bytes.Equal(hover, pressed))
		}
	}
}
