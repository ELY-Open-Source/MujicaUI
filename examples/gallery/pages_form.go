package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
)

// formStates shows a component's states side by side.
func formStates(c *ui.Context, fn func()) {
	ui.Column(c).Gap(8).Children(func() {
		ui.Text(c, "States").FontWeight(600)
		ui.Row(c).Gap(24).Wrap().AlignItems(ui.Start).Children(fn)
	})
}

// formToggles shows parameter check boxes in one row.
func formToggles(c *ui.Context, names []string, values ...*bool) {
	ui.Row(c).Gap(16).Wrap().Children(func() {
		for i, n := range names {
			ui.Checkbox(c, values[i], n)
		}
	})
}

type formDemo struct {
	st                     mujica.FormState
	name, email            string
	tags                   []string
	horizontal, slowServer bool
}

type fieldDemo struct {
	value                string
	required, horizontal bool
	withError            bool
}

type groupDemo struct {
	value                       string
	disabled, readOnly, invalid bool
}

type pinDemo struct {
	code              string
	length4, disabled bool
	invalid           bool
}

type tagsDemo struct {
	tags            []string
	limit, disabled bool
}

type ratingDemo struct {
	value                    float64
	half, readOnly, disabled bool
}

type colorDemo struct {
	color    ui.Color
	disabled bool
}

type pickerDemo struct {
	paths                       []string
	directory, multiple, filter bool
	disabled                    bool
}

type dropDemo struct {
	files    []mujica.DroppedFile
	images   bool
	disabled bool
}

type transferDemo struct {
	chosen   []string
	disabled bool
}

var transferDemoItems = []mujica.TransferItem{
	{ID: "sakiko", Label: "Sakiko Togawa"}, {ID: "uika", Label: "Uika Misumi"}, {ID: "mutsumi", Label: "Mutsumi Wakaba"},
	{ID: "umiri", Label: "Umiri Yahata"}, {ID: "nyamu", Label: "Nyamu Yutenji"}, {ID: "tomori", Label: "Tomori Takamatsu", Disabled: true},
	{ID: "anon", Label: "Anon Chihaya"}, {ID: "rana", Label: "Rana Kaname"}, {ID: "soyo", Label: "Soyo Nagasaki"},
}

type formAsync struct {
	version uint64
	err     string
}

