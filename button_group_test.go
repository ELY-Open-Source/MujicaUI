package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

var buttonGroupTestItems = []ButtonGroupItem{{Label: "Left"}, {Label: "Center", Disabled: true}, {Label: "Right"}}

func TestButtonGroupClickAndArrows(t *testing.T) {
	var got []int
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if i := ButtonGroup(c, buttonGroupTestItems, ButtonGroupOptions{Label: "Align"}).ClickedItem(); i >= 0 {
			got = append(got, i)
		}
	}, 400, 200)
	_ = tt.Click("Center")
	_ = tt.Click("Left")
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Right") {
		t.Fatal("Right arrow did not skip the disabled item to Right")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Left") {
		t.Fatal("arrows do not wrap round")
	}
	if len(got) != 2 || got[0] != 0 || got[1] != 2 {
		t.Fatalf("clicked %v, want [0 2]", got)
	}
}

func TestButtonGroupOneTabStop(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ButtonGroup(c, buttonGroupTestItems, ButtonGroupOptions{Vertical: true})
		Button(c, "After", ButtonOptions{})
	}, 400, 300)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Left") {
		t.Fatal("Tab did not enter the group")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Right") {
		t.Fatal("Down did not move in a vertical group")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("Tab did not leave the group")
	}
}

func TestButtonGroupDisabledGroup(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if ButtonGroup(c, buttonGroupTestItems, ButtonGroupOptions{Disabled: true}).ClickedItem() >= 0 {
			n++
		}
	}, 400, 200)
	_ = tt.Click("Left")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if n != 0 {
		t.Fatalf("disabled group clicked %d", n)
	}
}

func TestButtonGroupSnapshot(t *testing.T) {
	snapshot(t, "button_group", 420, 200, func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(24).AlignItems(ui.Start).Children(func() {
			ButtonGroup(c, []ButtonGroupItem{{Label: "Day"}, {Label: "Week"}, {Label: "Month", Disabled: true}}, ButtonGroupOptions{})
			ButtonGroup(c, []ButtonGroupItem{
				{Icon: icons.Must("bold"), Name: "Bold"},
				{Icon: icons.Must("italic"), Name: "Italic"},
				{Icon: icons.Must("underline"), Name: "Underline"},
			}, ButtonGroupOptions{Vertical: true})
		})
	})
}
