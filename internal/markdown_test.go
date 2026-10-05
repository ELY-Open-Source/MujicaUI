package internal

import (
	"reflect"
	"strings"
	"testing"
)

const mdSample = "# Title\n\nSome **bold** and *em* text with `code`\nand a [link](https://x.io).\n\n" +
	"- one\n- two\n  - nested\n3. three\n\n> quoted *line*\n\n```go\nfunc f() {}\n```\n\n" +
	"| a | b |\n|:--|--:|\n| 1 | 2 |\n| 3 |\n\n---\n\nend ~~old~~ new"

func TestMarkdownBlocks(t *testing.T) {
	var m Markdown
	m.Update(mdSample)
	var kinds []BlockKind
	for _, b := range m.Blocks {
		kinds = append(kinds, b.Kind)
	}
	want := []BlockKind{BlockHeading, BlockParagraph, BlockList, BlockQuote, BlockCode, BlockTable, BlockRule, BlockParagraph}
	if !reflect.DeepEqual(kinds, want) {
		t.Fatalf("kinds %v, want %v", kinds, want)
	}
	b := m.Blocks
	if b[0].Level != 1 || b[0].Spans[0].Text != "Title" {
		t.Errorf("heading %+v", b[0])
	}
	p := b[1].Spans
	if p[1] != (MdSpan{Text: "bold", Strong: true}) || p[3] != (MdSpan{Text: "em", Em: true}) || p[5] != (MdSpan{Text: "code", Code: true}) {
		t.Errorf("paragraph spans %+v", p)
	}
	if p[7] != (MdSpan{Text: "link", URL: "https://x.io"}) {
		t.Errorf("link %+v", p[7])
	}
	it := b[2].Items
	if len(it) != 4 || it[2].Depth != 1 || it[2].Spans[0].Text != "nested" || it[3].Number != 3 {
		t.Errorf("list items %+v", it)
	}
	if b[3].Spans[1] != (MdSpan{Text: "line", Em: true}) {
		t.Errorf("quote %+v", b[3].Spans)
	}
	if b[4].Lang != "go" || b[4].Code != "func f() {}" || b[4].Open {
		t.Errorf("code %+v", b[4])
	}
	tb := b[5]
	if !reflect.DeepEqual(tb.Header, []string{"a", "b"}) || !reflect.DeepEqual(tb.Align, []TableAlign{AlignLeft, AlignRight}) ||
		!reflect.DeepEqual(tb.Rows, [][]string{{"1", "2"}, {"3", ""}}) {
		t.Errorf("table %+v", tb)
	}
	if b[7].Spans[1] != (MdSpan{Text: "old", Strike: true}) {
		t.Errorf("strike %+v", b[7].Spans)
	}
}

func TestMarkdownInlineLiterals(t *testing.T) {
	for in, want := range map[string]string{
		"2 * 3 * 4":       "2 * 3 * 4",
		"snake_case_name": "snake_case_name",
		"**unclosed":      "**unclosed",
		`\*not em\*`:      "*not em*",
		"[no link]":       "[no link]",
	} {
		sp := ParseInline(in)
		if len(sp) != 1 || sp[0].Text != want || sp[0].Em || sp[0].Strong {
			t.Errorf("%q gave %+v", in, sp)
		}
	}
}

// TestMarkdownIncremental streams the sample byte by byte and checks finished blocks are kept.
func TestMarkdownIncremental(t *testing.T) {
	var m Markdown
	for i := 1; i <= len(mdSample); i++ {
		before := append([]MdBlock(nil), m.Blocks...)
		from := m.Update(mdSample[:i])
		if from < len(before)-1 {
			t.Fatalf("at %d reparsed from %d of %d blocks", i, from, len(before))
		}
		for j := 0; j < from; j++ {
			if !reflect.DeepEqual(m.Blocks[j], before[j]) {
				t.Fatalf("at %d finished block %d changed", i, j)
			}
		}
	}
	var whole Markdown
	whole.Update(mdSample)
	if !reflect.DeepEqual(m.Blocks, whole.Blocks) {
		t.Fatal("streamed parse differs from a whole parse")
	}
	if m.Update(mdSample) != len(m.Blocks) {
		t.Fatal("an unchanged source parsed again")
	}
	if m.Update("# Other") != 0 || len(m.Blocks) != 1 {
		t.Fatal("a replaced source kept old blocks")
	}
}

func TestMarkdownOpenFence(t *testing.T) {
	var m Markdown
	m.Update("```py\nprint(1)\n")
	if len(m.Blocks) != 1 || !m.Blocks[0].Open || m.Blocks[0].Code != "print(1)" {
		t.Fatalf("open fence %+v", m.Blocks)
	}
	m.Update("```py\nprint(1)\n```\nafter")
	if len(m.Blocks) != 2 || m.Blocks[0].Open || !strings.Contains(m.Blocks[1].Spans[0].Text, "after") {
		t.Fatalf("closed fence %+v", m.Blocks)
	}
}
