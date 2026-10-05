package messaging

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

func TestCallControlsClicksAndShortcuts(t *testing.T) {
	st := CallState{Muted: true}
	var k callCalls
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			CallControls(c, st, k.opts())
			input.Button(c, "Elsewhere", input.ButtonOptions{})
		})
	}, 500, 200)
	for _, name := range []string{"Unmute", "Turn camera off", "Share screen", "Leave"} {
		if err := tt.Click(name); err != nil {
			t.Fatal(err)
		}
	}
	if len(k.mute) != 1 || k.mute[0] || len(k.camera) != 1 || !k.camera[0] || len(k.share) != 1 || !k.share[0] || k.hang != 1 {
		t.Fatalf("clicks %+v", k)
	}
	if err := tt.Click("Elsewhere"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyD)
	tt.Key(ui.Cmd, ui.KeyE)
	tt.Key(ui.Cmd|ui.Shift, ui.KeyS)
	tt.Key(ui.Cmd|ui.Shift, ui.KeyH)
	if len(k.mute) != 2 || len(k.camera) != 2 || len(k.share) != 2 || k.hang != 2 {
		t.Fatalf("shortcuts %+v", k)
	}
}

func TestCallControlsDisabled(t *testing.T) {
	var k callCalls
	o := k.opts()
	o.Disabled = true
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CallControls(c, CallState{}, o)
	}, 500, 200)
	_ = tt.Click("Mute")
	tt.Key(ui.Cmd, ui.KeyD)
	tt.Key(ui.Cmd|ui.Shift, ui.KeyH)
	if len(k.mute) != 0 || k.hang != 0 {
		t.Fatalf("disabled controls reported %+v", k)
	}
	mediaPanics(t, "no callbacks", func(c *ui.Context) { CallControls(c, CallState{}, CallControlsOptions{}) })
}

func TestCallControlsSnapshot(t *testing.T) {
	var k callCalls
	snapshot(t, "call_controls", 420, 120, func(c *ui.Context) {
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() { CallControls(c, CallState{Muted: true, Sharing: true}, k.opts()) })
	})
}
