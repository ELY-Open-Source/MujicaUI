package navigation

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
)

func TestSegmentedControlClickAndKeys(t *testing.T) {
	sel, changes := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if SegmentedControl(c, &sel, []Segment{{Label: "Day"}, {Label: "Week"}, {Label: "Month"}}, SegmentedControlOptions{Label: "Range"}).Changed() {
			changes++
		}
	}, 400, 100)
	if err := tt.Click("Month"); err != nil {
		t.Fatal(err)
	}
	if sel != 2 || changes != 1 {
		t.Fatalf("click: sel %d changes %d", sel, changes)
	}
	tt.Key(0, ui.KeyRight)
	if sel != 0 || changes != 2 {
		t.Fatalf("Right wraps to %d, changes %d", sel, changes)
	}
	tt.Key(0, ui.KeyLeft)
	if sel != 2 {
		t.Fatalf("Left wraps to %d", sel)
	}
	tt.Key(0, ui.KeyHome)
	if sel != 0 || !tt.Focused("Day") {
		t.Fatalf("Home chose %d", sel)
	}
}

func TestSegmentedControlIconsAndDisabled(t *testing.T) {
	sel := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SegmentedControl(c, &sel, []Segment{
			{Label: "List", Icon: icons.Must("list"), HideLabel: true},
			{Label: "Grid", Icon: icons.Must("layout-grid"), HideLabel: true},
		}, SegmentedControlOptions{Disabled: true})
	}, 400, 100)
	if _, ok := tt.Find("Grid"); !ok {
		t.Fatal("icon segment has no name")
	}
	_ = tt.Click("Grid")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if sel != 0 {
		t.Fatalf("disabled control changed to %d", sel)
	}
}

func TestSegmentedControlEqualWidths(t *testing.T) {
	sel := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			SegmentedControl(c, &sel, []Segment{{Label: "A"}, {Label: "Much longer"}}, SegmentedControlOptions{Equal: true})
		})
	}, 400, 100)
	a, _ := tt.Find("A")
	b, _ := tt.Find("Much longer")
	// Equal segments center their text: A's center sits a segment width before the other's.
	da := (b.X + b.W/2) - (a.X + a.W/2)
	if da < b.W {
		t.Fatalf("segments not equal: A %+v, long %+v", a, b)
	}
}

func TestSegmentedControlInvalid(t *testing.T) {
	for name, segs := range map[string][]Segment{
		"empty":  nil,
		"label":  {{Icon: icons.Must("list")}},
		"noicon": {{Label: "X", HideLabel: true}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			sel := 0
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				SegmentedControl(c, &sel, segs, SegmentedControlOptions{})
			}, 300, 100, 1)
		})
	}
}

func TestSegmentedControlSnapshot(t *testing.T) {
	snapshot(t, "segmented_control", 480, 180, func(c *ui.Context) {
		a, b, e := 1, 0, 2
		ui.Column(c).Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
			SegmentedControl(c, &a, []Segment{{Label: "Day"}, {Label: "Week"}, {Label: "Month"}}, SegmentedControlOptions{})
			ui.Row(c).Gap(16).Children(func() {
				SegmentedControl(c, &b, []Segment{
					{Label: "List", Icon: icons.Must("list"), HideLabel: true},
					{Label: "Grid", Icon: icons.Must("layout-grid"), HideLabel: true},
				}, SegmentedControlOptions{})
				SegmentedControl(c, &b, []Segment{{Label: "Light", Icon: icons.Must("sun")}, {Label: "Dark", Icon: icons.Must("moon")}}, SegmentedControlOptions{Disabled: true})
			})
			SegmentedControl(c, &e, []Segment{{Label: "S"}, {Label: "Medium"}, {Label: "XL"}}, SegmentedControlOptions{Equal: true})
		})
	})
}
