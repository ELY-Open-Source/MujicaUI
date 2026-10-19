package datetime

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// timeGridConfig is what the week and day views share.
type timeGridConfig struct {
	name    string
	loc     *time.Location
	version int
	slot    int
	today   core.Date
	now     time.Time
	hour12  bool
}

// timeGesture is a drag in progress on the time grid.
type timeGesture struct {
	active bool
	kind   CalendarChangeKind
	id     string
	day    int
	anchor int
	dx, dy float32
	cur    int
}

// timeGridDay is one column's timed and all-day events.
type timeGridDay struct {
	timed  []timePlaced
	allDay []int
}

type timePlaced struct {
	ev        int
	from, to  int
	col, cols int
}

// timeGridSlot checks and defaults the slot length.
func timeGridSlot(name string, slot int) int {
	if slot == 0 {
		return 30
	}
	if slot < 0 || 60%slot != 0 && slot%60 != 0 {
		panic(fmt.Sprintf("mujica: %s SlotMinutes %d must divide an hour", name, slot))
	}
	return slot
}

// timeGrid draws day columns and turns gestures into changes.
func timeGrid(c *ui.Context, root ui.Element, days []core.Date, events []CalendarEvent, cfg timeGridConfig, r *CalendarViewResult) {
	k := core.Tokens(c)
	d := core.Density(c)
	hourH := d.ControlHeight() * 2
	slotH := hourH * float32(cfg.slot) / 60
	g := *ui.Local(root, "time-grid-gesture", func() *timeGesture { return &timeGesture{} })
	morning := func() []int {
		v := make([]int, len(days))
		for i := range v {
			v[i] = 8 * 60 / cfg.slot
		}
		return v
	}
	cursor := ui.Local(root, "time-grid-cursor", morning)
	if len(*cursor) != len(days) {
		*cursor = morning()
	}
	problem := ui.Local(root, "time-grid-problem", func() string { return "" })
	scroll := ui.Local(root, "time-grid-scroll", func() *ui.ScrollState { return &ui.ScrollState{Y: hourH * 8} })
	lay := eventsDerived(root, "time-grid-layout", events, cfg.version, eventDayStart(days[0], cfg.loc),
		eventDayStart(days[len(days)-1].AddDays(1), cfg.loc), func(hits []int) []timeGridDay { return timeGridLayout(events, hits, days, cfg.loc) })
	apply := func(ch CalendarChange, p string) {
		*problem = p
		if p == "" && !root.IsDisabled() {
			slog.Debug("mujica: "+cfg.name+" change", "kind", ch.Kind, "id", ch.ID, "start", ch.Start, "end", ch.End)
			r.change = &ch
		}
	}
	gutter := d.ControlHeight() * 1.6
	ui.Row(c).Children(func() {
		ui.Box(c).Width(gutter).Shrink(0)
		for _, day := range days {
			wd := core.Msg(c, "calendar.weekday_"+strconv.Itoa(int(day.Weekday())))
			b := ui.ButtonBase(c).Grow(1).Basis(0).MinHeight(theme.MinHitSize).Gap(d.Unit()).Radius(theme.ControlRadius).Label(core.DatetimeDayLabel(c, day))
			internal.Ring(b, k.Focus, theme.ControlRadius)
			if day == cfg.today {
				b.Background(k.Selection)
			} else if b.Hovered() {
				b.Background(k.SurfaceHover)
			}
			b.Children(func() {
				ui.Text(c, wd).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				t := ui.Text(c, strconv.Itoa(day.Day)).TextColor(k.Text).FontFeatures("tnum")
				if day == cfg.today {
					t.Bold().TextColor(k.AccentText)
				}
			})
			if b.Clicked() && !b.IsDisabled() {
				r.day = day
			}
		}
	})
	ui.Row(c).BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).Children(func() {
		ui.Text(c, core.Msg(c, "calendar_week_view.all_day")).Width(gutter).Shrink(0).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).
			Padding(d.Unit())
		for i, day := range days {
			cell := ui.Column(c).Grow(1).Basis(0).MinWidth(0).MinHeight(theme.MinHitSize).Padding(d.Unit() / 2).Gap(d.Unit() / 2)
			if id, ok := ui.Drop[monthDrag](cell); ok && !cell.IsDisabled() {
				for _, ev := range events {
					if ev.ID == string(id) {
						apply(eventShift(c, ev, monthDays(ev.Start.In(cfg.loc), day), cfg.loc))
					}
				}
			}
			cell.Children(func() {
				for _, h := range lay[i].allDay {
					monthChip(c, events[h], cfg.loc, CalendarMonthViewOptions{Hour12: cfg.hour12}, r, apply)
				}
			})
		}
	})
	ui.Scroll(c).TrackScroll(*scroll).Height(hourH * 9).Children(func() {
		ui.Row(c).Children(func() {
			ui.Column(c).Width(gutter).Shrink(0).Children(func() {
				for h := range 24 {
					ui.Text(c, timeText(c, TimeOfDay{h, 0, 0}, timeConfig{hour12: cfg.hour12})).Height(hourH).
						FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).FontFeatures("tnum").PaddingX(d.Unit())
				}
			})
			for i, day := range days {
				timeColumn(c, i, day, lay[i], events, cfg, hourH, slotH, g, &(*cursor)[i], r, apply)
			}
		})
	})
	core.DatetimeError(c, *problem)
}

