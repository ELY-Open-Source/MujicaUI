package input

import (
	"slices"
	"unicode"
)

// Mention is a whole "@Name" token in a MentionInput's text, at rune offsets [Start, End).
type Mention struct {
	ID         string
	Name       string
	Start, End int
}

// mentionDoc is a MentionInput's text, mentions, caret and selection anchor, in runes.
type mentionDoc struct {
	text          []rune
	mentions      []Mention
	caret, anchor int
}

func (d mentionDoc) clone() mentionDoc {
	d.text = slices.Clone(d.text)
	d.mentions = slices.Clone(d.mentions)
	return d
}

// span returns the selection, widened to whole mentions it touches.
func (d *mentionDoc) span() (int, int) {
	a, b := min(d.caret, d.anchor), max(d.caret, d.anchor)
	for _, m := range d.mentions {
		if m.Start < b && m.End > a {
			a, b = min(a, m.Start), max(b, m.End)
		}
	}
	return a, b
}

// replace puts ins in place of [a, b), dropping mentions inside and shifting later ones; the caret lands after ins.
func (d *mentionDoc) replace(a, b int, ins []rune) {
	d.text = slices.Concat(d.text[:a], ins, d.text[b:])
	shift := len(ins) - (b - a)
	kept := d.mentions[:0]
	for _, m := range d.mentions {
		switch {
		case m.End <= a:
			kept = append(kept, m)
		case m.Start >= b:
			m.Start, m.End = m.Start+shift, m.End+shift
			kept = append(kept, m)
		}
	}
	d.mentions = kept
	d.caret = a + len(ins)
	d.anchor = d.caret
}

// insert types s at the caret, over any selection.
func (d *mentionDoc) insert(s string) {
	a, b := d.span()
	d.replace(a, b, []rune(s))
}

// erase deletes the selection, else one rune or one whole mention back (or forward).
func (d *mentionDoc) erase(forward bool) bool {
	if d.caret != d.anchor {
		a, b := d.span()
		d.replace(a, b, nil)
		return true
	}
	for _, m := range d.mentions {
		if !forward && m.End == d.caret || forward && m.Start == d.caret {
			d.replace(m.Start, m.End, nil)
			return true
		}
	}
	switch {
	case !forward && d.caret > 0:
		d.replace(d.caret-1, d.caret, nil)
	case forward && d.caret < len(d.text):
		d.replace(d.caret, d.caret+1, nil)
	default:
		return false
	}
	return true
}

// move steps the caret by dir runes, over whole mentions; extend keeps the anchor.
func (d *mentionDoc) move(dir int, extend bool) {
	at := max(0, min(len(d.text), d.caret+dir))
	for _, m := range d.mentions {
		if at > m.Start && at < m.End {
			if dir < 0 {
				at = m.Start
			} else {
				at = m.End
			}
		}
	}
	d.moveTo(at, extend)
}

func (d *mentionDoc) moveTo(at int, extend bool) {
	d.caret = at
	if !extend {
		d.anchor = at
	}
}

// query returns where an "@" being typed starts and what follows it, or -1.
func (d *mentionDoc) query() (int, string) {
	if d.caret != d.anchor {
		return -1, ""
	}
	for i := d.caret - 1; i >= 0; i-- {
		r := d.text[i]
		if unicode.IsSpace(r) {
			return -1, ""
		}
		if r != '@' {
			continue
		}
		for _, m := range d.mentions {
			if i >= m.Start && i < m.End {
				return -1, ""
			}
		}
		if i > 0 && !unicode.IsSpace(d.text[i-1]) {
			return -1, ""
		}
		return i, string(d.text[i+1 : d.caret])
	}
	return -1, ""
}

// choose turns the "@query" before the caret into a mention of id, followed by a space.
func (d *mentionDoc) choose(at int, id, name string) {
	tok := []rune("@" + name)
	d.replace(at, d.caret, append(slices.Clone(tok), ' '))
	d.mentions = append(d.mentions, Mention{ID: id, Name: name, Start: at, End: at + len(tok)})
	slices.SortFunc(d.mentions, func(a, b Mention) int { return a.Start - b.Start })
}
