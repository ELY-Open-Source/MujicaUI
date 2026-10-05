package messaging

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var channelTestGroups = []ChannelGroup{
	{ID: "ch", Title: "Channels", Channels: []Channel{
		{ID: "general", Name: "general", Unread: 3},
		{ID: "rehearsal", Name: "rehearsal", Pinned: true},
		{ID: "random", Name: "random", Unread: 120, Muted: true},
	}},
	{ID: "dm", Title: "Direct messages", Channels: []Channel{{ID: "uika", Name: "Uika", Direct: true, Presence: PresenceOnline}}},
}

func channelTester(sel *string, disabled bool, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(260).Children(func() {
			if ChannelList(c, sel, channelTestGroups, ChannelListOptions{Disabled: disabled}).Changed() {
				*n++
			}
		})
	}, 320, 480)
}

func memberTestList(n int) []Member {
	out := make([]Member, n)
	for i := range out {
		out[i] = Member{ID: fmt.Sprint("m", i), Name: fmt.Sprint("Member ", i), Presence: PresenceStatus(i % 4)}
	}
	out[0].Role = "Owner"
	return out
}
