package agent

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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
func ToolCallCard(c *ui.Context, call ToolCall, opts ToolCallCardOptions) core.ChoiceResult {
	if call.Name == "" {
		panic("mujica: ToolCallCard needs a Name")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = AgentCard(c).Role(ui.RoleGroup).Label(call.Name + ", " + agentWord(c, call.State))
	st := ui.Local(r.Element, "mujica.toolcall", func() toolCallCardState {
		return toolCallCardState{opts.OpenArgs, opts.OpenResult}
	})
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).MinHeight(theme.MinHitSize).Children(func() {
			agentMark(c, call.State, 16, opts.Still)
			ui.Text(c, call.Name).Font(theme.MonoFont).FontWeight(600).SingleLine().Grow(1).Shrink(1)
			AgentCaption(c, agentWord(c, call.State))
			if call.State == AgentDone || call.State == AgentFailed {
				AgentCaption(c, AgentDuration(call.Duration))
			}
		})
		if call.Args != "" {
			r = core.NewChoiceResult(r.Element, data.Collapsible(c, &st.args, data.CollapsibleOptions{Title: core.Msg(c, "tool_call_card.args")}, func() {
				AgentCode(c, "json", agentPrettyJSON(call.Args), 240)
			}).Changed() || r.Changed())
		}
		body, title := call.Result, core.Msg(c, "tool_call_card.result")
		if call.Error != "" {
			body, title = call.Error, core.Msg(c, "tool_call_card.error")
		}
		if body != "" {
			r = core.NewChoiceResult(r.Element, data.Collapsible(c, &st.result, data.CollapsibleOptions{Title: title}, func() {
				if call.Error != "" {
					ui.Text(c, call.Error).TextColor(k.Danger).Selectable()
					return
				}
				AgentCode(c, "json", agentPrettyJSON(body), 320)
			}).Changed() || r.Changed())
		}
	})
	if r.Changed() {
		slog.Debug("mujica: tool_call_card toggle", "tool", call.Name, "args", st.args, "result", st.result)
	}
	return r
}
