package mujica

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

type popoverApp struct {
	open            bool
	top             float32
	panel, anchor   ui.Rect
	behind, applied int
}

func (a *popoverApp) view(c *ui.Context) {
	Use(c, Settings{})
	ui.Column(c).Padding(a.top, 0, 0, 40).Gap(8).Children(func() {
		b := Button(c, "Filter", ButtonOptions{})
		if b.Clicked() {
			a.open = !a.open
		}
		a.anchor = b.Bounds()
		if Button(c, "Behind", ButtonOptions{}).Clicked() {
			a.behind++
		}
		if p := Popover(c, b, &a.open, PopoverOptions{Title: "Filter rows"}, func() {
			if Button(c, "Apply", ButtonOptions{}).Clicked() {
				a.applied++
				a.open = false
			}
		}); p != nil {
			a.panel = p.Bounds()
		}
	})
}

func TestPopoverOpensBelowAndFlips(t *testing.T) {
	a := &popoverApp{top: 20}
	tt := ui.NewTester(a.view, 400, 300)
	_ = tt.Click("Filter")
	tt.Frame()
	if !a.open || a.panel.Y < a.anchor.Y+a.anchor.H {
		t.Fatalf("below: open %v panel %+v anchor %+v", a.open, a.panel, a.anchor)
	}
	a = &popoverApp{top: 240}
	tt = ui.NewTester(a.view, 400, 300)
	_ = tt.Click("Filter")
	tt.Frame()
	if a.panel.Y+a.panel.H > a.anchor.Y {
		t.Fatalf("no flip: panel %+v anchor %+v", a.panel, a.anchor)
	}
}

func TestPopoverOutsideClickAndEscape(t *testing.T) {
	a := &popoverApp{top: 20}
	tt := ui.NewTester(a.view, 400, 300)
	_ = tt.Click("Filter")
	tt.ClickAt(390, 290)
	if a.open || a.behind != 0 {
		t.Fatalf("outside click: open %v behind %d", a.open, a.behind)
	}
	_ = tt.Click("Filter")
	tt.Key(0, ui.KeyTab) // focus follows the anchor into the panel
	if !tt.Focused("Apply") {
		t.Fatal("Tab did not enter the popover")
	}
	tt.Key(0, ui.KeyEscape)
	if a.open || !tt.Focused("Filter") {
		t.Fatalf("Escape: open %v anchor focused %v", a.open, tt.Focused("Filter"))
	}
	_ = tt.Click("Filter")
	_ = tt.Click("Apply")
	if a.applied != 1 || a.open {
		t.Fatalf("applied %d open %v", a.applied, a.open)
	}
}

func TestPopoverFitsWindow(t *testing.T) {
	open := true
	var panel ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(30).Children(func() {
			b := Button(c, "Filter", ButtonOptions{})
			if p := Popover(c, b, &open, PopoverOptions{Width: 1000}, func() {
				ui.Text(c, strings.Repeat("Long row of guests. ", 200))
				Button(c, "Bottom", ButtonOptions{})
			}); p != nil {
				panel = p.Bounds()
			}
		})
	}, 400, 300)
	tt.Frame()
	if panel.W > 400 || panel.X < 0 || panel.Y+panel.H > 300 {
		t.Fatalf("panel %+v outside the 400×300 window", panel)
	}
}

func TestPopoverNeedsAnchor(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	open := true
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Popover(c, nil, &open, PopoverOptions{}, nil)
	}, 200, 200)
}

func TestPopoverSnapshot(t *testing.T) {
	snapshot(t, "popover", 420, 260, func(c *ui.Context) {
		open := true
		ui.Column(c).Padding(24).AlignItems(ui.Start).Children(func() {
			b := Button(c, "Guests", ButtonOptions{Variant: Secondary})
			Popover(c, b, &open, PopoverOptions{Title: "Guest list", Width: 240}, func() {
				ui.Text(c, "Twelve guests confirmed; three await a reply.").TextColor(Tokens(c).TextMuted)
			})
		})
	})
}
