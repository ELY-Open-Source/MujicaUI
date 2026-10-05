package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestOnlineStatusTextAndShape(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Row(c).Gap(12).Children(func() {
			for _, p := range []PresenceStatus{PresenceOnline, PresenceAway, PresenceBusy, PresenceOffline} {
				OnlineStatus(c, p, OnlineStatusOptions{})
			}
			OnlineStatus(c, PresenceBusy, OnlineStatusOptions{HideText: true}).Key("hidden")
		})
	}, 500, 80)
	for _, s := range []string{"Online", "Away", "Offline"} {
		if mediaTextCount(tt, s) != 2 {
			t.Fatalf("%s seen %d times: %v", s, mediaTextCount(tt, s), tt.Texts())
		}
	}
	if mediaTextCount(tt, "Busy") != 3 {
		t.Fatalf("Busy seen %d times", mediaTextCount(tt, "Busy"))
	}
	if _, ok := tt.Find("Busy"); !ok {
		t.Fatal("hidden text not kept as the name")
	}
	mediaPanics(t, "bad presence", func(c *ui.Context) { OnlineStatus(c, PresenceStatus(9), OnlineStatusOptions{}) })
}

func TestOnlineStatusSnapshot(t *testing.T) {
	snapshot(t, "online_status", 420, 60, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(16).Children(func() {
			for _, p := range []PresenceStatus{PresenceOnline, PresenceAway, PresenceBusy, PresenceOffline} {
				OnlineStatus(c, p, OnlineStatusOptions{})
			}
		})
	})
}
