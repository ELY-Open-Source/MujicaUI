package input

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestTextAreaNewlinesAndSubmit(t *testing.T) {
	s, submits := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if TextArea(c, &s, TextAreaOptions{Label: "Notes"}).Submitted() {
			submits++
		}
	}, 400, 300)
	_ = tt.Click("Notes")
	tt.Type("one")
	tt.Key(0, ui.KeyEnter)
	tt.Type("two")
	if s != "one\ntwo" || submits != 0 {
		t.Fatalf("value %q submits %d", s, submits)
	}
	tt.Key(ui.Cmd, ui.KeyEnter)
	if submits != 1 || s != "one\ntwo" {
		t.Fatalf("Cmd+Enter: value %q submits %d", s, submits)
	}
}

func TestTextAreaCustomSubmitKey(t *testing.T) {
	s, submits := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if TextArea(c, &s, TextAreaOptions{Label: "Chat", SubmitKey: ui.KeyEnter, SubmitMods: ui.Ctrl | ui.Shift}).Submitted() {
			submits++
		}
	}, 400, 300)
	_ = tt.Click("Chat")
	tt.Type("hi")
	tt.Key(ui.Cmd, ui.KeyEnter)
	tt.Key(ui.Ctrl|ui.Shift, ui.KeyEnter)
	if submits != 1 {
		t.Fatalf("submits %d", submits)
	}
}

func TestTextAreaHeightsAndCount(t *testing.T) {
	s := strings.Repeat("line\n", 40) + "end"
	var h float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).AlignItems(ui.Start).Width(300).Children(func() {
			r := TextArea(c, &s, TextAreaOptions{Label: "Log", MinHeight: 60, MaxHeight: 120, MaxLength: 100})
			h = r.Input.Bounds().H
		})
	}, 400, 400)
	tt.Frame()
	if h != 120 {
		t.Fatalf("height %v, want the 120 cap", h)
	}
	if !tt.HasText("203 of 100 characters") || !tt.HasText("103 characters over the limit") {
		t.Fatalf("count missing: %q", tt.Texts())
	}
	s = "ok"
	tt.Frame()
	if h < 60 || h >= 120 || !tt.HasText("2 of 100 characters") {
		t.Fatalf("height %v texts %q", h, tt.Texts())
	}
}

func TestTextAreaBadOptionsPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for MaxHeight below MinHeight")
		}
	}()
	s := ""
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TextArea(c, &s, TextAreaOptions{MinHeight: 100, MaxHeight: 50})
	}, 300, 200)
}

func TestTextAreaSnapshot(t *testing.T) {
	snapshot(t, "text_area", 420, 360, func(c *ui.Context) {
		a, b, e := "", "Les lys d’or sur fond d’azur.\nUne cour, un bal, un serment.", "Trop long"
		ui.Column(c).Padding(24).Gap(12).Width(360).Children(func() {
			TextArea(c, &a, TextAreaOptions{Label: "A", Placeholder: "Write a dedication…"})
			TextArea(c, &b, TextAreaOptions{Label: "B", MaxLength: 120})
			TextArea(c, &e, TextAreaOptions{Label: "E", MinHeight: 40, MaxLength: 5})
		})
	})
}
