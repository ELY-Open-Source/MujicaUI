package files

import (
	"log/slog"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/code"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// FileEntry identifies a file or folder by slash-separated Path.
type FileEntry struct {
	Path, Name string
	Dir        bool
	Size       int64
	Modified   time.Time
}

// FileExplorerOptions configure a FileExplorer.
type FileExplorerOptions struct {
	// Roots are the top folders of the tree.
	Roots []string
	// Folders returns subfolders; required, nil means a leaf.
	Folders func(dir string) []string
	// Entries returns a folder's files and folders; required.
	Entries  func(dir string) []FileEntry
	Disabled bool
	// Revision changes when entries are edited in place.
	Revision uint64
}

// FileExplorerState keeps navigation, selection and scroll position.
type FileExplorerState struct {
	// Current defaults to the first root.
	Current       string
	tree          data.TreeState[string]
	list          data.ListState[string]
	back, forward []string
	data          *FileEntry
	n             int
	revision      uint64
}

// Selected returns the chosen entries' paths.
func (s *FileExplorerState) Selected() []string {
	var out []string
	for p := range s.list.Selection.All() {
		out = append(out, p)
	}
	return out
}

// Navigate records the current folder and shows dir.
func (s *FileExplorerState) Navigate(dir string) {
	if dir == s.Current {
		return
	}
	if s.Current != "" {
		s.back = append(s.back, s.Current)
	}
	s.Current, s.forward = dir, nil
	s.list.Selection.Clear()
	slog.Debug("mujica: file_explorer navigate", "depth", strings.Count(dir, "/"))
}

func (s *FileExplorerState) step(from, to *[]string) {
	n := len(*from)
	if n == 0 {
		return
	}
	*to = append(*to, s.Current)
	s.Current, *from = (*from)[n-1], (*from)[:n-1]
	s.list.Selection.Clear()
}

// FileExplorerResult is a FileExplorer and what happened in this frame.
type FileExplorerResult struct {
	Element ui.Element
	// Navigated reports a new current folder; Changed a new selection.
	Navigated, Changed bool
	// Opened is a file the user opened.
	Opened *FileEntry
}

// FileExplorer combines folder navigation, history and file selection.
func FileExplorer(c *ui.Context, s *FileExplorerState, opts FileExplorerOptions) FileExplorerResult {
	if s == nil || opts.Folders == nil || opts.Entries == nil || len(opts.Roots) == 0 {
		panic("mujica: FileExplorer needs a state, Roots, Folders and Entries")
	}
	if s.Current == "" {
		s.Current = opts.Roots[0]
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var res FileExplorerResult
	res.Element = ui.Column(c).MinHeight(0).Gap(u * 2).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	if disabled {
		copy := *s
		copy.list.Selection = filesSelectionCopy(s.list.Selection)
		copy.tree.Selection = filesSelectionCopy(s.tree.Selection)
		copy.tree.Tree.Outline.Open = filesSelectionCopy(s.tree.Tree.Outline.Open)
		s = &copy
	}
	at := s.Current
	move := func(fn func()) {
		if !disabled {
			fn()
		}
	}
	mac := runtime.GOOS == "darwin"
	switch {
	case mac && res.Element.Shortcut(ui.Super, ui.KeyBracketLeft), !mac && res.Element.Shortcut(ui.Alt, ui.KeyLeft):
		move(func() { s.step(&s.back, &s.forward) })
	case mac && res.Element.Shortcut(ui.Super, ui.KeyBracketRight), !mac && res.Element.Shortcut(ui.Alt, ui.KeyRight):
		move(func() { s.step(&s.forward, &s.back) })
	case mac && res.Element.Shortcut(ui.Super, ui.KeyUp), !mac && res.Element.Shortcut(ui.Alt, ui.KeyUp):
		move(func() { s.Navigate(fileExplorerParent(s.Current, opts.Roots)) })
	}
	res.Element.Children(func() {
		ui.Row(c).Gap(u).AlignItems(ui.Center).Shrink(0).Children(func() {
			for _, b := range []struct {
				icon, key string
				on        bool
				do        func()
			}{
				{"chevron-left", "file_explorer.back", len(s.back) > 0, func() { s.step(&s.back, &s.forward) }},
				{"chevron-right", "file_explorer.forward", len(s.forward) > 0, func() { s.step(&s.forward, &s.back) }},
				{"arrow-up", "file_explorer.up", fileExplorerParent(s.Current, opts.Roots) != s.Current, func() { s.Navigate(fileExplorerParent(s.Current, opts.Roots)) }},
			} {
				a := core.IconAction(c, icons.Must(b.icon), core.Msg(c, b.key)).Disabled(!b.on)
				if a.Clicked() && b.on {
					move(b.do)
				}
			}
			if r := PathBar(c, s.Current, PathBarOptions{Disabled: opts.Disabled}); r.Navigated != "" {
				move(func() { s.Navigate(r.Navigated) })
			}
		})
		ui.Row(c).Grow(1).MinHeight(core.Density(c).ControlHeight() * 4).Gap(u * 2).AlignItems(ui.Stretch).Children(func() {
			fileExplorerTree(c, s, opts, move)
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				entries := opts.Entries(s.Current)
				var first *FileEntry
				if len(entries) > 0 {
					first = &entries[0]
				}
				if s.data != first || s.n != len(entries) || s.revision != opts.Revision {
					valid := make(map[string]bool, len(entries))
					for _, entry := range entries {
						if entry.Path == "" || valid[entry.Path] {
							panic("mujica: FileExplorer needs unique entry paths")
						}
						valid[entry.Path] = true
					}
					for key := range s.list.Selection.All() {
						if !valid[key] {
							s.list.Selection.Remove(key)
						}
					}
					s.data, s.n, s.revision = first, len(entries), opts.Revision
				}
				v := data.List(c, &s.list, len(entries), data.ListOptions[string]{Multiple: true,
					Key: func(i int) string { return entries[i].Path }, Label: func(i int) string { return entries[i].Name },
					Empty: core.Msg(c, "file_explorer.empty")}, func(i int) { fileExplorerRow(c, k, entries[i]) })
				v.Element.Grow(1)
				res.Changed = v.Changed() && !disabled
				if key, ok := s.list.Selected(); ok && v.Submitted() && !disabled {
					for _, e := range entries {
						if e.Path != key {
							continue
						}
						if e.Dir {
							s.Navigate(e.Path)
						} else {
							open := e
							res.Opened = &open
							slog.Debug("mujica: file_explorer open")
						}
					}
				}
			})
		})
	})
	res.Navigated = s.Current != at
	return res
}

