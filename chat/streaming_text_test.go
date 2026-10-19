package chat

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func streamingTester(text *string, opts StreamingTextOptions, reduce bool) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		StreamingText(c, *text, opts)
	}, 400, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: reduce, TextScale: 1})
	tt.Frame()
	return tt
}

// streamingShown returns the visible text of the stream.
func streamingShown(tt *ui.Tester) string {
	texts := tt.Texts()
	return texts[len(texts)-1]
}

func TestStreamingTextSmoothReveal(t *testing.T) {
	text := "Hello"
	tt := streamingTester(&text, StreamingTextOptions{Smooth: true, Streaming: true, Rate: 1}, false)
	if !strings.HasPrefix(streamingShown(tt), "Hello") || !strings.HasSuffix(streamingShown(tt), "▍") {
		t.Fatalf("first text %q", streamingShown(tt))
	}
	text += ", court"
	tt.Frame()
	if s := streamingShown(tt); strings.Contains(s, "court") {
		t.Fatalf("1 char/s revealed %q at once", s)
	}
}

func TestStreamingTextInstantAndReducedMotion(t *testing.T) {
	for _, tc := range []struct {
		opts   StreamingTextOptions
		reduce bool
	}{{StreamingTextOptions{}, false}, {StreamingTextOptions{Smooth: true, Rate: 1}, true}} {
		text := "Hello"
		tt := streamingTester(&text, tc.opts, tc.reduce)
		text += ", court"
		tt.Frame()
		if s := streamingShown(tt); s != "Hello, court" {
			t.Errorf("%+v reduce %v: shows %q", tc.opts, tc.reduce, s)
		}
	}
}

func TestStreamingTextReplacedRestarts(t *testing.T) {
	text := "abcdef"
	tt := streamingTester(&text, StreamingTextOptions{Smooth: true, Rate: 1}, false)
	text = "abXYZ"
	tt.Frame()
	if s := streamingShown(tt); s != "ab" {
		t.Fatalf("replaced text shows %q, want the shared prefix", s)
	}
}

func TestStreamingTextNegativeRate(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a negative rate")
		}
	}()
	text := ""
	streamingTester(&text, StreamingTextOptions{Rate: -1}, false)
}

func TestStreamingTextSnapshot(t *testing.T) {
	snapshot(t, "streaming_text", 360, 90, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			StreamingText(c, "The herald reads the decree, word by word", StreamingTextOptions{Streaming: true, Smooth: true})
		})
	})
}
