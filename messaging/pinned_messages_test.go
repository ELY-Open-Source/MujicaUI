package messaging

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestPinnedMessages(t *testing.T) {
	pins := []ChatEntry{{ID: "a", Author: "Sakiko", Text: "Rehearsal moves to nine."}, {ID: "b", Author: "Uika", Text: "Masks are in the hall."}}
	var jumps, unpins []string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(420).Children(func() {
			r := PinnedMessages(c, pins, PinnedMessagesOptions{})
			if r.Jump != "" {
				jumps = append(jumps, r.Jump)
			}
			if r.Unpinned != "" {
				unpins = append(unpins, r.Unpinned)
			}
		})
	}, 500, 400)
	if !tt.HasText("Pinned messages (2)") {
		t.Fatalf("texts %v", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if err := tt.Click("Unpin Uika's message"); err != nil {
		t.Fatal(err)
	}
	if len(jumps) != 1 || jumps[0] != "a" || len(unpins) != 1 || unpins[0] != "b" {
		t.Fatalf("jumps %v unpins %v", jumps, unpins)
	}
	pins = nil
	tt.Frame()
	if !tt.HasText("No pinned messages") {
		t.Fatal("no empty state")
	}
}

func TestPinnedMessagesSnapshot(t *testing.T) {
	pins := []ChatEntry{{ID: "a", Author: "Sakiko", Text: "Rehearsal moves to nine."}, {ID: "b", Author: "Uika", Text: "Masks are in the hall."}}
	snapshot(t, "pinned_messages", 440, 200, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(400).Children(func() { PinnedMessages(c, pins, PinnedMessagesOptions{}) })
	})
}
