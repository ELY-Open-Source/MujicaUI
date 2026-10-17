package datetime

import (
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// EventPopoverOptions configure an EventPopover.
type EventPopoverOptions struct {
	// Location is the zone times show in; required.
	Location *time.Location
	Hour12   bool
	// ReadOnly hides the edit and delete actions.
	ReadOnly bool
	Disabled bool
}

// EventPopoverResult is the popover and the action the user chose.
type EventPopoverResult struct {
	Element ui.Element
	edit    bool
	del     bool
}

// Edit reports a click on Edit.
func (r EventPopoverResult) Edit() bool { return r.edit }

// Delete reports a click on Delete.
func (r EventPopoverResult) Delete() bool { return r.del }

// EventPopover summarizes an event beside anchor while *open.
func EventPopover(c *ui.Context, anchor ui.Element, open *bool, ev CalendarEvent, opts EventPopoverOptions) EventPopoverResult {
	loc := eventLocation("EventPopover", opts.Location)
	k := core.Tokens(c)
	d := core.Density(c)
	var r EventPopoverResult
	if *open && (opts.Disabled || anchor.IsDisabled()) {
		slog.Debug("mujica: event_popover close", "disabled", true)
		*open = false
	}
	r.Element = ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		core.OverlayStyle(c, panel, core.OverlayEnter(c, panel))
		ww, wh := c.Size()
		panel.Margin(d.Unit(), 0, 0, 0).Padding(d.Unit() * 4).Gap(d.Unit() * 2).Width(min(d.ControlHeight()*9, ww-16)).MaxHeight(wh - 16).
			Role(ui.RoleDialog).Label(ev.Title)
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(10, 10).Shrink(0).Draw(func(p *ui.Painter, rr ui.Rect) {
				p.FillPath(internal.Diamond(rr.X+rr.W/2, rr.Y+rr.H/2, rr.W), EventTone(k, ev.Color))
			})
			ui.Text(c, ev.Title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).TextColor(k.Text).Grow(1).Shrink(1)
			if core.IconAction(c, icons.Must("x"), core.Msg(c, "event_popover.close")).Clicked() {
				*open = false
			}
		})
		ui.Scroll(c).Shrink(1).Children(func() {
			ui.Column(c).Gap(d.Unit() * 2).Children(func() {
				for _, line := range eventPopoverLines(c, ev, loc, opts.Hour12) {
					ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Start).Children(func() {
						ui.Icon(c, icons.Must(line[0])).Size(16, 16).TextColor(k.TextMuted)
						ui.Text(c, line[1]).TextColor(k.Text).Grow(1).Shrink(1)
					})
				}
			})
		})
		if opts.ReadOnly {
			return
		}
		core.DatetimeFooter(c, func() {
			if input.Button(c, core.Msg(c, "event_popover.delete"), input.ButtonOptions{Variant: input.Danger}).Clicked() {
				slog.Debug("mujica: event_popover delete", "id", ev.ID)
				r.del, *open = true, false
			}
			if input.Button(c, core.Msg(c, "event_popover.edit"), input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				slog.Debug("mujica: event_popover edit", "id", ev.ID)
				r.edit, *open = true, false
			}
		})
	})
	return r
}

// eventPopoverLines lists icon and text pairs for the summary.
func eventPopoverLines(c *ui.Context, ev CalendarEvent, loc *time.Location, hour12 bool) [][2]string {
	s := ev.Start.In(loc)
	lines := [][2]string{{"clock", core.DatetimeDayLabel(c, core.DateOf(s)) + " · " + eventTimeText(c, ev, loc, hour12)}}
	if ev.Location != "" {
		lines = append(lines, [2]string{"house", ev.Location})
	}
	if ev.Repeat.Unit != RepeatNone {
		lines = append(lines, [2]string{"refresh-cw", RecurrenceSummary(c, ev.Repeat)})
	}
	if len(ev.Attendees) > 0 {
		var names []string
		for _, a := range ev.Attendees {
			names = append(names, a.Name)
		}
		lines = append(lines, [2]string{"user", core.CalendarCount(c, "event_popover.guests", len(names)) + ": " + strings.Join(names, ", ")})
	}
	if len(ev.Reminders) > 0 {
		var rs []string
		for _, d := range ev.Reminders {
			rs = append(rs, ReminderText(c, d))
		}
		lines = append(lines, [2]string{"bell", strings.Join(rs, core.Msg(c, "recurrence_editor.list_sep"))})
	}
	return lines
}
