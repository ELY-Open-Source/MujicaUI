package mujica_test

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
)

func ExampleButton() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Button(c, "Save", mujica.ButtonOptions{Icon: icons.Must("check")}).Clicked() {
			fmt.Println("saved")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleIconButton() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.IconButton(c, icons.Must("settings"), "Settings", mujica.IconButtonOptions{Variant: mujica.Ghost}).Clicked() {
			fmt.Println("open settings")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleButtonGroup() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		items := []mujica.ButtonGroupItem{{Label: "Day"}, {Label: "Week"}, {Label: "Month"}}
		if i := mujica.ButtonGroup(c, items, mujica.ButtonGroupOptions{Label: "Range"}).ClickedItem(); i >= 0 {
			fmt.Println("chose", items[i].Label)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTextInput() {
	var name string
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.TextInput(c, &name, mujica.TextInputOptions{Label: "Name", Placeholder: "Your name", Clearable: true})
		if r.Submitted() {
			fmt.Println("hello", name)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTextArea() {
	var notes string
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.TextArea(c, &notes, mujica.TextAreaOptions{Label: "Notes", MaxHeight: 160, MaxLength: 200})
		if r.Submitted() {
			fmt.Println("send", notes)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleNumberInput() {
	copies := 1.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.NumberInput(c, &copies, mujica.NumberInputOptions{Min: 1, Max: 99, Label: "Copies"}).Changed() {
			fmt.Println("copies", copies)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePasswordInput() {
	var password string
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.PasswordInput(c, &password, mujica.PasswordInputOptions{Label: "Password"}).Submitted() {
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
		mujica.Use(c, mujica.Settings{})
		matches := 0
		if query != "" {
			matches = strings.Count(text, query)
		}
		r := mujica.SearchInput(c, &query, mujica.SearchInputOptions{Find: true, Matches: matches, Current: &current})
		if r.Moved() {
			fmt.Println("show match", current)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCheckbox() {
	var mail, calendar bool
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CheckboxGroup(c, "Notifications", []mujica.CheckboxItem{
			{Label: "Mail", Value: &mail},
			{Label: "Calendar", Value: &calendar},
		}, mujica.CheckboxGroupOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRadioGroup() {
	size := "m"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		items := []mujica.RadioItem[string]{{Value: "s", Label: "Small"}, {Value: "m", Label: "Medium"}, {Value: "l", Label: "Large"}}
		if mujica.RadioGroup(c, &size, items, mujica.RadioGroupOptions{Label: "Size"}).Changed() {
			fmt.Println("size", size)
		}
	}
	ui.Render(view, 400, 300, 1)
}
