package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// Source is a document an answer cites.
type Source struct {
	// Title and URL are required.
	Title, URL string
	Site       string
	Snippet    string
	Published  time.Time
}

// CitationState links citations and sources through caller-owned state.
type CitationState struct {
	// Active selects a 1-based source; zero clears selection.
	Active int
	hover  int
	owner  uint64
	seen   time.Time
	locate uint64
}

// lit reports whether source n is chosen or hovered anywhere.
func (s *CitationState) lit(n int) bool { return n > 0 && (s.Active == n || s.hover == n) }

// track refreshes pointer ownership for the current frame.
func (s *CitationState) track(c *ui.Context, e *ui.Element, n int) {
	switch over := e.Hovered() && !e.IsDisabled(); {
	case over && s.hover != n:
		s.hover, s.owner = n, e.ID()
		c.Invalidate()
	case !over && s.owner == e.ID() && s.hover != 0:
		s.hover, s.owner = 0, 0
		c.Invalidate()
	}
	if s.owner == e.ID() {
		s.seen = c.Now()
	}
}

func (s *CitationState) expire(c *ui.Context, frame time.Time) {
	if s.hover != 0 && s.seen != frame {
		s.hover, s.owner = 0, 0
		c.Invalidate()
	}
}

// CitationBadge previews and locates a numbered source.
func CitationBadge(c *ui.Context, n int, src Source, st *CitationState) ChoiceResult {
	if n < 1 || st == nil || src.Title == "" || src.URL == "" {
		panic("mujica: CitationBadge needs n >= 1, a state and a Source with Title and URL")
	}
	k := Tokens(c)
	label := fmt.Sprintf(msg(c, "citation_badge.label"), n, src.Title)
	size := FontSize(c, theme.CaptionSize-1)
	b := ui.ButtonBase(c).MinWidth(theme.MinHitSize).MinHeight(theme.MinHitSize).Padding(0, 4).Radius(theme.ControlRadius).Label(label).Transition(colorFade)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	st.track(c, b, n)
	if st.lit(n) {
		b.Background(k.Selection).Border(theme.BorderWidth, k.AccentText).TextColor(k.AccentText)
	} else {
		b.Background(k.SurfaceHover).Border(theme.BorderWidth, k.Border).TextColor(k.TextMuted)
	}
	frame := c.Now()
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		st.expire(c, frame)
		if st.lit(n) {
			p.Fill(ui.Rect{X: r.X + 4, Y: r.Y + r.H - 4, W: r.W - 8, H: 2}, k.AccentText, 0)
		}
	})
	b.Children(func() { ui.Text(c, fmt.Sprint(n)).FontSize(size).FontFeatures("tnum").FontWeight(600).SingleLine() })
	r := ChoiceResult{Element: b}
	if b.Clicked() && !b.IsDisabled() {
		st.Active, r.changed = n, true
		st.locate++
		slog.Debug("mujica: citation_badge activate", "n", n)
	}
	citationPreview(c, b, n, src)
	return r
}

// citationPreviewDelay sets the hover and focus wait.
const citationPreviewDelay = 500 * time.Millisecond

// citationPreview displays a source after resting hover or focus.
// Pass-through preserves badge clicks; Escape dismisses the preview.
func citationPreview(c *ui.Context, b *ui.Element, n int, src Source) {
	st := ui.Local(b, "mujica.citation", func() tooltipState { return tooltipState{} })
	size := ui.Local(b, "mujica.citation.size", func() ui.Rect { return ui.Rect{} })
	if b.IsDisabled() || !b.Hovered() && !b.FocusVisible() {
		*st = tooltipState{}
		return
	}
	if st.dismissed {
		return
	}
	if st.since.IsZero() {
		st.since = c.Now()
	}
	if wait := citationPreviewDelay - c.Now().Sub(st.since); wait > 0 {
		c.After(wait)
		return
	}
	if overlayEscape(c, st) {
		st.dismissed = true
		return
	}
	u := Density(c).Unit()
	ww, wh := c.Size()
	w, h := min(float32(300), max(0, ww-8)), max(0, wh-8)
	x, y := overlayPlace(c, b.Bounds(), w, min(size.H, h), 6)
	ui.Overlay(c, func() {
		panel := ui.Column(c).Absolute().Left(x).Top(y).Width(w).MaxHeight(h).Padding(u * 3).PassThrough().Role(ui.RoleTooltip).Label(msg(c, "citation_badge.preview"))
		overlayStyle(c, panel, overlayAnchored(c, panel, *size))
		*size = panel.Bounds()
		panel.Children(func() { ui.Scroll(c).Shrink(1).Children(func() { sourceSummary(c, n, src) }) })
	})
}

// sourceSummary shows a source's number, site, title, time and snippet.
func sourceSummary(c *ui.Context, n int, src Source) {
	k := Tokens(c)
	u := Density(c).Unit()
	small := FontSize(c, theme.CaptionSize)
	ui.Column(c).Gap(u).MinWidth(0).Grow(1).Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf("[%d]", n)).FontSize(small).FontFeatures("tnum").TextColor(k.AccentText).SingleLine()
			site := src.Site
			if site == "" {
				site = src.URL
			}
			ui.Text(c, site).FontSize(small).TextColor(k.TextMuted).SingleLine().Shrink(1).MinWidth(0)
			if !src.Published.IsZero() {
				ui.Text(c, src.Published.Format("2006-01-02")).FontSize(small).TextColor(k.TextMuted).FontFeatures("tnum").SingleLine()
			}
		})
		ui.Text(c, src.Title).FontWeight(600).MaxLines(2)
		if src.Snippet != "" {
			ui.Text(c, src.Snippet).FontSize(small).TextColor(k.TextMuted).MaxLines(3)
		}
	})
}
