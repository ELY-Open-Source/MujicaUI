package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSharedConversationViewShowsAndCopies(t *testing.T) {
	copies := 0
	tt := chatTester(560, 460, func(c *ui.Context) {
		if SharedConversationView(c, chatMessages(), SharedConversationViewOptions{Title: "Banquet seating", SharedBy: "Lady Ann", SharedAt: chatNow}).Changed() {
			copies++
		}
	})
	for _, s := range []string{"Banquet seating", "Shared by Lady Ann", "Read only", "You", "Assistant", "Lady Ann", "And the musicians?"} {
		if !tt.HasText(s) {
			t.Errorf("%q missing from %q", s, tt.Texts())
		}
	}
	_ = tt.Click("Copy link")
	tt.Key(0, ui.KeyEnter)
	if copies != 2 {
		t.Fatalf("copies %d", copies)
	}
}

func TestSharedConversationViewVirtualizes(t *testing.T) {
	msgs := make([]SessionMessage, 100000)
	for i := range msgs {
		msgs[i] = SessionMessage{Role: SessionRole(i % 2), Text: fmt.Sprintf("Message %d", i)}
	}
	tt := chatTester(560, 460, func(c *ui.Context) {
		SharedConversationView(c, msgs, SharedConversationViewOptions{Title: "Long"})
	})
	if n := len(tt.Texts()); n > 80 {
		t.Fatalf("built %d texts for 100000 messages", n)
	}
}

func TestSharedConversationViewDisabledAndChecks(t *testing.T) {
	copies := 0
	tt := chatTester(560, 460, func(c *ui.Context) {
		if SharedConversationView(c, chatMessages(), SharedConversationViewOptions{Title: "T", Disabled: true}).Changed() {
			copies++
		}
	})
	_ = tt.Click("Copy link")
	if copies != 0 || !tt.HasText("Shared by someone") {
		t.Fatalf("copies %d texts %q", copies, tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a missing title")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		SharedConversationView(c, nil, SharedConversationViewOptions{})
	}, 100, 100, 1)
}

func TestSharedConversationViewSnapshot(t *testing.T) {
	snapshot(t, "shared_conversation_view", 560, 380, func(c *ui.Context) {
		ui.Column(c).Padding(16).Height(380).Children(func() {
			SharedConversationView(c, chatMessages(), SharedConversationViewOptions{Title: "Banquet seating", SharedBy: "Lady Ann", SharedAt: chatNow})
		})
	})
}
