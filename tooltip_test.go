package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

const tooltipTestDelay = 20 * time.Millisecond

func tooltipTester(x, y float32) (*ui.Tester, *int) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(y, 0, 0, x).Children(func() {
			b := Button(c, "Share", ButtonOptions{})
			if b.Clicked() {
				clicks++
			}
			Tooltip(c, b, "Share with the court", TooltipOptions{Delay: tooltipTestDelay})
		})
	}, 400, 300)
	return tt, &clicks
}

func waitTooltip(tt *ui.Tester) {
	time.Sleep(tooltipTestDelay + 10*time.Millisecond)
	tt.Frame()
}

func TestTooltipHoverDelay(t *testing.T) {
	tt, clicks := tooltipTester(40, 40)
	r, _ := tt.Find("Share")
	tt.Move(r.X+4, r.Y+4)
	if tt.HasText("Share with the court") {
		t.Fatal("tip showed before the delay")
	}
	waitTooltip(tt)
	tip, ok := tt.Find("Share with the court")
	if !ok || tip.Y < r.Y+r.H {
		t.Fatalf("tip %+v ok %v below %+v", tip, ok, r)
	}
	// The tip passes the pointer through: clicking it clicks what is below.
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
	waitTooltip(tt)
	if !tt.HasText("Share with the court") {
		t.Fatal("no tip on keyboard focus")
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Share with the court") {
		t.Fatal("Escape did not hide the tip")
	}
	waitTooltip(tt)
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

func TestTooltipSnapshot(t *testing.T) {
	snapshot(t, "tooltip", 360, 140, func(c *ui.Context) {
		ui.Column(c).Padding(24).AlignItems(ui.Start).Children(func() {
			Button(c, "Seal", ButtonOptions{Variant: Secondary})
			tooltipShow(c, ui.Rect{X: 24, Y: 24, W: 62, H: 36}, "Seal the letter with wax")
		})
	})
}
