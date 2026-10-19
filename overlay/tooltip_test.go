package overlay

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

const tooltipTestDelay = 150 * time.Millisecond

func tooltipTester(x, y float32) (*ui.Tester, *int) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(y, 0, 0, x).Children(func() {
			b := input.Button(c, "Share", input.ButtonOptions{})
			if b.Clicked() {
				clicks++
			}
			Tooltip(c, b, "Share with the court", TooltipOptions{Delay: tooltipTestDelay})
		})
	}, 400, 300)
	return tt, &clicks
}

// waitTooltip waits until the tip shows, however slow the machine.
func waitTooltip(tt *ui.Tester) bool {
	return overlayUntil(tt, func() bool { return tt.HasText("Share with the court") })
}

func TestTooltipHoverDelay(t *testing.T) {
	tt, clicks := tooltipTester(40, 40)
	r, _ := tt.Find("Share")
	start := time.Now()
	tt.Move(r.X+4, r.Y+4)
	if time.Since(start) < tooltipTestDelay && tt.HasText("Share with the court") {
		t.Fatal("tip showed before the delay")
	}
	waitTooltip(tt)
	tip, ok := tt.Find("Share with the court")
	if !ok || tip.Y < r.Y+r.H {
		t.Fatalf("tip %+v ok %v below %+v", tip, ok, r)
	}
	tt.Move(400, 300)
	tt.Frame()
	if tt.HasText("Share with the court") {
		t.Fatal("tip stayed after leave")
	}
	_ = tt.Click("Share")
	if *clicks != 1 {
		t.Fatalf("clicks %d", *clicks)
	}
}

func TestTooltipFocusAndEscape(t *testing.T) {
	tt, _ := tooltipTester(40, 40)
	tt.Key(0, ui.KeyTab)
	if !waitTooltip(tt) {
		t.Fatal("no tip on keyboard focus")
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Share with the court") {
		t.Fatal("Escape did not hide the tip")
	}
	tt.Advance(tooltipTestDelay + 30*time.Millisecond)
	if tt.HasText("Share with the court") {
		t.Fatal("tip came back while still focused")
	}
}

func TestTooltipEdgeFlips(t *testing.T) {
	tt, _ := tooltipTester(330, 262)
	r, _ := tt.Find("Share")
	tt.Move(r.X+4, r.Y+4)
	waitTooltip(tt)
	tip, ok := tt.Find("Share with the court")
	if !ok || tip.Y > r.Y || tip.X+tip.W > 400 {
		t.Fatalf("tip %+v not above and inside for anchor %+v", tip, r)
	}
}

func TestTooltipScalesWithText(t *testing.T) {
	widths := map[float32]float32{}
	for _, scale := range []float32{1, 1.5} {
		tt, _ := tooltipTester(20, 20)
		tt.SetPreferences(ui.Preferences{TextScale: scale})
		tt.Key(0, ui.KeyTab)
		if !overlayUntil(tt, func() bool { return tt.HasText("Share with the court") }) {
			t.Fatalf("scale %v: no tip", scale)
		}
		r, _ := tt.Find("Share with the court")
		widths[scale] = r.W
	}
	if widths[1.5] < widths[1]*1.3 {
		t.Fatalf("tip text width %v at 1.5x, %v at 1x: not scaled", widths[1.5], widths[1])
	}
}

func TestTooltipSnapshot(t *testing.T) {
	snapshotAfter(t, "tooltip", 360, 140, func(c *ui.Context) {
		ui.Column(c).Padding(24).AlignItems(ui.Start).Children(func() {
			b := input.Button(c, "Seal", input.ButtonOptions{Variant: input.Secondary})
			Tooltip(c, b, "Seal the letter with wax", TooltipOptions{Delay: time.Millisecond})
		})
	}, func(tt *ui.Tester) {
		tt.Key(0, ui.KeyTab)
		if !overlayUntil(tt, func() bool { return tt.HasText("Seal the letter with wax") }) {
			t.Fatal("tip did not show")
		}
	})
}
