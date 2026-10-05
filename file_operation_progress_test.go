package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

type opProbe struct {
	op       FileOperation
	disabled bool
	canceled int
	answers  []ConflictChoice
	all      []bool
}

func (p *opProbe) view(c *ui.Context) {
	Use(c, Settings{})
	r := FileOperationProgress(c, p.op, FileOperationOptions{Disabled: p.disabled})
	r.Element.Width(480)
	if r.Canceled {
		p.canceled++
	}
	if r.Resolved != ConflictNone {
		p.answers = append(p.answers, r.Resolved)
		p.all = append(p.all, r.ApplyAll)
	}
}

func TestFileOperationProgressShowsCounts(t *testing.T) {
	p := &opProbe{op: FileOperation{Kind: FileMove, Done: 3, Total: 12, BytesDone: 12 << 20, BytesTotal: 48 << 20, Current: "decrees/42.pdf"}}
	tt := ui.NewTester(p.view, 520, 300)
	for _, s := range []string{"Moving 3 of 12 items", "25%", "12.0 MB of 48.0 MB", "decrees/42.pdf"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	if err := tt.Click("Cancel"); err != nil {
		t.Fatal(err)
	}
	if p.canceled != 1 {
		t.Fatal("cancel not reported")
	}
}

func TestFileOperationProgressConflict(t *testing.T) {
	p := &opProbe{op: FileOperation{Kind: FileCopy, Done: 1, Total: 4, Conflict: "seal.go"}}
	tt := ui.NewTester(p.view, 520, 300)
	if !tt.HasText("“seal.go” already exists here.") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Keep both")
	_ = tt.Click("Do this for all conflicts")
	_ = tt.Click("Replace")
	if len(p.answers) != 2 || p.answers[0] != ConflictKeepBoth || p.all[0] || p.answers[1] != ConflictReplace || !p.all[1] {
		t.Fatalf("answers %v all %v", p.answers, p.all)
	}
	// Keyboard: the answers are buttons in the Tab order.
	for range 5 {
		tt.Key(0, ui.KeyTab)
	}
	tt.Key(0, ui.KeyEnter)
	if len(p.answers) != 3 {
		t.Fatalf("keyboard answer missing: %v", p.answers)
	}
}

func TestFileOperationProgressDisabledAndInvalid(t *testing.T) {
	p := &opProbe{op: FileOperation{Kind: FileDelete, Total: 2, Conflict: "a"}, disabled: true}
	tt := ui.NewTester(p.view, 520, 300)
	_ = tt.Click("Cancel")
	_ = tt.Click("Skip")
	if p.canceled != 0 || len(p.answers) != 0 {
		t.Fatal("disabled panel answered")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("Done above Total did not panic")
		}
	}()
	ui.NewTester((&opProbe{op: FileOperation{Done: 5, Total: 2}}).view, 520, 300)
}

func TestFileOperationProgressSnapshot(t *testing.T) {
	snapshot(t, "file_operation_progress", 520, 300, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			FileOperationProgress(c, FileOperation{Kind: FileCopy, Done: 4, Total: 10, BytesDone: 30 << 20, BytesTotal: 80 << 20,
				Current: "chronicles/first.md", Conflict: "first.md"}, FileOperationOptions{})
		})
	})
}

func TestFileOperationProgressNegativeByteTotal(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("negative byte total was accepted")
		}
	}()
	ui.NewTester((&opProbe{op: FileOperation{BytesTotal: -1}}).view, 520, 300)
}
