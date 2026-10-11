package input

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MentionCandidate is someone a MentionInput may mention.
type MentionCandidate struct {
	ID, Name string
}

// MentionInputOptions configure a MentionInput.
type MentionInputOptions struct {
	// Candidates is required: it returns matches for the text after "@", called when that text changes.
	Candidates  func(query string) []MentionCandidate
	Placeholder string
	Label       string
	Disabled    bool
	// Error marks the value invalid and shows the message below.
	Error string
}

// mentionInputState is one field's document, history and popup.
type mentionInputState struct {
	doc                mentionDoc
	synced             string // the value as last read or written
	syncedMentions     []Mention
	undo, redo         []mentionDoc
	compose            string
	query              string
	at, pick           int
	found              []MentionCandidate
	dismissed          bool
	panel              ui.Rect // the popup's box last frame, which keeps it through a click
	offset             float32
	changed, submitted bool
}

// MentionInput edits *value with "@" mentions kept in *mentions; a mention is deleted and stepped over whole.
func MentionInput(c *ui.Context, value *string, mentions *[]Mention, opts MentionInputOptions) InputResult {
	if value == nil || mentions == nil || opts.Candidates == nil {
		panic("mujica: MentionInput needs a value, mentions and Candidates")
	}
	text := []rune(*value)
	for _, m := range *mentions {
		if m.Start < 0 || m.End > len(text) || m.Start >= m.End || string(text[m.Start:m.End]) != "@"+m.Name {
			panic(fmt.Sprintf("mujica: MentionInput mention %+v does not match its text", m))
		}
	}
	k := core.Tokens(c)
	var r InputResult
	look := &textInputLook{disabled: opts.Disabled, err: opts.Error}
	r.Element, r.Input = textInputFrame(c, look, func(frame *ui.Element) *ui.Element {
		in := ui.Box(c).Grow(1).MinWidth(core.Density(c).Unit() * 20).Focusable().Clip().Role(ui.RoleTextField).Label(opts.Label).Cursor(ui.CursorText)
		s := *ui.Local(in, "mention_input", func() *mentionInputState { return &mentionInputState{at: -1} })
		if s.doc.text == nil || s.synced != *value || !slices.Equal(s.syncedMentions, *mentions) {
			// Set from outside: start over at the end.
			s.doc = mentionDoc{text: text, mentions: slices.Clone(*mentions), caret: len(text), anchor: len(text)}
			s.undo, s.redo = nil, nil
			s.synced, s.syncedMentions = *value, slices.Clone(*mentions)
		}
		mentionInputQuery(s, opts)
		px, py, _ := frame.PointerPosition()
		fb := frame.Bounds()
		onPanel := s.panel.Contains(fb.X+px, fb.Y+py)
		popup := (in.Focused() || onPanel) && s.at >= 0 && !s.dismissed && len(s.found) > 0
		in.HandleInput(func(ev ui.InputEvent) bool {
			if in.IsDisabled() {
				return false
			}
			return mentionInputEvent(c, s, ev, popup, in)
		})
		if s.changed && !look.disabled {
			*value, *mentions = string(s.doc.text), slices.Clone(s.doc.mentions)
			s.synced, s.syncedMentions = *value, slices.Clone(*mentions)
			r.changed = true
		}
		r.submitted = s.submitted && !look.disabled
		s.changed, s.submitted = false, false
		mentionInputView(c, k, s, in, opts.Placeholder)
		s.panel = ui.Rect{}
		if popup {
			mentionInputPopup(c, k, s, frame, in)
		}
		return in
	}, nil)
	if r.changed {
		slog.Debug("mujica: mention_input change", "mentions", len(*mentions))
	}
	return r
}

// mentionInputQuery refreshes the candidates when the "@" text changes.
func mentionInputQuery(s *mentionInputState, opts MentionInputOptions) {
	at, q := s.doc.query()
	if at != s.at || q != s.query {
		s.at, s.query, s.pick, s.dismissed = at, q, 0, false
		s.found = nil
		if at >= 0 {
			s.found = opts.Candidates(q)
		}
	}
}

