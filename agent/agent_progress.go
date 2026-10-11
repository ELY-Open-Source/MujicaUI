package agent

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/feedback"
)

// AgentProgressOptions configure an AgentProgress.
type AgentProgressOptions struct {
	// Label names the work; empty uses "Progress".
	Label string
	// Action is what the agent does now.
	Action string
	// Indeterminate shows work of unknown length; the value is ignored.
	Indeterminate bool
	Still         bool
}

// AgentProgress shows progress 0 to 1 and the current action.
func AgentProgress(c *ui.Context, value float64, opts AgentProgressOptions) *ui.Element {
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "agent_progress.label")
	}
	u := core.Density(c).Unit()
	col := ui.Column(c).FillWidth().Shrink(0).Gap(u * 1.5)
	col.Children(func() {
		feedback.Progress(c, value, feedback.ProgressOptions{Label: label, ShowValue: !opts.Indeterminate, Indeterminate: opts.Indeterminate})
		if opts.Action != "" {
			ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).Role(ui.RoleStatus).Label(opts.Action).Children(func() {
				agentMark(c, AgentRunning, 14, opts.Still)
				AgentCaption(c, opts.Action)
			})
		}
	})
	return col
}
