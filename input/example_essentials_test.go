package input_test

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func ExampleCopyButton() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.CopyButton(c, "https://example.com/invite/42", input.CopyButtonOptions{Label: "Copy link"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleHoldToConfirm() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.HoldToConfirm(c, "Delete project", input.HoldToConfirmOptions{}).Changed() {
			// delete the project
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBulkActionBar() {
	selected := 3
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := input.BulkActionBar(c, selected, input.BulkActionBarOptions{}, func() {
			input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary})
		})
		if r.Changed() {
			selected = 0
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleScrubInput() {
	opacity := 0.6
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.ScrubInput(c, "Opacity", &opacity, input.ScrubInputOptions{Max: 1, Step: 0.1})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCurrencyInput() {
	cents := int64(123450)
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.CurrencyInput(c, &cents, input.CurrencyInputOptions{Symbol: "$", Label: "Price"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMaskedInput() {
	phone := ""
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.MaskedInput(c, &phone, input.MaskedInputOptions{Mask: "(999) 999-9999", Label: "Phone"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMentionInput() {
	text, mentions := "", []input.Mention(nil)
	people := []input.MentionCandidate{{ID: "u1", Name: "Alice"}, {ID: "u2", Name: "Bob"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.MentionInput(c, &text, &mentions, input.MentionInputOptions{
			Label: "Comment",
			Candidates: func(q string) []input.MentionCandidate {
				var out []input.MentionCandidate
				for _, p := range people {
					if strings.HasPrefix(strings.ToLower(p.Name), strings.ToLower(q)) {
						out = append(out, p)
					}
				}
				return out
			},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleChoiceChips() {
	var status []string
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.ChoiceChips(c, &status, []input.ToggleItem[string]{
			{Value: "open", Label: "Open"}, {Value: "closed", Label: "Closed"},
		}, input.ChoiceChipsOptions{Multiple: true, Label: "Status"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleKnob() {
	gain := 35.0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.Knob(c, &gain, input.KnobOptions{Label: "Gain", Ticks: 11})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVerticalSlider() {
	level := 40.0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.VerticalSlider(c, &level, input.VerticalSliderOptions{Step: 10, Label: "Level"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleColorPalette() {
	color := ui.Hex("#7D2034")
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.ColorPalette(c, &color, []input.ColorSwatch{
			{Color: ui.Hex("#7D2034"), Name: "Wine"},
			{Color: ui.Hex("#AD9164"), Name: "Old gold"},
			{Color: ui.Hex("#355E80"), Name: "Slate blue"},
		}, input.ColorPaletteOptions{Custom: true, Label: "Tag color"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSignaturePad() {
	var strokes [][]input.SignaturePoint
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.SignaturePad(c, &strokes, input.SignaturePadOptions{}).Changed() {
			_ = input.SignaturePath(strokes) // SVG path data to store
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBanner() {
	open := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := input.Banner(c, &open, "Maintenance tonight at 22:00.", input.BannerOptions{
			Severity: core.SeverityWarning,
			Action:   "Details",
		})
		if r.Action {
			// show the details
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMasonry() {
	titles := []string{"Sea", "Forest", "Desert", "Glacier"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.Masonry(c, len(titles), input.MasonryOptions{
			Key: func(i int) string { return titles[i] },
		}, func(i int) {
			ui.Text(c, titles[i])
		}).Size(380, 280)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleResizablePanelGroup() {
	var sizes []float32
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := input.ResizablePanelGroup(c, &sizes, []input.ResizablePanel{
			{Label: "Files", Collapsible: true, Build: func() { ui.Text(c, "Files") }},
			{Label: "Editor", Build: func() { ui.Text(c, "Editor") }},
		}, input.ResizablePanelGroupOptions{})
		if r.Committed {
			// save sizes
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDockPanel() {
	var layout input.DockLayout
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := input.DockPanel(c, &layout, []input.DockItem{
			{ID: "files", Title: "Files", Build: func() { ui.Text(c, "Files") }},
			{ID: "terminal", Title: "Terminal", Build: func() { ui.Text(c, "$") }},
		}, input.DockPanelOptions{})
		if r.Changed() {
			// save layout
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleQRCode() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		input.QRCode(c, "https://example.com/mujica", input.QRCodeOptions{Level: input.QRQuartile})
	}
	ui.Render(view, 400, 300, 1)
}
