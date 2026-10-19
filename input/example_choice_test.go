package input_test

import (
	"log"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func ExampleSwitch() {
	wifi := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Switch(c, &wifi, "Wi-Fi", input.SwitchOptions{}).Changed() {
			log.Println("Wi-Fi:", wifi)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToggle() {
	bold := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.Toggle(c, &bold, "", input.ToggleOptions{Icon: icons.Must("bold"), Label: "Bold"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToggleGroup() {
	align := []string{"left"}
	items := []input.ToggleItem[string]{
		{Value: "left", Label: "Left"},
		{Value: "center", Label: "Center"},
		{Value: "right", Label: "Right"},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.ToggleGroup(c, &align, items, input.ToggleGroupOptions{Label: "Align"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSlider() {
	volume := 40.0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.Slider(c, &volume, input.SliderOptions{Step: 5, ShowValue: true, Label: "Volume"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRangeSlider() {
	low, high := 20.0, 80.0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.RangeSlider(c, &low, &high, input.RangeSliderOptions{Step: 10, Label: "Price"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSelect() {
	plan := ""
	options := []input.SelectOption[string]{
		{Value: "free", Label: "Free"},
		{Value: "pro", Label: "Pro"},
		{Value: "team", Label: "Team", Disabled: true},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.Select(c, &plan, options, input.SelectOptions{Clearable: true, Label: "Plan"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMultiSelect() {
	tags := []int{1}
	options := []input.SelectOption[int]{
		{Value: 1, Label: "Rose"},
		{Value: 2, Label: "Lily"},
		{Value: 3, Label: "Iris"},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.MultiSelect(c, &tags, options, input.MultiSelectOptions{Max: 2, Label: "Flowers"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCombobox() {
	font := ""
	fonts := []string{"Garamond", "Georgia", "Didot"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := input.Combobox(c, &font, fonts, input.ComboboxOptions{AllowCustom: true, Label: "Font"})
		if r.QueryChanged() {
			log.Println("search:", r.Query())
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCascader() {
	var path []string
	level := func(p []string) input.CascaderLevel {
		if len(p) == 0 {
			return input.CascaderLevel{Nodes: []input.CascaderNode{{Value: "fr", Label: "France"}}}
		}
		return input.CascaderLevel{Nodes: []input.CascaderNode{{Value: "paris", Label: "Paris", Leaf: true}}}
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.Cascader(c, &path, input.CascaderOptions{Level: level, Label: "City"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTreeSelect() {
	var roles []string
	nodes := []input.TreeSelectNode{
		{ID: "court", Label: "Court", Children: []input.TreeSelectNode{
			{ID: "king", Label: "King"},
			{ID: "queen", Label: "Queen"},
		}},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.TreeSelect(c, &roles, nodes, input.TreeSelectOptions{Multiple: true, Label: "Roles"})
	}
	ui.Render(view, 400, 300, 1)
}
