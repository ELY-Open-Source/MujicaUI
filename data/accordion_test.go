package data

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func accordionHarness(open *[]string, opts *AccordionOptions, changes *int) func(c *ui.Context) {
	return func(c *ui.Context) {
		core.Use(c, core.Settings{})
		items := []AccordionItem{
			{ID: "a", Title: "Choir", Content: func() { ui.Text(c, "Choir body") }},
			{ID: "b", Title: "Nave", Content: func() { ui.Text(c, "Nave body") }},
			{ID: "x", Title: "Locked", Disabled: true, Content: func() { ui.Text(c, "Locked body") }},
			{ID: "c", Title: "Crypt", Content: func() { ui.Text(c, "Crypt body") }},
		}
		ui.Column(c).Padding(10).Children(func() {
			if Accordion(c, open, items, *opts).Changed() {
				*changes++
			}
		})
	}
}

func TestAccordionSingle(t *testing.T) {
	var open []string
	var opts AccordionOptions
	changes := 0
	tt := ui.NewTester(accordionHarness(&open, &opts, &changes), 300, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	_ = tt.Click("Choir")
	_ = tt.Click("Nave")
	if fmt.Sprint(open) != "[b]" || changes != 2 || tt.HasText("Choir body") || !tt.HasText("Nave body") {
		t.Fatalf("open %v changes %d", open, changes)
	}
	_ = tt.Click("Nave")
	if len(open) != 0 {
		t.Fatalf("could not close all: %v", open)
	}
	_ = tt.Click("Locked")
	if len(open) != 0 || changes != 3 {
		t.Fatalf("disabled opened: %v", open)
	}
}

func TestAccordionMultipleKeepOne(t *testing.T) {
	open := []string{"a"}
	opts := AccordionOptions{Multiple: true, KeepOne: true}
	changes := 0
	tt := ui.NewTester(accordionHarness(&open, &opts, &changes), 300, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	_ = tt.Click("Crypt")
	if fmt.Sprint(open) != "[a c]" {
		t.Fatalf("multiple: %v", open)
	}
	_ = tt.Click("Choir")
	_ = tt.Click("Crypt")
	if fmt.Sprint(open) != "[c]" || changes != 2 {
		t.Fatalf("keep one: %v changes %d", open, changes)
	}
}

func TestAccordionKeyboard(t *testing.T) {
	var open []string
	var opts AccordionOptions
	changes := 0
	tt := ui.NewTester(accordionHarness(&open, &opts, &changes), 300, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Choir") {
		t.Fatal("Tab missed the first header")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Nave") {
		t.Fatal("Down missed Nave")
	}
	tt.Key(0, ui.KeyDown) // skips the disabled header
	if !tt.Focused("Crypt") {
		t.Fatal("Down did not skip the disabled header")
	}
	tt.Key(0, ui.KeyHome)
	if !tt.Focused("Choir") {
		t.Fatal("Home missed")
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if fmt.Sprint(open) != "[c]" {
		t.Fatalf("Enter on End opened %v", open)
	}
	tt.Key(0, ui.KeyUp)
	if !tt.Focused("Nave") {
		t.Fatal("Up missed Nave")
	}
}

func TestAccordionBadItems(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic on duplicate IDs")
		}
	}()
	var open []string
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Accordion(c, &open, []AccordionItem{{ID: "a", Content: func() {}}, {ID: "a", Content: func() {}}}, AccordionOptions{})
	}, 100, 100)
}

func TestAccordionSnapshot(t *testing.T) {
	var opts AccordionOptions
	changes := 0
	open := []string{"b"}
	snapshot(t, "accordion", 400, 300, accordionHarnessNoUse(&open, &opts, &changes))
}

func accordionHarnessNoUse(open *[]string, opts *AccordionOptions, changes *int) func(c *ui.Context) {
	return func(c *ui.Context) {
		items := []AccordionItem{
			{ID: "a", Title: "Choir", Content: func() { ui.Text(c, "Twelve voices, two at the altar.") }},
			{ID: "b", Title: "Nave", Content: func() { ui.Text(c, "Pews of black oak, polished for the feast.") }},
			{ID: "x", Title: "Locked", Disabled: true, Content: func() {}},
			{ID: "c", Title: "Crypt", Content: func() {}},
		}
		ui.Column(c).Padding(16).Children(func() { Accordion(c, open, items, *opts) })
	}
}
