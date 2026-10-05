package feedback

import (
	"strconv"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func feedItems(n int) []Activity {
	base := time.Date(2026, 10, 5, 18, 0, 0, 0, time.UTC)
	out := make([]Activity, n)
	for i := range out {
		out[i] = Activity{ID: strconv.Itoa(i), Actor: "Mira", Action: "edited", Object: "Doc " + strconv.Itoa(i), Time: base.Add(-time.Duration(i) * 5 * time.Hour)}
	}
	return out
}

func TestActivityFeedGroupsByDay(t *testing.T) {
	items := feedItems(6) // 18:00, 13:00, 08:00, 03:00 on the 5th, then 22:00 and 17:00 on the 4th
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ActivityFeed(c, items, ActivityFeedOptions{Location: time.UTC, Today: core.Date{Year: 2026, Month: 10, Day: 5}}).Size(600, 500)
	}, 640, 540)
	tt.Frame()
	for _, s := range []string{"Today", "Yesterday", "Mira edited Doc 0, 18:00", "Mira edited Doc 4, 22:00"} {
		if _, ok := tt.Find(s); !ok {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	today, _ := tt.Find("Today")
	yest, _ := tt.Find("Yesterday")
	doc3, _ := tt.Find("Mira edited Doc 3, 03:00")
	doc4, _ := tt.Find("Mira edited Doc 4, 22:00")
	if !(today.Y < doc3.Y && doc3.Y < yest.Y && yest.Y < doc4.Y) {
		t.Fatalf("order: today %v, doc3 %v, yesterday %v, doc4 %v", today, doc3, yest, doc4)
	}
	// Another zone moves 03:00 UTC to the previous day.
	tt2 := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		ActivityFeed(c, items, ActivityFeedOptions{Location: time.FixedZone("W", -5*3600), Today: core.Date{Year: 2026, Month: 10, Day: 5}}).Size(600, 500)
	}, 640, 540)
	tt2.Frame()
	if _, ok := tt2.Find("Mira edited Doc 3, 22:00"); !ok || !tt2.HasText("今天") || !tt2.HasText("昨天") {
		t.Fatalf("zone texts %q", tt2.Texts())
	}
}

func TestActivityFeedVirtualizesAndEmpty(t *testing.T) {
	items := feedItems(100000)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ActivityFeed(c, items, ActivityFeedOptions{Location: time.UTC, Today: core.Date{Year: 2026, Month: 10, Day: 5}}).Size(600, 400)
	}, 640, 440)
	tt.Frame()
	if n := len(tt.Texts()); n == 0 || n > 120 {
		t.Fatalf("built %d texts for 100000 events", n)
	}
	tt.Scroll(100, 100, 0, 2000000)
	tt.Frame()
	if _, ok := tt.Find("Mira edited Doc 0, 18:00"); ok {
		t.Fatal("first event still built after scrolling far")
	}
	tt3 := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ActivityFeed(c, nil, ActivityFeedOptions{}).Size(400, 300)
	}, 440, 340)
	if !tt3.HasText("No activity yet") {
		t.Fatalf("empty: %q", tt3.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("event without ID accepted")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ActivityFeed(c, []Activity{{Actor: "x"}}, ActivityFeedOptions{})
	}, 200, 100)
}

func TestActivityFeedSnapshot(t *testing.T) {
	items := feedItems(5)
	items[1].Actor, items[1].Action, items[1].Object = "Lin", "commented on", "Quarterly report"
	items[3].Actor, items[3].Action, items[3].Object = "Sakiko", "approved", "Merge request !42"
	snapshot(t, "activity_feed", 520, 320, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			ActivityFeed(c, items, ActivityFeedOptions{Location: time.UTC, Today: core.Date{Year: 2026, Month: 10, Day: 5}}).Size(496, 296)
		})
	})
}
