package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// MarkdownResult is a MarkdownView's element and the link the user opened this frame.
type MarkdownResult struct {
	Element *ui.Element
	link    string
}

// Link returns the URL of the link clicked or activated this frame.
func (r MarkdownResult) Link() (string, bool) { return r.link, r.link != "" }

// MarkdownView renders a Markdown subset; a growing src parses only its last block again.
func MarkdownView(c *ui.Context, src string) MarkdownResult {
	u := Density(c).Unit()
	var r MarkdownResult
	r.Element = ui.Column(c).Gap(u * 3).MinWidth(0)
	doc := ui.Local(r.Element, "doc", func() *internal.Markdown { return &internal.Markdown{} })
	if from := (*doc).Update(src); from < len((*doc).Blocks) {
		slog.Debug("mujica: markdown_view parse", "from", from, "blocks", len((*doc).Blocks))
	}
	r.Element.Children(func() {
		for i := range (*doc).Blocks {
			b := &(*doc).Blocks[i]
			ui.Column(c).Key(i).MinWidth(0).Children(func() {
				if url := markdownBlock(c, b); url != "" {
					r.link = url
				}
			})
		}
	})
	if r.link != "" {
		slog.Debug("mujica: markdown_view link")
	}
	return r
}

// markdownBlock builds one block and returns a link opened in it.
func markdownBlock(c *ui.Context, b *internal.MdBlock) string {
	k := Tokens(c)
	u := Density(c).Unit()
	switch b.Kind {
	case internal.BlockHeading:
		size := []float32{theme.H1Size, theme.H2Size, theme.H3Size}[min(b.Level, 3)-1]
		t := ui.RichText(c).Font(theme.SerifFont).FontSize(FontSize(c, size)).FontWeight(600).Role(ui.RoleText)
		return markdownSpans(c, t, b.Spans)
	case internal.BlockParagraph:
		return markdownSpans(c, ui.RichText(c), b.Spans)
	case internal.BlockQuote:
		var link string
		ui.Row(c).Padding(0, 0, 0, u*3).BorderWidth(0, 0, 0, 3).BorderColor(k.Border).Children(func() {
			link = markdownSpans(c, ui.RichText(c).TextColor(k.TextMuted).Grow(1), b.Spans)
		})
		return link
	case internal.BlockList:
		var link string
		ui.Column(c).Gap(u).Role(ui.RoleList).Children(func() {
			for _, it := range b.Items {
				ui.Row(c).Gap(u*2).AlignItems(ui.Start).Padding(0, 0, 0, float32(it.Depth)*u*5).Role(ui.RoleListItem).Children(func() {
					mark := "•"
					if it.Number > 0 {
						mark = fmt.Sprintf("%d.", it.Number)
					}
					ui.Text(c, mark).TextColor(k.TextMuted).FontFeatures("tnum").MinWidth(u * 4).Shrink(0)
					if l := markdownSpans(c, ui.RichText(c).Grow(1).MinWidth(0), it.Spans); l != "" {
						link = l
					}
				})
			}
		})
		return link
	case internal.BlockCode:
		CodeBlock(c, b.Code, CodeBlockOptions{Language: b.Lang})
	case internal.BlockTable:
		align := make([]ui.Align, len(b.Align))
		for i, a := range b.Align {
			align[i] = []ui.Align{ui.Start, ui.Center, ui.End}[a]
		}
		TableBlock(c, b.Header, b.Rows, TableBlockOptions{Align: align})
	case internal.BlockRule:
		ui.Box(c).Height(u * 4).Draw(func(p *ui.Painter, r ui.Rect) {
			internal.RuleDiamond(p, r, u*2, 5, k.Ornament)
		})
	}
	return ""
}

// markdownSpans fills paragraph t with styled spans; links are focusable and return their URL when used.
func markdownSpans(c *ui.Context, t *ui.Element, spans []internal.MdSpan) string {
	k := Tokens(c)
	var link string
	t.MinWidth(0).Children(func() {
		for _, sp := range spans {
			e := ui.Text(c, sp.Text)
			if sp.Strong {
				e.FontWeight(700)
			}
			if sp.Em {
				e.Italic()
			}
			if sp.Strike {
				e.Strikethrough()
			}
			if sp.Code {
				e.Font(theme.MonoFont).TextBackground(k.SurfaceHover)
			}
			if sp.URL == "" {
				continue
			}
			e.TextColor(k.AccentText).Underline().Focusable().Role(ui.RoleLink).Cursor(ui.CursorPointer).Tooltip(sp.URL)
			if e.Clicked() && !e.IsDisabled() {
				link = sp.URL
			}
		}
	})
	return link
}
