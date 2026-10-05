package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ChangeKind tags a release note.
type ChangeKind int

const (
	ChangeNew ChangeKind = iota
	ChangeImproved
	ChangeFixed
)

// ReleaseChange is one line of release notes.
type ReleaseChange struct {
	Kind ChangeKind
	Text string
}

// Release is a version's notes.
type Release struct {
	Version, Date string
	Changes       []ReleaseChange
}

// WhatsNewDialogOptions configure a WhatsNewDialog.
type WhatsNewDialogOptions struct {
	// Title replaces the built-in "What's new".
	Title string
}

// WhatsNewDialog displays release notes and a persistent hide choice.
func WhatsNewDialog(c *ui.Context, open *bool, hide *bool, releases []Release, opts WhatsNewDialogOptions) ChoiceResult {
	if len(releases) == 0 {
		panic("mujica: WhatsNewDialog needs releases")
	}
	for _, rel := range releases {
		if rel.Version == "" {
			panic("mujica: Release needs a Version")
		}
		for _, change := range rel.Changes {
			if change.Kind < ChangeNew || change.Kind > ChangeFixed || change.Text == "" {
				panic("mujica: ReleaseChange needs a valid Kind and Text")
			}
		}
	}
	if ui.Box(c).IsDisabled() {
		*open = false
	}
	k := Tokens(c)
	d := Density(c)
	title := opts.Title
	if title == "" {
		title = msg(c, "whats_new_dialog.title")
	}
	var r ChoiceResult
	was := *open
	before := *hide
	r.Element = Dialog(c, open, DialogOptions{Title: title, Width: 480, Decorated: true, Actions: func() {
		Checkbox(c, hide, msg(c, "whats_new_dialog.hide"), CheckboxOptions{})
		ui.Box(c).Grow(1)
		if Button(c, msg(c, "whats_new_dialog.close"), ButtonOptions{}).AutoFocus().Clicked() {
			*open = false
		}
	}}, func() {
		for i, rel := range releases {
			if rel.Version == "" {
				panic("mujica: Release needs a Version")
			}
			ui.Column(c).Key(i).Gap(d.Unit() * 2).Role(ui.RoleGroup).Label(rel.Version).Children(func() {
				ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.End).Children(func() {
					ui.Text(c, rel.Version).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).FontWeight(600)
					if rel.Date != "" {
						ui.Text(c, rel.Date).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize))
					}
				})
				for _, ch := range rel.Changes {
					ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Start).Children(func() {
						tone := map[ChangeKind]BadgeTone{ChangeNew: BadgeAccent, ChangeImproved: BadgeInfo, ChangeFixed: BadgeSuccess}[ch.Kind]
						Badge(c, BadgeOptions{Text: msg(c, fmt.Sprintf("whats_new_dialog.kind%d", ch.Kind)), Tone: tone}, nil).Shrink(0)
						ui.Text(c, ch.Text).Grow(1).MinWidth(0)
					})
				}
			})
		}
	})
	r.changed = *hide != before
	if was && !*open {
		slog.Debug("mujica: whats_new_dialog close", "hide", *hide)
	}
	return r
}
