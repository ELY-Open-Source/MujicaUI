package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestChatEmptyStateKinds(t *testing.T) {
	for _, tc := range []struct {
		opts   ChatEmptyStateOptions
		texts  []string
		action string
	}{
		{ChatEmptyStateOptions{}, []string{"Start a conversation"}, "New conversation"},
		{ChatEmptyStateOptions{Kind: ChatEmptyNoResults, Query: "minuet"}, []string{"minuet"}, "Clear search"},
		{ChatEmptyStateOptions{Kind: ChatEmptyOffline}, []string{"You are offline"}, "Retry"},
		{ChatEmptyStateOptions{Kind: ChatEmptyOffline, Title: "Lost the court", Action: "Reconnect"}, []string{"Lost the court"}, "Reconnect"},
	} {
		acted := 0
		tt := chatTester(420, 320, func(c *ui.Context) {
			if ChatEmptyState(c, tc.opts).Action {
				acted++
			}
		})
		for _, s := range tc.texts {
			if !tt.HasText(s) {
				t.Errorf("%+v: %q missing from %q", tc.opts, s, tt.Texts())
			}
		}
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyEnter)
		if err := tt.Click(tc.action); err != nil || acted != 2 {
			t.Errorf("%+v: action %q acted %d (%v)", tc.opts, tc.action, acted, err)
		}
	}
}

func TestChatEmptyStateRejectsBadKind(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ChatEmptyState(c, ChatEmptyStateOptions{Kind: 9})
	}, 100, 100, 1)
}

func TestChatEmptyStateSnapshot(t *testing.T) {
	snapshot(t, "chat_empty_state", 820, 240, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(8).Children(func() {
			for _, k := range []ChatEmptyKind{ChatEmptyNew, ChatEmptyNoResults, ChatEmptyOffline} {
				ui.Column(c).Width(260).Children(func() { ChatEmptyState(c, ChatEmptyStateOptions{Kind: k, Query: "minuet"}) })
			}
		})
	})
}
