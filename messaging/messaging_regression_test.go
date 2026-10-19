package messaging

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestMessagingLargeListsBounded(t *testing.T) {
	const count = 100000
	channels := make([]Channel, count)
	pins := make([]ChatEntry, count)
	for i := range count {
		channels[i] = Channel{ID: fmt.Sprint(i), Name: fmt.Sprint("Channel ", i)}
		pins[i] = ChatEntry{ID: fmt.Sprint(i), Author: "Author", Text: fmt.Sprint("Pinned ", i)}
	}
	for _, name := range []string{"channels", "pins"} {
		t.Run(name, func(t *testing.T) {
			selected := ""
			groups := []ChannelGroup{{ID: "g", Title: "Group", Channels: channels}}
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ui.Column(c).Height(280).Children(func() {
					if name == "channels" {
						ChannelList(c, &selected, groups, ChannelListOptions{Height: 240})
					} else {
						PinnedMessages(c, pins, PinnedMessagesOptions{Height: 240})
					}
				})
			}, 400, 300)
			prefix := "Channel "
			if name == "pins" {
				prefix = "Pinned "
			}
			n := 0
			for _, text := range tt.Texts() {
				if len(text) >= len(prefix) && text[:len(prefix)] == prefix {
					n++
				}
			}
			if n == 0 || n > 60 {
				t.Fatalf("built %d of %d", n, count)
			}
		})
	}
}