var formPages = []page{
	{ID: "041", Name: "Form", ZhName: "表单", Category: "Forms", Code: formCode["Form"],
		Keys: "Tab moves between fields; leaving a field validates it. Enter in a field submits. A failed submit focuses the first error.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "041", func() formDemo {
				return formDemo{name: "Sakiko", tags: []string{"keyboard"}}
			})
			async := *ui.Local(c.Root(), "041.async", func() chan formAsync { return make(chan formAsync, 8) })
			saved := *ui.Local(c.Root(), "041.save", func() chan string { return make(chan string, 1) })
			for drained := false; !drained; {
				select {
				case r := <-async:
					log(fmt.Sprintf("async name check v%d: %q, kept %v", r.version, r.err, d.st.Resolve("name", r.version, r.err)))
				case e := <-saved:
					d.st.EndSubmit(e)
					log(fmt.Sprintf("submit ended, error %q", e))
				default:
					drained = true
				}
			}
			formToggles(c, []string{"Labels beside fields", "Server rejects the submit"}, &d.horizontal, &d.slowServer)
			res := mujica.Form(c, &d.st, mujica.FormOptions{Horizontal: d.horizontal}, func(f *mujica.FormScope) {
				mujica.FormItem(f, "name", &d.name, mujica.FormItemOptions[string]{Label: "Stage name", Required: true, Initial: "Sakiko",
					Description: "Checked against the archive; \"Oblivionis\" is taken.",
					Async: func(v string, version uint64) {
						log(fmt.Sprintf("async check started v%d", version))
						go func() {
							time.Sleep(700 * time.Millisecond)
							msg := ""
							if strings.EqualFold(v, "Oblivionis") {
								msg = "That name is taken."
							}
							async <- formAsync{version, msg}
							c.Invalidate()
						}()
					}}, func() *ui.Element {
					return mujica.InputGroup(c, &d.name, mujica.InputGroupOptions{Invalid: d.st.Error("name") != ""}).Input
				})
				mujica.FormItem(f, "email", &d.email, mujica.FormItemOptions[string]{Label: "Email", Required: true,
					Validate: func(v string) string {
						if !strings.Contains(v, "@") {
							return "Enter an email address."
						}
						return ""
					}}, func() *ui.Element {
					g := mujica.InputGroup(c, &d.email, mujica.InputGroupOptions{Invalid: d.st.Error("email") != ""})
					if g.Submitted() {
						d.st.Submit()
					}
					return g.Input
				})
				mujica.FormItem(f, "tags", &d.tags, mujica.FormItemOptions[[]string]{Label: "Roles", Required: true, Initial: []string{"keyboard"}}, func() *ui.Element {
					return mujica.TagsInput(c, &d.tags, mujica.TagsInputOptions{Max: 4, Invalid: d.st.Error("tags") != ""}).Input
				})
			})
			if res.Submitted() {
				log(fmt.Sprintf("submitted %q %q %q", d.name, d.email, d.tags))
				fail := ""
				if d.slowServer {
					fail = "The archive refused the entry. Your values are kept."
				}
				go func() {
					time.Sleep(900 * time.Millisecond)
					saved <- fail
					c.Invalidate()
				}()
			}
			if res.Reset() {
				log("reset to initial values")
			}
		}},
	{ID: "042", Name: "FormField", ZhName: "表单字段", Category: "Forms", Code: formCode["FormField"],
		Keys: "Clicking the label focuses the control; Tab reaches the control directly.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "042", func() fieldDemo { return fieldDemo{} })
			formToggles(c, []string{"Required", "Label beside", "Show error"}, &d.required, &d.horizontal, &d.withError)
			e := ""
			if d.withError {
				e = "This name is already in the guest list."
			}
			mujica.FormField(c, "Guest name", mujica.FormFieldOptions{Required: d.required, Horizontal: d.horizontal,
				Description: "As written on the invitation.", Error: e}, func() *ui.Element {
				g := mujica.InputGroup(c, &d.value, mujica.InputGroupOptions{Invalid: d.withError})
				if g.Changed() {
					log(fmt.Sprintf("value %q", d.value))
				}
				return g.Input
			})
		}},
	{ID: "043", Name: "InputGroup", ZhName: "输入组合", Category: "Forms", Code: formCode["InputGroup"],
		Keys: "Each button and the input are separate Tab stops. Enter in the input submits.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "043", func() groupDemo { return groupDemo{value: "mujica"} })
			formToggles(c, []string{"Disabled", "Read-only", "Invalid"}, &d.disabled, &d.readOnly, &d.invalid)
			opts := mujica.InputGroupOptions{Prefix: "https://", Suffix: ".org", Label: "Site", Disabled: d.disabled, ReadOnly: d.readOnly, Invalid: d.invalid,
				Buttons: []mujica.InputGroupButton{{Label: "Clear", Icon: icons.Must("x"), IconOnly: true}, {Label: "Copy", Icon: icons.Must("copy")}}}
			g := mujica.InputGroup(c, &d.value, opts)
			switch g.Clicked {
			case 0:
				d.value = ""
				log("clear")
			case 1:
				c.WriteClipboard("https://" + d.value + ".org")
				log("copied")
			}
			if g.Changed() {
				log(fmt.Sprintf("value %q", d.value))
			}
			if g.Submitted() {
				log("submitted")
			}
			formStates(c, func() {
				for _, o := range []mujica.InputGroupOptions{{Label: "Disabled", Disabled: true}, {Label: "Read-only", ReadOnly: true}, {Label: "Invalid", Invalid: true}} {
					v := o.Label
					mujica.InputGroup(c, &v, o)
				}
			})
		}},
	{ID: "044", Name: "PinInput", ZhName: "验证码输入", Category: "Forms", Code: formCode["PinInput"],
		Keys: "Type digits to fill cells, Backspace clears the last, Cmd/Ctrl+V pastes the whole code.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "044", func() pinDemo { return pinDemo{} })
			formToggles(c, []string{"4 digits", "Disabled", "Invalid"}, &d.length4, &d.disabled, &d.invalid)
			n := 6
			if d.length4 {
				n = 4
			}
			if len(d.code) > n {
				d.code = d.code[:n]
			}
			p := mujica.PinInput(c, &d.code, mujica.PinInputOptions{Length: n, Disabled: d.disabled, Invalid: d.invalid})
			if p.Changed() {
				log(fmt.Sprintf("code %q", d.code))
			}
			if p.Completed() {
				log(fmt.Sprintf("completed %q", d.code))
			}
			formStates(c, func() {
				a, b := "12", "0420"
				mujica.PinInput(c, &a, mujica.PinInputOptions{Length: 4, Disabled: true, Label: "Disabled code"})
				mujica.PinInput(c, &b, mujica.PinInputOptions{Length: 4, Invalid: true, Label: "Invalid code"})
			})
		}},
	{ID: "045", Name: "TagsInput", ZhName: "标签输入", Category: "Forms", Code: formCode["TagsInput"],
		Keys: "Enter or a comma adds a tag; Backspace in the empty input removes the last; Tab reaches each remove button.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "045", func() tagsDemo { return tagsDemo{tags: []string{"rose", "thorn"}} })
			formToggles(c, []string{"At most 5", "Disabled"}, &d.limit, &d.disabled)
			opts := mujica.TagsInputOptions{Placeholder: "Add a flower", Label: "Flowers", Disabled: d.disabled,
				Validate: func(s string) string {
					if len([]rune(s)) > 16 {
						return "Keep tags under 16 characters."
					}
					return ""
				}}
			if d.limit {
				opts.Max = 5
			}
			if mujica.TagsInput(c, &d.tags, opts).Changed() {
				log(fmt.Sprintf("tags %q", d.tags))
			}
		}},
	{ID: "046", Name: "Rating", ZhName: "评分", Category: "Forms", Code: formCode["Rating"],
		Keys: "Arrows change by one step, Home clears, End fills.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "046", func() ratingDemo { return ratingDemo{value: 3, half: true} })
			formToggles(c, []string{"Half stars", "Read-only", "Disabled"}, &d.half, &d.readOnly, &d.disabled)
			step := 1.0
			if d.half {
				step = 0.5
			}
			if !d.half {
				d.value = float64(int(d.value + 0.5))
			}
			if mujica.Rating(c, &d.value, mujica.RatingOptions{Step: step, ReadOnly: d.readOnly, Disabled: d.disabled, ShowValue: true}).Changed() {
				log(fmt.Sprintf("rating %v", d.value))
			}
			formStates(c, func() {
				a, b := 4.5, 2.0
				mujica.Rating(c, &a, mujica.RatingOptions{Step: 0.5, ReadOnly: true, Label: "Read-only rating"})
				mujica.Rating(c, &b, mujica.RatingOptions{Disabled: true, Label: "Disabled rating"})
			})
		}},
	{ID: "047", Name: "ColorPicker", ZhName: "颜色选择器", Category: "Forms", Code: formCode["ColorPicker"],
		Keys: "Enter or Space opens; arrows move in the square and sliders; Escape closes and returns focus.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "047", func() colorDemo { return colorDemo{color: ui.Hex("#7D2034")} })
			formToggles(c, []string{"Disabled"}, &d.disabled)
			if mujica.ColorPicker(c, &d.color, mujica.ColorPickerOptions{Label: "Banner color", Disabled: d.disabled}).Changed() {
				log(fmt.Sprintf("color %v", d.color))
			}
		}},
	{ID: "048", Name: "FilePicker", ZhName: "文件选择", Category: "Forms", Code: formCode["FilePicker"],
		Keys: "Enter or Space opens the system dialog; the dialog's own keys apply.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "048", func() pickerDemo { return pickerDemo{multiple: true} })
			formToggles(c, []string{"Folders", "Multiple", "PDF and MIDI only", "Disabled"}, &d.directory, &d.multiple, &d.filter, &d.disabled)
			opts := mujica.FilePickerOptions{Directory: d.directory, Multiple: d.multiple, Disabled: d.disabled, Title: "Choose scores"}
			if d.filter {
				opts.Filters = []mygo.FileFilter{{Name: "Scores", Extensions: []string{"pdf", "mid"}}}
			}
			p := mujica.FilePicker(c, &d.paths, opts)
			if p.Changed() {
				log(fmt.Sprintf("picked %d: %q", len(d.paths), d.paths))
			}
			if p.Canceled() {
				log("canceled (not an error)")
			}
		}},
	{ID: "049", Name: "FileDropZone", ZhName: "文件拖放区", Category: "Forms", Code: formCode["FileDropZone"],
		Keys: "Tab reaches Browse… and each remove button; dropping needs a pointer, Browse is the keyboard path.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "049", func() dropDemo { return dropDemo{} })
			formToggles(c, []string{"Images only", "Disabled"}, &d.images, &d.disabled)
			opts := mujica.FileDropZoneOptions{MaxFiles: 4, MaxSize: 10 << 20, Disabled: d.disabled}
			if d.images {
				opts.Extensions = []string{"png", "jpg", "jpeg", "gif"}
			}
			if mujica.FileDropZone(c, &d.files, opts).Changed() {
				log(fmt.Sprintf("files %d", len(d.files)))
			}
		}},
	{ID: "050", Name: "Transfer", ZhName: "穿梭框", Category: "Forms", Code: formCode["Transfer"],
		Keys: "Up and Down move within a list, Space checks; Tab reaches the move buttons.",
		Demo: func(c *ui.Context, log func(string)) {
			d := ui.Local(c.Root(), "050", func() transferDemo { return transferDemo{chosen: []string{"sakiko"}} })
			formToggles(c, []string{"Disabled"}, &d.disabled)
			if mujica.Transfer(c, &d.chosen, transferDemoItems, mujica.TransferOptions{SourceTitle: "Auditions", TargetTitle: "Band", Disabled: d.disabled}).Changed() {
				log(fmt.Sprintf("band %q", d.chosen))
			}
		}},
}

