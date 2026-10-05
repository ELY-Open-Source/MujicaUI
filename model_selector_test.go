package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var modelsSample = []ModelOption{
	{ID: "herald", Name: "Herald 2", Capabilities: []string{"Vision", "Tools"}, ContextTokens: 200000},
	{ID: "page", Name: "Page Mini", Capabilities: []string{"Fast"}, ContextTokens: 32000},
	{ID: "oracle", Name: "Oracle", ContextTokens: 1000000, Unavailable: "Needs the Court plan"},
}

func modelTester(v *string, changes *int, off *bool) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
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
