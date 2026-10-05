package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

var shortcutGroups = []ShortcutGroup{
	{Title: "General", Bindings: []ShortcutBinding{{"Command palette", Shortcut{ui.Ctrl | ui.Shift, ui.KeyP}}, {"Settings", Shortcut{ui.Ctrl, ui.KeyComma}}}},
	{Title: "Editing", Bindings: []ShortcutBinding{{"Duplicate line", Shortcut{ui.Alt | ui.Shift, ui.KeyDown}}, {"Toggle comment", Shortcut{}}}},
}

func TestKeyboardShortcutsListGroupsAndSearch(t *testing.T) {
	query := ""
	tt := accountTester(500, 500, func(c *ui.Context) { KeyboardShortcutsList(c, &query, shortcutGroups, KeyboardShortcutsListOptions{}) })
	for _, want := range []string{"General", "Editing", "Command palette", "Toggle comment", "—"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	if _, ok := tt.Find("Duplicate line, alt+shift+down"); !ok {
		t.Fatal("row not spoken with its keys")
	}
	_ = tt.Click("Search shortcuts")
	tt.Type("comment")
	if tt.HasText("General") || !tt.HasText("Toggle comment") {
		t.Fatalf("search: %q", tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("ctrl+,")
	if !tt.HasText("Settings") || tt.HasText("Editing") {
		t.Fatalf("key search: %q", tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("zzz")
	if !tt.HasText("No results") {
		t.Fatalf("empty: %q", tt.Texts())
	}
}

func TestKeyboardShortcutsListSnapshot(t *testing.T) {
	snapshot(t, "keyboard_shortcuts_list", 460, 320, func(c *ui.Context) {
		q := ui.Local(c.Root(), "q", func() string { return "" })
		ui.Column(c).Padding(16).Children(func() { KeyboardShortcutsList(c, q, shortcutGroups, KeyboardShortcutsListOptions{}) })
	})
}

func TestKeyboardShortcutsListLargeDataBuildsVisibleRows(t *testing.T) {
	items := make([]ShortcutBinding, 100000)
	for i := range items {
		items[i] = ShortcutBinding{Action: fmt.Sprintf("Action %d", i)}
	}
	query := ""
	tt := accountTester(500, 500, func(c *ui.Context) {
		KeyboardShortcutsList(c, &query, []ShortcutGroup{{Title: "General", Bindings: items}}, KeyboardShortcutsListOptions{})
	})
	if n := len(tt.Texts()); n > 100 {
		t.Fatalf("built %d text elements for a bounded viewport", n)
	}
	if err := tt.Click("Search shortcuts"); err != nil {
		t.Fatal(err)
	}
	tt.Type("99999")
	if !tt.HasText("Action 99999") || tt.HasText("Action 99998") {
		t.Fatal("large-list search does not reach last binding")
	}
}

func TestKeyboardShortcutsListDisabledSearchAndInvalidBinding(t *testing.T) {
	q := ""
	tt := accountTester(600, 500, func(c *ui.Context) {
		ui.Column(c).Disabled(true).Children(func() { KeyboardShortcutsList(c, &q, shortcutGroups, KeyboardShortcutsListOptions{}) })
	})
	if err := tt.Click("Search shortcuts"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Type("save")
	if q != "" {
		t.Fatal("disabled search edits")
	}
	accountPanics(t, func(c *ui.Context) {
		KeyboardShortcutsList(c, &q, []ShortcutGroup{{Bindings: []ShortcutBinding{{}}}}, KeyboardShortcutsListOptions{})
	})
}

func TestKeyboardShortcutsListCachedPlanSeesInPlaceChanges(t *testing.T) {
	st := shortcutsListState{}
	groups := []ShortcutGroup{{Title: "General", Bindings: []ShortcutBinding{{Action: "Save"}}}}
	shortcutsListPlan(&st, groups, "save")
	if len(st.rows) != 2 {
		t.Fatal("missing header or binding")
	}
	before := &st.rows[0]
	shortcutsListPlan(&st, groups, "save")
	if before != &st.rows[0] {
		t.Fatal("unchanged data rebuilt shortcut plan")
	}
	groups[0].Bindings[0].Action = "Open"
	shortcutsListPlan(&st, groups, "save")
	if len(st.rows) != 0 {
		t.Fatal("cached list misses in-place binding change")
	}
}
