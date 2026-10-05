package internal

import (
	"strings"
	"testing"
)

func TestHighlightGo(t *testing.T) {
	toks, open := Highlight("go", `func f(s string) int { return 42 // done`, false)
	if open {
		t.Fatal("block left open")
	}
	got := map[string]TokenKind{}
	for _, tk := range toks {
		got[strings.TrimSpace(tk.Text)] = tk.Kind
	}
	for text, want := range map[string]TokenKind{"func": TokKeyword, "string": TokType, "42": TokNumber, "// done": TokComment} {
		if got[text] != want {
			t.Errorf("%q kind %d, want %d (tokens %+v)", text, got[text], want, toks)
		}
	}
}

func TestHighlightStringsAndEscapes(t *testing.T) {
	toks, _ := Highlight("javascript", `const s = "a\"b"; x`, false)
	var str string
	for _, tk := range toks {
		if tk.Kind == TokString {
			str = tk.Text
		}
	}
	if str != `"a\"b"` {
		t.Fatalf("string token %q", str)
	}
}

func TestHighlightBlockCommentAcrossLines(t *testing.T) {
	_, open := Highlight("go", "x /* start", false)
	if !open {
		t.Fatal("block not open")
	}
	toks, open := Highlight("go", "still */ y", true)
	if open || toks[0].Kind != TokComment || toks[0].Text != "still */" {
		t.Fatalf("tokens %+v open %v", toks, open)
	}
}

func TestHighlightUnknownLanguageIsPlain(t *testing.T) {
	toks, _ := Highlight("cobol", "MOVE A TO B", false)
	if len(toks) != 1 || toks[0].Kind != TokPlain {
		t.Fatalf("tokens %+v", toks)
	}
}

func TestHighlightKeepsText(t *testing.T) {
	line := `select * from "users" where id = 7 -- note`
	toks, _ := Highlight("sql", line, false)
	var b strings.Builder
	for _, tk := range toks {
		b.WriteString(tk.Text)
	}
	if b.String() != line {
		t.Fatalf("round trip %q", b.String())
	}
}
