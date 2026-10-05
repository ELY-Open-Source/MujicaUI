package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// AttendeeStatus is an attendee's reply.
type AttendeeStatus int

const (
	AttendeePending AttendeeStatus = iota
	AttendeeAccepted
	AttendeeTentative
	AttendeeDeclined
)

// Attendee is one person invited to an event.
type Attendee struct {
	Name      string
	Email     string
	Status    AttendeeStatus
	Organizer bool
	Optional  bool
}

// AttendeeListOptions configure an AttendeeList.
type AttendeeListOptions struct {
	// Removable shows a remove button on everyone but the organizer.
	Removable bool
	Disabled  bool
	Label     string
}

// AttendeeListResult is the list and the attendee the user removed.
type AttendeeListResult struct {
	Element *ui.Element
	removed int
}

// Removed returns the index of the attendee the user removed.
func (r AttendeeListResult) Removed() (int, bool) { return r.removed, r.removed >= 0 }

// AttendeeList shows invitees, their replies and a summary.
func AttendeeList(c *ui.Context, attendees []Attendee, opts AttendeeListOptions) AttendeeListResult {
	k := Tokens(c)
	d := Density(c)
	var counts [4]int
	for _, a := range attendees {
		if a.Status < AttendeePending || a.Status > AttendeeDeclined {
			panic(fmt.Sprintf("mujica: AttendeeList status %d out of range", a.Status))
		}
		counts[a.Status]++
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "attendee_list.label")
	}
	r := AttendeeListResult{removed: -1}
	r.Element = ui.Column(c).Gap(d.Unit()).Role(ui.RoleList).Label(label).Disabled(opts.Disabled)
	r.Element.Children(func() {
		summary := fmt.Sprintf(msg(c, "attendee_list.summary"), counts[AttendeeAccepted], counts[AttendeeTentative],
			counts[AttendeeDeclined], counts[AttendeePending])
		ui.Text(c, summary).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Role(ui.RoleStatus)
		if len(attendees) == 0 {
			ui.Text(c, msg(c, "attendee_list.empty")).TextColor(k.TextMuted)
		}
		for i, a := range attendees {
			ui.Row(c).Key(i).MinHeight(d.ControlHeight()).Gap(d.Unit() * 2).AlignItems(ui.Center).Role(ui.RoleListItem).Children(func() {
				Avatar(c, a.Name, AvatarOptions{Size: 28})
				ui.Column(c).Grow(1).Shrink(1).Children(func() {
					name := a.Name
					if a.Organizer {
						name += " · " + msg(c, "attendee_list.organizer")
					}
					if a.Optional {
						name += " · " + msg(c, "attendee_list.optional")
					}
					ui.Text(c, name).TextColor(k.Text).SingleLine().Ellipsis("…")
					if a.Email != "" {
						ui.Text(c, a.Email).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine().Ellipsis("…")
					}
				})
				icon, fg := attendeeStatusLook(k, a.Status)
				ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Shrink(0).TextColor(fg).Children(func() {
					ui.Icon(c, icons.Must(icon)).Size(16, 16)
					ui.Text(c, msg(c, attendeeStatusKeys[a.Status])).FontSize(FontSize(c, theme.CaptionSize))
				})
				if opts.Removable && !a.Organizer {
					b := iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "attendee_list.remove"), a.Name))
					if b.Clicked() && !b.IsDisabled() {
						slog.Debug("mujica: attendee_list remove", "index", i)
						r.removed = i
					}
				}
			})
		}
	})
	return r
}

var attendeeStatusKeys = [...]string{"attendee_list.pending", "attendee_list.accepted", "attendee_list.tentative", "attendee_list.declined"}

func attendeeStatusLook(k theme.Tokens, s AttendeeStatus) (icon string, fg ui.Color) {
	switch s {
	case AttendeeAccepted:
		return "circle-check", k.Success
	case AttendeeTentative:
		return "circle-alert", k.Warning
	case AttendeeDeclined:
		return "circle-x", k.Danger
	}
	return "clock", k.TextMuted
}
