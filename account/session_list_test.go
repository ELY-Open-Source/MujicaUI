package account

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func sessionSample(now time.Time) []Session {
	return []Session{
		{ID: "s1", Device: "MacBook Pro · Safari", Location: "Tokyo", LastActive: now, Current: true},
		{ID: "s2", Device: "iPhone · App", Location: "Osaka", LastActive: time.Date(2026, 10, 1, 9, 30, 0, 0, time.Local)},
		{ID: "s3", Device: "Windows · Edge"},
	}
}

func TestSessionListRevokeOneAndOthers(t *testing.T) {
	var revoked []string
	others := 0
	disabled := false
	tt := accountTester(600, 500, func(c *ui.Context) {
		r := SessionList(c, sessionSample(c.Now()), SessionListOptions{Disabled: disabled})
		if id, ok := r.Revoked(); ok {
			revoked = append(revoked, id)
		}
		if r.OthersRevoked() {
			others++
		}
	})
	for _, want := range []string{"This device", "Tokyo · Active now", "Osaka · 2026-10-01 09:30", "Last active unknown"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	if _, ok := tt.Find("Sign out of MacBook Pro · Safari"); ok {
		t.Fatal("current session offers sign out")
	}
	if err := tt.Click("Sign out of iPhone · App"); err != nil {
		t.Fatal(err)
	}
	_ = tt.Click("Sign out of other sessions")
	if !tt.Focused("Cancel") || !tt.HasText("2 other sessions will need to sign in again.") {
		t.Fatalf("confirm: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if others != 0 {
		t.Fatal("Enter on Cancel signed out")
	}
	_ = tt.Click("Sign out of other sessions")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(revoked) != 1 || revoked[0] != "s2" || others != 1 {
		t.Fatalf("revoked %v others %d", revoked, others)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Sign out of Windows · Edge")
	if len(revoked) != 1 {
		t.Fatal("disabled list signed out")
	}
}

func TestSessionListSnapshot(t *testing.T) {
	snapshot(t, "session_list", 560, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { SessionList(c, sessionSample(c.Now()), SessionListOptions{}) })
	})
}

func TestSessionListLargeDataBuildsVisibleRows(t *testing.T) {
	items := make([]Session, 100000)
	for i := range items {
		items[i] = Session{ID: fmt.Sprint(i), Device: fmt.Sprintf("Device %d", i), Current: i == 0}
	}
	tt := accountTester(600, 500, func(c *ui.Context) { SessionList(c, items, SessionListOptions{}) })
	if n := len(tt.Texts()); n > 100 {
		t.Fatalf("built %d text elements for a bounded viewport", n)
	}
}
