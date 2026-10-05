package mujica

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// CalendarHeatmapOptions configure a CalendarHeatmap.
type CalendarHeatmapOptions struct {
	Frame ChartFrame
	// WeekStart picks each column's first weekday; zero follows the locale.
	WeekStart WeekStart
	Format    func(v float64) string
}

// CalendarHeatmap colors each day of *year, a column per week.
func CalendarHeatmap(c *ui.Context, year *int, values map[Date]float64, opts CalendarHeatmapOptions) ChartResult {
	if *year < 1 || *year > 9999 {
		panic("mujica: CalendarHeatmap year out of range")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	size := FontSize(c, theme.CaptionSize)
	first := Date{Year: *year, Month: time.January, Day: 1}
	days := first.AddMonths(12).civil().Sub(first.civil()).Hours() / 24
	n := int(days)
	lead := (int(first.Weekday()) - int(calendarFirstWeekday(c, opts.WeekStart)) + 7) % 7
	weeks := (lead + n + 6) / 7
	top := 0.0
	for i := range n {
		if v, ok := values[first.AddDays(i)]; ok {
			top = max(top, v)
		}
	}
	// level is a day's color step, 0 to 4.
	level := func(i int) int {
		v, ok := values[first.AddDays(i)]
		if !ok || v <= 0 || top <= 0 {
			return 0
		}
		return 1 + min(3, int(v/top*4))
	}
	var res ChartResult
	wrap := ui.Column(c).Gap(u * 2).MinWidth(0)
	wrap.Children(func() {
		ui.Row(c).Gap(u).AlignItems(ui.Center).Children(func() {
			prev := iconAction(c, icons.Must("chevron-left"), msg(c, "calendar_heatmap.prev")).Disabled(*year <= 1)
			ui.Text(c, strconv.Itoa(*year)).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).Width(64).TextAlign(ui.Center)
			next := iconAction(c, icons.Must("chevron-right"), msg(c, "calendar_heatmap.next")).Disabled(*year >= 9999)
			if prev.Clicked() && !prev.IsDisabled() {
				*year--
				res.value = true
			}
			if next.Clicked() && !next.IsDisabled() {
				*year++
				res.value = true
			}
		})
		ch := newChart(c, opts.Frame)
		geom := func(r ui.Rect) (x, y, step float32) {
			step = min((r.W-36)/float32(weeks), (r.H-48)/7)
			return r.X + 32, r.Y + 18, step
		}
		cell := func(r ui.Rect, i int) ui.Rect {
			x, y, step := geom(r)
			at := lead + i
			return ui.Rect{X: x + float32(at/7)*step, Y: y + float32(at%7)*step, W: step - 2, H: step - 2}
		}
		spec := chartSpec{n: n}
		spec.readout = func(i int) string {
			d := first.AddDays(i)
			v, ok := values[d]
			s := d.String() + " · "
			if !ok {
				return s + msg(c, "chart.missing")
			}
			return s + chartValue(c, v, opts.Format)
		}
		spec.move = func(i int, key ui.Key) int {
			step := map[ui.Key]int{ui.KeyLeft: -7, ui.KeyRight: 7, ui.KeyUp: -1, ui.KeyDown: 1}[key]
			return max(0, min(n-1, i+step))
		}
		spec.pick = func(r ui.Rect, x, y float32) int {
			x0, y0, step := geom(r)
			col, row := int((x-x0)/step), int((y-y0)/step)
			i := col*7 + row - lead
			if x < x0 || y < y0 || row > 6 || i < 0 || i >= n {
				return -1
			}
			return i
		}
		spec.draw = func(p *ui.Painter, r ui.Rect, active int) {
			x0, y0, step := geom(r)
			for i := range n {
				b := cell(r, i)
				if lv := level(i); lv == 0 {
					p.Fill(b, k.SurfaceHover, 2)
				} else {
					p.Fill(b, internal.Ramp(k, float64(lv)/4), 2)
				}
				if d := first.AddDays(i); d.Day == 1 {
					chartText(p, d.Format("Jan"), b.X, y0-4, size, k.TextMuted, 0, 1)
				}
				if i == active {
					p.Stroke(ui.Rect{X: b.X - 2, Y: b.Y - 2, W: b.W + 4, H: b.H + 4}, k.Text, 3, 2)
				}
			}
			for row := 0; row < 7; row += 2 {
				wd := time.Weekday((int(calendarFirstWeekday(c, opts.WeekStart)) + row) % 7)
				chartText(p, wd.String()[:3], x0-4, y0+float32(row)*step+step/2, size, k.TextMuted, 1, 0.5)
			}
			ly := y0 + 7*step + 8
			w, _ := chartText(p, msg(c, "calendar_heatmap.less"), x0, ly, size, k.TextMuted, 0, 0)
			for lv := range 5 {
				sw := ui.Rect{X: x0 + w + 6 + float32(lv)*(step+2), Y: ly + 2, W: step - 2, H: step - 2}
				if lv == 0 {
					p.Fill(sw, k.SurfaceHover, 2)
				} else {
					p.Fill(sw, internal.Ramp(k, float64(lv)/4), 2)
				}
			}
			chartText(p, msg(c, "calendar_heatmap.more"), x0+w+12+5*(step+2), ly, size, k.TextMuted, 0, 0)
		}
		res2, _ := ch.build(spec)
		res.Point, res.changed, res.retry = res2.Point, res2.changed, res2.retry
	})
	if res.value {
		slog.Debug("mujica: calendar_heatmap year", "year", *year)
	}
	res.Element = wrap
	return res
}
