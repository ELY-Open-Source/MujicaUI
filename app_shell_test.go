package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

type appShellProbe struct {
	collapsed bool
	toggled   int
	content   ui.ScrollState
	side      ui.Rect
}

func (p *appShellProbe) view(c *ui.Context) {
	Use(c, Settings{})
	r := AppShell(c, &p.collapsed, AppShellOptions{
		TitleBar: func() { TitleBar(c, "Court Records", TitleBarOptions{}) },
		Sidebar: func() {
			for i := range 30 {
				ui.Text(c, fmt.Sprintf("Nav %d", i)).Height(24)
			}
		},
		Content: func() {
			for i := range 40 {
				ui.Text(c, fmt.Sprintf("Entry %d", i)).Height(24)
			}
		},
		StatusBar:     func() { StatusBar(c, StatusBarOptions{Left: []StatusItem{{ID: "s", Text: "Ready"}}}) },
		ContentScroll: &p.content,
	})
	if r.Toggled {
		p.toggled++
	}
}

func appShellTester(p *appShellProbe) *ui.Tester {
	tt := ui.NewTester(p.view, 640, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	return tt
}

func TestAppShellSlots(t *testing.T) {
	tt := appShellTester(&appShellProbe{})
	title, _ := tt.Find("Court Records")
	nav, _ := tt.Find("Nav 0")
	entry, _ := tt.Find("Entry 0")
	ready, _ := tt.Find("Ready")
	if title.Y > nav.Y || nav.X >= entry.X || ready.Y < 360 {
		t.Fatalf("slots out of place: title %+v nav %+v entry %+v status %+v", title, nav, entry, ready)
	}
	if entry.X < 240 {
		t.Fatalf("content at %v overlaps the 240 DIP sidebar", entry.X)
	}
}

func TestAppShellCollapseByButton(t *testing.T) {
	p := &appShellProbe{}
	tt := appShellTester(p)
	if err := tt.Click("Collapse sidebar"); err != nil {
		t.Fatal(err)
	}
	if !p.collapsed || p.toggled != 1 || tt.HasText("Nav 0") {
		t.Fatalf("after collapse: collapsed=%v toggled=%d nav shown=%v", p.collapsed, p.toggled, tt.HasText("Nav 0"))
	}
	if entry, _ := tt.Find("Entry 0"); entry.X != 44 {
		t.Fatalf("content at %v beside the rail, want 44", entry.X)
	}
	if err := tt.Click("Expand sidebar"); err != nil {
		t.Fatal(err)
	}
	if p.collapsed || !tt.HasText("Nav 0") || p.toggled != 2 {
		t.Fatal("Expand did not reopen the sidebar")
	}
}

func TestAppShellCollapseByShortcut(t *testing.T) {
	p := &appShellProbe{}
	tt := appShellTester(p)
	tt.Key(ui.Cmd, ui.KeyB)
	if !p.collapsed {
		t.Fatal("Cmd+B did not collapse")
	}
	tt.Key(ui.Cmd, ui.KeyB)
	if p.collapsed || p.toggled != 2 {
		t.Fatalf("Cmd+B twice: collapsed=%v toggled=%d", p.collapsed, p.toggled)
	}
}

func TestAppShellContentScrollsAlone(t *testing.T) {
	p := &appShellProbe{}
	tt := appShellTester(p)
	nav, _ := tt.Find("Nav 0")
	tt.Scroll(400, 200, 0, 96)
	if p.content.Y != 96 {
		t.Fatalf("content offset %v, want 96", p.content.Y)
	}
	if r, _ := tt.Find("Nav 0"); r.Y != nav.Y {
		t.Fatalf("sidebar moved from %v to %v with the content", nav.Y, r.Y)
	}
	if r, _ := tt.Find("Court Records"); r.Y > 20 {
		t.Fatalf("title bar scrolled to %v", r.Y)
	}
	// Collapsing keeps the content's place.
	tt.Key(ui.Cmd, ui.KeyB)
	if p.content.Y != 96 {
		t.Fatalf("collapse reset the content offset to %v", p.content.Y)
	}
}

func TestAppShellNeedsContent(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("missing Content did not panic")
		}
	}()
	var collapsed bool
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		AppShell(c, &collapsed, AppShellOptions{})
	}, 100, 100)
}

func TestAppShellSnapshot(t *testing.T) {
	snapshot(t, "app_shell", 560, 260, func(c *ui.Context) {
		collapsed := ui.Local(c.Root(), "snap", func() bool { return false })
		AppShell(c, collapsed, AppShellOptions{
			TitleBar: func() { TitleBar(c, "Court Records", TitleBarOptions{}) },
			Sidebar: func() {
				for _, s := range []string{"Chronicles", "Letters", "Decrees"} {
					ui.Text(c, s).Padding(6, 12)
				}
			},
			Content: func() {
				Container(c, ContainerOptions{}, func() {
					Stack(c, StackOptions{Gap: SpaceM}, func() {
						ui.Text(c, "The Grand Hall").FontSize(18)
						Divider(c, DividerOptions{Ornament: true})
						ui.Text(c, "Content scrolls on its own, apart from the sidebar.")
					})
				})
			},
			StatusBar: func() { StatusBar(c, StatusBarOptions{Left: []StatusItem{{ID: "s", Text: "Ready"}}}) },
		})
	})
}
