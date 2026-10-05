package feedback

import (
	"cmp"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/display"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Activity is one event of an ActivityFeed: Actor did Action to Object at Time.
type Activity struct {
	ID     string
	Actor  string
	Action string
	Object string
	Time   time.Time
}

// ActivityFeedOptions configure an ActivityFeed.
type ActivityFeedOptions struct {
	// Location decides which day an event falls on; nil is time.Local.
	Location *time.Location
	// Today fixes the day called Today; zero uses the clock.
	Today core.Date
}

// activityRow is a day header (item -1) or an event.
type activityRow struct {
	day  core.Date
	item int
}

// activityFeedCache keeps the grouped rows until the count, first or last event, or location changes.
type activityFeedCache struct {
	n           int
	first, last string
	loc         *time.Location
	rows        []activityRow
	scroll      ui.ScrollState
	sized       bool
}

// ActivityFeed lists events grouped by day, newest first as given, building only the rows in view.
func ActivityFeed(c *ui.Context, items []Activity, opts ActivityFeedOptions) *ui.Element {
	loc := cmp.Or(opts.Location, time.Local)
	for _, a := range items {
		if a.ID == "" || a.Actor == "" || a.Action == "" || a.Time.IsZero() {
			panic("mujica: ActivityFeed event needs an ID, Actor, Action and Time")
		}
	}
	k := core.Tokens(c)
	d := core.Density(c)
	sc := ui.Scroll(c).MinHeight(0).Role(ui.RoleList).Label(core.Msg(c, "activity_feed.label"))
	if len(items) == 0 {
		sc.Children(func() {
			EmptyState(c, EmptyStateOptions{Title: core.Msg(c, "activity_feed.empty"), Description: core.Msg(c, "activity_feed.empty_hint"), Icon: icons.Must("clock")})
		})
		return sc
	}
	cache := *ui.Local(sc, "activity_feed", func() *activityFeedCache { return &activityFeedCache{} })
	sc.TrackScroll(&cache.scroll)
	if cache.n != len(items) || cache.first != items[0].ID || cache.last != items[len(items)-1].ID || cache.loc != loc {
		cache.n, cache.first, cache.last, cache.loc = len(items), items[0].ID, items[len(items)-1].ID, loc
		cache.rows = cache.rows[:0]
		var prev core.Date
		for i, a := range items {
			if day := core.DateOf(a.Time.In(loc)); day != prev || i == 0 {
				cache.rows = append(cache.rows, activityRow{day: day, item: -1})
				prev = day
			}
			cache.rows = append(cache.rows, activityRow{item: i})
		}
		slog.Debug("mujica: activity_feed grouped", "events", len(items), "rows", len(cache.rows))
	}
	rowH := d.ControlHeight() + 2*d.Unit()
	view := sc.Bounds()
	if view.H == 0 && !cache.sized {
		// Unmeasured: build again once the box has its height.
		cache.sized = true
		c.Invalidate()
	}
	first := max(0, int(cache.scroll.Y/rowH)-2)
	last := min(len(cache.rows), int((cache.scroll.Y+max(view.H, rowH))/rowH)+3)
	today := core.DatetimeToday(c, opts.Today)
	sc.Children(func() {
		ui.Box(c).FillWidth().Height(rowH * float32(len(cache.rows))).Children(func() {
			for r := first; r < last; r++ {
				row := cache.rows[r]
				box := ui.Row(c).Key(r).Absolute().Left(0).Right(0).Top(rowH*float32(r)).Height(rowH).
					AlignItems(ui.Center).Gap(d.Unit()*3).Padding(0, d.Unit()*3)
				if row.item < 0 {
					activityFeedHeader(c, k, box, row.day, today)
					continue
				}
				activityFeedEvent(c, k, box, items[row.item], loc)
			}
		})
	})
	return sc
}

// activityFeedHeader is a day's title: Today, Yesterday or the date, over an ornament rule.
func activityFeedHeader(c *ui.Context, k theme.Tokens, box *ui.Element, day, today core.Date) {
	label := core.DatetimeDayLabel(c, day)
	switch day {
	case today:
		label = core.Msg(c, "activity_feed.today")
	case today.AddDays(-1):
		label = core.Msg(c, "activity_feed.yesterday")
	}
	box.Role(ui.RoleText).AlignItems(ui.End).Padding(0, core.Density(c).Unit()*3, core.Density(c).Unit(), core.Density(c).Unit()*3)
	box.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, k.Ornament, 0) })
	box.Children(func() {
		ui.Text(c, label).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.CaptionSize+1)).FontWeight(600).TextColor(k.Text).SingleLine()
	})
}

// activityFeedEvent is one event: avatar, "Actor action Object", and the time.
func activityFeedEvent(c *ui.Context, k theme.Tokens, box *ui.Element, a Activity, loc *time.Location) {
	when := a.Time.In(loc).Format("15:04")
	box.Role(ui.RoleListItem).Label(a.Actor + " " + a.Action + " " + a.Object + ", " + when)
	box.Children(func() {
		display.Avatar(c, a.Actor, display.AvatarOptions{Size: 28})
		ui.Row(c).Grow(1).Shrink(1).MinWidth(0).Gap(4).Children(func() {
			ui.Text(c, a.Actor).Bold().TextColor(k.Text).SingleLine().Shrink(0)
			ui.Text(c, a.Action).TextColor(k.TextMuted).SingleLine().Shrink(0)
			if a.Object != "" {
				ui.Text(c, a.Object).TextColor(k.AccentText).SingleLine().Shrink(1).Ellipsis("…")
			}
		})
		ui.Text(c, when).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").SingleLine()
	})
}
