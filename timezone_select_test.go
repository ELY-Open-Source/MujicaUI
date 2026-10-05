package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var tzNow = time.Date(2026, time.July, 1, 12, 0, 0, 0, time.UTC)

func timezoneTester(zone *string, opts *TimezoneSelectOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		o := *opts
		o.Now = tzNow
		ui.Column(c).Padding(16).Children(func() {
			if TimezoneSelect(c, zone, o).Changed() {
				*n++
			}
		})
	}, 520, 640)
}

func TestTimezoneSelectSearchByNameAndOffset(t *testing.T) {
	zone, n := "", 0
	tt := timezoneTester(&zone, &TimezoneSelectOptions{}, &n)
	_ = tt.Click("Time zone")
	if !tt.HasText("Common") || !tt.HasText("UTC+09:00 Asia/Tokyo") || !tt.HasText("UTC−04:00 America/New York") {
		t.Fatalf("summer offsets and sections: %v", tt.Texts()[:12])
	}
	tt.Type("kolk")
	if !tt.HasText("UTC+05:30 Asia/Kolkata") || tt.HasText("UTC+09:00 Asia/Tokyo") {
		t.Fatalf("name search: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if zone != "Asia/Kolkata" || n != 1 || !tt.HasText("UTC+05:30 · Asia/Kolkata") {
		t.Fatalf("Enter picked %q", zone)
	}
	_ = tt.Click("Time zone")
	tt.Type("+12:45")
	tt.Key(0, ui.KeyEnter)
	if zone != "Pacific/Chatham" {
		t.Fatalf("offset search picked %q", zone)
	}
}

func TestTimezoneSelectKeysCommonClickEmpty(t *testing.T) {
	zone, n := "UTC", 0
	tt := timezoneTester(&zone, &TimezoneSelectOptions{Common: []string{"Europe/Paris"}, Zones: []string{"Asia/Tokyo", "Europe/Paris"}}, &n)
	_ = tt.Click("Time zone")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if zone != "Asia/Tokyo" || n != 1 {
		t.Fatalf("Down then Enter picked %q", zone)
	}
	_ = tt.Click("Time zone")
	_ = tt.Click("UTC+02:00 Europe/Paris")
	if zone != "Europe/Paris" {
		t.Fatalf("click picked %q", zone)
	}
	_ = tt.Click("Time zone")
	tt.Type("atlantis")
	if !tt.HasText("No matching time zone.") {
		t.Fatal("no empty state")
	}
}

func TestTimezoneSelectDisabledAndInvalid(t *testing.T) {
	zone, n := "UTC", 0
	tt := timezoneTester(&zone, &TimezoneSelectOptions{Disabled: true}, &n)
	_ = tt.Click("Time zone")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if tt.HasText("Common") {
		t.Fatal("disabled select opened")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for an unknown zone")
		}
	}()
	bad := "Mars/Olympus"
	timezoneTester(&bad, &TimezoneSelectOptions{}, &n)
}

func TestTimezoneOffsetHalfHours(t *testing.T) {
	if got := TimezoneOffset(timezoneLoad("America/St_Johns"), tzNow); got != "UTC−02:30" {
		t.Fatalf("offset %q", got)
	}
}

func TestTimezoneSelectSnapshot(t *testing.T) {
	zone := "Asia/Shanghai"
	snapshot(t, "timezone_select", 420, 80, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { TimezoneSelect(c, &zone, TimezoneSelectOptions{Now: tzNow}) })
	})
}
