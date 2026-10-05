package main

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

var (
	headingRe = regexp.MustCompile(`^(#{1,6})\s+(.*?)\s*#*\s*$`)
	listRe    = regexp.MustCompile(`^(\s*)([-*+]|\d+[.)])\s+(.*)$`)
	sepRe     = regexp.MustCompile(`^\s*\|?\s*:?-{2,}:?\s*(\|\s*:?-{2,}:?\s*)*\|?\s*$`)
	linkRe    = regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)\)`)
	boldRe    = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	slugRe    = regexp.MustCompile(`[^\p{L}\p{N}]+`)
)

// markdown converts the subset of Markdown used by the docs to HTML; link rewrites relative hrefs.
type markdown struct {
	link func(href string) string
}

func (m markdown) render(src string) string {
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")
	var b strings.Builder
	m.blocks(&b, lines)
	return b.String()
}

func (m markdown) blocks(b *strings.Builder, lines []string) {
	for i := 0; i < len(lines); {
		line := lines[i]
		trim := strings.TrimSpace(line)
		switch {
		case trim == "":
			i++
		case strings.HasPrefix(trim, "```"):
			i = m.fence(b, lines, i)
		case headingRe.MatchString(trim):
			g := headingRe.FindStringSubmatch(trim)
			n := len(g[1])
			text := m.inline(g[2])
			if n == 1 {
				b.WriteString("<h1>" + text + "</h1>\n")
			} else {
				id := slug(g[2])
				b.WriteString("<h" + strconv.Itoa(n) + ` id="` + id + `">` + text + "</h" + strconv.Itoa(n) + ">\n")
			}
			i++
		case trim == "---" || trim == "***":
			b.WriteString("<hr>\n")
			i++
		case strings.HasPrefix(trim, "|") && i+1 < len(lines) && sepRe.MatchString(lines[i+1]):
			i = m.table(b, lines, i)
		case listRe.MatchString(line):
			i = m.list(b, lines, i)
		case strings.HasPrefix(trim, ">"):
			var inner []string
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), ">") {
				inner = append(inner, strings.TrimPrefix(strings.TrimPrefix(strings.TrimSpace(lines[i]), ">"), " "))
				i++
			}
			b.WriteString("<blockquote>\n")
			m.blocks(b, inner)
			b.WriteString("</blockquote>\n")
		default:
			var para []string
			for i < len(lines) && strings.TrimSpace(lines[i]) != "" && !startsBlock(lines, i) {
				para = append(para, strings.TrimSpace(lines[i]))
				i++
			}
			if len(para) == 0 {
				para = append(para, trim)
				i++
			}
			b.WriteString("<p>" + m.inline(strings.Join(para, "\n")) + "</p>\n")
		}
	}
}

func startsBlock(lines []string, i int) bool {
	t := strings.TrimSpace(lines[i])
	return strings.HasPrefix(t, "```") || headingRe.MatchString(t) || listRe.MatchString(lines[i]) ||
		strings.HasPrefix(t, ">") || (strings.HasPrefix(t, "|") && i+1 < len(lines) && sepRe.MatchString(lines[i+1]))
}

func (m markdown) fence(b *strings.Builder, lines []string, i int) int {
	open := strings.TrimSpace(lines[i])
	lang := strings.TrimSpace(strings.TrimLeft(open, "`"))
	var code []string
	i++
	for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
		code = append(code, lines[i])
		i++
	}
	class := ""
	if lang != "" {
		class = ` class="lang-` + html.EscapeString(lang) + `"`
	}
	b.WriteString(`<div class="code"><button type="button" class="copy" aria-label="复制代码">复制</button><pre><code` +
		class + ">" + html.EscapeString(strings.Join(code, "\n")) + "</code></pre></div>\n")
	return i + 1
}

func (m markdown) table(b *strings.Builder, lines []string, i int) int {
	head := splitRow(lines[i])
	b.WriteString("<div class=\"table\"><table>\n<thead><tr>")
	for _, c := range head {
		b.WriteString("<th>" + m.inline(c) + "</th>")
	}
	b.WriteString("</tr></thead>\n<tbody>\n")
	i += 2
	for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "|") {
		b.WriteString("<tr>")
		for _, c := range splitRow(lines[i]) {
			b.WriteString("<td>" + m.inline(c) + "</td>")
		}
		b.WriteString("</tr>\n")
		i++
	}
	b.WriteString("</tbody></table></div>\n")
	return i
}

