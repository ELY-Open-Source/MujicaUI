package files

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// fileFS maps folders to entry names, with trailing directory slashes.
var fileFS = map[string][]string{
	"/court":             {"chronicles/", "decrees/", "seal.go", "notes.md"},
	"/court/chronicles":  {"first.md", "second.md"},
	"/court/decrees":     {"old/", "42.pdf"},
	"/court/decrees/old": {"7.pdf"},
}

func fileFSOptions() FileExplorerOptions {
	return FileExplorerOptions{
		Roots: []string{"/court"},
		Folders: func(dir string) []string {
			var out []string
			for _, n := range fileFS[dir] {
				if strings.HasSuffix(n, "/") {
					out = append(out, path.Join(dir, strings.TrimSuffix(n, "/")))
				}
			}
			return out
		},
		Entries: func(dir string) []FileEntry {
			var out []FileEntry
			for _, n := range fileFS[dir] {
				name := strings.TrimSuffix(n, "/")
				out = append(out, FileEntry{Path: path.Join(dir, name), Name: name, Dir: name != n, Size: 2048})
			}
			return out
		},
	}
}

type explorerProbe struct {
	s      FileExplorerState
	opts   FileExplorerOptions
	nav    int
	opened []string
}

func (p *explorerProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := FileExplorer(c, &p.s, p.opts)
	r.Element.Height(380)
	if r.Navigated {
		p.nav++
	}
	if r.Opened != nil {
		p.opened = append(p.opened, r.Opened.Path)
	}
}

func TestFileExplorerOpenFolderAndHistory(t *testing.T) {
	p := &explorerProbe{opts: fileFSOptions()}
	tt := ui.NewTester(p.view, 720, 400)
	if p.s.Current != "/court" || !tt.HasText("seal.go") {
		t.Fatalf("start %q", p.s.Current)
	}
	rows := 0
	for _, s := range tt.Texts() {
		if s == "decrees" {
			rows++
		}
	}
	if rows != 2 {
		t.Fatalf("decrees shows %d times, want tree and list", rows)
	}
	// The list's decrees row is the one above seal.go.
	seal, _ := tt.Find("seal.go")
	tt.ClickAt(seal.X, seal.Y+seal.H/2-36)
	tt.ClickAt(seal.X, seal.Y+seal.H/2-36)
	if p.s.Current != "/court/decrees" || !tt.HasText("42.pdf") || p.nav != 1 {
		t.Fatalf("open folder: current %q nav %d", p.s.Current, p.nav)
	}
	if err := tt.Click("Back"); err != nil {
		t.Fatal(err)
	}
	if p.s.Current != "/court" {
		t.Fatalf("back to %q", p.s.Current)
	}
	_ = tt.Click("Forward")
	if p.s.Current != "/court/decrees" {
		t.Fatalf("forward to %q", p.s.Current)
	}
	_ = tt.Click("Enclosing folder")
	if p.s.Current != "/court" {
		t.Fatalf("up to %q", p.s.Current)
	}
}

func TestFileExplorerTreeAndOpenFile(t *testing.T) {
	p := &explorerProbe{opts: fileFSOptions()}
	tt := ui.NewTester(p.view, 720, 400)
	p.s.Navigate("/court/decrees/old")
	tt.Frame()
	if !tt.HasText("old") || !tt.HasText("7.pdf") {
		t.Fatalf("tree did not open the ancestors: %q", tt.Texts())
	}
	if sel, _ := p.s.tree.Selected(); sel != "/court/decrees/old" {
		t.Fatalf("tree chose %q", sel)
	}
	r, _ := tt.Find("7.pdf")
	tt.ClickAt(r.X, r.Y+r.H/2)
	tt.ClickAt(r.X, r.Y+r.H/2)
	if len(p.opened) != 1 || p.opened[0] != "/court/decrees/old/7.pdf" {
		t.Fatalf("opened %v", p.opened)
	}
	if got := p.s.Selected(); len(got) != 1 {
		t.Fatalf("selection %v", got)
	}
}

