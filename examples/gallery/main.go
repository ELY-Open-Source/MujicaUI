// Gallery tours every MujicaUI component in a native MyGo window.
//
//	go run ./examples/gallery
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	app := newGallery()
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:     "MujicaUI Gallery",
			Width:     1200,
			Height:    800,
			MinWidth:  800,
			MinHeight: 600,
			Content:   ui.View(app.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
