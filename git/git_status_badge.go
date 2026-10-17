package git

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// GitStatus is a file's change state in version control.
type GitStatus int

const (
	GitModified GitStatus = iota
	GitAdded
	GitDeleted
	GitUntracked
	GitRenamed
	GitCopied
)

var gitStatusKeys = [...]string{"modified", "added", "deleted", "untracked", "renamed", "copied"}

// Letter is the git status code, as M.
func (s GitStatus) Letter() string { return string("MADURC"[gitStatusIndex(s)]) }

func gitStatusIndex(s GitStatus) int {
	if s < GitModified || s > GitCopied {
		panic("mujica: unknown GitStatus")
	}
	return int(s)
}

// GitStatusBadgeOptions configure a GitStatusBadge.
type GitStatusBadgeOptions struct {
	// ShowLabel writes the status name beside the letter.
	ShowLabel bool
}

// GitStatusBadge marks a status by letter, outline and color.
func GitStatusBadge(c *ui.Context, status GitStatus, opts GitStatusBadgeOptions) ui.Element {
	k := core.Tokens(c)
	name := core.Msg(c, "git_status_badge."+gitStatusKeys[gitStatusIndex(status)])
	col := gitStatusColor(k, status)
	row := ui.Row(c).AlignItems(ui.Center).Gap(core.Density(c).Unit() * 1.5).Shrink(0).Label(name).Role(ui.RoleImage)
	if !opts.ShowLabel {
		row.Tooltip(name)
	}
	row.Children(func() {
		ui.Box(c).Size(18, 18).Shrink(0).Center().Radius(theme.ControlRadius).Border(theme.BorderWidth, col).
			Children(func() {
				ui.Text(c, status.Letter()).Font(theme.MonoFont).FontWeight(700).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(col)
			})
		if opts.ShowLabel {
			ui.Text(c, name).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
		}
	})
	return row
}

func gitStatusColor(k theme.Tokens, s GitStatus) ui.Color {
	switch s {
	case GitAdded:
		return k.Success
	case GitDeleted:
		return k.Danger
	case GitUntracked, GitRenamed:
		return k.Info
	case GitCopied:
		return k.AccentText
	}
	return k.Warning
}
