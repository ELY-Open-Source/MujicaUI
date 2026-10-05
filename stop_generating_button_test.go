package mujica

import (
	"runtime"
	"testing"

	"github.com/egoist/mygo/ui"
)

func chatMod() ui.Modifiers {
	if runtime.GOOS == "darwin" {
		return ui.Cmd
	}
	return ui.Ctrl
}

func TestStopGeneratingButtonClickAndShortcut(t *testing.T) {
	generating, stops := true, 0
	tt := chatTester(400, 120, func(c *ui.Context) {
		if StopGeneratingButton(c, StopGeneratingButtonOptions{Generating: generating}).Changed() {
			stops++
			generating = false
		}
	})
	if err := tt.Click("Stop generating"); err != nil {
		t.Fatal(err)
	}
	if stops != 1 || tt.HasText("Stop generating") {
		t.Fatalf("stops %d, still shown %v", stops, tt.HasText("Stop generating"))
	}
	generating = true
	tt.Frame()
	tt.Key(chatMod(), ui.KeyPeriod)
	if stops != 2 {
		t.Fatalf("shortcut stops %d", stops)
	}
	tt.Key(chatMod(), ui.KeyPeriod)
	if stops != 2 {
		t.Fatal("shortcut acted while hidden")
	}
}

func TestStopGeneratingButtonDisabled(t *testing.T) {
	stops := 0
	tt := chatTester(400, 120, func(c *ui.Context) {
		if StopGeneratingButton(c, StopGeneratingButtonOptions{Generating: true, Disabled: true}).Changed() {
			stops++
		}
	})
	_ = tt.Click("Stop generating")
	tt.Key(chatMod(), ui.KeyPeriod)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if stops != 0 {
		t.Fatalf("disabled button stopped %d times", stops)
	}
}

func TestStopGeneratingButtonSnapshot(t *testing.T) {
	snapshot(t, "stop_generating_button", 300, 70, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			StopGeneratingButton(c, StopGeneratingButtonOptions{Generating: true})
		})
	})
}
