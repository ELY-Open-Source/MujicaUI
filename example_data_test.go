package mujica_test

import (
	"cmp"
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleCard() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Card(c, mujica.CardOptions{Title: "Archive", Decorated: true, Footer: func() {
			mujica.Button(c, "Open", mujica.ButtonOptions{})
		}}, func() {
			ui.Text(c, "Letters of the court, sealed in wax.")
		}).Width(280)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleList() {
	names := []string{"Matins", "Lauds", "Prime", "Terce"}
	var state mujica.ListState[string]
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		v := mujica.List(c, &state, len(names), mujica.ListOptions[string]{
			Key:   func(i int) string { return names[i] },
			Label: func(i int) string { return names[i] },
		}, func(i int) { ui.Text(c, names[i]) })
		v.Size(240, 200)
		if v.Changed() {
			key, _ := state.Selected()
			fmt.Println("chose", key)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDataTable() {
	type relic struct {
		ID   int
		Name string
		Age  int
	}
	rows := []relic{{1, "Censer", 300}, {2, "Altar", 120}, {3, "Bell", 80}}
	var state mujica.DataTableState[int]
	opts := mujica.DataTableOptions[relic, int]{
		Key: func(r relic) int { return r.ID },
		Columns: []mujica.DataColumn[relic]{
			{ID: "name", Title: "Name", Sortable: true,
				Text:    func(r relic) string { return r.Name },
				Compare: func(a, b relic) int { return cmp.Compare(a.Name, b.Name) }},
			{ID: "age", Title: "Age", Width: 80, Align: ui.End, Sortable: true,
				Text:    func(r relic) string { return strconv.Itoa(r.Age) },
				Compare: func(a, b relic) int { return cmp.Compare(a.Age, b.Age) }},
		},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DataTable(c, &state, rows, opts).Size(360, 240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTree() {
	kids := map[string][]string{"Nave": {"Pulpit", "Pews"}}
	var state mujica.TreeState[string]
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Tree(c, &state, mujica.TreeOptions[string]{
			Roots:    []string{"Nave", "Sacristy"},
			Children: func(n string) []string { return kids[n] },
			Label:    func(n string) string { return n },
		}).Size(280, 240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTreeTable() {
	kids := map[string][]string{"Vestry": {"Stole", "Alb"}}
	size := map[string]int{"Vestry": 9, "Stole": 2, "Alb": 5, "Organ": 7}
	var state mujica.TreeTableState[string]
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TreeTable(c, &state, mujica.TreeTableOptions[string]{
			Roots:    []string{"Vestry", "Organ"},
			Children: func(n string) []string { return kids[n] },
			Label:    func(n string) string { return n },
			Columns: []mujica.TreeColumn[string]{
				{ID: "name", Title: "Name", Cell: func(n string) { ui.Text(c, n) },
					Compare: cmp.Compare[string]},
				{ID: "size", Title: "Size", Width: 80, Align: ui.End,
					Cell: func(n string) { ui.Text(c, strconv.Itoa(size[n])) }},
			},
		}).Size(360, 240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAccordion() {
	open := []string{"choir"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Accordion(c, &open, []mujica.AccordionItem{
			{ID: "choir", Title: "Choir", Content: func() { ui.Text(c, "Twelve voices.") }},
			{ID: "nave", Title: "Nave", Content: func() { ui.Text(c, "Black oak pews.") }},
		}, mujica.AccordionOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCollapsible() {
	open := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Collapsible(c, &open, mujica.CollapsibleOptions{Title: "Details"}, func() {
			ui.Text(c, "Candles are lit at the seventh bell.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTimeline() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Timeline(c, []mujica.TimelineItem{
			{Time: "06:00", Title: "Matins sung", Status: mujica.TimelineDone},
			{Time: "09:00", Title: "Relics moved", Content: "From the crypt to the altar.", Status: mujica.TimelineActive},
			{Time: "18:00", Title: "Vespers"},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDescriptionList() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DescriptionList(c, []mujica.DescriptionItem{
			{Label: "Archivist", Value: "Sister Agnès", Copyable: true},
			{Label: "Shelf", Value: "North vault, row 7"},
			{Label: "Seal", Value: "Wine-red wax", Span: -1},
		}, mujica.DescriptionListOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStatistic() {
	change := 12.5
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Statistic(c, 12840, mujica.StatisticOptions{
			Title: "Pilgrims", Unit: "souls", Precision: 0,
			Change: &change, ChangeUnit: "%", ChangeNote: "vs last feast",
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleImageViewer() {
	picture := ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="240" height="160"><rect width="240" height="160" fill="#AD9164"/><path d="M120 20 L200 80 L120 140 L40 80 Z" fill="#7D2034"/></svg>`))
	open, index := true, 0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ImageViewer(c, &open, &index, []mujica.ViewerImage{{Source: picture, Alt: "Rose window"}})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCarousel() {
	index := 0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Carousel(c, &index, 3, mujica.CarouselOptions{Label: "Windows", Autoplay: true}, func(i int) {
			ui.Text(c, fmt.Sprintf("Rose window %d", i+1)).Padding(24)
		}).Size(360, 220)
	}
	ui.Render(view, 400, 300, 1)
}
