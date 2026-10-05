package chat

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// messageListApp is a chat of n messages with stable IDs.
type messageListApp struct {
	ids   []string
	state MessageListState
	built int
	res   MessageListResult
}

func newMessageListApp(n int) *messageListApp {
	a := &messageListApp{}
	for i := range n {
		a.ids = append(a.ids, fmt.Sprintf("m%d", i))
	}
	return a
}

func (a *messageListApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	a.built = 0
	day := time.Date(2026, 3, 1, 9, 0, 0, 0, time.UTC)
	a.res = MessageList(c, &a.state, len(a.ids), MessageListOptions{
		ID:   func(i int) string { return a.ids[i] },
		Date: func(i int) time.Time { return day.Add(time.Duration(i/3) * 24 * time.Hour) },
	}, func(i int) {
		a.built++
		ui.Text(c, "Message "+a.ids[i]).Padding(8)
	})
}

func TestMessageListVirtualizedAndFollowsEnd(t *testing.T) {
	a := newMessageListApp(100000)
	tt := ui.NewTester(a.view, 400, 300)
	tt.Frame()
	if a.built == 0 || a.built > 60 {
		t.Fatalf("built %d rows of 100000", a.built)
	}
	if !tt.HasText("Message m99999") || !a.state.AtEnd() {
		t.Fatal("list did not start at the newest message")
	}
	a.ids = append(a.ids, "m100000")
	tt.Frame()
	tt.Frame()
	if !tt.HasText("Message m100000") || a.res.Unread() != 0 {
		t.Fatalf("new message not followed, unread %d", a.res.Unread())
	}
}

func TestMessageListStopsFollowingAndCountsUnread(t *testing.T) {
	a := newMessageListApp(50)
	tt := ui.NewTester(a.view, 400, 300)
	tt.Frame()
	tt.Scroll(200, 150, 0, -2000)
	tt.Frame()
	if a.state.AtEnd() {
		t.Fatal("still at end after scrolling up")
	}
	a.ids = append(a.ids, "n1", "n2")
	tt.Frame()
	tt.Frame()
	if tt.HasText("Message n2") || a.res.Unread() != 2 {
		t.Fatalf("followed while scrolled up, unread %d", a.res.Unread())
	}
	if err := tt.Click("Scroll to latest, 2 new"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	tt.Frame()
	if !tt.HasText("Message n2") || a.res.Unread() != 0 {
		t.Fatalf("jump did not reach the end, unread %d", a.res.Unread())
	}
}

func TestMessageListKeepsPlaceByID(t *testing.T) {
	a := newMessageListApp(50)
	tt := ui.NewTester(a.view, 400, 300)
	tt.Frame()
	tt.Scroll(200, 150, 0, -600)
	tt.Frame()
	r, _ := tt.Find("Message m40")
	// Older messages load above; the one in view stays put.
	a.ids = append([]string{"old1", "old2", "old3"}, a.ids...)
	tt.Frame()
	r2, ok := tt.Find("Message m40")
	if !ok || r2.Y != r.Y {
		t.Fatalf("m40 moved from %v to %v (found %v)", r.Y, r2.Y, ok)
	}
}

func TestMessageListDaySeparators(t *testing.T) {
	a := newMessageListApp(6)
	tt := ui.NewTester(a.view, 400, 600)
	tt.Frame()
	for _, d := range []string{"Mar 1, 2026", "Mar 2, 2026"} {
		if !tt.HasText(d) {
			t.Errorf("missing separator %q in %q", d, tt.Texts())
		}
	}
}

func TestMessageListNeedsID(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without ID")
		}
	}()
	var s MessageListState
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		MessageList(c, &s, 1, MessageListOptions{}, func(int) {})
	}, 200, 100)
}

func TestMessageListSnapshot(t *testing.T) {
	snapshot(t, "message_list", 420, 320, func(c *ui.Context) {
		(*ui.Local(c.Root(), "app", func() *messageListApp { return newMessageListApp(6) })).view(c)
	})
}
