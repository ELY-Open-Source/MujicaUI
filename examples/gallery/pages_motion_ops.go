package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func motionInterruptDemo(c *ui.Context, log func(string)) {
	type state struct {
		running, stopping, confirm bool
		until                      time.Time
	}
	s := ui.Local(c.Root(), "281", func() state { return state{running: true, confirm: true} })
	if s.stopping && c.Now().After(s.until) {
		s.running, s.stopping = false, false
		log("Run stopped")
	}
	if s.stopping {
		c.After(s.until.Sub(c.Now()))
	}
	navigationDemo(c, func() {
		navigationCheck(c, &s.confirm, "Confirm", log)
		if mujica.Button(c, "Start run", mujica.ButtonOptions{Variant: mujica.Secondary, Disabled: s.running}).Clicked() {
			s.running = true
			log("Run started")
		}
	}, func() {
		ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
			if s.running && !s.stopping {
				mujica.DotsLoader(c, mujica.LoaderOptions{Label: "Agent running", ShowLabel: true})
			}
			if mujica.InterruptButton(c, s.running, mujica.InterruptButtonOptions{Confirm: s.confirm, Interrupting: s.stopping}).Changed() {
				log("Interrupt requested")
				s.stopping, s.until = true, c.Now().Add(1500*time.Millisecond)
			}
		})
	})
}

func motionTokenUsageDemo(c *ui.Context, log func(string)) {
	type state struct{ line, loading, empty, failed bool }
	s := ui.Local(c.Root(), "282", func() state { return state{} })
	usage := *ui.Local(c.Root(), "282.data", func() []mujica.TokenUsage {
		return []mujica.TokenUsage{{Label: "Plan", Input: 1200, Output: 300}, {Label: "Search", Input: 4800, Output: 650},
			{Label: "Read", Input: 6100, Output: 420}, {Label: "Write", Input: 2100, Output: 2900}, {Label: "Review", Input: 3300, Output: 800}}
	})
	navigationDemo(c, func() {
		navigationCheck(c, &s.line, "Lines", log)
		navigationCheck(c, &s.loading, "Loading", log)
		navigationCheck(c, &s.empty, "Empty", log)
		navigationCheck(c, &s.failed, "Error", log)
	}, func() {
		st := mujica.DataReady
		if s.loading {
			st = mujica.DataLoading
		}
		if s.failed {
			st = mujica.DataFailed
		}
		if s.empty {
			usage = nil
		}
		r := mujica.TokenUsageChart(c, usage, mujica.TokenUsageChartOptions{Line: s.line, Frame: mujica.ChartFrame{Status: st, Error: "Usage unavailable"}})
		if r.Changed() && r.Point >= 0 {
			log(fmt.Sprintf("Readout: mark %d", r.Point))
		}
		if name, ok := r.Toggled(); ok {
			log("Toggled " + name)
		}
	})
}

func motionCostDemo(c *ui.Context, log func(string)) {
	type state struct {
		table mujica.CostBreakdownState
		by    int
	}
	s := ui.Local(c.Root(), "283", func() state { return state{} })
	items := *ui.Local(c.Root(), "283.data", func() []mujica.CostItem {
		return []mujica.CostItem{
			{Model: "opus", Tool: "search", Step: "plan", InputTokens: 1000, OutputTokens: 200, Cost: 0.30},
			{Model: "haiku", Tool: "search", Step: "gather", InputTokens: 4000, OutputTokens: 500, Cost: 0.05},
			{Model: "opus", Tool: "write", Step: "write", InputTokens: 2000, OutputTokens: 3000, Cost: 1.15},
			{Model: "sonnet", Tool: "read", Step: "gather", InputTokens: 6000, OutputTokens: 400, Cost: 0.24},
		}
	})
	navigationDemo(c, func() {
		if mujica.SegmentedControl(c, &s.by, []mujica.Segment{{Label: "Model"}, {Label: "Tool"}, {Label: "Step"}}, mujica.SegmentedControlOptions{Label: "Group by"}).Changed() {
			log(fmt.Sprintf("Group by %d", s.by))
		}
	}, func() {
		v := mujica.CostBreakdown(c, &s.table, items, mujica.CostBreakdownOptions{GroupBy: mujica.CostGroup(s.by)})
		v.Element.Height(320)
		if v.SortChanged() {
			log("Sorted by " + s.table.Table.Sort.Column)
		}
		if key, ok := s.table.Table.Selected(); ok && v.Changed() {
			log("Selected " + key)
		}
	})
}