// timeGridLayout places each day's events in minutes with side-by-side columns.
func timeGridLayout(events []CalendarEvent, hits []int, days []core.Date, loc *time.Location) []timeGridDay {
	out := make([]timeGridDay, len(days))
	for i, day := range days {
		from, to := eventDayStart(day, loc), eventDayStart(day.AddDays(1), loc)
		var spans []eventSpan
		for _, h := range hits {
			ev := events[h]
			if !ev.Start.Before(to) || !ev.End.After(from) {
				continue
			}
			if ev.AllDay || ev.End.Sub(ev.Start) >= 24*time.Hour {
				out[i].allDay = append(out[i].allDay, h)
				continue
			}
			a, b := 0, 24*60
			if !ev.Start.Before(from) {
				a = timeWallMinutes(ev.Start.In(loc))
			}
			if ev.End.Before(to) {
				b = timeWallMinutes(ev.End.In(loc))
			}
			b = max(b, a+1)
			out[i].timed = append(out[i].timed, timePlaced{ev: h, from: a, to: b})
			spans = append(spans, eventSpan{float64(a), float64(b)})
		}
		col, cols := eventColumns(spans)
		for j := range out[i].timed {
			out[i].timed[j].col, out[i].timed[j].cols = col[j], cols[j]
		}
	}
	return out
}

func timeWallMinutes(t time.Time) int { return t.Hour()*60 + t.Minute() }

