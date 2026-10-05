package messaging

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

func TestUserProfileCard(t *testing.T) {
	clicked := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		UserProfileCard(c, UserProfile{Name: "Sakiko Togawa", Presence: PresenceBusy, Role: "Keyboard", Status: "In rehearsal"},
			UserProfileCardOptions{Actions: func() {
				if input.Button(c, "Message", input.ButtonOptions{}).Clicked() {
					clicked++
				}
			}})
	}, 400, 300)
	for _, s := range []string{"Sakiko Togawa", "Busy", "Keyboard", "In rehearsal"} {
		if !tt.HasText(s) {
			t.Fatalf("no %q in %v", s, tt.Texts())
		}
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if clicked != 1 {
		t.Fatal("action slot not reachable by keyboard")
	}
	mediaPanics(t, "no name", func(c *ui.Context) { UserProfileCard(c, UserProfile{}, UserProfileCardOptions{}) })
}

func TestUserProfileCardSnapshot(t *testing.T) {
	snapshot(t, "user_profile_card", 360, 260, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			UserProfileCard(c, UserProfile{Name: "Sakiko Togawa", Presence: PresenceOnline, Role: "Keyboard", Status: "Composing the second act",
				About: "Founder of the masquerade."}, UserProfileCardOptions{Actions: func() {
				input.Button(c, "Message", input.ButtonOptions{})
				input.Button(c, "Call", input.ButtonOptions{Variant: input.Secondary})
			}})
		})
	})
}
