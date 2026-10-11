package overlay

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type popconfirmApp struct {
	open                bool
	confirmed, canceled int
}

func (a *popconfirmApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Padding(30).AlignItems(ui.Start).Children(func() {
		b := input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary})
		if b.Clicked() {
			a.open = true
		}
		r := Popconfirm(c, b, &a.open, PopconfirmOptions{Title: "Archive this letter?", Description: "It leaves the inbox."})
		if r.Confirmed {
			a.confirmed++
		}
		if r.Canceled {
			a.canceled++
		}
	})
}

func TestPopconfirmConfirmByKeyboard(t *testing.T) {
	a := &popconfirmApp{}
	tt := ui.NewTester(a.view, 400, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !a.open || !tt.Focused("Confirm") {
		t.Fatalf("open %v, Confirm focused %v", a.open, tt.Focused("Confirm"))
	}
	tt.Key(0, ui.KeyEnter)
	if a.confirmed != 1 || a.canceled != 0 || a.open {
		t.Fatalf("confirmed %d canceled %d open %v", a.confirmed, a.canceled, a.open)
	}
	if !tt.Focused("Archive") {
		t.Fatal("focus not restored to the anchor")
	}
}

func TestPopconfirmDangerFocusesCancel(t *testing.T) {
	open, confirmed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		b := input.Button(c, "Dismiss", input.ButtonOptions{})
		if b.Clicked() {
			open = true
		}
		if Popconfirm(c, b, &open, PopconfirmOptions{Title: "Dismiss the guard?", Danger: true}).Confirmed {
			confirmed++
		}
	}, 400, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !open || !tt.Focused("Cancel") {
		t.Fatalf("open %v, Cancel focused %v", open, tt.Focused("Cancel"))
	}
	tt.Key(0, ui.KeyEnter)
	if confirmed != 0 || open {
		t.Fatalf("second Enter: confirmed %d open %v", confirmed, open)
	}
}

func TestPopconfirmCancelPaths(t *testing.T) {
	a := &popconfirmApp{}
	tt := ui.NewTester(a.view, 400, 300)
	_ = tt.Click("Archive")
	_ = tt.Click("Cancel")
	_ = tt.Click("Archive")
	tt.Key(0, ui.KeyEscape)
	_ = tt.Click("Archive")
	tt.ClickAt(390, 290)
	if a.canceled != 3 || a.confirmed != 0 || a.open {
		t.Fatalf("canceled %d confirmed %d open %v", a.canceled, a.confirmed, a.open)
	}
}

func TestPopconfirmKeepsActionsInWindow(t *testing.T) {
	long := strings.Repeat("The east gate stays unwatched tonight. ", 60)
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(30).Children(func() {
			b := input.Button(c, "Archive", input.ButtonOptions{})
			Popconfirm(c, b, &open, PopconfirmOptions{Title: "Archive?", Description: long})
		})
	}, 400, 300)
	tt.Frame()
	for _, label := range []string{"Confirm", "Cancel"} {
		r, ok := tt.Find(label)
		if !ok || r.Y < 0 || r.Y+r.H > 300 || r.X+r.W > 400 {
			t.Fatalf("%s at %+v outside the 400×300 window", label, r)
		}
	}
}

func TestPopconfirmSnapshot(t *testing.T) {
	snapshot(t, "popconfirm", 420, 220, func(c *ui.Context) {
		open := true
		ui.Column(c).Padding(24).AlignItems(ui.Start).Children(func() {
			b := input.Button(c, "Dismiss the guard", input.ButtonOptions{Variant: input.Secondary})
			Popconfirm(c, b, &open, PopconfirmOptions{Title: "Dismiss the guard?", Description: "The east gate stays unwatched tonight.", Danger: true})
		})
	})
}
