package code

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
)

var codeViewerSample = []string{
	"package court",
	"",
	"// Seal stamps a decree.",
	"func Seal(d string) int {",
	"\tif d == \"\" {",
	"\t\treturn 0",
	"\t}",
	"\treturn 42",
	"}",
}

func codeViewerTester(lines []string, s *CodeViewerState, disabled bool) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CodeViewer(c, lines, s, CodeViewerOptions{Language: "go", Label: "Source", Disabled: disabled}).Element.Height(400)
	}, 520, 420)
}

func TestCodeViewerFoldRanges(t *testing.T) {
	ends := codeViewerFolds(codeViewerSample, nil)
	if ends[3] != 7 || ends[4] != 5 || ends[0] != 0 || ends[8] != 8 {
		t.Fatalf("fold ends %v", ends)
	}
}

func TestCodeViewerFoldsAndNumbers(t *testing.T) {
	s := &CodeViewerState{}
	tt := codeViewerTester(codeViewerSample, s, false)
	if !tt.HasText("9") || !tt.HasText("\treturn 42") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if err := tt.Click("Fold line 4"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("\treturn 42") || !tt.HasText("}") || len(s.shown) != 5 {
		t.Fatalf("fold kept %d lines: %q", len(s.shown), tt.Texts())
	}
	if err := tt.Click("Unfold line 4"); err != nil {
		t.Fatal(err)
	}
	if len(s.shown) != 9 {
		t.Fatalf("unfold shows %d lines", len(s.shown))
	}
}

func TestCodeViewerKeysCurrentLineAndCopy(t *testing.T) {
	s := &CodeViewerState{}
	tt := codeViewerTester(codeViewerSample, s, false)
	for !tt.Focused("Source") {
		tt.Key(0, ui.KeyTab)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(ui.Shift, ui.KeyDown)
	tt.Key(ui.Shift, ui.KeyDown)
	tt.Key(ui.Cmd, ui.KeyC)
	if want := "package court\n\n// Seal stamps a decree."; tt.Clipboard() != want {
		t.Fatalf("copied %q", tt.Clipboard())
	}
	tt.Key(0, ui.KeyDown)
	if s.Line() != 3 {
		t.Fatalf("current line %d, want 3", s.Line())
	}
	tt.Key(0, ui.KeyLeft)
	if !s.Folded.Has(3) || len(s.shown) != 5 {
		t.Fatal("Left did not fold the block")
	}
	tt.Key(0, ui.KeyRight)
	if s.Folded.Has(3) {
		t.Fatal("Right did not unfold the block")
	}
}

func TestCodeViewerGoToOpensFold(t *testing.T) {
	s := &CodeViewerState{}
	tt := codeViewerTester(codeViewerSample, s, false)
	_ = tt.Click("Fold line 4")
	s.GoTo(7)
	tt.Frame()
	if s.Line() != 7 || s.Folded.Has(3) {
		t.Fatalf("GoTo gave line %d folded %v", s.Line(), s.Folded.Has(3))
	}
}

func TestCodeViewerHighlightColors(t *testing.T) {
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		k := core.Tokens(c)
		kw, str := CodeTokenColor(k, internal.TokKeyword), CodeTokenColor(k, internal.TokString)
		if kw != k.AccentText || str != k.Success {
			t.Errorf("keyword %v string %v", kw, str)
		}
	}, 10, 10)
}

func TestCodeViewerHundredThousandLines(t *testing.T) {
	lines := make([]string, 100_000)
	for i := range lines {
		lines[i] = fmt.Sprintf("x%d := %d // line", i, i)
	}
	s := &CodeViewerState{}
	tt := codeViewerTester(lines, s, false)
	rows := 0
	for _, t := range tt.Texts() {
		if strings.HasPrefix(t, "x") {
			rows++
		}
	}
	if rows == 0 || rows > 40 {
		t.Fatalf("built %d rows of 100000", rows)
	}
	s.GoTo(99_998)
	tt.Frame()
	if !tt.HasText("x99998 := 99998 // line") {
		t.Fatal("GoTo did not scroll to the end")
	}
}

func TestCodeViewerDisabledFold(t *testing.T) {
	s := &CodeViewerState{}
	tt := codeViewerTester(codeViewerSample, s, true)
	_ = tt.Click("Fold line 4")
	if s.Folded.Len() != 0 {
		t.Fatal("disabled viewer folded")
	}
}

func TestCodeViewerSnapshot(t *testing.T) {
	snapshot(t, "code_viewer", 480, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "snap", func() *CodeViewerState {
			s := &CodeViewerState{}
			s.GoTo(5)
			return s
		})
		ui.Column(c).Padding(12).Children(func() {
			CodeViewer(c, codeViewerSample, *s, CodeViewerOptions{Language: "go"}).Element.Height(276)
		})
	})
}
