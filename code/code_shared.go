package code

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// CodeSize is the design size of code and terminal text.
const CodeSize float32 = 13

// codeFont is the scaled monospace code font.
func codeFont(c *ui.Context) ui.Font {
	return ui.Font{Family: theme.MonoFont, Size: core.FontSize(c, CodeSize)}
}

// CodeMono is a single-line monospace text.
func CodeMono(c *ui.Context, s string) ui.Element {
	return ui.Text(c, s).Font(theme.MonoFont).FontSize(core.FontSize(c, CodeSize)).SingleLine()
}

// CodeCaption is muted caption text.
func CodeCaption(c *ui.Context, s string) ui.Element {
	return ui.Text(c, s).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(core.Tokens(c).TextMuted).SingleLine()
}

// CodeTokenColor colors a syntax token kind.
func CodeTokenColor(k theme.Tokens, kind internal.TokenKind) ui.Color {
	switch kind {
	case internal.TokKeyword:
		return k.AccentText
	case internal.TokString:
		return k.Success
	case internal.TokNumber:
		return k.Warning
	case internal.TokComment:
		return k.TextMuted
	case internal.TokType:
		return k.Info
	}
	return k.Text
}

// codeAnsiColor maps a terminal color onto the theme.
func codeAnsiColor(k theme.Tokens, a internal.AnsiColor) ui.Color {
	switch a.Kind {
	case internal.AnsiRGB:
		return ui.RGB(a.R, a.G, a.B)
	case internal.AnsiIndexed:
		if a.Index < 16 {
			base := [8]ui.Color{k.TextMuted, k.Danger, k.Success, k.Warning, k.Info, k.AccentText, k.Info.Mix(k.Success, 0.5), k.Text}
			return base[a.Index%8]
		}
		if a.Index >= 232 {
			v := uint8(8 + 10*int(a.Index-232))
			return ui.RGB(v, v, v)
		}
		i := int(a.Index) - 16
		step := func(n int) uint8 {
			if n == 0 {
				return 0
			}
			return uint8(55 + 40*n)
		}
		return ui.RGB(step(i/36), step(i/6%6), step(i%6))
	}
	return ui.Color{}
}

// CodeAnsiSpans turns ANSI text into themed spans.
func CodeAnsiSpans(c *ui.Context, line string) []ui.Span {
	k := core.Tokens(c)
	parts, _ := internal.ParseAnsi(line, internal.AnsiStyle{})
	spans := make([]ui.Span, 0, len(parts))
	for _, p := range parts {
		s := ui.Span{Text: p.Text, Color: codeAnsiColor(k, p.Style.FG), Underline: p.Style.Underline, Strikethrough: p.Style.Strike, Italic: p.Style.Italic}
		if p.Style.Bold {
			s.Weight = 700
		}
		if bg := codeAnsiColor(k, p.Style.BG); bg.A > 0 {
			s.Background = bg.Alpha(0.3)
		}
		spans = append(spans, s)
	}
	return spans
}

// CodeMatchSpans marks each query match in text.
func CodeMatchSpans(c *ui.Context, text, query string) []ui.Span {
	if query == "" {
		return []ui.Span{{Text: text}}
	}
	k := core.Tokens(c)
	lower, q := strings.ToLower(text), strings.ToLower(query)
	var spans []ui.Span
	for {
		i := strings.Index(lower, q)
		if i < 0 {
			break
		}
		if i > 0 {
			spans = append(spans, ui.Span{Text: text[:i]})
		}
		spans = append(spans, ui.Span{Text: text[i : i+len(q)], Background: k.Selection, Underline: true, DecorationColor: k.AccentText, Weight: 600})
		text, lower = text[i+len(q):], lower[i+len(q):]
	}
	if text != "" {
		spans = append(spans, ui.Span{Text: text})
	}
	return spans
}

// CodePanel is a framed surface holding a code tool.
func CodePanel(c *ui.Context) ui.Element {
	k := core.Tokens(c)
	return ui.Column(c).MinHeight(0).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius).Clip()
}

// CodeHeader is a panel's top strip of tools.
func CodeHeader(c *ui.Context, fn func()) ui.Element {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	return ui.Row(c).Shrink(0).MinHeight(theme.MinHitSize+2*u).Padding(u, u*2).Gap(u*2).AlignItems(ui.Center).
		BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).Children(fn)
}

// CodeSelectMark draws a chosen row's diamond.
func CodeSelectMark(c *ui.Context, on bool) {
	k := core.Tokens(c)
	ui.Box(c).Size(8, 8).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		if on {
			p.FillPath(internal.Diamond(internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2), r.W), k.AccentText)
		}
	})
}

// CodeRows builds a growing virtual list with marks.
func CodeRows(c *ui.Context, s *ui.ListState, n int, row func(i int)) ui.Element {
	d := core.Density(c)
	restore := data.DataSelectTheme(c, 0)
	defer restore()
	e := ui.List(c, s, n, func(i int) {
		ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			CodeSelectMark(c, s.Selected != nil && *s.Selected == i)
			row(i)
		})
	})
	k := core.Tokens(c)
	internal.Ring(e, k.Focus, theme.ControlRadius)
	return e.Grow(1).MinHeight(d.ControlHeight() * 3)
}

// CodeBadge is a kind letter in a framed square.
func CodeBadge(c *ui.Context, letter string, col ui.Color) {
	ui.Box(c).Size(18, 18).Shrink(0).Center().Radius(theme.ControlRadius).Border(theme.BorderWidth, col).
		Role(ui.RoleNone).Children(func() {
		ui.Text(c, letter).Font(theme.MonoFont).FontSize(core.FontSize(c, 11)).FontWeight(700).TextColor(col)
	})
}

// codeTextWidth returns the width of s in the code font.
func codeTextWidth(c *ui.Context, s string) float32 {
	w, _ := c.MeasureText(0, ui.Span{Text: s, Font: theme.MonoFont, Size: core.FontSize(c, CodeSize)})
	return w
}
