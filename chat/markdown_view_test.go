package chat

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

const markdownViewSample = "# Decree\n\nBy order of **the Queen**, *all* guests use `the east door`.\n\n" +
	"- Arrive at nine\n- Wear [court dress](https://example.com/dress)\n\n> Silence in the hall.\n\n" +
	"```go\nfunc bow() {}\n```\n\n| Room | Floor |\n|---|--:|\n| Hall | 1 |\n\n---\n\nEnd."

func TestMarkdownViewRendersSubset(t *testing.T) {
	src := markdownViewSample
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Children(func() { MarkdownView(c, src) })
	}, 600, 900)
	for _, s := range []string{"Decree", "the Queen", "all", "the east door", "Arrive at nine", "court dress",
		"Silence in the hall.", "func", "bow", "Room", "Hall", "End."} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
}

func TestMarkdownViewLinkByClickAndKeyboard(t *testing.T) {
	var got []string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if url, ok := MarkdownView(c, "See [the guide](https://x.io/guide).").Link(); ok {
			got = append(got, url)
		}
	}, 400, 200)
	if err := tt.Click("the guide"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnter)
	if len(got) != 2 || got[0] != "https://x.io/guide" || got[1] != got[0] {
		t.Fatalf("links %q", got)
	}
}

func TestMarkdownViewStreamsWithoutReparsingFinishedBlocks(t *testing.T) {
	src := ""
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Children(func() { MarkdownView(c, src) })
	}, 600, 900)
	for _, chunk := range strings.SplitAfter(markdownViewSample, " ") {
		src += chunk
		tt.Frame()
	}
	if !tt.HasText("End.") || !tt.HasText("Decree") {
		t.Fatalf("streamed render incomplete: %q", tt.Texts())
	}
}

func TestMarkdownViewSnapshot(t *testing.T) {
	snapshot(t, "markdown_view", 560, 720, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { MarkdownView(c, markdownViewSample) })
	})
}