func motionTraceDemo(c *ui.Context, log func(string)) {
	type state struct {
		viewer                           mujica.TraceViewerState
		empty, loading, failed, disabled bool
	}
	s := ui.Local(c.Root(), "284", func() state { return state{} })
	spans := *ui.Local(c.Root(), "284.data", func() []mujica.TraceSpan {
		return []mujica.TraceSpan{
			{ID: "run", Name: "agent.run", End: 900 * time.Millisecond, Attrs: []mujica.TraceAttr{{Key: "model", Value: "opus"}, {Key: "session", Value: "s-42"}}},
			{ID: "plan", Parent: "run", Name: "llm.plan", Start: 10 * time.Millisecond, End: 210 * time.Millisecond, Attrs: []mujica.TraceAttr{{Key: "tokens", Value: "1200"}}},
			{ID: "search", Parent: "run", Name: "tool.search", Start: 220 * time.Millisecond, End: 520 * time.Millisecond, Failed: true, Attrs: []mujica.TraceAttr{{Key: "error", Value: "timeout"}}},
			{ID: "fetch", Parent: "search", Name: "http.get", Start: 230 * time.Millisecond, End: 500 * time.Millisecond},
			{ID: "write", Parent: "run", Name: "llm.write", Start: 530 * time.Millisecond, End: 880 * time.Millisecond},
		}
	})
	navigationDemo(c, func() {
		navigationCheck(c, &s.empty, "Empty", log)
		navigationCheck(c, &s.loading, "Loading", log)
		navigationCheck(c, &s.failed, "Error", log)
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		status := mujica.DataReady
		if s.empty || s.loading || s.failed {
			spans = nil
		}
		if s.loading {
			status = mujica.DataLoading
		}
		if s.failed {
			status = mujica.DataFailed
		}
		ui.Column(c).Disabled(s.disabled).Children(func() {
			v := mujica.TraceViewer(c, &s.viewer, spans, mujica.TraceViewerOptions{Status: status, Error: "Trace unavailable"})
			v.Element.Height(440)
			if id, ok := s.viewer.Tree.Selected(); ok && v.Changed() {
				log("Selected " + id)
			}
			if v.SortChanged() {
				log("Sorted by " + s.viewer.Tree.Sort.Column)
			}
			if v.Retried() {
				log("Retry trace requested")
			}
		})
	})
}

func motionEvalDemo(c *ui.Context, log func(string)) {
	type state struct {
		table                            mujica.DataTableState[string]
		empty, loading, failed, disabled bool
	}
	s := ui.Local(c.Root(), "285", func() state { return state{} })
	cases := *ui.Local(c.Root(), "285.data", func() []mujica.EvalCase {
		return []mujica.EvalCase{
			{ID: "c1", Name: "Greets by title", Score: 0.98, Passed: true},
			{ID: "c2", Name: "Refuses forgery", Score: 0.41},
			{ID: "c3", Name: "Summarizes letter", Score: 0.87, Passed: true},
			{ID: "c4", Name: "Keeps the schedule", Score: 0.62},
		}
	})
	navigationDemo(c, func() {
		navigationCheck(c, &s.empty, "Empty", log)
		navigationCheck(c, &s.loading, "Loading", log)
		navigationCheck(c, &s.failed, "Error", log)
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		status := mujica.DataReady
		if s.empty || s.loading || s.failed {
			cases = nil
		}
		if s.loading {
			status = mujica.DataLoading
		}
		if s.failed {
			status = mujica.DataFailed
		}
		ui.Column(c).Disabled(s.disabled).Children(func() {
			r := mujica.EvalResultTable(c, &s.table, cases, mujica.EvalResultTableOptions{Status: status, Error: "Evaluation unavailable"})
			r.Element.Height(320)
			if id, ok := r.Diff(); ok {
				log("Diff " + id)
			}
			if r.Table.Changed() {
				id, _ := s.table.Selected()
				log("Selected " + id)
			}
			if r.Table.SortChanged() {
				log("Sorted by " + s.table.Sort.Column)
			}
			if r.Table.Retried() {
				log("Retry evaluation requested")
			}
		})
	})
}

func motionToolsDemo(c *ui.Context, log func(string)) {
	type state struct {
		list            mujica.ListState[string]
		tools           []mujica.AgentTool
		disabled, empty bool
	}
	s := ui.Local(c.Root(), "286", func() state {
		return state{tools: []mujica.AgentTool{
			{Name: "web_search", Description: "Search the web", Enabled: true, Params: []mujica.ToolParam{
				{Name: "query", Type: "string", Required: true, Description: "What to look for"}, {Name: "limit", Type: "integer"}}},
			{Name: "read_file", Description: "Read a local file", Params: []mujica.ToolParam{{Name: "path", Type: "string", Required: true}}},
			{Name: "run_shell", Description: "Run a shell command", Params: []mujica.ToolParam{{Name: "command", Type: "string", Required: true}, {Name: "timeout", Type: "number"}}},
		}}
	})
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
		navigationCheck(c, &s.empty, "Empty", log)
	}, func() {
		tools := s.tools
		if s.empty {
			tools = nil
		}
		r := mujica.ToolRegistryPanel(c, &s.list, tools, mujica.ToolRegistryPanelOptions{Disabled: s.disabled})
		r.Element.Height(320)
		if r.List.Changed() {
			name, _ := s.list.Selected()
			log("Selected " + name)
		}
		if name, on, ok := r.Toggled(); ok {
			log(fmt.Sprintf("%s enabled=%v", name, on))
			for i := range s.tools {
				if s.tools[i].Name == name {
					s.tools[i].Enabled = on
				}
			}
		}
	})
}
