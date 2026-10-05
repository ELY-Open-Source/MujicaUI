package feedback

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal/snaptest"
)

// snapshot compares view in both themes with testdata/snapshots/<os>/<name>-<mode>.png.
func snapshot(t *testing.T, name string, w, h int, view func(c *ui.Context)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, nil)
}

func encode(t *testing.T, img image.Image) []byte {
	t.Helper()
	return snaptest.Encode(t, img)
}

// snapshotAfter is snapshot with drive sending input before the capture.
func snapshotAfter(t *testing.T, name string, w, h int, view func(c *ui.Context), drive func(tt *ui.Tester)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, drive)
}

// tabTo presses Tab until the element labeled label has the focus.
func tabTo(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	for range 60 {
		if tt.Focused(label) {
			return
		}
		tt.Key(0, ui.KeyTab)
	}
	t.Fatalf("Tab never reached %q", label)
}

type chartHarness struct {
	res      chart.ChartResult
	changes  int
	disabled bool
	build    func(c *ui.Context) chart.ChartResult
}

// chartReadout returns the readout line, the text after the plot.
func chartReadout(tt *ui.Tester) string {
	texts := tt.Texts()
	return texts[len(texts)-1]
}

func snapshotTheme(c *ui.Context, mode string) {
	m := core.Light
	if mode == "dark" {
		m = core.Dark
	}
	core.Use(c, core.Settings{Mode: m})
}

func (h *chartHarness) view(c *ui.Context) {
	core.Use(c, core.Settings{Mode: core.Light})
	ui.Column(c).Padding(16).Width(480).Disabled(h.disabled).Children(func() {
		h.res = h.build(c)
		if h.res.Changed() {
			h.changes++
		}
	})
}
