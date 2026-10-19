package navigation

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal/snaptest"
)

// snapshot compares view in both themes with testdata/snapshots/<os>/<name>-<mode>.png.
func snapshot(t *testing.T, name string, w, h int, view func(c *ui.Context)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, nil)
}

type menubarApp struct {
	wrap   bool
	zoom   string
	chosen []string
}

// snapshotAfter is snapshot with drive sending input before the capture.
func snapshotAfter(t *testing.T, name string, w, h int, view func(c *ui.Context), drive func(tt *ui.Tester)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, drive)
}

func snapshotTheme(c *ui.Context, mode string) {
	m := core.Light
	if mode == "dark" {
		m = core.Dark
	}
	core.Use(c, core.Settings{Mode: m})
}

func (a *menubarApp) menus() []MenubarMenu {
	return []MenubarMenu{
		{Label: "File", Items: []MenubarItem{
			{ID: "new", Label: "New", Shortcut: MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyN}},
			{ID: "open", Label: "Open…", Shortcut: MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyO}},
			{Label: "Export", Items: []MenubarItem{
				{ID: "pdf", Label: "PDF"},
				{ID: "png", Label: "Image", Disabled: true},
				{ID: "svg", Label: "Vector"},
			}},
			{Kind: MenubarSeparator},
			{ID: "quit", Label: "Quit"},
		}},
		{Label: "View", Items: []MenubarItem{
			{ID: "wrap", Label: "Word wrap", Kind: MenubarCheck, Checked: a.wrap},
			{Kind: MenubarSeparator},
			{ID: "zoom-fit", Label: "Fit", Kind: MenubarRadio, Checked: a.zoom == "zoom-fit"},
			{ID: "zoom-100", Label: "Actual size", Kind: MenubarRadio, Checked: a.zoom == "zoom-100"},
		}},
		{Label: "Help", Items: []MenubarItem{{ID: "about", Label: "About"}}},
	}
}
