package main

const filesFileExplorerCode = `var state files.FileExplorerState
opts := files.FileExplorerOptions{
	Roots:   []string{"/court"},
	Folders: func(string) []string { return nil },
	Entries: func(string) []files.FileEntry {
		return []files.FileEntry{{Path: "/court/seal.go", Name: "seal.go", Size: 1024}}
	},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.FileExplorer(c, &state, opts).Element.Height(240)
}
ui.Render(view, 600, 300, 1)`

const filesFileGridCode = `entries := []files.FileEntry{{Path: "/court/seal.go", Name: "seal.go"}}
var state files.FileGridState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.FileGrid(c, entries, &state, files.FileGridOptions{}).Element.Height(240)
}
ui.Render(view, 400, 300, 1)`

const filesFileIconCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.FileIcon(c, "seal.go", files.FileIconOptions{Size: 32})
}
ui.Render(view, 100, 80, 1)`

const filesFilePreviewCode = `file := files.PreviewFile{Name: "seal.go", Text: "package court\nfunc Seal() {}"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.FilePreview(c, file, files.FilePreviewOptions{}).Element.Height(240)
}
ui.Render(view, 500, 300, 1)`

const filesPathBarCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.PathBar(c, "/court/chronicles/seals", files.PathBarOptions{})
}
ui.Render(view, 400, 100, 1)`

const filesRenameInlineCode = `name, editing := "seal.txt", true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.RenameInline(c, &name, &editing, files.RenameInlineOptions{})
}
ui.Render(view, 400, 100, 1)`

const filesFileOperationProgressCode = `op := files.FileOperation{Kind: files.FileCopy, Done: 2, Total: 10, Current: "seal.go"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.FileOperationProgress(c, op, files.FileOperationOptions{})
}
ui.Render(view, 500, 240, 1)`

const filesTransferQueueCode = `items := []files.QueuedTransfer{{ID: "a", Name: "ledger.zip", Upload: true, Status: files.TransferActive, Done: 25, Total: 100}}
var state files.TransferQueueState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.TransferQueue(c, items, &state, files.TransferQueueOptions{}).Element.Height(240)
}
ui.Render(view, 500, 300, 1)`

const filesStorageUsageCode = `segments := []files.StorageSegment{{Label: "Documents", Bytes: 24 << 30}, {Label: "Photos", Bytes: 40 << 30}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.StorageUsage(c, segments, files.StorageUsageOptions{Capacity: 128 << 30, Label: "Court disk"})
}
ui.Render(view, 500, 240, 1)`

const filesRecentFilesCode = `now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
recent := []files.RecentFile{{Path: "/court/seal.go", Opened: now}}
var state files.RecentFilesState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.RecentFiles(c, recent, &state, files.RecentFilesOptions{Now: now}).Element.Height(240)
}
ui.Render(view, 500, 300, 1)`

const filesArchiveViewerCode = `entries := []files.ArchiveEntry{{Path: "court/seal.go", Size: 4000, Packed: 1000}}
var state files.ArchiveViewerState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.ArchiveViewer(c, entries, &state, files.ArchiveViewerOptions{}).Element.Height(240)
}
ui.Render(view, 600, 300, 1)`

const filesPresenceAvatarsCode = `people := []files.Collaborator{{ID: "ada", Name: "Ada Lovelace"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.PresenceAvatars(c, people, files.PresenceAvatarsOptions{})
}
ui.Render(view, 300, 100, 1)`

const filesReactionsCode = `items := []files.Reaction{{Emoji: "\U0001f44d", Count: 3, Mine: true}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.Reactions(c, items, files.ReactionsOptions{Choices: []string{"\U0001f44d", "\U0001f389"}})
}
ui.Render(view, 400, 100, 1)`

const filesShareDialogCode = `open := false
members := []files.ShareMember{{ID: "ada", Name: "Ada", Permission: files.PermissionEdit, Owner: true}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.Button(c, "Share", input.ButtonOptions{}).Clicked() {
		open = true
	}
	files.ShareDialog(c, &open, &members, files.ShareDialogOptions{Title: "Ledger", Link: "https://court.org/ledger"})
}
ui.Render(view, 600, 400, 1)`

const filesPermissionSelectCode = `value := files.PermissionView
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.PermissionSelect(c, &value, files.PermissionSelectOptions{})
}
ui.Render(view, 300, 100, 1)`

const filesLiveIndicatorCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.LiveIndicator(c, files.LiveOffline, files.LiveIndicatorOptions{})
}
ui.Render(view, 500, 100, 1)`

const filesRemoteCursorCode = `carets := []files.RemoteCaret{{ID: "ada", Name: "Ada", Offset: 4, Anchor: 4}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Box(c).Size(300, 100).Children(func() {
		files.RemoteCursor(c, carets, files.RemoteCursorOptions{
			Caret: func(offset int) ui.Rect { return ui.Rect{X: float32(offset) * 8, Y: 20, H: 20} },
		})
	})
}
ui.Render(view, 320, 120, 1)`
