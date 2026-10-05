package internal

import "strings"

// mdStyle is the inline style in effect.
type mdStyle struct{ strong, em, strike bool }

// ParseInline splits text into styled spans: emphasis, code, strikethrough and links.
// An opener without its closer stays literal, as while a message streams in.
func ParseInline(s string) []MdSpan {
	var out []MdSpan
	parseInline(s, mdStyle{}, &out)
	return out
}

func mdEmit(out *[]MdSpan, sp MdSpan) {
	if sp.Text == "" {
		return
	}
	if n := len(*out); n > 0 {
		last := &(*out)[n-1]
		if sp.URL == "" && last.URL == "" && last.Strong == sp.Strong && last.Em == sp.Em && last.Code == sp.Code && last.Strike == sp.Strike {
			last.Text += sp.Text
			return
		}
	}
	*out = append(*out, sp)
}

func parseInline(s string, st mdStyle, out *[]MdSpan) {
	plain := func(t string) {
		mdEmit(out, MdSpan{Text: t, Strong: st.strong, Em: st.em, Strike: st.strike})
	}
	for i := 0; i < len(s); {
		switch c := s[i]; {
		case c == '\\' && i+1 < len(s) && strings.IndexByte("\\`*_~[]()#!>-|", s[i+1]) >= 0:
			plain(s[i+1 : i+2])
			i += 2
			continue
		case c == '`':
			if end := strings.IndexByte(s[i+1:], '`'); end >= 0 {
				mdEmit(out, MdSpan{Text: s[i+1 : i+1+end], Code: true, Strong: st.strong, Em: st.em, Strike: st.strike})
				i += end + 2
				continue
			}
		case c == '[':
			if text, url, n := mdLink(s[i:]); n > 0 {
				mdEmit(out, MdSpan{Text: text, URL: url, Strong: st.strong, Em: st.em, Strike: st.strike})
				i += n
				continue
			}
		case c == '*' || c == '_' || c == '~':
			if n := mdDelimited(s, i, st, out); n > 0 {
				i += n
				continue
			}
		}
		j := i + 1
		for j < len(s) && strings.IndexByte("\\`[*_~", s[j]) < 0 {
			j++
		}
		plain(s[i:j])
		i = j
	}
}

// mdDelimited parses emphasis or strikethrough opening at i and returns its length, or 0.
func mdDelimited(s string, i int, st mdStyle, out *[]MdSpan) int {
	c := s[i]
	n := 1
	if i+1 < len(s) && s[i+1] == c {
		n = 2
	}
	if c == '~' && n != 2 {
		return 0
	}
	// Intraword underscores stay literal, as in snake_case.
	if c == '_' && i > 0 && mdWord(s[i-1]) {
		return 0
	}
	marker := s[i : i+n]
	body := i + n
	if body >= len(s) || s[body] == ' ' {
		return 0
	}
	end := -1
	for j := body + 1; j+n <= len(s); j++ {
		if s[j:j+n] == marker && s[j-1] != ' ' && (n == 2 || j+1 >= len(s) || s[j+1] != c) {
			end = j
			break
		}
	}
	if end < 0 {
		return 0
	}
	inner := st
	switch {
	case c == '~':
		inner.strike = true
	case n == 2:
		inner.strong = true
	default:
		inner.em = true
	}
	parseInline(s[body:end], inner, out)
	return end + n - i
}

func mdWord(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
}

// mdLink parses [text](url) at the start of s and returns its length, or 0.
func mdLink(s string) (text, url string, n int) {
	close := strings.Index(s, "](")
	if close < 1 || strings.IndexByte(s[1:close], '[') >= 0 {
		return "", "", 0
	}
	end := strings.IndexByte(s[close+2:], ')')
	if end < 0 {
		return "", "", 0
	}
	url = strings.TrimSpace(s[close+2 : close+2+end])
	if url == "" || strings.ContainsAny(url, " \n") {
		return "", "", 0
	}
	return s[1:close], url, close + 3 + end
}
