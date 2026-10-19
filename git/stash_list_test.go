package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
)

func stashSample() []Stash {
	return []Stash{
		{Ref: "stash@{0}", Message: "WIP on main: graph lanes", Branch: "main"},
		{Ref: "stash@{1}", Message: "Experiment with blame colors", Branch: "fix/blame"},
	}
}

func TestStashListActions(t *testing.T) {
	var s data.ListState[string]
	var acted [][2]string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := StashList(c, &s, stashSample(), StashListOptions{})
		r.Element.Height(200)
		if a, item, ok := r.Action(); ok {
			acted = append(acted, [2]string{a, item})
		}
	}, 700, 240)
	r, _ := tt.Find("Experiment with blame colors")
	y := r.Y + r.H/2
	// Apply, Pop and Drop sit at the right end, 28 DIPs each.
	for _, x := range []float32{700 - 90, 700 - 60, 700 - 30} {
		tt.ClickAt(x, y)
	}
	want := [][2]string{{"apply", "stash@{1}"}, {"pop", "stash@{1}"}, {"drop", "stash@{1}"}}
	if len(acted) != 3 || acted[0] != want[0] || acted[1] != want[1] || acted[2] != want[2] {
		t.Fatalf("actions %q", acted)
	}
	for _, name := range []string{"Apply", "Pop", "Drop"} {
		if _, ok := tt.Find(name); !ok {
			t.Errorf("no %s button", name)
		}
	}
}

func TestStashListEmptyAndInvalid(t *testing.T) {
	var s data.ListState[string]
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		StashList(c, &s, nil, StashListOptions{})
	}, 400, 200)
	if !tt.HasText("没有贮藏") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		var s data.ListState[string]
		StashList(c, &s, []Stash{{Message: "x"}}, StashListOptions{}).Element.Height(100)
	}, 300, 120, 1)
}

func TestStashListSnapshot(t *testing.T) {
	snapshot(t, "stash_list", 640, 140, func(c *ui.Context) {
		s := ui.Local(c.Root(), "list", func() *data.ListState[string] { return &data.ListState[string]{} })
		ui.Column(c).Padding(16).Children(func() {
			StashList(c, *s, stashSample(), StashListOptions{}).Element.Height(100)
		})
	})
}
