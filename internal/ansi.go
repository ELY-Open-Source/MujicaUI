package internal

import (
	"log/slog"
	"strconv"
	"strings"
)

// AnsiColorKind says how an AnsiColor is given.
type AnsiColorKind uint8

const (
	AnsiDefault AnsiColorKind = iota
	AnsiIndexed
	AnsiRGB
)

// AnsiColor is a terminal color: default, a 256-color index, or true color.
type AnsiColor struct {
	Kind    AnsiColorKind
	Index   uint8
	R, G, B uint8
}

// AnsiStyle is the SGR state a span is drawn with.
type AnsiStyle struct {
	FG, BG                                AnsiColor
	Bold, Dim, Italic, Underline, Inverse bool
	Strike                                bool
}

// AnsiSpan is a run of text in one style.
type AnsiSpan struct {
	Text  string
	Style AnsiStyle
}

// ParseAnsi splits s into styled spans from st; unknown escapes stay literal.
func ParseAnsi(s string, st AnsiStyle) ([]AnsiSpan, AnsiStyle) {
	var spans []AnsiSpan
	var b strings.Builder
	flush := func() {
		if b.Len() > 0 {
			spans = append(spans, AnsiSpan{b.String(), st})
			b.Reset()
		}
	}
	for i := 0; i < len(s); {
		if s[i] != 0x1b || i+1 >= len(s) || s[i+1] != '[' {
			b.WriteByte(s[i])
			i++
			continue
		}
		j := i + 2
		for j < len(s) && (s[j] < 0x40 || s[j] > 0x7e) {
			j++
		}
		if j >= len(s) {
			b.WriteString(s[i:])
			break
		}
		if s[j] != 'm' {
			slog.Debug("mujica: ansi unsupported sequence kept", "seq", strconv.Quote(s[i:j+1]))
			b.WriteString(s[i : j+1])
			i = j + 1
			continue
		}
		flush()
		st = applySGR(st, s[i+2:j])
		i = j + 1
	}
	flush()
	return spans, st
}

func applySGR(st AnsiStyle, params string) AnsiStyle {
	if params == "" {
		return AnsiStyle{}
	}
	codes := strings.Split(params, ";")
	for k := 0; k < len(codes); k++ {
		n, err := strconv.Atoi(codes[k])
		if err != nil {
			slog.Debug("mujica: ansi bad SGR parameter", "params", params)
			return st
		}
		switch {
		case n == 0:
			st = AnsiStyle{}
		case n == 1:
			st.Bold = true
		case n == 2:
			st.Dim = true
		case n == 3:
			st.Italic = true
		case n == 4:
			st.Underline = true
		case n == 7:
			st.Inverse = true
		case n == 9:
			st.Strike = true
		case n == 22:
			st.Bold, st.Dim = false, false
		case n == 23:
			st.Italic = false
		case n == 24:
			st.Underline = false
		case n == 27:
			st.Inverse = false
		case n == 29:
			st.Strike = false
		case n >= 30 && n <= 37:
			st.FG = AnsiColor{Kind: AnsiIndexed, Index: uint8(n - 30)}
		case n >= 90 && n <= 97:
			st.FG = AnsiColor{Kind: AnsiIndexed, Index: uint8(n - 90 + 8)}
		case n == 39:
			st.FG = AnsiColor{}
		case n >= 40 && n <= 47:
			st.BG = AnsiColor{Kind: AnsiIndexed, Index: uint8(n - 40)}
		case n >= 100 && n <= 107:
			st.BG = AnsiColor{Kind: AnsiIndexed, Index: uint8(n - 100 + 8)}
		case n == 49:
			st.BG = AnsiColor{}
		case n == 38 || n == 48:
			c, used, ok := extendedColor(codes[k+1:])
			if !ok {
				slog.Debug("mujica: ansi bad extended color", "params", params)
				return st
			}
			if n == 38 {
				st.FG = c
			} else {
				st.BG = c
			}
			k += used
		}
	}
	return st
}

// extendedColor reads "5;n" or "2;r;g;b" after 38 or 48.
func extendedColor(rest []string) (AnsiColor, int, bool) {
	num := func(i int) (uint8, bool) {
		if i >= len(rest) {
			return 0, false
		}
		v, err := strconv.Atoi(rest[i])
		return uint8(v), err == nil && v >= 0 && v <= 255
	}
	if len(rest) == 0 {
		return AnsiColor{}, 0, false
	}
	switch rest[0] {
	case "5":
		v, ok := num(1)
		return AnsiColor{Kind: AnsiIndexed, Index: v}, 2, ok
	case "2":
		r, ok1 := num(1)
		g, ok2 := num(2)
		b, ok3 := num(3)
		return AnsiColor{Kind: AnsiRGB, R: r, G: g, B: b}, 4, ok1 && ok2 && ok3
	}
	return AnsiColor{}, 0, false
}
