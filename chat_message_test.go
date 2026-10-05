package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestChatContinues(t *testing.T) {
	at := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	a := ChatEntry{Author: "Uika", Time: at}
	for _, tc := range []struct {
		b    ChatEntry
		want bool
	}{
		{ChatEntry{Author: "Uika", Time: at.Add(time.Minute)}, true},
		{ChatEntry{Author: "Uika", Time: at.Add(6 * time.Minute)}, false},
		{ChatEntry{Author: "Mutsumi", Time: at.Add(time.Minute)}, false},
		{ChatEntry{Author: "Uika", Time: at.Add(-time.Minute)}, false},
	} {
		if got := ChatContinues(a, tc.b); got != tc.want {
			t.Errorf("%+v: %v", tc.b, got)
		}
	}
}

func TestChatMessageSlotsAndMerge(t *testing.T) {
	at := time.Date(2026, 10, 5, 21, 7, 0, 0, time.UTC)
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(500).Children(func() {
			ChatMessage(c, ChatEntry{ID: "1", Author: "Uika", Time: at, Text: "Rehearsal at nine."}, ChatMessageOptions{
				Reply:     func() { QuoteReply(c, "Sakiko", "When do we start?") },
				Reactions: func() { Tag(c, "👍 2", TagOptions{}) },
			})
			ChatMessage(c, ChatEntry{ID: "2", Author: "Uika", Time: at, Text: "Bring the masks.", Edited: true}, ChatMessageOptions{Continued: true})
		})
	}, 600, 400)
	if mediaTextCount(tt, "U") != 1 || mediaTextCount(tt, "Assistant, Uika") != 1 {
		t.Fatalf("continued message repeated the author: %v", tt.Texts())
	}
	for _, s := range []string{"21:07", "(edited)", "When do we start?", "👍 2"} {
		if !tt.HasText(s) {
			t.Fatalf("no %q in %v", s, tt.Texts())
		}
	}
}

func TestChatMessageSnapshot(t *testing.T) {
	at := time.Date(2026, 10, 5, 21, 7, 0, 0, time.UTC)
	snapshot(t, "chat_message", 520, 340, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(8).Width(480).Children(func() {
			ChatMessage(c, ChatEntry{Author: "Uika Misumi", Time: at, Text: "Rehearsal moves to nine."}, ChatMessageOptions{
				Reply: func() { QuoteReply(c, "Sakiko", "When do we start?") }})
			ChatMessage(c, ChatEntry{Author: "Uika Misumi", Time: at, Text: "Bring the masks.", Edited: true}, ChatMessageOptions{Continued: true})
			ChatMessage(c, ChatEntry{Author: "Me", Time: at, Text: "Understood.", Own: true}, ChatMessageOptions{})
		})
	})
}
