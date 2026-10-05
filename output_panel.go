package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// OutputChannel is a named stream of ANSI lines.
type OutputChannel struct {
	Name  string
	Lines []string
}

// OutputPanelState is the channel, following and place.
type OutputPanelState struct {
	Channel int
	// Follow keeps the newest line in view; it starts on.
	Follow  bool
	rows    ui.ListState
	shown   int
	started bool
}

// OutputPanelOptions configure an OutputPanel.
type OutputPanelOptions struct {
	Disabled bool
}

// OutputPanelResult is an OutputPanel and what happened in this frame.
type OutputPanelResult struct {
	Element *ui.Element
	// Cleared names a channel to clear.
	Changed bool
	Cleared string
}

// OutputPanel shows one channel, following its tail.
func OutputPanel(c *ui.Context, channels []OutputChannel, s *OutputPanelState, opts OutputPanelOptions) OutputPanelResult {
	if s == nil || len(channels) == 0 {
		panic("mujica: OutputPanel needs a state and a channel")
	}
	if s.Channel < 0 || s.Channel >= len(channels) {
		panic("mujica: OutputPanel Channel out of range")
	}
	if !s.started {
		s.started, s.Follow, s.shown = true, true, s.Channel
	}
	var res OutputPanelResult
	res.Element = codePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	res.Element.Children(func() {
		codeHeader(c, func() {
			options := make([]SelectOption[int], len(channels))
			for i, ch := range channels {
				options[i] = SelectOption[int]{Value: i, Label: ch.Name}
			}
			if Select(c, &s.Channel, options, SelectOptions{Label: msg(c, "output_panel.channel"), Disabled: opts.Disabled}).Changed() {
				res.Changed = true
				slog.Debug("mujica: output_panel channel", "channel", s.Channel)
			}
			ui.Spacer(c)
			was := s.Follow
			if Toggle(c, &s.Follow, msg(c, "log_viewer.follow"), ToggleOptions{Disabled: opts.Disabled}).Changed() && s.Follow && !was {
				s.rows.ScrollToEnd()
			}
			if iconAction(c, icons.Must("trash-2"), msg(c, "output_panel.clear")).Disabled(opts.Disabled).Clicked() && !disabled {
				res.Cleared = channels[s.Channel].Name
				slog.Debug("mujica: output_panel clear", "channel", s.Channel)
			}
		})
		if s.shown != s.Channel {
			s.shown, s.rows = s.Channel, ui.ListState{}
		}
		lines := channels[s.Channel].Lines
		if len(lines) == 0 {
			codeCaption(c, msg(c, "output_panel.empty")).Padding(Density(c).Unit() * 4)
			return
		}
		d := Density(c)
		s.rows.FollowEnd = s.Follow
		s.rows.Label = func(i int) string { return lines[i] }
		ui.List(c, &s.rows, len(lines), func(i int) {
			ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).AlignItems(ui.Center).Children(func() {
				ui.RichText(c, codeAnsiSpans(c, lines[i])...).Font(theme.MonoFont).FontSize(FontSize(c, codeSize)).Grow(1).MinWidth(0)
			})
		}).Grow(1).MinHeight(d.ControlHeight()*3).Padding(d.Unit(), 0)
	})
	return res
}
