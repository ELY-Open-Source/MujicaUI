package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestInlineDiffNamesAndShowsBothLines(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		InlineDiff(c, "The court meets at noon", "The court meets at dusk", InlineDiffOptions{})
	}, 600, 120)
	if _, ok := tt.Find("Old: The court meets at noon"); !ok {
		t.Fatalf("old line: %q", tt.Texts())
	}
	if _, ok := tt.Find("New: The court meets at dusk"); !ok {
		t.Fatalf("new line: %q", tt.Texts())
	}
}

func TestInlineDiffCachesWords(t *testing.T) {
	from, to := "a b", "a c"
	var cache *inlineDiffCache
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		e := InlineDiff(c, from, to, InlineDiffOptions{Mono: true})
		cache = *ui.Local(e, "words", func() *inlineDiffCache { return nil })
	}, 300, 100)
	first := &cache.a[0]
	tt.Frame()
	if &cache.a[0] != first {
		t.Fatal("word diff recomputed without a change")
	}
	to = "a d"
	tt.Frame()
	if cache.b[1].Text != "d" {
		t.Fatalf("new text not diffed: %+v", cache.b)
	}
}

func TestInlineDiffSnapshot(t *testing.T) {
	snapshot(t, "inline_diff", 560, 200, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			InlineDiff(c, "The court meets at noon in the east hall.", "The council meets at dusk in the east hall.", InlineDiffOptions{})
			InlineDiff(c, "return a + b", "return a - total", InlineDiffOptions{Mono: true})
		})
	})
}
