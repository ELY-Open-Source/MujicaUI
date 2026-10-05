package mujica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// AgentState is an agent or step state; zero is pending.
type AgentState int

const (
	AgentPending AgentState = iota
	AgentPlanning
	AgentRunning
	AgentWaiting
	AgentDone
	AgentFailed
	AgentSkipped
)

func (s AgentState) String() string {
	switch s {
	case AgentPending:
		return "pending"
	case AgentPlanning:
		return "planning"
	case AgentRunning:
		return "running"
	case AgentWaiting:
		return "waiting"
	case AgentDone:
		return "done"
	case AgentFailed:
		return "failed"
	case AgentSkipped:
		return "skipped"
	}
	panic(fmt.Sprintf("mujica: unknown AgentState %d", int(s)))
}

// agentTone returns the color and icon marking s.
func agentTone(k theme.Tokens, s AgentState) (ui.Color, *ui.SVG) {
	switch s {
	case AgentPlanning:
		return k.Info, icons.Must("list")
	case AgentRunning:
		return k.AccentText, icons.Must("loader-circle")
	case AgentWaiting:
		return k.Warning, icons.Must("clock")
	case AgentDone:
		return k.Success, icons.Must("circle-check")
	case AgentFailed:
		return k.Danger, icons.Must("circle-x")
	case AgentSkipped:
		return k.TextMuted, icons.Must("ban")
	}
	return k.TextMuted, icons.Must("circle")
}

// agentMark draws s; running spins unless still.
func agentMark(c *ui.Context, s AgentState, size float32, still bool) *ui.Element {
	tone, icon := agentTone(Tokens(c), s)
	if s == AgentRunning && !still {
		return Spinner(c, SpinnerOptions{Size: size, Color: tone, Label: msg(c, "agent_status.running")})
	}
	return ui.Icon(c, icon).Size(size, size).TextColor(tone).Shrink(0)
}

// agentWord is the built-in word for s.
func agentWord(c *ui.Context, s AgentState) string { return msg(c, "agent_status."+s.String()) }