// timeColumn draws one day: lines, events, preview and cursor.
func timeColumn(c *ui.Context, i int, day core.Date, lay timeGridDay, events []CalendarEvent, cfg timeGridConfig, hourH, slotH float32,
	g *timeGesture, cursor *int, r *CalendarViewResult, apply func(CalendarChange, string)) {
	k := core.Tokens(c)
	slots := 24 * 60 / cfg.slot
	col := ui.Box(c).Grow(1).Basis(0).MinWidth(0).Height(hourH*24).Focusable().
		Role(ui.RoleGroup).Label(fmt.Sprintf(core.Msg(c, "calendar_week_view.column"), core.DatetimeDayLabel(c, day))).BorderWidth(0, 0, 0, theme.BorderWidth).BorderColor(k.Border)
	disabled := col.IsDisabled()
	slotOf := func(y float32) int { return max(0, min(slots-1, int(y/slotH))) }
	if !disabled {
		switch {
		case col.Shortcut(0, ui.KeyUp):
			*cursor = max(0, *cursor-1)
		case col.Shortcut(0, ui.KeyDown):
			*cursor = min(slots-1, *cursor+1)
		case col.Shortcut(0, ui.KeyEnter):
			apply(timeCreate(c, day, *cursor, *cursor+1, cfg))
		}
		_, py, _ := col.PointerPosition()
		dx, dy, pressing := col.Dragged()
		if pressing && !g.active {
			*g = timeGesture{active: true, kind: CalendarCreate, day: i, anchor: slotOf(py - dy), cur: slotOf(py)}
		}
		if g.active && g.kind == CalendarCreate && g.day == i {
			if pressing {
				g.dx, g.dy = g.dx+dx, g.dy+dy
				g.cur = slotOf(py)
			} else {
				a, b := min(g.anchor, g.cur), max(g.anchor, g.cur)+1
				if math.Abs(float64(g.dy)) >= 4 {
					apply(timeCreate(c, day, a, b, cfg))
				}
				*cursor = g.anchor
				*g = timeGesture{}
			}
		}
	}
	focused := col.Focused()
	cur := *cursor
	creating := g.active && g.kind == CalendarCreate && g.day == i
	ga, gb := min(g.anchor, g.cur), max(g.anchor, g.cur)+1
	col.Draw(func(p *ui.Painter, rr ui.Rect) {
		for h := 1; h < 24; h++ {
			internal.Rule(p, rr.X, rr.X+rr.W, rr.Y+float32(h)*hourH, k.Border)
		}
		if focused {
			box := ui.Rect{X: rr.X + 2, Y: rr.Y + float32(cur)*slotH + 1, W: rr.W - 4, H: slotH - 2}
			p.Fill(box, k.Selection, theme.ControlRadius)
			core.DatetimeSelectedMark(p, box, k.Focus)
		}
		if creating {
			box := ui.Rect{X: rr.X + 2, Y: rr.Y + float32(ga)*slotH, W: rr.W - 4, H: float32(gb-ga) * slotH}
			p.Fill(box, k.Selection, theme.ControlRadius)
			p.StrokeDashed(box, k.AccentText, theme.ControlRadius, theme.BorderWidth)
		}
	})
	col.Children(func() {
		now := cfg.now
		if now.IsZero() {
			now = c.Now()
		}
		for _, pl := range lay.timed {
			timeChip(c, i, pl, events[pl.ev], cfg, hourH, slotH, g, r, apply)
		}
		if day == core.DateOf(now.In(cfg.loc)) {
			CurrentTimeIndicator(c, CurrentTimeIndicatorOptions{Location: cfg.loc, Day: day, Now: cfg.now, Hour12: cfg.hour12})
		}
	})
}

// timeChip places an event; drag moves, bottom edge resizes.
func timeChip(c *ui.Context, i int, pl timePlaced, ev CalendarEvent, cfg timeGridConfig, hourH, slotH float32,
	g *timeGesture, r *CalendarViewResult, apply func(CalendarChange, string)) {
	from, to := float32(pl.from)*hourH/60, float32(pl.to)*hourH/60
	mine := g.active && g.id == ev.ID && g.day == i
	snap := func(v float32) int { return int(math.Round(float64(v/slotH))) * cfg.slot }
	var shiftY, growY float32
	if mine && g.kind == CalendarMove {
		shiftY = float32(snap(g.dy)) * hourH / 60
	}
	if mine && g.kind == CalendarResize {
		growY = float32(snap(g.dy)) * hourH / 60
	}
	box := ui.Box(c.Key(ev.ID)).Absolute().Top(from + shiftY).Height(max(to-from+growY, theme.MinHitSize)).
		LeftPercent(100 * float32(pl.col) / float32(pl.cols)).WidthPercent(100 / float32(pl.cols)).Padding(1)
	box.Children(func() {
		chip := EventChip(c, ev, EventChipOptions{Location: cfg.loc, Hour12: cfg.hour12, Selected: mine, Inline: to-from < hourH*0.75}).Fill()
		handle := ui.Box(c).Absolute().Left(0).Right(0).Bottom(0).Height(6).Cursor(ui.CursorResizeNS).Label(fmt.Sprintf(core.Msg(c, "calendar_week_view.resize"), ev.Title))
		if chip.IsDisabled() {
			return
		}
		if chip.Clicked() && !mine {
			r.opened = ev.ID
		}
		for key, dm := range map[ui.Key]int{ui.KeyUp: -cfg.slot, ui.KeyDown: cfg.slot} {
			if chip.Shortcut(ui.Alt, key) {
				apply(timeMove(c, ev, 0, dm, cfg))
			}
			if chip.Shortcut(ui.Alt|ui.Shift, key) {
				apply(timeResize(c, ev, dm, cfg))
			}
		}
		for key, dd := range map[ui.Key]int{ui.KeyLeft: -1, ui.KeyRight: 1} {
			if chip.Shortcut(ui.Alt, key) {
				apply(timeMove(c, ev, dd, 0, cfg))
			}
		}
		for _, part := range []struct {
			e    ui.Element
			kind CalendarChangeKind
		}{{chip, CalendarMove}, {handle, CalendarResize}} {
			dx, dy, ok := part.e.Dragged()
			switch {
			case ok && (!g.active && (dx != 0 || dy != 0) || mine && g.kind == part.kind):
				if !g.active {
					*g = timeGesture{active: true, kind: part.kind, id: ev.ID, day: i}
				}
				g.dx, g.dy = g.dx+dx, g.dy+dy
			case !ok && mine && g.kind == part.kind:
				dm := snap(g.dy)
				if part.kind == CalendarResize {
					if dm != 0 {
						apply(timeResize(c, ev, dm, cfg))
					}
				} else if dd := int(math.Round(float64(g.dx / max(chip.Bounds().W*float32(pl.cols), 1)))); dm != 0 || dd != 0 {
					apply(timeMove(c, ev, dd, dm, cfg))
				}
				*g = timeGesture{}
			}
		}
	})
}

