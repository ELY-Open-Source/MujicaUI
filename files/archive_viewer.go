package files

import (
	"cmp"
	"fmt"
	"log/slog"
	"maps"
	"math"
	"path"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/code"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
)

// ArchiveEntry describes an archive member with a slash-separated Path.
type ArchiveEntry struct {
	Path string
	Dir  bool
	// Size is the unpacked size; Packed the compressed size.
	Size, Packed int64
}

// ArchiveViewerState preserves folders, selection, sort and indexed entries.
type ArchiveViewerState struct {
	table    data.TreeTableState[string]
	index    archiveIndex
	data     *ArchiveEntry
	n        int
	revision uint64
}

// archiveIndex is the tree built from the entries.
type archiveIndex struct {
	byPath map[string]*ArchiveEntry
	kids   map[string][]string
	roots  []string
	// sums totals a folder's files.
	sums       map[string][2]int64
	files      int
	size, pack int64
}

// ArchiveViewerOptions configure an ArchiveViewer.
type ArchiveViewerOptions struct {
	// Revision changes when entries are edited in place.
	Revision uint64
	// Status shows the caller's listing in progress or failed.
	Status data.DataStatus
	Error  string
}

// ArchiveViewerResult is the viewer and the entry opened this frame.
type ArchiveViewerResult struct {
	Element ui.Element
	Opened  *ArchiveEntry
	Retried bool
}

// ArchiveViewer virtualizes archive entries with sizes and compression ratios.
func ArchiveViewer(c *ui.Context, entries []ArchiveEntry, s *ArchiveViewerState, opts ArchiveViewerOptions) ArchiveViewerResult {
	if s == nil {
		panic("mujica: ArchiveViewer needs a state")
	}
	if s.revision != opts.Revision {
		s.data = nil
		s.index.byPath = nil
		s.revision = opts.Revision
	}
	archiveBuild(entries, s)
	ix := &s.index
	u := core.Density(c).Unit()
	size := func(p string) (int64, int64) {
		if e := ix.byPath[p]; e != nil && !e.Dir {
			return e.Size, e.Packed
		}
		t := ix.sums[p]
		return t[0], t[1]
	}
	cell := func(text string) { code.CodeMono(c, text).TextColor(core.Tokens(c).TextMuted) }
	var res ArchiveViewerResult
	res.Element = ui.Column(c).MinHeight(0).Gap(u * 2)
	if res.Element.IsDisabled() {
		copy := *s
		copy.table.Selection = filesSelectionCopy(s.table.Selection)
		copy.table.Tree.Outline.Open = filesSelectionCopy(s.table.Tree.Outline.Open)
		copy.table.Layout = ui.TableLayout{Order: slices.Clone(s.table.Layout.Order), Widths: maps.Clone(s.table.Layout.Widths)}
		copy.table.Tree.Outline.List.Columns = copy.table.Layout
		s = &copy
	}
	res.Element.Children(func() {
		code.CodeCaption(c, fmt.Sprintf(core.Msg(c, "archive_viewer.summary"), ix.files, code.ProcessListBytes(uint64(ix.size)),
			code.ProcessListBytes(uint64(ix.pack)), archiveRatio(ix.size, ix.pack)))
		v := data.TreeTable(c, &s.table, data.TreeTableOptions[string]{
			Roots: ix.roots, Children: func(p string) []string { return ix.kids[p] },
			Label: func(p string) string { return path.Base(p) }, Status: opts.Status, Error: opts.Error,
			Empty: core.Msg(c, "archive_viewer.empty"),
			Columns: []data.TreeColumn[string]{
				{ID: "name", Title: core.Msg(c, "archive_viewer.name"), Compare: strings.Compare, Cell: func(p string) {
					ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
						FileIcon(c, p, FileIconOptions{Folder: ix.kids[p] != nil})
						ui.Text(c, path.Base(p)).SingleLine()
					})
				}},
				{ID: "size", Title: core.Msg(c, "archive_viewer.size"), Width: 96, Align: ui.End,
					Compare: func(a, b string) int { x, _ := size(a); y, _ := size(b); return cmp.Compare(x, y) },
					Cell:    func(p string) { n, _ := size(p); cell(code.ProcessListBytes(uint64(n))) }},
				{ID: "packed", Title: core.Msg(c, "archive_viewer.packed"), Width: 96, Align: ui.End,
					Cell: func(p string) { _, n := size(p); cell(code.ProcessListBytes(uint64(n))) }},
				{ID: "ratio", Title: core.Msg(c, "archive_viewer.ratio"), Width: 72, Align: ui.End,
					Cell: func(p string) { a, b := size(p); cell(archiveRatio(a, b)) }},
			},
		})
		v.Element.Grow(1)
		res.Retried = v.Retried()
		if p, ok := s.table.Selected(); ok && v.Submitted() && !res.Element.IsDisabled() {
			if e := ix.byPath[p]; e != nil && !e.Dir {
				open := *e
				res.Opened = &open
				slog.Debug("mujica: archive_viewer open")
			}
		}
	})
	return res
}

// archiveRatio formats the saved-size percentage.
func archiveRatio(size, packed int64) string {
	if size <= 0 {
		return "—"
	}
	return fmt.Sprintf("%.0f%%", 100*(1-float64(packed)/float64(size)))
}

// archiveBuild caches entries and implied parent directories.
func archiveBuild(entries []ArchiveEntry, s *ArchiveViewerState) {
	var first *ArchiveEntry
	if len(entries) > 0 {
		first = &entries[0]
	}
	if s.data == first && s.n == len(entries) && s.index.byPath != nil {
		return
	}
	s.data, s.n = first, len(entries)
	ix := archiveIndex{byPath: map[string]*ArchiveEntry{}, kids: map[string][]string{}, sums: map[string][2]int64{}}
	var addDir func(p string)
	link := func(p string) {
		parent := path.Dir(p)
		if parent == "." {
			ix.roots = append(ix.roots, p)
			return
		}
		addDir(parent)
		ix.kids[parent] = append(ix.kids[parent], p)
	}
	addDir = func(p string) {
		if e := ix.byPath[p]; e != nil && !e.Dir {
			panic("mujica: ArchiveEntry file used as a folder: " + p)
		}
		if _, ok := ix.kids[p]; ok {
			return
		}
		ix.kids[p] = []string{}
		link(p)
	}
	for i := range entries {
		e := &entries[i]
		p := strings.Trim(e.Path, "/")
		if p == "" || e.Size < 0 || e.Packed < 0 {
			panic("mujica: ArchiveEntry needs a Path and sizes not negative")
		}
		if _, dup := ix.byPath[p]; dup {
			panic("mujica: duplicate ArchiveEntry " + p)
		}
		if !e.Dir && ix.kids[p] != nil {
			panic("mujica: ArchiveEntry folder used as a file: " + p)
		}
		ix.byPath[p] = e
		if e.Dir {
			addDir(p)
			continue
		}
		link(p)
		ix.files++
		if e.Size > math.MaxInt64-ix.size || e.Packed > math.MaxInt64-ix.pack {
			panic("mujica: ArchiveEntry totals overflow")
		}
		ix.size += e.Size
		ix.pack += e.Packed
		for d := path.Dir(p); d != "."; d = path.Dir(d) {
			t := ix.sums[d]
			ix.sums[d] = [2]int64{t[0] + e.Size, t[1] + e.Packed}
		}
	}
	s.index = ix
	slog.Debug("mujica: archive_viewer index", "files", ix.files)
}
