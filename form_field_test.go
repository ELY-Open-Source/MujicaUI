package mujica

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestFormFieldLabelFocusesControl(t *testing.T) {
	name := ""
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		FormField(c, "Name", FormFieldOptions{Required: true, Description: "As on your seal."}, func() *ui.Element {
			return InputGroup(c, &name, InputGroupOptions{}).Input
		})
	}, 400, 200)
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Sakiko")
	if name != "Sakiko" {
		t.Fatalf("label click did not focus the control: %q", name)
	}
	if !tt.Focused("Name") || !tt.HasText("As on your seal.") || !tt.HasText("*") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

// The label goes to the returned input, not the first focusable element before it.
func TestFormFieldTargetsReturnedControl(t *testing.T) {
	tags, site := []string{"rose"}, ""
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Gap(12).Children(func() {
			FormField(c, "Flowers", FormFieldOptions{Error: "Too few."}, func() *ui.Element {
				return TagsInput(c, &tags, TagsInputOptions{}).Input
			})
			FormField(c, "Site", FormFieldOptions{}, func() *ui.Element {
				return InputGroup(c, &site, InputGroupOptions{Buttons: []InputGroupButton{{Label: "Back", Leading: true}}}).Input
			})
		})
	}, 500, 300)
	_ = tt.Click("Flowers")
	tt.Type("ivy")
	tt.Key(0, ui.KeyEnter)
	if !slices.Equal(tags, []string{"rose", "ivy"}) || !tt.Focused("Flowers") {
		t.Fatalf("tags %q", tags)
	}
	if _, ok := tt.Find("Remove rose"); !ok {
		t.Fatal("remove button lost its own name")
	}
	_ = tt.Click("Site")
	tt.Type("court")
	if site != "court" || tt.Focused("Back") {
		t.Fatalf("site %q", site)
	}
}

func TestFormFieldErrorAndLayout(t *testing.T) {
	v := ""
	horizontal := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		FormField(c, "Email", FormFieldOptions{Error: "Enter an email address.", Horizontal: horizontal}, func() *ui.Element {
			return InputGroup(c, &v, InputGroupOptions{Prefix: "mailto:"}).Input
		})
	}, 500, 200)
	if !tt.HasText("Enter an email address.") {
		t.Fatal("error not shown")
	}
	label, _ := tt.Find("Email")
	input, _ := tt.Find("mailto:")
	if !(input.Y > label.Y+label.H/2) {
		t.Fatalf("vertical: label %v input %v", label, input)
	}
	horizontal = true
	tt.Frame()
	label, _ = tt.Find("Email")
	input, _ = tt.Find("mailto:")
	if !(input.X > label.X+label.W) || input.Y > label.Y+label.H {
		t.Fatalf("horizontal: label %v input %v", label, input)
	}
}

func TestFormFieldNeedsLabel(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an empty label")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		FormField(c, "", FormFieldOptions{}, func() *ui.Element { return ui.Box(c) })
	}, 100, 100, 1)
}

func TestFormFieldSnapshot(t *testing.T) {
	a, b := "Sakiko Togawa", "uika@"
	snapshot(t, "form_field", 520, 260, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(16).Children(func() {
			FormField(c, "Name", FormFieldOptions{Required: true, Description: "As written on the invitation."}, func() *ui.Element {
				return InputGroup(c, &a, InputGroupOptions{}).Input
			})
			FormField(c, "Email", FormFieldOptions{Horizontal: true, Error: "Enter a full email address."}, func() *ui.Element {
				return InputGroup(c, &b, InputGroupOptions{Invalid: true}).Input
			})
		})
	})
}
