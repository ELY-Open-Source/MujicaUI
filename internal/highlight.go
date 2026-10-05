package internal

import (
	"strings"
	"unicode"
)

// TokenKind classifies a piece of source code.
type TokenKind uint8

const (
	TokPlain TokenKind = iota
	TokKeyword
	TokString
	TokNumber
	TokComment
	TokType
	TokPunct
)

// Token is a run of source text of one kind.
type Token struct {
	Kind TokenKind
	Text string
}

// lang is a language's lexical rules.
type lang struct {
	keywords, types map[string]bool
	line            string
	blockOpen       string
	blockClose      string
	quotes          string
}

func words(s string) map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields(s) {
		m[w] = true
	}
	return m
}

var langs = map[string]lang{
	"go": {keywords: words("break case chan const continue default defer else fallthrough for func go goto if import interface map package range return select struct switch type var nil true false"),
		types: words("bool byte complex64 complex128 error float32 float64 int int8 int16 int32 int64 rune string uint uint8 uint16 uint32 uint64 uintptr any"),
		line:  "//", blockOpen: "/*", blockClose: "*/", quotes: "\"'`"},
	"javascript": {keywords: words("async await break case catch class const continue default delete do else export extends finally for from function if import in instanceof let new return static super switch this throw try typeof var void while yield null undefined true false"),
		types: words("string number boolean any unknown never object interface type enum"),
		line:  "//", blockOpen: "/*", blockClose: "*/", quotes: "\"'`"},
	"python": {keywords: words("and as assert async await break class continue def del elif else except finally for from global if import in is lambda nonlocal not or pass raise return try while with yield None True False"),
		types: words("int float str bool list dict set tuple bytes object"),
		line:  "#", quotes: "\"'"},
	"json":  {keywords: words("true false null"), quotes: "\""},
	"shell": {keywords: words("if then else elif fi for in do done while case esac function return export local echo cd"), line: "#", quotes: "\"'"},
	"sql": {keywords: words("select from where and or not insert into values update set delete create table index on join left right inner outer group by order having limit offset as null is in like distinct union all primary key foreign references"),
		line: "--", blockOpen: "/*", blockClose: "*/", quotes: "'\""},
}

// Highlight tokenizes one line of code; inBlock carries an open block comment between lines.
func Highlight(language, line string, inBlock bool) ([]Token, bool) {
	if language == "typescript" {
		language = "javascript"
	}
	l, ok := langs[language]
	if !ok {
		return []Token{{TokPlain, line}}, false
	}
	var out []Token
	emit := func(k TokenKind, s string) {
		if s == "" {
			return
		}
		if n := len(out); n > 0 && out[n-1].Kind == k {
			out[n-1].Text += s
			return
		}
		out = append(out, Token{k, s})
	}
	i := 0
	for i < len(line) {
		rest := line[i:]
		switch {
		case inBlock:
			end := strings.Index(rest, l.blockClose)
			if end < 0 {
				emit(TokComment, rest)
				return out, true
			}
			emit(TokComment, rest[:end+len(l.blockClose)])
			i += end + len(l.blockClose)
			inBlock = false
		case l.line != "" && strings.HasPrefix(rest, l.line):
			emit(TokComment, rest)
			return out, false
		case l.blockOpen != "" && strings.HasPrefix(rest, l.blockOpen):
			emit(TokComment, l.blockOpen)
			i += len(l.blockOpen)
			inBlock = true
		case strings.ContainsRune(l.quotes, rune(rest[0])):
			n := quotedLen(rest)
			emit(TokString, rest[:n])
			i += n
		case rest[0] >= '0' && rest[0] <= '9':
			n := spanLen(rest, func(r rune) bool { return unicode.IsDigit(r) || unicode.IsLetter(r) || r == '.' || r == '_' })
			emit(TokNumber, rest[:n])
			i += n
		case isIdent(rune(rest[0])):
			n := spanLen(rest, func(r rune) bool { return isIdent(r) || unicode.IsDigit(r) })
			w := rest[:n]
			switch {
			case l.keywords[w] || l.keywords[strings.ToLower(w)] && language == "sql":
				emit(TokKeyword, w)
			case l.types[w]:
				emit(TokType, w)
			default:
				emit(TokPlain, w)
			}
			i += n
		case strings.ContainsRune("{}()[];,.:=+-*/<>!&|%^~?", rune(rest[0])):
			emit(TokPunct, rest[:1])
			i++
		default:
			n := spanLen(rest, func(r rune) bool { return r < 128 && unicode.IsSpace(r) })
			if n == 0 {
				n = len(string([]rune(rest)[0]))
			}
			emit(TokPlain, rest[:n])
			i += n
		}
	}
	return out, inBlock
}

func isIdent(r rune) bool { return r == '_' || r == '$' || unicode.IsLetter(r) }

func spanLen(s string, ok func(rune) bool) int {
	for i, r := range s {
		if !ok(r) {
			return i
		}
	}
	return len(s)
}

// quotedLen measures a string literal at s[0], honoring backslash escapes.
func quotedLen(s string) int {
	q := s[0]
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '\\':
			if q != '`' {
				i++
			}
		case q:
			return i + 1
		}
	}
	return len(s)
}
