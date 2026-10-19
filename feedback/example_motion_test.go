package feedback_test

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/feedback"
)

func ExampleInterruptButton() {
	running, stopping := true, false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := feedback.InterruptButton(c, running, feedback.InterruptButtonOptions{Confirm: true, Interrupting: stopping})
		if r.Changed() {
			stopping = true // cancel the run, then clear running and stopping
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTokenUsageChart() {
	usage := []feedback.TokenUsage{{Label: "Plan", Input: 1200, Output: 300}, {Label: "Write", Input: 2100, Output: 2900}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.TokenUsageChart(c, usage, feedback.TokenUsageChartOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCostBreakdown() {
	var state feedback.CostBreakdownState
	items := []feedback.CostItem{
		{Model: "opus", Tool: "search", InputTokens: 1000, OutputTokens: 200, Cost: 0.30},
		{Model: "haiku", Tool: "search", InputTokens: 4000, OutputTokens: 500, Cost: 0.05},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.CostBreakdown(c, &state, items, feedback.CostBreakdownOptions{GroupBy: feedback.CostByModel}).Element.Height(260)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTraceViewer() {
	var state feedback.TraceViewerState
	spans := []feedback.TraceSpan{
		{ID: "run", Name: "agent.run", End: 900 * time.Millisecond},
		{ID: "plan", Parent: "run", Name: "llm.plan", Start: 10 * time.Millisecond, End: 210 * time.Millisecond,
			Attrs: []feedback.TraceAttr{{Key: "model", Value: "opus"}}},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.TraceViewer(c, &state, spans, feedback.TraceViewerOptions{}).Element.Height(260)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEvalResultTable() {
	var state data.DataTableState[string]
	cases := []feedback.EvalCase{{ID: "c1", Name: "Greets by title", Score: 0.98, Passed: true}, {ID: "c2", Name: "Refuses forgery", Score: 0.41}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := feedback.EvalResultTable(c, &state, cases, feedback.EvalResultTableOptions{})
		r.Table.Element.Height(260)
		if id, ok := r.Diff(); ok {
			_ = id // open expected and actual output side by side
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToolRegistryPanel() {
	var state data.ListState[string]
	tools := []feedback.AgentTool{{Name: "web_search", Description: "Search the web", Enabled: true,
		Params: []feedback.ToolParam{{Name: "query", Type: "string", Required: true}}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := feedback.ToolRegistryPanel(c, &state, tools, feedback.ToolRegistryPanelOptions{})
		r.Element.Height(260)
		if name, on, ok := r.Toggled(); ok && name == tools[0].Name {
			tools[0].Enabled = on
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDotsLoader() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.DotsLoader(c, feedback.LoaderOptions{Label: "Thinking", ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBarsLoader() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.BarsLoader(c, feedback.LoaderOptions{Size: 20, Label: "Indexing"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePulseLoader() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.PulseLoader(c, feedback.LoaderOptions{Label: "Listening", ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOrbitLoader() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.OrbitLoader(c, feedback.LoaderOptions{Size: 24, Label: "Connecting"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleWaveLoader() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.WaveLoader(c, feedback.LoaderOptions{Label: "Speaking", ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleShimmer() {
	loading := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.Shimmer(c, loading, feedback.ShimmerOptions{Height: 48}, func() {
			ui.Text(c, "Report ready")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCountUp() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.CountUp(c, 12840, feedback.CountUpOptions{Duration: 800 * time.Millisecond})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTypewriter() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := feedback.Typewriter(c, "The court assembles at dusk.", feedback.TypewriterOptions{Speed: 30, Skippable: true})
		if r.Changed() {
			// show the next message
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleShimmerText() {
	thinking := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.ShimmerText(c, "Thinking…", feedback.ShimmerTextOptions{Active: thinking})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStagger() {
	guests := []string{"Sakiko", "Uika", "Mutsumi"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.Stagger(c, len(guests), feedback.StaggerOptions{}, func(i int) {
			ui.Text(c, guests[i])
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleReorderTransition() {
	order := []string{"Sakiko", "Uika", "Mutsumi"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.ReorderTransition(c, order, feedback.ReorderTransitionOptions{Gap: 4}, func(name string) {
			ui.Text(c, name)
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLayoutTransition() {
	open := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.LayoutTransition(c, feedback.LayoutTransitionOptions{Clip: true}, func() {
			ui.Text(c, "Program")
			if open {
				ui.Text(c, "Overture · Masquerade · Finale")
			}
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePresence() {
	show := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.Presence(c, show, feedback.PresenceOptions{Kind: feedback.PresenceSlide}, func() {
			ui.Text(c, "Saved")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBlinkHighlight() {
	price := 101.5
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.BlinkHighlight(c, price, feedback.BlinkHighlightOptions{})
	}
	ui.Render(view, 400, 300, 1)
}
