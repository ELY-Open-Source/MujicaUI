package git

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Commit is one commit of a history, newest first.
type Commit struct {
	Hash    string
	Parents []string
	Author  string
	When    time.Time
	Subject string
	// Refs are branch and tag names pointing at the commit.
	Refs []string
}

// GitListResult is a built version-control list and its events.
type GitListResult struct {
	Element            *ui.Element
	changed, submitted bool
	action, item       string
}

// Changed reports a new selection.
func (r GitListResult) Changed() bool { return r.changed }

// Submitted reports Enter or a double click on a row.
func (r GitListResult) Submitted() bool { return r.submitted }

// Action returns this frame's action ID and its item.
func (r GitListResult) Action() (action, item string, ok bool) {
	return r.action, r.item, r.action != ""
}

// CommitListOptions configure a CommitList.
type CommitListOptions struct {
	// TimeFormat lays out commit times; empty means "2006-01-02 15:04".
	TimeFormat string
	Status     data.DataStatus
	Error      string
}

// CommitList shows commits to choose from; give it a height.
func CommitList(c *ui.Context, s *data.ListState[string], commits []Commit, opts CommitListOptions) GitListResult {
	layout := gitTimeLayout(opts.TimeFormat)
	v := data.List(c, s, len(commits), data.ListOptions[string]{
		Key:    func(i int) string { return gitCommitHash(commits[i]) },
		Label:  func(i int) string { return commits[i].Subject },
		Status: opts.Status, Error: opts.Error, Empty: core.Msg(c, "commit_list.empty"),
	}, func(i int) { commitListRow(c, commits[i], layout) })
	r := GitListResult{Element: v.Element, changed: v.Changed(), submitted: v.Submitted()}
	if r.changed || r.submitted {
		key, _ := s.Selected()
		slog.Debug("mujica: commit list choose", "hash", key, "open", r.submitted)
	}
	return r
}

func commitListRow(c *ui.Context, cm Commit, layout string) {
	k := core.Tokens(c)
	gitRefs(c, cm.Refs)
	ui.Text(c, cm.Subject).SingleLine().Grow(1).MinWidth(0)
	small := core.FontSize(c, theme.CaptionSize)
	ui.Text(c, cm.Author).FontSize(small).TextColor(k.TextMuted).SingleLine().MaxWidth(140).Shrink(0)
	ui.Text(c, gitShortHash(cm.Hash)).Font(theme.MonoFont).FontSize(small).TextColor(k.AccentText).Shrink(0)
	ui.Text(c, cm.When.Format(layout)).FontSize(small).FontFeatures("tnum").TextColor(k.TextMuted).Shrink(0)
}

// gitRefs shows branch and tag names as small outlined pills.
func gitRefs(c *ui.Context, refs []string) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	for _, ref := range refs {
		ui.Text(c, ref).FontSize(core.FontSize(c, theme.CaptionSize)).FontWeight(600).TextColor(k.AccentText).SingleLine().
			Padding(0, u*1.5).Radius(theme.ControlRadius).Border(theme.BorderWidth, k.AccentText).Shrink(0).MaxWidth(160)
	}
}

func gitCommitHash(cm Commit) string {
	if cm.Hash == "" {
		panic("mujica: Commit needs a Hash")
	}
	return cm.Hash
}

// gitShortHash is the first seven characters of a hash.
func gitShortHash(h string) string {
	if len(h) > 7 {
		return h[:7]
	}
	return h
}

func gitTimeLayout(f string) string {
	if f == "" {
		return "2006-01-02 15:04"
	}
	return f
}
