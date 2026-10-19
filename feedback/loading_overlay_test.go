package feedback

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

type overlayProbe struct {
	loading        bool
	saves, outside int
	cancels        int
	cancel, label  string
}

func overlayTester(p *overlayProbe) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Gap(12).Children(func() {
			r := LoadingOverlay(c, p.loading, LoadingOverlayOptions{Cancel: p.cancel, Label: p.label}, func() {
				ui.Column(c).Padding(24).Children(func() {
					if input.Button(c, "Save", input.ButtonOptions{}).Clicked() {
						p.saves++
					}
				})
			})
			if r.Canceled {
				p.cancels++
			}
			if input.Button(c, "Outside", input.ButtonOptions{}).Clicked() {
				p.outside++
			}
		})
	}, 500, 400)
}

func TestLoadingOverlayBlocksOnlyItsRegion(t *testing.T) {
	p := &overlayProbe{}
	tt := overlayTester(p)
	_ = tt.Click("Save")
	p.loading = true
	tt.Frame()
	r, _ := tt.Find("Save")
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.Key(0, ui.KeyEnter)
	if err := tt.Click("Outside"); err != nil {
		t.Fatal(err)
	}
	if p.saves != 1 || p.outside != 1 {
		t.Fatalf("saves %d outside %d", p.saves, p.outside)
	}
	if !tt.HasText("Loading…") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if a := tt.Announcements(); !slices.Contains(a, "Loading…") {
		t.Fatalf("announcements %q", a)
	}
	p.loading = false
	tt.Frame()
	_ = tt.Click("Save")
	if p.saves != 2 || tt.HasText("Loading…") {
		t.Fatalf("after loading saves %d", p.saves)
	}
}

func TestLoadingOverlayCancelAndFocus(t *testing.T) {
	p := &overlayProbe{cancel: "Stop", label: "Importing"}
	tt := overlayTester(p)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Save") {
		t.Fatal("Save not focused")
	}
	p.loading = true
	tt.Frame()
	if !tt.Focused("Stop") {
		t.Fatal("focus left the region instead of moving to cancel")
	}
	tt.Key(0, ui.KeyEscape)
	if p.cancels != 1 {
		t.Fatalf("Escape cancels %d", p.cancels)
	}
	if err := tt.Click("Stop"); err != nil {
		t.Fatal(err)
	}
	if p.cancels != 2 || p.saves != 0 {
		t.Fatalf("cancels %d saves %d", p.cancels, p.saves)
	}
}

func TestLoadingOverlayFocusWithoutCancel(t *testing.T) {
	p := &overlayProbe{}
	tt := overlayTester(p)
	tt.Key(0, ui.KeyTab)
	p.loading = true
	tt.Frame()
	tt.Frame()
	if !tt.Focused("Loading…") {
		t.Fatal("focus not kept in the region")
	}
	tt.Key(0, ui.KeyEscape)
	if p.cancels != 0 {
		t.Fatal("Escape canceled without a cancel action")
	}
}

func TestLoadingOverlaySnapshot(t *testing.T) {
	snapshot(t, "loading_overlay", 420, 220, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			LoadingOverlay(c, true, LoadingOverlayOptions{Cancel: "Cancel"}, func() {
				ui.Column(c).Padding(16).Gap(8).AlignItems(ui.Start).Background(core.Tokens(c).Surface).Radius(4).Children(func() {
					ui.Text(c, "Quarterly report")
					ui.Text(c, "Revenue, churn and pipeline for Q3.")
					input.Button(c, "Export", input.ButtonOptions{})
				})
			})
		})
	})
}
