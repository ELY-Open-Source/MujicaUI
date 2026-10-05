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
	Site      string
	Snippet   string
	Published time.Time
}

// CitationState links citations, source cards and the sources panel; the caller holds it.
type CitationState struct {
	// Active is the 1-based number of the source the user chose; 0 for none.
	Active int
	hover  int
	owner  uint64
}

// lit reports whether source n is chosen or hovered anywhere.
func (s *CitationState) lit(n int) bool { return n > 0 && (s.Active == n || s.hover == n) }

// track records that e, showing source n, is or is no longer under the pointer.
func (s *CitationState) track(c *ui.Context, e *ui.Element, n int) {
	switch over := e.Hovered(); {
	case over && s.hover != n:
		s.hover, s.owner = n, e.ID()
		c.Invalidate()
	case !over && s.owner == e.ID() && s.hover != 0:
		s.hover, s.owner = 0, 0
		c.Invalidate()
	}
}

// CitationBadge is an inline source number; hovering previews the source, clicking makes it Active.
func CitationBadge(c *ui.Context, n int, src Source, st *CitationState) ChoiceResult {
	if n < 1 || st == nil || src.Title == "" || src.URL == "" {
		panic("mujica: CitationBadge needs n >= 1, a state and a Source with Title and URL")
	}
	k := Tokens(c)
	label := fmt.Sprintf(msg(c, "citation_badge.label"), n, src.Title)
	size := FontSize(c, theme.CaptionSize-1)
	b := ui.ButtonBase(c).MinWidth(18).Height(18).Padding(0, 4).Radius(theme.ControlRadius).Label(label).Transition(colorFade)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	st.track(c, b, n)
	if st.lit(n) {
		b.Background(k.Selection).Border(theme.BorderWidth, k.AccentText).TextColor(k.AccentText)
	} else {
		b.Background(k.SurfaceHover).Border(theme.BorderWidth, k.Border).TextColor(k.TextMuted)
	}
	b.Children(func() { ui.Text(c, fmt.Sprint(n)).FontSize(size).FontFeatures("tnum").FontWeight(600).SingleLine() })
	r := ChoiceResult{Element: b}
	if b.Clicked() && !b.IsDisabled() && st.Active != n {
		st.Active, r.changed = n, true
		slog.Debug("mujica: citation_badge activate", "n", n)
	}
	citationPreview(c, b, n, src)
	return r
}

// citationPreviewDelay is how long hover or focus rests before the preview shows.
const citationPreviewDelay = 500 * time.Millisecond

// citationPreview shows the source by badge b once hover or keyboard focus rests on it.
// It lets the pointer through, so the badge stays clickable under it; Escape hides it.
func citationPreview(c *ui.Context, b *ui.Element, n int, src Source) {
	st := ui.Local(b, "mujica.citation", func() tooltipState { return tooltipState{} })
	size := ui.Local(b, "mujica.citation.size", func() ui.Rect { return ui.Rect{} })
	if !b.Hovered() && !b.FocusVisible() {
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
	w := float32(300)
	x, y := overlayPlace(c, b.Bounds(), w, size.H, 6)
	ui.Overlay(c, func() {
		panel := ui.Column(c).Absolute().Left(x).Top(y).Width(w).Padding(u * 3).PassThrough().Role(ui.RoleTooltip)
		overlayStyle(c, panel, overlayAnchored(c, panel, *size))
		*size = panel.Bounds()
		panel.Children(func() { sourceSummary(c, n, src) })
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
