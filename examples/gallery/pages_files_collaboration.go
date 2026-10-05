package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/files"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

var filesPeople = []files.Collaborator{
	{ID: "a", Name: "Ada Lovelace"}, {ID: "b", Name: "Blaise Pascal", Away: true},
	{ID: "c", Name: "Carl Gauss"}, {ID: "d", Name: "Dora Maar"}, {ID: "e", Name: "Emmy Noether"},
}

func presenceAvatarsDemo(c *ui.Context, log func(string)) {
	type state struct {
		following       string
		disabled, empty bool
		max             float64
	}
	s := ui.Local(c.Root(), "392", func() state { return state{max: 3} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	mediaCheck(c, log, &s.empty, "Empty")
	if input.NumberInput(c, &s.max, input.NumberInputOptions{Min: 1, Max: 5, Step: 1, Label: "Visible avatars"}).Changed() {
		log(fmt.Sprintf("Visible: %.0f", s.max))
	}
	people := filesPeople
	if s.empty {
		people = nil
	}
	r := files.PresenceAvatars(c, people, files.PresenceAvatarsOptions{Following: s.following, Max: int(s.max), Disabled: s.disabled})
	if r.Follow != "" {
		s.following = r.Follow
		log("Follow: " + r.Follow)
	}
	if r.Unfollow {
		s.following = ""
		log("Following stopped")
	}
}

var filesReactionChoices = []string{"\U0001f44d", "\u2764\ufe0f", "\U0001f389", "\U0001f440", "\U0001f64f"}

func reactionsDemo(c *ui.Context, log func(string)) {
	type state struct {
		items    []files.Reaction
		disabled bool
	}
	s := ui.Local(c.Root(), "394", func() state {
		return state{items: []files.Reaction{{Emoji: filesReactionChoices[0], Count: 3, Mine: true}, {Emoji: filesReactionChoices[2], Count: 1}}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	if input.Button(c, "Clear reactions", input.ButtonOptions{}).Clicked() {
		s.items = nil
		log("Reactions cleared")
	}
	r := files.Reactions(c, s.items, files.ReactionsOptions{Choices: filesReactionChoices, Disabled: s.disabled})
	emoji := r.Toggled
	if emoji == "" {
		emoji = r.Added
	}
	if emoji == "" {
		return
	}
	for i := range s.items {
		if s.items[i].Emoji == emoji {
			if s.items[i].Mine {
				s.items[i].Mine = false
				s.items[i].Count--
			} else {
				s.items[i].Mine = true
				s.items[i].Count++
			}
			if s.items[i].Count == 0 {
				s.items = append(s.items[:i:i], s.items[i+1:]...)
			}
			log("Reaction toggled")
			return
		}
	}
	s.items = append(s.items, files.Reaction{Emoji: emoji, Count: 1, Mine: true})
	log("Reaction added")
}

func shareDialogDemo(c *ui.Context, log func(string)) {
	type state struct {
		members        []files.ShareMember
		open, disabled bool
		next           int
	}
	s := ui.Local(c.Root(), "395", func() state {
		return state{members: []files.ShareMember{
			{ID: "owner", Name: "Ada Lovelace", Email: "ada@court.org", Owner: true, Permission: files.PermissionEdit},
			{ID: "b", Name: "Blaise Pascal", Email: "blaise@court.org", Permission: files.PermissionView},
		}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	if input.Button(c, "Share Ledger", input.ButtonOptions{Disabled: s.disabled}).Clicked() {
		s.open = true
		log("Share dialog opened")
	}
	was := s.open
	r := files.ShareDialog(c, &s.open, &s.members, files.ShareDialogOptions{Title: "Ledger", Link: "https://court.org/ledger", Disabled: s.disabled, OnCopy: func(link string) { c.WriteClipboard(link); log("Copy callback") }})
	if r.Invited != nil {
		s.next++
		s.members = append(s.members, files.ShareMember{ID: fmt.Sprintf("invite-%d", s.next), Name: r.Invited.Email, Email: r.Invited.Email, Permission: r.Invited.Permission})
		log("Member invited")
	}
	if r.Changed != "" {
		log("Permission changed: " + r.Changed)
	}
	if r.Removed != "" {
		log("Member removed: " + r.Removed)
	}
	if r.Copied {
		log("Link copied")
	}
	if was && !s.open {
		log("Share dialog closed")
	}
}

func permissionSelectDemo(c *ui.Context, log func(string)) {
	type state struct {
		value    files.Permission
		disabled bool
	}
	s := ui.Local(c.Root(), "396", func() state { return state{} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if files.PermissionSelect(c, &s.value, files.PermissionSelectOptions{Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("Permission: %d", s.value))
	}
}

func liveIndicatorDemo(c *ui.Context, log func(string)) {
	type state struct {
		status   files.LiveStatus
		showTime bool
	}
	s := ui.Local(c.Root(), "399", func() state { return state{showTime: true} })
	if input.Select(c, &s.status, []input.SelectOption[files.LiveStatus]{{Value: files.LiveSynced, Label: "Synced"}, {Value: files.LiveSyncing, Label: "Syncing"}, {Value: files.LiveReconnecting, Label: "Reconnecting"}, {Value: files.LiveOffline, Label: "Offline"}, {Value: files.LiveError, Label: "Error"}}, input.SelectOptions{Label: "Sync status"}).Changed() {
		log(fmt.Sprintf("Sync status: %d", s.status))
	}
	mediaCheck(c, log, &s.showTime, "Last sync time")
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	opts := files.LiveIndicatorOptions{Error: "Server unavailable", Now: now}
	if s.showTime {
		opts.LastSynced = now.Add(-3 * time.Minute)
	}
	if files.LiveIndicator(c, s.status, opts).Changed() {
		s.status = files.LiveReconnecting
		log("Sync retried")
	}
}

func remoteCursorDemo(c *ui.Context, log func(string)) {
	type state struct {
		offset            float64
		selection, second bool
	}
	s := ui.Local(c.Root(), "400", func() state { return state{offset: 48, selection: true, second: true} })
	if input.NumberInput(c, &s.offset, input.NumberInputOptions{Min: 0, Max: 79, Step: 1, Label: "Remote text offset"}).Changed() {
		log(fmt.Sprintf("Offset: %.0f", s.offset))
	}
	mediaCheck(c, log, &s.selection, "Selection")
	mediaCheck(c, log, &s.second, "Second collaborator")
	carets := []files.RemoteCaret{{ID: "ada", Name: "Ada", Offset: int(s.offset), Anchor: int(s.offset)}}
	if s.selection {
		carets[0].Anchor = max(0, carets[0].Offset-5)
	}
	if s.second {
		carets = append(carets, files.RemoteCaret{ID: "blaise", Name: "Blaise", Offset: 8, Anchor: 8})
	}
	ui.Box(c).Height(140).FillWidth().Background(core.Tokens(c).Surface).Children(func() {
		ui.Column(c).Children(func() {
			for _, line := range []string{"The seal of the court is kept", "by the chancellor, who answers"} {
				ui.Text(c, line).Font(theme.MonoFont).FontSize(core.FontSize(c, 13)).Height(20)
			}
		})
		at := func(o int) (float32, float32) { return float32(o%40) * 8, float32(o/40) * 20 }
		files.RemoteCursor(c, carets, files.RemoteCursorOptions{
			Caret: func(o int) ui.Rect { x, y := at(o); return ui.Rect{X: x, Y: y, H: 20} },
			Selection: func(from, to int) []ui.Rect {
				var out []ui.Rect
				for o := from; o < to; o++ {
					x, y := at(o)
					out = append(out, ui.Rect{X: x, Y: y, W: 8, H: 20})
				}
				return out
			},
		})
	})
}
