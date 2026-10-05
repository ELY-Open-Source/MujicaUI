package main

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
	"mujicaui/theme"
)

const filesID = "c-files"

type fileEntry struct {
	ID       int
	Name     string
	Size     int64
	Modified time.Time
}

type filesDemo struct {
	entries  []fileEntry
	next     int
	version  int
	query    string
	table    mujica.DataTableState[int]
	dropped  []mujica.DroppedFile
	picked   []string
	renaming bool
	renameID int
	newName  string
	deleting bool
	doomed   []int
}

func filesSample() []fileEntry {
	at := func(m time.Month, d, h int) time.Time { return time.Date(2026, m, d, h, 0, 0, 0, time.UTC) }
	return []fileEntry{
		{1, "setlist-masquerade.pdf", 182_400, at(9, 28, 21)},
		{2, "rehearsal-0412.wav", 48_211_968, at(4, 12, 19)},
		{3, "stage-plot.png", 2_406_400, at(8, 3, 14)},
		{4, "lyrics-draft.md", 6_912, at(10, 1, 23)},
		{5, "venue-contract.docx", 96_256, at(7, 19, 10)},
		{6, "overture-score.musicxml", 731_136, at(9, 30, 2)},
		{7, "press-kit.zip", 18_874_368, at(6, 6, 16)},
		{8, "budget-2026.xlsx", 41_984, at(5, 22, 9)},
	}
}

func newFilesDemo() filesDemo { return filesDemo{entries: filesSample(), next: 9} }

// fileSize writes a byte count the way a file manager does.
func fileSize(n int64) string {
	switch {
	case n < 1<<10:
		return fmt.Sprintf("%d B", n)
	case n < 1<<20:
		return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	}
}

func (g *gallery) fileColumns(c *ui.Context) []mujica.DataColumn[fileEntry] {
	return []mujica.DataColumn[fileEntry]{
		{ID: "name", Title: g.t("Name", "名称"), MinWidth: 200, Sortable: true,
			Text:    func(f fileEntry) string { return f.Name },
			Compare: func(a, b fileEntry) int { return cmp.Compare(a.Name, b.Name) },
			Cell: func(f fileEntry) {
				ui.Row(c).Gap(8).AlignItems(ui.Center).MinWidth(0).Children(func() {
					ui.Icon(c, icons.Must("file")).Size(16, 16).TextColor(mujica.Tokens(c).TextMuted)
					ui.Text(c, f.Name).SingleLine()
				})
			}},
		{ID: "size", Title: g.t("Size", "大小"), Width: 110, Align: ui.End, Sortable: true,
			Text:    func(f fileEntry) string { return fileSize(f.Size) },
			Compare: func(a, b fileEntry) int { return cmp.Compare(a.Size, b.Size) }},
		{ID: "modified", Title: g.t("Modified", "修改时间"), Width: 170, Sortable: true,
			Text:    func(f fileEntry) string { return f.Modified.Format("2006-01-02 15:04") },
			Compare: func(a, b fileEntry) int { return a.Modified.Compare(b.Modified) }},
	}
}

// targets are the checked rows, or else the focused row.
func (s *filesDemo) targets() []int {
	if s.table.Selection.Len() > 0 {
		return slices.Sorted(s.table.Selection.All())
	}
	if id, ok := s.table.Selected(); ok {
		return []int{id}
	}
	return nil
}

func (s *filesDemo) entry(id int) *fileEntry {
	i := slices.IndexFunc(s.entries, func(f fileEntry) bool { return f.ID == id })
	if i < 0 {
		panic(fmt.Sprintf("gallery: no file %d", id))
	}
	return &s.entries[i]
}

// add appends a file unless one with its name is listed; it reports whether it did.
func (s *filesDemo) add(path string, size int64, at time.Time) bool {
	name := filepath.Base(path)
	if slices.ContainsFunc(s.entries, func(f fileEntry) bool { return f.Name == name }) {
		return false
	}
	s.entries = append(s.entries, fileEntry{s.next, name, size, at})
	s.next++
	s.version++
	return true
}