// mentionInputEvent applies one input event to the document and reports whether it took it.
func mentionInputEvent(c *ui.Context, s *mentionInputState, ev ui.InputEvent, popup bool, in *ui.Element) bool {
	edit := func(f func()) {
		s.undo = append(s.undo, s.doc.clone())
		s.redo = nil
		f()
		s.changed = true
	}
	extend := ev.Mods&ui.Shift != 0
	switch ev.Kind {
	case ui.InputText:
		s.compose = ""
		edit(func() { s.doc.insert(ev.Text) })
	case ui.InputCompose:
		s.compose = ev.Text
	case ui.InputPointerDown:
		in.Focus()
		s.doc.moveTo(mentionInputHit(c, s, ev.X), extend)
	case ui.InputPointerMove:
		if ev.Button != 0 {
			return false
		}
		s.doc.moveTo(mentionInputHit(c, s, ev.X), true)
	case ui.InputCommand:
		a, b := s.doc.span()
		switch ev.Text {
		case "copy", "cut":
			c.WriteClipboard(string(s.doc.text[a:b]))
			if ev.Text == "cut" && a != b {
				edit(func() { s.doc.erase(false) })
			}
		case "paste":
			edit(func() { s.doc.insert(c.ReadClipboard()) })
		case "selectAll":
			s.doc.anchor, s.doc.caret = 0, len(s.doc.text)
		case "delete":
			edit(func() { s.doc.erase(true) })
		case "undo", "redo":
			from, to := &s.undo, &s.redo
			if ev.Text == "redo" {
				from, to = to, from
			}
			if n := len(*from); n > 0 {
				*to = append(*to, s.doc.clone())
				s.doc, *from = (*from)[n-1], (*from)[:n-1]
				s.changed = true
			}
		default:
			return false
		}
	case ui.InputKeyDown:
		return s.compose == "" && mentionInputKey(s, ev.Key, extend, popup, edit)
	default:
		return false
	}
	return true
}

// mentionInputKey handles editing and popup keys; others go on to Tab and shortcuts.
func mentionInputKey(s *mentionInputState, key ui.Key, extend, popup bool, edit func(func())) bool {
	if popup {
		switch key {
		case ui.KeyDown, ui.KeyUp:
			n := len(s.found)
			if key == ui.KeyDown {
				s.pick = (s.pick + 1) % n
			} else {
				s.pick = (s.pick + n - 1) % n
			}
			return true
		case ui.KeyEnter, ui.KeyTab:
			m := s.found[s.pick]
			edit(func() { s.doc.choose(s.at, m.ID, m.Name) })
			slog.Debug("mujica: mention_input mention", "id", m.ID)
			return true
		case ui.KeyEscape:
			s.dismissed = true
			return true
		}
	}
	switch key {
	case ui.KeyBackspace, ui.KeyDelete:
		if s.doc.caret == s.doc.anchor && (key == ui.KeyBackspace && s.doc.caret == 0 || key == ui.KeyDelete && s.doc.caret == len(s.doc.text)) {
			return true
		}
		edit(func() { s.doc.erase(key == ui.KeyDelete) })
	case ui.KeyLeft:
		s.doc.move(-1, extend)
	case ui.KeyRight:
		s.doc.move(1, extend)
	case ui.KeyHome:
		s.doc.moveTo(0, extend)
	case ui.KeyEnd:
		s.doc.moveTo(len(s.doc.text), extend)
	case ui.KeyEnter:
		s.submitted = true
	default:
		return false
	}
	return true
}

// mentionInputHit is the rune boundary nearest x in the field, outside any mention.
func mentionInputHit(c *ui.Context, s *mentionInputState, x float32) int {
	x += s.offset
	best, bestD := 0, float32(-1)
	for i := 0; i <= len(s.doc.text); i++ {
		w, _ := c.MeasureText(0, ui.Span{Text: string(s.doc.text[:i]), Size: core.FontSize(c, theme.BodySize)})
		if d := max(w-x, x-w); bestD < 0 || d < bestD {
			best, bestD = i, d
		}
	}
	for _, m := range s.doc.mentions {
		if best > m.Start && best < m.End {
			best = m.End
		}
	}
	return best
}

