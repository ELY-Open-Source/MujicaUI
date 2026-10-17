package navigation

import (
	"image/color"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
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
		core.Use(c, core.Settings{})
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
			core.Use(c, core.Settings{})
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
		core.Use(c, core.Settings{})
		ui.Column(c).Width(240).Children(func() {
			w = Tabs(c, &sel, tabs, TabsOptions{}).Element.Bounds().W
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
				core.Use(c, core.Settings{})
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
			core.Use(c, core.Settings{Mode: tabsMode(c), Density: d})
			tabsView(c)
		})
	}
}

func TestTabsStateSnapshots(t *testing.T) {
	center := func(tt *ui.Tester, s string) (float32, float32) {
		r, ok := tt.Find(s)
		if !ok {
			t.Fatalf("no %q", s)
		}
		return r.X + r.W/2, r.Y + r.H/2
	}
	snapshotAfter(t, "tabs-hover", 520, 300, tabsView, func(tt *ui.Tester) { tt.Move(center(tt, "Members")) })
	snapshotAfter(t, "tabs-pressed", 520, 300, tabsView, func(tt *ui.Tester) { tt.Press(center(tt, "Settings")) })
	snapshotAfter(t, "tabs-focus", 520, 300, tabsView, func(tt *ui.Tester) {
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyRight)
	})
}

// TestTabsStateColorsEase checks hover fades in over StateDuration, not at once.
func TestTabsStateColorsEase(t *testing.T) {
	sel := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		Tabs(c, &sel, []Tab{{Label: "One"}, {Label: "Two"}}, TabsOptions{})
	}, 300, 120)
	r, _ := tt.Find("Two")
	x, y := int(r.X-8), int(r.Y+r.H/2)
	hover := tabsRGBA(theme.Light().SurfaceHover)
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	if got := tt.Image().RGBAAt(x, y); got == hover {
		t.Fatal("hover color applied at once")
	}
	tt.Advance(theme.StateDuration + 40*time.Millisecond)
	if got := tt.Image().RGBAAt(x, y); got != hover {
		t.Fatalf("hover color %v after the transition, want %v", got, hover)
	}
}

func tabsRGBA(c ui.Color) color.RGBA { return color.RGBA{c.R, c.G, c.B, c.A} }

// tabsMode recovers the mode snapshot applied, to add a density.
func tabsMode(c *ui.Context) core.Mode {
	if core.Tokens(c).Dark {
		return core.Dark
	}
	return core.Light
}

func TestTabsAncestorDisabledNavigation(t *testing.T) {
	for _, key := range []ui.Key{ui.KeyRight, ui.KeyDown, ui.KeyEnd} {
		selected, changes := 0, 0
		disabled := false
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			ui.Column(c).Disabled(disabled).Children(func() {
				if Tabs(c, &selected, tabsSample(), TabsOptions{}).Changed() {
					changes++
				}
			})
		}, 500, 200)
		if err := tt.Click("General"); err != nil {
			t.Fatal(err)
		}
		disabled = true
		tt.Key(0, key)
		if selected != 0 || changes != 0 {
			t.Fatalf("disabled tab navigation: selected=%d events=%d", selected, changes)
		}
	}
}
