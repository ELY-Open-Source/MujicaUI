package main

import (
	"github.com/ELY-Open-Source/MujicaUI/feedback"
)

var motionPages = []page{
	{
		ID: "281", Name: "InterruptButton", ZhName: "中断按钮 停止智能体", Category: "Agent Ops & Motion",
		Demo: motionInterruptDemo,
		Code: `running, stopping := true, false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := feedback.InterruptButton(c, running, feedback.InterruptButtonOptions{Confirm: true, Interrupting: stopping})
	if r.Changed() {
		stopping = true // cancel the run, then clear running and stopping
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab focuses Stop; Enter or Space interrupts. With Confirm, the popup focuses Keep running; Tab reaches Stop; Escape cancels.",
	},
	{
		ID: "282", Name: "TokenUsageChart", ZhName: "Token 用量图 输入输出", Category: "Agent Ops & Motion",
		Demo: motionTokenUsageDemo,
		Code: `usage := []feedback.TokenUsage{{Label: "Plan", Input: 1200, Output: 300}, {Label: "Write", Input: 2100, Output: 2900}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.TokenUsageChart(c, usage, feedback.TokenUsageChartOptions{})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab focuses the chart; Left/Right move between bars or points and the readout is announced.",
	},
	{
		ID: "283", Name: "CostBreakdown", ZhName: "费用分解 成本 合计", Category: "Agent Ops & Motion",
		Demo: motionCostDemo,
		Code: `var state feedback.CostBreakdownState
items := []feedback.CostItem{
	{Model: "opus", Tool: "search", InputTokens: 1000, OutputTokens: 200, Cost: 0.30},
	{Model: "haiku", Tool: "search", InputTokens: 4000, OutputTokens: 500, Cost: 0.05},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.CostBreakdown(c, &state, items, feedback.CostBreakdownOptions{GroupBy: feedback.CostByModel}).Element.Height(260)
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a group; column headers sort with Enter; the column menu takes Tab.",
	},
	{
		ID: "284", Name: "TraceViewer", ZhName: "调用链 追踪 跨度 时间轴", Category: "Agent Ops & Motion",
		Demo: motionTraceDemo,
		Code: `var state feedback.TraceViewerState
spans := []feedback.TraceSpan{
	{ID: "run", Name: "agent.run", End: 900 * time.Millisecond},
	{ID: "plan", Parent: "run", Name: "llm.plan", Start: 10 * time.Millisecond, End: 210 * time.Millisecond,
		Attrs: []feedback.TraceAttr{{Key: "model", Value: "opus"}}},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.TraceViewer(c, &state, spans, feedback.TraceViewerOptions{}).Element.Height(260)
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a span, Right opens and Left closes it; typing jumps by name.",
	},
	{
		ID: "285", Name: "EvalResultTable", ZhName: "评测结果 得分 通过失败", Category: "Agent Ops & Motion",
		Demo: motionEvalDemo,
		Code: `var state data.DataTableState[string]
cases := []feedback.EvalCase{{ID: "c1", Name: "Greets by title", Score: 0.98, Passed: true}, {ID: "c2", Name: "Refuses forgery", Score: 0.41}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := feedback.EvalResultTable(c, &state, cases, feedback.EvalResultTableOptions{})
	r.Table.Element.Height(260)
	if id, ok := r.Diff(); ok {
		_ = id // open expected and actual output side by side
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a case; the diff button takes Tab, Enter or Space.",
	},
	{
		ID: "286", Name: "ToolRegistryPanel", ZhName: "工具注册 工具列表 启停", Category: "Agent Ops & Motion",
		Demo: motionToolsDemo,
		Code: `var state data.ListState[string]
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
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a tool; its switch takes Tab and Space.",
	},
	{
		ID: "287", Name: "DotsLoader", ZhName: "三点加载 加载中", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("287", feedback.DotsLoader),
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.DotsLoader(c, feedback.LoaderOptions{Label: "Thinking", ShowLabel: true})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "288", Name: "BarsLoader", ZhName: "条形加载", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("288", feedback.BarsLoader),
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.BarsLoader(c, feedback.LoaderOptions{Size: 20, Label: "Indexing"})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "289", Name: "PulseLoader", ZhName: "呼吸加载 脉冲", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("289", feedback.PulseLoader),
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.PulseLoader(c, feedback.LoaderOptions{Label: "Listening", ShowLabel: true})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "290", Name: "OrbitLoader", ZhName: "环绕加载 菱形", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("290", feedback.OrbitLoader),
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.OrbitLoader(c, feedback.LoaderOptions{Size: 24, Label: "Connecting"})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "291", Name: "WaveLoader", ZhName: "波形加载", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("291", feedback.WaveLoader),
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.WaveLoader(c, feedback.LoaderOptions{Label: "Speaking", ShowLabel: true})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "292", Name: "Shimmer", ZhName: "占位扫光 骨架 加载", Category: "Agent Ops & Motion",
		Demo: motionShimmerDemo,
		Code: `loading := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.Shimmer(c, loading, feedback.ShimmerOptions{Height: 48}, func() {
		ui.Text(c, "Report ready")
	})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress while loading.",
	},
	{
		ID: "293", Name: "CountUp", ZhName: "数字过渡 计数", Category: "Agent Ops & Motion",
		Demo: motionCountUpDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.CountUp(c, 12840, feedback.CountUpOptions{Duration: 800 * time.Millisecond})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; the final value is read out.",
	},
	{
		ID: "294", Name: "Typewriter", ZhName: "打字机 逐字显示", Category: "Agent Ops & Motion",
		Demo: motionTypewriterDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := feedback.Typewriter(c, "The court assembles at dusk.", feedback.TypewriterOptions{Speed: 30, Skippable: true})
	if r.Changed() {
		// show the next message
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Skip takes Tab, Enter or Space.",
	},
	{
		ID: "295", Name: "ShimmerText", ZhName: "文字扫光 思考中", Category: "Agent Ops & Motion",
		Demo: motionShimmerTextDemo,
		Code: `thinking := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.ShimmerText(c, "Thinking…", feedback.ShimmerTextOptions{Active: thinking})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as a status.",
	},
	{
		ID: "296", Name: "Stagger", ZhName: "依次进场 交错动画", Category: "Agent Ops & Motion",
		Demo: motionStaggerDemo,
		Code: `guests := []string{"Sakiko", "Uika", "Mutsumi"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.Stagger(c, len(guests), feedback.StaggerOptions{}, func(i int) {
		ui.Text(c, guests[i])
	})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable.",
	},
	{
		ID: "297", Name: "ReorderTransition", ZhName: "重排过渡 排序动画", Category: "Agent Ops & Motion",
		Demo: motionReorderDemo,
		Code: `order := []string{"Sakiko", "Uika", "Mutsumi"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.ReorderTransition(c, order, feedback.ReorderTransitionOptions{Gap: 4}, func(name string) {
		ui.Text(c, name)
	})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Focus inside an item follows it as it moves.",
	},
	{
		ID: "298", Name: "LayoutTransition", ZhName: "布局过渡 FLIP", Category: "Agent Ops & Motion",
		Demo: motionLayoutDemo,
		Code: `open := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.LayoutTransition(c, feedback.LayoutTransitionOptions{Clip: true}, func() {
		ui.Text(c, "Program")
		if open {
			ui.Text(c, "Overture · Masquerade · Finale")
		}
	})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable.",
	},
	{
		ID: "299", Name: "Presence", ZhName: "进出场 淡入淡出", Category: "Agent Ops & Motion",
		Demo: motionPresenceDemo,
		Code: `show := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.Presence(c, show, feedback.PresenceOptions{Kind: feedback.PresenceSlide}, func() {
		ui.Text(c, "Saved")
	})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Leaving content takes no focus.",
	},
	{
		ID: "300", Name: "BlinkHighlight", ZhName: "变化高亮 闪烁 涨跌", Category: "Agent Ops & Motion",
		Demo: motionBlinkDemo,
		Code: `price := 101.5
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.BlinkHighlight(c, price, feedback.BlinkHighlightOptions{})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; the value and direction are read out.",
	},
}
