package files

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
)

func archiveSample() []ArchiveEntry {
	return []ArchiveEntry{
		{Path: "court/seal.go", Size: 4000, Packed: 1000},
		{Path: "court/decrees/42.pdf", Size: 6000, Packed: 5000},
		{Path: "README.md", Size: 1000, Packed: 400},
		{Path: "empty/", Dir: true},
	}
}

type archiveProbe struct {
	entries []ArchiveEntry
	s       ArchiveViewerState
	opts    ArchiveViewerOptions
	opened  []string
}

func (p *archiveProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := ArchiveViewer(c, p.entries, &p.s, p.opts)
	r.Element.Height(340)
	if r.Opened != nil {
		p.opened = append(p.opened, r.Opened.Path)
	}
}

func TestArchiveViewerIndex(t *testing.T) {
	s := &ArchiveViewerState{}
	archiveBuild(archiveSample(), s)
	ix := s.index
	if ix.files != 3 || ix.size != 11000 || ix.pack != 6400 {
		t.Fatalf("totals %d %d %d", ix.files, ix.size, ix.pack)
	}
	if got := ix.sums["court"]; got != [2]int64{10000, 6000} {
		t.Fatalf("court sums %v", got)
	}
	if len(ix.roots) != 3 || len(ix.kids["court"]) != 2 || ix.kids["empty"] == nil {
		t.Fatalf("roots %v kids %v", ix.roots, ix.kids)
	}
	if archiveRatio(4000, 1000) != "75%" || archiveRatio(0, 0) != "—" {
		t.Fatal("ratio")
	}
}

func TestArchiveViewerExpandAndOpen(t *testing.T) {
	p := &archiveProbe{entries: archiveSample()}
	tt := ui.NewTester(p.view, 620, 360)
	if !tt.HasText("3 files, 11 KB unpacked, 6 KB packed (42% smaller)") {
		t.Fatalf("summary: %q", tt.Texts())
	}
	if !tt.HasText("40%") || tt.HasText("seal.go") {
		t.Fatalf("court row: %q", tt.Texts())
	}
	cr, _ := tt.Find("10 KB")
	tt.ClickAt(cr.X, cr.Y+cr.H/2)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("seal.go") || !tt.HasText("75%") {
		t.Fatalf("expand: %q", tt.Texts())
	}
	r, _ := tt.Find("seal.go")
	tt.ClickAt(r.X, r.Y+r.H/2)
	tt.ClickAt(r.X, r.Y+r.H/2)
	if len(p.opened) != 1 || p.opened[0] != "court/seal.go" {
		t.Fatalf("opened %v", p.opened)
	}
}

func TestArchiveViewerStatusAndInvalid(t *testing.T) {
	tt := ui.NewTester((&archiveProbe{opts: ArchiveViewerOptions{Status: data.DataFailed, Error: "bad archive"}}).view, 620, 360)
	if !tt.HasText("bad archive") {
		t.Fatalf("texts %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate entry did not panic")
		}
	}()
	archiveBuild([]ArchiveEntry{{Path: "a"}, {Path: "a"}}, &ArchiveViewerState{})
}

func TestArchiveViewerSnapshot(t *testing.T) {
	snapshot(t, "archive_viewer", 600, 320, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *archiveProbe {
			p := &archiveProbe{entries: archiveSample()}
			p.s.table.SetOpen("court", true)
			return p
		})
		ui.Column(c).Padding(12).Children(func() {
			ArchiveViewer(c, (*p).entries, &(*p).s, ArchiveViewerOptions{}).Element.Height(296)
		})
	})
}

func TestArchiveViewerRejectsFileParent(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("file used as a folder was accepted")
		}
	}()
	archiveBuild([]ArchiveEntry{{Path: "court"}, {Path: "court/a.txt"}}, &ArchiveViewerState{})
}
func TestArchiveViewerAncestorDisabled(t *testing.T) {
	s := ArchiveViewerState{}
	disabled, opened := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			r := ArchiveViewer(c, []ArchiveEntry{{Path: "a.txt", Size: 10}}, &s, ArchiveViewerOptions{})
			r.Element.Height(200)
			if r.Opened != nil {
				opened++
			}
		})
	}, 500, 250)
	if err := tt.Click("a.txt"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Key(0, ui.KeyEnter)
	if opened != 0 {
		t.Fatal("disabled archive opened")
	}
}

func TestArchiveViewerLargeAndRevision(t *testing.T) {
	entries := make([]ArchiveEntry, 100000)
	for i := range entries {
		entries[i] = ArchiveEntry{Path: fmt.Sprintf("f-%d.txt", i), Size: 10, Packed: 5}
	}
	s := ArchiveViewerState{}
	revision := uint64(0)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ArchiveViewer(c, entries, &s, ArchiveViewerOptions{Revision: revision}).Element.Height(220)
	}, 600, 250)
	n := 0
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "f-") {
			n++
		}
	}
	if n == 0 || n > 60 {
		t.Fatalf("built %d archive entries", n)
	}
	entries[0].Size = 20
	revision++
	tt.Frame()
	if s.index.size != 1000010 {
		t.Fatalf("revised size %d", s.index.size)
	}
}

func TestArchiveViewerTotalOverflow(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("overflowing archive total was accepted")
		}
	}()
	archiveBuild([]ArchiveEntry{{Path: "a", Size: 1 << 62}, {Path: "b", Size: 1 << 62}}, &ArchiveViewerState{})
}

func TestArchiveViewerDisablingBeforeDeliveredKey(t *testing.T) {
	s := ArchiveViewerState{}
	disabled := false
	entries := []ArchiveEntry{{Path: "a.txt", Size: 1}, {Path: "b.txt", Size: 1}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() { ArchiveViewer(c, entries, &s, ArchiveViewerOptions{}).Element.Height(200) })
	}, 600, 250)
	if err := tt.Click("a.txt"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Key(0, ui.KeyDown)
	if key, _ := s.table.Selected(); key != "a.txt" {
		t.Fatalf("disabled selection %s", key)
	}
}
