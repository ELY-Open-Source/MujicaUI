package git

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Branch is a local or remote branch.
type Branch struct {
	// Name is unique; remote names carry the remote, as origin/main.
	Name   string
	Remote bool
	// Current marks the checked-out branch.
	Current bool
	// Upstream is the tracked remote branch, if any.
	Upstream      string
	Ahead, Behind int
}

// GitAction is a row action of a version-control list.
type GitAction struct {
	ID    string
	Label string
	// Icon names an icon from icons/.
	Icon string
}

// BranchListOptions configure a BranchList.
type BranchListOptions struct {
	// Actions show on every row; Action returns their IDs.
	Actions []GitAction
	Status  data.DataStatus
	Error   string
}

// BranchList shows branches with upstream and ahead/behind counts.
func BranchList(c *ui.Context, s *data.ListState[string], branches []Branch, opts BranchListOptions) GitListResult {
	gitActionsCheck(opts.Actions)
	var r GitListResult
	v := data.List(c, s, len(branches), data.ListOptions[string]{
		Key:    func(i int) string { return branchKey(branches[i]) },
		Label:  func(i int) string { return branches[i].Name },
		Status: opts.Status, Error: opts.Error, Empty: core.Msg(c, "branch_list.empty"),
	}, func(i int) {
		b := branches[i]
		branchListRow(c, b)
		gitRowActions(c, opts.Actions, b.Name, &r)
	})
	r.Element, r.changed, r.submitted = v.Element, v.Changed(), v.Submitted()
	return r
}

func branchKey(b Branch) string {
	if b.Name == "" {
		panic("mujica: Branch needs a Name")
	}
	if b.Ahead < 0 || b.Behind < 0 {
		panic("mujica: Branch Ahead and Behind must not be negative")
	}
	return b.Name
}

func branchListRow(c *ui.Context, b Branch) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	mark := ui.Box(c).Size(16, 16).Shrink(0).Center().Role(ui.RoleNone)
	if b.Current {
		mark.TextColor(k.AccentText).Children(func() { ui.Icon(c, icons.Must("check")).Size(14, 14) })
	}
	ui.Icon(c, icons.Must("git-branch")).Size(14, 14).Shrink(0).TextColor(k.TextMuted)
	name := ui.Text(c, b.Name).SingleLine().MinWidth(0).Shrink(1)
	if b.Current {
		name.FontWeight(600)
		ui.Text(c, core.Msg(c, "branch_list.current")).FontSize(small).TextColor(k.AccentText).Shrink(0)
	}
	ui.Spacer(c)
	if b.Upstream != "" {
		ui.Row(c).Gap(u).AlignItems(ui.Center).Shrink(0).TextColor(k.TextMuted).Children(func() {
			ui.Text(c, b.Upstream).FontSize(small).SingleLine().MaxWidth(160)
			if b.Ahead > 0 {
				ui.Text(c, fmt.Sprintf("↑%d", b.Ahead)).FontSize(small).FontFeatures("tnum").Label(fmt.Sprintf(core.Msg(c, "branch_list.ahead"), b.Ahead))
			}
			if b.Behind > 0 {
				ui.Text(c, fmt.Sprintf("↓%d", b.Behind)).FontSize(small).FontFeatures("tnum").Label(fmt.Sprintf(core.Msg(c, "branch_list.behind"), b.Behind))
			}
		})
	}
}

// gitActionsCheck panics on incomplete or duplicate actions.
func gitActionsCheck(actions []GitAction) {
	seen := map[string]bool{}
	for _, a := range actions {
		if a.ID == "" || a.Label == "" || a.Icon == "" {
			panic("mujica: GitAction needs an ID, Label and Icon")
		}
		if seen[a.ID] {
			panic("mujica: duplicate GitAction " + a.ID)
		}
		seen[a.ID] = true
		icons.Must(a.Icon)
	}
}

// gitRowActions shows row actions and records the one used.
func gitRowActions(c *ui.Context, actions []GitAction, item string, r *GitListResult) {
	for _, a := range actions {
		b := core.IconAction(c, icons.Must(a.Icon), a.Label).Tooltip(a.Label).Shrink(0)
		if b.Clicked() && !b.IsDisabled() {
			r.action, r.item = a.ID, item
			slog.Debug("mujica: git list action", "action", a.ID, "item", item)
		}
	}
}
