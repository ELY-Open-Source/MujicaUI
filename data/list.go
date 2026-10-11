package data

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DataStatus is the load state of a data view's rows.
type DataStatus int

const (
	DataReady DataStatus = iota
	// DataLoading is the first load, before any row came.
	DataLoading
	// DataLoadingMore appends rows below those shown.
	DataLoadingMore
	// DataFailed shows the error with a retry action.
	DataFailed
	// DataUnloaded marks a tree item whose children were not asked for yet.
	DataUnloaded
)

// ListOptions configure a List of items keyed by K.
type ListOptions[K comparable] struct {
	// Key returns item i's stable identity; required.
	Key func(i int) K
	// Label returns item i's text, read out and matched by typing; required.
	Label    func(i int) string
	Multiple bool
	// Grid lays items out in columns at least ItemWidth wide, ItemHeight high.
	Grid                  bool
	ItemWidth, ItemHeight float32
	// Empty replaces the built-in "No items" text.
	Empty  string
	Status DataStatus
	// Error describes a DataFailed status; empty uses the built-in text.
	Error string
}

// ListState is a List's selection and place. Its zero value selects nothing.
type ListState[K comparable] struct {
	// Selection holds the chosen keys of a Multiple list.
	Selection ui.Selection[K]
	Rows      ui.ListState
	Grid      ui.GridState
	cursor    int
	started   bool
	key       K
	has       bool
}

// Selected returns the chosen item's key, or the last chosen of a Multiple list.
func (s *ListState[K]) Selected() (K, bool) { return s.key, s.has }

// ScrollTo brings item i to the top of the list.
func (s *ListState[K]) ScrollTo(i int) {
	s.Rows.ScrollTo(i, ui.Start)
	s.Grid.ScrollTo(i, ui.Start)
}

// ListView is a built List and what happened to it in this frame.
type ListView struct {
	Element                   *ui.Element
	changed, submitted, retry bool
}

// Changed reports a new selection.
func (v ListView) Changed() bool { return v.changed }

// Submitted reports Enter or a double click on an item.
func (v ListView) Submitted() bool { return v.submitted }

// Retried reports a click on a failed load's retry.
func (v ListView) Retried() bool { return v.retry }

// List shows n items, building only those in view, as rows or a grid.
func List[K comparable](c *ui.Context, s *ListState[K], n int, opts ListOptions[K], item func(i int)) ListView {
	if opts.Key == nil || opts.Label == nil {
		panic("mujica: List needs Key and Label")
	}
	if n < 0 {
		panic("mujica: List with negative item count")
	}
	if !s.started {
		s.started, s.cursor = true, -1
	}
	listLocate(s, n, opts.Key)
	k := core.Tokens(c)
	var v ListView
	v.Element = ui.Column(c).MinHeight(0)
	v.Element.Children(func() {
		if n == 0 {
			v.retry = DataPlaceholder(c, opts.Status, opts.Error, opts.Empty, core.Msg(c, "list.empty"))
			return
		}
		var e *ui.Element
		if opts.Grid {
			e = listGrid(c, s, n, opts, item)
		} else {
			e = listRows(c, s, n, opts, item)
		}
		e.Grow(1).MinHeight(0).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
		internal.Ring(e, k.Focus, theme.ControlRadius)
		v.changed, v.submitted = e.Changed(), e.Submitted()
		v.retry = dataFooter(c, opts.Status, opts.Error)
	})
	s.has = s.cursor >= 0 && s.cursor < n
	if s.has {
		s.key = opts.Key(s.cursor)
	}
	if v.changed {
		slog.Debug("mujica: list changed", "selected", s.cursor, "count", s.Selection.Len())
	}
	return v
}

// listLocate finds the chosen key again when items moved or came and went.
func listLocate[K comparable](s *ListState[K], n int, key func(int) K) {
	if !s.has || s.cursor >= 0 && s.cursor < n && key(s.cursor) == s.key {
		return
	}
	s.cursor = -1
	for i := range n {
		if key(i) == s.key {
			s.cursor = i
			break
		}
	}
	slog.Debug("mujica: list selection moved", "row", s.cursor)
}

func listRows[K comparable](c *ui.Context, s *ListState[K], n int, opts ListOptions[K], item func(i int)) *ui.Element {
	d := core.Density(c)
	s.Rows.Key = func(i int) any { return opts.Key(i) }
	s.Rows.Label = opts.Label
	s.Rows.Selected = &s.cursor
	if opts.Multiple {
		s.Rows.Selection = &s.Selection
	} else {
		s.Rows.Selection = nil
	}
	k := core.Tokens(c)
	restore := DataSelectTheme(c, 0)
	defer restore()
	return ui.List(c, &s.Rows, n, func(i int) {
		ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Row(c).Grow(1).MinWidth(0).AlignItems(ui.Center).Gap(d.Unit() * 2).Children(func() { item(i) })
			if listChosen(s, opts, i) {
				ui.Icon(c, icons.Must("check")).Size(16, 16).Shrink(0).TextColor(k.AccentText)
			}
		})
	}).Padding(d.Unit())
}

