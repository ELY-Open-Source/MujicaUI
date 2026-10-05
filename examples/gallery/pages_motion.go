package main

import mujica "mujicaui"

var motionPages = []page{
	{
		ID: "281", Name: "InterruptButton", ZhName: "中断按钮 停止智能体", Category: "Agent Ops & Motion",
		Demo: motionInterruptDemo,
		Code: `running, stopping := true, false
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.InterruptButton(c, running, mujica.InterruptButtonOptions{Confirm: true, Interrupting: stopping})
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
		Code: `usage := []mujica.TokenUsage{{Label: "Plan", Input: 1200, Output: 300}, {Label: "Write", Input: 2100, Output: 2900}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TokenUsageChart(c, usage, mujica.TokenUsageChartOptions{})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab focuses the chart; Left/Right move between bars or points and the readout is announced.",
	},
	{
		ID: "283", Name: "CostBreakdown", ZhName: "费用分解 成本 合计", Category: "Agent Ops & Motion",
		Demo: motionCostDemo,
		Code: `var state mujica.CostBreakdownState
items := []mujica.CostItem{
	{Model: "opus", Tool: "search", InputTokens: 1000, OutputTokens: 200, Cost: 0.30},
	{Model: "haiku", Tool: "search", InputTokens: 4000, OutputTokens: 500, Cost: 0.05},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.CostBreakdown(c, &state, items, mujica.CostBreakdownOptions{GroupBy: mujica.CostByModel}).Element.Height(260)
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a group; column headers sort with Enter; the column menu takes Tab.",
	},
	{
		ID: "284", Name: "TraceViewer", ZhName: "调用链 追踪 跨度 时间轴", Category: "Agent Ops & Motion",
		Demo: motionTraceDemo,
		Code: `var state mujica.TraceViewerState
spans := []mujica.TraceSpan{
	{ID: "run", Name: "agent.run", End: 900 * time.Millisecond},
	{ID: "plan", Parent: "run", Name: "llm.plan", Start: 10 * time.Millisecond, End: 210 * time.Millisecond,
		Attrs: []mujica.TraceAttr{{Key: "model", Value: "opus"}}},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TraceViewer(c, &state, spans, mujica.TraceViewerOptions{}).Element.Height(260)
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a span, Right opens and Left closes it; typing jumps by name.",
	},
	{
		ID: "285", Name: "EvalResultTable", ZhName: "评测结果 得分 通过失败", Category: "Agent Ops & Motion",
		Demo: motionEvalDemo,
		Code: `var state mujica.DataTableState[string]
cases := []mujica.EvalCase{{ID: "c1", Name: "Greets by title", Score: 0.98, Passed: true}, {ID: "c2", Name: "Refuses forgery", Score: 0.41}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.EvalResultTable(c, &state, cases, mujica.EvalResultTableOptions{})
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
		Code: `var state mujica.ListState[string]
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
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a tool; its switch takes Tab and Space.",
	},
	{
		ID: "287", Name: "DotsLoader", ZhName: "三点加载 加载中", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("287", mujica.DotsLoader),
		Code: `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.DotsLoader(c, mujica.LoaderOptions{Label: "Thinking", ShowLabel: true})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "288", Name: "BarsLoader", ZhName: "条形加载", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("288", mujica.BarsLoader),
		Code: `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.BarsLoader(c, mujica.LoaderOptions{Size: 20, Label: "Indexing"})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "289", Name: "PulseLoader", ZhName: "呼吸加载 脉冲", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("289", mujica.PulseLoader),
		Code: `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PulseLoader(c, mujica.LoaderOptions{Label: "Listening", ShowLabel: true})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "290", Name: "OrbitLoader", ZhName: "环绕加载 菱形", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("290", mujica.OrbitLoader),
		Code: `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.OrbitLoader(c, mujica.LoaderOptions{Size: 24, Label: "Connecting"})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "291", Name: "WaveLoader", ZhName: "波形加载", Category: "Agent Ops & Motion",
		Demo: motionLoaderDemo("291", mujica.WaveLoader),
		Code: `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.WaveLoader(c, mujica.LoaderOptions{Label: "Speaking", ShowLabel: true})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as progress.",
	},
	{
		ID: "292", Name: "Shimmer", ZhName: "占位扫光 骨架 加载", Category: "Agent Ops & Motion",
		Demo: motionShimmerDemo,
		Code: `loading := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Shimmer(c, loading, mujica.ShimmerOptions{Height: 48}, func() {
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
	mujica.Use(c, mujica.Settings{})
	mujica.CountUp(c, 12840, mujica.CountUpOptions{Duration: 800 * time.Millisecond})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; the final value is read out.",
	},
	{
		ID: "294", Name: "Typewriter", ZhName: "打字机 逐字显示", Category: "Agent Ops & Motion",
		Demo: motionTypewriterDemo,
		Code: `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.Typewriter(c, "The court assembles at dusk.", mujica.TypewriterOptions{Speed: 30, Skippable: true})
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
	mujica.Use(c, mujica.Settings{})
	mujica.ShimmerText(c, "Thinking…", mujica.ShimmerTextOptions{Active: thinking})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as a status.",
	},
	{
		ID: "296", Name: "Stagger", ZhName: "依次进场 交错动画", Category: "Agent Ops & Motion",
		Demo: motionStaggerDemo,
		Code: `guests := []string{"Sakiko", "Uika", "Mutsumi"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Stagger(c, len(guests), mujica.StaggerOptions{}, func(i int) {
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
	mujica.Use(c, mujica.Settings{})
	mujica.ReorderTransition(c, order, mujica.ReorderTransitionOptions{Gap: 4}, func(name string) {
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
	mujica.Use(c, mujica.Settings{})
	mujica.LayoutTransition(c, mujica.LayoutTransitionOptions{Clip: true}, func() {
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
	mujica.Use(c, mujica.Settings{})
	mujica.Presence(c, show, mujica.PresenceOptions{Kind: mujica.PresenceSlide}, func() {
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
	mujica.Use(c, mujica.Settings{})
	mujica.BlinkHighlight(c, price, mujica.BlinkHighlightOptions{})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; the value and direction are read out.",
	},
}
