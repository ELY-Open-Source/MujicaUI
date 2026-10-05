package mujica_test

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleSplitButton() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.SplitButton(c, "Save", mujica.SplitButtonOptions{}, func(m *mujica.PopupMenu) {
			if m.Item("Save as…", mujica.PopupMenuItemOptions{}) {
				// ask for a file name
			}
		})
		if r.Changed() {
			// save
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCopyButton() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CopyButton(c, "https://example.com/invite/42", mujica.CopyButtonOptions{Label: "Copy link"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleHoldToConfirm() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.HoldToConfirm(c, "Delete project", mujica.HoldToConfirmOptions{}).Changed() {
			// delete the project
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBulkActionBar() {
	selected := 3
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.BulkActionBar(c, selected, mujica.BulkActionBarOptions{}, func() {
			mujica.Button(c, "Archive", mujica.ButtonOptions{Variant: mujica.Secondary})
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
		mujica.Use(c, mujica.Settings{})
		mujica.ScrubInput(c, "Opacity", &opacity, mujica.ScrubInputOptions{Max: 1, Step: 0.1})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCurrencyInput() {
	cents := int64(123450)
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CurrencyInput(c, &cents, mujica.CurrencyInputOptions{Symbol: "$", Label: "Price"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMaskedInput() {
	phone := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MaskedInput(c, &phone, mujica.MaskedInputOptions{Mask: "(999) 999-9999", Label: "Phone"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMentionInput() {
	text, mentions := "", []mujica.Mention(nil)
	people := []mujica.MentionCandidate{{ID: "u1", Name: "Alice"}, {ID: "u2", Name: "Bob"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MentionInput(c, &text, &mentions, mujica.MentionInputOptions{
			Label: "Comment",
			Candidates: func(q string) []mujica.MentionCandidate {
				var out []mujica.MentionCandidate
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
		mujica.Use(c, mujica.Settings{})
		mujica.ChoiceChips(c, &status, []mujica.ToggleItem[string]{
			{Value: "open", Label: "Open"}, {Value: "closed", Label: "Closed"},
		}, mujica.ChoiceChipsOptions{Multiple: true, Label: "Status"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleKnob() {
	gain := 35.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Knob(c, &gain, mujica.KnobOptions{Label: "Gain", Ticks: 11})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVerticalSlider() {
	level := 40.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.VerticalSlider(c, &level, mujica.VerticalSliderOptions{Step: 10, Label: "Level"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleColorPalette() {
	color := ui.Hex("#7D2034")
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ColorPalette(c, &color, []mujica.ColorSwatch{
			{Color: ui.Hex("#7D2034"), Name: "Wine"},
			{Color: ui.Hex("#AD9164"), Name: "Old gold"},
			{Color: ui.Hex("#355E80"), Name: "Slate blue"},
		}, mujica.ColorPaletteOptions{Custom: true, Label: "Tag color"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSignaturePad() {
	var strokes [][]mujica.SignaturePoint
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.SignaturePad(c, &strokes, mujica.SignaturePadOptions{}).Changed() {
			_ = mujica.SignaturePath(strokes) // SVG path data to store
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBanner() {
	open := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.Banner(c, &open, "Maintenance tonight at 22:00.", mujica.BannerOptions{
			Severity: mujica.SeverityWarning,
			Action:   "Details",
		})
		if r.Action {
			// show the details
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleUndoToast() {
	open := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.UndoToast(c, &open, "Note deleted", mujica.UndoToastOptions{
			OnUndo: func() {
				// restore the note
			},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMasonry() {
	titles := []string{"Sea", "Forest", "Desert", "Glacier"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Masonry(c, len(titles), mujica.MasonryOptions{
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
		mujica.Use(c, mujica.Settings{})
		r := mujica.ResizablePanelGroup(c, &sizes, []mujica.ResizablePanel{
			{Label: "Files", Collapsible: true, Build: func() { ui.Text(c, "Files") }},
			{Label: "Editor", Build: func() { ui.Text(c, "Editor") }},
		}, mujica.ResizablePanelGroupOptions{})
		if r.Committed {
			// save sizes
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDockPanel() {
	var layout mujica.DockLayout
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.DockPanel(c, &layout, []mujica.DockItem{
			{ID: "files", Title: "Files", Build: func() { ui.Text(c, "Files") }},
			{ID: "terminal", Title: "Terminal", Build: func() { ui.Text(c, "$") }},
		}, mujica.DockPanelOptions{})
		if r.Changed() {
			// save layout
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleQRCode() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.QRCode(c, "https://example.com/mujica", mujica.QRCodeOptions{Level: mujica.QRQuartile})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleActivityFeed() {
	events := []mujica.Activity{
		{ID: "1", Actor: "Lin", Action: "commented on", Object: "Quarterly report", Time: time.Now()},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ActivityFeed(c, events, mujica.ActivityFeedOptions{}).Size(380, 280)
	}
	ui.Render(view, 400, 300, 1)
}
