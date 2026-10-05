package mujica

import (
	"runtime"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSearchInputSearchClearSubmit(t *testing.T) {
	q, changes, submits := "", 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := SearchInput(c, &q, SearchInputOptions{})
		if r.Changed() {
			changes++
		}
		if r.Submitted() {
			submits++
		}
	}, 400, 200)
	_ = tt.Click("Search")
	tt.Type("rose")
	tt.Key(0, ui.KeyEnter)
	if q != "rose" || changes != 1 || submits != 1 {
		t.Fatalf("query %q changes %d submits %d", q, changes, submits)
	}
	_ = tt.Click("Clear")
	if q != "" || changes != 2 {
		t.Fatalf("after clear %q changes %d", q, changes)
	}
	tt.Type("lys")
	tt.Key(0, ui.KeyEscape)
	if q != "" || changes != 4 {
		t.Fatalf("after Escape %q changes %d", q, changes)
	}
}

// searchInputFindTester counts matches of the query in a fixed text.
func searchInputFindTester(q *string, cur, moves *int, locale Locale) *ui.Tester {
	const text = "rose lys rose crown rose"
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Locale: locale})
		n := 0
		if *q != "" {
			n = strings.Count(text, *q)
		}
		if SearchInput(c, q, SearchInputOptions{Find: true, Matches: n, Current: cur, Label: "Find in page"}).Moved() {
			*moves++
		}
	}, 500, 200)
}

func TestSearchInputFindBarNavigation(t *testing.T) {
	q, cur, moves := "", 0, 0
	tt := searchInputFindTester(&q, &cur, &moves, "")
	_ = tt.Click("Find in page")
	tt.Type("rose")
	if !tt.HasText("1 of 3") {
		t.Fatalf("count missing: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEnter)
	if cur != 2 || !tt.HasText("3 of 3") {
		t.Fatalf("after two Enters current %d texts %q", cur, tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if cur != 0 {
		t.Fatalf("Enter did not wrap, current %d", cur)
	}
	tt.Key(ui.Shift, ui.KeyEnter)
	if cur != 2 || q != "rose" {
		t.Fatalf("Shift+Enter: current %d query %q", cur, q)
	}
	_ = tt.Click("Previous match")
	_ = tt.Click("Next match")
	_ = tt.Click("Next match")
	if cur != 0 || moves != 7 {
		t.Fatalf("buttons: current %d moves %d", cur, moves)
	}
	if runtime.GOOS == "darwin" {
		tt.Key(ui.Cmd, ui.KeyG)
	} else {
		tt.Key(0, ui.KeyF3)
	}
	if cur != 1 {
		t.Fatalf("find-next shortcut: current %d", cur)
	}
}

func TestSearchInputFindNoMatches(t *testing.T) {
	q, cur, moves := "", 0, 0
	tt := searchInputFindTester(&q, &cur, &moves, ZhCN)
	_ = tt.Click("Find in page")
	tt.Type("tulip")
	if !tt.HasText("无匹配") {
		t.Fatalf("no-match status missing: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("下一项")
	if cur != 0 || moves != 0 {
		t.Fatalf("moved without matches: current %d moves %d", cur, moves)
	}
}

func TestSearchInputFindNeedsCurrent(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without Current")
		}
	}()
	q := ""
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		SearchInput(c, &q, SearchInputOptions{Find: true})
	}, 200, 100)
}

func TestSearchInputSnapshot(t *testing.T) {
	snapshot(t, "search_input", 460, 220, func(c *ui.Context) {
		a, b, f, cur := "", "Marie", "rose", 1
		ui.Column(c).Padding(24).Gap(12).Width(400).Children(func() {
			SearchInput(c, &a, SearchInputOptions{})
			SearchInput(c, &b, SearchInputOptions{})
			SearchInput(c, &f, SearchInputOptions{Find: true, Matches: 3, Current: &cur})
		})
	})
}
