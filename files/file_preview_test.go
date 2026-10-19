package files

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
)

type previewProbe struct {
	f                       PreviewFile
	opts                    FilePreviewOptions
	closed, opened, retried int
}

func (p *previewProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := FilePreview(c, p.f, p.opts)
	r.Element.Size(480, 300)
	for _, e := range []struct {
		on bool
		n  *int
	}{{r.Closed, &p.closed}, {r.Opened, &p.opened}, {r.Retried, &p.retried}} {
		if e.on {
			*e.n++
		}
	}
}

func TestFilePreviewKinds(t *testing.T) {
	cases := []struct {
		f    PreviewFile
		want string
	}{
		{PreviewFile{Name: "seal.go", Text: "package court\nfunc Seal() {}"}, "func Seal() {}"},
		{PreviewFile{Name: "notes.md", Text: "# The court"}, "# The court"},
		{PreviewFile{Name: "deed.pdf", Size: 4096}, "No preview for PDF files"},
		{PreviewFile{Name: "crest.svg", Image: ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 4 4"><rect width="4" height="4"/></svg>`))}, "crest.svg"},
	}
	for _, tc := range cases {
		tt := ui.NewTester((&previewProbe{f: tc.f}).view, 500, 320)
		if !tt.HasText(tc.want) {
			t.Fatalf("%s: texts %q", tc.f.Name, tt.Texts())
		}
	}
	tt := ui.NewTester((&previewProbe{f: PreviewFile{Name: "deed.pdf", Size: 4096}}).view, 500, 320)
	if !tt.HasText("4 KB") {
		t.Fatal("size missing")
	}
}

func TestFilePreviewActions(t *testing.T) {
	p := &previewProbe{f: PreviewFile{Name: "notes.md", Text: "text"}}
	tt := ui.NewTester(p.view, 500, 320)
	_ = tt.Click("Open")
	_ = tt.Click("Close preview")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEscape)
	if p.opened != 1 || p.closed != 2 {
		t.Fatalf("opened %d closed %d", p.opened, p.closed)
	}
}

func TestFilePreviewStatusAndNoName(t *testing.T) {
	p := &previewProbe{f: PreviewFile{Name: "big.log"}, opts: FilePreviewOptions{Status: data.DataFailed, Error: "read denied"}}
	tt := ui.NewTester(p.view, 500, 320)
	if !tt.HasText("read denied") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Retry")
	if p.retried != 1 {
		t.Fatal("retry not reported")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("missing name did not panic")
		}
	}()
	ui.NewTester((&previewProbe{}).view, 500, 320)
}

func TestFilePreviewSnapshot(t *testing.T) {
	snapshot(t, "file_preview", 500, 340, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			FilePreview(c, PreviewFile{Name: "seal.go", Size: 912, Text: "package court\n\n// Seal stamps a decree.\nfunc Seal(d string) int {\n\treturn len(d)\n}"}, FilePreviewOptions{}).Element.Height(316)
		})
	})
}

func TestFilePreviewEmptyCode(t *testing.T) {
	for _, f := range []PreviewFile{{Name: "empty.go"}, {Name: "empty.html"}, {Name: "empty.custom", Language: "go"}} {
		tt := ui.NewTester((&previewProbe{f: f}).view, 500, 320)
		if tt.HasText("No preview for GO files") || tt.HasText("No preview for HTML files") || tt.HasText("No preview for CUSTOM files") {
			t.Fatalf("empty code shows unsupported type: %s", f.Name)
		}
	}
}

func TestFilePreviewLargeCodeCache(t *testing.T) {
	file := PreviewFile{Name: "large.go", Text: strings.Repeat("var seal = 1\n", 100000)}
	var state *previewCodeState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := FilePreview(c, file, FilePreviewOptions{})
		r.Element.Size(480, 300)
		state = ui.Local(r.Element, "file_preview.code", func() previewCodeState { return previewCodeState{} })
	}, 500, 320)
	before := state.viewer.Data
	tt.Frame()
	if before == nil || state.viewer.Data != before || len(state.lines) != 100001 {
		t.Fatal("unchanged preview rebuilt code cache")
	}
	n := 0
	for _, text := range tt.Texts() {
		if text == "var seal = 1" {
			n++
		}
	}
	if n == 0 || n > 60 {
		t.Fatalf("built %d code lines", n)
	}
	file.Text = "package court"
	tt.Frame()
	if len(state.lines) != 1 || state.lines[0] != "package court" || state.viewer.Data == before {
		t.Fatal("changed preview retained code cache")
	}
}
