package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/code"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func codeSymbolDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        code.SymbolOutlineState
		cursor   float64
		disabled bool
	}
	s := ui.Local(c.Root(), "333", func() *state { return &state{cursor: 15} })
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		ui.Text(c, fmt.Sprintf("Cursor line %.0f", st.cursor))
		if input.Slider(c, &st.cursor, input.SliderOptions{Min: 1, Max: 34, Step: 1}).Changed() {
			log(fmt.Sprintf("cursor %.0f", st.cursor))
		}
	})
	symbols := []code.Symbol{
		{Name: "Decree", Kind: code.SymbolType, Line: 6, EndLine: 9, Detail: "struct", Children: []code.Symbol{
			{Name: "ID", Kind: code.SymbolField, Line: 7, EndLine: 7, Detail: "int"},
			{Name: "Title", Kind: code.SymbolField, Line: 8, EndLine: 8, Detail: "string"},
		}},
		{Name: "Seal", Kind: code.SymbolFunction, Line: 13, EndLine: 21, Detail: "(d Decree) (string, error)"},
		{Name: "Herald", Kind: code.SymbolFunction, Line: 23, EndLine: 31, Detail: "(ds []Decree) int"},
	}
	r := code.SymbolOutline(c, symbols, &st.s, code.SymbolOutlineOptions{CursorLine: int(st.cursor), Disabled: st.disabled})
	r.Element.Size(420, 320)
	if r.Jumped != nil {
		log(fmt.Sprintf("jump to %s, line %d", r.Jumped.Name, r.Jumped.Line))
	}
}

func codeCompletionDemo(c *ui.Context, log func(string)) {
	type state struct {
		text     string
		open     bool
		s        code.CompletionMenuState
		disabled bool
	}
	s := ui.Local(c.Root(), "334", func() *state { return &state{} })
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	all := []code.Completion{
		{Label: "Seal", Kind: code.SymbolFunction, Detail: "func(d Decree) (string, error)", Doc: "Seal stamps a decree with the royal wax."},
		{Label: "Search", Kind: code.SymbolFunction, Detail: "func(q string) []Decree"},
		{Label: "Season", Kind: code.SymbolConstant, Detail: "= 4"},
		{Label: "Herald", Kind: code.SymbolFunction, Detail: "func(ds []Decree) int"},
		{Label: "Decree", Kind: code.SymbolType, Detail: "struct", Doc: "Decree is an order of the court."},
	}
	word := st.text[strings.LastIndexAny(st.text, " .(")+1:]
	var items []code.Completion
	for _, it := range all {
		if word != "" && strings.HasPrefix(strings.ToLower(it.Label), strings.ToLower(word)) {
			items = append(items, it)
		}
	}
	ui.Column(c).Width(420).Gap(8).Children(func() {
		ui.Text(c, "Type “se” or “de” to see suggestions.").TextColor(core.Tokens(c).TextMuted)
		in := input.TextInput(c, &st.text, input.TextInputOptions{Label: "Code", Placeholder: "court.", Disabled: st.disabled})
		if in.Changed() {
			st.open = len(items) > 0
		}
		r := code.CompletionMenu(c, in.Input, &st.open, items, &st.s, code.CompletionMenuOptions{Disabled: st.disabled})
		if r.Accepted != nil {
			st.text = st.text[:len(st.text)-len(word)] + r.Accepted.Label
			log("accept " + r.Accepted.Label)
		}
	})
}

func codeDebugDemo(c *ui.Context, log func(string)) {
	type state struct {
		dbg      code.DebugState
		disabled bool
	}
	s := ui.Local(c.Root(), "335", func() *state { return &state{dbg: code.DebugPaused} })
	st := *s
	names := []string{"Stopped", "Running", "Paused"}
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		ui.Text(c, "State: "+names[st.dbg])
	})
	a := code.DebugToolbar(c, st.dbg, code.DebugToolbarOptions{Disabled: st.disabled}).Action
	switch a {
	case code.DebugNone:
		return
	case code.DebugContinue, code.DebugRestart:
		st.dbg = code.DebugRunning
	case code.DebugPause, code.DebugStepOver, code.DebugStepInto, code.DebugStepOut:
		st.dbg = code.DebugPaused
	case code.DebugStop:
		st.dbg = code.DebugStopped
	}
	log(fmt.Sprintf("action %d, now %s", a, names[st.dbg]))
}

func codeBreakpointDemo(c *ui.Context, log func(string)) {
	type state struct {
		bps      []code.Breakpoint
		disabled bool
	}
	s := ui.Local(c.Root(), "336", func() *state {
		return &state{bps: []code.Breakpoint{
			{ID: "1", File: "court/seal.go", Line: 14, Enabled: true},
			{ID: "2", File: "court/ledger.go", Line: 40, Condition: "n > 3"},
			{ID: "3", File: "main.go", Line: 9, Enabled: true},
		}}
	})
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	r := code.BreakpointList(c, &st.bps, code.BreakpointListOptions{Disabled: st.disabled})
	r.Element.Size(480, 240)
	for _, e := range []struct{ what, id string }{{"toggled", r.Toggled}, {"condition", r.Edited}, {"open", r.Opened}} {
		if e.id != "" {
			log(e.what + " " + e.id)
		}
	}
}

func codeCallStackDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        code.CallStackState
		disabled bool
	}
	s := ui.Local(c.Root(), "337", func() *state { return &state{s: code.CallStackState{Current: "f1"}} })
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	threads := []code.DebugThread{
		{ID: "t1", Name: "main", Frames: []code.StackFrame{
			{ID: "f1", Name: "court.Seal", File: "court/seal.go", Line: 14},
			{ID: "f2", Name: "court.Herald", File: "court/seal.go", Line: 27},
			{ID: "f3", Name: "main.main", File: "main.go", Line: 9},
		}},
		{ID: "t2", Name: "ledger writer", Frames: []code.StackFrame{{ID: "f4", Name: "ledger.Flush", File: "ledger/write.go", Line: 51}}},
	}
	r := code.CallStack(c, threads, &st.s, code.CallStackOptions{Disabled: st.disabled})
	r.Element.Size(440, 300)
	if r.Changed() {
		log("frame " + st.s.Current)
	}
}

func codeVariablesDemo(c *ui.Context, log func(string)) {
	type state struct {
		vars     []code.Variable
		s        code.VariablesPanelState
		disabled bool
	}
	s := ui.Local(c.Root(), "338", func() *state {
		return &state{vars: []code.Variable{
			{Name: "d", Type: "court.Decree", Value: "{…}", Children: []code.Variable{
				{Name: "ID", Value: "42", Type: "int", Changed: true},
				{Name: "Title", Value: `"Of seals"`, Type: "string"},
			}},
			{Name: "i", Value: "2", Type: "int", Changed: true},
			{Name: "err", Value: "nil", Type: "error"},
		}}
	})
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	r := code.VariablesPanel(c, st.vars, &st.s, code.VariablesPanelOptions{Disabled: st.disabled})
	r.Element.Size(440, 260)
	if r.Edited != nil {
		log(fmt.Sprintf("set %s = %s", r.Edited.Path, r.Edited.Value))
		codeDemoSet(st.vars, strings.Split(r.Edited.Path, "."), r.Edited.Value)
	}
}

// codeDemoSet writes value at the dotted path.
func codeDemoSet(vars []code.Variable, path []string, value string) {
	for i := range vars {
		if vars[i].Name != path[0] {
			continue
		}
		if len(path) == 1 {
			vars[i].Value, vars[i].Changed = value, true
			return
		}
		codeDemoSet(vars[i].Children, path[1:], value)
	}
}

func codeHexDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        code.HexViewerState
		disabled bool
	}
	s := ui.Local(c.Root(), "339", func() *state { return &state{} })
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	data := []byte(strings.Join(codeDemoSource, "\n"))
	r := code.HexViewer(c, data, &st.s, code.HexViewerOptions{Disabled: st.disabled})
	r.Element.Size(720, 300)
	if r.Changed() {
		a, b := st.s.Selection()
		log(fmt.Sprintf("bytes %d..%d", a, b))
	}
}

func codeFlamegraphDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        code.FlamegraphState
		status   string
		disabled bool
	}
	s := ui.Local(c.Root(), "340", func() *state { return &state{status: "Ready"} })
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		opts := []input.SelectOption[string]{{Value: "Ready", Label: "Ready"}, {Value: "Loading", Label: "Loading"}, {Value: "Failed", Label: "Failed"}}
		if input.Select(c, &st.status, opts, input.SelectOptions{Label: "Status"}).Changed() {
			log("status " + st.status)
		}
	})
	root := &code.FlameNode{Name: "main", Value: 240, Children: []code.FlameNode{
		{Name: "court.Herald", Value: 150, Children: []code.FlameNode{
			{Name: "court.Seal", Value: 110, Children: []code.FlameNode{{Name: "wax.Melt", Value: 70}, {Name: "ink.Dry", Value: 25}}},
			{Name: "ledger.Write", Value: 30},
		}},
		{Name: "archive.Load", Value: 60, Children: []code.FlameNode{{Name: "json.Decode", Value: 45}}},
		{Name: "gc", Value: 20},
	}}
	status := map[string]data.DataStatus{"Ready": data.DataReady, "Loading": data.DataLoading, "Failed": data.DataFailed}[st.status]
	ui.Column(c).Width(680).Children(func() {
		if code.Flamegraph(c, root, &st.s, code.FlamegraphOptions{Unit: "ms", Status: status, Error: "Profile could not be read", Disabled: st.disabled}).Changed() {
			log(fmt.Sprintf("zoom %v focus %v", st.s.Zoom, st.s.Focus))
		}
	})
}
