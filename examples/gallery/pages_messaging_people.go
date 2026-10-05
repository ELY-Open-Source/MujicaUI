package main

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
)

func channelListDemo(c *ui.Context, log func(string)) {
	type state struct {
		selected string
		disabled bool
		revision int
		groups   []mujica.ChannelGroup
	}
	s := ui.Local(c.Root(), "401", func() state {
		return state{groups: []mujica.ChannelGroup{
			{ID: "channels", Title: "Channels", Channels: []mujica.Channel{
				{ID: "general", Name: "general", Unread: 3}, {ID: "rehearsal", Name: "rehearsal", Pinned: true},
				{ID: "random", Name: "random", Unread: 120, Muted: true},
			}},
			{ID: "direct", Title: "Direct messages", Channels: []mujica.Channel{{ID: "uika", Name: "Uika", Direct: true, Presence: mujica.PresenceOnline}}},
		}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	before := s.groups[0].Channels[0]
	mediaCheck(c, log, &s.groups[0].Channels[0].Muted, "Mute general")
	mediaCheck(c, log, &s.groups[0].Channels[0].Pinned, "Pin general")
	if before != s.groups[0].Channels[0] {
		s.revision++
	}
	if mujica.ChannelList(c, &s.selected, s.groups, mujica.ChannelListOptions{Disabled: s.disabled, Revision: s.revision}).Changed() {
		log("channel " + s.selected)
	}
	ui.Text(c, "Selected: "+s.selected)
}

func memberListDemo(c *ui.Context, log func(string)) {
	type state struct {
		list            mujica.MemberListState
		disabled, empty bool
		members         []mujica.Member
	}
	s := ui.Local(c.Root(), "405", func() state {
		return state{members: []mujica.Member{
			{ID: "1", Name: "Sakiko", Presence: mujica.PresenceOnline, Role: "Leader"},
			{ID: "2", Name: "Uika", Presence: mujica.PresenceBusy, Role: "Vocal"},
			{ID: "3", Name: "Mutsumi", Presence: mujica.PresenceAway}, {ID: "4", Name: "Umiri", Presence: mujica.PresenceOffline},
		}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	mediaCheck(c, log, &s.empty, "Empty")
	if mujica.Button(c, "Toggle Uika presence", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
		s.members = slices.Clone(s.members)
		if s.members[1].Presence == mujica.PresenceOffline {
			s.members[1].Presence = mujica.PresenceOnline
		} else {
			s.members[1].Presence = mujica.PresenceOffline
		}
		log("Uika presence changed")
	}
	members := s.members
	if s.empty {
		members = nil
	}
	r := mujica.MemberList(c, &s.list, members, mujica.MemberListOptions{Disabled: s.disabled, Height: 240})
	if r.Changed() {
		log("member " + s.list.Selected)
	}
	if r.Opened != "" {
		log("profile " + r.Opened)
	}
}

func userProfileCardDemo(c *ui.Context, log func(string)) {
	type state struct {
		presence          mujica.Presence
		actions, disabled bool
	}
	s := ui.Local(c.Root(), "406", func() state { return state{actions: true, presence: mujica.PresenceOnline} })
	mediaCheck(c, log, &s.actions, "Actions")
	mediaCheck(c, log, &s.disabled, "Disable actions")
	messagingPresence(c, log, &s.presence)
	opts := mujica.UserProfileCardOptions{}
	if s.actions {
		opts.Actions = func() {
			for _, action := range []string{"Message", "Call"} {
				if mujica.Button(c, action, mujica.ButtonOptions{Disabled: s.disabled}).Clicked() {
					log(action)
				}
			}
		}
	}
	mujica.UserProfileCard(c, mujica.UserProfile{Name: "Sakiko Togawa", Presence: s.presence, Role: "Keyboard", Status: "In rehearsal", About: "Founder of the masquerade."}, opts)
}

func onlineStatusDemo(c *ui.Context, log func(string)) {
	type state struct {
		presence mujica.Presence
		hidden   bool
	}
	s := ui.Local(c.Root(), "407", func() state { return state{presence: mujica.PresenceOnline} })
	mediaCheck(c, log, &s.hidden, "Hide text")
	messagingPresence(c, log, &s.presence)
	mujica.OnlineStatus(c, s.presence, mujica.OnlineStatusOptions{HideText: s.hidden})
	inputSection(c, "States", func() {
		for p := mujica.PresenceOffline; p <= mujica.PresenceBusy; p++ {
			mujica.OnlineStatus(c, p, mujica.OnlineStatusOptions{})
		}
	})
}

func statusSetterDemo(c *ui.Context, log func(string)) {
	type state struct {
		status   mujica.CustomStatus
		disabled bool
	}
	s := ui.Local(c.Root(), "408", func() state { return state{} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if mujica.StatusSetter(c, &s.status, mujica.StatusSetterOptions{Disabled: s.disabled}).Changed() {
		log("status changed")
	}
	ui.Text(c, fmt.Sprintf("Text: %s; expires: %v", s.status.Text, s.status.Until))
}

func messagingPresence(c *ui.Context, log func(string), value *mujica.Presence) {
	items := []mujica.SelectOption[mujica.Presence]{
		{Value: mujica.PresenceOnline, Label: "Online"}, {Value: mujica.PresenceAway, Label: "Away"},
		{Value: mujica.PresenceBusy, Label: "Busy"}, {Value: mujica.PresenceOffline, Label: "Offline"},
	}
	if mujica.Select(c, value, items, mujica.SelectOptions{Label: "Presence"}).Changed() {
		log(fmt.Sprintf("presence %d", *value))
	}
}
