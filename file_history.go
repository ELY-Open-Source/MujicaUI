package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// FileHistoryOptions configure a FileHistory.
type FileHistoryOptions struct {
	// TimeFormat lays out commit times; empty means "2006-01-02 15:04".
	TimeFormat string
	Status     DataStatus
	Error      string
}

// FileHistory lists commits that changed path; Action "diff" opens one.
func FileHistory(c *ui.Context, s *ListState[string], path string, commits []Commit, opts FileHistoryOptions) GitListResult {
	if path == "" {
		panic("mujica: FileHistory needs a path")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	layout := gitTimeLayout(opts.TimeFormat)
	var r GitListResult
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(u * 2)
	col.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(u * 2).Children(func() {
			ui.Icon(c, icons.Must("history")).Size(16, 16).TextColor(k.Ornament.Mix(k.Text, 0.35))
			ui.Text(c, path).Font(theme.MonoFont).FontWeight(600).SingleLine().MinWidth(0)
			ui.Text(c, msg(c, "file_history.title")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Shrink(0)
		})
		actions := []GitAction{{ID: "diff", Label: msg(c, "file_history.diff"), Icon: "arrow-left-right"}}
		v := List(c, s, len(commits), ListOptions[string]{
			Key:    func(i int) string { return gitCommitHash(commits[i]) },
			Label:  func(i int) string { return commits[i].Subject },
			Status: opts.Status, Error: opts.Error, Empty: msg(c, "file_history.empty"),
		}, func(i int) {
			commitListRow(c, commits[i], layout)
			gitRowActions(c, actions, commits[i].Hash, &r)
		})
		v.Element.Grow(1)
		r.changed, r.submitted = v.Changed(), v.Submitted()
	})
	r.Element = col
	return r
}
