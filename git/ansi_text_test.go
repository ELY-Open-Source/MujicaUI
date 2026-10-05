package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestAnsiTextStripsSGRKeepsUnknown(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		AnsiText(c, "\x1b[1;31merror\x1b[0m: bad\n\x1b[2Kcleared", AnsiTextOptions{})
	}, 400, 100)
	if !tt.HasText("error: bad") {
		t.Fatalf("styled line: %q", tt.Texts())
	}
	if !tt.HasText("\x1b[2Kcleared") {
		t.Fatalf("unknown escape not kept: %q", tt.Texts())
	}
}

func TestAnsiTextSpanStyles(t *testing.T) {
	k := theme.Light()
	parse := func(s string) ui.Span {
		spans, _ := internal.ParseAnsi(s, internal.AnsiStyle{})
		return ansiTextSpan(k, spans[0])
	}
	red, _ := internal.AnsiResolve(false, internal.AnsiColor{Kind: internal.AnsiIndexed, Index: 1})
	if s := parse("\x1b[1;3;4;31mx"); s.Color != red || s.Weight != 700 || !s.Italic || !s.Underline {
		t.Errorf("bold italic underline red: %+v", s)
	}
	if s := parse("\x1b[38;5;196mx"); s.Color != ui.RGB(255, 0, 0) {
		t.Errorf("256 color: %+v", s.Color)
	}
	if s := parse("\x1b[38;2;10;20;30;48;2;1;2;3mx"); s.Color != ui.RGB(10, 20, 30) || s.Background != ui.RGB(1, 2, 3) {
		t.Errorf("true color: %+v", s)
	}
	if s := parse("\x1b[7mx"); s.Color != k.Surface || s.Background != k.Text {
		t.Errorf("inverse: %+v", s)
	}
	if s := parse("\x1b[9mx"); !s.Strikethrough || s.Color != k.Text {
		t.Errorf("strike: %+v", s)
	}
}

func TestAnsiTextSnapshot(t *testing.T) {
	snapshot(t, "ansi_text", 520, 200, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			AnsiText(c, "\x1b[1mBuilding\x1b[0m mujica \x1b[2m(v1.1)\x1b[0m\n"+
				"\x1b[32m✓\x1b[0m 300 passed  \x1b[31m✗\x1b[0m 2 failed  \x1b[33m!\x1b[0m 1 skipped\n"+
				"\x1b[30m0\x1b[31m1\x1b[32m2\x1b[33m3\x1b[34m4\x1b[35m5\x1b[36m6\x1b[37m7\x1b[90m8\x1b[91m9\x1b[92mA\x1b[93mB\x1b[94mC\x1b[95mD\x1b[96mE\x1b[97mF\x1b[0m\n"+
				"\x1b[38;5;208m256-color\x1b[0m \x1b[38;2;125;32;52mtrue color\x1b[0m \x1b[4munderline\x1b[0m \x1b[3mitalic\x1b[0m \x1b[7m inverse \x1b[0m\n"+
				"\x1b[41;97m FAIL \x1b[0m kept: \x1b[2K", AnsiTextOptions{})
		})
	})
}
