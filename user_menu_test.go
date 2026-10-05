package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestUserMenuStatusSettingsSignOut(t *testing.T) {
	status := UserStatusOnline
	var changes, settings, leaves int
	tt := accountTester(500, 400, func(c *ui.Context) {
		r := UserMenu(c, &status, UserMenuOptions{Name: "Sakiko Togawa", Email: "sakiko@court"})
		if r.StatusChanged() {
			changes++
		}
		if r.Settings() {
			settings++
		}
		if r.SignOut() {
			leaves++
		}
	})
	if err := tt.Click("Sakiko Togawa, Online"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Sakiko Togawa · sakiko@court") {
		t.Fatalf("header: %q", tt.Texts())
	}
	_ = tt.Click("Do not disturb")
	if status != UserStatusBusy || changes != 1 {
		t.Fatalf("status %v", status)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	_ = tt.Click("Settings")
	_ = tt.Click("Sakiko Togawa, Do not disturb")
	_ = tt.Click("Sign out")
	if settings != 1 || leaves != 1 {
		t.Fatalf("settings %d sign out %d", settings, leaves)
	}
}

func TestUserMenuSnapshot(t *testing.T) {
	snapshot(t, "user_menu", 220, 90, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(16).Children(func() {
			for s := UserStatusOnline; s <= UserStatusInvisible; s++ {
				st := s
				ui.Box(c).Key(int(s)).Children(func() { UserMenu(c, &st, UserMenuOptions{Name: "Sakiko Togawa"}) })
			}
		})
	})
}

func TestUserMenuDisabledAndInvalidStatus(t *testing.T) {
	status, disabled, actions := UserStatusOnline, false, 0
	tt := accountTester(600, 500, func(c *ui.Context) {
		ui.Column(c).Disabled(disabled).Children(func() {
			r := UserMenu(c, &status, UserMenuOptions{Name: "User"})
			if r.StatusChanged() || r.Settings() || r.SignOut() {
				actions++
			}
		})
	})
	if err := tt.Click("User, Online"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if status != UserStatusOnline || actions != 0 || tt.HasText("Sign out") {
		t.Fatal("disabled user menu acts")
	}
	status = UserStatus(99)
	accountPanics(t, func(c *ui.Context) { UserMenu(c, &status, UserMenuOptions{Name: "User"}) })
}
