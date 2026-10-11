package media

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func scrubberTester(pos *time.Duration, opts VideoScrubberOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(40).Children(func() {
			if VideoScrubber(c, pos, opts).Changed() {
				*changes++
			}
		})
	}, 500, 300)
}

func TestVideoScrubberKeys(t *testing.T) {
	pos, n := 10*time.Second, 0
	tt := scrubberTester(&pos, VideoScrubberOptions{Duration: time.Minute}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if pos != 15*time.Second {
		t.Fatalf("Right: %v", pos)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	if pos != 5*time.Second {
		t.Fatalf("Left twice: %v", pos)
	}
	tt.Key(0, ui.KeyEnd)
	if pos != time.Minute {
		t.Fatalf("End: %v", pos)
	}
	tt.Key(0, ui.KeyHome)
	if pos != 0 || n != 5 {
		t.Fatalf("Home: %v, changes %d", pos, n)
	}
}

func TestVideoScrubberClickAndHover(t *testing.T) {
	pos, n := time.Duration(0), 0
	var asked time.Duration
	opts := VideoScrubberOptions{Duration: 100 * time.Second, Preview: func(at time.Duration) { asked = at }}
	tt := scrubberTester(&pos, opts, &n)
	r, ok := tt.Find("Timeline")
	if !ok {
		t.Fatal("no timeline")
	}
	mid := r.X + r.W/2
	tt.Move(mid, r.Y+r.H/2)
	if asked < 45*time.Second || asked > 55*time.Second || !tt.HasText(core.MediaClock(asked)) {
		t.Fatalf("hover preview at %v, texts %v", asked, tt.Texts())
	}
	tt.ClickAt(mid, r.Y+r.H/2)
	if pos < 45*time.Second || pos > 55*time.Second || n != 1 {
		t.Fatalf("click seek %v, changes %d", pos, n)
	}
}

func TestVideoScrubberDisabled(t *testing.T) {
	pos, n := 10*time.Second, 0
	tt := scrubberTester(&pos, VideoScrubberOptions{Duration: time.Minute, Disabled: true}, &n)
	r, _ := tt.Find("Timeline")
	tt.ClickAt(r.X+r.W-10, r.Y+r.H/2)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if pos != 10*time.Second || n != 0 {
		t.Fatalf("disabled moved to %v", pos)
	}
}

func TestVideoScrubberBadOptions(t *testing.T) {
	pos := time.Duration(0)
	mediaPanics(t, "no duration", func(c *ui.Context) { VideoScrubber(c, &pos, VideoScrubberOptions{}) })
	far := 2 * time.Minute
	mediaPanics(t, "past end", func(c *ui.Context) { VideoScrubber(c, &far, VideoScrubberOptions{Duration: time.Minute}) })
}

func TestVideoScrubberSnapshot(t *testing.T) {
	pos := 40 * time.Second
	snapshot(t, "video_scrubber", 420, 80, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			VideoScrubber(c, &pos, VideoScrubberOptions{Duration: 2 * time.Minute,
				Buffered: []MediaRange{{0, 70 * time.Second}, {90 * time.Second, 100 * time.Second}}})
		})
	})
}
