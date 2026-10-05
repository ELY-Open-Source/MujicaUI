package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func codeSymbolDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        mujica.SymbolOutlineState
		cursor   float64
		disabled bool
	}
	s := ui.Local(c.Root(), "333", func() *state { return &state{cursor: 15} })
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		ui.Text(c, fmt.Sprintf("Cursor line %.0f", st.cursor))
		if mujica.Slider(c, &st.cursor, mujica.SliderOptions{Min: 1, Max: 34, Step: 1}).Changed() {
			log(fmt.Sprintf("cursor %.0f", st.cursor))
		}
	})
	symbols := []mujica.Symbol{
		{Name: "Decree", Kind: mujica.SymbolType, Line: 6, EndLine: 9, Detail: "struct", Children: []mujica.Symbol{
			{Name: "ID", Kind: mujica.SymbolField, Line: 7, EndLine: 7, Detail: "int"},
			{Name: "Title", Kind: mujica.SymbolField, Line: 8, EndLine: 8, Detail: "string"},
		}},
		{Name: "Seal", Kind: mujica.SymbolFunction, Line: 13, EndLine: 21, Detail: "(d Decree) (string, error)"},
		{Name: "Herald", Kind: mujica.SymbolFunction, Line: 23, EndLine: 31, Detail: "(ds []Decree) int"},
	}
	r := mujica.SymbolOutline(c, symbols, &st.s, mujica.SymbolOutlineOptions{CursorLine: int(st.cursor), Disabled: st.disabled})
	r.Element.Size(420, 320)
	if r.Jumped != nil {
		log(fmt.Sprintf("jump to %s, line %d", r.Jumped.Name, r.Jumped.Line))
	}
}

func codeCompletionDemo(c *ui.Context, log func(string)) {
	type state struct {
		text     string
		open     bool
		s        mujica.CompletionMenuState
		disabled bool
	}
	s := ui.Local(c.Root(), "334", func() *state { return &state{} })
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	all := []mujica.Completion{
		{Label: "Seal", Kind: mujica.SymbolFunction, Detail: "func(d Decree) (string, error)", Doc: "Seal stamps a decree with the royal wax."},
		{Label: "Search", Kind: mujica.SymbolFunction, Detail: "func(q string) []Decree"},
		{Label: "Season", Kind: mujica.SymbolConstant, Detail: "= 4"},
		{Label: "Herald", Kind: mujica.SymbolFunction, Detail: "func(ds []Decree) int"},
		{Label: "Decree", Kind: mujica.SymbolType, Detail: "struct", Doc: "Decree is an order of the court."},
	}
	word := st.text[strings.LastIndexAny(st.text, " .(")+1:]
	var items []mujica.Completion
	for _, it := range all {
		if word != "" && strings.HasPrefix(strings.ToLower(it.Label), strings.ToLower(word)) {
			items = append(items, it)
		}
	}
	ui.Column(c).Width(420).Gap(8).Children(func() {
		ui.Text(c, "Type “se” or “de” to see suggestions.").TextColor(mujica.Tokens(c).TextMuted)
		in := mujica.TextInput(c, &st.text, mujica.TextInputOptions{Label: "Code", Placeholder: "court.", Disabled: st.disabled})
		if in.Changed() {
			st.open = len(items) > 0
		}
		r := mujica.CompletionMenu(c, in.Input, &st.open, items, &st.s, mujica.CompletionMenuOptions{Disabled: st.disabled})
		if r.Accepted != nil {
			st.text = st.text[:len(st.text)-len(word)] + r.Accepted.Label
			log("accept " + r.Accepted.Label)
		}
	})
}

func codeDebugDemo(c *ui.Context, log func(string)) {
	type state struct {
		dbg      mujica.DebugState
		disabled bool
	}
	s := ui.Local(c.Root(), "335", func() *state { return &state{dbg: mujica.DebugPaused} })
	st := *s
	names := []string{"Stopped", "Running", "Paused"}
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		ui.Text(c, "State: "+names[st.dbg])
	})
	a := mujica.DebugToolbar(c, st.dbg, mujica.DebugToolbarOptions{Disabled: st.disabled}).Action
	switch a {
	case mujica.DebugNone:
		return
	case mujica.DebugContinue, mujica.DebugRestart:
		st.dbg = mujica.DebugRunning
	case mujica.DebugPause, mujica.DebugStepOver, mujica.DebugStepInto, mujica.DebugStepOut:
		st.dbg = mujica.DebugPaused
	case mujica.DebugStop:
		st.dbg = mujica.DebugStopped
	}
	log(fmt.Sprintf("action %d, now %s", a, names[st.dbg]))
}

