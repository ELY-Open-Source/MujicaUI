package messaging

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/display"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Member supplies a unique ID, name, presence and role.
type Member struct {
	ID, Name string
	Image    *ui.Bitmap
	Presence PresenceStatus
	// Role, as "Owner"; empty shows none.
	Role string
}

// MemberListOptions configure a MemberList.
type MemberListOptions struct {
	// Height is the list's height; zero means 360.
	Height   float32
	Disabled bool
}

// MemberListState keeps a MemberList's place and choice.
type MemberListState struct {
	rows   ui.ListState
	cursor int
	// Selected is the chosen member's ID; empty is none.
	Selected  string
	cache     memberRows
	looked    string
	lookedGen int
}

type memberRows struct {
	first   *Member
	n       int
	rows    []int
	online  int
	offline int
	gen     int
}

// MemberListResult is a MemberList's element and this frame's events.
type MemberListResult struct {
	Element ui.Element
	changed bool
	// Opened contains the ID activated by Enter or double-click.
	Opened string
}

// Changed reports a new selection.
func (r MemberListResult) Changed() bool { return r.changed }

// MemberList virtualizes online and offline groups.
func MemberList(c *ui.Context, s *MemberListState, members []Member, opts MemberListOptions) MemberListResult {
	if opts.Height < 0 {
		panic("mujica: MemberList Height must not be negative")
	}
	h := opts.Height
	if h == 0 {
		h = 360
	}
	memberGroup(s, members)
	k := core.Tokens(c)
	d := core.Density(c)
	rows := s.cache.rows
	if s.looked != s.Selected || s.lookedGen != s.cache.gen {
		s.looked, s.lookedGen, s.cursor = s.Selected, s.cache.gen, -1
		for i, r := range rows {
			if s.Selected != "" && r >= 0 && members[r].ID == s.Selected {
				s.cursor = i
				break
			}
		}
	}
	s.rows.Selected = &s.cursor
	s.rows.Header = func(i int) bool { return rows[i] < 0 }
	s.rows.Key = func(i int) any {
		if rows[i] < 0 {
			return rows[i]
		}
		return members[rows[i]].ID
	}
	s.rows.Label = func(i int) string {
		if rows[i] < 0 {
			return memberHeader(c, s, rows[i])
		}
		return members[rows[i]].Name
	}
	var r MemberListResult
	box := ui.Column(c).Height(h).Disabled(opts.Disabled)
	r.Element = box
	box.Children(func() {
		if len(rows) == 0 {
			ui.Text(c, core.Msg(c, "member_list.empty")).TextColor(k.TextMuted)
			return
		}
		restore := data.DataSelectTheme(c, 0)
		l := ui.List(c, &s.rows, len(rows), func(i int) {
			if rows[i] < 0 {
				ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).AlignItems(ui.Center).Background(k.Surface).Children(func() {
					ui.Text(c, memberHeader(c, s, rows[i])).FontSize(core.FontSize(c, theme.CaptionSize)).FontWeight(600).TextColor(k.TextMuted)
				})
				return
			}
			m := members[rows[i]]
			ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Shrink(0).Children(func() {
					display.Avatar(c, m.Name, display.AvatarOptions{Image: m.Image, Size: 28})
					ui.Box(c).Absolute().Right(-2).Bottom(-2).Padding(1.5).Radius(6).Background(k.Surface).Children(func() {
						OnlineStatus(c, m.Presence, OnlineStatusOptions{HideText: true, Size: 9})
					})
				})
				t := ui.Text(c, m.Name).SingleLine().Grow(1).MinWidth(0)
				if m.Presence == PresenceOffline {
					t.TextColor(k.TextMuted)
				}
				if m.Role != "" {
					display.Badge(c, display.BadgeOptions{Text: m.Role, Tone: display.BadgeNeutral}, nil)
				}
				if i == s.cursor {
					ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(k.AccentText)
				}
			})
		})
		restore()
		l.Grow(1).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius).Padding(d.Unit())
		if l.IsDisabled() {
			return
		}
		if l.Changed() && s.cursor >= 0 && rows[s.cursor] >= 0 {
			s.Selected, r.changed = members[rows[s.cursor]].ID, true
			s.looked = s.Selected
		}
		if l.Submitted() && s.cursor >= 0 && rows[s.cursor] >= 0 {
			r.Opened = members[rows[s.cursor]].ID
		}
	})
	if r.changed || r.Opened != "" {
		slog.Debug("mujica: member_list choose", "id", s.Selected, "opened", r.Opened)
	}
	return r
}

func memberHeader(c *ui.Context, s *MemberListState, row int) string {
	if row == -1 {
		return fmt.Sprintf(core.Msg(c, "member_list.online"), s.cache.online)
	}
	return fmt.Sprintf(core.Msg(c, "member_list.offline"), s.cache.offline)
}

// memberGroup rebuilds the sections only when members changed.
func memberGroup(s *MemberListState, members []Member) {
	var first *Member
	if len(members) > 0 {
		first = &members[0]
	}
	if s.cache.rows != nil && s.cache.first == first && s.cache.n == len(members) {
		return
	}
	seen := make(map[string]bool, len(members))
	on, off := []int{-1}, []int{-2}
	for i, m := range members {
		if m.ID == "" || m.Name == "" || seen[m.ID] {
			panic("mujica: MemberList members need unique IDs and names: " + m.ID)
		}
		seen[m.ID] = true
		m.Presence.key()
		if m.Presence == PresenceOffline {
			off = append(off, i)
		} else {
			on = append(on, i)
		}
	}
	s.cache = memberRows{first: first, n: len(members), online: len(on) - 1, offline: len(off) - 1, gen: s.cache.gen + 1}
	if len(on) > 1 {
		s.cache.rows = append(s.cache.rows, on...)
	}
	if len(off) > 1 {
		s.cache.rows = append(s.cache.rows, off...)
	}
	if s.cache.rows == nil {
		s.cache.rows = []int{}
	}
	slog.Debug("mujica: member_list grouped", "online", s.cache.online, "offline", s.cache.offline)
}
