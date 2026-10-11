package account

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

// accountTester runs a themed view inside a padded column.
func accountTester(w, h int, view func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Fill().Children(func() { view(c) })
	}, w, h)
}

var settingsPages = []SettingsPage{{ID: "general", Label: "General"}, {ID: "privacy", Label: "Privacy"}}

func settingsView(c *ui.Context, page *string, sound, track *bool, changes *int) {
	if SettingsLayout(c, page, settingsPages, SettingsLayoutOptions{}, func(id string) {
		switch id {
		case "general":
			SettingsSection(c, "Notifications", func() {
				SettingsRow(c, "Sound", SettingsRowOptions{ID: "sound", Description: "Chime on new messages."}, func() *ui.Element {
					return input.Switch(c, sound, "", input.SwitchOptions{})
				})
			})
		case "privacy":
			SettingsSection(c, "Data", func() {
				SettingsRow(c, "Usage reports", SettingsRowOptions{ID: "track"}, func() *ui.Element {
					return input.Switch(c, track, "", input.SwitchOptions{})
				})
			})
		}
	}).Changed() {
		*changes++
	}
}

func TestSettingsLayoutNavigatesAndLabelsRows(t *testing.T) {
	page, sound, track, n := "", false, false, 0
	tt := accountTester(900, 500, func(c *ui.Context) { settingsView(c, &page, &sound, &track, &n) })
	if page != "general" || !tt.HasText("Chime on new messages.") {
		t.Fatalf("first page: %q", page)
	}
	if err := tt.Click("Chime on new messages."); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeySpace)
	if !tt.Focused("Sound") || !sound {
		t.Fatal("row label does not reach its control")
	}
	if err := tt.Click("Privacy"); err != nil {
		t.Fatal(err)
	}
	if page != "privacy" || n != 1 || !tt.HasText("Usage reports") || tt.HasText("Sound") {
		t.Fatalf("page %q changes %d texts %q", page, n, tt.Texts())
	}
}

func TestSettingsLayoutRejectsUnknownPage(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("unknown page accepted")
		}
	}()
	page := "nope"
	accountTester(400, 300, func(c *ui.Context) { SettingsLayout(c, &page, settingsPages, SettingsLayoutOptions{}, func(string) {}) })
}

func TestSettingsLayoutSnapshot(t *testing.T) {
	snapshot(t, "settings_layout", 860, 360, func(c *ui.Context) {
		st := ui.Local(c.Root(), "snap", func() [3]bool { return [3]bool{false, true, false} })
		page := "general"
		ui.Box(c).Fill().Children(func() { settingsView(c, &page, &st[1], &st[2], new(int)) })
	})
}

func accountPanics(t *testing.T, view func(*ui.Context)) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("invalid options accepted")
		}
	}()
	accountTester(600, 500, view)
}

func TestSettingsLayoutDisabledNavigationAndControl(t *testing.T) {
	page, sound, track, changes := "general", false, false, 0
	tt := accountTester(900, 500, func(c *ui.Context) {
		ui.Column(c).Fill().Disabled(true).Children(func() { settingsView(c, &page, &sound, &track, &changes) })
	})
	if err := tt.Click("Privacy"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	if page != "general" || sound || track || changes != 0 {
		t.Fatal("disabled settings change state")
	}
}
