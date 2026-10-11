package chat

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MessageListState keeps a MessageList's place and what the user has seen.
type MessageListState struct {
	list    ui.ListState
	seen    int
	started bool
}

// ScrollToEnd shows the newest message and follows new ones again.
func (s *MessageListState) ScrollToEnd() { s.list.ScrollToEnd() }

// AtEnd reports whether the newest message is in view.
func (s *MessageListState) AtEnd() bool { return s.list.AtEnd() }

// MessageListOptions configure a MessageList.
type MessageListOptions struct {
	// ID returns message i's stable identity; required.
	ID func(i int) string
	// Date returns message i's time; set, a separator heads each new day.
	Date func(i int) time.Time
	// Label returns message i's text for assistive technology and typing to find.
	Label func(i int) string
}

// MessageListResult is a MessageList's element and its unread count.
type MessageListResult struct {
	Element *ui.Element
	unread  int
}

// Unread returns how many messages arrived since the user scrolled away from the end.
func (r MessageListResult) Unread() int { return r.unread }

// MessageList shows n messages, building only those in view, and follows the end until the user scrolls up.
func MessageList(c *ui.Context, s *MessageListState, n int, opts MessageListOptions, row func(i int)) MessageListResult {
	if opts.ID == nil || row == nil {
		panic("mujica: MessageList needs ID and row")
	}
	if n < 0 {
		panic("mujica: MessageList with negative message count")
	}
	u := core.Density(c).Unit()
	s.list.Key = func(i int) any { return opts.ID(i) }
	s.list.Label = opts.Label
	s.list.FollowEnd = true
	at := !s.started || s.list.AtEnd()
	if at || n < s.seen {
		s.seen = n
	}
	s.started = true
	r := MessageListResult{unread: n - s.seen}
	r.Element = ui.Box(c).Grow(1).MinHeight(0)
	r.Element.Children(func() {
		ui.List(c, &s.list, n, func(i int) {
			ui.Column(c).Padding(u*2, u*4).Gap(u * 2).Children(func() {
				if opts.Date != nil && (i == 0 || !messageListSameDay(opts.Date(i-1), opts.Date(i))) {
					messageListDay(c, opts.Date(i))
				}
				row(i)
			})
		}).Fill().Justify(ui.End).Label(core.Msg(c, "message_list.label"))
		jump := ScrollToBottomButton(c, !at, r.unread)
		if jump.Clicked() && !r.Element.IsDisabled() {
			s.list.ScrollToEnd()
			s.seen = n
			slog.Debug("mujica: message_list jump", "unread", r.unread)
		}
	})
	return r
}

func messageListSameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// messageListDay is a separator naming the day: today, yesterday or the date.
func messageListDay(c *ui.Context, t time.Time) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := t.Format(core.Msg(c, "message_list.date"))
	now := c.Now().In(t.Location())
	switch {
	case messageListSameDay(t, now):
		label = core.Msg(c, "message_list.today")
	case messageListSameDay(t, now.AddDate(0, 0, -1)):
		label = core.Msg(c, "message_list.yesterday")
	}
	ui.Row(c).Gap(u*3).AlignItems(ui.Center).Padding(u*2, 0).Role(ui.RoleText).Label(label).Children(func() {
		rule := func() {
			ui.Box(c).Grow(1).Height(theme.BorderWidth).Background(k.Border)
		}
		rule()
		ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
		rule()
	})
}
