package mujica

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestThinkingIndicatorPhase(t *testing.T) {
	tt := chatTester(400, 100, func(c *ui.Context) { ThinkingIndicator(c, ThinkingIndicatorOptions{}) })
	if !tt.HasText("Thinking") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt = chatTester(400, 100, func(c *ui.Context) {
		ThinkingIndicator(c, ThinkingIndicatorOptions{Phase: "Reading the archives"})
	})
	if !tt.HasText("Reading the archives") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestThinkingIndicatorReducedMotionIsStill(t *testing.T) {
	frame := func(reduce bool, at time.Duration) []byte {
		tt := chatTester(200, 60, func(c *ui.Context) { ThinkingIndicator(c, ThinkingIndicatorOptions{}) })
		tt.SetPreferences(ui.Preferences{ReduceMotion: reduce, TextScale: 1})
		tt.Frame()
		if at > 0 {
			time.Sleep(at)
			tt.Frame()
		}
		return encode(t, tt.Image())
	}
	if !bytes.Equal(frame(true, 0), frame(true, 300*time.Millisecond)) {
		t.Fatal("reduced motion still moves")
	}
	tt := chatTester(200, 60, func(c *ui.Context) { ThinkingIndicator(c, ThinkingIndicatorOptions{}) })
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if !tt.HasText("Thinking…") {
		t.Fatalf("static label missing: %q", tt.Texts())
	}
}

func TestThinkingIndicatorSnapshot(t *testing.T) {
	snapshot(t, "thinking_indicator", 360, 70, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ThinkingIndicator(c, ThinkingIndicatorOptions{Phase: "Searching the archives"})
		})
	})
}
