package mujica

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
)

// UptimeDay is one day's availability, 0 to 1; NaN means no data.
type UptimeDay struct {
	Date   Date
	Uptime float64
	// Note says what happened, as "12 min outage".
	Note string
}

// UptimeBarOptions configure an UptimeBar; Frame.Height 0 means 36.
type UptimeBarOptions struct {
	Frame ChartFrame
	// Degraded and Down are the uptime limits below which a day is marked; zero means 0.999 and 0.95.
	Degraded, Down float64
}

// UptimeBar draws one bar per day by availability; hover or arrows read a day's detail.
func UptimeBar(c *ui.Context, days []UptimeDay, opts UptimeBarOptions) ChartResult {
	deg, down := opts.Degraded, opts.Down
	if deg == 0 {
		deg = 0.999
	}
	if down == 0 {
		down = 0.95
	}
	if !(down < deg && deg <= 1 && down > 0) {
		panic("mujica: UptimeBar needs 0 < Down < Degraded <= 1")
	}
	for i, d := range days {
		if !math.IsNaN(d.Uptime) && (d.Uptime < 0 || d.Uptime > 1) {
			panic(fmt.Sprintf("mujica: UptimeBar day %d uptime %v outside 0..1", i, d.Uptime))
		}
	}
	if opts.Frame.Height == 0 {
		opts.Frame.Height = 36
	}
	k := Tokens(c)
	ch := newChart(c, opts.Frame)
	n := len(days)
	band := func(r ui.Rect) internal.Band { return internal.Band{N: n, R0: r.X, R1: r.X + r.W, Pad: 0.25} }
	level := func(v float64) int { // 0 ok, 1 degraded, 2 down, 3 no data
		switch {
		case math.IsNaN(v):
			return 3
		case v < down:
			return 2
		case v < deg:
			return 1
		}
		return 0
	}
	words := [...]string{msg(c, "uptime_bar.ok"), msg(c, "uptime_bar.degraded"), msg(c, "uptime_bar.down"), msg(c, "uptime_bar.nodata")}
	res, _ := ch.build(chartSpec{
		n: n,
		draw: func(p *ui.Painter, r ui.Rect, active int) {
			b := band(r)
			for i, d := range days {
				x, w := internal.Snap(p, b.Start(i)), max(b.Width(), 1)
				lv := level(d.Uptime)
				col := [...]ui.Color{k.Success, k.Warning, k.Danger, k.Border}[lv]
				// Shorter bars mark bad days, so color is not the only cue.
				h := r.H * [...]float32{1, 0.75, 0.5, 1}[lv]
				p.Fill(ui.Rect{X: x, Y: r.Y + r.H - h, W: w, H: h}, col, 1)
				if i == active {
					p.Stroke(ui.Rect{X: x - 2, Y: r.Y - 2, W: w + 4, H: r.H + 4}, k.Focus, 2, 1.5)
				}
			}
		},
		pick: func(r ui.Rect, x, _ float32) int {
			b := band(r)
			return max(0, min(n-1, int((x-r.X)/b.Step())))
		},
		readout: func(i int) string {
			d := days[i]
			s := fmt.Sprintf("%s · %s", d.Date.String(), words[level(d.Uptime)])
			if !math.IsNaN(d.Uptime) {
				s += fmt.Sprintf(" · %.2f%%", d.Uptime*100)
			}
			if d.Note != "" {
				s += " · " + d.Note
			}
			return s
		},
	})
	return res
}
