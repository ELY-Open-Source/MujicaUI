package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// FileGridOptions configure a FileGrid.
type FileGridOptions struct {
	// Revision changes when entries are edited in place.
	Revision uint64
	// Thumbnail builds a preview; nil or false shows FileIcon.
	Thumbnail func(e FileEntry) bool
	// ItemWidth defaults to 128 DIP, minimum 28.
	ItemWidth float32
	Disabled  bool
}

// FileGridResult is a FileGrid and what happened in this frame.
type FileGridResult struct {
	Element *ui.Element
	// Changed reports a new selection; Opened is the entry opened.
	Changed bool
	Opened  *FileEntry
}

// FileGridState is a FileGrid's selection and place.
type FileGridState struct {
	list     ListState[string]
	data     *FileEntry
	n        int
	revision uint64
	paths    []string
}

// Selected returns the chosen entries' paths.
func (s *FileGridState) Selected() []string {
	var out []string
	for p := range s.list.Selection.All() {
		out = append(out, p)
	}
	return out
}

// FileGrid virtualizes thumbnails, file selection and opening.
func FileGrid(c *ui.Context, entries []FileEntry, s *FileGridState, opts FileGridOptions) FileGridResult {
	if s == nil {
		panic("mujica: FileGrid needs a state")
	}
	if opts.ItemWidth < 0 {
		panic("mujica: FileGrid ItemWidth is negative")
	}
	w := opts.ItemWidth
	if w == 0 {
		w = 128
	}
	k := Tokens(c)
	u := Density(c).Unit()
	if w < theme.MinHitSize {
		panic("mujica: FileGrid ItemWidth is below the minimum hit size")
	}
	thumb := max(0, w-4*u)
	var res FileGridResult
	res.Element = ui.Column(c).MinHeight(0).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	if disabled {
		copy := *s
		copy.list.Selection = filesSelectionCopy(s.list.Selection)
		s = &copy
	}
	var first *FileEntry
	if len(entries) > 0 {
		first = &entries[0]
	}
	if s.data != first || s.n != len(entries) || s.revision != opts.Revision {
		valid := make(map[string]bool, len(entries))
		moved := len(s.paths) != len(entries)
		for i, entry := range entries {
			if entry.Path == "" || valid[entry.Path] {
				panic("mujica: FileGrid needs unique entry paths")
			}
			valid[entry.Path] = true
			if i >= len(s.paths) || s.paths[i] != entry.Path {
				moved = true
			}
		}
		for key := range s.list.Selection.All() {
			if !valid[key] {
				s.list.Selection.Remove(key)
			}
		}
		if moved {
			s.list.grid = ui.GridState{}
			s.paths = make([]string, len(entries))
			for i, entry := range entries {
				s.paths[i] = entry.Path
			}
		}
		s.data, s.n, s.revision = first, len(entries), opts.Revision
	}
	res.Element.Children(func() {
		v := List(c, &s.list, len(entries), ListOptions[string]{Multiple: true, Grid: true, ItemWidth: w, ItemHeight: thumb*0.75 + 56,
			Key: func(i int) string { return entries[i].Path }, Label: func(i int) string { return entries[i].Name },
			Empty: msg(c, "file_grid.empty")}, func(i int) {
			e := entries[i]
			ui.Box(c).Height(thumb * 0.75).Shrink(0).Center().Clip().Radius(theme.ControlRadius).Background(k.Background).Children(func() {
				if opts.Thumbnail == nil || !opts.Thumbnail(e) {
					FileIcon(c, e.Name, FileIconOptions{Folder: e.Dir, Size: 40})
				}
			})
			ui.Text(c, e.Name).SingleLine().TextAlign(ui.Center).FontSize(FontSize(c, theme.CaptionSize))
		})
		v.Element.Grow(1)
		res.Changed = v.Changed() && !disabled
		if key, ok := s.list.Selected(); ok && v.Submitted() && !res.Element.IsDisabled() {
			for _, e := range entries {
				if e.Path == key {
					open := e
					res.Opened = &open
					slog.Debug("mujica: file_grid open", "dir", e.Dir)
				}
			}
		}
	})
	return res
}
