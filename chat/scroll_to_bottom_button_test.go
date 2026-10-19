package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestScrollToBottomButtonShowsWhenAway(t *testing.T) {
	away, unread, clicks := false, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Box(c).Fill().Children(func() {
			if ScrollToBottomButton(c, away, unread).Clicked() {
				clicks++
			}
		})
	}, 300, 200)
	if _, ok := tt.Find("Scroll to latest"); ok {
		t.Fatal("shown while at the bottom")
	}
	away, unread = true, 4
	tt.Frame()
	if !tt.HasText("4") {
		t.Fatalf("unread count missing: %q", tt.Texts())
	}
	_ = tt.Click("Scroll to latest, 4 new")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if clicks != 2 {
		t.Fatalf("clicks %d", clicks)
	}
}

func TestScrollToBottomButtonNegativeUnread(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for negative unread")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ScrollToBottomButton(c, true, -1)
	}, 100, 100)
}

func TestScrollToBottomButtonSnapshot(t *testing.T) {
	snapshot(t, "scroll_to_bottom_button", 240, 80, func(c *ui.Context) {
		ui.Box(c).Fill().Children(func() { ScrollToBottomButton(c, true, 12) })
	})
}
