package mujica_test

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleTerminal() {
	var term mujica.TerminalState
	term.Write("\x1b[1m~/court\x1b[0m $ ls\r\n\x1b[34mdecrees\x1b[0m  ledger.txt\r\n$ ")
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Terminal(c, &term, mujica.TerminalOptions{
			Input:   func(data string) { fmt.Printf("to pty: %q\n", data) },
			Resized: func(cols, rows int) { fmt.Println("pty size", cols, rows) },
		}).Size(380, 200)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTerminalTabs() {
	tabs := []mujica.TerminalTab{{ID: "a", Title: "zsh"}, {ID: "b", Title: "build", Busy: true}}
	active := "a"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.TerminalTabs(c, &tabs, &active, mujica.TerminalTabsOptions{}, func(id string) {
			ui.Text(c, "terminal "+id)
		})
		if r.Closed != "" {
			fmt.Println("close pty", r.Closed)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTerminalSearch() {
	var term mujica.TerminalState
	term.Write("build ok\r\nerror: missing seal\r\n")
	var search mujica.TerminalSearchState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Column(c).Gap(8).Children(func() {
			mujica.TerminalSearch(c, &term, &search, mujica.TerminalSearchOptions{})
			mujica.Terminal(c, &term, mujica.TerminalOptions{}).Height(160)
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommandHistory() {
	history := []mujica.HistoryEntry{{Command: "git status"}, {Command: "go test ./...", ExitCode: 1}}
	var state mujica.CommandHistoryState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.CommandHistory(c, history, &state, mujica.CommandHistoryOptions{})
		r.Element.Height(240)
		if r.Replayed != nil {
			fmt.Println("run", r.Replayed.Command)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLogViewer() {
	lines := []mujica.LogLine{{Level: mujica.LogInfo, Text: "court opened"}, {Level: mujica.LogError, Text: "seal missing"}}
	var state mujica.LogViewerState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.LogViewer(c, lines, &state, mujica.LogViewerOptions{}).Element.Height(240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProcessList() {
	procs := []mujica.Process{{PID: 7, Name: "archivist", CPU: 42, Memory: 50 << 20}, {PID: 9, Name: "herald", CPU: 3, Memory: 8 << 20}}
	var state mujica.ProcessListState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.ProcessList(c, procs, &state, mujica.ProcessListOptions{})
		r.Element.Height(240)
		if r.End != nil {
			fmt.Println("kill", r.End.PID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCodeViewer() {
	lines := []string{"func Seal(d string) int {", "\treturn len(d)", "}"}
	var state mujica.CodeViewerState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CodeViewer(c, lines, &state, mujica.CodeViewerOptions{Language: "go"}).Element.Height(240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMinimap() {
	lines := []string{"package court", "", "func Seal() {}", "func Herald() {}"}
	view := mujica.MinimapView{Top: 0, Visible: 2}
	draw := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Minimap(c, lines, &view, mujica.MinimapOptions{Language: "go"}).Changed() {
			fmt.Println("scroll editor to", view.Top)
		}
	}
	ui.Render(draw, 400, 300, 1)
}

func ExampleFindWidget() {
	text := "the seal and the seal"
	state := mujica.FindState{Query: "seal", Replacement: "signet", ShowReplace: true}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.FindWidget(c, &text, &state, mujica.FindWidgetOptions{}).Replaced > 0 {
			fmt.Println(text)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSearchPanel() {
	results := []mujica.SearchFileResult{{Path: "court/seal.go", Matches: []mujica.SearchMatch{{Line: 12, Text: "func Seal() {", Start: 5, End: 9}}}}
	state := mujica.SearchPanelState{Query: "seal"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.SearchPanel(c, results, &state, mujica.SearchPanelOptions{})
		r.Element.Height(280)
		if r.Opened != nil {
			fmt.Println("open", r.Opened.Path, r.Opened.Line)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProblemsPanel() {
	problems := []mujica.Problem{{Severity: mujica.ProblemError, File: "court/seal.go", Line: 12, Col: 5, Message: "undefined: Decree"}}
	var state mujica.ProblemsPanelState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.ProblemsPanel(c, problems, &state, mujica.ProblemsPanelOptions{})
		r.Element.Height(240)
		if r.Opened != nil {
			fmt.Println("go to", r.Opened.File, r.Opened.Line)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOutputPanel() {
	channels := []mujica.OutputChannel{{Name: "Build", Lines: []string{"compiling court", "\x1b[32mok\x1b[0m"}}}
	var state mujica.OutputPanelState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.OutputPanel(c, channels, &state, mujica.OutputPanelOptions{})
		r.Element.Height(240)
		if r.Cleared != "" {
			channels[state.Channel].Lines = nil
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSymbolOutline() {
	symbols := []mujica.Symbol{{Name: "Court", Kind: mujica.SymbolType, Line: 3, EndLine: 40, Children: []mujica.Symbol{
		{Name: "Seal", Kind: mujica.SymbolMethod, Line: 10, EndLine: 20},
	}}}
	var state mujica.SymbolOutlineState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.SymbolOutline(c, symbols, &state, mujica.SymbolOutlineOptions{CursorLine: 12})
		r.Element.Height(240)
		if r.Jumped != nil {
			fmt.Println("go to line", r.Jumped.Line)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCompletionMenu() {
	text := "se"
	open := true
	items := []mujica.Completion{{Label: "Seal", Kind: mujica.SymbolMethod, Detail: "func()"}, {Label: "Search", Kind: mujica.SymbolFunction}}
	var state mujica.CompletionMenuState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		in := mujica.TextInput(c, &text, mujica.TextInputOptions{Label: "Code"})
		if r := mujica.CompletionMenu(c, in.Input, &open, items, &state, mujica.CompletionMenuOptions{}); r.Accepted != nil {
			text = r.Accepted.Label
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDebugToolbar() {
	state := mujica.DebugPaused
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		switch mujica.DebugToolbar(c, state, mujica.DebugToolbarOptions{}).Action {
		case mujica.DebugContinue:
			state = mujica.DebugRunning
		case mujica.DebugStop:
			state = mujica.DebugStopped
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBreakpointList() {
	bps := []mujica.Breakpoint{{ID: "a", File: "court/seal.go", Line: 12, Enabled: true}, {ID: "b", File: "main.go", Line: 9, Condition: "n > 3"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.BreakpointList(c, &bps, mujica.BreakpointListOptions{})
		r.Element.Height(200)
		if r.Toggled != "" || r.Edited != "" {
			fmt.Println("sync breakpoints with the debugger")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCallStack() {
	threads := []mujica.DebugThread{{ID: "t1", Name: "main", Frames: []mujica.StackFrame{
		{ID: "f1", Name: "court.Seal", File: "court/seal.go", Line: 14},
		{ID: "f2", Name: "main.main", File: "main.go", Line: 9},
	}}}
	state := mujica.CallStackState{Current: "f1"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.CallStack(c, threads, &state, mujica.CallStackOptions{})
		r.Element.Height(200)
		if r.Changed() {
			fmt.Println("show frame", state.Current)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVariablesPanel() {
	vars := []mujica.Variable{{Name: "decree", Value: "{…}", Children: []mujica.Variable{{Name: "id", Value: "42", Changed: true}}}, {Name: "count", Value: "7"}}
	var state mujica.VariablesPanelState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.VariablesPanel(c, vars, &state, mujica.VariablesPanelOptions{})
		r.Element.Height(200)
		if r.Edited != nil {
			fmt.Println("set", r.Edited.Path, "=", r.Edited.Value)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleHexViewer() {
	data := []byte("MUJICA court ledger\x00\x01")
	var state mujica.HexViewerState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.HexViewer(c, data, &state, mujica.HexViewerOptions{}).Element.Height(160)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFlamegraph() {
	root := &mujica.FlameNode{Name: "main", Value: 100, Children: []mujica.FlameNode{
		{Name: "court.Seal", Value: 60}, {Name: "ledger.Write", Value: 30},
	}}
	var state mujica.FlamegraphState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Flamegraph(c, root, &state, mujica.FlamegraphOptions{Unit: "ms"})
	}
	ui.Render(view, 400, 300, 1)
}
