package agent

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// AgentStatusOptions configure an AgentStatus.
type AgentStatusOptions struct {
	// Label replaces the built-in word, as "Writing the report".
	Label string
	// Still stops the running animation.
	Still bool
}

// AgentStatus shows a state by icon shape and word.
func AgentStatus(c *ui.Context, state AgentState, opts AgentStatusOptions) *ui.Element {
	label := opts.Label
	if label == "" {
		label = agentWord(c, state) // panics on an unknown state
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	tone, _ := agentTone(k, state)
	pill := ui.Row(c).Shrink(0).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(u*1.5).Padding(u, u*2.5).
		Radius(theme.CardRadius).Border(theme.BorderWidth, tone).Background(k.Surface).
		Role(ui.RoleStatus).Label(label)
	pill.Children(func() {
		agentMark(c, state, 14, opts.Still)
		ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize+1)).TextColor(k.Text).SingleLine()
	})
	return pill
}
