package files

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

var recentNow = time.Date(2026, 10, 5, 15, 0, 0, 0, time.UTC)

func recentSample() []RecentFile {
	return []RecentFile{
		{"/court/old/charter.pdf", recentNow.AddDate(0, -2, 0)},
		{"/court/seal.go", recentNow.Add(-2 * time.Hour)},
		{"/court/notes.md", recentNow.Add(-26 * time.Hour)},
		{"/court/crest.png", recentNow.AddDate(0, 0, -4)},
		{"/court/ledger.csv", recentNow.Add(-10 * time.Minute)},
	}
}

type recentProbe struct {
	files    []RecentFile
	s        RecentFilesState
	disabled bool
	opened   []string
}

func (p *recentProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := RecentFiles(c, p.files, &p.s, RecentFilesOptions{Now: recentNow, Disabled: p.disabled})
	r.Element.Height(420)
	if r.Opened != nil {
		p.opened = append(p.opened, r.Opened.Path)
	}
	if r.Removed != "" {
		for i, f := range p.files {
			if f.Path == r.Removed {
				p.files = append(p.files[:i:i], p.files[i+1:]...)
				break
			}
		}
	}
	if r.Cleared {
		p.files = nil
	}
}

func TestRecentFilesGroups(t *testing.T) {
	p := &recentProbe{files: recentSample()}
	ui.NewTester(p.view, 480, 440)
	want := []recentRow{{0, -1}, {0, 4}, {0, 1}, {1, -1}, {1, 2}, {2, -1}, {2, 3}, {3, -1}, {3, 0}}
	if len(p.s.flat) != len(want) {
		t.Fatalf("rows %v", p.s.flat)
	}
	for i := range want {
		if p.s.flat[i] != want[i] {
			t.Fatalf("rows %v, want %v", p.s.flat, want)
		}
	}
}

func TestRecentFilesOpenRemoveClear(t *testing.T) {
	p := &recentProbe{files: recentSample()}
	tt := ui.NewTester(p.view, 480, 440)
	for _, s := range []string{"Today", "Yesterday", "This week", "Earlier", "14:50", "Oct 1"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	r, _ := tt.Find("seal.go")
	tt.ClickAt(r.X, r.Y+r.H/2)
	tt.ClickAt(r.X, r.Y+r.H/2)
	if len(p.opened) != 1 || p.opened[0] != "/court/seal.go" {
		t.Fatalf("opened %v", p.opened)
	}
	if err := tt.Click("Remove notes.md from recent"); err != nil {
		t.Fatal(err)
	}
	if len(p.files) != 4 || tt.HasText("Yesterday") {
		t.Fatalf("remove left %d files", len(p.files))
	}
	_ = tt.Click("Clear")
	if len(p.files) != 0 || !tt.HasText("No recent files") {
		t.Fatal("clear did not empty the list")
	}
}

func TestRecentFilesKeyboardAndDisabled(t *testing.T) {
	p := &recentProbe{files: recentSample()}
	tt := ui.NewTester(p.view, 480, 440)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(p.opened) != 1 || p.opened[0] != "/court/ledger.csv" {
		t.Fatalf("keyboard opened %v", p.opened)
	}
	p = &recentProbe{files: recentSample(), disabled: true}
	tt = ui.NewTester(p.view, 480, 440)
	_ = tt.Click("Clear")
	if len(p.files) != 5 {
		t.Fatal("disabled list cleared")
	}
}

func TestRecentFilesSnapshot(t *testing.T) {
	snapshot(t, "recent_files", 460, 440, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *recentProbe { return &recentProbe{files: recentSample()} })
		ui.Column(c).Padding(12).Children(func() {
			RecentFiles(c, (*p).files, &(*p).s, RecentFilesOptions{Now: recentNow}).Element.Height(416)
		})
	})
}

func TestRecentFilesLargeAndRevision(t *testing.T) {
	files := make([]RecentFile, 100000)
	for i := range files {
		files[i] = RecentFile{Path: fmt.Sprintf("/f-%d.txt", i), Opened: recentNow}
	}
	s := RecentFilesState{}
	revision := uint64(0)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		RecentFiles(c, files, &s, RecentFilesOptions{Now: recentNow, Revision: revision}).Element.Height(220)
	}, 500, 250)
	n := 0
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "f-") {
			n++
		}
	}
	if n == 0 || n > 60 {
		t.Fatalf("built %d recent files", n)
	}
	files[0].Opened = recentNow.AddDate(0, -1, 0)
	revision++
	tt.Frame()
	if r := s.flat[len(s.flat)-1]; r.group != 3 || r.file != 0 {
		t.Fatalf("revised last row %+v", r)
	}
}

func TestRecentFilesDisablingBeforeDeliveredKey(t *testing.T) {
	p := &recentProbe{files: recentSample()}
	tt := ui.NewTester(p.view, 480, 440)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	before := p.s.selected
	p.disabled = true
	tt.Key(0, ui.KeyDown)
	if p.s.selected != before {
		t.Fatalf("disabled row %d, want %d", p.s.selected, before)
	}
}

func TestRecentFilesRemovedSelectionDoesNotOpenAnother(t *testing.T) {
	p := &recentProbe{files: []RecentFile{{Path: "a.txt", Opened: recentNow}, {Path: "b.txt", Opened: recentNow}}}
	tt := ui.NewTester(p.view, 480, 440)
	if err := tt.Click("b.txt"); err != nil {
		t.Fatal(err)
	}
	p.files = p.files[:1]
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if len(p.opened) != 0 {
		t.Fatalf("opened %v after removing selection", p.opened)
	}
}
