package agent

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/overlay"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ToolApprovalOptions configure a ToolApprovalDialog. Tool is required.
type ToolApprovalOptions struct {
	Tool string
	// Args is the call's arguments, JSON when valid.
	Args string
	// Reason says why the agent wants to run the tool.
	Reason string
}

// ToolApprovalResult reports the choice made in this frame.
type ToolApprovalResult struct {
	Element *ui.Element
	Allowed bool
	Denied  bool
}

// ToolApprovalDialog asks to approve a tool call; Deny has focus.
func ToolApprovalDialog(c *ui.Context, open *bool, opts ToolApprovalOptions) ToolApprovalResult {
	if opts.Tool == "" {
		panic("mujica: ToolApprovalDialog needs a Tool")
	}
	var r ToolApprovalResult
	was := *open
	u := core.Density(c).Unit()
	r.Element = overlay.Dialog(c, open, overlay.DialogOptions{
		Title: core.Msg(c, "tool_approval_dialog.title"), Description: opts.Reason, KeepOnBackdrop: true, HideClose: true,
		Actions: func() {
			// Deny first: a reflexive Enter never runs the tool.
			if input.Button(c, core.Msg(c, "tool_approval_dialog.deny"), input.ButtonOptions{Variant: input.Secondary}).AutoFocus().Clicked() {
				r.Denied = true
			}
			if input.Button(c, core.Msg(c, "tool_approval_dialog.allow"), input.ButtonOptions{}).Clicked() {
				r.Allowed = true
			}
		},
	}, func() {
		ui.Column(c).Gap(u * 2).Children(func() {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
				AgentCaption(c, core.Msg(c, "tool_approval_dialog.tool"))
				ui.Text(c, opts.Tool).Font(theme.MonoFont).FontWeight(600)
			})
			if opts.Args != "" {
				AgentCode(c, "json", agentPrettyJSON(opts.Args), 260)
			}
		})
	})
	if was && !*open && !r.Allowed {
		r.Denied = true // Escape
	}
	if r.Allowed || r.Denied {
		slog.Debug("mujica: tool_approval_dialog choose", "tool", opts.Tool, "allowed", r.Allowed)
		*open = false
	}
	return r
}
