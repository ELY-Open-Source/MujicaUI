package agent

import (
	"fmt"
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

// screenshotFrames makes n solid frames of shifting colors.
func screenshotFrames(n int) []ScreenshotFrame {
	out := make([]ScreenshotFrame, n)
	for i := range out {
		img := image.NewRGBA(image.Rect(0, 0, 160, 90))
		for p := 0; p < len(img.Pix); p += 4 {
			img.Pix[p], img.Pix[p+1], img.Pix[p+2], img.Pix[p+3] = uint8(60+i*40), 90, uint8(160-i*20), 255
		}
		img.Set(10, 10, color.White)
		out[i] = ScreenshotFrame{Image: ui.NewBitmap(img), Caption: fmt.Sprintf("Step %d", i+1), At: time.Duration(i) * time.Second}
	}
	return out
}

func screenshotTester(st *ScreenshotStreamState, frames *[]ScreenshotFrame, opts ScreenshotStreamOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ScreenshotStream(c, st, *frames, opts).Changed() {
			*changes++
		}
	}, 600, 420)
}

func TestScreenshotStreamFollowsPausesAndSteps(t *testing.T) {
	var st ScreenshotStreamState
	frames := screenshotFrames(3)
	changes := 0
	tt := screenshotTester(&st, &frames, ScreenshotStreamOptions{}, &changes)
	if st.Index != 2 || !tt.HasText("Frame 3 of 3") || !tt.HasText("Live") {
		t.Fatalf("live: %+v %q", st, tt.Texts())
	}
	frames = append(frames, screenshotFrames(4)[3])
	tt.Frame()
	if st.Index != 3 || !tt.HasText("Step 4") {
		t.Fatalf("new frame not followed: %+v", st)
	}
	_ = tt.Click("Previous frame")
	if !st.Paused || st.Index != 2 || changes != 1 || !tt.HasText("Step 3") {
		t.Fatalf("step back: %+v", st)
	}
	frames = append(frames, screenshotFrames(5)[4])
	tt.Frame()
	if st.Index != 2 {
		t.Fatal("paused viewer moved")
	}
	_ = tt.Click("Go live")
	if st.Paused || st.Index != 4 {
		t.Fatalf("go live: %+v", st)
	}
}

func TestScreenshotStreamKeyboard(t *testing.T) {
	var st ScreenshotStreamState
	frames := screenshotFrames(3)
	n := 0
	tt := screenshotTester(&st, &frames, ScreenshotStreamOptions{}, &n)
	tt.Key(0, ui.KeyTab) // the viewport
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft) // stops at the first
	if st.Index != 0 || !st.Paused {
		t.Fatalf("left: %+v", st)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	if st.Paused || st.Index != 2 {
		t.Fatalf("space: %+v", st)
	}
}

func TestScreenshotStreamEmptyAndDisabled(t *testing.T) {
	var st ScreenshotStreamState
	var frames []ScreenshotFrame
	n := 0
	tt := screenshotTester(&st, &frames, ScreenshotStreamOptions{}, &n)
	if !tt.HasText("Waiting for frames") || !tt.HasText("Frame 0 of 0") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	frames = screenshotFrames(3)
	st = ScreenshotStreamState{}
	tt = screenshotTester(&st, &frames, ScreenshotStreamOptions{Disabled: true}, &n)
	_ = tt.Click("Previous frame")
	_ = tt.Click("Pause")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyLeft)
	if n != 0 || st.Paused || st.Index != 2 {
		t.Fatalf("disabled moved: %+v", st)
	}
}

func TestScreenshotStreamSnapshot(t *testing.T) {
	snapshot(t, "screenshot_stream", 520, 330, func(c *ui.Context) {
		st := ScreenshotStreamState{Index: 1, Paused: true}
		ui.Column(c).Padding(16).Children(func() {
			ScreenshotStream(c, &st, screenshotFrames(3), ScreenshotStreamOptions{Height: 200})
		})
	})
}
