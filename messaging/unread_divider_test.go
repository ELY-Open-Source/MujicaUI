package messaging

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestUnreadDivider(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Children(func() {
			UnreadDivider(c, 12)
			UnreadDivider(c, 0)
		})
	}, 400, 120)
	if !tt.HasText("12 new messages") || !tt.HasText("New") {
		t.Fatalf("texts %v", tt.Texts())
	}
	mediaPanics(t, "negative", func(c *ui.Context) { UnreadDivider(c, -1) })
}

func TestUnreadDividerSnapshot(t *testing.T) {
	snapshot(t, "unread_divider", 420, 70, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() { UnreadDivider(c, 3) })
	})
}