func listGrid[K comparable](c *ui.Context, s *ListState[K], n int, opts ListOptions[K], item func(i int)) *ui.Element {
	d := core.Density(c)
	w, h := opts.ItemWidth, opts.ItemHeight
	if w <= 0 {
		w = 160
	}
	if h <= 0 {
		h = 120
	}
	s.Grid.Key = func(i int) any { return opts.Key(i) }
	s.Grid.Label = opts.Label
	s.Grid.Selected = &s.cursor
	if opts.Multiple {
		s.Grid.Selection = &s.Selection
	} else {
		s.Grid.Selection = nil
	}
	// Chosen items take an AccentText ring over a Selection face.
	k := core.Tokens(c)
	saved := c.Theme()
	t := *saved
	t.Accent = k.AccentText
	c.SetTheme(&t)
	e := ui.GridView(c, &s.Grid, n, w, h, func(i int) {
		chosen := listChosen(s, opts, i)
		face := ui.Column(c).Grow(1).Padding(d.Unit() * 2).Gap(d.Unit()).Radius(theme.ControlRadius)
		if chosen {
			face.Background(k.Selection)
		}
		face.Children(func() {
			item(i)
			if chosen {
				ui.Icon(c, icons.Must("check")).Size(16, 16).Absolute().Top(d.Unit()).Right(d.Unit()).TextColor(k.AccentText)
			}
		})
	}).Padding(d.Unit() * 2)
	c.SetTheme(saved)
	return e
}

// DataSelectTheme shows MyGo's chosen rows as Selection with Text; spacing > 0 sets the unit.
func DataSelectTheme(c *ui.Context, spacing float32) func() {
	k := core.Tokens(c)
	saved := c.Theme()
	t := *saved
	t.Accent, t.AccentText = k.Selection, k.Text
	if spacing > 0 {
		t.Spacing = spacing
	}
	c.SetTheme(&t)
	return func() { c.SetTheme(saved) }
}

func listChosen[K comparable](s *ListState[K], opts ListOptions[K], i int) bool {
	if opts.Multiple {
		return s.Selection.Has(opts.Key(i))
	}
	return s.cursor == i
}

// DataPlaceholder fills an empty view with its status and reports a retry.
func DataPlaceholder(c *ui.Context, st DataStatus, errText, empty, emptyDefault string) bool {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	retry := false
	ui.Column(c).Grow(1).MinHeight(u * 30).Center().Gap(u * 3).Padding(u * 6).Children(func() {
		switch st {
		case DataLoading, DataLoadingMore, DataUnloaded:
			core.Spinner(c, core.SpinnerOptions{Label: core.Msg(c, "list.loading"), ShowLabel: true})
		case DataFailed:
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).TextColor(k.Danger).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(20, 20)
				ui.Text(c, DataError(c, errText)).Role(ui.RoleStatus)
			})
			retry = input.Button(c, core.Msg(c, "list.retry"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("refresh-cw")}).Clicked()
		default:
			if empty == "" {
				empty = emptyDefault
			}
			ui.Icon(c, icons.Must("inbox")).Size(24, 24).TextColor(k.TextMuted)
			ui.Text(c, empty).TextColor(k.TextMuted)
		}
	})
	return retry
}

// dataFooter shows appending and failed loads below rows and reports a retry.
func dataFooter(c *ui.Context, st DataStatus, errText string) bool {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	retry := false
	switch st {
	case DataLoadingMore:
		ui.Row(c).Padding(u*2, u*3).Shrink(0).Children(func() {
			core.Spinner(c, core.SpinnerOptions{Label: core.Msg(c, "list.loadingMore"), ShowLabel: true})
		})
	case DataFailed:
		ui.Row(c).Padding(u*2, u*3).Gap(u * 3).AlignItems(ui.Center).Shrink(0).Children(func() {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).TextColor(k.Danger).Grow(1).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(16, 16)
				ui.Text(c, DataError(c, errText)).Role(ui.RoleStatus)
			})
			retry = input.Button(c, core.Msg(c, "list.retry"), input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("refresh-cw")}).Clicked()
		})
	}
	return retry
}

func DataError(c *ui.Context, s string) string {
	if s == "" {
		return core.Msg(c, "list.failed")
	}
	return s
}
