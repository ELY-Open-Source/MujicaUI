// Gallery tours every MujicaUI component in a native MyGo window.
//
//	go run ./examples/gallery            native window
//	go run ./examples/gallery -web :8080 browser mirror
package main

import (
	"flag"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	web := flag.String("web", "", "serve the gallery to a browser at this address instead of opening a window")
	flag.Parse()
	app := newGallery()
	if *web != "" {
		log.Fatal(serveWeb(*web, app.view))
	}
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