func (g *gallery) filesView(c *ui.Context, log func(string)) {
	s := &g.cFiles
	k := mujica.Tokens(c)
	cols := g.fileColumns(c)
	targets := s.targets()
	ui.Row(c).FillWidth().Gap(12).AlignItems(ui.Center).Children(func() {
		mujica.Toolbar(c, mujica.ToolbarOptions{Label: g.t("File actions", "文件操作")}, func() {
			g.fileActions(c, log, targets, func(m *mujica.PopupMenu) { g.fileMenu(m, log, targets) })
		}).Grow(1)
		before := s.query
		mujica.SearchInput(c, &s.query, mujica.SearchInputOptions{Placeholder: g.t("Filter files", "筛选文件")}).Element.Width(240)
		if s.query != before {
			log(fmt.Sprintf("filter %q", s.query))
		}
	})
	s.table.Filter = s.query
	shown := len(mujica.DataTableRows(s.entries, cols, s.table.Sort, s.query))
	switch {
	case len(s.entries) == 0:
		if mujica.EmptyState(c, mujica.EmptyStateOptions{Kind: mujica.EmptyNoData, Title: g.t("No files", "没有文件"),
			Description: g.t("Drop files below or choose some to add them.", "将文件拖到下方，或选取文件以添加。"),
			Action:      g.t("Restore samples", "恢复示例")}).Action {
			g.restoreFiles(log)
		}
	case shown == 0:
		if mujica.EmptyState(c, mujica.EmptyStateOptions{Kind: mujica.EmptyNoResults, Query: s.query, Action: g.t("Clear filter", "清除筛选")}).Action {
			s.query = ""
			log("filter cleared")
		}
	default:
		v := mujica.DataTable(c, &s.table, s.entries, mujica.DataTableOptions[fileEntry, int]{
			Columns: cols, Key: func(f fileEntry) int { return f.ID }, Multiple: true, Version: s.version,
		})
		v.Element.FillWidth().Height(430)
		mujica.ContextMenu(c, v.Element, mujica.ContextMenuOptions{}, func(m *mujica.PopupMenu) { g.fileMenu(m, log, targets) })
		if v.Changed() {
			log(fmt.Sprintf("selection %d checked", s.table.Selection.Len()))
		}
		if v.SortChanged() {
			log(fmt.Sprintf("sort %s descending=%v", s.table.Sort.Column, s.table.Sort.Descending))
		}
		if v.Submitted() && len(targets) == 1 {
			g.startRename(log, targets[0])
		}
	}
	var size int64
	for _, f := range s.entries {
		size += f.Size
	}
	mujica.StatusBar(c, mujica.StatusBarOptions{
		Left: []mujica.StatusItem{
			{ID: "count", Icon: icons.Must("folder"), Text: fmt.Sprintf(g.t("%d files · %s", "%d 个文件 · %s"), len(s.entries), fileSize(size))},
			{ID: "shown", Text: fmt.Sprintf(g.t("%d shown", "显示 %d 个"), shown), Secondary: true},
		},
		Right: []mujica.StatusItem{{ID: "selected", Text: fmt.Sprintf(g.t("%d selected", "已选 %d 个"), s.table.Selection.Len())}},
	}).Element.Radius(theme.CardRadius).Border(theme.BorderWidth, k.Border)
	g.fileAdders(c, log)
	g.renameDialog(c, log)
	g.deleteDialog(c, log)
}

