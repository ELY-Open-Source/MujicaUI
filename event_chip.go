package mujica

import (
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// EventChipOptions configure an EventChip.
type EventChipOptions struct {
	// Location is the zone the time shows in; required.
	Location *time.Location
	Hour12   bool
	// Inline puts the start time and title on one line.
	Inline   bool
	Selected bool
	Disabled bool
}

// EventChip is a clickable, truncating event block.
func EventChip(c *ui.Context, ev CalendarEvent, opts EventChipOptions) *ui.Element {
	loc := eventLocation("EventChip", opts.Location)
	k := Tokens(c)
	d := Density(c)
	col := eventColor(k, ev.Color)
	when := eventTimeText(c, ev, loc, opts.Hour12)
	title := ev.Title
	if title == "" {
		title = msg(c, "event_chip.untitled")
	}
	b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).MinWidth(theme.MinHitSize).Justify(ui.Start).AlignItems(ui.Start).
		Radius(theme.ControlRadius).Padding(d.Unit(), d.Unit()*2).PaddingX(d.Unit() * 2).Gap(d.Unit()).
		Label(title + ", " + when).Disabled(opts.Disabled).Clip()
	internal.Ring(b, k.Focus, theme.ControlRadius)
	bg := col.Alpha(0.14).Over(k.Surface)
	switch {
	case opts.Selected:
		bg = k.Selection
	case b.Pressed():
		bg = col.Alpha(0.24).Over(k.Surface)
	case b.Hovered():
		bg = col.Alpha(0.2).Over(k.Surface)
	}
	b.Background(bg).Transition(colorFade)
	if !opts.Inline {
		b.Column()
	}
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: 3, H: r.H}, col, 0)
	})
	selected := opts.Selected
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if selected {
			datetimeSelectedMark(p, r, col)
		}
	})
	b.Children(func() {
		small := FontSize(c, theme.CaptionSize)
		if opts.Inline {
			ui.Text(c, eventStartText(c, ev, loc, opts.Hour12)).FontSize(small).TextColor(k.TextMuted).SingleLine().Shrink(0).FontFeatures("tnum")
		}
		ui.Text(c, title).FontSize(small).Bold().TextColor(k.Text).SingleLine().Ellipsis("…").Shrink(1)
		if !opts.Inline {
			ui.Text(c, when).FontSize(small).TextColor(k.TextMuted).SingleLine().Ellipsis("…").FontFeatures("tnum")
		}
	})
	return b
}

// eventStartText is the start time alone, or "All day".
func eventStartText(c *ui.Context, ev CalendarEvent, loc *time.Location, hour12 bool) string {
	if ev.AllDay {
		return msg(c, "event_chip.all_day")
	}
	s := ev.Start.In(loc)
	return timeText(c, TimeOfDay{s.Hour(), s.Minute(), 0}, timeConfig{hour12: hour12})
}
