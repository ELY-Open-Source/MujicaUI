package display

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func headingHeights(t *testing.T, p ui.Preferences) (h1, h2, h3 float32) {
	t.Helper()
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			h1 = Heading(c, "Court", HeadingOptions{}).Bounds().H
			h2 = Heading(c, "Court", HeadingOptions{Level: 2}).Bounds().H
			h3 = Heading(c, "Court", HeadingOptions{Level: 3}).Bounds().H
		})
	}, 400, 400)
	tt.SetPreferences(p)
	tt.Frame()
	return
}

func TestHeadingLevels(t *testing.T) {
	h1, h2, h3 := headingHeights(t, ui.Preferences{TextScale: 1})
	if !(h1 > h2 && h2 > h3 && h3 > 0) {
		t.Fatalf("heights %v %v %v not decreasing", h1, h2, h3)
	}
}

func TestHeadingTextScale(t *testing.T) {
	h1, _, _ := headingHeights(t, ui.Preferences{TextScale: 1})
	big, _, _ := headingHeights(t, ui.Preferences{TextScale: 2})
	if big < h1*1.8 {
		t.Fatalf("scaled height %v, base %v", big, h1)
	}
}

func TestHeadingSubtitleAndDivider(t *testing.T) {
	var plain, full float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			plain = Heading(c, "Archive", HeadingOptions{Level: 2}).Bounds().H
			full = Heading(c, "Ledger", HeadingOptions{Level: 2, Subtitle: "Entries since 1802", Divider: true}).Bounds().H
		})
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("Entries since 1802") {
		t.Fatal("subtitle missing")
	}
	if full < plain+20 {
		t.Fatalf("subtitle and divider add %v DIPs", full-plain)
	}
}

func TestHeadingBadLevelPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Heading(c, "x", HeadingOptions{Level: 4})
	}, 100, 100)
	tt.Frame()
}

func TestHeadingSnapshot(t *testing.T) {
	snapshot(t, "heading", 420, 260, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(16).Width(380).Children(func() {
			Heading(c, "The Grand Hall", HeadingOptions{Subtitle: "Ceremonies of the winter season", Divider: true})
			Heading(c, "Guest Ledger", HeadingOptions{Level: 2})
			Heading(c, "宫廷礼仪", HeadingOptions{Level: 3, Subtitle: "第三章"})
		})
	})
}
