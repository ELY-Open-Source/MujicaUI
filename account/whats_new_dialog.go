package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/display"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/overlay"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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
func WhatsNewDialog(c *ui.Context, open *bool, hide *bool, releases []Release, opts WhatsNewDialogOptions) core.ChoiceResult {
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
	k := core.Tokens(c)
	d := core.Density(c)
	title := opts.Title
	if title == "" {
		title = core.Msg(c, "whats_new_dialog.title")
	}
	var r core.ChoiceResult
	was := *open
	before := *hide
	r.Element = overlay.Dialog(c, open, overlay.DialogOptions{Title: title, Width: 480, Decorated: true, Actions: func() {
		input.Checkbox(c, hide, core.Msg(c, "whats_new_dialog.hide"), input.CheckboxOptions{})
		ui.Box(c).Grow(1)
		if input.Button(c, core.Msg(c, "whats_new_dialog.close"), input.ButtonOptions{}).AutoFocus().Clicked() {
			*open = false
		}
	}}, func() {
		for i, rel := range releases {
			if rel.Version == "" {
				panic("mujica: Release needs a Version")
			}
			ui.Column(c).Key(i).Gap(d.Unit() * 2).Role(ui.RoleGroup).Label(rel.Version).Children(func() {
				ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.End).Children(func() {
					ui.Text(c, rel.Version).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600)
					if rel.Date != "" {
						ui.Text(c, rel.Date).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize))
					}
				})
				for _, ch := range rel.Changes {
					ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Start).Children(func() {
						tone := map[ChangeKind]display.BadgeTone{ChangeNew: display.BadgeAccent, ChangeImproved: display.BadgeInfo, ChangeFixed: display.BadgeSuccess}[ch.Kind]
						display.Badge(c, display.BadgeOptions{Text: core.Msg(c, fmt.Sprintf("whats_new_dialog.kind%d", ch.Kind)), Tone: tone}, nil).Shrink(0)
						ui.Text(c, ch.Text).Grow(1).MinWidth(0)
					})
				}
			})
		}
	})
	r = core.NewChoiceResult(r.Element, *hide != before)
	if was && !*open {
		slog.Debug("mujica: whats_new_dialog close", "hide", *hide)
	}
	return r
}
