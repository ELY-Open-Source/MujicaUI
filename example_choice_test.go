package mujica_test

import (
	"log"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
)

func ExampleSwitch() {
	wifi := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Switch(c, &wifi, "Wi-Fi", mujica.SwitchOptions{}).Changed() {
			log.Println("Wi-Fi:", wifi)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToggle() {
	bold := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Toggle(c, &bold, "", mujica.ToggleOptions{Icon: icons.Must("bold"), Label: "Bold"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToggleGroup() {
	align := []string{"left"}
	items := []mujica.ToggleItem[string]{
		{Value: "left", Label: "Left"},
		{Value: "center", Label: "Center"},
		{Value: "right", Label: "Right"},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ToggleGroup(c, &align, items, mujica.ToggleGroupOptions{Label: "Align"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSlider() {
	volume := 40.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Slider(c, &volume, mujica.SliderOptions{Step: 5, ShowValue: true, Label: "Volume"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRangeSlider() {
	low, high := 20.0, 80.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.RangeSlider(c, &low, &high, mujica.RangeSliderOptions{Step: 10, Label: "Price"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSelect() {
	plan := ""
	options := []mujica.SelectOption[string]{
		{Value: "free", Label: "Free"},
		{Value: "pro", Label: "Pro"},
		{Value: "team", Label: "Team", Disabled: true},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Select(c, &plan, options, mujica.SelectOptions{Clearable: true, Label: "Plan"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMultiSelect() {
	tags := []int{1}
	options := []mujica.SelectOption[int]{
		{Value: 1, Label: "Rose"},
		{Value: 2, Label: "Lily"},
		{Value: 3, Label: "Iris"},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MultiSelect(c, &tags, options, mujica.MultiSelectOptions{Max: 2, Label: "Flowers"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCombobox() {
	font := ""
	fonts := []string{"Garamond", "Georgia", "Didot"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.Combobox(c, &font, fonts, mujica.ComboboxOptions{AllowCustom: true, Label: "Font"})
		if r.QueryChanged() {
			log.Println("search:", r.Query())
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCascader() {
	var path []string
	level := func(p []string) mujica.CascaderLevel {
		if len(p) == 0 {
			return mujica.CascaderLevel{Nodes: []mujica.CascaderNode{{Value: "fr", Label: "France"}}}
		}
		return mujica.CascaderLevel{Nodes: []mujica.CascaderNode{{Value: "paris", Label: "Paris", Leaf: true}}}
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Cascader(c, &path, mujica.CascaderOptions{Level: level, Label: "City"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTreeSelect() {
	var roles []string
	nodes := []mujica.TreeSelectNode{
		{ID: "court", Label: "Court", Children: []mujica.TreeSelectNode{
			{ID: "king", Label: "King"},
			{ID: "queen", Label: "Queen"},
		}},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TreeSelect(c, &roles, nodes, mujica.TreeSelectOptions{Multiple: true, Label: "Roles"})
	}
	ui.Render(view, 400, 300, 1)
}
