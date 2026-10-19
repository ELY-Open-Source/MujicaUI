package input_test

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func ExampleButton() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Save", input.ButtonOptions{Icon: icons.Must("check")}).Clicked() {
			fmt.Println("saved")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleIconButton() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.IconButton(c, icons.Must("settings"), "Settings", input.IconButtonOptions{Variant: input.Ghost}).Clicked() {
			fmt.Println("open settings")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleButtonGroup() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		items := []input.ButtonGroupItem{{Label: "Day"}, {Label: "Week"}, {Label: "Month"}}
		if i := input.ButtonGroup(c, items, input.ButtonGroupOptions{Label: "Range"}).ClickedItem(); i >= 0 {
			fmt.Println("chose", items[i].Label)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTextInput() {
	var name string
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := input.TextInput(c, &name, input.TextInputOptions{Label: "Name", Placeholder: "Your name", Clearable: true})
		if r.Submitted() {
			fmt.Println("hello", name)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTextArea() {
	var notes string
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := input.TextArea(c, &notes, input.TextAreaOptions{Label: "Notes", MaxHeight: 160, MaxLength: 200})
		if r.Submitted() {
			fmt.Println("send", notes)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleNumberInput() {
	copies := 1.0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.NumberInput(c, &copies, input.NumberInputOptions{Min: 1, Max: 99, Label: "Copies"}).Changed() {
			fmt.Println("copies", copies)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePasswordInput() {
	var password string
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.PasswordInput(c, &password, input.PasswordInputOptions{Label: "Password"}).Submitted() {
			fmt.Println("signing in")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSearchInput() {
	var query string
	current := 0
	text := "rose lys rose crown rose"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		matches := 0
		if query != "" {
			matches = strings.Count(text, query)
		}
		r := input.SearchInput(c, &query, input.SearchInputOptions{Find: true, Matches: matches, Current: &current})
		if r.Moved() {
			fmt.Println("show match", current)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCheckbox() {
	var mail, calendar bool
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.CheckboxGroup(c, "Notifications", []input.CheckboxItem{
			{Label: "Mail", Value: &mail},
			{Label: "Calendar", Value: &calendar},
		}, input.CheckboxGroupOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRadioGroup() {
	size := "m"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		items := []input.RadioItem[string]{{Value: "s", Label: "Small"}, {Value: "m", Label: "Medium"}, {Value: "l", Label: "Large"}}
		if input.RadioGroup(c, &size, items, input.RadioGroupOptions{Label: "Size"}).Changed() {
			fmt.Println("size", size)
		}
	}
	ui.Render(view, 400, 300, 1)
}
