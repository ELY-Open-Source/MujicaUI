package account

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/input"
)

func shortcutRecorderTester(v *Shortcut, opts *ShortcutRecorderOptions, n *int) *ui.Tester {
	return accountTester(400, 160, func(c *ui.Context) {
		if ShortcutRecorder(c, v, *opts).Changed() {
			*n++
		}
	})
}

func TestShortcutRecorderRecordsAndCancels(t *testing.T) {
	v, n := Shortcut{}, 0
	opts := ShortcutRecorderOptions{Label: "Open palette"}
	tt := shortcutRecorderTester(&v, &opts, &n)
	if _, ok := tt.Find("Open palette: none"); !ok {
		t.Fatal("empty value not named")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Press keys…") {
		t.Fatalf("not recording: %q", tt.Texts())
	}
	tt.TypeKey(0, ui.KeyK, "k")
	if !tt.HasText("Add Ctrl, Alt or Command; F-keys may stand alone.") || n != 0 {
		t.Fatal("bare letter accepted")
	}
	tt.Key(ui.Ctrl|ui.Shift, ui.KeyK)
	if v != (Shortcut{ui.Ctrl | ui.Shift, ui.KeyK}) || n != 1 || tt.HasText("Press keys…") {
		t.Fatalf("recorded %+v", v)
	}
	if v.String() != "ctrl+shift+K" {
		t.Fatalf("spelled %q", v.String())
	}
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyEscape)
	tt.Key(0, ui.KeyF5)
	if v.Key != ui.KeyK || tt.HasText("Press keys…") {
		t.Fatal("escape did not cancel")
	}
}

func TestShortcutRecorderConflictAndDisabled(t *testing.T) {
	v, n := Shortcut{ui.Ctrl, ui.KeyP}, 0
	opts := ShortcutRecorderOptions{Label: "Print", Taken: []ShortcutBinding{{Action: "Save", Shortcut: Shortcut{ui.Ctrl, ui.KeyS}}}}
	tt := shortcutRecorderTester(&v, &opts, &n)
	_ = tt.Click("Print: ctrl+P")
	tt.Key(ui.Ctrl, ui.KeyS)
	if v.Key != ui.KeyP || n != 0 || !tt.HasText("ctrl+S is already used by Save") {
		t.Fatalf("conflict: %+v %q", v, tt.Texts())
	}
	tt.Key(0, ui.KeySpace)
	opts.Disabled = true
	tt.Key(ui.Ctrl, ui.KeyJ)
	if v.Key != ui.KeyP {
		t.Fatal("disabled recorder recorded")
	}
}

func TestShortcutRecorderSnapshot(t *testing.T) {
	snapshot(t, "shortcut_recorder", 300, 110, func(c *ui.Context) {
		a := ui.Local(c.Root(), "a", func() Shortcut { return Shortcut{ui.Ctrl | ui.Shift, ui.KeyK} })
		b := ui.Local(c.Root(), "b", func() Shortcut { return Shortcut{} })
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			ShortcutRecorder(c, a, ShortcutRecorderOptions{Label: "Command palette"})
			ShortcutRecorder(c, b, ShortcutRecorderOptions{Label: "Zen mode"})
		})
	})
}

func TestShortcutRecorderFunctionKeyAndTabCancellation(t *testing.T) {
	v, n := Shortcut{}, 0
	opts := ShortcutRecorderOptions{Label: "Refresh"}
	tt := accountTester(400, 200, func(c *ui.Context) {
		if ShortcutRecorder(c, &v, opts).Changed() {
			n++
		}
		input.Button(c, "Next control", input.ButtonOptions{})
	})
	if err := tt.Click("Refresh: none"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyF5)
	if v != (Shortcut{Key: ui.KeyF5}) || n != 1 {
		t.Fatal("standalone function key fails")
	}
	if err := tt.Click("Refresh: F5"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Ctrl, ui.KeyP)
	if v.Key != ui.KeyF5 || n != 1 {
		t.Fatal("recorder retains keys after Tab")
	}
	accountPanics(t, func(c *ui.Context) { ShortcutRecorder(c, &v, ShortcutRecorderOptions{}) })
}

func TestShortcutRecorderNamesEveryDeliveredKey(t *testing.T) {
	for _, key := range []ui.Key{ui.KeyInsert, ui.KeyContextMenu, ui.KeyBack, ui.KeyForward} {
		t.Run(fmt.Sprint(key), func(t *testing.T) {
			defer func() {
				if recover() != nil {
					t.Fatal("recording a valid MyGo key panicked")
				}
			}()
			v, n := Shortcut{}, 0
			opts := ShortcutRecorderOptions{Label: "Action"}
			tt := shortcutRecorderTester(&v, &opts, &n)
			if err := tt.Click("Action: none"); err != nil {
				t.Fatal(err)
			}
			tt.Key(ui.Ctrl, key)
			if v != (Shortcut{ui.Ctrl, key}) || n != 1 || v.String() == "" {
				t.Fatal("valid key was not recorded and named")
			}
		})
	}
}
