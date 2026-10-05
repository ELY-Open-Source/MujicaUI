package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// ThemeSelectorOptions configure a ThemeSelector.
type ThemeSelectorOptions struct {
	Disabled bool
}

// ThemeSelector selects an appearance using preview cards.
func ThemeSelector(c *ui.Context, mode *Mode, opts ThemeSelectorOptions) ChoiceResult {
	if *mode != System && *mode != Light && *mode != Dark {
		panic(fmt.Sprintf("mujica: ThemeSelector mode %d is not System, Light or Dark", *mode))
	}
	k := Tokens(c)
	d := Density(c)
	var r ChoiceResult
	r.Element = ui.Row(c).Gap(d.Unit() * 3).Wrap().Role(ui.RoleRadioGroup).Label(msg(c, "theme_selector.label")).
		FocusGroup(ui.Horizontal).Disabled(opts.Disabled)
	selected := mode
	if r.Element.IsDisabled() {
		copy := *mode
		selected = &copy
	}
	r.Element.Children(func() {
		for _, m := range []Mode{System, Light, Dark} {
			card := ui.RadioBase(c, selected, m).Column().Gap(d.Unit() * 2).Padding(d.Unit() * 2).Radius(theme.CardRadius).AlignItems(ui.Stretch).Transition(colorFade)
			if card.Changed() && !card.IsDisabled() {
				r.changed = true
				slog.Debug("mujica: theme_selector chose", "mode", m.String())
			}
			on := *mode == m
			internal.Ring(card, k.Focus, theme.CardRadius)
			border, width := k.Border, theme.BorderWidth
			switch {
			case on:
				border, width = k.AccentText, 2
				card.Background(k.Selection)
			case card.Hovered():
				card.Background(k.SurfaceHover)
			}
			card.Border(width, border)
			card.Children(func() {
				ui.Box(c).Size(132, 84).Role(ui.RoleNone).Draw(func(p *ui.Painter, rect ui.Rect) {
					switch m {
					case Light:
						themePreview(p, rect, theme.Light())
					case Dark:
						themePreview(p, rect, theme.Dark())
					default:
						half := ui.Rect{X: rect.X, Y: rect.Y, W: rect.W / 2, H: rect.H}
						p.Clip(half, 0, func() { themePreview(p, rect, theme.Light()) })
						p.Clip(ui.Rect{X: half.X + half.W, Y: rect.Y, W: rect.W / 2, H: rect.H}, 0, func() { themePreview(p, rect, theme.Dark()) })
					}
				})
				ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Children(func() {
					if on {
						ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(k.AccentText)
					}
					ui.Text(c, msg(c, "theme_selector."+m.String())).SingleLine()
				})
			})
		}
	})
	return r
}

// themePreview paints a miniature themed window.
func themePreview(p *ui.Painter, r ui.Rect, t theme.Tokens) {
	p.Fill(r, t.Background, theme.ControlRadius)
	p.Stroke(r, t.Border, theme.ControlRadius, 1)
	p.Fill(ui.Rect{X: r.X + 1, Y: r.Y + 1, W: r.W - 2, H: 12}, t.Surface, theme.ControlRadius)
	p.Fill(ui.Rect{X: r.X + 1, Y: r.Y + 13, W: r.W * 0.28, H: r.H - 14}, t.Surface, 0)
	x := r.X + r.W*0.36
	for i, w := range []float32{0.5, 0.38, 0.44} {
		p.Fill(ui.Rect{X: x, Y: r.Y + 22 + float32(i)*10, W: r.W * w, H: 4}, t.TextMuted, 2)
	}
	p.Fill(ui.Rect{X: x, Y: r.Y + r.H - 20, W: 30, H: 12}, t.Accent, 2)
}
