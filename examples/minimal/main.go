package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type app struct{ name string }

func (a *app) view(c *ui.Context) {
	core.Use(c, core.Settings{Locale: core.ZhCN})
	ui.Column(c).Fill().Center().Gap(12).Children(func() {
		input.TextInput(c, &a.name, input.TextInputOptions{Placeholder: "姓名"})
		if input.Button(c, "觐见", input.ButtonOptions{}).Clicked() {
			log.Println("hello", a.name)
		}
	})
}

func main() {
	a := &app{}
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{Title: "Mujica", Width: 480, Height: 320, Content: ui.View(a.view)})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
