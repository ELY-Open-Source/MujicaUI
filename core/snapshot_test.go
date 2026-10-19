package core

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/internal/snaptest"
)

// snapshot compares view in both themes with testdata/snapshots/<os>/<name>-<mode>.png.
func snapshot(t *testing.T, name string, w, h int, view func(c *ui.Context)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, nil)
}

// snapshotAfter is snapshot with drive sending input before the capture.
func snapshotAfter(t *testing.T, name string, w, h int, view func(c *ui.Context), drive func(tt *ui.Tester)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, drive)
}

func snapshotTheme(c *ui.Context, mode string) {
	m := Light
	if mode == "dark" {
		m = Dark
	}
	Use(c, Settings{Mode: m})
}

func encode(t *testing.T, img image.Image) []byte {
	t.Helper()
	return snaptest.Encode(t, img)
}
