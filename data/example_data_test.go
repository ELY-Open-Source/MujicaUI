package data_test

import (
	"cmp"
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func ExampleCard() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.Card(c, data.CardOptions{Title: "Archive", Decorated: true, Footer: func() {
			input.Button(c, "Open", input.ButtonOptions{})
		}}, func() {
			ui.Text(c, "Letters of the court, sealed in wax.")
		}).Width(280)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleList() {
	names := []string{"Matins", "Lauds", "Prime", "Terce"}
	var state data.ListState[string]
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		v := data.List(c, &state, len(names), data.ListOptions[string]{
			Key:   func(i int) string { return names[i] },
			Label: func(i int) string { return names[i] },
		}, func(i int) { ui.Text(c, names[i]) })
		v.Element.Size(240, 200)
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
	var state data.DataTableState[int]
	opts := data.DataTableOptions[relic, int]{
		Key: func(r relic) int { return r.ID },
		Columns: []data.DataColumn[relic]{
			{ID: "name", Title: "Name", Sortable: true,
				Text:    func(r relic) string { return r.Name },
				Compare: func(a, b relic) int { return cmp.Compare(a.Name, b.Name) }},
			{ID: "age", Title: "Age", Width: 80, Align: ui.End, Sortable: true,
				Text:    func(r relic) string { return strconv.Itoa(r.Age) },
				Compare: func(a, b relic) int { return cmp.Compare(a.Age, b.Age) }},
		},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.DataTable(c, &state, rows, opts).Element.Size(360, 240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTree() {
	kids := map[string][]string{"Nave": {"Pulpit", "Pews"}}
	var state data.TreeState[string]
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.Tree(c, &state, data.TreeOptions[string]{
			Roots:    []string{"Nave", "Sacristy"},
			Children: func(n string) []string { return kids[n] },
			Label:    func(n string) string { return n },
		}).Element.Size(280, 240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTreeTable() {
	kids := map[string][]string{"Vestry": {"Stole", "Alb"}}
	size := map[string]int{"Vestry": 9, "Stole": 2, "Alb": 5, "Organ": 7}
	var state data.TreeTableState[string]
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.TreeTable(c, &state, data.TreeTableOptions[string]{
			Roots:    []string{"Vestry", "Organ"},
			Children: func(n string) []string { return kids[n] },
			Label:    func(n string) string { return n },
			Columns: []data.TreeColumn[string]{
				{ID: "name", Title: "Name", Cell: func(n string) { ui.Text(c, n) },
					Compare: cmp.Compare[string]},
				{ID: "size", Title: "Size", Width: 80, Align: ui.End,
					Cell: func(n string) { ui.Text(c, strconv.Itoa(size[n])) }},
			},
		}).Element.Size(360, 240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAccordion() {
	open := []string{"choir"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.Accordion(c, &open, []data.AccordionItem{
			{ID: "choir", Title: "Choir", Content: func() { ui.Text(c, "Twelve voices.") }},
			{ID: "nave", Title: "Nave", Content: func() { ui.Text(c, "Black oak pews.") }},
		}, data.AccordionOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCollapsible() {
	open := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.Collapsible(c, &open, data.CollapsibleOptions{Title: "Details"}, func() {
			ui.Text(c, "Candles are lit at the seventh bell.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTimeline() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.Timeline(c, []data.TimelineItem{
			{Time: "06:00", Title: "Matins sung", Status: data.TimelineDone},
			{Time: "09:00", Title: "Relics moved", Content: "From the crypt to the altar.", Status: data.TimelineActive},
			{Time: "18:00", Title: "Vespers"},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDescriptionList() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.DescriptionList(c, []data.DescriptionItem{
			{Label: "Archivist", Value: "Sister Agnès", Copyable: true},
			{Label: "Shelf", Value: "North vault, row 7"},
			{Label: "Seal", Value: "Wine-red wax", Span: -1},
		}, data.DescriptionListOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStatistic() {
	change := 12.5
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.Statistic(c, 12840, data.StatisticOptions{
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
		core.Use(c, core.Settings{})
		data.ImageViewer(c, &open, &index, []data.ViewerImage{{Source: picture, Alt: "Rose window"}})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCarousel() {
	index := 0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		data.Carousel(c, &index, 3, data.CarouselOptions{Label: "Windows", Autoplay: true}, func(i int) {
			ui.Text(c, fmt.Sprintf("Rose window %d", i+1)).Padding(24)
		}).Element.Size(360, 220)
	}
	ui.Render(view, 400, 300, 1)
}
