package data

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestStatisticFormat(t *testing.T) {
	for _, tc := range []struct {
		v    float64
		p    int
		want string
	}{
		{1234567.891, 2, "1,234,567.89"},
		{999, 0, "999"},
		{1000, 0, "1,000"},
		{-42.5, 1, "−42.5"},
		{-0.001, 1, "0.0"},
		{0, 3, "0.000"},
		{math.Inf(1), 0, "+Inf"},
	} {
		if got := StatisticFormat(tc.v, tc.p); got != tc.want {
			t.Errorf("format(%v, %d) = %q, want %q", tc.v, tc.p, got, tc.want)
		}
	}
}

func TestStatisticChange(t *testing.T) {
	for _, tc := range []struct {
		change float64
		text   string
		word   string
	}{{3.26, "+3.3%", "up"}, {-1.5, "−1.5%", "down"}, {0, "0.0%", "unchanged"}} {
		ch := tc.change
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			Statistic(c, 12840, StatisticOptions{Title: "Pilgrims", Prefix: "≈", Unit: "souls", Precision: 1, Change: &ch, ChangeUnit: "%", ChangeNote: "vs last feast"})
		}, 400, 200)
		for _, s := range []string{"Pilgrims", "≈12,840.0", "souls", tc.text, "vs last feast"} {
			if !tt.HasText(s) {
				t.Fatalf("%v: missing %q in %q", tc.change, s, tt.Texts())
			}
		}
		want := "Pilgrims ≈12,840.0 souls, " + tc.word + " " + tc.text
		if _, ok := tt.Find(want); !ok {
			t.Fatalf("spoken label %q missing", want)
		}
	}
}

func TestStatisticCustomFormatAndBadPrecision(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Statistic(c, 0.5, StatisticOptions{Format: func(v float64) string { return "half" }})
	}, 200, 100)
	if !tt.HasText("half") {
		t.Fatalf("custom format: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Statistic(c, 1, StatisticOptions{Precision: -1})
	}, 100, 100)
}

func TestStatisticSnapshot(t *testing.T) {
	up, down := 12.5, -3.0
	snapshot(t, "statistic", 560, 140, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(40).AlignItems(ui.Start).Children(func() {
			Statistic(c, 12840, StatisticOptions{Title: "Pilgrims", Unit: "souls", Change: &up, ChangeUnit: "%", Precision: 1, ChangeNote: "vs last feast"})
			Statistic(c, 3.75, StatisticOptions{Title: "Tithe", Prefix: "£", Precision: 2, Change: &down, ChangeUnit: "%"})
			Statistic(c, 7, StatisticOptions{Title: "Bells"})
		})
	})
}

func TestStatisticTextScale(t *testing.T) {
	var h float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).AlignItems(ui.Start).Children(func() { h = Statistic(c, 7, StatisticOptions{}).Bounds().H })
	}, 300, 300)
	tt.Frame()
	base := h
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	tt.Frame()
	if h < base*1.8 {
		t.Fatalf("height %v at 200%%, %v at 100%%", h, base)
	}
}
