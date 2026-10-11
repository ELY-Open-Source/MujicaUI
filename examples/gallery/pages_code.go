package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/code"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

const codeCategory = "Terminal & Code"

var codePages = []page{
	{ID: "321", Name: "Terminal", ZhName: "终端", Category: codeCategory, Demo: codeTerminalDemo, Code: codeTerminalCode,
		Keys: "Focused: keys and text go to the PTY; Ctrl+letters send control codes; Cmd+C/V (Ctrl+Shift+C/V elsewhere) copy and paste; Shift+PageUp/PageDown scroll back; drag selects."},
	{ID: "322", Name: "TerminalTabs", ZhName: "终端标签", Category: codeCategory, Demo: codeTerminalTabsDemo, Code: codeTerminalTabsCode,
		Keys: "Left/Right, Home/End switch tabs; F2 or double-click renames (Enter saves, Escape cancels); Delete closes, asking first while a process runs."},
	{ID: "323", Name: "TerminalSearch", ZhName: "终端搜索", Category: codeCategory, Demo: codeTerminalSearchDemo, Code: codeTerminalSearchCode,
		Keys: "Enter / Shift+Enter go to the next / previous match; Cmd+G (F3 elsewhere) too."},
	{ID: "324", Name: "CommandHistory", ZhName: "命令历史", Category: codeCategory, Demo: codeCommandHistoryDemo, Code: codeCommandHistoryCode,
		Keys: "Up/Down choose; Enter or double-click replays; typing in the list jumps by command."},
	{ID: "325", Name: "LogViewer", ZhName: "日志查看器", Category: codeCategory, Demo: codeLogViewerDemo, Code: codeLogViewerCode,
		Keys: "Space toggles a level or Follow; arrows, PageUp/PageDown, Home/End scroll the log."},
	{ID: "326", Name: "ProcessList", ZhName: "进程列表", Category: codeCategory, Demo: codeProcessListDemo, Code: codeProcessListCode,
		Keys: "Up/Down choose; Delete ends the chosen process after confirming; click a header to sort, again to reverse."},
	{ID: "327", Name: "CodeViewer", ZhName: "代码查看器", Category: codeCategory, Demo: codeCodeViewerDemo, Code: codeCodeViewerCode,
		Keys: "Up/Down move the current line; Shift extends; Cmd/Ctrl+C copies; Left folds, Right unfolds."},
	{ID: "328", Name: "Minimap", ZhName: "代码缩略图", Category: codeCategory, Demo: codeMinimapDemo, Code: codeMinimapCode,
		Keys: "Up/Down move a line, PageUp/PageDown a screen, Home/End to the ends; click or drag jumps."},
	{ID: "329", Name: "FindWidget", ZhName: "查找替换", Category: codeCategory, Demo: codeFindWidgetDemo, Code: codeFindWidgetCode,
		Keys: "Enter / Shift+Enter next / previous match; Enter in the replace field replaces the current match."},
	{ID: "330", Name: "SearchPanel", ZhName: "多文件搜索", Category: codeCategory, Demo: codeSearchPanelDemo, Code: codeSearchPanelCode,
		Keys: "Up/Down choose; Enter opens a match or folds a file; Left folds, Right unfolds a file."},
	{ID: "331", Name: "ProblemsPanel", ZhName: "问题面板", Category: codeCategory, Demo: codeProblemsDemo, Code: codeProblemsPanelCode,
		Keys: "Space toggles a severity; Up/Down choose; Enter opens a problem or folds a file; Left/Right fold and unfold."},
	{ID: "332", Name: "OutputPanel", ZhName: "输出面板", Category: codeCategory, Demo: codeOutputDemo, Code: codeOutputPanelCode,
		Keys: "Tab reaches the channel menu, Follow and Clear; arrows choose a channel in the open menu."},
	{ID: "333", Name: "SymbolOutline", ZhName: "符号大纲", Category: codeCategory, Demo: codeSymbolDemo, Code: codeSymbolOutlineCode,
		Keys: "Up/Down move; Right opens, Left closes; Enter or double-click jumps; typing jumps by name."},
	{ID: "334", Name: "CompletionMenu", ZhName: "补全菜单", Category: codeCategory, Demo: codeCompletionDemo, Code: codeCompletionMenuCode,
		Keys: "While typing: Down/Up move, Tab or Enter accepts, Escape closes; the editor keeps the focus."},
	{ID: "335", Name: "DebugToolbar", ZhName: "调试工具栏", Category: codeCategory, Demo: codeDebugDemo, Code: codeDebugToolbarCode,
		Keys: "F5 continue, F6 pause, F10 step over, F11 step into, Shift+F11 step out, Cmd/Ctrl+Shift+F5 restart, Shift+F5 stop; Left/Right inside the toolbar."},
	{ID: "336", Name: "BreakpointList", ZhName: "断点列表", Category: codeCategory, Demo: codeBreakpointDemo, Code: codeBreakpointListCode,
		Keys: "Space toggles; Enter on a location opens it; the condition field saves with Enter, cancels with Escape."},
	{ID: "337", Name: "CallStack", ZhName: "调用栈", Category: codeCategory, Demo: codeCallStackDemo, Code: codeCallStackCode,
		Keys: "Up/Down choose a frame; Enter folds a thread; Left/Right fold and unfold."},
	{ID: "338", Name: "VariablesPanel", ZhName: "变量面板", Category: codeCategory, Demo: codeVariablesDemo, Code: codeVariablesPanelCode,
		Keys: "Up/Down move; Right opens, Left closes; Enter, F2 or double-click edits a leaf; Enter saves, Escape cancels."},
	{ID: "339", Name: "HexViewer", ZhName: "十六进制查看器", Category: codeCategory, Demo: codeHexDemo, Code: codeHexViewerCode,
		Keys: "Arrows move by byte and row; Shift extends; Home/End go to the row ends; Cmd/Ctrl+C copies hex."},
	{ID: "340", Name: "Flamegraph", ZhName: "火焰图", Category: codeCategory, Demo: codeFlamegraphDemo, Code: codeFlamegraphCode,
		Keys: "Up parent, Down first child, Left/Right siblings; Enter zooms in, Escape zooms out; the readout is announced."},
}

