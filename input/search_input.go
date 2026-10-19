package input

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// SearchInputOptions configure a SearchInput.
type SearchInputOptions struct {
	// Placeholder is empty for the built-in "Search" or "Find".
	Placeholder string
	// Label names the field; empty uses the placeholder.
	Label    string
	Disabled bool
	// Find turns on find-bar mode: match count, previous and next.
	Find bool
	// Matches is how many matches the caller found for the query.
	Matches int
	// Current is the chosen match index; Find mode requires it.
	Current *int
}

// SearchInputResult adds find-bar moves to InputResult.
type SearchInputResult struct {
	InputResult
	moved bool
}

// Moved reports that the user went to another match this frame.
func (r SearchInputResult) Moved() bool { return r.moved }

// SearchInput is a search field with clear, or a find bar with match navigation.
func SearchInput(c *ui.Context, query *string, opts SearchInputOptions) SearchInputResult {
	if query == nil {
		panic("mujica: SearchInput needs a query")
	}
	if opts.Find && opts.Current == nil {
		panic("mujica: SearchInput Find mode needs Current")
	}
	if opts.Matches < 0 {
		panic("mujica: SearchInput Matches is negative")
	}
	k := core.Tokens(c)
	place := opts.Placeholder
	if place == "" {
		place = core.Msg(c, "search_input.search")
		if opts.Find {
			place = core.Msg(c, "search_input.find")
		}
	}
	name := opts.Label
	if name == "" {
		name = place
	}
	var r SearchInputResult
	look := &textInputLook{disabled: opts.Disabled}
	step := func(d int) {
		if opts.Matches == 0 || look.disabled {
			return
		}
		*opts.Current = (*opts.Current + d + opts.Matches) % opts.Matches
		r.moved = true
		slog.Debug("mujica: search_input move", "current", *opts.Current)
	}
	if opts.Find {
		*opts.Current = max(0, min(*opts.Current, opts.Matches-1))
	}
	r.Element, r.Input = textInputFrame(c, look, func(frame ui.Element) ui.Element {
		if opts.Find {
			// Built inside the frame, where step knows of disabled ancestors.
			mac := runtime.GOOS == "darwin"
			switch {
			case mac && c.Shortcut(ui.Cmd, ui.KeyG), !mac && c.Shortcut(0, ui.KeyF3):
				step(1)
			case mac && c.Shortcut(ui.Shift|ui.Cmd, ui.KeyG), !mac && c.Shortcut(ui.Shift, ui.KeyF3):
				step(-1)
			}
		}
		ui.Icon(c, icons.Must("search")).Size(16, 16).Shrink(0).TextColor(k.TextMuted)
		back := ui.Local(frame, "back", func() bool { return false })
		in := TextInputEditor(c, frame, query, look.disabled, ui.TextInputBase).Grow(1).MinWidth(core.Density(c).Unit() * 10).Placeholder(place).Label(name)
		r.changed = in.Changed() && !look.disabled
		if opts.Find {
			// Shift+Enter would reach the editor as a plain submit.
			in.HandleInput(func(ev ui.InputEvent) bool {
				if ev.Kind == ui.InputKeyDown && ev.Key == ui.KeyEnter && ev.Mods == ui.Shift {
					*back = true
					return true
				}
				return false
			})
			switch {
			case *back:
				*back = false
				step(-1)
			case in.Submitted():
				step(1)
			}
		} else {
			r.submitted = in.Submitted() && !look.disabled
		}
		if *query != "" && !look.disabled && textInputClear(c, in) {
			*query = ""
			r.changed = true
			slog.Debug("mujica: search_input clear")
		}
		if opts.Find {
			searchInputFindControls(c, *query, opts.Matches, *opts.Current, step)
		}
		return in
	}, nil)
	return r
}

// searchInputFindControls shows the match count with previous and next buttons.
func searchInputFindControls(c *ui.Context, query string, matches, current int, step func(int)) {
	k := core.Tokens(c)
	status := core.Msg(c, "search_input.none")
	switch {
	case query == "":
		status = ""
	case matches > 0:
		status = fmt.Sprintf(core.Msg(c, "search_input.count"), current+1, matches)
	}
	ui.Text(c, status).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").
		SingleLine().Shrink(0).Role(ui.RoleStatus)
	ui.Box(c).Width(theme.BorderWidth).Height(16).Background(k.Border).Shrink(0)
	opts := IconButtonOptions{Variant: Ghost, Disabled: matches == 0}
	if IconButton(c, icons.Must("chevron-up"), core.Msg(c, "search_input.previous"), opts).Size(theme.MinHitSize, theme.MinHitSize).MinHeight(theme.MinHitSize).Clicked() {
		step(-1)
	}
	if IconButton(c, icons.Must("chevron-down"), core.Msg(c, "search_input.next"), opts).Size(theme.MinHitSize, theme.MinHitSize).MinHeight(theme.MinHitSize).Clicked() {
		step(1)
	}
}
