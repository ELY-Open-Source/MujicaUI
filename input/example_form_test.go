package input_test

import (
	"fmt"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func ExampleForm() {
	var state input.FormState
	name := ""
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res := input.Form(c, &state, input.FormOptions{}, func(f *input.FormScope) {
			input.FormItem(c, f, "name", &name, input.FormItemOptions[string]{Label: "Name", Required: true}, func() ui.Element {
				return input.InputGroup(c, &name, input.InputGroupOptions{Invalid: state.Error("name") != ""}).Input
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
		core.Use(c, core.Settings{})
		input.FormField(c, "Email", input.FormFieldOptions{Required: true, Description: "For receipts."}, func() ui.Element {
			return input.InputGroup(c, &email, input.InputGroupOptions{}).Input
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleInputGroup() {
	site := ""
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		g := input.InputGroup(c, &site, input.InputGroupOptions{
			Prefix: "https://", Label: "Site",
			Buttons: []input.InputGroupButton{{Label: "Open"}},
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
		core.Use(c, core.Settings{})
		if input.PinInput(c, &code, input.PinInputOptions{Length: 6}).Completed() {
			fmt.Println("verify", code)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTagsInput() {
	tags := []string{"gothic"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.TagsInput(c, &tags, input.TagsInputOptions{Max: 5, Placeholder: "Add a tag", Label: "Tags"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRating() {
	stars := 3.5
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.Rating(c, &stars, input.RatingOptions{Step: 0.5, ShowValue: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleColorPicker() {
	tint := ui.Hex("#7D2034")
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.ColorPicker(c, &tint, input.ColorPickerOptions{Label: "Tint"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFilePicker() {
	var paths []string
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.FilePicker(c, &paths, input.FilePickerOptions{
			Multiple: true,
			Filters:  []mygo.FileFilter{{Name: "Scores", Extensions: []string{"pdf", "mid"}}},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFileDropZone() {
	var files []input.DroppedFile
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.FileDropZone(c, &files, input.FileDropZoneOptions{
			Extensions: []string{"png", "jpg"}, MaxFiles: 3, MaxSize: 5 << 20,
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTransfer() {
	chosen := []string{"uika"}
	items := []input.TransferItem{{ID: "uika", Label: "Uika"}, {ID: "mutsumi", Label: "Mutsumi"}, {ID: "nyamu", Label: "Nyamu", Disabled: true}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.Transfer(c, &chosen, items, input.TransferOptions{})
	}
	ui.Render(view, 400, 300, 1)
}
