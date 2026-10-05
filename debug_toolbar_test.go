package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func debugTester(state *DebugState, disabled bool, got *[]DebugAction) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if a := DebugToolbar(c, *state, DebugToolbarOptions{Disabled: disabled}).Action; a != DebugNone {
			*got = append(*got, a)
		}
	}, 480, 120)
}

func TestDebugEnabledRules(t *testing.T) {
	cases := []struct {
		state DebugState
		on    []DebugAction
	}{
		{DebugStopped, []DebugAction{DebugContinue}},
		{DebugRunning, []DebugAction{DebugPause, DebugRestart, DebugStop}},
		{DebugPaused, []DebugAction{DebugContinue, DebugStepOver, DebugStepInto, DebugStepOut, DebugRestart, DebugStop}},
	}
	for _, tc := range cases {
		n := 0
		for a := DebugContinue; a <= DebugStop; a++ {
			if DebugEnabled(tc.state, a) {
				n++
			}
		}
		if n != len(tc.on) {
			t.Fatalf("state %d enables %d actions, want %d", tc.state, n, len(tc.on))
		}
		for _, a := range tc.on {
			if !DebugEnabled(tc.state, a) {
				t.Fatalf("state %d disables %d", tc.state, a)
			}
		}
	}
}

func TestDebugToolbarClicksFollowState(t *testing.T) {
	state := DebugPaused
	var got []DebugAction
	tt := debugTester(&state, false, &got)
	for _, name := range []string{"Step over", "Step into", "Step out", "Continue"} {
		if err := tt.Click(name); err != nil {
			t.Fatal(err)
		}
	}
	want := []DebugAction{DebugStepOver, DebugStepInto, DebugStepOut, DebugContinue}
	if len(got) != 4 || got[0] != want[0] || got[3] != want[3] {
		t.Fatalf("actions %v", got)
	}
	state = DebugRunning
	tt.Frame()
	_ = tt.Click("Step over")
	if err := tt.Click("Pause"); err != nil {
		t.Fatal(err)
	}
	if len(got) != 5 || got[4] != DebugPause {
		t.Fatalf("running: actions %v", got)
	}
	state = DebugStopped
	tt.Frame()
	if _, ok := tt.Find("Start"); !ok {
		t.Fatal("stopped toolbar does not offer Start")
	}
}

func TestDebugToolbarShortcuts(t *testing.T) {
	state := DebugPaused
	var got []DebugAction
	tt := debugTester(&state, false, &got)
	tt.Key(0, ui.KeyF10)
	tt.Key(ui.Shift, ui.KeyF11)
	tt.Key(ui.Shift, ui.KeyF5)
	state = DebugStopped
	tt.Frame()
	tt.Key(0, ui.KeyF10)
	tt.Key(0, ui.KeyF5)
	want := []DebugAction{DebugStepOver, DebugStepOut, DebugStop, DebugContinue}
	if len(got) != len(want) {
		t.Fatalf("actions %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("actions %v, want %v", got, want)
		}
	}
}

func TestDebugToolbarDisabled(t *testing.T) {
	state := DebugPaused
	var got []DebugAction
	tt := debugTester(&state, true, &got)
	_ = tt.Click("Continue")
	tt.Key(0, ui.KeyF10)
	if len(got) != 0 {
		t.Fatalf("disabled toolbar acted: %v", got)
	}
}

func TestDebugToolbarSnapshot(t *testing.T) {
	snapshot(t, "debug_toolbar", 340, 70, func(c *ui.Context) {
		ui.Row(c).Padding(12).Children(func() {
			DebugToolbar(c, DebugPaused, DebugToolbarOptions{})
		})
	})
}
