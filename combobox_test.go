package mujica

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

var comboboxTestFonts = []string{"Garamond", "Georgia", "Baskerville", "Didot", "Bodoni"}

type comboboxRun struct {
	value, query string
	changes      int
	queries      []string
}

func comboboxTester(r *comboboxRun, options func(query string) []string, opts ComboboxOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(10).Gap(10).Children(func() {
			res := Combobox(c, &r.value, options(r.query), opts)
			if res.Changed() {
				r.changes++
			}
			if res.QueryChanged() {
				r.queries = append(r.queries, res.Query())
			}
			r.query = res.Query()
			Button(c, "Other", ButtonOptions{})
		})
	}, 400, 400)
}

func fixedFonts(string) []string { return comboboxTestFonts }

func TestComboboxFiltersAndChooses(t *testing.T) {
	r := &comboboxRun{}
	tt := comboboxTester(r, fixedFonts, ComboboxOptions{Label: "Font"})
	tt.Key(0, ui.KeyTab)
	tt.Type("bo")
	if !tt.HasText("Bodoni") || tt.HasText("Garamond") {
		t.Fatalf("shown %v", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if r.value != "Bodoni" || r.changes != 1 {
		t.Fatalf("value %q changes %d", r.value, r.changes)
	}
}

func TestComboboxArrowsAndClick(t *testing.T) {
	r := &comboboxRun{}
	tt := comboboxTester(r, fixedFonts, ComboboxOptions{Label: "Font"})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if r.value != "Georgia" {
		t.Fatalf("arrows chose %q", r.value)
	}
	tt.Key(0, ui.KeyDown)
	if err := tt.Click("Didot"); err != nil {
		t.Fatal(err)
	}
	if r.value != "Didot" || r.query != "Didot" {
		t.Fatalf("click chose %q, field %q", r.value, r.query)
	}
}

func TestComboboxRevertsUnknownText(t *testing.T) {
	r := &comboboxRun{value: "Didot"}
	tt := comboboxTester(r, fixedFonts, ComboboxOptions{})
	tt.Key(0, ui.KeyTab)
	tt.Type("zzz")
	if !tt.HasText("No matches") {
		t.Fatal("no-match state missing")
	}
	tt.Key(0, ui.KeyTab)
	if r.value != "Didot" || r.query != "Didot" || r.changes != 0 {
		t.Fatalf("value %q field %q changes %d", r.value, r.query, r.changes)
	}
}

func TestComboboxCustomValue(t *testing.T) {
	r := &comboboxRun{}
	tt := comboboxTester(r, fixedFonts, ComboboxOptions{AllowCustom: true})
	tt.Key(0, ui.KeyTab)
	tt.Type("Caslon")
	if !tt.HasText("Use “Caslon”") {
		t.Fatalf("custom option missing: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if r.value != "Caslon" {
		t.Fatalf("value %q", r.value)
	}
	tt.Type(" Pro")
	tt.Key(0, ui.KeyTab)
	if r.value != "Caslon Pro" {
		t.Fatalf("blur kept %q", r.value)
	}
}

func TestComboboxAsyncResults(t *testing.T) {
	r := &comboboxRun{}
	loading := true
	remote := func(q string) []string {
		if loading || q == "" {
			return nil
		}
		var out []string
		for _, f := range comboboxTestFonts {
			if strings.Contains(strings.ToLower(f), strings.ToLower(q)) {
				out = append(out, f+" (remote)")
			}
		}
		return out
	}
	opts := ComboboxOptions{ExternalFilter: true}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		opts.Loading = loading
		res := Combobox(c, &r.value, remote(r.query), opts)
		if res.QueryChanged() {
			r.queries = append(r.queries, res.Query())
		}
		r.query = res.Query()
	}, 400, 400)
	tt.Key(0, ui.KeyTab)
	tt.Type("di")
	if len(r.queries) != 1 || r.queries[0] != "di" {
		t.Fatalf("queries %v", r.queries)
	}
	if !tt.HasText("Loading…") {
		t.Fatal("loading state missing")
	}
	loading = false
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if r.value != "Didot (remote)" {
		t.Fatalf("value %q", r.value)
	}
}

func TestComboboxDisabled(t *testing.T) {
	r := &comboboxRun{value: "Didot"}
	tt := comboboxTester(r, fixedFonts, ComboboxOptions{Disabled: true})
	_ = tt.Click("Didot")
	tt.Type("x")
	if r.value != "Didot" || r.query != "Didot" || tt.HasText("Bodoni") {
		t.Fatalf("disabled edited: %q %q", r.value, r.query)
	}
}

func TestComboboxSnapshot(t *testing.T) {
	snapshot(t, "combobox", 460, 110, func(c *ui.Context) {
		a, b := "", "Didot"
		ui.Row(c).Padding(20).Gap(16).Children(func() {
			Combobox(c, &a, comboboxTestFonts, ComboboxOptions{Placeholder: "Typeface"})
			Combobox(c, &b, comboboxTestFonts, ComboboxOptions{Loading: true})
		})
	})
}
