package mujica

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ToolCallGroupOptions configure a ToolCallGroup.
type ToolCallGroupOptions struct {
	// Title replaces the built-in "N tool calls".
	Title string
	Still bool
}

// ToolCallGroup folds tool calls under a count, time, failures summary.
func ToolCallGroup(c *ui.Context, open *bool, calls []ToolCall, opts ToolCallGroupOptions) ChoiceResult {
	if open == nil {
		panic("mujica: ToolCallGroup needs open")
	}
	seen := make(map[string]bool, len(calls))
	var total time.Duration
	failed, running := 0, false
	for _, call := range calls {
		if call.ID == "" || seen[call.ID] {
			panic(fmt.Sprintf("mujica: ToolCallGroup call IDs must be unique and set, got %q", call.ID))
		}
		seen[call.ID] = true
		total += call.Duration
		failed += agentB2I(call.State == AgentFailed)
		running = running || call.State == AgentRunning
	}
	k := Tokens(c)
	u := Density(c).Unit()
	title := opts.Title
	if title == "" {
		title = fmt.Sprintf(msg(c, "tool_call_group.count"), len(calls))
	}
	state := AgentDone
	switch {
	case running:
		state = AgentRunning
	case failed > 0:
		state = AgentFailed
	}
	summary := title + " · " + agentDuration(total)
	if failed > 0 {
		summary += " · " + fmt.Sprintf(msg(c, "tool_call_group.failed"), failed)
	}
	return Collapsible(c, open, CollapsibleOptions{Title: summary, Trigger: func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			agentMark(c, state, 16, opts.Still)
			ui.Text(c, title).FontWeight(600).SingleLine()
			agentCaption(c, agentDuration(total))
			if failed > 0 {
				ui.Text(c, fmt.Sprintf(msg(c, "tool_call_group.failed"), failed)).TextColor(k.Danger).
					FontSize(FontSize(c, theme.CaptionSize+1)).SingleLine()
			}
		})
	}}, func() {
		for _, call := range calls {
			ui.Column(c).Key(call.ID).Children(func() { ToolCallCard(c, call, ToolCallCardOptions{Still: opts.Still}) })
		}
	})
}

func agentB2I(b bool) int {
	if b {
		return 1
	}
	return 0
}
