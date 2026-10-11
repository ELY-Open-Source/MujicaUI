package main

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/messaging"
)

func channelListDemo(c *ui.Context, log func(string)) {
	type state struct {
		selected string
		disabled bool
		revision int
		groups   []messaging.ChannelGroup
	}
	s := ui.Local(c.Root(), "401", func() state {
		return state{groups: []messaging.ChannelGroup{
			{ID: "channels", Title: "Channels", Channels: []messaging.Channel{
				{ID: "general", Name: "general", Unread: 3}, {ID: "rehearsal", Name: "rehearsal", Pinned: true},
				{ID: "random", Name: "random", Unread: 120, Muted: true},
			}},
			{ID: "direct", Title: "Direct messages", Channels: []messaging.Channel{{ID: "uika", Name: "Uika", Direct: true, Presence: messaging.PresenceOnline}}},
		}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	before := s.groups[0].Channels[0]
	mediaCheck(c, log, &s.groups[0].Channels[0].Muted, "Mute general")
	mediaCheck(c, log, &s.groups[0].Channels[0].Pinned, "Pin general")
	if before != s.groups[0].Channels[0] {
		s.revision++
	}
	if messaging.ChannelList(c, &s.selected, s.groups, messaging.ChannelListOptions{Disabled: s.disabled, Revision: s.revision}).Changed() {
		log("channel " + s.selected)
	}
	ui.Text(c, "Selected: "+s.selected)
}

func memberListDemo(c *ui.Context, log func(string)) {
	type state struct {
		list            messaging.MemberListState
		disabled, empty bool
		members         []messaging.Member
	}
	s := ui.Local(c.Root(), "405", func() state {
		return state{members: []messaging.Member{
			{ID: "1", Name: "Sakiko", Presence: messaging.PresenceOnline, Role: "Leader"},
			{ID: "2", Name: "Uika", Presence: messaging.PresenceBusy, Role: "Vocal"},
			{ID: "3", Name: "Mutsumi", Presence: messaging.PresenceAway}, {ID: "4", Name: "Umiri", Presence: messaging.PresenceOffline},
		}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	mediaCheck(c, log, &s.empty, "Empty")
	if input.Button(c, "Toggle Uika presence", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
		s.members = slices.Clone(s.members)
		if s.members[1].Presence == messaging.PresenceOffline {
			s.members[1].Presence = messaging.PresenceOnline
		} else {
			s.members[1].Presence = messaging.PresenceOffline
		}
		log("Uika presence changed")
	}
	members := s.members
	if s.empty {
		members = nil
	}
	r := messaging.MemberList(c, &s.list, members, messaging.MemberListOptions{Disabled: s.disabled, Height: 240})
	if r.Changed() {
		log("member " + s.list.Selected)
	}
	if r.Opened != "" {
		log("profile " + r.Opened)
	}
}

func userProfileCardDemo(c *ui.Context, log func(string)) {
	type state struct {
		presence          messaging.PresenceStatus
		actions, disabled bool
	}
	s := ui.Local(c.Root(), "406", func() state { return state{actions: true, presence: messaging.PresenceOnline} })
	mediaCheck(c, log, &s.actions, "Actions")
	mediaCheck(c, log, &s.disabled, "Disable actions")
	messagingPresence(c, log, &s.presence)
	opts := messaging.UserProfileCardOptions{}
	if s.actions {
		opts.Actions = func() {
			for _, action := range []string{"Message", "Call"} {
				if input.Button(c, action, input.ButtonOptions{Disabled: s.disabled}).Clicked() {
					log(action)
				}
			}
		}
	}
	messaging.UserProfileCard(c, messaging.UserProfile{Name: "Sakiko Togawa", Presence: s.presence, Role: "Keyboard", Status: "In rehearsal", About: "Founder of the masquerade."}, opts)
}

func onlineStatusDemo(c *ui.Context, log func(string)) {
	type state struct {
		presence messaging.PresenceStatus
		hidden   bool
	}
	s := ui.Local(c.Root(), "407", func() state { return state{presence: messaging.PresenceOnline} })
	mediaCheck(c, log, &s.hidden, "Hide text")
	messagingPresence(c, log, &s.presence)
	messaging.OnlineStatus(c, s.presence, messaging.OnlineStatusOptions{HideText: s.hidden})
	inputSection(c, "States", func() {
		for p := messaging.PresenceOffline; p <= messaging.PresenceBusy; p++ {
			messaging.OnlineStatus(c, p, messaging.OnlineStatusOptions{})
		}
	})
}

func statusSetterDemo(c *ui.Context, log func(string)) {
	type state struct {
		status   messaging.CustomStatus
		disabled bool
	}
	s := ui.Local(c.Root(), "408", func() state { return state{} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if messaging.StatusSetter(c, &s.status, messaging.StatusSetterOptions{Disabled: s.disabled}).Changed() {
		log("status changed")
	}
	ui.Text(c, fmt.Sprintf("Text: %s; expires: %v", s.status.Text, s.status.Until))
}

func messagingPresence(c *ui.Context, log func(string), value *messaging.PresenceStatus) {
	items := []input.SelectOption[messaging.PresenceStatus]{
		{Value: messaging.PresenceOnline, Label: "Online"}, {Value: messaging.PresenceAway, Label: "Away"},
		{Value: messaging.PresenceBusy, Label: "Busy"}, {Value: messaging.PresenceOffline, Label: "Offline"},
	}
	if input.Select(c, value, items, input.SelectOptions{Label: "Presence"}).Changed() {
		log(fmt.Sprintf("presence %d", *value))
	}
}
