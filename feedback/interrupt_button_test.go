package feedback

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestInterruptButtonDirectAndConfirmed(t *testing.T) {
	running, interrupting, confirm := true, false, false
	stops := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Children(func() {
			if InterruptButton(c, running, InterruptButtonOptions{Confirm: confirm, Interrupting: interrupting}).Changed() {
				stops++
			}
		})
	}, 400, 300)
	if err := tt.Click("Stop"); err != nil {
		t.Fatal(err)
	}
	if stops != 1 {
		t.Fatalf("direct stop %d", stops)
	}
	confirm = true
	tt.Frame()
	_ = tt.Click("Stop")
	if stops != 1 || !tt.HasText("Stop the agent?") {
		t.Fatalf("confirm not asked: %q", tt.Texts())
	}
	_ = tt.Click("Keep running")
	if stops != 1 || tt.HasText("Stop the agent?") {
		t.Fatal("cancel interrupted")
	}
	_ = tt.Click("Stop")
	tt.Key(0, ui.KeyTab) // Danger focuses Cancel first; Tab reaches Stop
	tt.Key(0, ui.KeyEnter)
	if stops != 2 {
		t.Fatalf("confirmed stop %d: %q", stops, tt.Texts())
	}
	interrupting = true
	tt.Frame()
	_ = tt.Click("Stopping")
	if stops != 2 || !tt.HasText("Stopping") {
		t.Fatalf("interrupting: stops %d %q", stops, tt.Texts())
	}
	interrupting, running = false, false
	tt.Frame()
	_ = tt.Click("Stop")
	if stops != 2 {
		t.Fatal("idle button interrupted")
	}
}

func TestInterruptButtonInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		InterruptButton(c, false, InterruptButtonOptions{Interrupting: true})
	}, 100, 50, 1)
}

func TestInterruptButtonSnapshot(t *testing.T) {
	snapshot(t, "interrupt_button", 460, 90, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(12).Children(func() {
			InterruptButton(c, true, InterruptButtonOptions{})
			InterruptButton(c, true, InterruptButtonOptions{Interrupting: true})
			InterruptButton(c, false, InterruptButtonOptions{Label: "Cancel run"})
		})
	})
}

func TestInterruptButtonAncestorDisabledClosesConfirmation(t *testing.T) {
	disabled := false
	stops := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			if InterruptButton(c, true, InterruptButtonOptions{Confirm: true}).Changed() {
				stops++
			}
		})
	}, 400, 300)
	_ = tt.Click("Stop")
	disabled = true
	tt.Frame()
	if tt.HasText("Stop the agent?") || stops != 0 {
		t.Fatal("disabled confirmation retained")
	}
	_ = tt.Click("Stop")
	tt.Key(0, ui.KeyEnter)
	if stops != 0 {
		t.Fatal("disabled interrupt delivered")
	}
}
