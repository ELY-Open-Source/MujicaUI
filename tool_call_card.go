package mujica

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ToolCall is one tool invocation; the caller owns it.
type ToolCall struct {
	ID   string
	Name string
	// Args and Result show as JSON when valid.
	Args   string
	Result string
	// Error replaces Result when the call failed.
	Error    string
	State    AgentState
	Duration time.Duration
}

// ToolCallCardOptions configure a ToolCallCard.
type ToolCallCardOptions struct {
	// OpenArgs and OpenResult start the sections expanded.
	OpenArgs, OpenResult bool
	// Still stops the running animation.
	Still bool
}

type toolCallCardState struct{ args, result bool }

// ToolCallCard shows a tool call with foldable arguments and result.
func ToolCallCard(c *ui.Context, call ToolCall, opts ToolCallCardOptions) ChoiceResult {
	if call.Name == "" {
		panic("mujica: ToolCallCard needs a Name")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r ChoiceResult
	r.Element = agentCard(c).Role(ui.RoleGroup).Label(call.Name + ", " + agentWord(c, call.State))
	st := ui.Local(r.Element, "mujica.toolcall", func() toolCallCardState {
		return toolCallCardState{opts.OpenArgs, opts.OpenResult}
	})
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).MinHeight(theme.MinHitSize).Children(func() {
			agentMark(c, call.State, 16, opts.Still)
			ui.Text(c, call.Name).Font(theme.MonoFont).FontWeight(600).SingleLine().Grow(1).Shrink(1)
			agentCaption(c, agentWord(c, call.State))
			if call.State == AgentDone || call.State == AgentFailed {
				agentCaption(c, agentDuration(call.Duration))
			}
		})
		if call.Args != "" {
			r.changed = Collapsible(c, &st.args, CollapsibleOptions{Title: msg(c, "tool_call_card.args")}, func() {
				agentCode(c, "json", agentPrettyJSON(call.Args), 240)
			}).Changed() || r.changed
		}
		body, title := call.Result, msg(c, "tool_call_card.result")
		if call.Error != "" {
			body, title = call.Error, msg(c, "tool_call_card.error")
		}
		if body != "" {
			r.changed = Collapsible(c, &st.result, CollapsibleOptions{Title: title}, func() {
				if call.Error != "" {
					ui.Text(c, call.Error).TextColor(k.Danger).Selectable()
					return
				}
				agentCode(c, "json", agentPrettyJSON(body), 320)
			}).Changed() || r.changed
		}
	})
	if r.changed {
		slog.Debug("mujica: tool_call_card toggle", "tool", call.Name, "args", st.args, "result", st.result)
	}
	return r
}
