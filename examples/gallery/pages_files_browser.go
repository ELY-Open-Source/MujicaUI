package main

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
)

var filesFolders = map[string][]string{
	"/court":            {"chronicles/", "decrees/", "seal.go", "notes.md"},
	"/court/chronicles": {"first.md", "second.md"},
	"/court/decrees":    {"42.pdf", "crest.png"},
}

func filesEntries(dir string) []mujica.FileEntry {
	var entries []mujica.FileEntry
	for _, n := range filesFolders[dir] {
		name := strings.TrimSuffix(n, "/")
		entries = append(entries, mujica.FileEntry{Path: path.Join(dir, name), Name: name, Dir: name != n, Size: 2048})
	}
	return entries
}

func fileExplorerDemo(c *ui.Context, log func(string)) {
	type state struct {
		explorer        mujica.FileExplorerState
		disabled, empty bool
	}
	s := ui.Local(c.Root(), "381", func() state { return state{} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	mediaCheck(c, log, &s.empty, "Empty folder")
	r := mujica.FileExplorer(c, &s.explorer, mujica.FileExplorerOptions{
		Roots: []string{"/court"}, Disabled: s.disabled,
		Folders: func(dir string) []string {
			var dirs []string
			for _, e := range filesEntries(dir) {
				if e.Dir {
					dirs = append(dirs, e.Path)
				}
			}
			return dirs
		},
		Entries: func(dir string) []mujica.FileEntry {
			if s.empty {
				return nil
			}
			return filesEntries(dir)
		},
	})
	r.Element.Height(320).FillWidth()
	if r.Navigated {
		log("Navigate: " + s.explorer.Current)
	}
	if r.Changed {
		log(fmt.Sprintf("Selected: %v", s.explorer.Selected()))
	}
	if r.Opened != nil {
		log("Open: " + r.Opened.Path)
	}
}

func fileGridDemo(c *ui.Context, log func(string)) {
	type state struct {
		grid                    mujica.FileGridState
		entries                 []mujica.FileEntry
		disabled, thumbs, empty bool
	}
	s := ui.Local(c.Root(), "382", func() state {
		return state{thumbs: true, entries: append(filesEntries("/court"), filesEntries("/court/decrees")...)}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	mediaCheck(c, log, &s.thumbs, "Thumbnails")
	mediaCheck(c, log, &s.empty, "Empty")
	entries := s.entries
	if s.empty {
		entries = nil
	}
	r := mujica.FileGrid(c, entries, &s.grid, mujica.FileGridOptions{Disabled: s.disabled, Thumbnail: func(e mujica.FileEntry) bool {
		if !s.thumbs || mujica.FileKindOf(e.Name) != mujica.FileImage {
			return false
		}
		mujica.Image(c, mediaBitmaps[0], mujica.ImageOptions{Alt: e.Name}).Element.Fill()
		return true
	}})
	r.Element.Height(320).FillWidth()
	if r.Changed {
		log(fmt.Sprintf("Selected: %v", s.grid.Selected()))
	}
	if r.Opened != nil {
		log("Open: " + r.Opened.Path)
	}
}

func fileIconDemo(c *ui.Context, log func(string)) {
	type state struct {
		name         string
		folder, open bool
		size         float64
	}
	s := ui.Local(c.Root(), "383", func() state { return state{name: "seal.go", size: 32} })
	if mujica.TextInput(c, &s.name, mujica.TextInputOptions{Label: "File name"}).Changed() {
		log("File name changed")
	}
	mediaCheck(c, log, &s.folder, "Folder")
	if s.folder {
		mediaCheck(c, log, &s.open, "Open folder")
	}
	if mujica.NumberInput(c, &s.size, mujica.NumberInputOptions{Min: 20, Max: 64, Step: 4, Label: "Icon size"}).Changed() {
		log(fmt.Sprintf("Size: %.0f", s.size))
	}
	mujica.FileIcon(c, s.name, mujica.FileIconOptions{Folder: s.folder, Open: s.folder && s.open, Size: float32(s.size)})
	ui.Row(c).Gap(12).Children(func() {
		for _, name := range []string{"seal.go", "notes.md", "crest.png", "deed.pdf", "ledger.zip", "hymn.mp3", "rite.mov", "blob"} {
			mujica.FileIcon(c, name, mujica.FileIconOptions{Size: 32})
		}
	})
}

func filePreviewDemo(c *ui.Context, log func(string)) {
	type state struct {
		kind   int
		status mujica.DataStatus
		closed bool
	}
	s := ui.Local(c.Root(), "384", func() state { return state{} })
	if mujica.Select(c, &s.kind, []mujica.SelectOption[int]{{Value: 0, Label: "Code"}, {Value: 1, Label: "Text"}, {Value: 2, Label: "Image"}, {Value: 3, Label: "Unsupported"}}, mujica.SelectOptions{Label: "Preview type"}).Changed() {
		log("Preview type changed")
		s.closed = false
	}
	filesStatus(c, log, &s.status)
	if s.closed {
		if mujica.Button(c, "Show preview", mujica.ButtonOptions{}).Clicked() {
			s.closed = false
			log("Preview reopened")
		}
		return
	}
	f := []mujica.PreviewFile{{Name: "seal.go", Text: "package court\n\nfunc Seal() {}"}, {Name: "notes.md", Text: "The seal of the court."}, {Name: "crest.png", Image: mediaBitmaps[0]}, {Name: "deed.pdf", Size: 4096}}[s.kind]
	r := mujica.FilePreview(c, f, mujica.FilePreviewOptions{Status: s.status, Error: "Reading failed"})
	r.Element.Height(280).FillWidth()
	if r.Closed {
		s.closed = true
		log("Preview closed")
	}
	if r.Opened {
		log("Open: " + f.Name)
	}
	if r.Retried {
		s.status = mujica.DataReady
		log("Read retried")
	}
}

func filesStatus(c *ui.Context, log func(string), value *mujica.DataStatus) {
	if mujica.Select(c, value, []mujica.SelectOption[mujica.DataStatus]{{Value: mujica.DataReady, Label: "Ready"}, {Value: mujica.DataLoading, Label: "Loading"}, {Value: mujica.DataFailed, Label: "Failed"}}, mujica.SelectOptions{Label: "Load status"}).Changed() {
		log(fmt.Sprintf("Status: %d", *value))
	}
}

func pathBarDemo(c *ui.Context, log func(string)) {
	type state struct {
		path     string
		disabled bool
	}
	s := ui.Local(c.Root(), "385", func() state { return state{path: "/court/chronicles/royal/archives/seals"} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if mujica.TextInput(c, &s.path, mujica.TextInputOptions{Label: "Path"}).Changed() {
		log("Path changed")
	}
	if strings.Trim(s.path, "/") == "" {
		ui.Text(c, "Enter a path")
		return
	}
	if r := mujica.PathBar(c, s.path, mujica.PathBarOptions{Disabled: s.disabled}); r.Navigated != "" {
		s.path = r.Navigated
		log("Navigate: " + r.Navigated)
	}
}

func renameInlineDemo(c *ui.Context, log func(string)) {
	type state struct {
		name                      string
		editing, disabled, folder bool
	}
	s := ui.Local(c.Root(), "386", func() state { return state{name: "seal.txt"} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	mediaCheck(c, log, &s.folder, "Folder")
	if mujica.Button(c, "Rename", mujica.ButtonOptions{Disabled: s.disabled}).Clicked() {
		s.editing = true
		log("Rename started")
	}
	r := mujica.RenameInline(c, &s.name, &s.editing, mujica.RenameInlineOptions{Disabled: s.disabled, Folder: s.folder, Validate: func(name string) string {
		if strings.ContainsAny(name, "/\\") {
			return "Name contains a path separator"
		}
		return ""
	}})
	if r.Renamed {
		log("Renamed: " + s.name)
	}
	if r.Canceled {
		log("Rename canceled")
	}
}

func fileOperationProgressDemo(c *ui.Context, log func(string)) {
	type state struct {
		op                 mujica.FileOperation
		disabled, conflict bool
		canceled           bool
	}
	s := ui.Local(c.Root(), "387", func() state {
		return state{op: mujica.FileOperation{Kind: mujica.FileCopy, Total: 10, Current: "seal.go"}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	mediaCheck(c, log, &s.conflict, "Conflict")
	if mujica.Select(c, &s.op.Kind, []mujica.SelectOption[mujica.FileOperationKind]{{Value: mujica.FileCopy, Label: "Copy"}, {Value: mujica.FileMove, Label: "Move"}, {Value: mujica.FileDelete, Label: "Delete"}}, mujica.SelectOptions{Label: "Operation"}).Changed() {
		log("Operation changed")
	}
	if mujica.Button(c, "Advance", mujica.ButtonOptions{Disabled: s.op.Done == s.op.Total}).Clicked() {
		s.op.Done++
		log(fmt.Sprintf("Done: %d", s.op.Done))
	}
	if mujica.Button(c, "Reset", mujica.ButtonOptions{}).Clicked() {
		s.op.Done = 0
		s.canceled = false
		log("Operation reset")
	}
	op := s.op
	if s.conflict {
		op.Conflict = "seal.go"
	}
	r := mujica.FileOperationProgress(c, op, mujica.FileOperationOptions{Disabled: s.disabled || s.canceled})
	if r.Canceled {
		s.canceled = true
		log("Operation canceled")
	}
	if r.Resolved != mujica.ConflictNone {
		s.conflict = false
		log(fmt.Sprintf("Conflict: %d; apply all: %v", r.Resolved, r.ApplyAll))
	}
}

func filesTransfers() []mujica.QueuedTransfer {
	return []mujica.QueuedTransfer{
		{ID: "a", Name: "ledger.zip", Upload: true, Status: mujica.TransferActive, Done: 30, Total: 120},
		{ID: "b", Name: "crest.png", Status: mujica.TransferPaused, Done: 1, Total: 4},
		{ID: "c", Name: "deed.pdf", Upload: true, Status: mujica.TransferFailed, Total: 2, Error: "Network lost"},
	}
}

func transferQueueDemo(c *ui.Context, log func(string)) {
	type state struct {
		items    []mujica.QueuedTransfer
		queue    mujica.TransferQueueState
		disabled bool
		revision uint64
	}
	s := ui.Local(c.Root(), "388", func() state { return state{items: filesTransfers()} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if mujica.Button(c, "Reset", mujica.ButtonOptions{}).Clicked() {
		s.items = filesTransfers()
		log("Queue reset")
	}
	r := mujica.TransferQueue(c, s.items, &s.queue, mujica.TransferQueueOptions{Disabled: s.disabled, Revision: s.revision})
	r.Element.Height(260).FillWidth()
	if r.Action == mujica.TransferNone {
		return
	}
	for i := range s.items {
		if s.items[i].ID == r.ID {
			switch r.Action {
			case mujica.TransferPause:
				s.items[i].Status = mujica.TransferPaused
			case mujica.TransferResume, mujica.TransferRetry:
				s.items[i].Status = mujica.TransferActive
			case mujica.TransferRemove:
				s.items = append(s.items[:i:i], s.items[i+1:]...)
			}
			break
		}
	}
	s.revision++
	log(fmt.Sprintf("Transfer %s: %d", r.ID, r.Action))
}

func storageUsageDemo(c *ui.Context, log func(string)) {
	type state struct{ photos float64 }
	s := ui.Local(c.Root(), "389", func() state { return state{photos: 40} })
	if mujica.NumberInput(c, &s.photos, mujica.NumberInputOptions{Min: 0, Max: 92, Step: 1, Label: "Photos (GB)"}).Changed() {
		log(fmt.Sprintf("Photos: %.0f GB", s.photos))
	}
	mujica.StorageUsage(c, []mujica.StorageSegment{{Label: "Documents", Bytes: 24 << 30}, {Label: "Photos", Bytes: int64(s.photos) * (1 << 30)}, {Label: "Apps", Bytes: 12 << 30}}, mujica.StorageUsageOptions{Capacity: 128 << 30, Label: "Court disk"}).FillWidth()
}

func filesRecent() []mujica.RecentFile {
	now := time.Date(2026, 10, 5, 15, 0, 0, 0, time.UTC)
	return []mujica.RecentFile{{Path: "/court/seal.go", Opened: now}, {Path: "/court/notes.md", Opened: now.AddDate(0, 0, -1)}, {Path: "/court/crest.png", Opened: now.AddDate(0, 0, -4)}, {Path: "/court/charter.pdf", Opened: now.AddDate(0, -1, 0)}}
}

func recentFilesDemo(c *ui.Context, log func(string)) {
	type state struct {
		files    []mujica.RecentFile
		recent   mujica.RecentFilesState
		disabled bool
	}
	s := ui.Local(c.Root(), "390", func() state { return state{files: filesRecent()} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if mujica.Button(c, "Reset", mujica.ButtonOptions{}).Clicked() {
		s.files = filesRecent()
		log("Recent files reset")
	}
	r := mujica.RecentFiles(c, s.files, &s.recent, mujica.RecentFilesOptions{Disabled: s.disabled, Now: time.Date(2026, 10, 5, 15, 0, 0, 0, time.UTC)})
	r.Element.Height(320).FillWidth()
	if r.Opened != nil {
		log("Open: " + r.Opened.Path)
	}
	if r.Cleared {
		s.files = nil
		log("Recent files cleared")
	}
	if r.Removed != "" {
		for i, f := range s.files {
			if f.Path == r.Removed {
				s.files = append(s.files[:i:i], s.files[i+1:]...)
				break
			}
		}
		log("Removed: " + r.Removed)
	}
}

func archiveViewerDemo(c *ui.Context, log func(string)) {
	type state struct {
		archive mujica.ArchiveViewerState
		entries []mujica.ArchiveEntry
		status  mujica.DataStatus
		empty   bool
	}
	s := ui.Local(c.Root(), "391", func() state {
		return state{entries: []mujica.ArchiveEntry{{Path: "court/seal.go", Size: 4000, Packed: 1000}, {Path: "court/decrees/42.pdf", Size: 6000, Packed: 5000}, {Path: "README.md", Size: 1000, Packed: 400}}}
	})
	filesStatus(c, log, &s.status)
	mediaCheck(c, log, &s.empty, "Empty")
	entries := s.entries
	if s.empty {
		entries = nil
	}
	r := mujica.ArchiveViewer(c, entries, &s.archive, mujica.ArchiveViewerOptions{Status: s.status, Error: "Invalid archive"})
	r.Element.Height(280).FillWidth()
	if r.Opened != nil {
		log("Open archive entry: " + r.Opened.Path)
	}
	if r.Retried {
		s.status = mujica.DataReady
		log("Archive read retried")
	}
}
