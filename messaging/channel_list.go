package messaging

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/display"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Channel supplies channel or direct-conversation data.
type Channel struct {
	ID, Name string
	// Direct marks a private conversation, shown with the person's presence.
	Direct   bool
	Presence PresenceStatus
	Unread   int
	Muted    bool
	Pinned   bool
}

// ChannelGroup is a titled, collapsible section of channels.
type ChannelGroup struct {
	ID, Title string
	Channels  []Channel
}

// ChannelListOptions configure a ChannelList.
type ChannelListOptions struct {
	// Height defaults to 360 DIP.
	Height float32
	// Revision invalidates grouped rows after in-place data edits.
	Revision int
	Disabled bool
}

type channelListKey struct{}

type channelIndex struct{ group, item int }

type channelListState struct {
	list            ui.ListState
	folded          map[string]bool
	first           *ChannelGroup
	count, revision int
	rows            []channelIndex
}

// ChannelList chooses a channel from collapsible, virtualized groups.
func ChannelList(c *ui.Context, selected *string, groups []ChannelGroup, opts ChannelListOptions) core.ChoiceResult {
	if opts.Height < 0 {
		panic("mujica: ChannelList Height must be nonnegative")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	col := ui.Column(c).Disabled(opts.Disabled)
	st := ui.Local(col, channelListKey{}, func() channelListState {
		return channelListState{folded: map[string]bool{}}
	})
	var first *ChannelGroup
	if len(groups) > 0 {
		first = &groups[0]
	}
	if st.rows == nil || st.first != first || st.count != len(groups) || st.revision != opts.Revision {
		st.first, st.count, st.revision = first, len(groups), opts.Revision
		channelRows(st, groups)
	}
	st.list.Key = func(i int) any {
		r := st.rows[i]
		if r.item < 0 {
			return "group:" + groups[r.group].ID
		}
		return "channel:" + groups[r.group].Channels[r.item].ID
	}
	h := opts.Height
	if h == 0 {
		h = 360
	}
	h = min(h, float32(len(st.rows))*d.ControlHeight())
	old := *selected
	col.Children(func() {
		ui.List(c, &st.list, len(st.rows), func(i int) {
			r := st.rows[i]
			g := groups[r.group]
			if r.item < 0 {
				open := !st.folded[g.ID]
				if channelHeader(c, k, g, open) && !col.IsDisabled() {
					st.folded[g.ID] = open
					st.first = nil
					st.count = -1
					c.Invalidate()
					slog.Debug("mujica: channel_list fold", "group", g.ID, "open", !open)
				}
				return
			}
			ch := g.Channels[r.item]
			if channelRow(c, k, ch, ch.ID == *selected) && !col.IsDisabled() {
				*selected = ch.ID
			}
		}).Height(h).FocusGroup(ui.Vertical).Role(ui.RoleTree).Label(core.Msg(c, "channel_list.label"))
	})
	changed := *selected != old
	if changed {
		slog.Debug("mujica: channel_list select", "id", *selected)
	}
	return core.NewChoiceResult(col, changed)
}

func channelRows(st *channelListState, groups []ChannelGroup) {
	seen, groupIDs := map[string]bool{}, map[string]bool{}
	rows := []channelIndex{}
	for gi, g := range groups {
		if g.ID == "" || g.Title == "" || groupIDs[g.ID] {
			panic("mujica: ChannelList needs unique group IDs and titles")
		}
		groupIDs[g.ID] = true
		rows = append(rows, channelIndex{gi, -1})
		for _, ch := range g.Channels {
			if ch.ID == "" || ch.Name == "" || seen[ch.ID] || ch.Unread < 0 {
				panic("mujica: ChannelList needs unique channel IDs, names and nonnegative Unread")
			}
			ch.Presence.key()
			seen[ch.ID] = true
		}
		if st.folded[g.ID] {
			continue
		}
		for _, pinned := range []bool{true, false} {
			for ci, ch := range g.Channels {
				if ch.Pinned == pinned {
					rows = append(rows, channelIndex{gi, ci})
				}
			}
		}
	}
	st.rows = rows
}

func channelHeader(c *ui.Context, k theme.Tokens, g ChannelGroup, open bool) bool {
	u := core.Density(c).Unit()
	b := ui.ButtonBase(c).MinHeight(core.Density(c).ControlHeight()).Justify(ui.Start).Gap(u).Padding(0, u*2).Radius(theme.ControlRadius).
		Label(g.Title).Role(ui.RoleTreeItem).TextColor(k.TextMuted).Transition(core.ColorFade)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	arrow := "chevron-down"
	if !open {
		arrow = "chevron-right"
	}
	b.Children(func() {
		ui.Icon(c, icons.Must(arrow)).Size(14, 14)
		ui.Text(c, g.Title).FontSize(core.FontSize(c, theme.CaptionSize)).FontWeight(600).SingleLine()
	})
	if b.IsDisabled() {
		return false
	}
	toggled := b.Clicked()
	if open && b.Shortcut(0, ui.KeyLeft) || !open && b.Shortcut(0, ui.KeyRight) {
		toggled = true
	}
	return toggled
}

// channelRow is one channel button; true means it was chosen.
func channelRow(c *ui.Context, k theme.Tokens, ch Channel, on bool) bool {
	d := core.Density(c)
	u := d.Unit()
	label := ch.Name
	if ch.Unread > 0 {
		label = fmt.Sprintf(core.Msg(c, "channel_list.unread"), ch.Name, ch.Unread)
	}
	b := ui.ButtonBase(c).Key(ch.ID).MinHeight(d.ControlHeight()).Justify(ui.Start).Gap(u*2).Padding(0, u*2, 0, u*5).
		Radius(theme.ControlRadius).Label(label).Role(ui.RoleTreeItem).Selected(on).Background(ui.Transparent).
		TextColor(k.Text).Transition(core.ColorFade)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case on:
		b.Background(k.Selection).BorderWidth(0, 0, 0, 3).BorderColor(k.AccentText)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	if ch.Muted && !on {
		b.TextColor(k.TextMuted)
	}
	b.Children(func() {
		ui.Box(c).Width(16).Center().Shrink(0).Children(func() {
			if ch.Direct {
				OnlineStatus(c, ch.Presence, OnlineStatusOptions{HideText: true, Size: 9})
			} else {
				ui.Text(c, "#").TextColor(k.TextMuted)
			}
		})
		name := ui.Text(c, ch.Name).SingleLine().Grow(1).MinWidth(0)
		if ch.Unread > 0 && !ch.Muted {
			name.FontWeight(600)
		}
		if ch.Pinned {
			ui.Icon(c, msgPin).Size(13, 13).TextColor(k.TextMuted).Label(core.Msg(c, "channel_list.pinned"))
		}
		if ch.Muted {
			ui.Icon(c, msgBellOff).Size(13, 13).TextColor(k.TextMuted).Label(core.Msg(c, "channel_list.muted"))
		}
		if ch.Unread > 0 {
			tone := display.BadgeAccent
			if ch.Muted {
				tone = display.BadgeNeutral
			}
			display.Badge(c, display.BadgeOptions{Count: ch.Unread, Tone: tone}, nil)
		}
	})
	return b.Clicked()
}