// splitRow splits a table row on pipes outside code spans; `\|` is a literal pipe.
func splitRow(line string) []string {
	s := strings.TrimSpace(line)
	s = strings.TrimPrefix(s, "|")
	if strings.HasSuffix(s, "|") && !strings.HasSuffix(s, `\|`) {
		s = s[:len(s)-1]
	}
	var cells []string
	var cur strings.Builder
	inCode := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch == '\\' && i+1 < len(s) && s[i+1] == '|':
			cur.WriteByte('|')
			i++
		case ch == '`':
			inCode = !inCode
			cur.WriteByte(ch)
		case ch == '|' && !inCode:
			cells = append(cells, strings.TrimSpace(cur.String()))
			cur.Reset()
		default:
			cur.WriteByte(ch)
		}
	}
	return append(cells, strings.TrimSpace(cur.String()))
}

func (m markdown) list(b *strings.Builder, lines []string, i int) int {
	first := listRe.FindStringSubmatch(lines[i])
	indent := len(first[1])
	ordered := first[2][0] >= '0' && first[2][0] <= '9'
	tag := "ul"
	if ordered {
		tag = "ol"
	}
	b.WriteString("<" + tag + ">\n")
	for i < len(lines) {
		g := listRe.FindStringSubmatch(lines[i])
		if g == nil || len(g[1]) != indent {
			break
		}
		body := []string{g[3]}
		i++
		for i < len(lines) {
			l := lines[i]
			if strings.TrimSpace(l) == "" {
				if i+1 < len(lines) && leading(lines[i+1]) > indent {
					body = append(body, "")
					i++
					continue
				}
				break
			}
			if leading(l) <= indent {
				break
			}
			body = append(body, strings.TrimLeft(l, " \t"))
			i++
		}
		b.WriteString("<li>")
		if len(body) == 1 || !hasBlock(body) {
			b.WriteString(m.inline(strings.Join(body, "\n")))
		} else {
			var inner strings.Builder
			m.blocks(&inner, body)
			b.WriteString(inner.String())
		}
		b.WriteString("</li>\n")
		if i < len(lines) && strings.TrimSpace(lines[i]) == "" && i+1 < len(lines) {
			if g := listRe.FindStringSubmatch(lines[i+1]); g != nil && len(g[1]) == indent {
				i++
			}
		}
	}
	b.WriteString("</" + tag + ">\n")
	return i
}

func hasBlock(body []string) bool {
	for i := range body {
		if i > 0 && (strings.TrimSpace(body[i]) == "" || startsBlock(body, i)) {
			return true
		}
	}
	return false
}

func leading(s string) int {
	return len(s) - len(strings.TrimLeft(s, " \t"))
}

// inline renders code spans, links and bold; everything else is escaped text.
func (m markdown) inline(s string) string {
	var b strings.Builder
	for {
		start := strings.IndexByte(s, '`')
		if start < 0 {
			break
		}
		ticks := 1
		for start+ticks < len(s) && s[start+ticks] == '`' {
			ticks++
		}
		delim := strings.Repeat("`", ticks)
		end := strings.Index(s[start+ticks:], delim)
		if end < 0 {
			break
		}
		b.WriteString(m.text(s[:start]))
		code := strings.TrimSpace(s[start+ticks : start+ticks+end])
		b.WriteString("<code>" + html.EscapeString(code) + "</code>")
		s = s[start+ticks+end+ticks:]
	}
	b.WriteString(m.text(s))
	return b.String()
}

func (m markdown) text(s string) string {
	s = html.EscapeString(s)
	s = linkRe.ReplaceAllStringFunc(s, func(x string) string {
		g := linkRe.FindStringSubmatch(x)
		href := html.UnescapeString(g[2])
		if m.link != nil {
			href = m.link(href)
		}
		attr := ""
		if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
			attr = ` rel="noopener"`
		}
		return `<a href="` + html.EscapeString(href) + `"` + attr + ">" + g[1] + "</a>"
	})
	s = boldRe.ReplaceAllString(s, "<strong>$1</strong>")
	return strings.ReplaceAll(s, "\n", " ")
}

func slug(s string) string {
	s = strings.Trim(slugRe.ReplaceAllString(strings.ToLower(s), "-"), "-")
	if s == "" {
		return "section"
	}
	return s
}
