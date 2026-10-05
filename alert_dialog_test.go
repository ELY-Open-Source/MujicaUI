package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

type alertApp struct {
	open                 bool
	opts                 AlertDialogOptions
	confirmed, cancelled int
}

func (a *alertApp) view(c *ui.Context) {
	Use(c, Settings{})
	if Button(c, "Remove", ButtonOptions{}).Clicked() {
		a.open = true
	}
	r := AlertDialog(c, &a.open, a.opts)
	if r.Confirmed {
		a.confirmed++
	}
	if r.Cancelled {
		a.cancelled++
	}
}

func newAlertApp(opts AlertDialogOptions) (*alertApp, *ui.Tester) {
	opts.Title = "Delete the archive?"
	a := &alertApp{opts: opts}
	tt := ui.NewTester(a.view, 600, 400)
	_ = tt.Click("Remove")
	return a, tt
}

func TestAlertDialogFocusesCancelAndEnterCancels(t *testing.T) {
	a, tt := newAlertApp(AlertDialogOptions{Description: "This cannot be undone."})
	if !a.open || !tt.Focused("Cancel") {
		t.Fatalf("open %v, Cancel focused %v", a.open, tt.Focused("Cancel"))
	}
	tt.Key(0, ui.KeyEnter)
	if a.confirmed != 0 || a.cancelled != 1 || a.open {
		t.Fatalf("Enter: confirmed %d cancelled %d open %v", a.confirmed, a.cancelled, a.open)
	}
	if !tt.Focused("Remove") {
		t.Fatal("focus not restored")
	}
}

func TestAlertDialogIgnoresBackdropEscapeCancels(t *testing.T) {
	a, tt := newAlertApp(AlertDialogOptions{})
	tt.ClickAt(8, 392)
	if !a.open || a.cancelled != 0 {
		t.Fatalf("backdrop click: open %v cancelled %d", a.open, a.cancelled)
	}
	tt.Key(0, ui.KeyEscape)
	if a.open || a.cancelled != 1 || a.confirmed != 0 {
		t.Fatalf("Escape: open %v cancelled %d confirmed %d", a.open, a.cancelled, a.confirmed)
	}
}

func TestAlertDialogConfirmOnce(t *testing.T) {
	a, tt := newAlertApp(AlertDialogOptions{ConfirmLabel: "Erase"})
	if err := tt.Click("Erase"); err != nil {
		t.Fatal(err)
	}
	if a.confirmed != 1 || a.cancelled != 0 || a.open {
		t.Fatalf("confirmed %d cancelled %d open %v", a.confirmed, a.cancelled, a.open)
	}
}

func TestAlertDialogConfirmTextGatesEnter(t *testing.T) {
	a, tt := newAlertApp(AlertDialogOptions{ConfirmText: "archive"})
	_ = tt.Click("Delete") // locked
	if a.confirmed != 0 || !a.open {
		t.Fatal("locked Delete confirmed")
	}
	if err := tt.Click("Type “archive” to confirm"); err != nil {
		t.Fatal(err)
	}
	tt.Type("archiv")
	tt.Key(0, ui.KeyEnter)
	if a.confirmed != 0 || !a.open {
		t.Fatalf("before the target matched: confirmed %d open %v", a.confirmed, a.open)
	}
	tt.Type("e")
	tt.Key(0, ui.KeyEnter)
	if a.confirmed != 0 || !a.open {
		t.Fatalf("Enter in the field confirmed: %d", a.confirmed)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Delete") {
		t.Fatal("Delete not reachable once unlocked")
	}
	tt.Key(0, ui.KeyEnter)
	if a.confirmed != 1 || a.open {
		t.Fatalf("typed target: confirmed %d open %v", a.confirmed, a.open)
	}
}

// MyGo commits a composition on Enter and then handles the key (ui/editor.go:688).
func TestAlertDialogEnterDuringCompositionDoesNotConfirm(t *testing.T) {
	a, tt := newAlertApp(AlertDialogOptions{ConfirmText: "archive"})
	if err := tt.Click("Type “archive” to confirm"); err != nil {
		t.Fatal(err)
	}
	tt.Compose("archive", 7)
	tt.Key(0, ui.KeyEnter)
	if a.confirmed != 0 || !a.open {
		t.Fatalf("Enter while composing confirmed: %d", a.confirmed)
	}
}

func TestAlertDialogSnapshot(t *testing.T) {
	snapshot(t, "alert_dialog", 560, 320, func(c *ui.Context) {
		open := true
		AlertDialog(c, &open, AlertDialogOptions{
			Title: "Burn the old letters?", Description: "Twelve letters will be destroyed. This cannot be undone.",
			ConfirmText: "letters",
		})
	})
}
