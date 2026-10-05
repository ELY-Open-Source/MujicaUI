package mujica

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/egoist/mygo/ui"
)

// DebugState is where a debug session is.
type DebugState int

const (
	DebugStopped DebugState = iota
	DebugRunning
	DebugPaused
)

// DebugAction is what the user asked the debugger to do.
type DebugAction int

const (
	DebugNone DebugAction = iota
	// DebugContinue resumes a paused session or starts a stopped one.
	DebugContinue
	DebugPause
	DebugStepOver
	DebugStepInto
	DebugStepOut
	DebugRestart
	DebugStop
)

// DebugToolbarOptions configure a DebugToolbar.
type DebugToolbarOptions struct {
	Disabled bool
}

// DebugToolbarResult carries the requested action.
type DebugToolbarResult struct {
	Element *ui.Element
	Action  DebugAction
}

// debugIcon draws a toolbar glyph in a 24-unit box.
func debugIcon(body string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
}

var debugIcons = map[DebugAction]*ui.SVG{
	DebugContinue: debugIcon(`<path d="M7 4l12 8-12 8z" fill="currentColor"/>`),
	DebugPause:    debugIcon(`<path d="M8 5v14M16 5v14"/>`),
	DebugStepOver: debugIcon(`<path d="M4 13a8 8 0 0 1 15-3"/><path d="M20 5v5h-5"/><circle cx="12" cy="19" r="2" fill="currentColor"/>`),
	DebugStepInto: debugIcon(`<path d="M12 3v11"/><path d="M7 9l5 5 5-5"/><circle cx="12" cy="20" r="2" fill="currentColor"/>`),
	DebugStepOut:  debugIcon(`<path d="M12 15V4"/><path d="M7 9l5-5 5 5"/><circle cx="12" cy="20" r="2" fill="currentColor"/>`),
	DebugRestart:  debugIcon(`<path d="M4 12a8 8 0 1 0 3-6.2"/><path d="M4 4v5h5"/>`),
	DebugStop:     debugIcon(`<rect x="6" y="6" width="12" height="12" rx="1" fill="currentColor"/>`),
}

// debugShortcut is an action's key, as in common IDEs.
type debugShortcut struct {
	mods ui.Modifiers
	key  ui.Key
	name string
}

func debugShortcuts() map[DebugAction]debugShortcut {
	cmd := "Ctrl"
	if runtime.GOOS == "darwin" {
		cmd = "⌘"
	}
	return map[DebugAction]debugShortcut{
		DebugContinue: {0, ui.KeyF5, "F5"},
		DebugPause:    {0, ui.KeyF6, "F6"},
		DebugStepOver: {0, ui.KeyF10, "F10"},
		DebugStepInto: {0, ui.KeyF11, "F11"},
		DebugStepOut:  {ui.Shift, ui.KeyF11, "Shift+F11"},
		DebugRestart:  {ui.Cmd | ui.Shift, ui.KeyF5, cmd + "+Shift+F5"},
		DebugStop:     {ui.Shift, ui.KeyF5, "Shift+F5"},
	}
}

// DebugEnabled reports whether action applies in state.
func DebugEnabled(state DebugState, a DebugAction) bool {
	switch a {
	case DebugContinue:
		return state != DebugRunning
	case DebugPause:
		return state == DebugRunning
	case DebugStepOver, DebugStepInto, DebugStepOut:
		return state == DebugPaused
	case DebugRestart, DebugStop:
		return state != DebugStopped
	}
	return false
}

// DebugToolbar shows debugger actions enabled by state.
func DebugToolbar(c *ui.Context, state DebugState, opts DebugToolbarOptions) DebugToolbarResult {
	if state < DebugStopped || state > DebugPaused {
		panic("mujica: DebugState out of range")
	}
	keys := debugShortcuts()
	var res DebugToolbarResult
	res.Element = Toolbar(c, ToolbarOptions{Label: msg(c, "debug_toolbar.label")}, func() {
		first := DebugContinue
		if state == DebugRunning {
			first = DebugPause
		}
		for _, a := range []DebugAction{first, DebugStepOver, DebugStepInto, DebugStepOut, DebugRestart, DebugStop} {
			name := msg(c, debugKeys[a])
			if a == DebugContinue && state == DebugStopped {
				name = msg(c, "debug_toolbar.start")
			}
			on := DebugEnabled(state, a) && !opts.Disabled
			b := ToolbarButton(c, name, ToolbarButtonOptions{Icon: debugIcons[a], Disabled: !on})
			b.Tooltip(fmt.Sprintf("%s (%s)", name, keys[a].name))
			if b.Clicked() && on && !b.IsDisabled() {
				res.Action = a
			}
		}
	}).Disabled(opts.Disabled)
	if res.Action == DebugNone && !res.Element.IsDisabled() {
		for _, a := range []DebugAction{DebugContinue, DebugPause, DebugStepOver, DebugStepInto, DebugStepOut, DebugRestart, DebugStop} {
			if k := keys[a]; DebugEnabled(state, a) && c.Shortcut(k.mods, k.key) {
				res.Action = a
				break
			}
		}
	}
	if res.Action != DebugNone {
		slog.Debug("mujica: debug_toolbar action", "action", res.Action, "state", state)
	}
	return res
}

var debugKeys = map[DebugAction]string{
	DebugContinue: "debug_toolbar.continue", DebugPause: "debug_toolbar.pause", DebugStepOver: "debug_toolbar.stepOver",
	DebugStepInto: "debug_toolbar.stepInto", DebugStepOut: "debug_toolbar.stepOut", DebugRestart: "debug_toolbar.restart", DebugStop: "debug_toolbar.stop",
}
