package git

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// AnsiTextOptions configure an AnsiText.
type AnsiTextOptions struct {
	// Wrap breaks long lines at the width.
	Wrap bool
}

// AnsiText renders terminal output with its SGR styles.
func AnsiText(c *ui.Context, text string, opts AnsiTextOptions) *ui.Element {
	k := core.Tokens(c)
	col := ui.Column(c).Font(theme.MonoFont).TextColor(k.Text).Shrink(0)
	col.Children(func() {
		var st internal.AnsiStyle
		for _, line := range strings.Split(text, "\n") {
			var spans []internal.AnsiSpan
			spans, st = internal.ParseAnsi(line, st)
			out := make([]ui.Span, 0, len(spans))
			for _, s := range spans {
				out = append(out, ansiTextSpan(k, s))
			}
			if len(out) == 0 {
				out = append(out, ui.Span{Text: " "})
			}
			t := ui.RichText(c, out...)
			if !opts.Wrap {
				t.NoWrap()
			}
		}
	})
	return col
}

// ansiTextSpan maps one parsed span onto theme-readable colors.
func ansiTextSpan(k theme.Tokens, s internal.AnsiSpan) ui.Span {
	fg, hasFG := internal.AnsiResolve(k.Dark, s.Style.FG)
	if !hasFG {
		fg = k.Text
	}
	bg, hasBG := internal.AnsiResolve(k.Dark, s.Style.BG)
	if s.Style.Inverse {
		if !hasBG {
			bg = k.Surface
		}
		fg, bg, hasBG = bg, fg, true
	}
	if s.Style.Dim {
		fg = fg.Mix(k.Surface, 0.4)
	}
	out := ui.Span{Text: s.Text, Color: fg, Italic: s.Style.Italic, Underline: s.Style.Underline, Strikethrough: s.Style.Strike}
	if s.Style.Bold {
		out.Weight = 700
	}
	if hasBG {
		out.Background = bg
	}
	return out
}