// formCode holds each page's copyable snippet, the same as its doc and Example.
var formCode = map[string]string{
	"Form": `var state mujica.FormState
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
ui.Render(view, 400, 300, 1)`,
	"FormField": `email := ""
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FormField(c, "Email", mujica.FormFieldOptions{Required: true, Description: "For receipts."}, func() *ui.Element {
		return mujica.InputGroup(c, &email, mujica.InputGroupOptions{}).Input
	})
}
ui.Render(view, 400, 300, 1)`,
	"InputGroup": `site := ""
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
ui.Render(view, 400, 300, 1)`,
	"PinInput": `code := ""
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.PinInput(c, &code, mujica.PinInputOptions{Length: 6}).Completed() {
		fmt.Println("verify", code)
	}
}
ui.Render(view, 400, 300, 1)`,
	"TagsInput": `tags := []string{"gothic"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TagsInput(c, &tags, mujica.TagsInputOptions{Max: 5, Placeholder: "Add a tag", Label: "Tags"})
}
ui.Render(view, 400, 300, 1)`,
	"Rating": `stars := 3.5
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Rating(c, &stars, mujica.RatingOptions{Step: 0.5, ShowValue: true})
}
ui.Render(view, 400, 300, 1)`,
	"ColorPicker": `tint := ui.Hex("#7D2034")
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ColorPicker(c, &tint, mujica.ColorPickerOptions{Label: "Tint"})
}
ui.Render(view, 400, 300, 1)`,
	"FilePicker": `var paths []string
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FilePicker(c, &paths, mujica.FilePickerOptions{
		Multiple: true,
		Filters:  []mygo.FileFilter{{Name: "Scores", Extensions: []string{"pdf", "mid"}}},
	})
}
ui.Render(view, 400, 300, 1)`,
	"FileDropZone": `var files []mujica.DroppedFile
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FileDropZone(c, &files, mujica.FileDropZoneOptions{
		Extensions: []string{"png", "jpg"}, MaxFiles: 3, MaxSize: 5 << 20,
	})
}
ui.Render(view, 400, 300, 1)`,
	"Transfer": `chosen := []string{"uika"}
items := []mujica.TransferItem{{ID: "uika", Label: "Uika"}, {ID: "mutsumi", Label: "Mutsumi"}, {ID: "nyamu", Label: "Nyamu", Disabled: true}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Transfer(c, &chosen, items, mujica.TransferOptions{})
}
ui.Render(view, 400, 300, 1)`,
}
