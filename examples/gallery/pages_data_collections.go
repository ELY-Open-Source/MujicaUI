package main

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

var dataOffices = []string{"Matins", "Lauds", "Prime", "Terce", "Sext", "None", "Vespers", "Compline"}

func dataListDemo(c *ui.Context, log func(string)) {
	type state struct {
		big, grid, multiple bool
		status, query       string
		list                mujica.ListState[int]
	}
	s := ui.Local(c.Root(), "070", func() state { return state{status: "Ready"} })
	n := 24
	if s.big {
		n = 100_000
	}
	label := func(i int) string { return fmt.Sprintf("%s %d", dataOffices[i%len(dataOffices)], i+1) }
	st, show := dataStatus(s.status)
	if !show {
		n = 0
	}
	dataSection(c, "List", func() {
		dataControls(c, func() {
			ui.Checkbox(c, &s.big, "100,000 items")
			ui.Checkbox(c, &s.grid, "Grid")
			ui.Checkbox(c, &s.multiple, "Multiple")
			ui.Select(c, &s.status, dataStatuses)
		})
		dataControls(c, func() {
			ui.TextInput(c, &s.query).Width(200).Placeholder("Find by label")
			if mujica.Button(c, "Find", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
				for i := range n {
					if strings.Contains(strings.ToLower(label(i)), strings.ToLower(s.query)) {
						s.list.ScrollTo(i)
						log(fmt.Sprintf("found %q at %d", s.query, i))
						break
					}
				}
			}
		})
		v := mujica.List(c, &s.list, n, mujica.ListOptions[int]{
			Key: func(i int) int { return i }, Label: label, Grid: s.grid, Multiple: s.multiple,
			ItemWidth: 140, ItemHeight: 72, Status: st,
		}, func(i int) {
			ui.Text(c, label(i)).SingleLine()
			if s.grid {
				ui.Text(c, "Office").TextColor(mujica.Tokens(c).TextMuted)
			}
		})
		v.Size(480, 320)
		if v.Changed() {
			k, _ := s.list.Selected()
			log(fmt.Sprintf("selected %d, %d chosen", k, s.list.Selection.Len()))
		}
		if v.Submitted() {
			k, _ := s.list.Selected()
			log(fmt.Sprintf("submitted %d", k))
		}
		if v.Retried() {
			log("retry")
		}
	})
}

type dataRelic struct {
	ID     int
	Name   string
	Office string
	Year   int
	Sealed bool
}

func dataRelics(n int) []dataRelic {
	rows := make([]dataRelic, n)
	for i := range rows {
		rows[i] = dataRelic{ID: i + 1, Name: fmt.Sprintf("Relic %05d", (i*7919)%n+1), Office: dataOffices[i%len(dataOffices)],
			Year: 1200 + (i*37)%600, Sealed: i%3 == 0}
	}
	return rows
}

