package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

func TestEmptyStateKinds(t *testing.T) {
	opts := EmptyStateOptions{}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		EmptyState(c, opts)
	}, 500, 400)
	if !tt.HasText("Nothing here yet") || !tt.HasText("Items you add will appear here.") {
		t.Fatalf("no data: %q", tt.Texts())
	}
	opts = EmptyStateOptions{Kind: EmptyNoResults, Query: "invoice"}
	tt.Frame()
	if !tt.HasText("No results") || !tt.HasText("Nothing matches “invoice”.") {
		t.Fatalf("no results: %q", tt.Texts())
	}
	opts = EmptyStateOptions{Kind: EmptyNoResults}
	tt.Frame()
	if !tt.HasText("Try other words or clear the filters.") {
		t.Fatalf("no results without query: %q", tt.Texts())
	}
}

func TestEmptyStateActionKeyboard(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Locale: ZhCN})
		if EmptyState(c, EmptyStateOptions{Kind: EmptyNoResults, Query: "发票", Action: "清除搜索"}).Action {
			n++
		}
	}, 500, 400)
	if !tt.HasText("没有与“发票”匹配的内容。") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if err := tt.Click("清除搜索"); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("actions %d", n)
	}
}

func TestEmptyStatePanics(t *testing.T) {
	for name, o := range map[string]EmptyStateOptions{
		"icon and no icon":   {Icon: icons.Must("inbox"), NoIcon: true},
		"query without kind": {Query: "x"},
		"unknown kind":       {Kind: 9},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); EmptyState(c, o) }, 300, 200)
		}()
	}
}

func TestEmptyStateSnapshot(t *testing.T) {
	snapshot(t, "empty_state", 640, 280, func(c *ui.Context) {
		ui.Row(c).Padding(12).Gap(12).Children(func() {
			ui.Column(c).Grow(1).Basis(0).Children(func() {
				EmptyState(c, EmptyStateOptions{Action: "New project"})
			})
			ui.Column(c).Grow(1).Basis(0).Children(func() {
				EmptyState(c, EmptyStateOptions{Kind: EmptyNoResults, Query: "quarterly", Action: "Clear search"})
			})
		})
	})
}
