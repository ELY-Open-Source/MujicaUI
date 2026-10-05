package mujica

import (
	"github.com/egoist/mygo/ui"
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
		label = msg(c, "agent_progress.label")
	}
	u := Density(c).Unit()
	col := ui.Column(c).FillWidth().Shrink(0).Gap(u * 1.5)
	col.Children(func() {
		Progress(c, value, ProgressOptions{Label: label, ShowValue: !opts.Indeterminate, Indeterminate: opts.Indeterminate})
		if opts.Action != "" {
			ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).Role(ui.RoleStatus).Label(opts.Action).Children(func() {
				agentMark(c, AgentRunning, 14, opts.Still)
				agentCaption(c, opts.Action)
			})
		}
	})
	return col
}