func codeBreakpointDemo(c *ui.Context, log func(string)) {
	type state struct {
		bps      []mujica.Breakpoint
		disabled bool
	}
	s := ui.Local(c.Root(), "336", func() *state {
		return &state{bps: []mujica.Breakpoint{
			{ID: "1", File: "court/seal.go", Line: 14, Enabled: true},
			{ID: "2", File: "court/ledger.go", Line: 40, Condition: "n > 3"},
			{ID: "3", File: "main.go", Line: 9, Enabled: true},
		}}
	})
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	r := mujica.BreakpointList(c, &st.bps, mujica.BreakpointListOptions{Disabled: st.disabled})
	r.Element.Size(480, 240)
	for _, e := range []struct{ what, id string }{{"toggled", r.Toggled}, {"condition", r.Edited}, {"open", r.Opened}} {
		if e.id != "" {
			log(e.what + " " + e.id)
		}
	}
}

func codeCallStackDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        mujica.CallStackState
		disabled bool
	}
	s := ui.Local(c.Root(), "337", func() *state { return &state{s: mujica.CallStackState{Current: "f1"}} })
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	threads := []mujica.DebugThread{
		{ID: "t1", Name: "main", Frames: []mujica.StackFrame{
			{ID: "f1", Name: "court.Seal", File: "court/seal.go", Line: 14},
			{ID: "f2", Name: "court.Herald", File: "court/seal.go", Line: 27},
			{ID: "f3", Name: "main.main", File: "main.go", Line: 9},
		}},
		{ID: "t2", Name: "ledger writer", Frames: []mujica.StackFrame{{ID: "f4", Name: "ledger.Flush", File: "ledger/write.go", Line: 51}}},
	}
	r := mujica.CallStack(c, threads, &st.s, mujica.CallStackOptions{Disabled: st.disabled})
	r.Element.Size(440, 300)
	if r.Changed() {
		log("frame " + st.s.Current)
	}
}

func codeVariablesDemo(c *ui.Context, log func(string)) {
	type state struct {
		vars     []mujica.Variable
		s        mujica.VariablesPanelState
		disabled bool
	}
	s := ui.Local(c.Root(), "338", func() *state {
		return &state{vars: []mujica.Variable{
			{Name: "d", Type: "court.Decree", Value: "{…}", Children: []mujica.Variable{
				{Name: "ID", Value: "42", Type: "int", Changed: true},
				{Name: "Title", Value: `"Of seals"`, Type: "string"},
			}},
			{Name: "i", Value: "2", Type: "int", Changed: true},
			{Name: "err", Value: "nil", Type: "error"},
		}}
	})
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	r := mujica.VariablesPanel(c, st.vars, &st.s, mujica.VariablesPanelOptions{Disabled: st.disabled})
	r.Element.Size(440, 260)
	if r.Edited != nil {
		log(fmt.Sprintf("set %s = %s", r.Edited.Path, r.Edited.Value))
		codeDemoSet(st.vars, strings.Split(r.Edited.Path, "."), r.Edited.Value)
	}
}

// codeDemoSet writes value at the dotted path.
func codeDemoSet(vars []mujica.Variable, path []string, value string) {
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
		s        mujica.HexViewerState
		disabled bool
	}
	s := ui.Local(c.Root(), "339", func() *state { return &state{} })
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	data := []byte(strings.Join(codeDemoSource, "\n"))
	r := mujica.HexViewer(c, data, &st.s, mujica.HexViewerOptions{Disabled: st.disabled})
	r.Element.Size(720, 300)
	if r.Changed() {
		a, b := st.s.Selection()
		log(fmt.Sprintf("bytes %d..%d", a, b))
	}
}

func codeFlamegraphDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        mujica.FlamegraphState
		status   string
		disabled bool
	}
	s := ui.Local(c.Root(), "340", func() *state { return &state{status: "Ready"} })
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		opts := []mujica.SelectOption[string]{{Value: "Ready", Label: "Ready"}, {Value: "Loading", Label: "Loading"}, {Value: "Failed", Label: "Failed"}}
		if mujica.Select(c, &st.status, opts, mujica.SelectOptions{Label: "Status"}).Changed() {
			log("status " + st.status)
		}
	})
	root := &mujica.FlameNode{Name: "main", Value: 240, Children: []mujica.FlameNode{
		{Name: "court.Herald", Value: 150, Children: []mujica.FlameNode{
			{Name: "court.Seal", Value: 110, Children: []mujica.FlameNode{{Name: "wax.Melt", Value: 70}, {Name: "ink.Dry", Value: 25}}},
			{Name: "ledger.Write", Value: 30},
		}},
		{Name: "archive.Load", Value: 60, Children: []mujica.FlameNode{{Name: "json.Decode", Value: 45}}},
		{Name: "gc", Value: 20},
	}}
	status := map[string]mujica.DataStatus{"Ready": mujica.DataReady, "Loading": mujica.DataLoading, "Failed": mujica.DataFailed}[st.status]
	ui.Column(c).Width(680).Children(func() {
		if mujica.Flamegraph(c, root, &st.s, mujica.FlamegraphOptions{Unit: "ms", Status: status, Error: "Profile could not be read", Disabled: st.disabled}).Changed() {
			log(fmt.Sprintf("zoom %v focus %v", st.s.Zoom, st.s.Focus))
		}
	})
}
