package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// ChatMode is how the assistant answers.
type ChatMode int

const (
	ModeChat ChatMode = iota
	ModeSearch
	ModeResearch
	ModeAgent
)

// chatModes names each mode's copy key and icon, in order.
var chatModes = [...][2]string{
	{"mode_selector.chat", "message-square"},
	{"mode_selector.search", "globe"},
	{"mode_selector.research", "telescope"},
	{"mode_selector.agent", "bot"},
}

// ModeSelectorOptions configure a ModeSelector.
type ModeSelectorOptions struct {
	// IconsOnly hides the labels, which then name the icons.
	IconsOnly bool
	Disabled  bool
}

// ModeSelector switches conversation modes through clicks or arrow keys.
func ModeSelector(c *ui.Context, mode *ChatMode, opts ModeSelectorOptions) ChoiceResult {
	if mode == nil || *mode < ModeChat || *mode > ModeAgent {
		panic("mujica: unknown ChatMode")
	}
	k, u := Tokens(c), Density(c).Unit()
	old, pick := *mode, int(*mode)
	wrapper := ui.Row(c).Shrink(0).Disabled(opts.Disabled)
	var track *ui.Element
	wrapper.Children(func() {
		parts := ui.SegmentedBase(c, &pick, len(chatModes))
		track = parts.Track.Padding(2).Gap(2).Radius(theme.ControlRadius+2).Background(k.Background).
			Border(theme.BorderWidth, k.ControlBorder).MinHeight(Density(c).ControlHeight()).Label(msg(c, "mode_selector.label"))
		track.Children(func() {
			for i, m := range chatModes {
				label := msg(c, m[0])
				seg := parts.Segment(i).MinWidth(theme.MinHitSize).Padding(0, u*3).Gap(u * 1.5).Radius(theme.ControlRadius).Label(label)
				on := i == pick
				fg := k.TextMuted
				switch {
				case on:
					seg.Background(k.Selection).FontWeight(600).Draw(func(p *ui.Painter, r ui.Rect) {
						p.Fill(ui.Rect{X: r.X + 4, Y: r.Y + r.H - 3, W: r.W - 8, H: 2}, k.AccentText, 0)
					})
					fg = k.AccentText
				case seg.Pressed():
					seg.Background(k.Selection)
					fg = k.Text
				case seg.Hovered():
					seg.Background(k.SurfaceHover)
					fg = k.Text
				}
				internal.Ease(seg, Motion(c, theme.StateDuration), fg)
				internal.Ring(seg, k.Focus, theme.ControlRadius)
				if opts.IconsOnly {
					seg.Tooltip(label)
				}
				seg.Children(func() {
					ui.Icon(c, icons.Must(m[1])).Size(16, 16)
					if !opts.IconsOnly {
						ui.Text(c, label).SingleLine()
					}
				})
			}
		})
	})
	changed := ChatMode(pick) != old && !wrapper.IsDisabled()
	if changed {
		*mode = ChatMode(pick)
		slog.Debug("mujica: mode_selector changed", "mode", pick)
	}
	return ChoiceResult{Element: track, changed: changed}
}
