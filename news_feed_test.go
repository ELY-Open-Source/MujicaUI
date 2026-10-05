package mujica

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestNewsFeedReadSourceSymbolsAndKeyboard(t *testing.T) {
	var s NewsFeedState
	var got NewsFeedView
	rows := []NewsItem{{ID: "a", Title: "Alpha", Source: "Wire", Time: mktNow, Symbols: []string{"AAPL"}}, {ID: "b", Title: "Beta", Source: "Wire", Time: mktNow}}
	tt := chatTester(640, 420, func(c *ui.Context) {
		v := NewsFeed(c, &s, rows, NewsFeedOptions{Today: DateOf(mktNow)})
		v.Element.Height(380)
		if v.Opened != "" || v.Symbol != "" {
			got = v
		}
	})
	if !tt.HasText("2 unread") {
		t.Fatalf("count %q", tt.Texts())
	}
	_ = tt.Click("AAPL")
	if got.Symbol != "AAPL" || s.Read["a"] {
		t.Fatalf("symbol %+v read %v", got, s.Read)
	}
	_ = tt.Click("Alpha")
	if !s.Read["a"] || got.Opened != "a" {
		t.Fatalf("read %v", s.Read)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if !s.Read["b"] {
		t.Fatal("keyboard read missing")
	}
	tt.Frame()
	if !tt.HasText("0 unread") {
		t.Fatalf("count %q", tt.Texts())
	}
}
func TestNewsFeedSnapshot(t *testing.T) {
	snapshot(t, "news_feed", 560, 360, func(c *ui.Context) {
		ns := ui.Local(c.Root(), "snapshot.news_feed", func() NewsFeedState { return NewsFeedState{} })
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			NewsFeed(c, ns, []NewsItem{{ID: "1", Title: "Apple unveils new chips", Source: "Reuters", Time: mktNow, Symbols: []string{"AAPL"}},
				{ID: "2", Title: "Oil falls as supply grows", Source: "Bloomberg", Time: mktNow.Add(-26 * 3600e9), Symbols: []string{"XOM", "CVX"}}}, NewsFeedOptions{Today: DateOf(mktNow)}).Element.Height(300)

		})
	})
}
