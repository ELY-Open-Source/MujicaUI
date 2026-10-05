package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
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

// ModeSelector switches among chat, search, deep research and agent modes; arrows move the choice.
func ModeSelector(c *ui.Context, mode *ChatMode, opts ModeSelectorOptions) ChoiceResult {
	if *mode < ModeChat || *mode > ModeAgent {
		panic("mujica: unknown ChatMode")
	}
	segs := make([]Segment, len(chatModes))
	for i, m := range chatModes {
		segs[i] = Segment{Label: msg(c, m[0]), Icon: icons.Must(m[1]), HideLabel: opts.IconsOnly}
	}
	i := int(*mode)
	e := SegmentedControl(c, &i, segs, SegmentedControlOptions{Label: msg(c, "mode_selector.label"), Disabled: opts.Disabled})
	r := ChoiceResult{Element: e, changed: ChatMode(i) != *mode}
	if r.changed {
		*mode = ChatMode(i)
		slog.Debug("mujica: mode_selector changed", "mode", i)
	}
	return r
}
