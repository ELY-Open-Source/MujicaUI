package core_test

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func ExampleSpinner() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		core.Spinner(c, core.SpinnerOptions{Size: 20, Label: "Syncing", ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}
