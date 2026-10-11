package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestMessageActionsByRole(t *testing.T) {
	for role, want := range map[MessageRole][]string{
		MessageUser:      {"Copy", "Edit"},
		MessageAssistant: {"Copy", "Regenerate", "Good response", "Bad response", "Read aloud", "Share"},
		MessageSystem:    {"Copy"},
		MessageError:     {"Copy", "Regenerate"},
	} {
		fb := FeedbackNone
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			MessageActions(c, role, MessageActionsOptions{Feedback: &fb})
		}, 400, 100)
		for _, name := range want {
			if _, ok := tt.Find(name); !ok {
				t.Errorf("role %v: missing %q", role, name)
			}
		}
		if _, ok := tt.Find("Edit"); ok != (role == MessageUser) {
			t.Errorf("role %v: edit shown %v", role, ok)
		}
	}
}

func TestMessageActionsChooseAndRate(t *testing.T) {
	fb := FeedbackNone
	var got []MessageAction
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if a, ok := MessageActions(c, MessageAssistant, MessageActionsOptions{Feedback: &fb}).Action(); ok {
			got = append(got, a)
		}
	}, 400, 100)
	_ = tt.Click("Regenerate")
	_ = tt.Click("Good response")
	if fb != FeedbackLiked {
		t.Fatalf("feedback %v after like", fb)
	}
	_ = tt.Click("Bad response")
	if fb != FeedbackDisliked {
		t.Fatalf("feedback %v after dislike", fb)
	}
	_ = tt.Click("Bad response")
	if fb != FeedbackNone {
		t.Fatalf("feedback %v after undoing dislike", fb)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	// Tab returns to the last focused action, Dislike; Right moves to Read aloud.
	want := []MessageAction{ActionRegenerate, ActionLike, ActionDislike, ActionDislike, ActionSpeak}
	if len(got) != len(want) {
		t.Fatalf("actions %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("actions %v, want %v", got, want)
		}
	}
}

func TestMessageActionsHideRatingAndDisable(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(true).Children(func() {
			if _, ok := MessageActions(c, MessageAssistant, MessageActionsOptions{}).Action(); ok {
				n++
			}
		})
	}, 400, 100)
	if _, ok := tt.Find("Good response"); ok {
		t.Fatal("rating shown without Feedback")
	}
	_ = tt.Click("Copy")
	if n != 0 {
		t.Fatal("disabled actions fired")
	}
}

func TestMessageActionsSnapshot(t *testing.T) {
	fb := FeedbackLiked
	snapshot(t, "message_actions", 300, 70, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { MessageActions(c, MessageAssistant, MessageActionsOptions{Feedback: &fb}) })
	})
}
