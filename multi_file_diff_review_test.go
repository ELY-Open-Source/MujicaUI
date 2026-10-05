package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func reviewFiles() []FileChange {
	return []FileChange{
		fileChangeSample,
		{Path: "court/menu.md", Old: "soup\nroast\n", New: "soup\nroast\ntart\n"},
		{Path: "court/old.md", Old: "gone\n"},
	}
}

func reviewTester(d *[]FileDecision, files []FileChange, opts MultiFileDiffReviewOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(8).Children(func() {
			if MultiFileDiffReview(c, d, files, opts).Changed() {
				*changes++
			}
		})
	}, 760, 520)
}

func TestMultiFileDiffReviewPerFileAndAll(t *testing.T) {
	files := reviewFiles()
	d := make([]FileDecision, len(files))
	changes := 0
	tt := reviewTester(&d, files, MultiFileDiffReviewOptions{}, &changes)
	if !tt.HasText("0 of 3 files reviewed") || !tt.HasText("Guests: 42") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("court/menu.md")
	if !tt.HasText("tart") || tt.HasText("Guests: 42") {
		t.Fatal("choosing a file did not show its diff")
	}
	_ = tt.Click("Reject")
	if d[1] != FileRejected || changes != 1 || !tt.HasText("1 of 3 files reviewed") {
		t.Fatalf("reject: %v", d)
	}
	_ = tt.Click("Accept all")
	if d[0] != FileAccepted || d[1] != FileRejected || d[2] != FileAccepted || !tt.HasText("3 of 3 files reviewed") {
		t.Fatalf("accept all: %v", d)
	}
}

func TestMultiFileDiffReviewKeyboard(t *testing.T) {
	files := reviewFiles()
	d := make([]FileDecision, len(files))
	n := 0
	tt := reviewTester(&d, files, MultiFileDiffReviewOptions{}, &n)
	_ = tt.Click("court/ledger.md")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("gone") || tt.HasText("Guests: 42") {
		t.Fatalf("keyboard choice: %q", tt.Texts())
	}
}

func TestMultiFileDiffReviewDisabledAndEmpty(t *testing.T) {
	files := reviewFiles()
	d := make([]FileDecision, len(files))
	n := 0
	tt := reviewTester(&d, files, MultiFileDiffReviewOptions{Disabled: true}, &n)
	_ = tt.Click("Accept all")
	_ = tt.Click("Accept")
	_ = tt.Click("court/menu.md")
	if n != 0 || d[0] != FilePending || tt.HasText("tart") {
		t.Fatalf("disabled review acted: %v", d)
	}
	var none []FileDecision
	tt = reviewTester(&none, nil, MultiFileDiffReviewOptions{}, &n)
	if !tt.HasText("No changes to review") {
		t.Fatalf("empty: %q", tt.Texts())
	}
}

func TestMultiFileDiffReviewManyFilesVirtualized(t *testing.T) {
	files := make([]FileChange, 100000)
	for i := range files {
		files[i] = FileChange{Path: fmt.Sprintf("f%06d.txt", i), Old: "a\n", New: "b\n"}
	}
	d := make([]FileDecision, len(files))
	n := 0
	tt := reviewTester(&d, files, MultiFileDiffReviewOptions{}, &n)
	if got := len(tt.Texts()); got > 400 {
		t.Fatalf("%d texts built for 100000 files", got)
	}
	if !tt.HasText("f000000.txt") {
		t.Fatal("first file missing")
	}
}

func TestMultiFileDiffReviewNeedsMatchingDecisions(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	d := []FileDecision{}
	n := 0
	reviewTester(&d, reviewFiles(), MultiFileDiffReviewOptions{}, &n)
}

func TestMultiFileDiffReviewSnapshot(t *testing.T) {
	snapshot(t, "multi_file_diff_review", 760, 420, func(c *ui.Context) {
		files := reviewFiles()
		d := []FileDecision{FilePending, FileAccepted, FileRejected}
		ui.Column(c).Padding(12).Children(func() {
			MultiFileDiffReview(c, &d, files, MultiFileDiffReviewOptions{Height: 390})
		})
	})
}
