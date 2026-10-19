package media

import (
	"log/slog"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// PlaylistItem is one track; ID must be unique and stable.
type PlaylistItem struct {
	ID, Title, Subtitle string
	Duration            time.Duration
}

// PlaylistOptions configure a Playlist.
type PlaylistOptions struct {
	// Current is the playing item's ID; empty is none.
	Current  string
	Disabled bool
	// Height caps the list; zero means 320.
	Height float32
}

// PlaylistResult is the playlist's element and this frame's events.
type PlaylistResult struct {
	Element ui.Element
	// Played is the ID chosen by Enter or double click.
	Played string
	// Removed is the ID taken out of *items.
	Removed string
	// Moved reports a reorder of *items.
	Moved bool
}

type playlistKey struct{}

// Playlist lists *items, building only rows in view.
func Playlist(c *ui.Context, items *[]PlaylistItem, opts PlaylistOptions) PlaylistResult {
	if opts.Height < 0 {
		panic("mujica: Playlist Height must not be negative")
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	seen := make(map[string]bool, len(*items))
	for _, it := range *items {
		if it.ID == "" || seen[it.ID] {
			panic("mujica: Playlist item IDs must be unique and non-empty: " + it.ID)
		}
		seen[it.ID] = true
	}
	k := core.Tokens(c)
	var r PlaylistResult
	box := ui.Column(c).Height(h).Disabled(opts.Disabled)
	r.Element = box
	st := ui.Local(box, playlistKey{}, func() *data.ListState[string] { return &data.ListState[string]{} })
	s := *st
	list := *items
	live := !box.IsDisabled()
	s.Rows.Reorder = nil
	if live {
		s.Rows.Reorder = func(rows []int, to int) {
			*items = playlistMove(*items, rows, to)
			r.Moved = true
		}
	}
	removeAt := -1
	box.Children(func() {
		v := data.List(c, s, len(list), data.ListOptions[string]{
			Key:   func(i int) string { return list[i].ID },
			Label: func(i int) string { return list[i].Title },
			Empty: core.Msg(c, "playlist.empty"),
		}, func(i int) {
			playlistRow(c, k, list[i], i, list[i].ID == opts.Current)
			if core.IconAction(c, icons.Must("x"), core.Msg(c, "playlist.remove")+" "+list[i].Title).Clicked() && live {
				removeAt = i
			}
		})
		v.Element.Grow(1)
		if removeAt >= 0 {
			r.Removed = list[removeAt].ID
			*items = slices.Delete(*items, removeAt, removeAt+1)
			return
		}
		sel, ok := s.Selected()
		i := slices.IndexFunc(*items, func(it PlaylistItem) bool { return it.ID == sel })
		if !live || !ok || i < 0 {
			return
		}
		switch {
		case v.Submitted():
			r.Played = sel
		case v.Element.Shortcut(0, ui.KeyDelete), v.Element.Shortcut(0, ui.KeyBackspace):
			r.Removed = sel
			*items = slices.Delete(*items, i, i+1)
		case v.Element.Shortcut(ui.Alt, ui.KeyUp) && i > 0:
			*items, r.Moved = playlistMove(*items, []int{i}, i-1), true
		case v.Element.Shortcut(ui.Alt, ui.KeyDown) && i < len(*items)-1:
			*items, r.Moved = playlistMove(*items, []int{i}, i+2), true
		}
	})
	switch {
	case r.Played != "":
		slog.Debug("mujica: playlist play", "id", r.Played)
	case r.Removed != "":
		slog.Debug("mujica: playlist remove", "id", r.Removed)
	case r.Moved:
		slog.Debug("mujica: playlist move")
	}
	return r
}

// playlistRow shows a track's mark, title, subtitle and length.
func playlistRow(c *ui.Context, k theme.Tokens, it PlaylistItem, i int, playing bool) {
	ui.Icon(c, icons.Must("grip-vertical")).Size(14, 14).TextColor(k.TextMuted).Shrink(0)
	ui.Box(c).Width(20).Center().Shrink(0).Children(func() {
		if playing {
			ui.Icon(c, core.MediaPlay).Size(14, 14).TextColor(k.AccentText).Label(core.Msg(c, "playlist.playing"))
		} else {
			ui.Textf(c, "%d", i+1).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).FontFeatures("tnum")
		}
	})
	ui.Column(c).Grow(1).MinWidth(0).Children(func() {
		t := ui.Text(c, it.Title).SingleLine()
		if playing {
			t.TextColor(k.AccentText).FontWeight(600)
		}
		if it.Subtitle != "" {
			ui.Text(c, it.Subtitle).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
		}
	})
	if it.Duration > 0 {
		ui.Text(c, core.MediaClock(it.Duration)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).FontFeatures("tnum")
	}
}

// playlistMove puts rows before index to; len means the end.
func playlistMove(items []PlaylistItem, rows []int, to int) []PlaylistItem {
	moved := make([]PlaylistItem, 0, len(rows))
	rest := make([]PlaylistItem, 0, len(items))
	at := to
	for i, it := range items {
		if slices.Contains(rows, i) {
			moved = append(moved, it)
			if i < to {
				at--
			}
			continue
		}
		rest = append(rest, it)
	}
	return slices.Concat(rest[:at], moved, rest[at:])
}
