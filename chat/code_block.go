package chat

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CopyResult is a block's element and whether its copy action ran this frame.
type CopyResult struct {
	Element *ui.Element
	copied  bool
}

// Copied reports that the block wrote its content to the clipboard this frame.
func (r CopyResult) Copied() bool { return r.copied }

// CodeBlockOptions configure a CodeBlock.
type CodeBlockOptions struct {
	// Language picks the highlighting: go, javascript, typescript, python, json, shell, sql.
	Language    string
	LineNumbers bool
	// Wrap and Collapsed set the toggles' first state.
	Wrap, Collapsed bool
}

// codeBlockState holds the user's toggles and the copy confirmation.
type codeBlockState struct {
	wrap, collapsed bool
	copiedAt        time.Time
}

// codeBlockCopiedFor is how long the copy button shows its check.
const codeBlockCopiedFor = 2 * time.Second

// CodeBlock shows code with a language label, highlighting, copy, collapse and wrap toggles.
func CodeBlock(c *ui.Context, code string, opts CodeBlockOptions) CopyResult {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r CopyResult
	r.Element = ui.Column(c).Radius(theme.CardRadius).Border(theme.BorderWidth, k.Border).Background(k.Surface).Clip().MinWidth(0)
	st := ui.Local(r.Element, "code", func() codeBlockState { return codeBlockState{wrap: opts.Wrap, collapsed: opts.Collapsed} })
	disabled := r.Element.IsDisabled()
	lines := strings.Split(strings.TrimSuffix(code, "\n"), "\n")
	r.Element.Children(func() {
		ui.Row(c).Padding(0, u, 0, u*3).Gap(u).AlignItems(ui.Center).Background(k.SurfaceHover).Children(func() {
			lang := opts.Language
			if lang == "" {
				lang = core.Msg(c, "code_block.plain")
			}
			ui.Text(c, lang).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine().Grow(1)
			if codeBlockToggle(c, "wrap-text", core.Msg(c, "code_block.wrap"), &st.wrap) {
				slog.Debug("mujica: code_block wrap", "on", st.wrap)
			}
			icon, name := "copy", core.Msg(c, "code_block.copy")
			if left := codeBlockCopiedFor - c.Now().Sub(st.copiedAt); left > 0 {
				icon, name = "check", core.Msg(c, "code_block.copied")
				c.After(left)
			}
			if core.IconAction(c, icons.Must(icon), name).Clicked() && !disabled {
				c.WriteClipboard(code)
				st.copiedAt, r.copied = c.Now(), true
				c.Announce(core.Msg(c, "code_block.copied"))
				slog.Debug("mujica: code_block copy", "lines", len(lines))
			}
			fold, foldName := "chevron-up", core.Msg(c, "code_block.collapse")
			if st.collapsed {
				fold, foldName = "chevron-down", fmt.Sprintf(core.Msg(c, "code_block.expand"), len(lines))
			}
			if core.IconAction(c, icons.Must(fold), foldName).Clicked() && !disabled {
				st.collapsed = !st.collapsed
				slog.Debug("mujica: code_block collapse", "on", st.collapsed)
			}
		})
		if st.collapsed {
			return
		}
		var body *ui.Element
		if st.wrap {
			body = ui.Column(c).Padding(u*2, u*3).MinWidth(0)
		} else {
			body = ui.ScrollHorizontal(c).Padding(u*2, u*3).AlignItems(ui.Start)
		}
		body.Font(theme.MonoFont).FontSize(core.FontSize(c, 13)).Children(func() {
			ui.Column(c).MinWidth(0).Grow(1).Children(func() { codeBlockLines(c, lines, opts, st.wrap) })
		})
	})
	return r
}

// codeBlockLines builds numbered, highlighted lines.
func codeBlockLines(c *ui.Context, lines []string, opts CodeBlockOptions, wrap bool) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	digits := len(fmt.Sprint(len(lines)))
	inBlock := false
	for i, line := range lines {
		var toks []internal.Token
		toks, inBlock = internal.Highlight(opts.Language, line, inBlock)
		ui.Row(c).Gap(u * 3).AlignItems(ui.Start).Children(func() {
			if opts.LineNumbers {
				ui.Text(c, fmt.Sprintf("%*d", digits, i+1)).TextColor(k.TextMuted).Shrink(0).NoWrap()
			}
			spans := make([]ui.Span, len(toks))
			for j, t := range toks {
				spans[j] = ui.Span{Text: t.Text, Color: codeBlockColor(k, t.Kind), Italic: t.Kind == internal.TokComment}
			}
			if len(spans) == 0 {
				spans = []ui.Span{{Text: " "}}
			}
			t := ui.RichText(c, spans...).MinWidth(0)
			if wrap {
				t.Grow(1)
			} else {
				t.NoWrap()
			}
		})
	}
}

// codeBlockColor colors a token from the theme; comments are also italic.
func codeBlockColor(k theme.Tokens, kind internal.TokenKind) ui.Color {
	switch kind {
	case internal.TokKeyword:
		return k.AccentText
	case internal.TokString:
		return k.Success
	case internal.TokNumber:
		return k.Warning
	case internal.TokComment, internal.TokPunct:
		return k.TextMuted
	case internal.TokType:
		return k.Info
	}
	return k.Text
}

// codeBlockToggle is a pressed-state icon button; on shows a Selection face and frame.
func codeBlockToggle(c *ui.Context, icon, label string, on *bool) bool {
	k := core.Tokens(c)
	b := ui.ToggleBase(c, on).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).Label(label).TextColor(k.TextMuted)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case *on:
		b.Background(k.Selection).Border(theme.BorderWidth, k.AccentText).TextColor(k.AccentText)
	case b.Hovered():
		b.Background(k.SurfaceHover).TextColor(k.Text)
	}
	b.Children(func() { ui.Icon(c, icons.Must(icon)).Size(16, 16) })
	return b.Changed()
}
