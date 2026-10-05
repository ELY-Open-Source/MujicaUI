package mujica

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

type popconfirmApp struct {
	open                 bool
	confirmed, cancelled int
}

func (a *popconfirmApp) view(c *ui.Context) {
	Use(c, Settings{})
	ui.Column(c).Padding(30).AlignItems(ui.Start).Children(func() {
		b := Button(c, "Archive", ButtonOptions{Variant: Secondary})
		if b.Clicked() {
			a.open = true
		}
		r := Popconfirm(c, b, &a.open, PopconfirmOptions{Title: "Archive this letter?", Description: "It leaves the inbox."})
		if r.Confirmed {
			a.confirmed++
		}
		if r.Cancelled {
			a.cancelled++
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
	if a.confirmed != 1 || a.cancelled != 0 || a.open {
		t.Fatalf("confirmed %d cancelled %d open %v", a.confirmed, a.cancelled, a.open)
	}
	if !tt.Focused("Archive") {
		t.Fatal("focus not restored to the anchor")
	}
}

func TestPopconfirmDangerFocusesCancel(t *testing.T) {
	open, confirmed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		b := Button(c, "Dismiss", ButtonOptions{})
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
	if a.cancelled != 3 || a.confirmed != 0 || a.open {
		t.Fatalf("cancelled %d confirmed %d open %v", a.cancelled, a.confirmed, a.open)
	}
}

func TestPopconfirmKeepsActionsInWindow(t *testing.T) {
	long := strings.Repeat("The east gate stays unwatched tonight. ", 60)
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(30).Children(func() {
			b := Button(c, "Archive", ButtonOptions{})
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
			b := Button(c, "Dismiss the guard", ButtonOptions{Variant: Secondary})
			Popconfirm(c, b, &open, PopconfirmOptions{Title: "Dismiss the guard?", Description: "The east gate stays unwatched tonight.", Danger: true})
		})
	})
}
