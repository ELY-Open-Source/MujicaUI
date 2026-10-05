package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestPasswordInputMaskPasteAndReveal(t *testing.T) {
	pw, changes, submits := "", 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := PasswordInput(c, &pw, PasswordInputOptions{Label: "Password"})
		if r.Changed() {
			changes++
		}
		if r.Submitted() {
			submits++
		}
	}, 400, 200)
	_ = tt.Click("Password")
	tt.SetClipboard("s3cr")
	tt.Command("paste")
	tt.Type("et")
	if pw != "s3cret" || changes != 2 {
		t.Fatalf("value matched %v (length %d), changes %d", pw == "s3cret", len(pw), changes)
	}
	tt.SetClipboard("")
	tt.Command("selectAll")
	tt.Command("copy")
	if tt.Clipboard() != "" {
		t.Fatalf("masked field copied %d characters", len(tt.Clipboard()))
	}
	for _, s := range tt.Texts() {
		if s == "s3cret" {
			t.Fatal("masked password shows in the texts")
		}
	}
	if err := tt.Click("Show password"); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.Find("Hide password"); !ok || pw != "s3cret" {
		t.Fatalf("after reveal: value matched %v, hide toggle shown %v", pw == "s3cret", ok)
	}
	_ = tt.Click("Password")
	tt.Command("selectAll")
	tt.Command("copy")
	if tt.Clipboard() != "s3cret" {
		t.Fatalf("revealed copy matched %v (length %d)", tt.Clipboard() == "s3cret", len(tt.Clipboard()))
	}
	tt.Key(0, ui.KeyEnter)
	if submits != 1 {
		t.Fatalf("submits %d", submits)
	}
	_ = tt.Click("Hide password")
	if _, ok := tt.Find("Show password"); !ok {
		t.Fatal("hide did not mask again")
	}
}

func TestPasswordInputToggleByKeyboard(t *testing.T) {
	pw := "x"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		PasswordInput(c, &pw, PasswordInputOptions{Label: "PIN"})
	}, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if _, ok := tt.Find("Hide password"); !ok {
		t.Fatal("Space on the toggle did not reveal")
	}
}

func TestPasswordInputDisabled(t *testing.T) {
	pw := "locked"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		PasswordInput(c, &pw, PasswordInputOptions{Label: "Off", Disabled: true})
	}, 400, 200)
	_ = tt.Click("Show password")
	_ = tt.Click("Off")
	tt.Type("zz")
	if pw != "locked" {
		t.Fatalf("disabled value matched %v (length %d)", pw == "locked", len(pw))
	}
	if _, ok := tt.Find("Show password"); !ok {
		t.Fatal("disabled toggle revealed")
	}
}

func TestPasswordInputSnapshot(t *testing.T) {
	snapshot(t, "password_input", 420, 220, func(c *ui.Context) {
		a, b, e := "", "versailles", "short"
		ui.Column(c).Padding(24).Gap(12).Width(320).Children(func() {
			PasswordInput(c, &a, PasswordInputOptions{Label: "A", Placeholder: "Password"})
			PasswordInput(c, &b, PasswordInputOptions{Label: "B"})
			PasswordInput(c, &e, PasswordInputOptions{Label: "E", Error: "At least 8 characters"})
		})
	})
}
