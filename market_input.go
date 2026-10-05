package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"
)

// input pans, zooms, drags alerts and drawings, and keeps the focused candle in view.
func (m *mktCore) input(plot *ui.Element) {
	if plot.IsDisabled() {
		return
	}
	n := len(m.candles)
	room := m.room()
	changed := false
	plot.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll || plot.IsDisabled() || ev.DY == 0 {
			return false
		}
		f := 0.9
		if ev.DY > 0 {
			f = 1 / 0.9
		}
		marketZoom(m.view, n, m.room(), f)
		return true
	})
	switch {
	case plot.Shortcut(ui.Shift, ui.KeyLeft):
		marketPan(m.view, n, room, -max(1, (m.to-m.from)/10))
		changed = true
	case plot.Shortcut(ui.Shift, ui.KeyRight):
		marketPan(m.view, n, room, max(1, (m.to-m.from)/10))
		changed = true
	case plot.Shortcut(0, ui.KeyEqual):
		marketZoom(m.view, n, room, 0.8)
		changed = true
	case plot.Shortcut(0, ui.KeyMinus):
		marketZoom(m.view, n, room, 1.25)
		changed = true
	}
	// The keyboard mark may leave the window: the view follows it.
	f := m.ch.st.focus
	if f >= 0 && f != m.st.focus && !changed && (f < m.from || f >= m.to) {
		if f < m.from {
			marketPan(m.view, n, room, f-m.from)
		} else {
			marketPan(m.view, n, room, f-m.to+1)
		}
		changed = true
	}
	m.st.focus = f
	b := plot.Bounds()
	g := m.geom(ui.Rect{X: b.X, Y: b.Y, W: b.W, H: b.H})
	px, py, over := plot.PointerPosition()
	x, y := b.X+px, b.Y+py
	dx, dy, held := plot.Dragged()
	used := false
	if d := m.opts.Drawings; d != nil {
		used = d.input(m, plot, g, x, y, over, held, dx, dy)
	}
	if !used {
		used = m.alertsInput(g, x, y, held, dy)
	}
	if !used && held && dx != 0 && g.band.Step() > 0 {
		m.st.dragRest += dx
		if steps := int(m.st.dragRest / g.band.Step()); steps != 0 {
			m.st.dragRest -= float32(steps) * g.band.Step()
			marketPan(m.view, n, room, -steps)
			changed = true
		}
	}
	if !held {
		m.st.dragRest = 0
	}
	if changed {
		slog.Debug("mujica: candlestick_chart view", "count", m.view.Count, "end", m.view.End)
		m.c.Invalidate()
	}
}
