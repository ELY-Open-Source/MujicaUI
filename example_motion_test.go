package mujica_test

import (
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleInterruptButton() {
	running, stopping := true, false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.InterruptButton(c, running, mujica.InterruptButtonOptions{Confirm: true, Interrupting: stopping})
		if r.Changed() {
			stopping = true // cancel the run, then clear running and stopping
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTokenUsageChart() {
	usage := []mujica.TokenUsage{{Label: "Plan", Input: 1200, Output: 300}, {Label: "Write", Input: 2100, Output: 2900}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TokenUsageChart(c, usage, mujica.TokenUsageChartOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCostBreakdown() {
	var state mujica.CostBreakdownState
	items := []mujica.CostItem{
		{Model: "opus", Tool: "search", InputTokens: 1000, OutputTokens: 200, Cost: 0.30},
		{Model: "haiku", Tool: "search", InputTokens: 4000, OutputTokens: 500, Cost: 0.05},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CostBreakdown(c, &state, items, mujica.CostBreakdownOptions{GroupBy: mujica.CostByModel}).Element.Height(260)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTraceViewer() {
	var state mujica.TraceViewerState
	spans := []mujica.TraceSpan{
		{ID: "run", Name: "agent.run", End: 900 * time.Millisecond},
		{ID: "plan", Parent: "run", Name: "llm.plan", Start: 10 * time.Millisecond, End: 210 * time.Millisecond,
			Attrs: []mujica.TraceAttr{{Key: "model", Value: "opus"}}},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TraceViewer(c, &state, spans, mujica.TraceViewerOptions{}).Element.Height(260)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEvalResultTable() {
	var state mujica.DataTableState[string]
	cases := []mujica.EvalCase{{ID: "c1", Name: "Greets by title", Score: 0.98, Passed: true}, {ID: "c2", Name: "Refuses forgery", Score: 0.41}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.EvalResultTable(c, &state, cases, mujica.EvalResultTableOptions{})
		r.Table.Element.Height(260)
		if id, ok := r.Diff(); ok {
			_ = id // open expected and actual output side by side
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToolRegistryPanel() {
	var state mujica.ListState[string]
	tools := []mujica.AgentTool{{Name: "web_search", Description: "Search the web", Enabled: true,
		Params: []mujica.ToolParam{{Name: "query", Type: "string", Required: true}}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.ToolRegistryPanel(c, &state, tools, mujica.ToolRegistryPanelOptions{})
		r.Element.Height(260)
		if name, on, ok := r.Toggled(); ok && name == tools[0].Name {
			tools[0].Enabled = on
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDotsLoader() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DotsLoader(c, mujica.LoaderOptions{Label: "Thinking", ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBarsLoader() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.BarsLoader(c, mujica.LoaderOptions{Size: 20, Label: "Indexing"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePulseLoader() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.PulseLoader(c, mujica.LoaderOptions{Label: "Listening", ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOrbitLoader() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.OrbitLoader(c, mujica.LoaderOptions{Size: 24, Label: "Connecting"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleWaveLoader() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.WaveLoader(c, mujica.LoaderOptions{Label: "Speaking", ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleShimmer() {
	loading := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Shimmer(c, loading, mujica.ShimmerOptions{Height: 48}, func() {
			ui.Text(c, "Report ready")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCountUp() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CountUp(c, 12840, mujica.CountUpOptions{Duration: 800 * time.Millisecond})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTypewriter() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.Typewriter(c, "The court assembles at dusk.", mujica.TypewriterOptions{Speed: 30, Skippable: true})
		if r.Changed() {
			// show the next message
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleShimmerText() {
	thinking := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ShimmerText(c, "Thinking…", mujica.ShimmerTextOptions{Active: thinking})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStagger() {
	guests := []string{"Sakiko", "Uika", "Mutsumi"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Stagger(c, len(guests), mujica.StaggerOptions{}, func(i int) {
			ui.Text(c, guests[i])
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleReorderTransition() {
	order := []string{"Sakiko", "Uika", "Mutsumi"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ReorderTransition(c, order, mujica.ReorderTransitionOptions{Gap: 4}, func(name string) {
			ui.Text(c, name)
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLayoutTransition() {
	open := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.LayoutTransition(c, mujica.LayoutTransitionOptions{Clip: true}, func() {
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
		mujica.Use(c, mujica.Settings{})
		mujica.Presence(c, show, mujica.PresenceOptions{Kind: mujica.PresenceSlide}, func() {
			ui.Text(c, "Saved")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBlinkHighlight() {
	price := 101.5
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.BlinkHighlight(c, price, mujica.BlinkHighlightOptions{})
	}
	ui.Render(view, 400, 300, 1)
}
