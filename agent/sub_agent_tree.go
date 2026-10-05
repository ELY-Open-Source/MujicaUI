package agent

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SubAgent is one agent in a hierarchy.
type SubAgent struct {
	ID       string
	Name     string
	Task     string
	State    AgentState
	Children []SubAgent
}

// SubAgentTreeOptions configure a SubAgentTree.
type SubAgentTreeOptions struct {
	Disabled bool
	Still    bool
}

// SubAgentTree shows nested agents and their states.
func SubAgentTree(c *ui.Context, roots []SubAgent, opts SubAgentTreeOptions) core.ChoiceResult {
	seen := map[string]bool{}
	var check func([]SubAgent)
	check = func(list []SubAgent) {
		for _, a := range list {
			if a.ID == "" || a.Name == "" || seen[a.ID] {
				panic(fmt.Sprintf("mujica: SubAgentTree agent %q needs a unique ID and a Name", a.ID))
			}
			seen[a.ID] = true
			check(a.Children)
		}
	}
	check(roots)
	var r core.ChoiceResult
	r.Element = ui.Column(c).Shrink(0).FocusGroup(ui.Vertical).Disabled(opts.Disabled)
	open := ui.Local(r.Element, "mujica.subagents.open", func() map[string]bool { return map[string]bool{} })
	r.Element.Children(func() {
		var rows func([]SubAgent, int)
		rows = func(list []SubAgent, depth int) {
			for _, a := range list {
				if subAgentRow(c, a, depth, (*open)[a.ID], opts.Still) {
					(*open)[a.ID] = !(*open)[a.ID]
					r = core.NewChoiceResult(r.Element, true)
					slog.Debug("mujica: sub_agent_tree toggle", "id", a.ID, "open", (*open)[a.ID])
				}
				if (*open)[a.ID] {
					rows(a.Children, depth+1)
				}
			}
		}
		rows(roots, 0)
	})
	for id := range *open {
		if !seen[id] {
			delete(*open, id)
		}
	}
	return r
}

// subAgentRow builds one row, reporting an expand or collapse.
func subAgentRow(c *ui.Context, a SubAgent, depth int, open, still bool) bool {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	parent := len(a.Children) > 0
	row := ui.ButtonBase(c).Key(a.ID).Justify(ui.Start).Gap(u*2).MinHeight(d.ControlHeight()).
		Padding(0, u*2, 0, u*2+float32(depth)*u*5).Radius(theme.ControlRadius).
		Label(a.Name + ", " + agentWord(c, a.State))
	internal.Ring(row, k.Focus, theme.ControlRadius)
	if row.Hovered() {
		row.Background(k.SurfaceHover)
	}
	row.Children(func() {
		chev := ui.Box(c).Size(14, 14).Shrink(0)
		if parent {
			name := "chevron-right"
			if open {
				name = "chevron-down"
			}
			chev.Children(func() { ui.Icon(c, icons.Must(name)).Size(14, 14).TextColor(k.TextMuted) })
		}
		agentMark(c, a.State, 16, still)
		ui.Text(c, a.Name).FontWeight(600).SingleLine().Shrink(0)
		ui.Text(c, a.Task).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize+1)).SingleLine().Grow(1).Shrink(1)
		AgentCaption(c, agentWord(c, a.State))
		if parent {
			AgentCaption(c, fmt.Sprintf(core.Msg(c, "sub_agent_tree.children"), len(a.Children)))
		}
	})
	if !parent || row.IsDisabled() {
		return false
	}
	return row.Clicked() || !open && row.Shortcut(0, ui.KeyRight) || open && row.Shortcut(0, ui.KeyLeft)
}
