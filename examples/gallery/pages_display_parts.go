package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/theme"
)

// displaySection shows a titled block of a demo.
func displaySection(c *ui.Context, title string, fn func()) {
	ui.Column(c).Gap(10).Margin(0, 0, 20, 0).Children(func() {
		mujica.Heading(c, title, mujica.HeadingOptions{Level: 3})
		fn()
	})
}

// displayControls lays out a row of parameter controls.
func displayControls(c *ui.Context, fn func()) {
	ui.Row(c).Gap(16).Wrap().AlignItems(ui.Center).Children(fn)
}

// displayCode shows a page's snippet with a copy button.
func displayCode(c *ui.Context, code string, log func(string)) {
	displaySection(c, "Code", func() {
		if mujica.Button(c, "Copy code", mujica.ButtonOptions{Variant: mujica.Secondary}).AlignSelf(ui.Start).Clicked() {
			c.WriteClipboard(code)
			log("code copied")
		}
		mujica.Text(c, code, mujica.TextOptions{Selectable: true}).Font(theme.MonoFont).FontSize(12)
	})
}

// displayLog logs a control's new value when it changed.
func displayLog[T any](log func(string), e *ui.Element, name string, v *T) {
	if e.Changed() {
		log(fmt.Sprintf("%s = %v", name, *v))
	}
}
