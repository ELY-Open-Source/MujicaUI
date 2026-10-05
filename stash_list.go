package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// Stash is one stash entry.
type Stash struct {
	// Ref names the entry, as stash@{0}.
	Ref     string
	Message string
	Branch  string
}

// StashListOptions configure a StashList.
type StashListOptions struct {
	// Status and Error show the load state.
	Status DataStatus
	Error  string
}

// StashList shows stashes with apply, pop and drop actions.
func StashList(c *ui.Context, s *ListState[string], stashes []Stash, opts StashListOptions) GitListResult {
	var r GitListResult
	k := Tokens(c)
	actions := []GitAction{
		{ID: "apply", Label: msg(c, "stash_list.apply"), Icon: "check"},
		{ID: "pop", Label: msg(c, "stash_list.pop"), Icon: "upload"},
		{ID: "drop", Label: msg(c, "stash_list.drop"), Icon: "trash-2"},
	}
	v := List(c, s, len(stashes), ListOptions[string]{
		Key:    func(i int) string { return stashListKey(stashes[i]) },
		Label:  func(i int) string { return stashes[i].Message },
		Status: opts.Status, Error: opts.Error, Empty: msg(c, "stash_list.empty"),
	}, func(i int) {
		st := stashes[i]
		small := FontSize(c, theme.CaptionSize)
		ui.Icon(c, icons.Must("archive")).Size(14, 14).Shrink(0).TextColor(k.TextMuted)
		ui.Text(c, st.Ref).Font(theme.MonoFont).FontSize(small).TextColor(k.AccentText).Shrink(0)
		ui.Text(c, st.Message).SingleLine().Grow(1).MinWidth(0)
		if st.Branch != "" {
			ui.Text(c, st.Branch).FontSize(small).TextColor(k.TextMuted).SingleLine().MaxWidth(140).Shrink(0)
		}
		gitRowActions(c, actions, st.Ref, &r)
	})
	r.Element, r.changed, r.submitted = v.Element, v.Changed(), v.Submitted()
	return r
}

func stashListKey(s Stash) string {
	if s.Ref == "" {
		panic("mujica: Stash needs a Ref")
	}
	return s.Ref
}
