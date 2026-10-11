package media

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func stripFrames(n int) []StripFrame {
	f := make([]StripFrame, n)
	for i := range f {
		f[i].At = time.Duration(i) * time.Second
	}
	return f
}

func stripTester(cur *int, frames []StripFrame, disabled bool, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(520).Children(func() {
			if VideoThumbnailStrip(c, cur, frames, VideoThumbnailStripOptions{Disabled: disabled}).Changed() {
				*n++
			}
		})
	}, 600, 160)
}

// stripBuilt counts the time labels the strip built.
func stripBuilt(tt *ui.Tester) int {
	k := 0
	for _, s := range tt.Texts() {
		if strings.Contains(s, ":") && !strings.Contains(s, " ") {
			k++
		}
	}
	return k
}

func TestVideoThumbnailStripKeysAndClick(t *testing.T) {
	cur, n := 0, 0
	tt := stripTester(&cur, stripFrames(40), false, &n)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Frames") {
		t.Fatal("Tab did not reach the frames")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if cur != 2 {
		t.Fatalf("Right twice: %d", cur)
	}
	tt.Key(0, ui.KeyEnd)
	if cur != 39 || !tt.HasText("0:39") {
		t.Fatalf("End: %d, %v", cur, tt.Texts())
	}
	r, _ := tt.Find("0:37")
	tt.ClickAt(r.X+20, r.Y-10)
	if cur != 37 || n != 4 {
		t.Fatalf("click: %d, changes %d", cur, n)
	}
}

func TestVideoThumbnailStripBuildsOnlyVisible(t *testing.T) {
	cur, n := 50000, 0
	tt := stripTester(&cur, stripFrames(100000), false, &n)
	tt.Frame()
	if b := stripBuilt(tt); b < 3 || b > 10 {
		t.Fatalf("built %d frames of 100000", b)
	}
	if !tt.HasText(core.MediaClock(50000 * time.Second)) {
		t.Fatal("current frame not built")
	}
}

func TestVideoThumbnailStripDisabledAndBad(t *testing.T) {
	cur, n := 3, 0
	tt := stripTester(&cur, stripFrames(10), true, &n)
	_ = tt.Click("Next frame")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if cur != 3 || n != 0 {
		t.Fatal("disabled strip moved")
	}
	bad := 5
	mediaPanics(t, "out of range", func(c *ui.Context) {
		VideoThumbnailStrip(c, &bad, stripFrames(3), VideoThumbnailStripOptions{})
	})
}

func TestVideoThumbnailStripSnapshot(t *testing.T) {
	frames := stripFrames(12)
	for i := range frames {
		frames[i].Image = mediaTestBitmap(64, 36)
	}
	cur := 2
	snapshot(t, "video_thumbnail_strip", 560, 110, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(520).Children(func() {
			VideoThumbnailStrip(c, &cur, frames, VideoThumbnailStripOptions{})
		})
	})
}
