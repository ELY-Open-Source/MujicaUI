package files

import (
	"fmt"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestFileGridFreshSliceKeepsRangeAndScroll(t *testing.T) {
	entries := fileGridEntries(6)
	var s FileGridState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		FileGrid(c, slices.Clone(entries), &s, FileGridOptions{}).Element.Height(320)
	}, 600, 340)
	if err := tt.Click("crest-3.png"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if err := tt.ClickWith(ui.Shift, "crest-5.png"); err != nil {
		t.Fatal(err)
	}
	if got := s.Selected(); len(got) != 3 {
		t.Fatalf("range selection %v", got)
	}
	entries = fileGridEntries(100000)
	tt.Frame()
	tt.Scroll(300, 200, 0, 10000)
	tt.Frame()
	if tt.HasText("sketches") {
		t.Fatal("unchanged entries reset the scroll position")
	}
}

func fileGridEntries(n int) []FileEntry {
	out := make([]FileEntry, n)
	for i := range out {
		name := fmt.Sprintf("crest-%d.png", i)
		out[i] = FileEntry{Path: "/art/" + name, Name: name}
	}
	out[0] = FileEntry{Path: "/art/sketches", Name: "sketches", Dir: true}
	return out
}

type gridProbe struct {
	entries  []FileEntry
	s        FileGridState
	disabled bool
	thumbs   int
	opened   []string
}

func (p *gridProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := FileGrid(c, p.entries, &p.s, FileGridOptions{Disabled: p.disabled, Thumbnail: func(e FileEntry) bool {
		if FileKindOf(e.Name) != FileImage {
			return false
		}
		p.thumbs++
		ui.Box(c).Fill().Background(core.Tokens(c).Info.Alpha(0.25))
		return true
	}})
	r.Element.Height(320)
	if r.Opened != nil {
		p.opened = append(p.opened, r.Opened.Name)
	}
}

func TestFileGridMultiSelectAndOpen(t *testing.T) {
	p := &gridProbe{entries: fileGridEntries(6)}
	tt := ui.NewTester(p.view, 600, 340)
	if err := tt.Click("crest-1.png"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ClickWith(ui.Cmd, "crest-3.png"); err != nil {
		t.Fatal(err)
	}
	if got := p.s.Selected(); len(got) != 2 {
		t.Fatalf("selection %v", got)
	}
	r, _ := tt.Find("sketches")
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if len(p.opened) != 1 || p.opened[0] != "sketches" {
		t.Fatalf("opened %v", p.opened)
	}
}

func TestFileGridThumbnailSlotAndVirtual(t *testing.T) {
	p := &gridProbe{entries: fileGridEntries(100_000)}
	tt := ui.NewTester(p.view, 600, 340)
	p.thumbs = 0
	tt.Frame()
	if p.thumbs == 0 || p.thumbs > 60 {
		t.Fatalf("built %d thumbnails of 100000", p.thumbs)
	}
	if _, ok := tt.Find("Folder"); !ok {
		t.Fatal("folder without thumbnail shows no icon")
	}
}

func TestFileGridKeysAndDisabled(t *testing.T) {
	p := &gridProbe{entries: fileGridEntries(6)}
	tt := ui.NewTester(p.view, 600, 340)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(p.opened) != 1 {
		t.Fatalf("keyboard opened %v", p.opened)
	}
	p = &gridProbe{entries: fileGridEntries(6), disabled: true}
	tt = ui.NewTester(p.view, 600, 340)
	_ = tt.Click("crest-1.png")
	if len(p.s.Selected()) != 0 {
		t.Fatal("disabled grid selected")
	}
}

func TestFileGridSnapshot(t *testing.T) {
	snapshot(t, "file_grid", 560, 340, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *gridProbe {
			p := &gridProbe{entries: append(fileGridEntries(3), FileEntry{Path: "/art/deed.pdf", Name: "deed.pdf"}, FileEntry{Path: "/art/x.zip", Name: "archive.zip"})}
			p.s.list.Selection.Add("/art/crest-1.png")
			return p
		})
		ui.Column(c).Padding(12).Children(func() { FileGrid(c, (*p).entries, &(*p).s, FileGridOptions{}).Element.Height(316) })
	})
}

func TestFileGridRejectsSmallWidth(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("width below hit size was accepted")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		FileGrid(c, fileGridEntries(3), &FileGridState{}, FileGridOptions{ItemWidth: 1})
	}, 400, 300)
}

func TestFileGridDisablingBeforeDeliveredKey(t *testing.T) {
	p := &gridProbe{entries: fileGridEntries(6)}
	tt := ui.NewTester(p.view, 600, 340)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	before := p.s.Selected()
	p.disabled = true
	tt.Key(0, ui.KeyDown)
	after := p.s.Selected()
	if len(after) != len(before) || after[0] != before[0] {
		t.Fatalf("disabled selection %v, want %v", after, before)
	}
}

func TestFileGridShrinkThenShiftSelect(t *testing.T) {
	p := &gridProbe{entries: fileGridEntries(6)}
	tt := ui.NewTester(p.view, 600, 340)
	if err := tt.Click("crest-5.png"); err != nil {
		t.Fatal(err)
	}
	p.entries = p.entries[:1]
	tt.Frame()
	if err := tt.ClickWith(ui.Shift, "sketches"); err != nil {
		t.Fatal(err)
	}
	if got := p.s.Selected(); len(got) != 1 || got[0] != "/art/sketches" {
		t.Fatalf("selection %v", got)
	}
}
