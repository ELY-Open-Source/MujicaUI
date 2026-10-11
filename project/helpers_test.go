package project

import (
	"fmt"
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal/snaptest"
)

// snapshot compares view in both themes with testdata/snapshots/<os>/<name>-<mode>.png.
func snapshot(t *testing.T, name string, w, h int, view func(c *ui.Context)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, nil)
}

func oct(day int) core.Date { return core.Date{Year: 2026, Month: time.October, Day: day} }

var calendarToday = core.Date{Year: 2026, Month: time.October, Day: 5}

// datetimeTestLabel is the en-US spoken name of a day cell.
func datetimeTestLabel(d core.Date) string {
	return fmt.Sprintf("%v, %v %d, %d", d.Weekday(), d.Month, d.Day, d.Year)
}

func datetimeNear(c color.RGBA, w ui.Color) bool {
	d := func(a, b uint8) int {
		if a > b {
			return int(a - b)
		}
		return int(b - a)
	}
	return d(c.R, w.R)+d(c.G, w.G)+d(c.B, w.B) < 60
}

func calDrag(tt *ui.Tester, x0, y0, x1, y1 float32) {
	tt.Press(x0, y0)
	for i := 1; i <= 4; i++ {
		f := float32(i) / 4
		tt.Move(x0+(x1-x0)*f, y0+(y1-y0)*f)
	}
	tt.Release(x1, y1)
}

func snapshotTheme(c *ui.Context, mode string) {
	m := core.Light
	if mode == "dark" {
		m = core.Dark
	}
	core.Use(c, core.Settings{Mode: m})
}