// fileActions fills the toolbar; menu fills the More dropdown, shared with the context menu.
func (g *gallery) fileActions(c *ui.Context, log func(string), targets []int, menu func(m *mujica.PopupMenu)) {
	mujica.ToolbarGroup(c, func() {
		if mujica.ToolbarButton(c, g.t("Rename", "重命名"), mujica.ToolbarButtonOptions{Icon: icons.Must("file-text"), ShowLabel: true, Disabled: len(targets) != 1}).Clicked() {
			g.startRename(log, targets[0])
		}
		if mujica.ToolbarButton(c, g.t("Delete", "删除"), mujica.ToolbarButtonOptions{Icon: icons.Must("trash-2"), ShowLabel: true, Disabled: len(targets) == 0}).Clicked() {
			g.startDelete(log, targets)
		}
	})
	mujica.ToolbarSeparator(c)
	mujica.DropdownMenu(c, g.t("Actions", "操作"), mujica.DropdownMenuOptions{Icon: icons.Must("ellipsis")}, menu)
}

// fileMenu holds every file action; the dropdown and the context menu share it.
func (g *gallery) fileMenu(m *mujica.PopupMenu, log func(string), targets []int) {
	s := &g.cFiles
	if m.Item(g.t("Rename…", "重命名…"), mujica.PopupMenuItemOptions{Disabled: len(targets) != 1}) {
		g.startRename(log, targets[0])
	}
	if m.Item(g.t("Delete…", "删除…"), mujica.PopupMenuItemOptions{Disabled: len(targets) == 0}) {
		g.startDelete(log, targets)
	}
	if m.Item(g.t("Select all", "全选"), mujica.PopupMenuItemOptions{Disabled: len(s.entries) == 0}) {
		for _, f := range s.entries {
			s.table.Selection.Add(f.ID)
		}
		log(fmt.Sprintf("selected all %d", len(s.entries)))
	}
	if m.Item(g.t("Clear selection", "清除选择"), mujica.PopupMenuItemOptions{Disabled: s.table.Selection.Len() == 0}) {
		s.table.Selection.Clear()
		log("selection cleared")
	}
	if m.Item(g.t("Restore samples", "恢复示例"), mujica.PopupMenuItemOptions{}) {
		g.restoreFiles(log)
	}
}

func (g *gallery) restoreFiles(log func(string)) {
	s := &g.cFiles
	s.entries, s.next, s.query = filesSample(), 9, ""
	s.table.Selection.Clear()
	s.version++
	log("samples restored")
}

// fileAdders are the drop zone and the picker; both only list files, nothing is uploaded.
func (g *gallery) fileAdders(c *ui.Context, log func(string)) {
	s := &g.cFiles
	mujica.Card(c, mujica.CardOptions{Title: g.t("Add files", "添加文件")}, func() {
		ui.Text(c, g.t("Files are checked and listed, never read or uploaded.", "文件只做校验并列出，不会读取或上传。")).TextColor(mujica.Tokens(c).TextMuted)
		if mujica.FileDropZone(c, &s.dropped, mujica.FileDropZoneOptions{MaxFiles: 8, MaxSize: 50 << 20,
			Hint: g.t("Drop files here to list them", "将文件拖到此处以列出")}).Changed() {
			for _, f := range s.dropped {
				g.addFile(c, log, f.Path, f.Size, time.Now())
			}
			s.dropped = nil
		}
		p := mujica.FilePicker(c, &s.picked, mujica.FilePickerOptions{Multiple: true, Label: g.t("Choose files…", "选取文件…")})
		if p.Canceled() {
			log("picker canceled")
		}
		if p.Changed() {
			for _, path := range s.picked {
				info, err := os.Stat(path)
				if err != nil {
					mujica.ShowToast(c, mujica.ToastMessage{Severity: mujica.SeverityError, Message: err.Error()})
					log("stat failed: " + err.Error())
					continue
				}
				g.addFile(c, log, path, info.Size(), info.ModTime())
			}
			s.picked = nil
		}
	})
}

func (g *gallery) addFile(c *ui.Context, log func(string), path string, size int64, at time.Time) {
	if g.cFiles.add(path, size, at) {
		log("added " + filepath.Base(path))
		return
	}
	log("skipped duplicate " + filepath.Base(path))
	mujica.ShowToast(c, mujica.ToastMessage{ID: "dup-" + path, Severity: mujica.SeverityWarning,
		Message: fmt.Sprintf(g.t("%s is already listed", "%s 已在列表中"), filepath.Base(path))})
}

