package chat

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var modelsSample = []ModelOption{
	{ID: "herald", Name: "Herald 2", Capabilities: []string{"Vision", "Tools"}, ContextTokens: 200000},
	{ID: "page", Name: "Page Mini", Capabilities: []string{"Fast"}, ContextTokens: 32000},
	{ID: "oracle", Name: "Oracle", ContextTokens: 1000000, Unavailable: "Needs the Court plan"},
}

func modelTester(v *string, changes *int, off *bool) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(16).Children(func() {
			if ModelSelector(c, v, modelsSample, ModelSelectorOptions{Disabled: *off}).Changed() {
				*changes++
			}
		})
	}, 500, 500)
}

func TestModelSelectorShowsDetailsAndChooses(t *testing.T) {
	v, n, off := "herald", 0, false
	tt := modelTester(&v, &n, &off)
	_ = tt.Click("Model")
	for _, s := range []string{"200K context", "Vision", "1M context", "Needs the Court plan", "32K context"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	_ = tt.Click("Oracle")
	if v != "herald" || n != 0 {
		t.Fatalf("unavailable model chosen: %q", v)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if v != "page" || n != 1 {
		t.Fatalf("keyboard choice %q changes %d (Down stops before the unavailable model)", v, n)
	}
}

func TestModelSelectorDisabled(t *testing.T) {
	v, n, off := "herald", 0, true
	tt := modelTester(&v, &n, &off)
	_ = tt.Click("Model")
	if tt.HasText("Page Mini") {
		t.Fatal("disabled selector opened")
	}
}

func TestModelSelectorTokens(t *testing.T) {
	for n, want := range map[int]string{512: "512", 8000: "8K", 200000: "200K", 1000000: "1M", 1500000: "1.5M"} {
		if got := modelSelectorTokens(n); got != want {
			t.Errorf("%d: %q, want %q", n, got, want)
		}
	}
}

func TestModelSelectorSnapshot(t *testing.T) {
	snapshotAfter(t, "model_selector", 440, 330, func(c *ui.Context) {
		v := ui.Local(c.Root(), "v", func() string { return "herald" })
		ui.Column(c).Padding(16).Children(func() { ModelSelector(c, v, modelsSample, ModelSelectorOptions{}) })
	}, func(tt *ui.Tester) { _ = tt.Click("Model") })
}

func TestModelSelectorAncestorDisabledCloses(t *testing.T) {
	v, off, changes := "herald", false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(off).Children(func() {
			if ModelSelector(c, &v, modelsSample, ModelSelectorOptions{}).Changed() {
				changes++
			}
		})
	}, 400, 400)
	_ = tt.Click("Model")
	off = true
	tt.Frame()
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if tt.HasText("Page Mini") || v != "herald" || changes != 0 {
		t.Fatalf("disabled popup: value %q changes %d texts %q", v, changes, tt.Texts())
	}
}

func TestModelSelectorLargeDataAndEnd(t *testing.T) {
	models := make([]ModelOption, 100000)
	for i := range models {
		models[i] = ModelOption{ID: fmt.Sprint(i), Name: fmt.Sprintf("Model %d", i)}
	}
	v := "0"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ModelSelector(c, &v, models, ModelSelectorOptions{})
	}, 400, 300)
	_ = tt.Click("Model")
	if n := len(tt.Texts()); n > 200 {
		t.Errorf("built %d labels, want <= 200", n)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if v != "99999" {
		t.Errorf("End choice %q, want 99999", v)
	}
}

func TestModelSelectorPopupFitsWindow(t *testing.T) {
	value := "0"
	models := make([]ModelOption, 100)
	for i := range models {
		models[i] = ModelOption{ID: fmt.Sprint(i), Name: fmt.Sprintf("Model %d", i)}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(16).Children(func() { ModelSelector(c, &value, models, ModelSelectorOptions{}) })
	}, 440, 330)
	_ = tt.Click("Model")
	tt.Key(0, ui.KeyEnd)
	b, ok := tt.Find("Model 99")
	if !ok || b.H == 0 || b.Y < 0 || b.Y+b.H > 330 {
		t.Fatalf("last available model outside window %v found %v", b, ok)
	}
}
