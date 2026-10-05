package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

type pinInputProbe struct {
	value              string
	changes, completes int
}

func (p *pinInputProbe) tester(opts PinInputOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := PinInput(c, &p.value, opts)
		if r.Changed() {
			p.changes++
		}
		if r.Completed() {
			p.completes++
		}
	}, 480, 120)
}

func TestPinInputTypingAndBackspace(t *testing.T) {
	p := &pinInputProbe{}
	tt := p.tester(PinInputOptions{Length: 4, Label: "Code"})
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Code") {
		t.Fatal("not focusable")
	}
	tt.TypeKey(0, ui.Key0, "0")
	tt.Type("a")
	tt.Type("07")
	if p.value != "007" || p.completes != 0 {
		t.Fatalf("value %q completes %d", p.value, p.completes)
	}
	tt.Type("9")
	tt.Type("5")
	if p.value != "0079" || p.completes != 1 {
		t.Fatalf("full: value %q completes %d", p.value, p.completes)
	}
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyBackspace)
	if p.value != "00" {
		t.Fatalf("backspace: %q", p.value)
	}
	tt.Type("12")
	if p.value != "0012" || p.completes != 2 {
		t.Fatalf("refill: %q completes %d", p.value, p.completes)
	}
}

func TestPinInputPasteAndCompose(t *testing.T) {
	p := &pinInputProbe{value: "1"}
	tt := p.tester(PinInputOptions{Label: "Code"})
	tt.Key(0, ui.KeyTab)
	tt.SetClipboard(" 012-345 9")
	tt.Command("paste")
	if p.value != "012345" || p.completes != 1 || p.changes != 1 {
		t.Fatalf("paste: %q completes %d changes %d", p.value, p.completes, p.changes)
	}
	tt.Key(0, ui.KeyBackspace)
	tt.Compose("4", 1)
	if p.value != "01234" || !tt.HasText("4") {
		t.Fatalf("compose committed early: %q", p.value)
	}
	tt.Type("8")
	if p.value != "012348" {
		t.Fatalf("commit: %q", p.value)
	}
}

func TestPinInputDisabled(t *testing.T) {
	p := &pinInputProbe{value: "12"}
	tt := p.tester(PinInputOptions{Label: "Code", Disabled: true})
	tt.Key(0, ui.KeyTab)
	tt.Type("3")
	if tt.Focused("Code") || p.value != "12" {
		t.Fatalf("disabled took input: %q", p.value)
	}
}

func TestPinInputRejectsBadValue(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("letters accepted")
		}
	}()
	p := &pinInputProbe{value: "12a"}
	p.tester(PinInputOptions{})
}

func TestPinInputSnapshot(t *testing.T) {
	a, b := "0420", "13"
	snapshot(t, "pin_input", 420, 150, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(16).Children(func() {
			PinInput(c, &a, PinInputOptions{}).Element.Focus()
			PinInput(c, &b, PinInputOptions{Length: 4, Invalid: true})
		})
	})
}
