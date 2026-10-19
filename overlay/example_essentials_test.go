package overlay_test

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/overlay"
)

func ExampleSplitButton() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := overlay.SplitButton(c, "Save", overlay.SplitButtonOptions{}, func(m *overlay.PopupMenu) {
			if m.Item("Save as…", overlay.PopupMenuItemOptions{}) {
				// ask for a file name
			}
		})
		if r.Changed() {
			// save
		}
	}
	ui.Render(view, 400, 300, 1)
}