// timeCreate makes an event over slots [a, b).
func timeCreate(c *ui.Context, day core.Date, a, b int, cfg timeGridConfig) (CalendarChange, string) {
	s, sn := timeAt(day, a*cfg.slot, cfg.loc)
	e, en := timeAt(day, b*cfg.slot, cfg.loc)
	if p := eventWallProblem(c, timeDay(day, a*cfg.slot), timeWall(a*cfg.slot), cfg.loc, sn); p != "" {
		return CalendarChange{}, p
	}
	if p := eventWallProblem(c, timeDay(day, b*cfg.slot), timeWall(b*cfg.slot), cfg.loc, en); p != "" {
		return CalendarChange{}, p
	}
	return CalendarChange{Kind: CalendarCreate, Start: s, End: e}, ""
}

// timeMove shifts ev by days and minutes of wall time.
func timeMove(c *ui.Context, ev CalendarEvent, days, minutes int, cfg timeGridConfig) (CalendarChange, string) {
	s := ev.Start.In(cfg.loc)
	day := core.DateOf(s).AddDays(days)
	m := timeWallMinutes(s) + minutes
	at, n := timeAt(day, m, cfg.loc)
	if p := eventWallProblem(c, timeDay(day, m), timeWall(m), cfg.loc, n); p != "" {
		return CalendarChange{}, p
	}
	return CalendarChange{Kind: CalendarMove, ID: ev.ID, Start: at, End: at.Add(ev.End.Sub(ev.Start))}, ""
}

// timeResize moves ev's end by minutes of wall time.
func timeResize(c *ui.Context, ev CalendarEvent, minutes int, cfg timeGridConfig) (CalendarChange, string) {
	e := ev.End.In(cfg.loc)
	day, m := core.DateOf(e), timeWallMinutes(e)+minutes
	at, n := timeAt(day, m, cfg.loc)
	if p := eventWallProblem(c, timeDay(day, m), timeWall(m), cfg.loc, n); p != "" {
		return CalendarChange{}, p
	}
	if !at.After(ev.Start) {
		return CalendarChange{}, core.Msg(c, "calendar_week_view.too_short")
	}
	return CalendarChange{Kind: CalendarResize, ID: ev.ID, Start: ev.Start, End: at}, ""
}

// timeAt resolves minute m of day as wall time.
func timeAt(day core.Date, m int, loc *time.Location) (time.Time, int) {
	return dateTimeResolve(timeDay(day, m), timeWall(m), loc)
}

func timeDay(day core.Date, m int) core.Date {
	return day.AddDays(int(math.Floor(float64(m) / (24 * 60))))
}

func timeWall(m int) TimeOfDay {
	m = ((m % (24 * 60)) + 24*60) % (24 * 60)
	return TimeOfDay{m / 60, m % 60, 0}
}