func TestFileExplorerShortcutsAndDisabled(t *testing.T) {
	p := &explorerProbe{opts: fileFSOptions()}
	tt := ui.NewTester(p.view, 720, 400)
	p.s.Navigate("/court/chronicles")
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	back, up := ui.Modifiers(ui.Alt), ui.Modifiers(ui.Alt)
	bk, uk := ui.KeyLeft, ui.KeyUp
	if runtime.GOOS == "darwin" {
		back, up, bk = ui.Super, ui.Super, ui.KeyBracketLeft
	}
	tt.Key(back, bk)
	if p.s.Current != "/court" {
		t.Fatalf("back shortcut gave %q", p.s.Current)
	}
	p.s.Navigate("/court/decrees/old")
	tt.Frame()
	tt.Key(up, uk)
	if p.s.Current != "/court/decrees" {
		t.Fatalf("up shortcut gave %q", p.s.Current)
	}
	p.opts.Disabled = true
	tt.Frame()
	_ = tt.Click("Back")
	if p.s.Current != "/court/decrees" {
		t.Fatal("disabled explorer moved")
	}
}

func TestFileExplorerSnapshot(t *testing.T) {
	snapshot(t, "file_explorer", 640, 300, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *explorerProbe {
			p := &explorerProbe{opts: fileFSOptions()}
			p.s.Navigate("/court/decrees")
			return p
		})
		ui.Column(c).Padding(12).Children(func() {
			FileExplorer(c, &(*p).s, (*p).opts).Element.Height(276)
		})
	})
}

func TestFileExplorerDisablingBeforeDeliveredKey(t *testing.T) {
	p := &explorerProbe{opts: fileFSOptions()}
	tt := ui.NewTester(p.view, 720, 400)
	if err := tt.Click("seal.go"); err != nil {
		t.Fatal(err)
	}
	before := p.s.Selected()
	p.opts.Disabled = true
	tt.Key(0, ui.KeyDown)
	after := p.s.Selected()
	if len(after) != len(before) || after[0] != before[0] {
		t.Fatalf("disabled selection %v, want %v", after, before)
	}
}

func TestFileExplorerPrunesRemovedSelection(t *testing.T) {
	entries := []FileEntry{{Path: "/court/a.txt", Name: "a.txt"}, {Path: "/court/b.txt", Name: "b.txt"}}
	var s FileExplorerState
	opts := FileExplorerOptions{Roots: []string{"/court"}, Folders: func(string) []string { return nil }, Entries: func(string) []FileEntry { return entries }}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		FileExplorer(c, &s, opts).Element.Height(380)
	}, 720, 400)
	if err := tt.Click("b.txt"); err != nil {
		t.Fatal(err)
	}
	entries = entries[:1]
	tt.Frame()
	if got := s.Selected(); len(got) != 0 {
		t.Fatalf("removed selection %v", got)
	}
}

func TestFileExplorerLargeAndRevision(t *testing.T) {
	entries := make([]FileEntry, 100000)
	for i := range entries {
		name := fmt.Sprintf("f-%d.txt", i)
		entries[i] = FileEntry{Path: "/court/" + name, Name: name}
	}
	var s FileExplorerState
	opts := FileExplorerOptions{Roots: []string{"/court"}, Folders: func(string) []string { return nil }, Entries: func(string) []FileEntry { return entries }}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		FileExplorer(c, &s, opts).Element.Height(380)
	}, 720, 400)
	n := 0
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "f-") {
			n++
		}
	}
	if n == 0 || n > 60 {
		t.Fatalf("built %d file labels", n)
	}
	if err := tt.Click("f-0.txt"); err != nil {
		t.Fatal(err)
	}
	entries[0].Path = "/court/replaced.txt"
	opts.Revision++
	tt.Frame()
	if got := s.Selected(); len(got) != 0 {
		t.Fatalf("replaced selection %v", got)
	}
}
