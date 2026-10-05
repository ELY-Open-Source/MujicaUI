package mujica

import (
	"fmt"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

func watchlistTester(s *WatchlistState, groups *[]WatchGroup, opts WatchlistOptions, got *WatchlistView) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(10).Height(400).Children(func() {
			v := Watchlist(c, s, *groups, opts)
			if v.Changed || v.Submitted || v.Moved {
				*got = v
			}
			if v.Moved {
				g := &(*groups)[s.Group]
				g.Items = watchlistMove(g.Items, v.Rows, v.To)
			}
		})
	}, 700, 420)
	tt.Frame()
	return tt
}

// watchlistMove is how a caller applies a Moved event.
func watchlistMove(items []WatchItem, rows []int, to int) []WatchItem {
	var moved, rest []WatchItem
	at := to
	for i, it := range items {
		if slices.Contains(rows, i) {
			moved = append(moved, it)
			if i < to {
				at--
			}
		} else {
			rest = append(rest, it)
		}
	}
	return slices.Insert(rest, at, moved...)
}

func TestWatchlistSelectKeepsThroughUpdates(t *testing.T) {
	var s WatchlistState
	groups := mktWatch()
	var got WatchlistView
	tt := watchlistTester(&s, &groups, WatchlistOptions{}, &got)
	if err := tt.Click("MSFT"); err != nil {
		t.Fatal(err)
	}
	if s.Selected != "MSFT" || !got.Changed {
		t.Fatalf("selected %q", s.Selected)
	}
	// A refresh with new prices in a new order keeps the choice on MSFT.
	groups[0].Items = []WatchItem{groups[0].Items[2], groups[0].Items[1], groups[0].Items[0]}
	groups[0].Items[1].Price = 420
	got = WatchlistView{}
	tt.Frame()
	tt.Key(0, ui.KeyDown)
	if s.Selected != "AAPL" {
		t.Fatalf("Down after refresh chose %q", s.Selected)
	}
	tt.Key(0, ui.KeyEnter)
	if !got.Submitted {
		t.Fatal("Enter did not submit")
	}
}

func TestWatchlistColumnsAndGroups(t *testing.T) {
	var s WatchlistState
	groups := mktWatch()
	var got WatchlistView
	tt := watchlistTester(&s, &groups, WatchlistOptions{}, &got)
	if !tt.HasText("51,000,000") {
		t.Fatalf("volume missing: %q", tt.Texts())
	}
	_ = tt.Click("Columns")
	_ = tt.Click("Volume")
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("51,000,000") || !s.Hidden[WatchVolume] {
		t.Fatal("volume column still shown")
	}
	_ = tt.Click("Energy")
	if s.Group != 1 || !tt.HasText("XOM") || tt.HasText("AAPL") {
		t.Fatalf("group %d texts %q", s.Group, tt.Texts())
	}
}

func TestWatchlistDragReorders(t *testing.T) {
	var s WatchlistState
	groups := mktWatch()
	var got WatchlistView
	tt := watchlistTester(&s, &groups, WatchlistOptions{}, &got)
	from, _ := tt.Find("NVDA")
	to, _ := tt.Find("AAPL")
	tt.Press(from.X, from.Y+from.H/2)
	tt.Move(from.X, from.Y-10)
	tt.Move(to.X, to.Y)
	tt.Release(to.X, to.Y)
	order := []string{}
	for _, it := range groups[0].Items {
		order = append(order, it.Symbol)
	}
	if !got.Moved || order[0] != "NVDA" {
		t.Fatalf("moved %v order %v", got.Moved, order)
	}
}

func TestWatchlistDisabledAndLarge(t *testing.T) {
	var s WatchlistState
	groups := mktWatch()
	var got WatchlistView
	tt := watchlistTester(&s, &groups, WatchlistOptions{Disabled: true}, &got)
	_ = tt.Click("MSFT")
	if s.Selected != "" || got.Changed {
		t.Fatalf("disabled list chose %q", s.Selected)
	}
	big := []WatchGroup{{Name: "All"}}
	for i := range 100000 {
		big[0].Items = append(big[0].Items, WatchItem{Symbol: fmt.Sprint("S", i), Price: 1})
	}
	var b WatchlistState
	tt = watchlistTester(&b, &big, WatchlistOptions{}, &got)
	if n := len(tt.Texts()); n > 200 {
		t.Fatalf("built %d texts for 100000 rows", n)
	}
}

func TestWatchlistSnapshot(t *testing.T) {
	snapshot(t, "watchlist", 640, 260, func(c *ui.Context) {
		s := &WatchlistState{Selected: "MSFT"}
		ui.Column(c).Padding(16).Height(260).Children(func() { Watchlist(c, s, mktWatch(), WatchlistOptions{}) })
	})
}
