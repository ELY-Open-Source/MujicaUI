package mujica

func init() {
	addMessages(map[string][2]string{
		"interrupt_button.stop":     {"Stop", "停止"},
		"interrupt_button.stopping": {"Stopping", "正在停止"},
		"interrupt_button.confirm":  {"Stop the agent?", "停止智能体？"},
		"interrupt_button.detail":   {"Work in progress will be lost.", "进行中的工作将会丢失。"},
		"interrupt_button.keep":     {"Keep running", "继续运行"},

		"token_usage_chart.label":  {"Token usage", "Token 用量"},
		"token_usage_chart.input":  {"Input", "输入"},
		"token_usage_chart.output": {"Output", "输出"},

		"cost_breakdown.model":      {"Model", "模型"},
		"cost_breakdown.tool":       {"Tool", "工具"},
		"cost_breakdown.step":       {"Step", "步骤"},
		"cost_breakdown.calls":      {"Calls", "调用"},
		"cost_breakdown.input":      {"Input tokens", "输入 token"},
		"cost_breakdown.output":     {"Output tokens", "输出 token"},
		"cost_breakdown.cost":       {"Cost", "费用"},
		"cost_breakdown.empty":      {"No costs yet", "暂无费用"},
		"cost_breakdown.unnamed":    {"(unnamed)", "（未命名）"},
		"cost_breakdown.total":      {"Total", "合计"},
		"cost_breakdown.total_name": {"Total %s over %d calls", "合计 %s，共 %d 次调用"},
		"cost_breakdown.tokens":     {"%d in · %d out", "输入 %d · 输出 %d"},

		"trace_viewer.span":     {"Span", "跨度"},
		"trace_viewer.duration": {"Duration", "耗时"},
		"trace_viewer.timeline": {"Timeline", "时间轴"},
		"trace_viewer.start":    {"Start", "开始"},
		"trace_viewer.failed":   {"failed", "失败"},
		"trace_viewer.empty":    {"No spans", "没有跨度"},

		"eval_result_table.case":    {"Case", "用例"},
		"eval_result_table.score":   {"Score", "得分"},
		"eval_result_table.result":  {"Result", "结果"},
		"eval_result_table.pass":    {"Pass", "通过"},
		"eval_result_table.fail":    {"Fail", "失败"},
		"eval_result_table.diff":    {"View differences for %s", "查看 %s 的差异"},
		"eval_result_table.view":    {"View differences", "查看差异"},
		"eval_result_table.summary": {"%d of %d passed", "通过 %d / %d"},
		"eval_result_table.empty":   {"No eval cases", "没有评测用例"},

		"tool_registry_panel.empty":     {"No tools", "没有工具"},
		"tool_registry_panel.choose":    {"Choose a tool to see its parameters", "选择工具查看参数"},
		"tool_registry_panel.no_params": {"No parameters", "无参数"},
		"tool_registry_panel.enable":    {"Enable %s", "启用 %s"},
		"tool_registry_panel.required":  {"required", "必填"},

		"typewriter.skip": {"Skip", "跳过"},

		"blink_highlight.up":   {"up", "上升"},
		"blink_highlight.down": {"down", "下降"},
	})
}
