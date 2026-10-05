package data

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

type carouselHarness struct {
	index   int
	opts    CarouselOptions
	changes int
}

func (h *carouselHarness) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Padding(10).Size(400, 260).Children(func() {
		v := Carousel(c, &h.index, 4, h.opts, func(i int) {
			ui.Text(c, fmt.Sprintf("Window %d", i+1)).Padding(20)
		})
		v.Element.Grow(1)
		if v.Changed() {
			h.changes++
		}
	})
}

func TestCarouselButtonsAndDots(t *testing.T) {
	h := &carouselHarness{}
	tt := ui.NewTester(h.view, 420, 280)
	_ = tt.Click("Previous slide")
	if h.index != 0 || h.changes != 0 {
		t.Fatal("previous moved past the first slide")
	}
	_ = tt.Click("Next slide")
	if h.index != 1 || !tt.HasText("Window 2") || h.changes != 1 {
		t.Fatalf("next: %d changes %d %q", h.index, h.changes, tt.Texts())
	}
	_ = tt.Click("Slide 4 of 4")
	if h.index != 3 {
		t.Fatalf("dot: %d", h.index)
	}
	_ = tt.Click("Next slide")
	if h.index != 3 {
		t.Fatal("next moved past the last slide without Loop")
	}
	h.opts.Loop = true
	tt.Frame()
	_ = tt.Click("Next slide")
	if h.index != 0 {
		t.Fatalf("loop: %d", h.index)
	}
}

func TestCarouselKeyboard(t *testing.T) {
	h := &carouselHarness{}
	tt := ui.NewTester(h.view, 420, 280)
	tt.Key(0, ui.KeyTab) // previous is disabled: focus lands on next
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if h.index != 2 {
		t.Fatalf("Right×2: %d", h.index)
	}
	tt.Key(0, ui.KeyEnd)
	if h.index != 3 {
		t.Fatalf("End: %d", h.index)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyLeft)
	if h.index != 0 {
		t.Fatalf("Home, Left: %d", h.index)
	}
	if got := tt.Announcements(); len(got) == 0 || got[len(got)-1] != "Slide 1 of 4" {
		t.Fatalf("announcements %q", got)
	}
}

func TestCarouselNoAutoplayByDefault(t *testing.T) {
	h := &carouselHarness{opts: CarouselOptions{Interval: 10 * time.Millisecond}}
	tt := ui.NewTester(h.view, 420, 280)
	time.Sleep(30 * time.Millisecond)
	tt.Frame()
	if h.index != 0 || tt.HasText("Pause autoplay") {
		t.Fatalf("advanced without Autoplay: %d", h.index)
	}
}

func TestCarouselAutoplayPause(t *testing.T) {
	h := &carouselHarness{opts: CarouselOptions{Autoplay: true, Interval: 20 * time.Millisecond}}
	tt := ui.NewTester(h.view, 420, 280)
	time.Sleep(30 * time.Millisecond)
	tt.Frame()
	if h.index != 1 {
		t.Fatalf("autoplay did not advance: %d", h.index)
	}
	_ = tt.Click("Pause autoplay")
	tt.Move(1, 1) // away from the carousel
	tt.SetFocused(false)
	tt.SetFocused(true)
	at := h.index
	time.Sleep(50 * time.Millisecond)
	tt.Frame()
	if h.index != at || !tt.HasText("") {
		t.Fatalf("paused carousel advanced: %d -> %d", at, h.index)
	}
	if _, ok := tt.Find("Resume autoplay"); !ok {
		t.Fatal("no resume button")
	}
	r, _ := tt.Find("Window " + fmt.Sprint(at+1))
	_ = tt.Click("Resume autoplay")
	tt.Move(r.X+2, r.Y+2) // hover pauses too
	time.Sleep(50 * time.Millisecond)
	tt.Frame()
	if h.index != at {
		t.Fatalf("hovered carousel advanced: %d -> %d", at, h.index)
	}
}

func TestCarouselSnapshot(t *testing.T) {
	snapshot(t, "carousel", 440, 240, func(c *ui.Context) {
		i := ui.Local(c.Root(), "i", func() *int { v := 1; return &v })
		ui.Column(c).Padding(16).Size(440, 240).Children(func() {
			Carousel(c, *i, 4, CarouselOptions{Autoplay: true, Label: "Windows"}, func(n int) {
				ui.Column(c).Grow(1).Center().Children(func() {
					ui.Text(c, fmt.Sprintf("Rose window %d", n+1)).FontSize(core.FontSize(c, 18))
				})
			}).Element.Grow(1)
		})
	})
}

func TestCarouselEmptyAutoplayKeys(t *testing.T) {
	index := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Carousel(c, &index, 0, CarouselOptions{Autoplay: true}, func(int) {}).Element.Size(300, 200)
	}, 320, 240)
	for range 4 {
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyRight)
		tt.Key(0, ui.KeyLeft)
	}
	if index != 0 {
		t.Fatalf("index %d", index)
	}
}

func TestCarouselDotsFollowArrows(t *testing.T) {
	h := &carouselHarness{}
	tt := ui.NewTester(h.view, 420, 280)
	for range 10 {
		if tt.Focused("Slide 1 of 4") {
			break
		}
		tt.Key(0, ui.KeyTab)
	}
	if !tt.Focused("Slide 1 of 4") {
		t.Fatal("Tab never reached the first dot")
	}
	tt.Key(0, ui.KeyRight)
	if h.index != 1 || !tt.Focused("Slide 2 of 4") {
		t.Fatalf("Right: index %d, focus on dot 2 %v", h.index, tt.Focused("Slide 2 of 4"))
	}
	tt.Key(0, ui.KeyEnter)
	if h.index != 1 {
		t.Fatalf("Enter on the focused dot moved to %d", h.index)
	}
	tt.Key(0, ui.KeyEnd)
	if h.index != 3 || !tt.Focused("Slide 4 of 4") {
		t.Fatalf("End: index %d", h.index)
	}
}
