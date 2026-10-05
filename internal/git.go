package internal

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CodeSpans colors a code line; inBlock carries block comments.
func CodeSpans(k theme.Tokens, language, line string, inBlock bool) ([]ui.Span, bool) {
	toks, open := Highlight(language, line, inBlock)
	spans := make([]ui.Span, 0, len(toks))
	for _, t := range toks {
		s := ui.Span{Text: t.Text, Color: k.Text}
		switch t.Kind {
		case TokKeyword:
			s.Color, s.Weight = k.AccentText, 600
		case TokString:
			s.Color = k.Success
		case TokNumber:
			s.Color = k.Info
		case TokComment:
			s.Color, s.Italic = k.TextMuted, true
		case TokType:
			s.Color = k.Warning
		case TokPunct:
			s.Color = k.TextMuted
		}
		spans = append(spans, s)
	}
	return spans, open
}

// ansiLight and ansiDark are 16 colors readable on Surface.
var (
	ansiLight = [16]string{"#211C1D", "#A12635", "#326247", "#805714", "#355E80", "#7A3B78", "#2F6B6B", "#6A5B55",
		"#6A5B55", "#B0303F", "#3B7352", "#8E6417", "#3F6E94", "#8A4787", "#367A7A", "#211C1D"}
	ansiDark = [16]string{"#A89A95", "#F0A1AA", "#94C6A4", "#DFC080", "#9DBDDA", "#D9A6D6", "#8FCFCF", "#F2EBDD",
		"#BCAEAA", "#F5B7BE", "#ADD6B9", "#E9CF99", "#B5CDE4", "#E4BDE2", "#A9DCDC", "#FFFFFF"}
	ansiCube = [6]uint8{0, 95, 135, 175, 215, 255}
)

// AnsiResolve resolves a color; ok is false for default.
func AnsiResolve(dark bool, col AnsiColor) (ui.Color, bool) {
	switch col.Kind {
	case AnsiRGB:
		return ui.RGB(col.R, col.G, col.B), true
	case AnsiIndexed:
		i := col.Index
		switch {
		case i < 16 && dark:
			return ui.Hex(ansiDark[i]), true
		case i < 16:
			return ui.Hex(ansiLight[i]), true
		case i < 232:
			i -= 16
			return ui.RGB(ansiCube[i/36], ansiCube[i/6%6], ansiCube[i%6]), true
		default:
			g := 8 + 10*(i-232)
			return ui.RGB(g, g, g), true
		}
	}
	return ui.Color{}, false
}

// StripAnsi drops SGR escapes from s, keeping its text.
func StripAnsi(s string) string {
	spans, _ := ParseAnsi(s, AnsiStyle{})
	var b strings.Builder
	for _, sp := range spans {
		b.WriteString(sp.Text)
	}
	return b.String()
}
