package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// TestNavigationCaptionsFollowTextScale checks small text grows with the desktop text scale.
func TestNavigationCaptionsFollowTextScale(t *testing.T) {
	for _, tc := range []struct {
		text string
		view func(c *ui.Context)
	}{
		{"Name and seal", func(c *ui.Context) {
			cur := 0
			Steps(c, &cur, []Step{{Title: "Account", Description: "Name and seal"}}, StepsOptions{})
		}},
		{"12", func(c *ui.Context) {
			sel := "a"
			Sidebar(c, &sel, []SidebarSection{{Title: "Mail", Items: []SidebarItem{{ID: "a", Label: "Inbox", Icon: icons.Must("inbox"), Count: 12}}}}, SidebarOptions{})
		}},
		{"Verse of the court", func(c *ui.Context) {
			navigationMenuItem(c, NavigationLink{Label: "Poems", Path: "/p", Description: "Verse of the court"}, false, false)
		}},
		{"Esc close", func(c *ui.Context) {
			open := true
			CommandPalette(c, &open, []Command{{ID: "a", Label: "A"}}, CommandPaletteOptions{})
		}},
	} {
		height := func(scale float32) float32 {
			tt := ui.NewTester(func(c *ui.Context) {
				Use(c, Settings{})
				ui.Column(c).Fill().AlignItems(ui.Start).Children(func() { tc.view(c) })
			}, 900, 700)
			tt.SetPreferences(ui.Preferences{TextScale: scale})
			tt.Frame()
			r, ok := tt.Find(tc.text)
			if !ok {
				t.Fatalf("no %q", tc.text)
			}
			return r.H
		}
		if h1, h2 := height(1), height(2); h2 < h1*1.8 {
			t.Errorf("%q: height %v at 2x text, %v at 1x", tc.text, h2, h1)
		}
	}
}
