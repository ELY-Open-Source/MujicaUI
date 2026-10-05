package mujica

import (
	"log/slog"
	"sort"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChartSync shares time windows and crosshairs across linked charts.
type ChartSync struct {
	// Count is linked zoom; zero defaults to 120.
	Count int
	// End is inclusive; zero follows each chart's latest candle.
	End time.Time
	// Start is the inclusive left edge across candle intervals.
	Start time.Time
	// Cursor is the crosshair time; zero shows none.
	Cursor time.Time
	// Unlinked charts keep their own view and crosshair.
	Unlinked   map[string]bool
	owner      string
	rangeOwner string
	loaded     map[string]MarketView
}

func (s *ChartSync) linked(id string) bool { return !s.Unlinked[id] }

// load applies the shared time window to linked charts.
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
	if s.rangeOwner == "" {
		s.rangeOwner = id
	}
	if len(candles) > 0 && (s.Start.IsZero() || s.End.IsZero() && s.rangeOwner == id) {
		from, _ := marketWindow(*view, len(candles), len(candles))
		s.Start = candles[from].Time
	}
	if !s.Start.IsZero() && len(candles) > 0 {
		from := sort.Search(len(candles), func(i int) bool { return !candles[i].Time.Before(s.Start) })
		to := len(candles)
		if view.End > 0 {
			to = view.End
		}
		if from < to {
			view.Count = to - from
		}
	}
	if s.loaded == nil {
		s.loaded = map[string]MarketView{}
	}
	s.loaded[id] = *view
}

// store updates the shared window and active crosshair.
func (s *ChartSync) store(id string, candles []Candle, view *MarketView, active int) {
	if !s.linked(id) {
		return
	}
	loaded, ok := s.loaded[id]
	if !ok || *view != loaded {
		s.Count, s.End, s.rangeOwner = view.Count, time.Time{}, id
		if view.End > 0 && view.End <= len(candles) {
			s.End = candles[view.End-1].Time
		}
		if len(candles) > 0 {
			from, _ := marketWindow(*view, len(candles), len(candles))
			s.Start = candles[from].Time
		}
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

// Link toggles synchronization for a stable chart ID.
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