// mentionInputView draws the text with mention chips, the selection and the caret, scrolled to keep the caret in view.
func mentionInputView(c *ui.Context, k theme.Tokens, s *mentionInputState, in *ui.Element, placeholder string) {
	size := core.FontSize(c, theme.BodySize)
	measure := func(i int) float32 {
		w, _ := c.MeasureText(0, ui.Span{Text: string(s.doc.text[:i]), Size: size})
		return w
	}
	_, h := c.MeasureText(0, ui.Span{Text: "@", Size: size})
	x := measure(s.doc.caret)
	if box := in.Bounds().W; box > 0 {
		s.offset = max(0, min(s.offset, x), x-box+2)
	}
	in.MinHeight(h).TextCaret(ui.Rect{X: x - s.offset, Y: 0, W: 1, H: h})
	a, b := min(s.doc.caret, s.doc.anchor), max(s.doc.caret, s.doc.anchor)
	sel := [2]float32{measure(a), measure(b)}
	chips := make([][2]float32, len(s.doc.mentions))
	for i, m := range s.doc.mentions {
		chips[i] = [2]float32{measure(m.Start), measure(m.End)}
	}
	spans := mentionInputSpans(k, s, size, placeholder)
	in.Draw(func(p *ui.Painter, r ui.Rect) {
		y := r.Y + (r.H-h)/2
		for _, ch := range chips {
			p.Fill(ui.Rect{X: r.X + ch[0] - s.offset - 1, Y: y, W: ch[1] - ch[0] + 2, H: h}, k.Selection, theme.ControlRadius)
		}
		if a != b && in.Focused() {
			p.Fill(ui.Rect{X: r.X + sel[0] - s.offset, Y: y, W: sel[1] - sel[0], H: h}, k.Selection.Mix(k.AccentText, 0.25), 0)
		}
		p.RichText(r.X-s.offset, y, 0, spans...)
	})
	in.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if in.Focused() && s.compose == "" {
			p.Fill(ui.Rect{X: r.X + x - s.offset, Y: r.Y + (r.H-h)/2, W: 1, H: h}, k.Text, 0)
		}
	})
}

// mentionInputSpans is the field's text: mentions in AccentText, a composition underlined.
func mentionInputSpans(k theme.Tokens, s *mentionInputState, size float32, placeholder string) []ui.Span {
	if len(s.doc.text) == 0 && s.compose == "" {
		return []ui.Span{{Text: placeholder, Size: size, Color: k.TextMuted}}
	}
	var spans []ui.Span
	at := 0
	seg := func(to int, fg ui.Color) {
		if to > at {
			spans = append(spans, ui.Span{Text: string(s.doc.text[at:to]), Size: size, Color: fg})
			at = to
		}
	}
	composed := s.compose == ""
	for _, m := range s.doc.mentions {
		if !composed && s.doc.caret <= m.Start {
			seg(s.doc.caret, k.Text)
			spans = append(spans, ui.Span{Text: s.compose, Size: size, Color: k.Text, Underline: true})
			composed = true
		}
		seg(m.Start, k.Text)
		seg(m.End, k.AccentText)
	}
	if !composed {
		seg(s.doc.caret, k.Text)
		spans = append(spans, ui.Span{Text: s.compose, Size: size, Color: k.Text, Underline: true})
	}
	seg(len(s.doc.text), k.Text)
	return spans
}

// mentionInputPopup lists the candidates under the field; a click or Enter picks one.
func mentionInputPopup(c *ui.Context, k theme.Tokens, s *mentionInputState, frame, in *ui.Element) {
	b := frame.Bounds()
	u := core.Density(c).Unit()
	ui.Overlay(c, func() {
		ww, wh := c.Size()
		h := float32(len(s.found)) * core.Density(c).ControlHeight()
		x, y := core.PlaceOverlay(c, b, min(b.W, 320), min(h+2*u, wh/2), 4)
		panel := ui.Column(c).Absolute().Left(x).Top(y).Width(min(max(b.W, 200), ww-8)).MaxHeight(wh / 2).Padding(u).Role(ui.RoleList)
		core.OverlayStyle(c, panel, core.OverlayEnter(c, panel))
		s.panel = panel.Bounds()
		panel.Children(func() {
			ui.Scroll(c).Shrink(1).Children(func() {
				for i, m := range s.found {
					it := ui.Row(c).Key(m.ID).MinHeight(core.Density(c).ControlHeight()).AlignItems(ui.Center).Padding(0, u*2).
						Radius(theme.ControlRadius).Role(ui.RoleListItem).Label(m.Name).Cursor(ui.CursorPointer)
					if i == s.pick {
						it.Selected(true).Background(k.Selection).TextColor(k.AccentText)
						it.Draw(func(p *ui.Painter, r ui.Rect) {
							p.Fill(ui.Rect{X: r.X, Y: r.Y + 6, W: 2, H: r.H - 12}, k.AccentText, 0)
						})
					} else if it.Hovered() {
						it.Background(k.SurfaceHover)
					}
					it.Children(func() { ui.Text(c, "@"+m.Name).SingleLine() })
					if it.Clicked() {
						s.undo = append(s.undo, s.doc.clone())
						s.doc.choose(s.at, m.ID, m.Name)
						s.changed = true
						in.Focus()
						c.Invalidate()
					}
				}
			})
		})
	})
}
