package code_test

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/code"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func ExampleTerminal() {
	var term code.TerminalState
	term.Write("\x1b[1m~/court\x1b[0m $ ls\r\n\x1b[34mdecrees\x1b[0m  ledger.txt\r\n$ ")
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		code.Terminal(c, &term, code.TerminalOptions{
			Input:   func(data string) { fmt.Printf("to pty: %q\n", data) },
			Resized: func(cols, rows int) { fmt.Println("pty size", cols, rows) },
		}).Size(380, 200)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTerminalTabs() {
	tabs := []code.TerminalTab{{ID: "a", Title: "zsh"}, {ID: "b", Title: "build", Busy: true}}
	active := "a"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.TerminalTabs(c, &tabs, &active, code.TerminalTabsOptions{}, func(id string) {
			ui.Text(c, "terminal "+id)
		})
		if r.Closed != "" {
			fmt.Println("close pty", r.Closed)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTerminalSearch() {
	var term code.TerminalState
	term.Write("build ok\r\nerror: missing seal\r\n")
	var search code.TerminalSearchState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Gap(8).Children(func() {
			code.TerminalSearch(c, &term, &search, code.TerminalSearchOptions{})
			code.Terminal(c, &term, code.TerminalOptions{}).Height(160)
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommandHistory() {
	history := []code.HistoryEntry{{Command: "git status"}, {Command: "go test ./...", ExitCode: 1}}
	var state code.CommandHistoryState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.CommandHistory(c, history, &state, code.CommandHistoryOptions{})
		r.Element.Height(240)
		if r.Replayed != nil {
			fmt.Println("run", r.Replayed.Command)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLogViewer() {
	lines := []code.LogLine{{Level: code.LogInfo, Text: "court opened"}, {Level: code.LogError, Text: "seal missing"}}
	var state code.LogViewerState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		code.LogViewer(c, lines, &state, code.LogViewerOptions{}).Element.Height(240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProcessList() {
	procs := []code.Process{{PID: 7, Name: "archivist", CPU: 42, Memory: 50 << 20}, {PID: 9, Name: "herald", CPU: 3, Memory: 8 << 20}}
	var state code.ProcessListState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.ProcessList(c, procs, &state, code.ProcessListOptions{})
		r.Element.Height(240)
		if r.End != nil {
			fmt.Println("kill", r.End.PID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCodeViewer() {
	lines := []string{"func Seal(d string) int {", "\treturn len(d)", "}"}
	var state code.CodeViewerState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		code.CodeViewer(c, lines, &state, code.CodeViewerOptions{Language: "go"}).Element.Height(240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMinimap() {
	lines := []string{"package court", "", "func Seal() {}", "func Herald() {}"}
	view := code.MinimapView{Top: 0, Visible: 2}
	draw := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if code.Minimap(c, lines, &view, code.MinimapOptions{Language: "go"}).Changed() {
			fmt.Println("scroll editor to", view.Top)
		}
	}
	ui.Render(draw, 400, 300, 1)
}

func ExampleFindWidget() {
	text := "the seal and the seal"
	state := code.FindState{Query: "seal", Replacement: "signet", ShowReplace: true}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if code.FindWidget(c, &text, &state, code.FindWidgetOptions{}).Replaced > 0 {
			fmt.Println(text)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSearchPanel() {
	results := []code.SearchFileResult{{Path: "court/seal.go", Matches: []code.SearchMatch{{Line: 12, Text: "func Seal() {", Start: 5, End: 9}}}}
	state := code.SearchPanelState{Query: "seal"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.SearchPanel(c, results, &state, code.SearchPanelOptions{})
		r.Element.Height(280)
		if r.Opened != nil {
			fmt.Println("open", r.Opened.Path, r.Opened.Line)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProblemsPanel() {
	problems := []code.Problem{{Severity: code.ProblemError, File: "court/seal.go", Line: 12, Col: 5, Message: "undefined: Decree"}}
	var state code.ProblemsPanelState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.ProblemsPanel(c, problems, &state, code.ProblemsPanelOptions{})
		r.Element.Height(240)
		if r.Opened != nil {
			fmt.Println("go to", r.Opened.File, r.Opened.Line)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOutputPanel() {
	channels := []code.OutputChannel{{Name: "Build", Lines: []string{"compiling court", "\x1b[32mok\x1b[0m"}}}
	var state code.OutputPanelState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.OutputPanel(c, channels, &state, code.OutputPanelOptions{})
		r.Element.Height(240)
		if r.Cleared != "" {
			channels[state.Channel].Lines = nil
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSymbolOutline() {
	symbols := []code.Symbol{{Name: "Court", Kind: code.SymbolType, Line: 3, EndLine: 40, Children: []code.Symbol{
		{Name: "Seal", Kind: code.SymbolMethod, Line: 10, EndLine: 20},
	}}}
	var state code.SymbolOutlineState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.SymbolOutline(c, symbols, &state, code.SymbolOutlineOptions{CursorLine: 12})
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
	items := []code.Completion{{Label: "Seal", Kind: code.SymbolMethod, Detail: "func()"}, {Label: "Search", Kind: code.SymbolFunction}}
	var state code.CompletionMenuState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		in := input.TextInput(c, &text, input.TextInputOptions{Label: "Code"})
		if r := code.CompletionMenu(c, in.Input, &open, items, &state, code.CompletionMenuOptions{}); r.Accepted != nil {
			text = r.Accepted.Label
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDebugToolbar() {
	state := code.DebugPaused
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		switch code.DebugToolbar(c, state, code.DebugToolbarOptions{}).Action {
		case code.DebugContinue:
			state = code.DebugRunning
		case code.DebugStop:
			state = code.DebugStopped
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBreakpointList() {
	bps := []code.Breakpoint{{ID: "a", File: "court/seal.go", Line: 12, Enabled: true}, {ID: "b", File: "main.go", Line: 9, Condition: "n > 3"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.BreakpointList(c, &bps, code.BreakpointListOptions{})
		r.Element.Height(200)
		if r.Toggled != "" || r.Edited != "" {
			fmt.Println("sync breakpoints with the debugger")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCallStack() {
	threads := []code.DebugThread{{ID: "t1", Name: "main", Frames: []code.StackFrame{
		{ID: "f1", Name: "court.Seal", File: "court/seal.go", Line: 14},
		{ID: "f2", Name: "main.main", File: "main.go", Line: 9},
	}}}
	state := code.CallStackState{Current: "f1"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.CallStack(c, threads, &state, code.CallStackOptions{})
		r.Element.Height(200)
		if r.Changed() {
			fmt.Println("show frame", state.Current)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVariablesPanel() {
	vars := []code.Variable{{Name: "decree", Value: "{…}", Children: []code.Variable{{Name: "id", Value: "42", Changed: true}}}, {Name: "count", Value: "7"}}
	var state code.VariablesPanelState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := code.VariablesPanel(c, vars, &state, code.VariablesPanelOptions{})
		r.Element.Height(200)
		if r.Edited != nil {
			fmt.Println("set", r.Edited.Path, "=", r.Edited.Value)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleHexViewer() {
	data := []byte("MUJICA court ledger\x00\x01")
	var state code.HexViewerState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		code.HexViewer(c, data, &state, code.HexViewerOptions{}).Element.Height(160)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFlamegraph() {
	root := &code.FlameNode{Name: "main", Value: 100, Children: []code.FlameNode{
		{Name: "court.Seal", Value: 60}, {Name: "ledger.Write", Value: 30},
	}}
	var state code.FlamegraphState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		code.Flamegraph(c, root, &state, code.FlamegraphOptions{Unit: "ms"})
	}
	ui.Render(view, 400, 300, 1)
}
