package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var releases = []Release{
	{Version: "1.1", Date: "2026-10-05", Changes: []ReleaseChange{{ChangeNew, "500 components."}, {ChangeFixed, "Snapshots no longer drift."}}},
	{Version: "1.0", Changes: []ReleaseChange{{ChangeImproved, "Court-gothic theme."}}},
}

func TestWhatsNewDialogListsAndHides(t *testing.T) {
	open, hide, n := true, false, 0
	tt := accountTester(700, 600, func(c *ui.Context) {
		if WhatsNewDialog(c, &open, &hide, releases, WhatsNewDialogOptions{}).Changed() {
			n++
		}
	})
	for _, want := range []string{"What's new", "1.1", "2026-10-05", "New", "Fixed", "Improved", "Court-gothic theme."} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	_ = tt.Click("Don't show again")
	if !hide || n != 1 {
		t.Fatal("hide not set")
	}
	_ = tt.Click("Got it")
	if open {
		t.Fatal("not closed")
	}
	open = true
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	if open || !hide {
		t.Fatal("escape did not close or lost hide")
	}
}

func TestWhatsNewDialogSnapshot(t *testing.T) {
	snapshot(t, "whats_new_dialog", 560, 460, func(c *ui.Context) {
		open := ui.Local(c.Root(), "open", func() bool { return true })
		hide := ui.Local(c.Root(), "hide", func() bool { return false })
		WhatsNewDialog(c, open, hide, releases, WhatsNewDialogOptions{})
	})
}

func TestWhatsNewDialogDisabledParentCloses(t *testing.T) {
	open, hide, disabled := true, false, false
	tt := accountTester(700, 600, func(c *ui.Context) {
		ui.Column(c).Disabled(disabled).Children(func() { WhatsNewDialog(c, &open, &hide, releases, WhatsNewDialogOptions{}) })
	})
	disabled = true
	tt.Frame()
	if open || hide || tt.HasText("What's new") {
		t.Fatal("disabled parent retains dialog")
	}
}

func TestWhatsNewDialogInvalidChange(t *testing.T) {
	open, hide := true, false
	accountPanics(t, func(c *ui.Context) {
		WhatsNewDialog(c, &open, &hide, []Release{{Version: "1", Changes: []ReleaseChange{{Kind: ChangeKind(99), Text: "Change"}}}}, WhatsNewDialogOptions{})
	})
}