func (g *gallery) startRename(log func(string), id int) {
	s := &g.cFiles
	s.renaming, s.renameID, s.newName = true, id, s.entry(id).Name
	log("rename " + s.newName)
}

// renameError checks a new name against the others in the list.
func (g *gallery) renameError() string {
	s := &g.cFiles
	name := strings.TrimSpace(s.newName)
	switch {
	case name == "":
		return g.t("Enter a name.", "请输入名称。")
	case strings.ContainsAny(name, `/\`):
		return g.t("A name cannot contain / or \\.", "名称不能包含 / 或 \\。")
	case slices.ContainsFunc(s.entries, func(f fileEntry) bool { return f.Name == name && f.ID != s.renameID }):
		return g.t("Another file has this name.", "已有同名文件。")
	}
	return ""
}

func (g *gallery) renameDialog(c *ui.Context, log func(string)) {
	s := &g.cFiles
	if !s.renaming {
		return
	}
	bad := g.renameError()
	commit := func() {
		if bad != "" {
			return
		}
		f := s.entry(s.renameID)
		log(fmt.Sprintf("renamed %s to %s", f.Name, strings.TrimSpace(s.newName)))
		f.Name = strings.TrimSpace(s.newName)
		s.version++
		s.renaming = false
	}
	mujica.Dialog(c, &s.renaming, mujica.DialogOptions{Title: g.t("Rename file", "重命名文件"), Actions: func() {
		if mujica.Button(c, g.t("Cancel", "取消"), mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.renaming = false
			log("rename canceled")
		}
		if mujica.Button(c, g.t("Rename", "重命名"), mujica.ButtonOptions{Disabled: bad != ""}).Clicked() {
			commit()
		}
	}}, func() {
		mujica.FormField(c, g.t("New name", "新名称"), mujica.FormFieldOptions{Error: bad}, func() *ui.Element {
			r := mujica.TextInput(c, &s.newName, mujica.TextInputOptions{Error: bad})
			if r.Submitted() {
				commit()
			}
			return r.Input.AutoFocus()
		})
	})
}

func (g *gallery) startDelete(log func(string), ids []int) {
	s := &g.cFiles
	s.deleting, s.doomed = true, ids
	log(fmt.Sprintf("delete asked for %d", len(ids)))
}

func (g *gallery) deleteDialog(c *ui.Context, log func(string)) {
	s := &g.cFiles
	if !s.deleting {
		return
	}
	desc := fmt.Sprintf(g.t("%d files will be removed from the list.", "将从列表中移除 %d 个文件。"), len(s.doomed))
	if len(s.doomed) == 1 {
		desc = fmt.Sprintf(g.t("“%s” will be removed from the list.", "将从列表中移除“%s”。"), s.entry(s.doomed[0]).Name)
	}
	r := mujica.AlertDialog(c, &s.deleting, mujica.AlertDialogOptions{Title: g.t("Delete files?", "删除文件？"), Description: desc,
		ConfirmLabel: g.t("Delete permanently", "永久删除")})
	switch {
	case r.Confirmed:
		s.entries = slices.DeleteFunc(slices.Clone(s.entries), func(f fileEntry) bool { return slices.Contains(s.doomed, f.ID) })
		for _, id := range s.doomed {
			s.table.Selection.Remove(id)
		}
		s.version++
		log(fmt.Sprintf("deleted %d", len(s.doomed)))
		mujica.ShowToast(c, mujica.ToastMessage{ID: "files-deleted", Severity: mujica.SeveritySuccess,
			Message: fmt.Sprintf(g.t("Deleted %d files", "已删除 %d 个文件"), len(s.doomed))})
		s.doomed = nil
	case r.Canceled:
		s.doomed = nil
		log("delete canceled")
	}
}
