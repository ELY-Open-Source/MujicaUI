package mujica

import (
	"image/color"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// A popup left open must stop taking choices once its control is disabled.
func TestChoicePopupsCloseWhenDisabled(t *testing.T) {
	disabled := false
	var multi []int
	var paths, ids []string
	f := newCascaderFake()
	cases := []struct {
		name, open, pick string
		build            func(c *ui.Context)
		same             func() bool
	}{
		{"multi", "Select…", "Lily", func(c *ui.Context) {
			MultiSelect(c, &multi, multiSelectTestOptions, MultiSelectOptions{Disabled: disabled})
		}, func() bool { return len(multi) == 0 }},
		{"cascader", "Select…", "Europe", func(c *ui.Context) {
			Cascader(c, &paths, CascaderOptions{Level: f.level, Load: f.load, AnyLevel: true, Disabled: disabled})
		}, func() bool { return len(paths) == 0 }},
		{"tree", "Select…", "Guard", func(c *ui.Context) {
			TreeSelect(c, &ids, treeSelectTestNodes, TreeSelectOptions{Disabled: disabled})
		}, func() bool { return len(ids) == 0 }},
	}
	for _, tc := range cases {
		disabled = false
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			ui.Column(c).Padding(10).Children(func() { tc.build(c) })
		}, 700, 500)
		if err := tt.Click(tc.open); err != nil {
			t.Fatal(err)
		}
		if !tt.HasText(tc.pick) {
			t.Fatalf("%s: popup did not open", tc.name)
		}
		disabled = true
		tt.Frame()
		_ = tt.Click(tc.pick)
		if !tc.same() || tt.HasText(tc.pick) {
			t.Errorf("%s: disabled control took a choice or kept its popup", tc.name)
		}
	}
}

func TestTreeSelectSingleMarksOwnNode(t *testing.T) {
	single := TreeSelectOptions{}
	court := &treeSelectTestNodes[0]
	if got := treeSelectCheck(court, []string{"court"}, single); got != internal.Checked {
		t.Fatalf("chosen parent shows %v", got)
	}
	if got := treeSelectCheck(court, []string{"king", "queen"}, single); got != internal.Unchecked {
		t.Fatalf("parent of chosen leaves shows %v", got)
	}
	var sel []string
	changes := 0
	tt := treeSelectTester(&sel, single, &changes)
	_ = tt.Click("Select…")
	_ = tt.Click("Court")
	if !slices.Equal(sel, []string{"court"}) {
		t.Fatalf("sel %v", sel)
	}
}

func TestTreeSelectRightSkipsHiddenChildren(t *testing.T) {
	var sel []string
	changes := 0
	tt := treeSelectTester(&sel, TreeSelectOptions{Multiple: true}, &changes)
	_ = tt.Click("Select…")
	tt.Type("queen")
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Court") {
		t.Fatal("Down did not reach Court")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Queen") {
		t.Fatal("Right did not reach the only shown child")
	}
}

// pixel reads the light theme color at (x, y) of the last frame.
func pixel(tt *ui.Tester, x, y float32) color.RGBA {
	return tt.Image().RGBAAt(int(x), int(y))
}

func TestSelectChosenRowHasSelectionFill(t *testing.T) {
	v, changes := "banana", 0
	tt := selectTester(&v, selectTestOptions, SelectOptions{}, &changes)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Move(0, 0)
	tt.Key(0, ui.KeyDown)
	blue, ok := tt.Find("Blueberry")
	if !ok {
		t.Fatal("popup closed")
	}
	h := theme.Default.ControlHeight()
	got := pixel(tt, blue.X+blue.W+20, blue.Y+blue.H/2-h)
	s := theme.Light().Selection
	if got.R != s.R || got.G != s.G || got.B != s.B {
		t.Fatalf("chosen row %v, want selection %v", got, s)
	}
}

func TestChoiceHitAreas(t *testing.T) {
	v, changes := "apple", 0
	tt := selectTester(&v, selectTestOptions, SelectOptions{Clearable: true}, &changes)
	if r, _ := tt.Find("Clear"); r.W < theme.MinHitSize || r.H < theme.MinHitSize {
		t.Errorf("clear button %v", r)
	}
	var sel []string
	tt = treeSelectTester(&sel, TreeSelectOptions{}, &changes)
	_ = tt.Click("Select…")
	if r, _ := tt.Find("Expand Court"); r.W < theme.MinHitSize || r.H < theme.MinHitSize {
		t.Errorf("tree arrow %v", r)
	}
	font := ""
	tt = ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Combobox(c, &font, comboboxTestFonts, ComboboxOptions{})
	}, 400, 300)
	if r, _ := tt.Find("Show options"); r.W < theme.MinHitSize || r.H < theme.MinHitSize {
		t.Errorf("combobox arrow %v", r)
	}
}

func TestChoiceRowsFollowDensity(t *testing.T) {
	for _, d := range []theme.Density{theme.Compact, theme.Default, theme.Comfortable} {
		v, on := "apple", false
		var sw *ui.Element
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Density: d})
			ui.Column(c).Padding(10).AlignItems(ui.Start).Children(func() {
				sw = Switch(c, &on, "Wi-Fi", SwitchOptions{})
				Select(c, &v, selectTestOptions, SelectOptions{Label: "Fruit"})
			})
		}, 400, 500)
		if h := sw.Bounds().H; h < d.ControlHeight() {
			t.Errorf("%v: switch %v high", d, h)
		}
		_ = tt.Click("Fruit")
		a, _ := tt.Find("Cherry")
		b, _ := tt.Find("Date")
		if b.Y-a.Y < d.ControlHeight() {
			t.Errorf("%v: option rows %v apart", d, b.Y-a.Y)
		}
	}
}

func TestMultiSelectCaptionFollowsTextScale(t *testing.T) {
	height := func(scale float32) float32 {
		values := []int{1}
		tt := multiSelectTester(&values, MultiSelectOptions{Max: 2}, new(int))
		tt.SetPreferences(ui.Preferences{TextScale: scale})
		_ = tt.Click("Rose")
		r, ok := tt.Find("1 of 2 chosen")
		if !ok {
			t.Fatal("limit note missing")
		}
		return r.H
	}
	if one, two := height(1), height(2); two < one*1.5 {
		t.Fatalf("caption %v high at 2x, %v at 1x", two, one)
	}
}

func TestComboboxCaretAfterChoice(t *testing.T) {
	r := &comboboxRun{}
	tt := comboboxTester(r, fixedFonts, ComboboxOptions{})
	tt.Key(0, ui.KeyTab)
	tt.Type("di")
	tt.Key(0, ui.KeyEnter)
	if r.value != "Didot" {
		t.Fatalf("value %q", r.value)
	}
	tt.Type("x")
	if r.query != "Didotx" {
		t.Fatalf("typing after a choice gave %q", r.query)
	}
}

// The first Tab enters a toggle group at its first enabled item; later at the last focused one.
func TestToggleGroupTabEntry(t *testing.T) {
	value := []string{"right"}
	tt := choiceBetween(func(c *ui.Context) {
		ToggleGroup(c, &value, toggleGroupTestItems, ToggleGroupOptions{})
	})
	_ = tt.Click("Before")
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Left") {
		t.Fatal("first Tab did not land on the first item")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyTab)
	_ = tt.Click("Before")
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Center") {
		t.Fatal("Tab did not return to the last focused item")
	}
}
