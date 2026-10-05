package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// SuggestionChipsOptions configure SuggestionChips.
type SuggestionChipsOptions struct {
	// Send makes a click send the suggestion; by default it fills the input.
	Send bool
	// Label names the group, such as "Follow-up questions".
	Label    string
	Disabled bool
}

// SuggestionChipsView is built chips and the suggestion used this frame.
type SuggestionChipsView struct {
	Element *ui.Element
	// Chosen is the suggestion clicked, empty for none.
	Chosen string
	// Sent reports that Chosen should be sent now; otherwise it is already in the input.
	Sent bool
}

// SuggestionChips offers starter prompts or follow-ups; a click fills *input, or with Send asks to send.
func SuggestionChips(c *ui.Context, input *string, suggestions []string, opts SuggestionChipsOptions) SuggestionChipsView {
	if input == nil {
		panic("mujica: SuggestionChips needs an input")
	}
	k := Tokens(c)
	d := Density(c)
	row := ui.Row(c).Wrap().Gap(d.Unit() * 2).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(opts.Label).Disabled(opts.Disabled)
	v := SuggestionChipsView{Element: row}
	icon := icons.Must("chevron-right")
	if opts.Send {
		icon = icons.Must("arrow-up")
	}
	row.Children(func() {
		for i, s := range suggestions {
			if s == "" {
				panic("mujica: SuggestionChips with an empty suggestion")
			}
			b := ui.ButtonBase(c).Key(i).MinHeight(max(theme.MinHitSize, d.ControlHeight()-4)).Padding(0, d.Unit()*3).
				Gap(d.Unit()*2).Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).Label(s)
			internal.Ring(b, k.Focus, theme.ControlRadius)
			switch {
			case b.Pressed():
				b.Background(k.Selection)
			case b.Hovered():
				b.Background(k.SurfaceHover)
			}
			b.Transition(colorFade)
			b.Children(func() {
				ui.Text(c, s).SingleLine()
				ui.Icon(c, icon).Size(14, 14).TextColor(k.TextMuted)
			})
			if b.Clicked() && !b.IsDisabled() {
				v.Chosen = s
			}
		}
	})
	if v.Chosen != "" {
		v.Sent = opts.Send
		if !opts.Send {
			*input = v.Chosen
		}
		slog.Debug("mujica: suggestion_chips choose", "send", v.Sent)
	}
	return v
}
