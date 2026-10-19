package code

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func logViewerLines(n int) []LogLine {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	lines := make([]LogLine, n)
	for i := range lines {
		lines[i] = LogLine{Time: at.Add(time.Duration(i) * time.Millisecond), Level: LogLevel(i % 4), Text: fmt.Sprintf("event %d", i)}
	}
	return lines
}

func logViewerTester(lines *[]LogLine, s *LogViewerState, disabled bool) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		LogViewer(c, *lines, s, LogViewerOptions{Disabled: disabled}).Element.Height(400)
	}, 760, 420)
}

func logViewerRows(tt *ui.Tester) int {
	n := 0
	for _, s := range tt.Texts() {
		if strings.HasPrefix(s, "event ") {
			n++
		}
	}
	return n
}

func TestLogViewerHundredThousandLines(t *testing.T) {
	lines := logViewerLines(100_000)
	s := &LogViewerState{}
	tt := logViewerTester(&lines, s, false)
	if !tt.HasText("event 99999") {
		t.Fatal("follow did not start at the tail")
	}
	if n := logViewerRows(tt); n == 0 || n > 40 {
		t.Fatalf("built %d rows of 100000", n)
	}
	// Hiding a level refilters every line once.
	if err := tt.Click("ERROR"); err != nil {
		t.Fatal(err)
	}
	if len(s.shown) != 75_000 || tt.HasText("event 99999") {
		t.Fatalf("shown %d after hiding ERROR", len(s.shown))
	}
	if n := logViewerRows(tt); n > 40 {
		t.Fatalf("built %d rows after filtering", n)
	}
}

func TestLogViewerFollowsAppends(t *testing.T) {
	lines := logViewerLines(50)
	s := &LogViewerState{}
	tt := logViewerTester(&lines, s, false)
	lines = append(lines, LogLine{Level: LogInfo, Text: "event fresh"})
	tt.Frame()
	if !tt.HasText("event fresh") || len(s.shown) != 51 {
		t.Fatalf("appended line not followed; shown %d", len(s.shown))
	}
	if err := tt.Click("Follow"); err != nil {
		t.Fatal(err)
	}
	if s.Follow {
		t.Fatal("Follow toggle did not turn off")
	}
	tt.Scroll(300, 200, 0, -2000)
	lines = append(lines, LogLine{Level: LogInfo, Text: "event later"})
	tt.Frame()
	if tt.HasText("event later") {
		t.Fatal("paused viewer jumped to the new line")
	}
}

func TestLogViewerHighlight(t *testing.T) {
	lines := []LogLine{{Level: LogWarn, Text: "disk nearly full"}}
	s := &LogViewerState{Query: "FULL"}
	tt := logViewerTester(&lines, s, false)
	if !tt.HasText("disk nearly full") || !tt.HasText("WARN") {
		t.Fatalf("texts %q", tt.Texts())
	}
	var spans []ui.Span
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		spans = CodeMatchSpans(c, "disk nearly full", "FULL")
	}, 10, 10)
	if len(spans) != 2 || spans[1].Text != "full" || !spans[1].Underline || spans[1].Background.A == 0 {
		t.Fatalf("highlight spans %+v", spans)
	}
}

func TestLogViewerDisabledAndEmpty(t *testing.T) {
	lines := logViewerLines(8)
	s := &LogViewerState{}
	tt := logViewerTester(&lines, s, true)
	_ = tt.Click("ERROR")
	if s.Hidden[LogError] {
		t.Fatal("disabled viewer changed its filter")
	}
	var none []LogLine
	tt = logViewerTester(&none, &LogViewerState{}, false)
	if !tt.HasText("No log lines") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestLogViewerSnapshot(t *testing.T) {
	snapshot(t, "log_viewer", 640, 230, func(c *ui.Context) {
		s := ui.Local(c.Root(), "snap", func() *LogViewerState { return &LogViewerState{Query: "seal"} })
		at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
		lines := []LogLine{
			{at, LogDebug, "loading ledger index"},
			{at.Add(12 * time.Millisecond), LogInfo, "court opened by the chancellor"},
			{at.Add(40 * time.Millisecond), LogWarn, "seal of decree 41 is faded"},
			{at.Add(91 * time.Millisecond), LogError, "seal of decree 42 missing"},
		}
		ui.Column(c).Padding(12).Children(func() {
			LogViewer(c, lines, *s, LogViewerOptions{}).Element.Height(206)
		})
	})
}
