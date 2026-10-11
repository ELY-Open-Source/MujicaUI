package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
)

// displaySection shows a titled block of a demo.
func displaySection(c *ui.Context, title string, fn func()) {
	ui.Column(c).Gap(10).Margin(0, 0, 20, 0).Children(func() {
		display.Heading(c, title, display.HeadingOptions{Level: 3})
		fn()
	})
}

// displayControls lays out a row of parameter controls.
func displayControls(c *ui.Context, fn func()) {
	ui.Row(c).Gap(16).Wrap().AlignItems(ui.Center).Children(fn)
}

// displayLog logs a control's new value when it changed.
func displayLog[T any](log func(string), e *ui.Element, name string, v *T) {
	if e.Changed() {
		log(fmt.Sprintf("%s = %v", name, *v))
	}
}
