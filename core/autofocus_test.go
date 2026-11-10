package core

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestAutoFocusBuildsTheFocusedFrame(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		b := ui.ButtonBase(c).Label("Seal")
		state := "idle"
		if b.Focused() {
			state = "focused"
		}
		b.Children(func() { ui.Text(c, state) })
		AutoFocus(c, b)
	}, 200, 100)
	if !tt.Focused("Seal") || !tt.HasText("focused") {
		t.Fatalf("focused %v, texts %q: the view never saw its focus", tt.Focused("Seal"), tt.Texts())
	}
	if tt.FramePending() {
		t.Error("still asking for frames")
	}
}
