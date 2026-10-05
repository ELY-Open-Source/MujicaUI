package internal

import (
	"strconv"
	"strings"
)

// BlockKind is the kind of a Markdown block.
type BlockKind uint8

const (
	BlockParagraph BlockKind = iota
	BlockHeading
	BlockList
	BlockQuote
	BlockCode
	BlockTable
	BlockRule
)

// MdSpan is a run of inline text with one style.
type MdSpan struct {
	Text string
	// URL is set for a link.
	URL                      string
	Strong, Em, Code, Strike bool
}

// MdItem is one list item.
type MdItem struct {
	// Number is the ordinal of an ordered item, 0 for a bullet.
	Number int
	Depth  int
	Spans  []MdSpan
}

// TableAlign is a table column's alignment.
type TableAlign uint8

const (
	AlignLeft TableAlign = iota
	AlignCenter
	AlignRight
)

// MdBlock is one parsed block.
type MdBlock struct {
	Kind BlockKind
	// Level is a heading's level, 1 to 6.
	Level int
	// Spans hold a paragraph's, heading's or quote's text.
	Spans []MdSpan
	Items []MdItem
	// Lang and Code hold a fenced code block; Open is true until its fence closes.
	Lang, Code string
	Open       bool
	Header     []string
	Align      []TableAlign
	Rows       [][]string
}

// Markdown is a document parsed incrementally as its source grows.
type Markdown struct {
	src    string
	Blocks []MdBlock
	starts []int
}

// Update parses src and returns the index of the first block parsed anew.
// When src extends the last source, blocks before the last stay as they are.
func (m *Markdown) Update(src string) int {
	if src == m.src && m.starts != nil {
		return len(m.Blocks)
	}
	from := 0
	if len(m.Blocks) > 0 && strings.HasPrefix(src, m.src) {
		from = len(m.Blocks) - 1
	}
	off := 0
	if from > 0 {
		off = m.starts[from]
	}
	m.Blocks, m.starts = m.Blocks[:from], m.starts[:from]
	m.src = src
	parseBlocks(m, src, off)
	return from
}

// mdLine is one source line and its byte offset.
type mdLine struct {
	text string
	off  int
}

func parseBlocks(m *Markdown, src string, off int) {
	var lines []mdLine
	for off < len(src) {
		end := strings.IndexByte(src[off:], '\n')
		if end < 0 {
			lines = append(lines, mdLine{src[off:], off})
			break
		}
		lines = append(lines, mdLine{strings.TrimRight(src[off:off+end], "\r"), off})
		off += end + 1
	}
	for i := 0; i < len(lines); {
		l := lines[i].text
		if strings.TrimSpace(l) == "" {
			i++
			continue
		}
		start := lines[i].off
		var b MdBlock
		switch {
		case mdFence(l) != "":
			b, i = parseFenced(lines, i)
		case mdHeadingLevel(l) > 0:
			lv := mdHeadingLevel(l)
			b = MdBlock{Kind: BlockHeading, Level: lv, Spans: ParseInline(strings.TrimSpace(strings.TrimRight(l[lv:], "# ")))}
			i++
		case mdIsRule(l):
			b = MdBlock{Kind: BlockRule}
			i++
		case strings.HasPrefix(strings.TrimLeft(l, " "), ">"):
			b, i = parseQuote(lines, i)
		case mdListMarker(l) != nil:
			b, i = parseList(lines, i)
		case i+1 < len(lines) && strings.Contains(l, "|") && mdDelimiter(lines[i+1].text) != nil:
			b, i = parseTable(lines, i)
		default:
			b, i = parseParagraph(lines, i)
		}
		m.Blocks = append(m.Blocks, b)
		m.starts = append(m.starts, start)
	}
}

// mdStarts reports whether line l begins a block other than a paragraph.
func mdStarts(lines []mdLine, i int) bool {
	l := lines[i].text
	return strings.TrimSpace(l) == "" || mdFence(l) != "" || mdHeadingLevel(l) > 0 || mdIsRule(l) ||
		strings.HasPrefix(strings.TrimLeft(l, " "), ">") || mdListMarker(l) != nil ||
		i+1 < len(lines) && strings.Contains(l, "|") && mdDelimiter(lines[i+1].text) != nil
}

func parseParagraph(lines []mdLine, i int) (MdBlock, int) {
	var parts []string
	for ; i < len(lines) && (len(parts) == 0 || !mdStarts(lines, i)); i++ {
		parts = append(parts, strings.TrimSpace(lines[i].text))
	}
	return MdBlock{Kind: BlockParagraph, Spans: ParseInline(strings.Join(parts, " "))}, i
}

// mdFence returns the fence a line opens, "```" or "~~~", or "".
func mdFence(l string) string {
	t := strings.TrimLeft(l, " ")
	for _, f := range []string{"```", "~~~"} {
		if strings.HasPrefix(t, f) {
			return f
		}
	}
	return ""
}

