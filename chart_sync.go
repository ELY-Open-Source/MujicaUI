package mujica

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChartSync links charts by time: the same window, zoom and crosshair. The caller holds it.
type ChartSync struct {
	// Count is the linked zoom in candles; zero is the default 120.
	Count int
	// End is the newest linked time shown; zero follows each chart's newest candle.
	End time.Time
	// Cursor is the crosshair time; zero shows none.
	Cursor time.Time
	// Unlinked charts keep their own view and crosshair.
	Unlinked map[string]bool
	owner    string
}

func (s *ChartSync) linked(id string) bool { return !s.Unlinked[id] }

// load puts the linked window into view before a chart builds.
func (s *ChartSync) load(id string, candles []Candle, view *MarketView) {
	if !s.linked(id) {
		return
	}
	view.Count, view.End = s.Count, 0
	if !s.End.IsZero() {
		if i := marketIndexAt(candles, s.End); i+1 < len(candles) {
			view.End = max(1, i+1)
		}
	}
}

// store writes a linked chart's window and crosshair back after it builds.
func (s *ChartSync) store(id string, candles []Candle, view *MarketView, active int) {
	if !s.linked(id) {
		return
	}
	s.Count = view.Count
	s.End = time.Time{}
	if view.End > 0 && view.End <= len(candles) {
		s.End = candles[view.End-1].Time
	}
	switch {
	case active >= 0:
		s.Cursor, s.owner = candles[active].Time, id
	case s.owner == id:
		s.Cursor, s.owner = time.Time{}, ""
	}
}

// cursorIndex is the candle at the linked crosshair, or -1.
func (s *ChartSync) cursorIndex(id string, candles []Candle) int {
	if !s.linked(id) || s.Cursor.IsZero() || s.owner == id {
		return -1
	}
	return marketIndexAt(candles, s.Cursor)
}

// Link builds a toggle that links or unlinks chart id; Changed reports a switch.
func (s *ChartSync) Link(c *ui.Context, id string) ChoiceResult {
	if id == "" {
		panic("mujica: ChartSync.Link needs an id")
	}
	k := Tokens(c)
	on := s.linked(id)
	pressed := on
	b := ui.ToggleBase(c, &pressed).MinHeight(theme.MinHitSize).Padding(0, Density(c).Unit()*2).Gap(Density(c).Unit()).
		Radius(theme.ControlRadius).Label(msg(c, "chart_sync.link"))
	internal.Ring(b, k.Focus, theme.ControlRadius)
	changed := b.Changed() && !b.IsDisabled()
	if changed {
		if s.Unlinked == nil {
			s.Unlinked = map[string]bool{}
		}
		on = pressed
		s.Unlinked[id] = !on
		if s.owner == id {
			s.Cursor, s.owner = time.Time{}, ""
		}
		slog.Debug("mujica: chart_sync link", "chart", id, "linked", on)
	}
	fg, bg, word := k.TextMuted, ui.Transparent, msg(c, "chart_sync.unlinked")
	if on {
		fg, bg, word = k.AccentText, k.Selection, msg(c, "chart_sync.linked")
	}
	if b.Hovered() && !on {
		bg = k.SurfaceHover
	}
	b.Background(bg).TextColor(fg).Border(theme.BorderWidth, fg).Transition(colorFade)
	b.Children(func() {
		mark := "circle"
		if on {
			mark = "check"
		}
		ui.Box(c).Size(14, 14).Draw(func(p *ui.Painter, r ui.Rect) { internal.Mark(p, r, mark, fg) })
		ui.Text(c, word).FontSize(FontSize(c, theme.CaptionSize)).SingleLine()
	})
	return ChoiceResult{b, changed}
}
