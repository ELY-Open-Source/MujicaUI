package messaging

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestQuotedTextAndStrip(t *testing.T) {
	if got := QuoteStrip("> a\n>b\n> > c"); got != "a\nb\n> c" {
		t.Fatalf("strip %q", got)
	}
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		QuotedText(c, "> first\n> second", &open)
	}, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !open || !tt.HasText("first\nsecond") || !tt.HasText("Hide quoted text") {
		t.Fatalf("not expanded: %v", tt.Texts())
	}
	tt.Key(0, ui.KeySpace)
	if open {
		t.Fatal("Space did not fold")
	}
}

func TestQuotedTextSnapshot(t *testing.T) {
	a, b := false, true
	snapshot(t, "quoted_text", 420, 180, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(8).Children(func() {
			QuotedText(c, "> Can you send the score?\n> Thanks", &a)
			QuotedText(c, "> Rehearsal moves to nine.\n> Bring the masks.", &b)
		})
	})
}

func TestQuotedTextSingleLineCopy(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		QuotedText(c, "> One line", &open)
	}, 400, 200)
	if !tt.HasText("Show quoted text (1 line)") {
		t.Fatal("singular quote label missing")
	}
}
