package mujica

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestBannerActionCloseAndAnnounce(t *testing.T) {
	open := true
	var last AlertResult
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if r := Banner(c, &open, "Maintenance tonight at 22:00", BannerOptions{Severity: SeverityWarning, Action: "Details"}); r.Action || r.Closed {
			last = r
		}
	}, 700, 200)
	if a := tt.Announcements(); !slices.Equal(a, []string{"Warning: Maintenance tonight at 22:00"}) {
		t.Fatalf("announcements %q", a)
	}
	if err := tt.Click("Details"); err != nil {
		t.Fatal(err)
	}
	if !last.Action || !open {
		t.Fatalf("action: %+v open %v", last, open)
	}
	tt.Key(0, ui.KeyTab) // from Details to Dismiss
	tt.Key(0, ui.KeyEnter)
	if open || !last.Closed || tt.HasText("Maintenance") {
		t.Fatalf("keyboard close: open %v %+v", open, last)
	}
}

func TestBannerDisabledAndFixed(t *testing.T) {
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(true).Children(func() {
			Banner(c, &open, "Read only", BannerOptions{Severity: SeverityError})
		})
		Banner(c, nil, "Always here", BannerOptions{})
	}, 700, 200)
	_ = tt.Click("Dismiss")
	if !open {
		t.Fatal("disabled banner closed")
	}
	if !tt.HasText("Always here") {
		t.Fatal("banner without open missing")
	}
	for _, f := range []func(c *ui.Context){
		func(c *ui.Context) { Banner(c, nil, "", BannerOptions{}) },
		func(c *ui.Context) { Banner(c, nil, "x", BannerOptions{Severity: 7}) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("no panic")
				}
			}()
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); f(c) }, 300, 100)
		}()
	}
}

func TestBannerSnapshot(t *testing.T) {
	open := true
	snapshot(t, "banner", 640, 200, func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			Banner(c, &open, "A new version is available.", BannerOptions{Action: "Update"})
			Banner(c, &open, "Your trial ends in 3 days.", BannerOptions{Severity: SeverityWarning})
			Banner(c, nil, "Sync is failing; changes stay on this device.", BannerOptions{Severity: SeverityError, Action: "Retry"})
		})
	})
}
