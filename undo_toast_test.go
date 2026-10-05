package mujica

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func undoTester(open *bool, life time.Duration, undos *int, last *UndoToastResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if Button(c, "Delete", ButtonOptions{}).Clicked() {
			*open = true
		}
		if r := UndoToast(c, open, "File deleted", UndoToastOptions{Duration: life, OnUndo: func() { *undos++ }}); r.Undone || r.Expired {
			*last = r
		}
	}, 600, 400)
}

func TestUndoToastUndoAndFocus(t *testing.T) {
	open, undos := false, 0
	var last UndoToastResult
	tt := undoTester(&open, 0, &undos, &last)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("File deleted") || !tt.HasText("Undo (5)") || !tt.Focused("Delete") {
		t.Fatalf("shown %q, focus kept %v", tt.Texts(), tt.Focused("Delete"))
	}
	if !slices.Contains(tt.Announcements(), "File deleted") {
		t.Fatal("not announced")
	}
	if err := tt.Click("Undo (5)"); err != nil {
		t.Fatal(err)
	}
	if undos != 1 || open || !last.Undone || tt.HasText("File deleted") {
		t.Fatalf("undo: %d %v %+v", undos, open, last)
	}
}

func TestUndoToastCountdownPauseAndClose(t *testing.T) {
	open, undos := false, 0
	var last UndoToastResult
	tt := undoTester(&open, 60*time.Millisecond, &undos, &last)
	_ = tt.Click("Delete")
	r, _ := tt.Find("File deleted")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	time.Sleep(90 * time.Millisecond)
	tt.Frame()
	if !open {
		t.Fatal("hovered countdown ran out")
	}
	tt.Move(5, 5)
	time.Sleep(90 * time.Millisecond)
	tt.Frame()
	if open || !last.Expired || undos != 0 {
		t.Fatalf("expiry: %v %+v %d", open, last, undos)
	}
	last = UndoToastResult{}
	_ = tt.Click("Delete")
	_ = tt.Click("Dismiss notification")
	if open || !last.Expired || undos != 0 {
		t.Fatalf("close: %v %+v", open, last)
	}
}

func TestUndoToastRejects(t *testing.T) {
	open := true
	for name, f := range map[string]func(c *ui.Context){
		"no undo":  func(c *ui.Context) { UndoToast(c, &open, "x", UndoToastOptions{}) },
		"negative": func(c *ui.Context) { UndoToast(c, &open, "x", UndoToastOptions{OnUndo: func() {}, Duration: -1}) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); f(c) }, 300, 200)
		}()
	}
}

func TestUndoToastSnapshot(t *testing.T) {
	open := true
	snapshot(t, "undo_toast", 520, 160, func(c *ui.Context) {
		UndoToast(c, &open, "3 messages archived", UndoToastOptions{OnUndo: func() {}, Duration: time.Hour})
	})
}
