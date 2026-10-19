package git

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
)

const diffOld = "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hi\")\n}\n"
const diffNew = "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tname := \"court\"\n\tfmt.Println(\"hi\", name)\n}\n"

func TestDiffViewerRowsFoldAndPair(t *testing.T) {
	var a, b []string
	for i := range 30 {
		a = append(a, fmt.Sprintf("line %d", i))
	}
	b = append(b, a...)
	b[15] = "changed"
	d := internal.LineDiff(a, b)
	rows := diffViewerRows(d, 3, map[int]bool{}, false)
	// 12 hidden, 3 context, delete, insert, 3 context, 11 hidden.
	if len(rows) != 10 || rows[0].hidden != 12 || rows[9].hidden != 11 {
		t.Fatalf("unified rows %+v", rows)
	}
	split := diffViewerRows(d, 3, map[int]bool{}, true)
	if len(split) != 9 || split[4].left < 0 || split[4].right < 0 || d[split[4].left].Op != internal.DiffDelete {
		t.Fatalf("split rows %+v", split)
	}
	open := diffViewerRows(d, 3, map[int]bool{0: true}, false)
	if len(open) != 21 || open[0].hidden != 0 {
		t.Fatalf("expanded rows %d", len(open))
	}
}

func TestDiffViewerUnfoldAndModes(t *testing.T) {
	var s DiffViewerState
	from := strings.Repeat("same\n", 20) + "old\n"
	to := strings.Repeat("same\n", 20) + "new\n"
	split := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		DiffViewer(c, &s, from, to, DiffViewerOptions{SideBySide: split}).Height(760)
	}, 700, 800)
	if !tt.HasText("17 unchanged lines") {
		t.Fatalf("fold: %q", tt.Texts())
	}
	if _, ok := tt.Find("Deleted: old"); !ok {
		t.Fatal("deleted line not named")
	}
	if _, ok := tt.Find("Added: new"); !ok {
		t.Fatal("added line not named")
	}
	if err := tt.Click("17 unchanged lines"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("17 unchanged lines") || !tt.HasText("17") {
		t.Fatalf("unfold: %q", tt.Texts())
	}
	split = true
	tt.Frame()
	if !tt.HasText("−") || !tt.HasText("+") {
		t.Fatalf("split: %q", tt.Texts())
	}
}

func TestDiffViewerVirtualized(t *testing.T) {
	var a, b strings.Builder
	for i := range 100000 {
		fmt.Fprintf(&a, "line %d\n", i)
		if i%1000 == 0 {
			fmt.Fprintf(&b, "edit %d\n", i)
		} else {
			fmt.Fprintf(&b, "line %d\n", i)
		}
	}
	from, to := a.String(), b.String()
	var s DiffViewerState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		DiffViewer(c, &s, from, to, DiffViewerOptions{}).Height(400)
	}, 700, 440)
	if n := len(tt.Texts()); n > 400 {
		t.Fatalf("%d texts for 100,000 lines", n)
	}
	rows := &s.rows[0]
	tt.Frame()
	if &s.rows[0] != rows {
		t.Fatal("rows rebuilt without a change")
	}
}

func TestDiffViewerSnapshot(t *testing.T) {
	snapshot(t, "diff_viewer", 640, 520, func(c *ui.Context) {
		s := ui.Local(c.Root(), "u", func() *DiffViewerState { return &DiffViewerState{} })
		p := ui.Local(c.Root(), "p", func() *DiffViewerState { return &DiffViewerState{} })
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			DiffViewer(c, *s, diffOld, diffNew, DiffViewerOptions{Language: "go", Context: 1}).Height(220)
			DiffViewer(c, *p, diffOld, diffNew, DiffViewerOptions{Language: "go", SideBySide: true, Context: 1}).Height(220)
		})
	})
}
