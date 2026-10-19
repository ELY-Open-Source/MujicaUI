package code

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/overlay"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Completion is one candidate of a CompletionMenu.
type Completion struct {
	Label  string
	Kind   SymbolKind
	Detail string
	// Doc describes the chosen candidate below the list.
	Doc string
}

// CompletionMenuState is the highlighted candidate and place of a CompletionMenu.
type CompletionMenuState struct {
	Selected int
	rows     ui.ListState
	pending  ui.Key
}

// CompletionMenuOptions configure a CompletionMenu.
type CompletionMenuOptions struct {
	Disabled bool
}

// CompletionMenuResult carries the accepted candidate.
type CompletionMenuResult struct {
	Element  ui.Element
	Accepted *Completion
}

// completionRows is how many candidates show before the list scrolls.
const completionRows = 8

// CompletionMenu shows candidates under the editor anchor.
func CompletionMenu(c *ui.Context, anchor ui.Element, open *bool, items []Completion, s *CompletionMenuState, opts CompletionMenuOptions) CompletionMenuResult {
	if !anchor.Valid() || open == nil || s == nil {
		panic("mujica: CompletionMenu needs an anchor, open and a state")
	}
	for _, it := range items {
		if it.Kind < SymbolFunction || it.Kind > SymbolModule {
			panic("mujica: Completion kind out of range")
		}
	}
	var res CompletionMenuResult
	if opts.Disabled || anchor.IsDisabled() || len(items) == 0 {
		*open = false
	}
	if !*open {
		return res
	}
	s.Selected = max(0, min(s.Selected, len(items)-1))
	accept := func(i int) {
		it := items[i]
		res.Accepted = &it
		*open = false
		slog.Debug("mujica: completion_menu accept", "index", i)
	}
	move := func(d int) {
		s.Selected = (s.Selected + d + len(items)) % len(items)
		s.rows.ScrollIntoView(s.Selected)
	}
	// Take the menu keys before the editor does.
	disabled := anchor.IsDisabled()
	anchor.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputKeyDown || ev.Mods != 0 || !*open || disabled {
			return false
		}
		switch ev.Key {
		case ui.KeyDown, ui.KeyUp, ui.KeyTab, ui.KeyEnter, ui.KeyEscape:
			s.pending = ev.Key
			return true
		}
		return false
	})
	key := s.pending
	s.pending = 0
	switch key {
	case ui.KeyDown:
		move(1)
	case ui.KeyUp:
		move(-1)
	case ui.KeyTab, ui.KeyEnter:
		accept(s.Selected)
		return res
	case ui.KeyEscape:
		*open = false
		slog.Debug("mujica: completion_menu close")
		return res
	}
	k := core.Tokens(c)
	d := core.Density(c)
	b := anchor.Bounds()
	core.OverlayLayer(c, open)
	res.Element = ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		core.OverlayStyle(c, panel, overlay.OverlayAnchored(c, panel, b))
		panel.Margin(4, 0, 0, 0).Padding(d.Unit()).Width(overlay.PopoverWide(c, 360)).MaxHeight(overlay.PopoverRoom(c, b)).
			Label(core.Msg(c, "completion_menu.label")).Role(ui.RoleList)
		chosen := s.Selected
		s.rows.Selected = &chosen
		s.rows.Label = func(i int) string { return items[i].Label }
		restore := data.DataSelectTheme(c, 0)
		list := ui.List(c, &s.rows, len(items), func(i int) {
			it := items[i]
			ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				CodeSelectMark(c, i == s.Selected)
				CodeBadge(c, symbolBadges[it.Kind], k.AccentText)
				CodeMono(c, it.Label).Grow(1).MinWidth(0)
				if it.Detail != "" {
					CodeCaption(c, it.Detail).MinWidth(0).Shrink(1)
				}
			})
		}).MaxHeight(d.ControlHeight() * completionRows).Shrink(1)
		restore()
		if list.Changed() && chosen >= 0 {
			accept(chosen)
		}
		if doc := items[s.Selected].Doc; doc != "" {
			ui.Text(c, doc).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Padding(d.Unit()*2).
				BorderWidth(theme.BorderWidth, 0, 0, 0).BorderColor(k.Border).Shrink(0)
		}
	})
	return res
}
