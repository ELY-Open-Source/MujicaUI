package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/internal"
)

type c2PartsState struct {
	status, kind                          string
	disabled, crosshair, draggable, fixed bool
	minor                                 float64
	notes                                 []mujica.ChartNote
}

func c2TooltipDemo(c *ui.Context, log func(string))    { c2PartsDemo(c, log, "122", "ChartTooltip") }
func c2CrosshairDemo(c *ui.Context, log func(string))  { c2PartsDemo(c, log, "123", "ChartCrosshair") }
func c2AxisDemo(c *ui.Context, log func(string))       { c2PartsDemo(c, log, "124", "ChartAxis") }
func c2GridDemo(c *ui.Context, log func(string))       { c2PartsDemo(c, log, "125", "ChartGrid") }
func c2AnnotationDemo(c *ui.Context, log func(string)) { c2PartsDemo(c, log, "127", "ChartAnnotation") }

func c2PartsDemo(c *ui.Context, log func(string), id, name string) {
	s := ui.Local(c.Root(), id, func() c2PartsState {
		return c2PartsState{status: "Ready", kind: "Linear", crosshair: true, minor: 2, draggable: true, notes: []mujica.ChartNote{
			{Kind: mujica.NoteHLine, Value: 50, Text: "Goal"},
			{Kind: mujica.NoteVLine, Value: 3, Text: "Advent"},
			{Kind: mujica.NoteBand, Value: 6, End: 8, Text: "Lent"},
			{Kind: mujica.NoteText, Value: 1, End: 90, Text: "Feast"},
		}}
	})
	c1Section(c, name, func() {
		c2Check(c, &s.disabled, "Disabled", log)
		switch id {
		case "123":
			c2Check(c, &s.crosshair, "Crosshair", log)
		case "124":
			c2Select(c, &s.kind, "Axis kind", []string{"Linear", "Log", "Time", "Category"}, log)
			c2Check(c, &s.fixed, "Fixed ticks", log)
		case "125":
			c2Number(c, &s.minor, "Minor", 0, 8, 1, log)
		case "127":
			c2Check(c, &s.draggable, "Draggable", log)
		}
	}, func() {
		f, show := c2Frame(c, name, &s.status, log)
		ui.Column(c).Disabled(s.disabled).Children(func() {
			c1Log(log, name, c2PartsChart(c, f, show, s, id, log))
		})
		if id == "127" {
			ui.Text(c, fmt.Sprintf("Goal %.1f · Advent %.1f · Lent %.1f–%.1f · Feast %.1f, %.1f", s.notes[0].Value, s.notes[1].Value, s.notes[2].Value, s.notes[2].End, s.notes[3].Value, s.notes[3].End))
		}
	})
}

func c2PartsChart(c *ui.Context, frame mujica.ChartFrame, show bool, s *c2PartsState, id string, log func(string)) mujica.ChartResult {
	xo := mujica.ChartAxisOptions{Max: 10}
	yo := mujica.ChartAxisOptions{Side: mujica.SideLeft, Max: 100}
	pts := [][2]float64{{2, 30}, {5, 80}, {8, 55}}
	if id == "124" {
		switch s.kind {
		case "Log":
			yo.Kind, yo.Min = mujica.AxisLog, 1
		case "Time":
			xo.Kind = mujica.AxisTime
			xo.From = time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
			xo.To = xo.From.Add(10 * time.Hour)
			for i := range pts {
				pts[i][0] = float64(xo.From.Unix()) + pts[i][0]*3600
			}
		case "Category":
			xo.Kind, xo.Categories = mujica.AxisCategory, []string{"Nave", "Choir", "Crypt"}
			for i := range pts {
				pts[i][0] = float64(i)
			}
		}
		if s.fixed {
			yo.Ticks = []float64{10, 50, 90}
		}
	}
	x, y := mujica.ChartAxis(c, xo), mujica.ChartAxis(c, yo)
	grid := mujica.ChartGrid(c, mujica.ChartGridOptions{Y: &y, Minor: int(s.minor)})
	hair := mujica.ChartCrosshair(c, mujica.ChartCrosshairOptions{X: &x, Y: &y})
	var notes []mujica.ChartNote
	if id == "127" {
		notes = s.notes
		for i := range notes {
			notes[i].Draggable = s.draggable
		}
	}
	annotation := mujica.ChartAnnotation(c, notes, mujica.ChartAnnotationOptions{X: &x, Y: &y})
	plotOf := func(r ui.Rect) ui.Rect {
		return ui.Rect{X: r.X + 64, Y: r.Y + 16, W: max(1, r.W-100), H: max(1, r.H-56)}
	}
	at := func(r ui.Rect, i int) (float32, float32) {
		plot := plotOf(r)
		return x.Map(plot, pts[i][0]), y.Map(plot, pts[i][1])
	}
	marks := len(pts)
	if !show {
		marks = 0
	}
	colors := internal.Palette(mujica.Tokens(c))
	return mujica.Chart(c, frame, mujica.ChartSpec{
		Marks: marks,
		Pick: func(r ui.Rect, px, py float32) int {
			return hair.Snap(len(pts), func(i int) (float32, float32) { return at(r, i) }, px, py)
		},
		Anchor:  at,
		Readout: func(i int) string { return fmt.Sprintf("x %s · y %s", x.Label(pts[i][0]), y.Label(pts[i][1])) },
		Tooltip: func(i int) []mujica.ChartTooltipLine {
			return []mujica.ChartTooltipLine{{Text: x.Label(pts[i][0])}, {Text: "Pews " + y.Label(pts[i][1]), Color: colors[0]}, {Text: "Candles " + y.Label(pts[i][1]/2), Color: colors[1], Shape: 1}}
		},
		Input: func(el *ui.Element, r ui.Rect) {
			if annotation.Drag(el, plotOf(r)) {
				log("ChartAnnotation: moved")
			}
		},
		Draw: func(p *ui.Painter, r ui.Rect, active int) {
			plot := plotOf(r)
			grid.Draw(p, plot)
			x.Draw(p, plot)
			y.Draw(p, plot)
			annotation.Draw(p, plot)
			for i := range pts {
				px, py := at(r, i)
				internal.Marker(p, i, px, py, 10, colors[i])
			}
			if id == "123" && s.crosshair && active >= 0 {
				px, py := at(r, active)
				hair.Draw(p, plot, px, py)
			}
		},
	})
}
