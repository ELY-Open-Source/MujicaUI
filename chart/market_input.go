package chart

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

// input handles pan, zoom, alerts, drawings, and keyboard following.
func (m *mktCore) input(plot ui.Element) {
	if plot.IsDisabled() {
		return
	}
	n := len(m.candles)
	room := m.room()
	changed := false
	plot.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll || ev.DY == 0 {
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
	// Keep the keyboard mark visible.
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
		// A click is placed by input at the pointer; only a key press lands here.
		if i := m.act - 1; i >= 0 && !used && !d.Locked && d.Tool != ToolSelect {
			d.place(m, DrawingPoint{m.candles[i].Time, m.candles[i].Close})
			used = true
		}
	}
	if !used {
		m.alertsKeyboard(plot)
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

func (m *mktCore) alertsKeyboard(plot ui.Element) {
	if m.opts.Alerts == nil || len(*m.opts.Alerts) == 0 || m.base != 0 {
		return
	}
	a := *m.opts.Alerts
	m.st.alertFocus = min(m.st.alertFocus, len(a)-1)
	selected := false
	if plot.Shortcut(0, ui.KeyPageDown) {
		m.st.alertFocus = (m.st.alertFocus + 1) % len(a)
		selected = true
	}
	if plot.Shortcut(0, ui.KeyPageUp) {
		m.st.alertFocus = (max(m.st.alertFocus, 0) - 1 + len(a)) % len(a)
		selected = true
	}
	if selected && m.opts.Drawings != nil {
		d := m.opts.Drawings
		d.sel, d.pending, d.Tool = 0, nil, ToolSelect
	}
	delta := 0
	if m.opts.Drawings == nil || m.opts.Drawings.Selected() < 0 {
		if plot.Shortcut(ui.Alt, ui.KeyUp) {
			delta = 1
		}
		if plot.Shortcut(ui.Alt, ui.KeyDown) {
			delta = -1
		}
	}
	if i := m.st.alertFocus; i >= 0 && (selected || delta != 0) {
		spec := input.NumberInputResolve(input.NumberInputOptions{Step: math.Pow10(-m.dec), Precision: m.dec})
		if delta != 0 {
			a[i].Price = spec.StepBy(a[i].Price, delta)
			slog.Debug("mujica: price_alert_line moved", "id", a[i].ID)
		}
		m.c.Announce(fmt.Sprintf(core.Msg(m.c, "price_alert_line.tag"), core.PriceUngrouped(a[i].Price, m.dec)))
		m.c.Invalidate()
	}
}
