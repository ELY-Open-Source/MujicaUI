package messaging

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestChannelListSelectFoldAndOrder(t *testing.T) {
	sel, n := "", 0
	tt := channelTester(&sel, false, &n)
	if err := tt.Click("random, 120 unread"); err != nil {
		t.Fatal(err)
	}
	if sel != "random" || n != 1 || !tt.HasText("99+") {
		t.Fatalf("selected %q, changes %d", sel, n)
	}
	pinned, _ := tt.Find("rehearsal")
	general, _ := tt.Find("general, 3 unread")
	if pinned.Y >= general.Y {
		t.Fatal("pinned channel not first")
	}
	if err := tt.Click("Direct messages"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Uika") {
		t.Fatal("folded group still shows its channels")
	}
}

func TestChannelListKeyboard(t *testing.T) {
	sel, n := "", 0
	tt := channelTester(&sel, false, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if sel != "rehearsal" {
		t.Fatalf("Down, Enter chose %q", sel)
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyLeft)
	if tt.HasText("general") {
		t.Fatal("Left on a header did not fold it")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("general") {
		t.Fatal("Right did not unfold")
	}
}

func TestChannelListDisabledAndBad(t *testing.T) {
	sel, n := "", 0
	tt := channelTester(&sel, true, &n)
	_ = tt.Click("rehearsal")
	if sel != "" || n != 0 {
		t.Fatal("disabled list chose")
	}
	dup := []ChannelGroup{{ID: "a", Title: "A", Channels: []Channel{{ID: "x", Name: "x"}, {ID: "x", Name: "y"}}}}
	mediaPanics(t, "duplicate", func(c *ui.Context) { ChannelList(c, &sel, dup, ChannelListOptions{}) })
}

func TestChannelListSnapshot(t *testing.T) {
	sel := "general"
	snapshot(t, "channel_list", 300, 300, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(260).Children(func() { ChannelList(c, &sel, channelTestGroups, ChannelListOptions{}) })
	})
}
