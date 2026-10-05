package mujica

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestFeedbackFormReasonsAndSubmit(t *testing.T) {
	var v FeedbackValue
	var last FeedbackFormView
	tt := chatTester(480, 420, func(c *ui.Context) {
		if r := FeedbackForm(c, &v, FeedbackFormOptions{Rating: FeedbackDown}); r.Submitted || r.Cancelled {
			last = r
		}
	})
	if !tt.HasText("What went wrong?") || !tt.HasText("Harmful or unsafe") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Incorrect")
	_ = tt.Click("Not helpful")
	_ = tt.Click("Incorrect")
	if !slices.Equal(v.Reasons, []string{"Not helpful"}) {
		t.Fatalf("reasons %q", v.Reasons)
	}
	_ = tt.Click("Tell us more (optional)")
	tt.Type("Wrong date")
	_ = tt.Click("Submit")
	if !last.Submitted || v.Comment != "Wrong date" {
		t.Fatalf("submit %+v comment %q", last, v.Comment)
	}
}

func TestFeedbackFormCancelAndUpReasons(t *testing.T) {
	var v FeedbackValue
	var last FeedbackFormView
	tt := chatTester(480, 420, func(c *ui.Context) {
		if r := FeedbackForm(c, &v, FeedbackFormOptions{Reasons: []string{"Elegant"}}); r.Submitted || r.Cancelled {
			last = r
		}
	})
	if !tt.HasText("What did you like?") || !tt.HasText("Elegant") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEscape)
	if !last.Cancelled {
		t.Fatal("Escape did not cancel")
	}
}

func TestFeedbackFormBlocked(t *testing.T) {
	for _, opts := range []FeedbackFormOptions{{Disabled: true}, {Submitting: true}} {
		var v FeedbackValue
		submitted := false
		tt := chatTester(480, 420, func(c *ui.Context) {
			if FeedbackForm(c, &v, opts).Submitted {
				submitted = true
			}
		})
		_ = tt.Click("Accurate")
		_ = tt.Click("Submit")
		if submitted || len(v.Reasons) != 0 {
			t.Fatalf("%+v: blocked form acted", opts)
		}
	}
}

func TestFeedbackFormSnapshot(t *testing.T) {
	snapshot(t, "feedback_form", 480, 300, func(c *ui.Context) {
		v := FeedbackValue{Reasons: []string{"Incorrect"}}
		ui.Column(c).Padding(16).Children(func() { FeedbackForm(c, &v, FeedbackFormOptions{Rating: FeedbackDown}) })
	})
}
