package messaging

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/media"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// MailItem supplies a unique ID and mail preview data.
type MailItem struct {
	ID, From, Subject, Snippet string
	Time                       time.Time
	Unread, Starred            bool
	Attachments                int
}

// MailListOptions configure a MailList.
type MailListOptions struct {
	// Height is the list's height; zero means 420.
	Height   float32
	Disabled bool
}

// MailListResult is a MailList's element and this frame's events.
type MailListResult struct {
	Element ui.Element
	changed bool
	// Opened contains the ID activated by Enter or double-click.
	Opened string
	// Starred is the ID whose star the user toggled.
	Starred string
}

// Changed reports a new selection; s.Selection holds the chosen IDs.
func (r MailListResult) Changed() bool { return r.changed }

// MailList virtualizes previews, selection and star actions.
func MailList(c *ui.Context, s *data.ListState[string], mails []MailItem, opts MailListOptions) MailListResult {
	if opts.Height < 0 {
		panic("mujica: MailList Height must not be negative")
	}
	h := opts.Height
	if h == 0 {
		h = 420
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	var r MailListResult
	box := ui.Column(c).Height(h).Disabled(opts.Disabled)
	r.Element = box
	box.Children(func() {
		v := data.List(c, s, len(mails), data.ListOptions[string]{
			Key:      func(i int) string { return mails[i].ID },
			Label:    func(i int) string { return mailLabel(c, mails[i]) },
			Multiple: true,
			Empty:    core.Msg(c, "mail_list.empty"),
		}, func(i int) {
			m := mails[i]
			if m.ID == "" {
				panic("mujica: MailList item without an ID")
			}
			ui.Box(c).Width(8).Shrink(0).Children(func() {
				if m.Unread {
					ui.Box(c).Size(8, 8).Radius(4).Background(k.AccentText)
				}
			})
			star := media.MediaToggleAction(c, icons.Must("star"), fmt.Sprintf(core.Msg(c, "mail_list.star"), m.Subject), m.Starred)
			if m.Starred {
				star.TextColor(k.Warning)
			}
			if star.Clicked() && !star.IsDisabled() {
				r.Starred = m.ID
			}
			ui.Column(c).Grow(1).MinWidth(0).Gap(2).Padding(u, 0).Children(func() {
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					from := ui.Text(c, m.From).SingleLine().Grow(1).MinWidth(0)
					if m.Unread {
						from.FontWeight(700)
					}
					if m.Attachments > 0 {
						ui.Icon(c, msgPaperclip).Size(13, 13).TextColor(k.TextMuted).Label(fmt.Sprintf(core.Msg(c, "mail_list.attachments"), m.Attachments))
					}
					if !m.Time.IsZero() {
						ui.Text(c, mailTime(c, m.Time)).FontSize(small).TextColor(k.TextMuted).FontFeatures("tnum").SingleLine()
					}
				})
				subj := ui.Text(c, m.Subject).SingleLine()
				if m.Unread {
					subj.FontWeight(700)
				}
				ui.Text(c, m.Snippet).FontSize(small).TextColor(k.TextMuted).SingleLine()
			})
		})
		v.Element.Grow(1)
		r.changed = v.Changed()
		if v.Element.IsDisabled() {
			return
		}
		if id, ok := s.Selected(); ok {
			if v.Submitted() {
				r.Opened = id
			}
			if v.Element.Shortcut(0, ui.KeyS) {
				r.Starred = id
			}
		}
	})
	if r.Opened != "" || r.Starred != "" {
		slog.Debug("mujica: mail_list action", "opened", r.Opened, "starred", r.Starred)
	}
	return r
}

// mailLabel reads a row: unread, starred, from, subject.
func mailLabel(c *ui.Context, m MailItem) string {
	s := m.From + ", " + m.Subject
	if m.Starred {
		s = core.Msg(c, "mail_list.starred") + ", " + s
	}
	if m.Unread {
		s = core.Msg(c, "mail_list.unread") + ", " + s
	}
	return s
}

// mailTime shows today's mails by time, older ones by date.
func mailTime(c *ui.Context, t time.Time) string {
	now := c.Now().In(t.Location())
	if y, m, d := t.Date(); y == now.Year() && m == now.Month() && d == now.Day() {
		return t.Format("15:04")
	}
	return t.Format(core.Msg(c, "mail_list.date"))
}
