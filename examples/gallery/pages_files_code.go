package main

const filesFileExplorerCode = `var state mujica.FileExplorerState
opts := mujica.FileExplorerOptions{
	Roots:   []string{"/court"},
	Folders: func(string) []string { return nil },
	Entries: func(string) []mujica.FileEntry {
		return []mujica.FileEntry{{Path: "/court/seal.go", Name: "seal.go", Size: 1024}}
	},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FileExplorer(c, &state, opts).Element.Height(240)
}
ui.Render(view, 600, 300, 1)`

const filesFileGridCode = `entries := []mujica.FileEntry{{Path: "/court/seal.go", Name: "seal.go"}}
var state mujica.FileGridState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FileGrid(c, entries, &state, mujica.FileGridOptions{}).Element.Height(240)
}
ui.Render(view, 400, 300, 1)`

const filesFileIconCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FileIcon(c, "seal.go", mujica.FileIconOptions{Size: 32})
}
ui.Render(view, 100, 80, 1)`

const filesFilePreviewCode = `file := mujica.PreviewFile{Name: "seal.go", Text: "package court\nfunc Seal() {}"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FilePreview(c, file, mujica.FilePreviewOptions{}).Element.Height(240)
}
ui.Render(view, 500, 300, 1)`

const filesPathBarCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PathBar(c, "/court/chronicles/seals", mujica.PathBarOptions{})
}
ui.Render(view, 400, 100, 1)`

const filesRenameInlineCode = `name, editing := "seal.txt", true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.RenameInline(c, &name, &editing, mujica.RenameInlineOptions{})
}
ui.Render(view, 400, 100, 1)`

const filesFileOperationProgressCode = `op := mujica.FileOperation{Kind: mujica.FileCopy, Done: 2, Total: 10, Current: "seal.go"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FileOperationProgress(c, op, mujica.FileOperationOptions{})
}
ui.Render(view, 500, 240, 1)`

const filesTransferQueueCode = `items := []mujica.QueuedTransfer{{ID: "a", Name: "ledger.zip", Upload: true, Status: mujica.TransferActive, Done: 25, Total: 100}}
var state mujica.TransferQueueState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TransferQueue(c, items, &state, mujica.TransferQueueOptions{}).Element.Height(240)
}
ui.Render(view, 500, 300, 1)`

const filesStorageUsageCode = `segments := []mujica.StorageSegment{{Label: "Documents", Bytes: 24 << 30}, {Label: "Photos", Bytes: 40 << 30}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.StorageUsage(c, segments, mujica.StorageUsageOptions{Capacity: 128 << 30, Label: "Court disk"})
}
ui.Render(view, 500, 240, 1)`

const filesRecentFilesCode = `now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
files := []mujica.RecentFile{{Path: "/court/seal.go", Opened: now}}
var state mujica.RecentFilesState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.RecentFiles(c, files, &state, mujica.RecentFilesOptions{Now: now}).Element.Height(240)
}
ui.Render(view, 500, 300, 1)`

const filesArchiveViewerCode = `entries := []mujica.ArchiveEntry{{Path: "court/seal.go", Size: 4000, Packed: 1000}}
var state mujica.ArchiveViewerState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ArchiveViewer(c, entries, &state, mujica.ArchiveViewerOptions{}).Element.Height(240)
}
ui.Render(view, 600, 300, 1)`

const filesPresenceAvatarsCode = `people := []mujica.Collaborator{{ID: "ada", Name: "Ada Lovelace"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PresenceAvatars(c, people, mujica.PresenceAvatarsOptions{})
}
ui.Render(view, 300, 100, 1)`

const filesReactionsCode = `items := []mujica.Reaction{{Emoji: "\U0001f44d", Count: 3, Mine: true}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Reactions(c, items, mujica.ReactionsOptions{Choices: []string{"\U0001f44d", "\U0001f389"}})
}
ui.Render(view, 400, 100, 1)`

const filesShareDialogCode = `open := false
members := []mujica.ShareMember{{ID: "ada", Name: "Ada", Permission: mujica.PermissionEdit, Owner: true}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.Button(c, "Share", mujica.ButtonOptions{}).Clicked() {
		open = true
	}
	mujica.ShareDialog(c, &open, &members, mujica.ShareDialogOptions{Title: "Ledger", Link: "https://court.org/ledger"})
}
ui.Render(view, 600, 400, 1)`

const filesPermissionSelectCode = `value := mujica.PermissionView
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PermissionSelect(c, &value, mujica.PermissionSelectOptions{})
}
ui.Render(view, 300, 100, 1)`

const filesLiveIndicatorCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.LiveIndicator(c, mujica.LiveOffline, mujica.LiveIndicatorOptions{})
}
ui.Render(view, 500, 100, 1)`

const filesRemoteCursorCode = `carets := []mujica.RemoteCaret{{ID: "ada", Name: "Ada", Offset: 4, Anchor: 4}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	ui.Box(c).Size(300, 100).Children(func() {
		mujica.RemoteCursor(c, carets, mujica.RemoteCursorOptions{
			Caret: func(offset int) ui.Rect { return ui.Rect{X: float32(offset) * 8, Y: 20, H: 20} },
		})
	})
}
ui.Render(view, 320, 120, 1)`