// agentDuration formats d like "320 ms" or "4.2 s".
func agentDuration(d time.Duration) string {
	if d < 0 {
		panic(fmt.Sprintf("mujica: negative duration %v", d))
	}
	switch {
	case d < time.Second:
		return fmt.Sprintf("%d ms", d.Milliseconds())
	case d < time.Minute:
		return fmt.Sprintf("%.1f s", d.Seconds())
	case d < time.Hour:
		return fmt.Sprintf("%d m %02d s", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%d h %02d m", int(d.Hours()), int(d.Minutes())%60)
}

// agentCard is the bordered surface every agent card sits on.
func agentCard(c *ui.Context) *ui.Element {
	k := Tokens(c)
	u := Density(c).Unit()
	return ui.Column(c).Shrink(0).Padding(u*3).Gap(u*2).Radius(theme.CardRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.Border)
}

// agentCaption is small muted text.
func agentCaption(c *ui.Context, s string) *ui.Element {
	return ui.Text(c, s).FontSize(FontSize(c, theme.CaptionSize)).TextColor(Tokens(c).TextMuted)
}

// agentWell styles a sunken monospace well.
func agentWell(c *ui.Context, l *ui.Element, maxHeight float32) *ui.Element {
	k := Tokens(c)
	u := Density(c).Unit()
	return l.Shrink(1).MaxHeight(maxHeight).Padding(u*2, u*3).Radius(theme.ControlRadius).
		Background(k.Background).Border(theme.BorderWidth, k.Border).Font(theme.MonoFont).
		FontSize(FontSize(c, theme.CaptionSize+1))
}

// agentWellRows builds n rows in a well.
func agentWellRows(c *ui.Context, n int, row func(i int), maxHeight float32) *ui.Element {
	return agentWell(c, agentRows(c, n, row), maxHeight)
}

// agentRows builds n scrolling rows, virtualized when many.
func agentRows(c *ui.Context, n int, row func(i int)) *ui.Element {
	if n > agentVirtualRows {
		return ui.List(c, nil, n, row)
	}
	return ui.Scroll(c).Children(func() {
		for i := range n {
			row(i)
		}
	})
}

// agentVirtualRows is the row count above which wells virtualize.
const agentVirtualRows = 200

// agentCodeCache keeps lines and comment states per text.
type agentCodeCache struct {
	language, text string
	lines          []string
	open           []bool
}

func (cc *agentCodeCache) update(language, text string) {
	if cc.lines != nil && cc.language == language && cc.text == text {
		return
	}
	cc.language, cc.text = language, text
	cc.lines = internal.SplitLines(text)
	cc.open = make([]bool, len(cc.lines))
	block := false
	for i, line := range cc.lines {
		cc.open[i] = block
		_, block = internal.Highlight(language, line, block)
	}
}

// agentPrettyJSON indents valid JSON; anything else comes back unchanged.
func agentPrettyJSON(s string) string {
	var b bytes.Buffer
	if json.Indent(&b, []byte(s), "", "  ") != nil {
		return s
	}
	return b.String()
}

// agentTokenColor colors a highlight token from the theme.
func agentTokenColor(k theme.Tokens, kind internal.TokenKind) ui.Color {
	switch kind {
	case internal.TokKeyword:
		return k.AccentText
	case internal.TokString:
		return k.Success
	case internal.TokNumber, internal.TokType:
		return k.Info
	case internal.TokComment, internal.TokPunct:
		return k.TextMuted
	}
	return k.Text
}

// agentCodeLine builds one highlighted line of language.
func agentCodeLine(c *ui.Context, language, line string, inBlock bool) {
	k := Tokens(c)
	toks, _ := internal.Highlight(language, line, inBlock)
	spans := make([]ui.Span, 0, len(toks))
	for _, t := range toks {
		spans = append(spans, ui.Span{Text: t.Text, Color: agentTokenColor(k, t.Kind)})
	}
	if len(spans) == 0 {
		spans = append(spans, ui.Span{Text: " "})
	}
	ui.RichText(c, spans...)
}

// agentCode shows highlighted text in a capped well.
func agentCode(c *ui.Context, language, text string, maxHeight float32) *ui.Element {
	holder := ui.Column(c).Shrink(1)
	cc := ui.Local(holder, "mujica.agent.code", func() *agentCodeCache { return &agentCodeCache{} })
	(*cc).update(language, text)
	lines, open := (*cc).lines, (*cc).open
	row := func(i int) { agentCodeLine(c, language, lines[i], open[i]) }
	holder.Children(func() {
		agentWellRows(c, len(lines), row, maxHeight)
	})
	return holder
}

// agentAnsiColor maps a terminal color onto the theme.
func agentAnsiColor(k theme.Tokens, a internal.AnsiColor, fallback ui.Color) ui.Color {
	switch a.Kind {
	case internal.AnsiRGB:
		return ui.RGB(a.R, a.G, a.B)
	case internal.AnsiIndexed:
		switch a.Index % 8 {
		case 0:
			return k.TextMuted
		case 1:
			return k.Danger
		case 2:
			return k.Success
		case 3:
			return k.Warning
		case 4, 6:
			return k.Info
		case 5:
			return k.AccentText
		case 7:
			return k.Text
		}
	}
	return fallback
}

// agentAnsiSpans turns parsed terminal spans into rich text spans.
func agentAnsiSpans(k theme.Tokens, in []internal.AnsiSpan) []ui.Span {
	out := make([]ui.Span, 0, len(in)+1)
	for _, s := range in {
		fg := agentAnsiColor(k, s.Style.FG, k.Text)
		bg := agentAnsiColor(k, s.Style.BG, ui.Transparent)
		if s.Style.Inverse {
			fg, bg = agentAnsiColor(k, s.Style.BG, k.Background), fg
		}
		if s.Style.Dim {
			fg = fg.Alpha(0.7)
		}
		sp := ui.Span{Text: s.Text, Color: fg, Background: bg, Italic: s.Style.Italic,
			Underline: s.Style.Underline, Strikethrough: s.Style.Strike}
		if s.Style.Bold {
			sp.Weight = 700
		}
		out = append(out, sp)
	}
	if len(out) == 0 {
		out = append(out, ui.Span{Text: " "})
	}
	return out
}

// agentBar is a thin progress bar named label.
func agentBar(c *ui.Context, frac float64, label string) *ui.Element {
	if frac < 0 || frac > 1 || frac != frac {
		panic(fmt.Sprintf("mujica: progress %v outside 0..1", frac))
	}
	k := Tokens(c)
	return ui.Box(c).FillWidth().Height(4).Shrink(0).Role(ui.RoleProgress).Label(label).
		Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, k.Border, 2)
			if frac > 0 {
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * float32(frac), H: r.H}, k.AccentText, 2)
			}
		})
}

// agentToggle is a small text button reporting a click.
func agentToggle(c *ui.Context, label string, icon string, disabled bool) bool {
	k := Tokens(c)
	u := Density(c).Unit()
	b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).MinWidth(theme.MinHitSize).Padding(0, u*2).Gap(u).
		Radius(theme.ControlRadius).Disabled(disabled).TextColor(k.AccentText)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if !disabled && b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.Children(func() {
		if icon != "" {
			ui.Icon(c, icons.Must(icon)).Size(14, 14)
		}
		ui.Text(c, label).FontSize(FontSize(c, theme.CaptionSize+1)).SingleLine()
	})
	return !b.IsDisabled() && b.Clicked()
}
