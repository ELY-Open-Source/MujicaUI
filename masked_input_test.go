package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestMaskFormat(t *testing.T) {
	for _, tc := range [][3]string{
		{"(999) 999-9999", "555123", "(555) 123-____"},
		{"9999 9999 9999 9999", "", "____ ____ ____ ____"},
		{"aa-99", "AB12", "AB-12"},
	} {
		if got := MaskFormat(tc[0], tc[1]); got != tc[2] {
			t.Errorf("MaskFormat(%q, %q) = %q, want %q", tc[0], tc[1], got, tc[2])
		}
	}
}

func maskedTester(v *string, dis *bool, n, submits *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := MaskedInput(c, v, MaskedInputOptions{Mask: "(999) 999-9999", Label: "Phone", Disabled: *dis})
		if r.Changed() {
			*n++
		}
		if r.Submitted() {
			*submits++
		}
	}, 400, 120)
}

func TestMaskedInputTypingEditingAndPaste(t *testing.T) {
	v := ""
	dis := false
	n, submits := 0, 0
	tt := maskedTester(&v, &dis, &n, &submits)
	tt.Key(0, ui.KeyTab)
	tt.Type("55x5-12")
	if v != "55512" {
		t.Fatalf("typed: %q", v)
	}
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	tt.Type("9")
	if v != "55951" {
		t.Fatalf("insert at caret: %q", v)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyDelete)
	if v != "5951" {
		t.Fatalf("Delete at start: %q", v)
	}
	tt.Key(0, ui.KeyEnd)
	tt.SetClipboard("123-456-789")
	tt.Command("paste")
	if v != "5951123456" {
		t.Fatalf("paste fills to the end: %q", v)
	}
	tt.Type("7")
	if v != "5951123456" {
		t.Fatalf("typed past the mask: %q", v)
	}
	tt.Key(0, ui.KeyEnter)
	if submits != 1 || n == 0 {
		t.Fatalf("submits %d changes %d", submits, n)
	}
	if _, ok := tt.Find("Phone"); !ok {
		t.Fatal("no field name")
	}
}

func TestMaskedInputCompositionAndDisabled(t *testing.T) {
	v := "12"
	dis := false
	n, submits := 0, 0
	tt := maskedTester(&v, &dis, &n, &submits)
	tt.Key(0, ui.KeyTab)
	tt.Compose("3", 1)
	if v != "12" || !tt.HasText("3") {
		t.Fatalf("composition committed early: %q", v)
	}
	tt.Key(0, ui.KeyBackspace) // the input method's, while composing
	if v != "12" {
		t.Fatalf("Backspace during composition edited: %q", v)
	}
	tt.Type("3")
	if v != "123" {
		t.Fatalf("commit: %q", v)
	}
	dis = true
	tt.Frame()
	tt.Type("4")
	tt.Key(0, ui.KeyBackspace)
	if v != "123" {
		t.Fatalf("disabled edit: %q", v)
	}
}

func TestMaskedInputRejects(t *testing.T) {
	for name, tc := range map[string][2]string{"no slots": {"---", ""}, "too long": {"99", "123"}, "wrong kind": {"99", "1a"}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			v := tc[1]
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); MaskedInput(c, &v, MaskedInputOptions{Mask: tc[0]}) }, 200, 100)
		}()
	}
}

func TestMaskedInputSnapshot(t *testing.T) {
	a, b, d := "5551234", "", "4111111111111111"
	snapshot(t, "masked_input", 520, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(12).Children(func() {
			MaskedInput(c, &a, MaskedInputOptions{Mask: "(999) 999-9999"}).Element.Width(160)
			MaskedInput(c, &b, MaskedInputOptions{Mask: "99/99/9999"}).Element.Width(130)
			MaskedInput(c, &d, MaskedInputOptions{Mask: "9999 9999 9999 9999", Disabled: true}).Element.Width(190)
		})
	})
}
