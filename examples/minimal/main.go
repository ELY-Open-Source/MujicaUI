package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

type app struct{ name string }

func (a *app) view(c *ui.Context) {
	mujica.Use(c, mujica.Settings{Locale: mujica.ZhCN})
	ui.Column(c).Fill().Center().Gap(12).Children(func() {
		mujica.TextInput(c, &a.name, mujica.TextInputOptions{Placeholder: "姓名"})
		if mujica.Button(c, "觐见", mujica.ButtonOptions{}).Clicked() {
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
