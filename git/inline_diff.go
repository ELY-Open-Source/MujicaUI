package git

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// InlineDiffOptions configure an InlineDiff.
type InlineDiffOptions struct {
	// Mono sets the lines in the code font.
	Mono bool
}

type inlineDiffCache struct {
	from, to string
	a, b     []internal.WordSpan
	done     bool
}

// InlineDiff shows old and new lines with changed words marked.
func InlineDiff(c *ui.Context, from, to string, opts InlineDiffOptions) ui.Element {
	k := core.Tokens(c)
	col := ui.Column(c).AlignItems(ui.Stretch).Shrink(0).Radius(theme.ControlRadius).Clip()
	cache := ui.Local(col, "words", func() *inlineDiffCache { return &inlineDiffCache{} })
	w := *cache
	if !w.done || w.from != from || w.to != to {
		w.a, w.b = internal.WordDiff(from, to)
		w.from, w.to, w.done = from, to, true
	}
	if opts.Mono {
		col.Font(theme.MonoFont)
	}
	col.Children(func() {
		inlineDiffLine(c, "−", core.Msg(c, "inline_diff.removed"), w.a, k.Danger, k)
		inlineDiffLine(c, "+", core.Msg(c, "inline_diff.added"), w.b, k.Success, k)
	})
	return col
}

// inlineDiffLine shows one side with its changed words marked.
func inlineDiffLine(c *ui.Context, sign, name string, spans []internal.WordSpan, col ui.Color, k theme.Tokens) {
	u := core.Density(c).Unit()
	out := make([]ui.Span, 0, len(spans))
	plain := ""
	for _, s := range spans {
		plain += s.Text
		sp := ui.Span{Text: s.Text}
		switch s.Op {
		case internal.DiffDelete:
			sp.Background, sp.Strikethrough, sp.DecorationColor = col.Mix(k.Surface, 0.7), true, col
		case internal.DiffInsert:
			sp.Background, sp.Underline, sp.DecorationColor = col.Mix(k.Surface, 0.7), true, col
		}
		out = append(out, sp)
	}
	if len(out) == 0 {
		out = append(out, ui.Span{Text: " "})
	}
	ui.Row(c).AlignItems(ui.Start).Gap(u*2).Padding(u, u*3).Background(col.Mix(k.Surface, 0.9)).Label(name + ": " + plain).
		Children(func() {
			ui.Text(c, sign).Font(theme.MonoFont).FontWeight(700).TextColor(col).Shrink(0)
			ui.RichText(c, out...).TextColor(k.Text).MinWidth(0)
		})
}
