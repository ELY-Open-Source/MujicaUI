package agent

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// FileChange is a proposed edit; empty Old adds a file.
type FileChange struct {
	Path     string
	Old, New string
}

// FileDecision is a review state; zero is FilePending.
type FileDecision int

const (
	FilePending FileDecision = iota
	FileAccepted
	FileRejected
)

// agentDiffCache holds a change's line diff until its text changes.
type agentDiffCache struct {
	old, new  string
	lines     []internal.DiffLine
	add, del  int
	firstEdit int
	done      bool
}

func (dc *agentDiffCache) update(ch FileChange) {
	if dc.done && dc.old == ch.Old && dc.new == ch.New {
		return
	}
	dc.old, dc.new, dc.done = ch.Old, ch.New, true
	dc.lines = internal.LineDiff(internal.SplitLines(ch.Old), internal.SplitLines(ch.New))
	dc.add, dc.del = internal.DiffStats(dc.lines)
	dc.firstEdit = 0
	for i, l := range dc.lines {
		if l.Op != internal.DiffEqual {
			dc.firstEdit = i
			break
		}
	}
}

// agentDiffOf returns the cached diff of ch kept on holder.
func agentDiffOf(holder ui.Element, ch FileChange) *agentDiffCache {
	dc := ui.Local(holder, "mujica.agent.diff", func() *agentDiffCache { return &agentDiffCache{} })
	(*dc).update(ch)
	return *dc
}

// agentDiffStats shows signed added and removed counts.
func agentDiffStats(c *ui.Context, add, del int) {
	k := core.Tokens(c)
	size := core.FontSize(c, theme.CaptionSize+1)
	ui.Row(c).Gap(core.Density(c).Unit() * 1.5).Shrink(0).Label(fmt.Sprintf(core.Msg(c, "file_change_card.stats"), add, del)).Children(func() {
		ui.Text(c, fmt.Sprintf("+%d", add)).TextColor(k.Success).FontSize(size).Font(theme.MonoFont)
		ui.Text(c, fmt.Sprintf("−%d", del)).TextColor(k.Danger).FontSize(size).Font(theme.MonoFont)
	})
}

// agentDiffRow draws one numbered, signed, tinted diff line.
func agentDiffRow(c *ui.Context, l internal.DiffLine) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	sign, bg, fg := " ", ui.Transparent, k.TextMuted
	switch l.Op {
	case internal.DiffInsert:
		sign, bg, fg = "+", k.Success.Alpha(0.14), k.Success
	case internal.DiffDelete:
		sign, bg, fg = "−", k.Danger.Alpha(0.14), k.Danger
	}
	num := func(n int) string {
		if n == 0 {
			return ""
		}
		return fmt.Sprint(n)
	}
	ui.Row(c).Gap(u*2).Background(bg).Padding(0, u).Children(func() {
		ui.Text(c, num(l.Old)).Width(36).TextAlign(ui.End).TextColor(k.TextMuted).Shrink(0)
		ui.Text(c, num(l.New)).Width(36).TextAlign(ui.End).TextColor(k.TextMuted).Shrink(0)
		ui.Text(c, sign).TextColor(fg).Width(10).Shrink(0)
		text := l.Text
		if text == "" {
			text = " "
		}
		ui.Text(c, text).Grow(1).Shrink(1)
	})
}

// agentDiffView shows lines[from:to] in a capped well, virtualized when long.
func agentDiffView(c *ui.Context, lines []internal.DiffLine, maxHeight float32) ui.Element {
	w := AgentWellRows(c, len(lines), func(i int) { agentDiffRow(c, lines[i]) }, maxHeight)
	return w.Padding(core.Density(c).Unit(), 0)
}

// agentDecisionMark shows a decision by icon and word.
func agentDecisionMark(c *ui.Context, d FileDecision) {
	ui.Row(c).Gap(core.Density(c).Unit()).AlignItems(ui.Center).Shrink(0).Children(func() {
		AgentCaption(c, agentDecisionIcon(c, d))
	})
}

// agentDecisionIcon draws a named decision icon, returning its word.
func agentDecisionIcon(c *ui.Context, d FileDecision) string {
	k := core.Tokens(c)
	tone, icon, word := k.TextMuted, "circle", core.Msg(c, "file_change_card.pending")
	switch d {
	case FileAccepted:
		tone, icon, word = k.Success, "circle-check", core.Msg(c, "file_change_card.accepted")
	case FileRejected:
		tone, icon, word = k.Danger, "circle-x", core.Msg(c, "file_change_card.rejected")
	case FilePending:
	default:
		panic(fmt.Sprintf("mujica: unknown FileDecision %d", int(d)))
	}
	ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(tone).Shrink(0).Label(word)
	return word
}
