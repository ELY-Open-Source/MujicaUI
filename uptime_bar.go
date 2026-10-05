package mujica

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
)

// UptimeDay holds availability; NaN marks missing data.
type UptimeDay struct {
	Date   Date
	Uptime float64
	// Note says what happened, as "12 min outage".
	Note string
}

// UptimeBarOptions configure an UptimeBar; Frame.Height 0 means 36.
type UptimeBarOptions struct {
	Frame ChartFrame
	// Degraded and Down default to 0.999 and 0.95.
	Degraded, Down float64
}

// UptimeBar shows daily availability with pointer and keyboard readouts.
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

	if opts.Frame.Height == 0 {
		opts.Frame.Height = 36
	}
	k := Tokens(c)
	ch := newChart(c, opts.Frame)
	n := len(days)
	cache := ui.Local(ch.res.Element, "mujica.uptime", func() uptimeCache { return uptimeCache{} })
	if cache.key != chartID(days) {
		for i, d := range days {
			if !math.IsNaN(d.Uptime) && (d.Uptime < 0 || d.Uptime > 1) {
				panic(fmt.Sprintf("mujica: UptimeBar day %d outside 0..1", i))
			}
		}
		cache.key, cache.width = chartID(days), -1
	}
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
	pal := internal.Palette(k)
	words := [...]string{msg(c, "uptime_bar.ok"), msg(c, "uptime_bar.degraded"), msg(c, "uptime_bar.down"), msg(c, "uptime_bar.nodata")}
	res, _ := ch.build(chartSpec{
		n: n,
		draw: func(p *ui.Painter, r ui.Rect, active int) {
			b := band(r)
			cache.prepare(days, max(1, int(math.Ceil(float64(r.W*p.Scale())))))
			for _, slot := range cache.slots {
				d := days[slot.index]
				x, w := internal.Snap(p, b.Start(slot.first)), max(b.Start(slot.last)+b.Width()-b.Start(slot.first), 1)
				lv := level(d.Uptime)
				col := [...]ui.Color{pal[2], pal[3], k.Danger, k.TextMuted}[lv]
				// Shorter bars distinguish degraded days and outages.
				h := r.H * [...]float32{1, 0.75, 0.5, 1}[lv]
				bar := ui.Rect{X: x, Y: r.Y + r.H - h, W: w, H: h}
				if lv == 3 {
					p.Stroke(bar, col, 1, 1)
				} else {
					p.Fill(bar, col, 1)
				}
				if active >= slot.first && active <= slot.last {
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

type uptimeSlot struct{ first, last, index int }
type uptimeCache struct {
	key   string
	width int
	slots []uptimeSlot
}

func (s *uptimeCache) prepare(days []UptimeDay, width int) {
	if s.width == width {
		return
	}
	s.width, s.slots = width, s.slots[:0]
	n := len(days)
	for first := 0; first < n; {
		last := min(n-1, max(first, ((first*width/n+1)*n-1)/width))
		index := first
		for i := first + 1; i <= last; i++ {
			if math.IsNaN(days[index].Uptime) || !math.IsNaN(days[i].Uptime) && days[i].Uptime < days[index].Uptime {
				index = i
			}
		}
		s.slots = append(s.slots, uptimeSlot{first, last, index})
		first = last + 1
	}
}
