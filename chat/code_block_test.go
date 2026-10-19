package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

const codeBlockSample = "// greet the court\nfunc greet(name string) int {\n\treturn len(\"hi \" + name) + 42\n}"

func codeBlockTester(opts CodeBlockOptions, copies *int, disabled *bool) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(360).Disabled(*disabled).Children(func() {
			if CodeBlock(c, codeBlockSample, opts).Copied() {
				*copies++
			}
		})
	}, 400, 300)
}

func TestCodeBlockLabelNumbersAndCopy(t *testing.T) {
	copies, off := 0, false
	tt := codeBlockTester(CodeBlockOptions{Language: "go", LineNumbers: true}, &copies, &off)
	for _, s := range []string{"go", "1", "4", "func", "greet", "42"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	if err := tt.Click("Copy code"); err != nil {
		t.Fatal(err)
	}
	if tt.Clipboard() != codeBlockSample || copies != 1 {
		t.Fatalf("copied %d, clipboard matches %v", copies, tt.Clipboard() == codeBlockSample)
	}
	if _, ok := tt.Find("Copied"); !ok {
		t.Fatal("no copied confirmation")
	}
}

func TestCodeBlockCollapseAndWrapByKeyboard(t *testing.T) {
	copies, off := 0, false
	tt := codeBlockTester(CodeBlockOptions{Language: "go"}, &copies, &off)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !tt.Focused("Wrap lines") {
		t.Fatal("Tab did not reach the wrap toggle")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if tt.HasText("greet") {
		t.Fatal("collapse left the code built")
	}
	if err := tt.Click("Show 4 lines"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("greet") {
		t.Fatal("expand did not show the code")
	}
}

func TestCodeBlockDisabledBlocksActions(t *testing.T) {
	copies, off := 0, true
	tt := codeBlockTester(CodeBlockOptions{Collapsed: true}, &copies, &off)
	_ = tt.Click("Copy code")
	_ = tt.Click("Show 4 lines")
	if copies != 0 || tt.Clipboard() != "" || tt.HasText("greet") {
		t.Fatalf("disabled block acted: copies %d", copies)
	}
}

func TestCodeBlockSnapshot(t *testing.T) {
	snapshot(t, "code_block", 440, 200, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			CodeBlock(c, codeBlockSample, CodeBlockOptions{Language: "go", LineNumbers: true})
		})
	})
}