func parseFenced(lines []mdLine, i int) (MdBlock, int) {
	fence := mdFence(lines[i].text)
	b := MdBlock{Kind: BlockCode, Open: true, Lang: strings.ToLower(strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(lines[i].text), fence[:1])))}
	var code []string
	for i++; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i].text), fence) && strings.Trim(strings.TrimSpace(lines[i].text), fence[:1]) == "" {
			b.Open = false
			i++
			break
		}
		code = append(code, lines[i].text)
	}
	b.Code = strings.Join(code, "\n")
	return b, i
}

func mdHeadingLevel(l string) int {
	n := 0
	for n < len(l) && l[n] == '#' {
		n++
	}
	if n == 0 || n > 6 || n < len(l) && l[n] != ' ' {
		return 0
	}
	return n
}

func mdIsRule(l string) bool {
	t := strings.ReplaceAll(strings.TrimSpace(l), " ", "")
	if len(t) < 3 {
		return false
	}
	return strings.Count(t, t[:1]) == len(t) && strings.ContainsAny(t[:1], "-*_")
}

func parseQuote(lines []mdLine, i int) (MdBlock, int) {
	var parts []string
	for ; i < len(lines); i++ {
		t := strings.TrimLeft(lines[i].text, " ")
		if !strings.HasPrefix(t, ">") {
			break
		}
		parts = append(parts, strings.TrimSpace(t[1:]))
	}
	return MdBlock{Kind: BlockQuote, Spans: ParseInline(strings.Join(parts, " "))}, i
}

// mdMarker is a list item's marker.
type mdMarker struct {
	number, depth, width int
}

func mdListMarker(l string) *mdMarker {
	indent := len(l) - len(strings.TrimLeft(l, " "))
	t := l[indent:]
	if len(t) >= 2 && strings.ContainsRune("-*+", rune(t[0])) && t[1] == ' ' {
		return &mdMarker{depth: indent / 2, width: indent + 2}
	}
	d := 0
	for d < len(t) && t[d] >= '0' && t[d] <= '9' {
		d++
	}
	if d == 0 || d > 9 || d+1 >= len(t) || !strings.ContainsRune(".)", rune(t[d])) || t[d+1] != ' ' {
		return nil
	}
	n, _ := strconv.Atoi(t[:d])
	return &mdMarker{number: max(n, 1), depth: indent / 2, width: indent + d + 2}
}

func parseList(lines []mdLine, i int) (MdBlock, int) {
	b := MdBlock{Kind: BlockList}
	var texts []string
	for ; i < len(lines); i++ {
		l := lines[i].text
		if mk := mdListMarker(l); mk != nil {
			b.Items = append(b.Items, MdItem{Number: mk.number, Depth: mk.depth})
			texts = append(texts, strings.TrimSpace(l[mk.width:]))
			continue
		}
		if mdStarts(lines, i) {
			break
		}
		// A lazy line continues the item above.
		texts[len(texts)-1] += " " + strings.TrimSpace(l)
	}
	for j := range b.Items {
		b.Items[j].Spans = ParseInline(texts[j])
	}
	return b, i
}

// mdCells splits a table row into trimmed cells.
func mdCells(l string) []string {
	t := strings.TrimSpace(l)
	t = strings.TrimSuffix(strings.TrimPrefix(t, "|"), "|")
	cells := strings.Split(t, "|")
	for i := range cells {
		cells[i] = strings.TrimSpace(cells[i])
	}
	return cells
}

// mdDelimiter parses a table's delimiter row, or returns nil.
func mdDelimiter(l string) []TableAlign {
	if !strings.Contains(l, "-") {
		return nil
	}
	var out []TableAlign
	for _, c := range mdCells(l) {
		left, right := strings.HasPrefix(c, ":"), strings.HasSuffix(c, ":")
		if strings.Trim(c, ":") == "" || strings.Trim(c, ":-") != "" {
			return nil
		}
		switch {
		case left && right:
			out = append(out, AlignCenter)
		case right:
			out = append(out, AlignRight)
		default:
			out = append(out, AlignLeft)
		}
	}
	return out
}

func parseTable(lines []mdLine, i int) (MdBlock, int) {
	b := MdBlock{Kind: BlockTable, Header: mdCells(lines[i].text), Align: mdDelimiter(lines[i+1].text)}
	n := len(b.Header)
	for len(b.Align) < n {
		b.Align = append(b.Align, AlignLeft)
	}
	b.Align = b.Align[:n]
	for i += 2; i < len(lines) && strings.Contains(lines[i].text, "|") && strings.TrimSpace(lines[i].text) != ""; i++ {
		row := mdCells(lines[i].text)
		for len(row) < n {
			row = append(row, "")
		}
		b.Rows = append(b.Rows, row[:n])
	}
	return b, i
}
