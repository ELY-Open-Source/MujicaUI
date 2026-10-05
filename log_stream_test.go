package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func logStreamLines(n int) []LogLine {
	out := make([]LogLine, n)
	for i := range out {
		out[i] = LogLine{Level: LogLevel(i % 4), Text: fmt.Sprintf("\x1b[1mline\x1b[0m %d", i)}
	}
	return out
}

func logStreamTester(lines *[]LogLine, s *LogStreamState, opts LogStreamOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if LogStream(c, *lines, s, opts).Changed() {
			*changes++
		}
	}, 700, 400)
}

func TestLogStreamFollowPauseResume(t *testing.T) {
	lines := logStreamLines(10)
	var s LogStreamState
	n := 0
	tt := logStreamTester(&lines, &s, LogStreamOptions{}, &n)
	if !tt.HasText("line 9") || tt.HasText("\x1b[1mline") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Pause")
	lines = append(lines, logStreamLines(13)[10:]...)
	tt.Frame()
	if !s.Paused || tt.HasText("line 12") || !tt.HasText("3 new lines") {
		t.Fatalf("paused: %q", tt.Texts())
	}
	_ = tt.Click("Resume")
	if s.Paused || !tt.HasText("line 12") || n != 2 {
		t.Fatalf("resume: %q changes %d", tt.Texts(), n)
	}
}

func TestLogStreamLevelFilter(t *testing.T) {
	lines := logStreamLines(8)
	var s LogStreamState
	n := 0
	tt := logStreamTester(&lines, &s, LogStreamOptions{}, &n)
	_ = tt.Click("WARN")
	if s.MinLevel != LogWarn || tt.HasText("line 5") || !tt.HasText("line 6") || !tt.HasText("line 7") {
		t.Fatalf("filter: %v %q", s.MinLevel, tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	if s.MinLevel != LogError || tt.HasText("line 6") {
		t.Fatalf("keyboard filter: %v", s.MinLevel)
	}
}

func TestLogStreamDisabledEmptyMany(t *testing.T) {
	lines := logStreamLines(4)
	var s LogStreamState
	n := 0
	tt := logStreamTester(&lines, &s, LogStreamOptions{Disabled: true}, &n)
	_ = tt.Click("Pause")
	_ = tt.Click("WARN")
	if n != 0 || s.Paused || s.MinLevel != LogDebug {
		t.Fatalf("disabled changed: %+v", s)
	}
	none := []LogLine{}
	tt = logStreamTester(&none, &LogStreamState{}, LogStreamOptions{}, &n)
	if !tt.HasText("No log lines") {
		t.Fatal("empty text missing")
	}
	many := logStreamLines(100000)
	var ms LogStreamState
	tt = logStreamTester(&many, &ms, LogStreamOptions{}, &n)
	if got := len(tt.Texts()); got > 300 || !tt.HasText("line 99999") {
		t.Fatalf("%d texts; last line shown %v", got, tt.HasText("line 99999"))
	}
}

func TestLogStreamSnapshot(t *testing.T) {
	lines := []LogLine{
		{Level: LogInfo, Text: "server listening on :8080"},
		{Level: LogDebug, Text: "cache warm: 312 keys"},
		{Level: LogWarn, Text: "slow query \x1b[33m412 ms\x1b[0m on guests"},
		{Level: LogError, Text: "\x1b[31mpayment provider unreachable\x1b[0m"},
		{Level: LogInfo, Text: "retry 1 of 3 succeeded"},
	}
	snapshot(t, "log_stream", 640, 260, func(c *ui.Context) {
		var s LogStreamState
		ui.Column(c).Padding(16).Children(func() { LogStream(c, lines, &s, LogStreamOptions{Height: 220}) })
	})
}