// codeDemoCheck is a switch that logs its changes.
func codeDemoCheck(c *ui.Context, v *bool, label string, log func(string)) {
	if input.Switch(c, v, label, input.SwitchOptions{}).Changed() {
		log(fmt.Sprintf("%s: %v", label, *v))
	}
}

func codeDemoControls(c *ui.Context, fn func()) {
	ui.Row(c).Gap(16).Wrap().AlignItems(ui.Center).Margin(0, 0, 16, 0).Children(fn)
}

func codeTerminalDemo(c *ui.Context, log func(string)) {
	type state struct {
		term     code.TerminalState
		line     string
		disabled bool
	}
	s := ui.Local(c.Root(), "321", func() *state {
		s := &state{}
		s.term.Write("Welcome to the \x1b[1;35mcourt\x1b[0m shell. Type and press Enter.\r\n$ ")
		return s
	})
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		if input.Button(c, "Clear screen", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			st.term.Write("\x1b[2J\x1b[H$ ")
			log("clear screen")
		}
	})
	// A tiny echo shell stands in for the PTY.
	code.Terminal(c, &st.term, code.TerminalOptions{Disabled: st.disabled,
		Input: func(data string) {
			log(fmt.Sprintf("input %q", data))
			for _, r := range data {
				switch r {
				case '\r':
					st.term.Write("\r\n" + codeShell(st.line) + "$ ")
					st.line = ""
				case 0x7f:
					if st.line != "" {
						st.line = st.line[:len(st.line)-1]
						st.term.Write("\b \b")
					}
				default:
					if r >= ' ' {
						st.line += string(r)
						st.term.Write(string(r))
					}
				}
			}
		},
		Resized: func(cols, rows int) { log(fmt.Sprintf("resized %d×%d", cols, rows)) },
	}).Size(640, 260)
}

// codeShell answers a few commands for the demo terminal.
func codeShell(line string) string {
	switch f := strings.Fields(line); {
	case len(f) == 0:
		return ""
	case f[0] == "ls":
		return "\x1b[34mchronicles\x1b[0m  decrees.txt  \x1b[32mseal.sh\x1b[0m\r\n"
	case f[0] == "echo":
		return strings.Join(f[1:], " ") + "\r\n"
	case f[0] == "date":
		return time.Now().Format(time.RFC1123) + "\r\n"
	}
	return "\x1b[31m" + line + ": command not found\x1b[0m\r\n"
}

