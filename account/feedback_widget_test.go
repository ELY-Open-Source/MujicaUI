package account

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func TestFeedbackWidgetValidatesSubmitsOnce(t *testing.T) {
	var s FeedbackWidgetState
	var sent []Feedback
	attached := false
	tt := accountTester(500, 500, func(c *ui.Context) {
		r := FeedbackWidget(c, &s, FeedbackWidgetOptions{Attachments: func() { attached = true; ui.Text(c, "screenshot slot") }})
		if f, ok := r.Submitted(); ok {
			sent = append(sent, f)
		}
	})
	if !attached || !tt.HasText("screenshot slot") {
		t.Fatal("attachment slot not built")
	}
	_ = tt.Click("Idea")
	_ = tt.Click("Message")
	tt.Type("short")
	_ = tt.Click("Send feedback")
	if len(sent) != 0 || !tt.HasText("Write at least 10 characters.") {
		t.Fatalf("short sent: %q", tt.Texts())
	}
	_ = tt.Click("Message")
	tt.Type(" but now long enough")
	_ = tt.Click("Send feedback")
	_ = tt.Click("Send feedback")
	if len(sent) != 1 || sent[0].Kind != FeedbackIdea || sent[0].Text != "short but now long enough" {
		t.Fatalf("sent %+v", sent)
	}
	s.EndSubmit("")
	tt.Frame()
	if s.Draft.Text != "" {
		t.Fatal("sent draft kept")
	}
}

func TestFeedbackWidgetDisabled(t *testing.T) {
	s := FeedbackWidgetState{Draft: Feedback{Text: "A long enough message."}}
	sent := 0
	tt := accountTester(500, 500, func(c *ui.Context) {
		if _, ok := FeedbackWidget(c, &s, FeedbackWidgetOptions{Disabled: true}).Submitted(); ok {
			sent++
		}
	})
	_ = tt.Click("Bug")
	_ = tt.Click("Praise")
	_ = tt.Click("Send feedback")
	if sent != 0 || s.Draft.Kind != FeedbackBug {
		t.Fatal("disabled widget acted")
	}
}

func TestFeedbackWidgetSnapshot(t *testing.T) {
	snapshot(t, "feedback_widget", 460, 400, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() FeedbackWidgetState {
			return FeedbackWidgetState{Draft: Feedback{Kind: FeedbackBug, Text: "The ticker stutters on resize."}}
		})
		ui.Column(c).Padding(16).Children(func() {
			FeedbackWidget(c, s, FeedbackWidgetOptions{Attachments: func() {
				input.Button(c, "Attach screenshot", input.ButtonOptions{Variant: input.Secondary})
			}})
		})
	})
}

func TestFeedbackWidgetPendingFreezesDraftAndAttachments(t *testing.T) {
	s := FeedbackWidgetState{Draft: Feedback{Kind: FeedbackIdea, Text: "A long enough message."}}
	attachments := 0
	tt := accountTester(500, 500, func(c *ui.Context) {
		FeedbackWidget(c, &s, FeedbackWidgetOptions{Attachments: func() {
			if input.Button(c, "Attach", input.ButtonOptions{}).Clicked() {
				attachments++
			}
		}})
	})
	if err := tt.Click("Send feedback"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Praise"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Message"); err != nil {
		t.Fatal(err)
	}
	tt.Type("modified")
	if err := tt.Click("Attach"); err != nil {
		t.Fatal(err)
	}
	if s.Draft.Kind != FeedbackIdea || s.Draft.Text != "A long enough message." || attachments != 0 {
		t.Fatal("pending feedback changes its draft or attachments")
	}
}

func TestFeedbackWidgetFailurePreservesDraftAndInvalidLength(t *testing.T) {
	s := FeedbackWidgetState{Draft: Feedback{Kind: FeedbackIdea, Text: "A long enough message."}}
	tt := accountTester(600, 500, func(c *ui.Context) { FeedbackWidget(c, &s, FeedbackWidgetOptions{}) })
	if err := tt.Click("Send feedback"); err != nil {
		t.Fatal(err)
	}
	s.EndSubmit("Send failed.")
	tt.Frame()
	if s.Draft.Text != "A long enough message." || !tt.HasText("Send failed.") {
		t.Fatal("failed feedback loses its draft or message")
	}
	accountPanics(t, func(c *ui.Context) { FeedbackWidget(c, &s, FeedbackWidgetOptions{MinLength: 2001}) })
}

func TestFeedbackWidgetDisabledKeepsFormState(t *testing.T) {
	s := FeedbackWidgetState{Draft: Feedback{Text: "A long enough message."}}
	s.Form.Submit()
	accountTester(600, 500, func(c *ui.Context) { FeedbackWidget(c, &s, FeedbackWidgetOptions{Disabled: true}) })
	if s.Form.Submitting() || !s.Form.Requested || s.Form.Fields != nil {
		t.Fatal("disabled frame edits form state")
	}
}