func dataTableDemo(c *ui.Context, log func(string)) {
	type state struct {
		big, multiple bool
		status        string
		small, large  []dataRelic
		table         mujica.DataTableState[int]
	}
	s := ui.Local(c.Root(), "071", func() state { return state{status: "Ready", small: dataRelics(60)} })
	rows := s.small
	if s.big {
		if s.large == nil {
			s.large = dataRelics(100_000)
		}
		rows = s.large
	}
	st, show := dataStatus(s.status)
	if !show {
		rows = nil
	}
	k := mujica.Tokens(c)
	cols := []mujica.DataColumn[dataRelic]{
		{ID: "id", Title: "ID", Width: 80, Align: ui.End, Sortable: true, Text: func(r dataRelic) string { return strconv.Itoa(r.ID) },
			Compare: func(a, b dataRelic) int { return cmp.Compare(a.ID, b.ID) }},
		{ID: "name", Title: "Name", MinWidth: 120, Sortable: true, Text: func(r dataRelic) string { return r.Name },
			Compare: func(a, b dataRelic) int { return cmp.Compare(a.Name, b.Name) }},
		{ID: "office", Title: "Office", Width: 110, Sortable: true, Text: func(r dataRelic) string { return r.Office },
			Compare: func(a, b dataRelic) int { return cmp.Compare(a.Office, b.Office) }},
		{ID: "year", Title: "Year", Width: 80, Align: ui.End, Sortable: true, Text: func(r dataRelic) string { return strconv.Itoa(r.Year) },
			Compare: func(a, b dataRelic) int { return cmp.Compare(a.Year, b.Year) }},
		{ID: "seal", Title: "Seal", Width: 100, Text: func(r dataRelic) string { return map[bool]string{true: "Sealed", false: "Open"}[r.Sealed] },
			Cell: func(r dataRelic) {
				color := k.Success
				if r.Sealed {
					color = k.Warning
				}
				ui.Text(c, map[bool]string{true: "◆ Sealed", false: "◇ Open"}[r.Sealed]).TextColor(color).SingleLine()
			}},
	}
	dataSection(c, "DataTable", func() {
		dataControls(c, func() {
			ui.Checkbox(c, &s.big, "100,000 rows")
			ui.Checkbox(c, &s.multiple, "Multiple")
			ui.Select(c, &s.status, dataStatuses)
			ui.TextInput(c, &s.table.Filter).Width(180).Placeholder("Filter")
			mujica.DataTableColumnPicker(c, &s.table, cols)
		})
		v := mujica.DataTable(c, &s.table, rows, mujica.DataTableOptions[dataRelic, int]{
			Columns: cols, Key: func(r dataRelic) int { return r.ID }, Multiple: s.multiple, Status: st,
		})
		v.Size(620, 360)
		ui.Text(c, fmt.Sprintf("%d of %d rows shown", s.table.Rows(), len(rows))).TextColor(k.TextMuted)
		if v.SortChanged() {
			log(fmt.Sprintf("sort %s descending=%v", s.table.Sort.Column, s.table.Sort.Descending))
		}
		if v.Changed() {
			id, _ := s.table.Selected()
			log(fmt.Sprintf("selected row %d, %d chosen", id, s.table.Selection.Len()))
		}
		if v.Submitted() {
			id, _ := s.table.Selected()
			log(fmt.Sprintf("opened row %d", id))
		}
		if v.Retried() {
			log("retry")
		}
	})
}

// dataLazy simulates children that arrive 600 ms after they are asked for.
type dataLazy struct {
	kids   map[string][]string
	status map[string]mujica.DataStatus
	due    map[string]time.Time
	failed map[string]bool
}

func newDataLazy() *dataLazy {
	return &dataLazy{
		kids: map[string][]string{"Nave": {"Pulpit", "Pews", "Font"}, "Crypt": {}, "Archive": {},
			"Pulpit": nil, "Pews": nil, "Font": nil, "Sacristy": nil},
		status: map[string]mujica.DataStatus{"Crypt": mujica.DataUnloaded, "Archive": mujica.DataUnloaded},
		due:    map[string]time.Time{}, failed: map[string]bool{},
	}
}

func (d *dataLazy) load(c *ui.Context, log func(string)) func(string) {
	return func(n string) {
		d.status[n] = mujica.DataLoading
		d.due[n] = c.Now().Add(600 * time.Millisecond)
		log("load " + n)
	}
}

// tick finishes due loads; the Archive fails once to show retry.
func (d *dataLazy) tick(c *ui.Context) {
	for n, at := range d.due {
		if c.Now().Before(at) {
			c.After(at.Sub(c.Now()))
			continue
		}
		delete(d.due, n)
		if n == "Archive" && !d.failed[n] {
			d.failed[n], d.status[n] = true, mujica.DataFailed
			continue
		}
		d.status[n] = mujica.DataReady
		d.kids[n] = []string{n + " north", n + " south"}
	}
}

