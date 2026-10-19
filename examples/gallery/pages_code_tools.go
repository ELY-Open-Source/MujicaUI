package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/code"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func codeProcessListDemo(c *ui.Context, log func(string)) {
	type state struct {
		procs    []code.Process
		s        code.ProcessListState
		disabled bool
	}
	s := ui.Local(c.Root(), "326", func() *state {
		return &state{procs: []code.Process{
			{PID: 101, Name: "chancellor", CPU: 3.5, Memory: 200 << 20}, {PID: 7, Name: "archivist", CPU: 42.1, Memory: 50 << 20},
			{PID: 330, Name: "herald", CPU: 0.2, Memory: 2 << 30}, {PID: 412, Name: "scribe", CPU: 12.8, Memory: 96 << 20},
		}}
	})
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	r := code.ProcessList(c, st.procs, &st.s, code.ProcessListOptions{Disabled: st.disabled})
	r.Element.Size(560, 260)
	if r.End != nil {
		log(fmt.Sprintf("end %s (PID %d)", r.End.Name, r.End.PID))
		var kept []code.Process
		for _, p := range st.procs {
			if p.PID != r.End.PID {
				kept = append(kept, p)
			}
		}
		st.procs = kept
	}
}

var codeDemoSource = strings.Split(`package court

import "errors"

// Decree is an order of the court.
type Decree struct {
	ID    int
	Title string
}

/* Seal stamps a decree
   with the royal wax. */
func Seal(d Decree) (string, error) {
	if d.Title == "" {
		return "", errors.New("untitled decree")
	}
	for i := 0; i < 3; i++ {
		d.ID += i
	}
	return "sealed: " + d.Title, nil
}

func Herald(ds []Decree) int {
	n := 0
	for _, d := range ds {
		if _, err := Seal(d); err == nil {
			n++
		}
	}
	return n
}`, "\n")

func codeCodeViewerDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        code.CodeViewerState
		lang     string
		disabled bool
	}
	s := ui.Local(c.Root(), "327", func() *state { return &state{lang: "go"} })
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		langs := []input.SelectOption[string]{{Value: "go", Label: "Go"}, {Value: "", Label: "Plain"}}
		if input.Select(c, &st.lang, langs, input.SelectOptions{Label: "Language"}).Changed() {
			log("language " + st.lang)
		}
		if input.Button(c, "Go to line 20", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			st.s.GoTo(19)
			log("go to 20")
		}
	})
	r := code.CodeViewer(c, codeDemoSource, &st.s, code.CodeViewerOptions{Language: st.lang, Label: "Source", Disabled: st.disabled})
	r.Element.Size(600, 360)
	if r.Changed() {
		log(fmt.Sprintf("line %d", st.s.Line()+1))
	}
}

func codeMinimapDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        code.CodeViewerState
		view     code.MinimapView
		lines    []string
		disabled bool
	}
	s := ui.Local(c.Root(), "328", func() *state {
		var lines []string
		for range 12 {
			lines = append(lines, codeDemoSource...)
		}
		return &state{lines: lines}
	})
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	first, last := st.s.Visible()
	st.view = code.MinimapView{Top: first, Visible: max(last-first+1, 1)}
	ui.Row(c).Gap(8).Height(360).AlignItems(ui.Stretch).Children(func() {
		code.CodeViewer(c, st.lines, &st.s, code.CodeViewerOptions{Language: "go"}).Element.Width(520)
		if code.Minimap(c, st.lines, &st.view, code.MinimapOptions{Language: "go", Disabled: st.disabled}).Changed() {
			st.s.ScrollTo(st.view.Top)
			log(fmt.Sprintf("jump to line %d", st.view.Top+1))
		}
	})
}

func codeFindWidgetDemo(c *ui.Context, log func(string)) {
	type state struct {
		text               string
		s                  code.FindState
		readOnly, disabled bool
	}
	s := ui.Local(c.Root(), "329", func() *state {
		return &state{text: "The seal of the court; seals are kept sealed. Seal 42 and seal 7.", s: code.FindState{Query: "seal"}}
	})
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.readOnly, "Read-only", log)
		codeDemoCheck(c, &st.disabled, "Disabled", log)
	})
	ui.Column(c).Width(620).Gap(12).Children(func() {
		r := code.FindWidget(c, &st.text, &st.s, code.FindWidgetOptions{ReadOnly: st.readOnly, Disabled: st.disabled})
		if r.Moved {
			log(fmt.Sprintf("%d matches, at %d", len(st.s.Matches()), st.s.Current+1))
		}
		if r.Replaced > 0 {
			log(fmt.Sprintf("replaced %d", r.Replaced))
		}
		var spans []ui.Span
		at := 0
		for i, m := range st.s.Matches() {
			spans = append(spans, ui.Span{Text: st.text[at:m[0]]})
			sp := ui.Span{Text: st.text[m[0]:m[1]], Background: core.Tokens(c).Selection, Underline: true}
			if i == st.s.Current {
				sp.Weight = 700
			}
			spans = append(spans, sp)
			at = m[1]
		}
		spans = append(spans, ui.Span{Text: st.text[at:]})
		ui.RichText(c, spans...)
	})
}

func codeSearchPanelDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        code.SearchPanelState
		results  []code.SearchFileResult
		disabled bool
	}
	s := ui.Local(c.Root(), "330", func() *state { return &state{s: code.SearchPanelState{Query: "seal"}} })
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	// The demo searches its own two files.
	files := map[string][]string{"court/seal.go": codeDemoSource, "docs/ledger.md": {"# Ledger", "The seal is kept by the chancellor.", "Seals expire after a year."}}
	if st.results == nil || st.s.Query != "" && len(st.results) == 0 {
		st.results = codeDemoSearch(files, st.s.Query, st.s.Include, st.s.Exclude)
	}
	r := code.SearchPanel(c, st.results, &st.s, code.SearchPanelOptions{Disabled: st.disabled})
	r.Element.Size(560, 380)
	if r.Changed {
		st.results = codeDemoSearch(files, st.s.Query, st.s.Include, st.s.Exclude)
		log(fmt.Sprintf("search %q include %q exclude %q", st.s.Query, st.s.Include, st.s.Exclude))
	}
	if r.Opened != nil {
		log(fmt.Sprintf("open %s:%d", r.Opened.Path, r.Opened.Line))
	}
}

// codeDemoSearch finds q in the demo files.
func codeDemoSearch(files map[string][]string, q, include, exclude string) []code.SearchFileResult {
	out := []code.SearchFileResult{}
	if q == "" {
		return out
	}
	for _, path := range []string{"court/seal.go", "docs/ledger.md"} {
		if include != "" && !strings.Contains(path, strings.Trim(include, "*")) || exclude != "" && strings.Contains(path, strings.Trim(exclude, "*")) {
			continue
		}
		f := code.SearchFileResult{Path: path}
		for i, line := range files[path] {
			if at := strings.Index(strings.ToLower(line), strings.ToLower(q)); at >= 0 {
				f.Matches = append(f.Matches, code.SearchMatch{Line: i + 1, Text: line, Start: at, End: at + len(q)})
			}
		}
		if len(f.Matches) > 0 {
			out = append(out, f)
		}
	}
	return out
}

func codeProblemsDemo(c *ui.Context, log func(string)) {
	type state struct {
		s        code.ProblemsPanelState
		disabled bool
	}
	s := ui.Local(c.Root(), "331", func() *state { return &state{} })
	st := *s
	codeDemoControls(c, func() { codeDemoCheck(c, &st.disabled, "Disabled", log) })
	problems := []code.Problem{
		{Severity: code.ProblemError, File: "court/seal.go", Line: 12, Col: 5, Message: "undefined: Decree", Source: "compiler"},
		{Severity: code.ProblemWarning, File: "court/seal.go", Line: 30, Col: 2, Message: "unused variable wax", Source: "vet"},
		{Severity: code.ProblemInfo, File: "docs/ledger.md", Line: 4, Col: 1, Message: "line is longer than 80", Source: "lint"},
		{Severity: code.ProblemError, File: "main.go", Line: 9, Col: 14, Message: "missing return", Source: "compiler"},
	}
	r := code.ProblemsPanel(c, problems, &st.s, code.ProblemsPanelOptions{Disabled: st.disabled})
	r.Element.Size(620, 320)
	if r.Opened != nil {
		log(fmt.Sprintf("go to %s:%d:%d", r.Opened.File, r.Opened.Line, r.Opened.Col))
	}
}

func codeOutputDemo(c *ui.Context, log func(string)) {
	type state struct {
		channels []code.OutputChannel
		s        code.OutputPanelState
		disabled bool
	}
	s := ui.Local(c.Root(), "332", func() *state {
		var build []string
		for i := range 300 {
			build = append(build, fmt.Sprintf("compiling unit %d of 300", i+1))
		}
		return &state{channels: []code.OutputChannel{{Name: "Build", Lines: build},
			{Name: "Tests", Lines: []string{"\x1b[32mPASS\x1b[0m court/seal", "\x1b[31mFAIL\x1b[0m court/ledger: expected 3 seals", "\x1b[1mok\x1b[0m  2 of 3 packages"}}}}
	})
	st := *s
	codeDemoControls(c, func() {
		codeDemoCheck(c, &st.disabled, "Disabled", log)
		if input.Button(c, "Write line", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			ch := &st.channels[st.s.Channel]
			ch.Lines = append(ch.Lines, fmt.Sprintf("line %d", len(ch.Lines)+1))
			log("write to " + ch.Name)
		}
	})
	r := code.OutputPanel(c, st.channels, &st.s, code.OutputPanelOptions{Disabled: st.disabled})
	r.Element.Size(560, 300)
	if r.Changed {
		log("channel " + st.channels[st.s.Channel].Name)
	}
	if r.Cleared != "" {
		st.channels[st.s.Channel].Lines = nil
		log("cleared " + r.Cleared)
	}
}
