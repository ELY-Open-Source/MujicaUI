package navigation

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
)

func breadcrumbPath(names ...string) []BreadcrumbItem {
	var items []BreadcrumbItem
	for _, n := range names {
		items = append(items, BreadcrumbItem{Label: n})
	}
	return items
}

func TestBreadcrumbChoose(t *testing.T) {
	got := -1
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if i, ok := Breadcrumb(c, breadcrumbPath("Home", "Library", "Poems"), BreadcrumbOptions{}).Chosen(); ok {
			got = i
		}
	}, 500, 100)
	if err := tt.Click("Library"); err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Fatalf("chose %d, want 1", got)
	}
	got = -1
	_ = tt.Click("Poems")
	if got != -1 {
		t.Fatalf("current item acted as a link: %d", got)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if got != 0 {
		t.Fatalf("Enter on first link chose %d", got)
	}
}

func TestBreadcrumbFolds(t *testing.T) {
	got := -1
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if i, ok := Breadcrumb(c, breadcrumbPath("Root", "Alpha", "Beta", "Gamma", "Omega", "Leaf"), BreadcrumbOptions{}).Chosen(); ok {
			got = i
		}
	}, 600, 300)
	for _, hidden := range []string{"Alpha", "Beta", "Gamma"} {
		if tt.HasText(hidden) {
			t.Fatalf("%s should be folded: %q", hidden, tt.Texts())
		}
	}
	if !tt.HasText("Omega") || !tt.HasText("Root") {
		t.Fatalf("ends missing: %q", tt.Texts())
	}
	if err := tt.Click("Show hidden path"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Beta") || !tt.Focused("Alpha") {
		t.Fatalf("menu did not open with focus on first: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Beta") {
		t.Fatal("Escape did not close the menu")
	}
	_ = tt.Click("Show hidden path")
	if err := tt.Click("Gamma"); err != nil {
		t.Fatal(err)
	}
	if got != 3 || tt.HasText("Beta") {
		t.Fatalf("chose %d, want 3, menu closed", got)
	}
}

func TestBreadcrumbInvalid(t *testing.T) {
	for name, tc := range map[string]struct {
		items []BreadcrumbItem
		opts  BreadcrumbOptions
	}{
		"empty": {nil, BreadcrumbOptions{}},
		"max":   {breadcrumbPath("a", "b"), BreadcrumbOptions{MaxItems: 2}},
		"label": {[]BreadcrumbItem{{Icon: icons.Must("house")}}, BreadcrumbOptions{}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				Breadcrumb(c, tc.items, tc.opts)
			}, 300, 100, 1)
		})
	}
}

func TestBreadcrumbSnapshot(t *testing.T) {
	snapshot(t, "breadcrumb", 560, 120, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
			Breadcrumb(c, []BreadcrumbItem{{Label: "Home", Icon: icons.Must("house")}, {Label: "Library"}, {Label: "Collected poems"}}, BreadcrumbOptions{})
			Breadcrumb(c, breadcrumbPath("Archive", "2026", "October", "Correspondence of the court", "Letter to the Duchess of Mujica"), BreadcrumbOptions{MaxItemWidth: 140})
		})
	})
}