func dataTreeDemo(c *ui.Context, log func(string)) {
	type state struct {
		big, multiple bool
		lazy          *dataLazy
		tree          mujica.TreeState[string]
		bigTree       mujica.TreeState[string]
	}
	s := ui.Local(c.Root(), "072", func() state { return state{lazy: newDataLazy()} })
	s.lazy.tick(c)
	opts := mujica.TreeOptions[string]{
		Roots: []string{"Nave", "Crypt", "Archive", "Sacristy"}, Children: func(n string) []string { return s.lazy.kids[n] },
		Label: func(n string) string { return n }, ItemStatus: func(n string) mujica.DataStatus { return s.lazy.status[n] },
		Load: s.lazy.load(c, log), Multiple: s.multiple,
	}
	ts := &s.tree
	if s.big {
		opts = dataBigTree(s.multiple)
		ts = &s.bigTree
	}
	dataSection(c, "Tree", func() {
		dataControls(c, func() {
			ui.Checkbox(c, &s.big, "100,000 nodes")
			ui.Checkbox(c, &s.multiple, "Multiple")
		})
		v := mujica.Tree(c, ts, opts)
		v.Size(420, 340)
		if v.Changed() {
			n, _ := ts.Selected()
			log(fmt.Sprintf("selected %s, %d chosen", n, ts.Selection.Len()))
		}
		if v.Submitted() {
			n, _ := ts.Selected()
			log("submitted " + n)
		}
	})
}

// dataBigTree has 100 roots of 999 children: 100,000 nodes, built as they open.
func dataBigTree(multiple bool) mujica.TreeOptions[string] {
	roots := make([]string, 100)
	for i := range roots {
		roots[i] = fmt.Sprintf("Wing %d", i+1)
	}
	return mujica.TreeOptions[string]{
		Roots: roots, Label: func(n string) string { return n }, Multiple: multiple,
		Children: func(n string) []string {
			if !strings.HasPrefix(n, "Wing ") || strings.Contains(n, "·") {
				return nil
			}
			return dataWing(n)
		},
	}
}

var dataWings = map[string][]string{}

func dataWing(n string) []string {
	if kids, ok := dataWings[n]; ok {
		return kids
	}
	kids := make([]string, 999)
	for i := range kids {
		kids[i] = fmt.Sprintf("%s · cell %d", n, i+1)
	}
	dataWings[n] = kids
	return kids
}

func dataTreeTableDemo(c *ui.Context, log func(string)) {
	type state struct {
		multiple bool
		lazy     *dataLazy
		table    mujica.TreeTableState[string]
	}
	s := ui.Local(c.Root(), "073", func() state { return state{lazy: newDataLazy()} })
	s.lazy.tick(c)
	size := func(n string) int { return len(n)*7 + len(s.lazy.kids[n])*13 }
	text := func(n string) { ui.Text(c, n).SingleLine() }
	dataSection(c, "TreeTable", func() {
		dataControls(c, func() { ui.Checkbox(c, &s.multiple, "Multiple") })
		v := mujica.TreeTable(c, &s.table, mujica.TreeTableOptions[string]{
			Roots: []string{"Nave", "Crypt", "Archive", "Sacristy"}, Children: func(n string) []string { return s.lazy.kids[n] },
			Label: func(n string) string { return n }, ItemStatus: func(n string) mujica.DataStatus { return s.lazy.status[n] },
			Load: s.lazy.load(c, log), Multiple: s.multiple,
			Columns: []mujica.TreeColumn[string]{
				{ID: "name", Title: "Name", MinWidth: 160, Cell: text, Compare: cmp.Compare[string]},
				{ID: "size", Title: "Size", Width: 90, Align: ui.End, Cell: func(n string) { ui.Text(c, fmt.Sprintf("%d m²", size(n))) },
					Compare: func(a, b string) int { return cmp.Compare(size(a), size(b)) }},
				{ID: "kids", Title: "Rooms", Width: 80, Align: ui.End, Cell: func(n string) { ui.Text(c, strconv.Itoa(len(s.lazy.kids[n]))) }},
			},
		})
		v.Size(520, 320)
		if v.SortChanged() {
			log(fmt.Sprintf("sort %s descending=%v", s.table.Sort.Column, s.table.Sort.Descending))
		}
		if v.Changed() {
			n, _ := s.table.Selected()
			log(fmt.Sprintf("selected %s, %d chosen", n, s.table.Selection.Len()))
		}
	})
}
