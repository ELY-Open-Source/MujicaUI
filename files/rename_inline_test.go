package files

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

type renameProbe struct {
	name              string
	editing, disabled bool
	folder            bool
	renamed, canceled int
}

func (p *renameProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Padding(12).Gap(8).Children(func() {
		r := RenameInline(c, &p.name, &p.editing, RenameInlineOptions{Folder: p.folder, Disabled: p.disabled,
			Validate: func(n string) string {
				if strings.ContainsRune(n, '/') {
					return "No slashes"
				}
				return ""
			}})
		if r.Renamed {
			p.renamed++
		}
		if r.Canceled {
			p.canceled++
		}
		input.Button(c, "Elsewhere", input.ButtonOptions{})
	})
}

func TestRenameStem(t *testing.T) {
	for in, want := range map[string]string{"seal.txt": "seal", "a.tar.gz": "a.tar", ".env": ".env", "README": "README"} {
		if got := renameStem(in, false); got != want {
			t.Errorf("%s: stem %q, want %q", in, got, want)
		}
	}
	if renameStem("v1.2", true) != "v1.2" {
		t.Fatal("folder stem dropped its dot part")
	}
}

func TestRenameInlineTypingReplacesStem(t *testing.T) {
	p := &renameProbe{name: "seal.txt", editing: true}
	tt := ui.NewTester(p.view, 400, 200)
	if !tt.Focused("Name") {
		t.Fatal("editor did not take the focus")
	}
	tt.Type("wax")
	tt.Type("y")
	tt.Key(0, ui.KeyEnter)
	if p.name != "waxy.txt" || p.renamed != 1 || p.editing {
		t.Fatalf("name %q renamed %d editing %v", p.name, p.renamed, p.editing)
	}
}

func TestRenameInlineKeysAroundSelection(t *testing.T) {
	p := &renameProbe{name: "seal.txt", editing: true}
	tt := ui.NewTester(p.view, 400, 200)
	tt.Key(0, ui.KeyLeft)
	tt.Type("royal ")
	tt.Key(0, ui.KeyEnter)
	if p.name != "royal seal.txt" {
		t.Fatalf("Left then typing gave %q", p.name)
	}
	p.editing = true
	tt.Frame()
	tt.Key(0, ui.KeyBackspace)
	tt.Type("crest")
	tt.Key(0, ui.KeyEnter)
	if p.name != "crest.txt" {
		t.Fatalf("Backspace then typing gave %q", p.name)
	}
	p.editing = true
	tt.Frame()
	tt.Key(0, ui.KeyRight)
	tt.Type("s")
	tt.Key(0, ui.KeyEnter)
	if p.name != "crests.txt" {
		t.Fatalf("Right then typing gave %q", p.name)
	}
}

func TestRenameInlineValidationAndCancel(t *testing.T) {
	p := &renameProbe{name: "seal.txt", editing: true}
	tt := ui.NewTester(p.view, 400, 200)
	tt.Type("a/b")
	tt.Key(0, ui.KeyEnter)
	if !p.editing || p.name != "seal.txt" || !tt.HasText("No slashes") {
		t.Fatalf("invalid name accepted: %q %q", p.name, tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	if p.editing || p.name != "seal.txt" || p.canceled != 1 {
		t.Fatalf("Escape: name %q editing %v canceled %d", p.name, p.editing, p.canceled)
	}
}

func TestRenameInlineBlurCommitsAndFolder(t *testing.T) {
	p := &renameProbe{name: "v1.2", editing: true, folder: true}
	tt := ui.NewTester(p.view, 400, 200)
	tt.Type("archive")
	if err := tt.Click("Elsewhere"); err != nil {
		t.Fatal(err)
	}
	if p.name != "archive" || p.editing {
		t.Fatalf("blur gave %q editing %v", p.name, p.editing)
	}
}

func TestRenameInlineDoubleClickAndDisabled(t *testing.T) {
	p := &renameProbe{name: "seal.txt"}
	tt := ui.NewTester(p.view, 400, 200)
	r, _ := tt.Find("seal.txt")
	tt.ClickAt(r.X+4, r.Y+4)
	tt.ClickAt(r.X+4, r.Y+4)
	if !p.editing {
		t.Fatal("double click did not start editing")
	}
	p.disabled = true
	tt.Frame()
	if p.editing || p.canceled != 1 {
		t.Fatal("disabling did not end the edit")
	}
}

func TestRenameInlineSnapshot(t *testing.T) {
	snapshot(t, "rename_inline", 360, 110, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *renameProbe { return &renameProbe{name: "royal seal.txt", editing: true} })
		ui.Column(c).Padding(12).Gap(8).Children(func() {
			RenameInline(c, &(*p).name, &(*p).editing, RenameInlineOptions{})
			name, off := "ledger.md", false
			RenameInline(c, &name, &off, RenameInlineOptions{})
		})
	})
}

func TestRenameInlineCompositionCommit(t *testing.T) {
	p := &renameProbe{name: "seal.txt", editing: true}
	tt := ui.NewTester(p.view, 400, 200)
	tt.Compose("印", 1)
	tt.Type("印章")
	tt.Key(0, ui.KeyEnter)
	if p.name != "印章.txt" || p.renamed != 1 {
		t.Fatalf("committed name %q", p.name)
	}
}

func TestRenameInlineNativeTypingAndPaste(t *testing.T) {
	p := &renameProbe{name: "seal.txt", editing: true}
	tt := ui.NewTester(p.view, 400, 200)
	tt.TypeKey(0, ui.KeyW, "w")
	tt.Key(0, ui.KeyEnter)
	if p.name != "w.txt" {
		t.Fatalf("native typing renamed %q", p.name)
	}
	p.editing = true
	tt.Frame()
	tt.SetClipboard("wax")
	tt.Key(ui.Cmd, ui.KeyV)
	tt.Key(0, ui.KeyEnter)
	if p.name != "wax.txt" {
		t.Fatalf("paste renamed %q", p.name)
	}
}