func filesSelectionCopy(src ui.Selection[string]) ui.Selection[string] {
	var dst ui.Selection[string]
	for key := range src.All() {
		dst.Add(key)
	}
	return dst
}

// fileExplorerTree shows the folders, choosing the current one.
func fileExplorerTree(c *ui.Context, s *FileExplorerState, opts FileExplorerOptions, move func(func())) {
	for _, a := range fileExplorerAncestors(s.Current, opts.Roots) {
		s.tree.SetOpen(a, true)
	}
	s.tree.Tree.Sel, s.tree.Tree.Has = s.Current, true
	v := data.Tree(c, &s.tree, data.TreeOptions[string]{Roots: opts.Roots, Children: opts.Folders,
		Label: func(p string) string { return path.Base(p) }})
	v.Element.Width(220).Shrink(0)
	if v.Changed() {
		if dir, ok := s.tree.Selected(); ok && dir != s.Current {
			move(func() { s.Navigate(dir) })
		}
	}
}

// fileExplorerAncestors lists the folders above dir, inside its root.
func fileExplorerAncestors(dir string, roots []string) []string {
	var out []string
	for p := dir; ; {
		parent := fileExplorerParent(p, roots)
		if parent == p {
			return out
		}
		out = append(out, parent)
		p = parent
	}
}

// fileExplorerParent is dir's parent, or dir itself at a root.
func fileExplorerParent(dir string, roots []string) string {
	for _, r := range roots {
		if dir == r {
			return dir
		}
	}
	p := path.Dir(dir)
	if p == "." || p == dir {
		return dir
	}
	return p
}

func fileExplorerRow(c *ui.Context, k theme.Tokens, e FileEntry) {
	FileIcon(c, e.Name, FileIconOptions{Folder: e.Dir})
	ui.Text(c, e.Name).SingleLine().Grow(1).MinWidth(0)
	if !e.Dir {
		ui.Text(c, code.ProcessListBytes(uint64(max(e.Size, 0)))).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).Shrink(0)
	}
	if !e.Modified.IsZero() {
		ui.Text(c, e.Modified.Format("2006-01-02 15:04")).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).Shrink(0)
	}
}
