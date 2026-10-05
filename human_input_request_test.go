package mujica

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func humanTester(answer *string, opts *HumanInputRequestOptions, got *[]HumanInputResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := HumanInputRequest(c, answer, *opts)
		if r.Submitted || r.TimedOut {
			*got = append(*got, r)
		}
	}, 560, 400)
}

func TestHumanInputRequestOptionAndFree(t *testing.T) {
	answer := ""
	opts := HumanInputRequestOptions{Question: "Which hall for the ball?", Options: []string{"East wing", "Great hall"}, AllowFree: true}
	var got []HumanInputResult
	tt := humanTester(&answer, &opts, &got)
	_ = tt.Click("Send answer") // disabled while blank
	if len(got) != 0 {
		t.Fatal("blank answer sent")
	}
	_ = tt.Click("Great hall")
	if answer != "Great hall" {
		t.Fatalf("option: %q", answer)
	}
	_ = tt.Click("Or write your own answer")
	tt.Type("The garden")
	if answer != "The garden" {
		t.Fatalf("free: %q", answer)
	}
	tt.Key(0, ui.KeyEnter)
	if len(got) != 1 || !got[0].Submitted {
		t.Fatalf("Enter: %+v", got)
	}
}

func TestHumanInputRequestDeadline(t *testing.T) {
	answer := "East wing"
	opts := HumanInputRequestOptions{Question: "Which hall?", Options: []string{"East wing"}, Deadline: time.Now().Add(75 * time.Second)}
	var got []HumanInputResult
	tt := humanTester(&answer, &opts, &got)
	found := false
	for _, s := range tt.Texts() {
		found = found || strings.HasPrefix(s, "Answer within 1 m 1")
	}
	if !found {
		t.Fatalf("countdown missing: %q", tt.Texts())
	}
	opts.Deadline = time.Now().Add(-time.Second)
	tt.Frame()
	tt.Frame()
	if len(got) != 1 || !got[0].TimedOut || !tt.HasText("The agent stopped waiting") {
		t.Fatalf("timeout: %+v", got)
	}
	_ = tt.Click("Send answer")
	if len(got) != 1 {
		t.Fatal("answer sent after the deadline")
	}
}

func TestHumanInputRequestDisabledAndRules(t *testing.T) {
	answer := ""
	opts := HumanInputRequestOptions{Question: "Which hall?", Options: []string{"East wing"}, Disabled: true}
	var got []HumanInputResult
	tt := humanTester(&answer, &opts, &got)
	_ = tt.Click("East wing")
	if answer != "" {
		t.Fatal("disabled request changed the answer")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no options and no free answer did not panic")
		}
	}()
	humanTester(&answer, &HumanInputRequestOptions{Question: "x"}, &got)
}

func TestHumanInputRequestSnapshot(t *testing.T) {
	snapshot(t, "human_input_request", 520, 300, func(c *ui.Context) {
		answer := "Great hall"
		ui.Column(c).Padding(16).Children(func() {
			HumanInputRequest(c, &answer, HumanInputRequestOptions{Question: "Which hall should host the masquerade?",
				Options: []string{"East wing", "Great hall"}, AllowFree: true})
		})
	})
}
