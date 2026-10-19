package devtools

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal/snaptest"
)

// snapshot compares view in both themes with testdata/snapshots/<os>/<name>-<mode>.png.
func snapshot(t *testing.T, name string, w, h int, view func(c *ui.Context)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, nil)
}

// overlayUntil advances the Tester clock until cond holds, for at most 5 s.
func overlayUntil(tt *ui.Tester, cond func() bool) bool {
	end := tt.Now().Add(5 * time.Second)
	for !cond() {
		if tt.Now().After(end) {
			return false
		}
		tt.Advance(2 * time.Millisecond)
	}
	return true
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
