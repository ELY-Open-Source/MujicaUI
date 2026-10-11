package overlay_test

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
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
