package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var fileChangeSample = FileChange{Path: "court/ledger.md",
	Old: "# Ledger\nGuests: 40\nWine: claret\nMusic: none\n",
	New: "# Ledger\nGuests: 42\nWine: claret\nMusic: quartet\nFlowers: lilies\n"}

func fileChangeTester(d *FileDecision, ch FileChange, opts FileChangeCardOptions, got *[]FileChangeResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := FileChangeCard(c, d, ch, opts)
		if r.Decided || r.OpenDiff {
			*got = append(*got, r)
		}
	}, 560, 400)
}

func TestFileChangeCardStatsPreviewAndDecide(t *testing.T) {
	d := FilePending
	var got []FileChangeResult
	tt := fileChangeTester(&d, fileChangeSample, FileChangeCardOptions{}, &got)
	for _, s := range []string{"court/ledger.md", "+3", "−2", "Pending", "Guests: 40", "Guests: 42", "Flowers: lilies"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	_ = tt.Click("Open diff")
	if len(got) != 1 || !got[0].OpenDiff || d != FilePending {
		t.Fatalf("open diff: %+v", got)
	}
	_ = tt.Click("Accept")
	if d != FileAccepted || !got[1].Decided || !tt.HasText("Accepted") {
		t.Fatalf("accept: %v", d)
	}
	_ = tt.Click("Undo")
	if d != FilePending {
		t.Fatalf("undo: %v", d)
	}
	tt = fileChangeTester(&d, fileChangeSample, FileChangeCardOptions{}, &got)
	for range 3 {
		tt.Key(0, ui.KeyTab)
		if tt.Focused("Reject") {
			break
		}
	}
	tt.Key(0, ui.KeyEnter)
	if d != FileRejected || !tt.HasText("Rejected") {
		t.Fatalf("keyboard reject: %v", d)
	}
}

func TestFileChangeCardPreviewLimitAndNewFile(t *testing.T) {
	d := FilePending
	var got []FileChangeResult
	tt := fileChangeTester(&d, FileChange{Path: "a.txt", New: "1\n2\n3\n4\n5\n"}, FileChangeCardOptions{Preview: 2}, &got)
	if !tt.HasText("New file") || !tt.HasText("+5") || !tt.HasText("2") || tt.HasText("3") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestFileChangeCardDisabled(t *testing.T) {
	d := FilePending
	var got []FileChangeResult
	tt := fileChangeTester(&d, fileChangeSample, FileChangeCardOptions{Disabled: true}, &got)
	_ = tt.Click("Accept")
	_ = tt.Click("Open diff")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(got) != 0 || d != FilePending {
		t.Fatalf("disabled card acted: %+v %v", got, d)
	}
}

func TestFileChangeCardSnapshot(t *testing.T) {
	snapshot(t, "file_change_card", 520, 330, func(c *ui.Context) {
		d := FilePending
		ui.Column(c).Padding(16).Children(func() { FileChangeCard(c, &d, fileChangeSample, FileChangeCardOptions{}) })
	})
}