func codeTerminalTabsDemo(c *ui.Context, log func(string)) {
	type state struct {
		tabs     []code.TerminalTab
		active   string
		next     int
		disabled bool
	}
	s := ui.Local(c.Root(), "322", func() *state {
		return &state{tabs: []code.TerminalTab{{ID: "1", Title: "zsh"}, {ID: "2", Title: "build", Busy: true}, {ID: "3", Title: "logs"}}, active: "1", next: 4}
	})
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		if input.Button(c, "New terminal", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			id := fmt.Sprint(st.next)
			st.next++
			st.tabs = append(st.tabs, code.TerminalTab{ID: id, Title: "shell " + id})
			st.active = id
			log("new " + id)
		}
	})
	r := code.TerminalTabs(c, &st.tabs, &st.active, code.TerminalTabsOptions{Disabled: st.disabled}, func(id string) {
		ui.Text(c, "Terminal "+id+" content").TextColor(core.Tokens(c).TextMuted)
	})
	r.Element.Width(560)
	if r.Changed {
		log("active " + st.active)
	}
	if r.Renamed != "" {
		log("renamed " + r.Renamed)
	}
	if r.Closed != "" {
		log("closed " + r.Closed)
	}
}

func codeTerminalSearchDemo(c *ui.Context, log func(string)) {
	type state struct {
		term     code.TerminalState
		search   code.TerminalSearchState
		disabled bool
	}
	s := ui.Local(c.Root(), "323", func() *state {
		s := &state{}
		for i := range 60 {
			if i%9 == 4 {
				s.term.Write(fmt.Sprintf("[%02d] \x1b[31merror\x1b[0m: seal %d missing\r\n", i, i))
			} else {
				s.term.Write(fmt.Sprintf("[%02d] ledger page %d copied\r\n", i, i))
			}
		}
		return s
	})
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	ui.Column(c).Width(600).Gap(8).Children(func() {
		if code.TerminalSearch(c, &st.term, &st.search, code.TerminalSearchOptions{Disabled: st.disabled}).Changed() {
			log(fmt.Sprintf("query %q, match %d", st.search.Query, st.search.Current+1))
		}
		code.Terminal(c, &st.term, code.TerminalOptions{}).Height(220)
	})
}

func codeCommandHistoryDemo(c *ui.Context, log func(string)) {
	type state struct {
		entries  []code.HistoryEntry
		s        code.CommandHistoryState
		disabled bool
	}
	s := ui.Local(c.Root(), "324", func() *state {
		at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
		var e []code.HistoryEntry
		for i, cmd := range []string{"git pull", "go build ./...", "go test ./...", "make seal", "ls -la", "git status"} {
			e = append(e, code.HistoryEntry{Command: cmd, Time: at.Add(time.Duration(i) * time.Minute)})
		}
		e[3].ExitCode = 2
		return &state{entries: e}
	})
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	r := code.CommandHistory(c, st.entries, &st.s, code.CommandHistoryOptions{Disabled: st.disabled})
	r.Element.Size(480, 300)
	if r.Replayed != nil {
		log("replay " + r.Replayed.Command)
		st.entries = append(st.entries, code.HistoryEntry{Command: r.Replayed.Command, Time: time.Now()})
	}
}

func codeLogViewerDemo(c *ui.Context, log func(string)) {
	type state struct {
		lines    []code.LogLine
		s        code.LogViewerState
		disabled bool
	}
	s := ui.Local(c.Root(), "325", func() *state {
		at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
		lines := make([]code.LogLine, 100_000)
		words := []string{"ledger synced", "seal applied", "herald dispatched", "archive rotated"}
		for i := range lines {
			lines[i] = code.LogLine{Time: at.Add(time.Duration(i) * time.Second), Level: code.LogLevel(i % 7 % 4), Text: fmt.Sprintf("%s #%d", words[i%4], i)}
		}
		return &state{lines: lines}
	})
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		if input.Button(c, "Append line", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			st.lines = append(st.lines, code.LogLine{Time: time.Now(), Level: code.LogWarn, Text: fmt.Sprintf("late seal #%d", len(st.lines))})
			log("append")
		}
	})
	r := code.LogViewer(c, st.lines, &st.s, code.LogViewerOptions{Disabled: st.disabled})
	r.Element.Size(720, 360)
	if r.Changed() {
		log(fmt.Sprintf("hidden %v, follow %v", st.s.Hidden, st.s.Follow))
	}
}
