package mujica

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// overlayLongText is about 372 characters of description.
var overlayLongText = strings.Repeat("The east wing stays sealed until the court returns. ", 7)

// overlayInWindow fails unless every label lies fully inside a w×h window.
func overlayInWindow(t *testing.T, tt *ui.Tester, w, h float32, labels ...string) {
	t.Helper()
	for _, l := range labels {
		r, ok := tt.Find(l)
		if !ok || r.W == 0 || r.H == 0 || r.Y < 0 || r.X < 0 || r.Y+r.H > h || r.X+r.W > w {
			t.Errorf("%q at %+v (found %v) outside the %v×%v window", l, r, ok, w, h)
		}
	}
}

func TestDialogActionsStayInWindowAtLargeText(t *testing.T) {
	open, name := true, "Ophelia"
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Dialog(c, &open, DialogOptions{Title: "Rename", Description: overlayLongText, Actions: func() {
			Button(c, "Cancel", ButtonOptions{Variant: Secondary})
			Button(c, "Save", ButtonOptions{})
		}}, func() { ui.TextInput(c, &name) })
	}, 800, 600)
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	tt.Frame()
	overlayInWindow(t, tt, 800, 600, "Cancel", "Save")
}

func TestAlertDialogActionsStayInWindowAtLargeText(t *testing.T) {
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		AlertDialog(c, &open, AlertDialogOptions{Title: "Burn?", Description: overlayLongText, ConfirmText: "letters"})
	}, 800, 600)
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	tt.Frame()
	overlayInWindow(t, tt, 800, 600, "Cancel", "Delete", "Type “letters” to confirm")
}

func TestOverlayPanelFitsWindow(t *testing.T) {
	for _, modal := range []bool{true, false} {
		open := true
		var panel ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			if p := Overlay(c, &open, OverlayOptions{Modal: modal}, func() {
				ui.Text(c, strings.Repeat(overlayLongText, 6))
				Button(c, "Done", ButtonOptions{})
			}); p != nil {
				panel = p.Bounds()
			}
		}, 600, 400)
		tt.Frame()
		if panel.X < 16 || panel.Y < 16 || panel.X+panel.W > 584 || panel.Y+panel.H > 384 {
			t.Errorf("modal %v: panel %+v outside the window minus Inset", modal, panel)
		}
	}
}
