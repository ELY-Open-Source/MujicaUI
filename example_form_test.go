package mujica_test

import (
	"fmt"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleForm() {
	var state mujica.FormState
	name := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		res := mujica.Form(c, &state, mujica.FormOptions{}, func(f *mujica.FormScope) {
			mujica.FormItem(f, "name", &name, mujica.FormItemOptions[string]{Label: "Name", Required: true}, func() *ui.Element {
				return mujica.InputGroup(c, &name, mujica.InputGroupOptions{Invalid: state.Error("name") != ""}).Input
			})
		})
		if res.Submitted() {
			state.EndSubmit("") // after saving name
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFormField() {
	email := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.FormField(c, "Email", mujica.FormFieldOptions{Required: true, Description: "For receipts."}, func() {
			mujica.InputGroup(c, &email, mujica.InputGroupOptions{})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleInputGroup() {
	site := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		g := mujica.InputGroup(c, &site, mujica.InputGroupOptions{
			Prefix: "https://", Label: "Site",
			Buttons: []mujica.InputGroupButton{{Label: "Open"}},
		})
		if g.Clicked == 0 {
			c.OpenURL("https://" + site)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePinInput() {
	code := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.PinInput(c, &code, mujica.PinInputOptions{Length: 6}).Completed() {
			fmt.Println("verify", code)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTagsInput() {
	tags := []string{"gothic"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TagsInput(c, &tags, mujica.TagsInputOptions{Max: 5, Placeholder: "Add a tag", Label: "Tags"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRating() {
	stars := 3.5
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Rating(c, &stars, mujica.RatingOptions{Step: 0.5, ShowValue: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleColorPicker() {
	tint := ui.Hex("#7D2034")
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ColorPicker(c, &tint, mujica.ColorPickerOptions{Label: "Tint"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFilePicker() {
	var paths []string
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.FilePicker(c, &paths, mujica.FilePickerOptions{
			Multiple: true,
			Filters:  []mygo.FileFilter{{Name: "Scores", Extensions: []string{"pdf", "mid"}}},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFileDropZone() {
	var files []mujica.DroppedFile
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.FileDropZone(c, &files, mujica.FileDropZoneOptions{
			Extensions: []string{"png", "jpg"}, MaxFiles: 3, MaxSize: 5 << 20,
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTransfer() {
	chosen := []string{"uika"}
	items := []mujica.TransferItem{{ID: "uika", Label: "Uika"}, {ID: "mutsumi", Label: "Mutsumi"}, {ID: "nyamu", Label: "Nyamu", Disabled: true}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Transfer(c, &chosen, items, mujica.TransferOptions{})
	}
	ui.Render(view, 400, 300, 1)
}
