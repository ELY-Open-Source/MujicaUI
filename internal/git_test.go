package internal

import (
	"testing"

	"mujicaui/theme"
)

func TestAnsiPaletteReadable(t *testing.T) {
	for _, k := range []theme.Tokens{theme.Light(), theme.Dark()} {
		for i := range 16 {
			col, ok := AnsiResolve(k.Dark, AnsiColor{Kind: AnsiIndexed, Index: uint8(i)})
			if !ok {
				t.Fatalf("index %d unresolved", i)
			}
			if r := theme.Contrast(col, k.Surface); r < 4.5 {
				t.Errorf("dark=%v color %d contrast %.2f", k.Dark, i, r)
			}
		}
	}
	if c, _ := AnsiResolve(false, AnsiColor{Kind: AnsiIndexed, Index: 196}); c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("cube 196 = %v", c)
	}
	if c, _ := AnsiResolve(false, AnsiColor{Kind: AnsiIndexed, Index: 244}); c.R != 128 {
		t.Errorf("gray 244 = %v", c)
	}
	if _, ok := AnsiResolve(true, AnsiColor{}); ok {
		t.Error("default color resolved")
	}
}

func TestCodeSpansColorKeywords(t *testing.T) {
	k := theme.Light()
	spans, _ := CodeSpans(k, "go", "func main() // hi", false)
	if spans[0].Text != "func" || spans[0].Color != k.AccentText {
		t.Fatalf("keyword span %+v", spans[0])
	}
	if last := spans[len(spans)-1]; last.Color != k.TextMuted || !last.Italic {
		t.Fatalf("comment span %+v", last)
	}
}
