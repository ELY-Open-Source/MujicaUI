package mujica

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

func tabsSample() []Tab {
	return []Tab{
		{Label: "General", Panel: func() {}},
		{Label: "Fonts", Disabled: true},
		{Label: "Colors"},
		{Label: "Advanced"},
	}
}

func TestTabsClickAndPanel(t *testing.T) {
	sel, changes := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ts := tabsSample()
		ts[0].Panel = func() { ui.Text(c, "General page") }
		ts[2].Panel = func() { ui.Text(c, "Colors page") }
		if Tabs(c, &sel, ts, TabsOptions{Label: "Settings"}).Changed() {
			changes++
		}
	}, 500, 200)
	if !tt.HasText("General page") || tt.HasText("Colors page") {
		t.Fatalf("panel of tab 0 not shown: %q", tt.Texts())
	}
	if err := tt.Click("Colors"); err != nil {
		t.Fatal(err)
	}
	if sel != 2 || changes != 1 || !tt.HasText("Colors page") {
		t.Fatalf("click Colors: sel %d changes %d texts %q", sel, changes, tt.Texts())
	}
	_ = tt.Click("Fonts")
	if sel != 2 || changes != 1 {
		t.Fatalf("disabled tab chosen: sel %d changes %d", sel, changes)
	}
}

func TestTabsKeyboardSkipsDisabled(t *testing.T) {
	for _, vertical := range []bool{false, true} {
		sel := 0
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			Tabs(c, &sel, tabsSample(), TabsOptions{Vertical: vertical})
		}, 500, 300)
		tt.Key(0, ui.KeyTab)
		if !tt.Focused("General") {
			t.Fatalf("vertical=%v: Tab did not focus the chosen tab", vertical)
		}
		next := ui.KeyRight
		if vertical {
			next = ui.KeyDown
		}
		tt.Key(0, next)
		if sel != 2 || !tt.Focused("Colors") {
			t.Fatalf("vertical=%v: arrow chose %d, want 2 skipping disabled", vertical, sel)
		}
		tt.Key(0, ui.KeyEnd)
		if sel != 3 {
			t.Fatalf("vertical=%v: End chose %d", vertical, sel)
		}
		tt.Key(0, next)
		if sel != 0 {
			t.Fatalf("vertical=%v: wrap chose %d", vertical, sel)
		}
	}
}

func TestTabsOverflowScrolls(t *testing.T) {
	sel := 0
	var tabs []Tab
	for _, l := range []string{"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf", "Hotel"} {
		tabs = append(tabs, Tab{Label: l})
	}
	var w float32
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(240).Children(func() {
			w = Tabs(c, &sel, tabs, TabsOptions{}).Bounds().W
		})
	}, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Frame()
	r, ok := tt.Find("Hotel")
	if sel != 7 || !ok || r.X+r.W > 240 {
		t.Fatalf("last tab not scrolled into view: sel %d rect %+v", sel, r)
	}
	if w > 240 {
		t.Fatalf("tabs widened their parent to %v", w)
	}
}

func TestTabsInvalid(t *testing.T) {
	for name, tc := range map[string]struct {
		sel  int
		tabs []Tab
	}{
		"empty":    {0, nil},
		"range":    {4, tabsSample()},
		"disabled": {1, tabsSample()},
		"label":    {0, []Tab{{Icon: icons.Must("star")}}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			ui.Render(func(c *ui.Context) {
				Use(c, Settings{})
				Tabs(c, &tc.sel, tc.tabs, TabsOptions{})
			}, 300, 100, 1)
		})
	}
}

func tabsView(c *ui.Context) {
	sel, vsel := 0, 2
	ui.Column(c).Padding(16).Gap(16).Children(func() {
		Tabs(c, &sel, []Tab{
			{Label: "Overview", Icon: icons.Must("house"), Panel: func() {
				ui.Text(c, "The chosen tab's panel.")
			}},
			{Label: "Members", Icon: icons.Must("user")},
			{Label: "Archive", Disabled: true},
			{Label: "Settings", Icon: icons.Must("settings")},
		}, TabsOptions{})
		Tabs(c, &vsel, []Tab{
			{Label: "Profile"},
			{Label: "Security", Disabled: true},
			{Label: "Billing", Panel: func() { ui.Text(c, "Vertical layout.") }},
		}, TabsOptions{Vertical: true})
	})
}

func TestTabsSnapshot(t *testing.T) {
	snapshot(t, "tabs", 520, 300, tabsView)
	for _, d := range []theme.Density{theme.Compact, theme.Comfortable} {
		snapshot(t, "tabs-"+strings.ToLower(d.String()), 520, 340, func(c *ui.Context) {
			Use(c, Settings{Mode: tabsMode(c), Density: d})
			tabsView(c)
		})
	}
}

// tabsMode recovers the mode snapshot applied, to add a density.
func tabsMode(c *ui.Context) Mode {
	if Tokens(c).Dark {
		return Dark
	}
	return Light
}
