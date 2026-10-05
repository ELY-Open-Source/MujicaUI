package internal

import (
	"slices"
	"testing"
)

func TestWordTokens(t *testing.T) {
	got := wordTokens("let x_1 = foo(bar);  ok")
	want := []string{"let", " ", "x_1", " ", "=", " ", "foo", "(", "bar", ")", ";", "  ", "ok"}
	if !slices.Equal(got, want) {
		t.Fatalf("tokens %q", got)
	}
	if wordTokens("") != nil {
		t.Fatal("empty line has tokens")
	}
}

func TestWordDiff(t *testing.T) {
	from, to := WordDiff("return a + b", "return a - total")
	wantFrom := []WordSpan{{DiffEqual, "return a "}, {DiffDelete, "+"}, {DiffEqual, " "}, {DiffDelete, "b"}}
	wantTo := []WordSpan{{DiffEqual, "return a "}, {DiffInsert, "-"}, {DiffEqual, " "}, {DiffInsert, "total"}}
	if !slices.Equal(from, wantFrom) || !slices.Equal(to, wantTo) {
		t.Fatalf("from %+v to %+v", from, to)
	}
	join := func(s []WordSpan) (out string) {
		for _, w := range s {
			out += w.Text
		}
		return
	}
	for _, p := range [][2]string{{"", "new text"}, {"old text", ""}, {"same", "same"}, {"中文 词语", "中文 句子"}} {
		f, to := WordDiff(p[0], p[1])
		if join(f) != p[0] || join(to) != p[1] {
			t.Errorf("%q→%q rebuilt %q→%q", p[0], p[1], join(f), join(to))
		}
	}
}
