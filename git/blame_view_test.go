package git

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func blameSample() []BlameLine {
	a := BlameLine{Hash: gitTestHash(1), Author: "Sakiko", When: gitTestTime, Summary: "Add the graph"}
	b := BlameLine{Hash: gitTestHash(2), Author: "Uika", When: gitTestTime, Summary: "Comment the lanes"}
	lines := []string{"package mujica", "", "/* lanes", "   per row */", "func lanes() int {", "\treturn 2", "}"}
	out := make([]BlameLine, len(lines))
	for i, s := range lines {
		ln := a
		if i >= 2 && i <= 3 {
			ln = b
		}
		ln.Text = s
		out[i] = ln
	}
	return out
}

func TestBlameViewGroupsCommits(t *testing.T) {
	var s BlameViewState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		BlameView(c, &s, blameSample(), BlameViewOptions{Language: "go"}).Height(300)
	}, 700, 340)
	count := 0
	for _, txt := range tt.Texts() {
		if txt == "Sakiko" || txt == "Uika" {
			count++
		}
	}
	// Three runs: lines 1-2, 3-4 and 5-7.
	if count != 3 {
		t.Fatalf("author shown %d times, want 3: %q", count, tt.Texts())
	}
	if _, ok := tt.Find(fmt.Sprintf("Line 4, Uika, %s", gitShortHash(gitTestHash(2)))); !ok {
		t.Fatal("line 4 gutter not named")
	}
	if !s.blocks[3] || s.blocks[4] {
		t.Fatalf("block comment state %v", s.blocks)
	}
}

func TestBlameViewHoverDetails(t *testing.T) {
	var s BlameViewState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		BlameView(c, &s, blameSample(), BlameViewOptions{}).Height(300)
	}, 700, 340)
	r, _ := tt.Find("Uika")
	tt.Move(r.X+2, r.Y+2)
	tt.Advance(650 * time.Millisecond)
	tt.Move(r.X+3, r.Y+2)
	tt.Frame()
	if tt.HasText(gitTestHash(2) + "\nUika · 2026-10-05\nComment the lanes") {
		return
	}
	t.Fatalf("no hover details: %q", tt.Texts())
}

func TestBlameViewVirtualized(t *testing.T) {
	lines := make([]BlameLine, 100000)
	for i := range lines {
		lines[i] = BlameLine{Hash: gitTestHash(i / 3), Author: "Sakiko", When: gitTestTime, Text: fmt.Sprintf("line %d", i)}
	}
	var s BlameViewState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		BlameView(c, &s, lines, BlameViewOptions{}).Height(400)
	}, 700, 440)
	if n := len(tt.Texts()); n > 400 {
		t.Fatalf("%d texts for 100,000 lines", n)
	}
	blocks := &s.blocks[0]
	tt.Frame()
	if &s.blocks[0] != blocks {
		t.Fatal("highlight state recomputed without a data change")
	}
}

func TestBlameViewSnapshot(t *testing.T) {
	lines := blameSample()
	snapshot(t, "blame_view", 640, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "blame", func() *BlameViewState { return &BlameViewState{} })
		ui.Column(c).Padding(16).Children(func() {
			BlameView(c, *s, lines, BlameViewOptions{Language: "go"}).Height(260)
		})
	})
}
