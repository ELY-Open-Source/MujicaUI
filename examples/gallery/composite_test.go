package main

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// compositeTester opens the gallery on one composition at 1200×800.
func compositeTester(id string) (*gallery, *ui.Tester) {
	g := newGallery()
	g.selected = id
	return g, ui.NewTester(g.view, 1200, 800)
}

func click(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatal(err)
	}
}

// await renders frames until done holds or 5 seconds pass.
func await(tt *ui.Tester, done func() bool) {
	for deadline := time.Now().Add(5 * time.Second); !done() && time.Now().Before(deadline); {
		tt.Advance(20 * time.Millisecond)
	}
}

func TestCompositeSettingsSaveAndReset(t *testing.T) {
	g, tt := compositeTester(settingsID)
	s := &g.cSettings
	click(t, tt, "Appearance")
	click(t, tt, "Reduce motion")
	if !s.draft.reduce || s.saved.reduce {
		t.Fatalf("draft reduce %v, saved reduce %v", s.draft.reduce, s.saved.reduce)
	}
	if !tt.HasText("Unsaved changes") {
		t.Fatalf("no unsaved flag in %q", tt.Texts())
	}
	click(t, tt, "Save")
	if s.saved != s.draft || !s.saved.reduce {
		t.Fatalf("saved %+v, draft %+v", s.saved, s.draft)
	}
	if !tt.HasText("Settings saved") || !tt.HasText("All changes saved") {
		t.Fatalf("no saved toast in %q", tt.Texts())
	}
	click(t, tt, "Reset")
	if s.draft != settingsInitial || !s.saved.reduce {
		t.Fatalf("reset draft %+v, saved %+v", s.draft, s.saved)
	}
	if !tt.HasText("Unsaved changes") {
		t.Fatal("reset away from the saved values is not flagged")
	}
}

func TestCompositeFilesFilterSelectDelete(t *testing.T) {
	g, tt := compositeTester(filesID)
	s := &g.cFiles
	click(t, tt, "Filter files")
	tt.Type("rehearsal")
	if s.query != "rehearsal" || s.table.Rows() != 1 {
		t.Fatalf("query %q shows %d rows", s.query, s.table.Rows())
	}
	click(t, tt, "Select row")
	if s.table.Selection.Len() != 1 || !s.table.Selection.Has(2) {
		t.Fatalf("selection %v", slices.Collect(s.table.Selection.All()))
	}
	click(t, tt, "Delete")
	if !s.deleting || !slices.Equal(s.doomed, []int{2}) {
		t.Fatalf("deleting %v, doomed %v", s.deleting, s.doomed)
	}
	click(t, tt, "Delete permanently")
	if len(s.entries) != 7 || slices.ContainsFunc(s.entries, func(f fileEntry) bool { return f.ID == 2 }) {
		t.Fatalf("entries after delete %+v", s.entries)
	}
	if s.deleting || s.table.Selection.Len() != 0 {
		t.Fatalf("deleting %v, selection %d", s.deleting, s.table.Selection.Len())
	}
	if !tt.HasText("Nothing matches “rehearsal”.") {
		t.Fatalf("no empty search state in %q", tt.Texts())
	}
}

func TestCompositeFormInvalidThenValid(t *testing.T) {
	g, tt := compositeTester(formID)
	s := &g.cForm
	click(t, tt, "Submit entry")
	if got := s.st.Error("name"); got != "Required" {
		t.Fatalf("name error %q", got)
	}
	if got := s.st.Error("terms"); got != "Accept the terms to continue." {
		t.Fatalf("terms error %q", got)
	}
	if !tt.Focused("Stage name") {
		t.Fatal("the first invalid field has no focus")
	}
	tt.Type("Mortis")
	click(t, tt, "Password")
	tt.Type("hunter22secret")
	click(t, tt, "Debut date") // the FormItem label now reaches the picker's text box
	tt.Type("2026-12-24")
	tt.Key(0, ui.KeyEnter)
	click(t, tt, "Roles")
	tt.Type("guitar")
	tt.Key(0, ui.KeyEnter)
	click(t, tt, "I accept the covenant of the masquerade")
	if s.name != "Mortis" || s.password != "hunter22secret" || s.debut.Format("2006-01-02") != "2026-12-24" ||
		!slices.Equal(s.tags, []string{"guitar"}) || !s.terms {
		t.Fatalf("values %q, password ok %v, %v %q %v", s.name, s.password == "hunter22secret", s.debut, s.tags, s.terms)
	}
	click(t, tt, "Submit entry")
	await(tt, func() bool { return tt.HasText("Submitting…") }) // the name check finishes first
	if !s.st.Submitting() {
		t.Fatalf("no loading overlay in %q", tt.Texts())
	}
	click(t, tt, "Submit entry") // the overlay blocks a second submit
	await(tt, func() bool { return s.phase == formSucceeded })
	if s.phase != formSucceeded || !tt.HasText("Entry received") {
		t.Fatalf("phase %d, texts %q", s.phase, tt.Texts())
	}
	submits := 0
	for i, e := range g.events {
		if strings.Contains(e.text, "hunter22secret") {
			t.Fatalf("password logged in event %d of page %s", i, e.page)
		}
		if e.page == formID && strings.HasPrefix(e.text, "submitted ") {
			submits++
		}
	}
	if submits != 1 {
		t.Fatalf("%d submits, want 1", submits)
	}
}
