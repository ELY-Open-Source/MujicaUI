package internal

import "testing"

func TestParseAnsi(t *testing.T) {
	spans, end := ParseAnsi("ok \x1b[1;31mfail\x1b[0m \x1b[38;2;10;20;30mrgb\x1b[38;5;200mix", AnsiStyle{})
	want := []AnsiSpan{
		{"ok ", AnsiStyle{}},
		{"fail", AnsiStyle{Bold: true, FG: AnsiColor{Kind: AnsiIndexed, Index: 1}}},
		{" ", AnsiStyle{}},
		{"rgb", AnsiStyle{FG: AnsiColor{Kind: AnsiRGB, R: 10, G: 20, B: 30}}},
		{"ix", AnsiStyle{FG: AnsiColor{Kind: AnsiIndexed, Index: 200}}},
	}
	if len(spans) != len(want) {
		t.Fatalf("spans %+v", spans)
	}
	for i := range want {
		if spans[i] != want[i] {
			t.Errorf("span %d = %+v, want %+v", i, spans[i], want[i])
		}
	}
	if end.FG.Index != 200 {
		t.Errorf("end style %+v", end)
	}
}

func TestParseAnsiCarriesStyleAcrossLines(t *testing.T) {
	_, st := ParseAnsi("\x1b[92mgreen", AnsiStyle{})
	spans, _ := ParseAnsi("still", st)
	if spans[0].Style.FG != (AnsiColor{Kind: AnsiIndexed, Index: 10}) {
		t.Fatalf("style lost: %+v", spans[0].Style)
	}
}

func TestParseAnsiKeepsUnknownSequences(t *testing.T) {
	spans, _ := ParseAnsi("a\x1b[2Kb", AnsiStyle{})
	if len(spans) != 1 || spans[0].Text != "a\x1b[2Kb" {
		t.Fatalf("spans %+v", spans)
	}
}
