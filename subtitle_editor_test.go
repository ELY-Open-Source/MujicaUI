package mujica

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func subtitleTester(cues *[]SubtitleCue, disabled bool, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(20).Width(620).Children(func() {
			if SubtitleEditor(c, cues, SubtitleEditorOptions{Height: 260, Disabled: disabled}).Changed() {
				*n++
			}
		})
	}, 700, 400)
}

func TestCheckSubtitles(t *testing.T) {
	cues := []SubtitleCue{{0, 2 * time.Second, "a"}, {time.Second, 3 * time.Second, "b"}, {4 * time.Second, 4 * time.Second, "c"}}
	want := []SubtitleIssue{{1, SubtitleOverlap}, {2, SubtitleEndBeforeStart}}
	if got := CheckSubtitles(cues); !slices.Equal(got, want) {
		t.Fatalf("issues %v", got)
	}
}

func TestSubtitleParse(t *testing.T) {
	for s, want := range map[string]time.Duration{"0:01.5": 1500 * time.Millisecond, "1:05.250": 65250 * time.Millisecond, "1:00:00": time.Hour, "2:03": 123 * time.Second} {
		if got, ok := subtitleParse(s); !ok || got != want {
			t.Errorf("%q: %v %v", s, got, ok)
		}
	}
	for _, s := range []string{"", "5", "1:60", "1:5", "a:00", "1:00.1234", "1:-1", "+1:00"} {
		if _, ok := subtitleParse(s); ok {
			t.Errorf("%q parsed", s)
		}
	}
}

func TestSubtitleEditorEditTimesAndText(t *testing.T) {
	cues := []SubtitleCue{{0, 2 * time.Second, "Welcome"}, {3 * time.Second, 5 * time.Second, "to the court"}}
	n := 0
	tt := subtitleTester(&cues, false, &n)
	if err := tt.Click("End 1"); err != nil {
		t.Fatal(err)
	}
	tt.Command("selectAll")
	tt.Type("0:04")
	if cues[0].End != 4*time.Second || !tt.HasText("Overlaps cue 1") {
		t.Fatalf("end edit %+v, texts %v", cues[0], tt.Texts())
	}
	tt.Command("selectAll")
	tt.Type("zz")
	if cues[0].End != 4*time.Second || !tt.HasText("Use m:ss.mmm, as 1:05.250") {
		t.Fatalf("bad time kept? %+v", cues[0])
	}
	if err := tt.Click("Text 2"); err != nil {
		t.Fatal(err)
	}
	tt.Type("!")
	if cues[1].Text != "to the court!" || n < 2 {
		t.Fatalf("text %q changes %d", cues[1].Text, n)
	}
}

func TestSubtitleEditorAddRemove(t *testing.T) {
	cues := []SubtitleCue{{0, 2 * time.Second, "a"}}
	n := 0
	tt := subtitleTester(&cues, false, &n)
	if err := tt.Click("Add cue"); err != nil {
		t.Fatal(err)
	}
	if len(cues) != 2 || cues[1].Start != 2*time.Second || cues[1].End != 4*time.Second {
		t.Fatalf("added %+v", cues)
	}
	if err := tt.Click("Remove cue 1"); err != nil {
		t.Fatal(err)
	}
	if len(cues) != 1 || cues[0].Start != 2*time.Second || n != 2 {
		t.Fatalf("removed %+v, changes %d", cues, n)
	}
}

func TestSubtitleEditorBuildsOnlyVisible(t *testing.T) {
	cues := make([]SubtitleCue, 100000)
	for i := range cues {
		cues[i] = SubtitleCue{time.Duration(i) * time.Second, time.Duration(i)*time.Second + 900*time.Millisecond, "line"}
	}
	n := 0
	tt := subtitleTester(&cues, false, &n)
	tt.Frame()
	built := 0
	for _, x := range tt.Texts() {
		if strings.HasPrefix(x, "Remove cue ") {
			built++
		}
	}
	if built == 0 || built > 20 {
		t.Fatalf("built %d rows of 100000", built)
	}
}

func TestSubtitleEditorDisabled(t *testing.T) {
	cues := []SubtitleCue{{0, 2 * time.Second, "a"}}
	n := 0
	tt := subtitleTester(&cues, true, &n)
	_ = tt.Click("Add cue")
	_ = tt.Click("Remove cue 1")
	_ = tt.Click("Text 1")
	tt.Type("x")
	if len(cues) != 1 || cues[0].Text != "a" || n != 0 {
		t.Fatalf("disabled editor changed %+v", cues)
	}
}

func TestSubtitleEditorSnapshot(t *testing.T) {
	cues := []SubtitleCue{
		{0, 2500 * time.Millisecond, "The curtain rises."},
		{2 * time.Second, 5 * time.Second, "Welcome to the masquerade."},
		{6 * time.Second, 9 * time.Second, "Ave Mujica."},
	}
	snapshot(t, "subtitle_editor", 640, 300, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(600).Children(func() {
			SubtitleEditor(c, &cues, SubtitleEditorOptions{Height: 200})
		})
	})
}
