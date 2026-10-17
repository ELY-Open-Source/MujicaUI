package feedback

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func sampleNotifications() []Notification {
	return []Notification{
		{ID: "a", Title: "Build failed", Body: "main · 3 tests failing", Time: "2 min ago", Action: "View log"},
		{ID: "b", Title: "Invite accepted", Body: "Mira joined the workspace", Time: "1 h ago"},
		{ID: "c", Title: "Backup complete", Time: "Yesterday", Read: true},
	}
}

func notificationTester(items *[]Notification, last *NotificationCenterResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Height(560).Children(func() {
			r := NotificationCenter(c, items, NotificationCenterOptions{})
			r.Element = ui.Element{}
			if r != (NotificationCenterResult{}) {
				*last = r
			}
		})
	}, 640, 600)
}

func TestNotificationCenterReadAndAction(t *testing.T) {
	items := sampleNotifications()
	var last NotificationCenterResult
	tt := notificationTester(&items, &last)
	if !tt.HasText("2 unread") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if err := tt.Click("Unread: Invite accepted"); err != nil {
		t.Fatal(err)
	}
	if !items[1].Read || last.Opened != "b" || !tt.HasText("1 unread") {
		t.Fatalf("open: %+v %+v", items[1], last)
	}
	if err := tt.Click("View log"); err != nil {
		t.Fatal(err)
	}
	if last.Action != "a" || items[0].Read {
		t.Fatalf("action: %+v %+v", last, items[0])
	}
	if err := tt.Click("Mark unread"); err != nil {
		t.Fatal(err)
	}
	if items[1].Read || last.Toggled != "b" {
		t.Fatalf("toggle: %+v %+v", items[1], last)
	}
	if err := tt.Click("Mark all read"); err != nil {
		t.Fatal(err)
	}
	for _, n := range items {
		if !n.Read {
			t.Fatalf("mark all: %+v", n)
		}
	}
	if !last.MarkedAllRead || tt.HasText("1 unread") || tt.HasText("2 unread") {
		t.Fatalf("mark all: %+v %q", last, tt.Texts())
	}
}

func TestNotificationCenterRemoveAndClear(t *testing.T) {
	items := sampleNotifications()
	var last NotificationCenterResult
	tt := notificationTester(&items, &last)
	if err := tt.Click("Clear"); err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].ID != "b" || last.Removed != "a" {
		t.Fatalf("remove: %+v %+v", items, last)
	}
	if err := tt.Click("Clear all"); err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 || !last.ClearedAll || !tt.HasText("You're all caught up") {
		t.Fatalf("clear all: %+v %q", items, tt.Texts())
	}
}

func TestNotificationCenterKeyboard(t *testing.T) {
	items := sampleNotifications()
	var last NotificationCenterResult
	tt := notificationTester(&items, &last)
	tt.Key(0, ui.KeyTab) // Mark all read
	tt.Key(0, ui.KeyTab) // Clear all
	tt.Key(0, ui.KeyTab) // first notification
	if !tt.Focused("Unread: Build failed") {
		t.Fatal("Tab did not enter the list")
	}
	tt.Key(0, ui.KeyEnter)
	if !items[0].Read || last.Opened != "a" {
		t.Fatalf("Enter: %+v", items[0])
	}
	tt.Key(0, ui.KeyDelete)
	if len(items) != 2 || items[0].ID != "b" {
		t.Fatalf("Delete: %+v", items)
	}
	tt.Key(0, ui.KeyTab)
	if tt.Focused("Unread: Invite accepted") {
		t.Fatal("the list takes more than one Tab stop")
	}
}

func TestNotificationCenterPanicsOnRepeatedID(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	items := []Notification{{ID: "x"}, {ID: "x"}}
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		NotificationCenter(c, &items, NotificationCenterOptions{})
	}, 300, 200)
}

func TestNotificationCenterSnapshot(t *testing.T) {
	items := sampleNotifications()
	snapshot(t, "notification_center", 520, 400, func(c *ui.Context) {
		ui.Column(c).Padding(16).Height(400).Children(func() {
			NotificationCenter(c, &items, NotificationCenterOptions{})
		})
	})
}
