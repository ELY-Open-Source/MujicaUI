package mujica

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func chatTester(w, h int, view func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(10).Fill().Children(func() { view(c) })
	}, w, h)
}

func TestConversationListGroups(t *testing.T) {
	var s ConversationListState
	items := chatConversations()
	tt := chatTester(400, 600, func(c *ui.Context) {
		ConversationList(c, &s, items, ConversationListOptions{Now: chatNow}, nil).Element.Height(560)
	})
	want := []string{"Pinned", "Royal letters", "Today", "Garden maze", "Court etiquette notes", "Yesterday",
		"Ballroom seating plan", "Previous 7 days", "Harpsichord repair", "Older", "Coronation menu"}
	got := slices.Compact(slices.DeleteFunc(tt.Texts(), func(s string) bool { return !slices.Contains(want, s) }))
	if !slices.Equal(got, want) {
		t.Fatalf("rows %q", got)
	}
}

func TestConversationListSelectAndKeys(t *testing.T) {
	var s ConversationListState
	items := chatConversations()
	changes, opened := 0, 0
	tt := chatTester(400, 600, func(c *ui.Context) {
		v := ConversationList(c, &s, items, ConversationListOptions{Now: chatNow, Label: "History"}, nil)
		v.Element.Height(560)
		if v.Changed() {
			changes++
		}
		if v.Submitted() {
			opened++
		}
	})
	if err := tt.Click("Garden maze"); err != nil {
		t.Fatal(err)
	}
	if s.Selected != "f" || changes != 1 {
		t.Fatalf("selected %q changes %d", s.Selected, changes)
	}
	tt.Key(0, ui.KeyDown)
	if s.Selected != "a" {
		t.Fatalf("Down chose %q", s.Selected)
	}
	tt.Key(0, ui.KeyDown)
	if s.Selected != "b" {
		t.Fatalf("Down over a header chose %q", s.Selected)
	}
	tt.Key(0, ui.KeyEnter)
	if opened != 1 {
		t.Fatalf("Enter opened %d", opened)
	}
}

func TestConversationListDisabled(t *testing.T) {
	var s ConversationListState
	items := chatConversations()
	tt := chatTester(400, 600, func(c *ui.Context) {
		ConversationList(c, &s, items, ConversationListOptions{Now: chatNow, Disabled: true}, nil).Element.Height(560)
	})
	_ = tt.Click("Garden maze")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if s.Selected != "" {
		t.Fatalf("disabled list chose %q", s.Selected)
	}
}

func TestConversationListVirtualizesAndGroupsOnce(t *testing.T) {
	items := make([]ChatConversation, 100000)
	for i := range items {
		items[i] = ChatConversation{ID: fmt.Sprint(i), Title: fmt.Sprintf("Conversation %d", i), Updated: chatNow.Add(-time.Duration(i) * time.Minute)}
	}
	var s ConversationListState
	tt := chatTester(400, 600, func(c *ui.Context) {
		ConversationList(c, &s, items, ConversationListOptions{Now: chatNow}, nil).Element.Height(560)
	})
	if n := len(tt.Texts()); n > 80 {
		t.Fatalf("built %d texts for 100000 conversations", n)
	}
	// A marked plan survives frames that bring no new data.
	s.plan.rows[1].item = -7
	tt.Frame()
	if s.plan.rows[1].item != -7 {
		t.Fatal("regrouped without a data change")
	}
	s.plan.rows[1].item = 0
	// Today, Yesterday, Previous 7 days and Older head the rows.
	if len(s.plan.rows) != len(items)+4 {
		t.Fatalf("%d rows", len(s.plan.rows))
	}
}

func TestConversationListSnapshot(t *testing.T) {
	items := chatConversations()
	snapshot(t, "conversation_list", 360, 420, func(c *ui.Context) {
		s := &ConversationListState{Selected: "b"}
		ui.Column(c).Padding(16).Children(func() {
			ConversationList(c, s, items, ConversationListOptions{Now: chatNow, Label: "History"}, nil).Element.Height(380)
		})
	})
}
