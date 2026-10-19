package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

const (
	mergeBase   = "a\nb\nc\nd\ne\nf"
	mergeOurs   = "a\nB1\nc\nd\ne\nf\ng"
	mergeTheirs = "a\nB2\nc\nD\ne\nf"
)

func TestMergeChunks(t *testing.T) {
	chunks := mergeChunks([]string{"a", "b", "c", "d", "e", "f"}, []string{"a", "B1", "c", "d", "e", "f", "g"}, []string{"a", "B2", "c", "D", "e", "f"})
	// a | b conflict | c | d→D clean | e f | +g clean
	if len(chunks) != 6 {
		t.Fatalf("%d chunks: %+v", len(chunks), chunks)
	}
	if !chunks[1].conflict || chunks[1].ours[0] != "B1" || chunks[1].theirs[0] != "B2" {
		t.Errorf("conflict chunk %+v", chunks[1])
	}
	if chunks[3].conflict || chunks[3].auto[0] != "D" {
		t.Errorf("clean remote change %+v", chunks[3])
	}
	if chunks[5].conflict || len(chunks[5].auto) != 1 || chunks[5].auto[0] != "g" {
		t.Errorf("clean local append %+v", chunks[5])
	}
	same := mergeChunks([]string{"x"}, []string{"y"}, []string{"y"})
	if len(same) != 1 || same[0].conflict {
		t.Errorf("identical changes conflict: %+v", same)
	}
}

func TestThreeWayMergeTakeSides(t *testing.T) {
	var s ThreeWayMergeState
	changes := 0
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := ThreeWayMerge(c, &s, mergeBase, mergeOurs, mergeTheirs, ThreeWayMergeOptions{Disabled: disabled})
		r.Element.Height(560)
		if r.Changed() {
			changes++
		}
	}, 900, 600)
	if !tt.HasText("0 of 1 conflicts resolved") || !tt.HasText("Unresolved") {
		t.Fatalf("start: %q", tt.Texts())
	}
	if got := s.Result(); got != "a\n<<<<<<< ours\nB1\n=======\nB2\n>>>>>>> theirs\nc\nD\ne\nf\ng" {
		t.Fatalf("unresolved result %q", got)
	}
	if err := tt.Click("Take remote"); err != nil {
		t.Fatal(err)
	}
	if s.Result() != "a\nB2\nc\nD\ne\nf\ng" || changes != 1 || !tt.HasText("1 of 1 conflicts resolved") {
		t.Fatalf("take remote: %q changes %d", s.Result(), changes)
	}
	tt.Key(0, ui.KeyTab)
	_ = tt.Click("Take both")
	if s.Result() != "a\nB1\nB2\nc\nD\ne\nf\ng" {
		t.Fatalf("take both: %q", s.Result())
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Keep base")
	if s.Result() != "a\nB1\nB2\nc\nD\ne\nf\ng" || changes != 2 {
		t.Fatalf("disabled merge changed: %q", s.Result())
	}
}

func TestThreeWayMergeSnapshot(t *testing.T) {
	snapshot(t, "three_way_merge", 860, 420, func(c *ui.Context) {
		s := ui.Local(c.Root(), "m", func() *ThreeWayMergeState { return &ThreeWayMergeState{} })
		ui.Column(c).Padding(16).Children(func() {
			ThreeWayMerge(c, *s, mergeBase, mergeOurs, mergeTheirs, ThreeWayMergeOptions{}).Element.Height(380)
		})
	})
}
