package main

// Snippets are the bodies of the Example functions in example_data_test.go.

const dataCardCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Card(c, data.CardOptions{Title: "Archive", Decorated: true, Footer: func() {
		input.Button(c, "Open", input.ButtonOptions{})
	}}, func() {
		ui.Text(c, "Letters of the court, sealed in wax.")
	}).Width(280)
}
ui.Render(view, 400, 300, 1)`

const dataListCode = `names := []string{"Matins", "Lauds", "Prime", "Terce"}
var state data.ListState[string]
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	v := data.List(c, &state, len(names), data.ListOptions[string]{
		Key:   func(i int) string { return names[i] },
		Label: func(i int) string { return names[i] },
	}, func(i int) { ui.Text(c, names[i]) })
	v.Size(240, 200)
	if v.Changed() {
		key, _ := state.Selected()
		fmt.Println("chose", key)
	}
}
ui.Render(view, 400, 300, 1)`

const dataTableCode = `type relic struct {
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
	data.DataTable(c, &state, rows, opts).Size(360, 240)
}
ui.Render(view, 400, 300, 1)`

const dataTreeCode = `kids := map[string][]string{"Nave": {"Pulpit", "Pews"}}
var state data.TreeState[string]
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Tree(c, &state, data.TreeOptions[string]{
		Roots:    []string{"Nave", "Sacristy"},
		Children: func(n string) []string { return kids[n] },
		Label:    func(n string) string { return n },
	}).Size(280, 240)
}
ui.Render(view, 400, 300, 1)`

const dataTreeTableCode = `kids := map[string][]string{"Vestry": {"Stole", "Alb"}}
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
	}).Size(360, 240)
}
ui.Render(view, 400, 300, 1)`

const dataAccordionCode = `open := []string{"choir"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Accordion(c, &open, []data.AccordionItem{
		{ID: "choir", Title: "Choir", Content: func() { ui.Text(c, "Twelve voices.") }},
		{ID: "nave", Title: "Nave", Content: func() { ui.Text(c, "Black oak pews.") }},
	}, data.AccordionOptions{})
}
ui.Render(view, 400, 300, 1)`

const dataCollapsibleCode = `open := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Collapsible(c, &open, data.CollapsibleOptions{Title: "Details"}, func() {
		ui.Text(c, "Candles are lit at the seventh bell.")
	})
}
ui.Render(view, 400, 300, 1)`

const dataTimelineCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Timeline(c, []data.TimelineItem{
		{Time: "06:00", Title: "Matins sung", Status: data.TimelineDone},
		{Time: "09:00", Title: "Relics moved", Content: "From the crypt to the altar.", Status: data.TimelineActive},
		{Time: "18:00", Title: "Vespers"},
	})
}
ui.Render(view, 400, 300, 1)`

const dataDescriptionCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.DescriptionList(c, []data.DescriptionItem{
		{Label: "Archivist", Value: "Sister Agnès", Copyable: true},
		{Label: "Shelf", Value: "North vault, row 7"},
		{Label: "Seal", Value: "Wine-red wax", Span: -1},
	}, data.DescriptionListOptions{})
}
ui.Render(view, 400, 300, 1)`

const dataStatisticCode = `change := 12.5
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Statistic(c, 12840, data.StatisticOptions{
		Title: "Pilgrims", Unit: "souls", Precision: 0,
		Change: &change, ChangeUnit: "%", ChangeNote: "vs last feast",
	})
}
ui.Render(view, 400, 300, 1)`

const dataViewerCode = "picture := ui.MustParseSVG([]byte(`<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"240\" height=\"160\"><rect width=\"240\" height=\"160\" fill=\"#AD9164\"/><path d=\"M120 20 L200 80 L120 140 L40 80 Z\" fill=\"#7D2034\"/></svg>`))\nopen, index := true, 0\nview := func(c *ui.Context) {\n\tcore.Use(c, core.Settings{})\n\tdata.ImageViewer(c, &open, &index, []data.ViewerImage{{Source: picture, Alt: \"Rose window\"}})\n}\nui.Render(view, 400, 300, 1)"

const dataCarouselCode = `index := 0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Carousel(c, &index, 3, data.CarouselOptions{Label: "Windows", Autoplay: true}, func(i int) {
		ui.Text(c, fmt.Sprintf("Rose window %d", i+1)).Padding(24)
	}).Size(360, 220)
}
ui.Render(view, 400, 300, 1)`
