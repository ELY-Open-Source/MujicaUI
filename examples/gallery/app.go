package main

import (
	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

type gallery struct {
	pages    []page
	selected int
	events   []string
}

func newGallery() *gallery { return &gallery{pages: catalog()} }

func (g *gallery) log(s string) { g.events = append(g.events, s) }

func (g *gallery) view(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	ui.Row(c).Fill().Children(func() {
		ui.Scroll(c).Width(240).FillHeight().Children(func() {
			for i, p := range g.pages {
				if mujica.Button(c, p.ID+" "+p.Name, mujica.ButtonOptions{Variant: mujica.Ghost}).Key(p.ID).Clicked() {
					g.selected = i
				}
			}
		})
		ui.Scroll(c).Fill().Padding(24).Children(func() {
			if len(g.pages) > 0 {
				g.pages[g.selected].Demo(c, g.log)
			}
		})
	})
}
