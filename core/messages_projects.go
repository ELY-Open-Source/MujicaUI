package core

func init() {
	addMessages(map[string][2]string{
		"projects.overdue":      {"Overdue by %d days", "逾期 %d 天"},
		"projects.overdue_one":  {"Overdue by %d day", "逾期 %d 天"},
		"projects.due_today":    {"Due today", "今天到期"},
		"projects.due_tomorrow": {"Due tomorrow", "明天到期"},
		"projects.due":          {"Due %s %d", "%s%d日到期"},

		"priority_indicator.none":   {"No priority", "无优先级"},
		"priority_indicator.low":    {"Low", "低"},
		"priority_indicator.medium": {"Medium", "中"},
		"priority_indicator.high":   {"High", "高"},
		"priority_indicator.urgent": {"Urgent", "紧急"},

		"issue_card.description": {"%s, %s", "%s，%s"},

		"issue_id_badge.copy":         {"Copy %s", "复制 %s"},
		"issue_id_badge.copied":       {"Copied", "已复制"},
		"issue_id_badge.copied_label": {"%s copied", "已复制 %s"},

		"status_select.label":     {"Status", "状态"},
		"status_select.backlog":   {"Backlog", "待规划"},
		"status_select.unstarted": {"To do", "未开始"},
		"status_select.started":   {"In progress", "进行中"},
		"status_select.completed": {"Done", "已完成"},
		"status_select.canceled":  {"Canceled", "已取消"},

		"assignee_picker.label":       {"Assignees", "负责人"},
		"assignee_picker.placeholder": {"Unassigned", "未分配"},

		"due_date_picker.label":     {"Due date", "截止日期"},
		"due_date_picker.today":     {"Today", "今天"},
		"due_date_picker.tomorrow":  {"Tomorrow", "明天"},
		"due_date_picker.next_week": {"Next week", "下周"},

		"task_item.rename": {"Task title", "任务标题"},

		"subtask_list.label":    {"Subtasks", "子任务"},
		"subtask_list.progress": {"%d of %d done", "已完成 %d / %d"},
		"subtask_list.drag":     {"Drag %s", "拖动 %s"},
		"subtask_list.remove":   {"Remove %s", "删除 %s"},
		"subtask_list.add":      {"Add a subtask", "添加子任务"},

		"label_manager.label":   {"Labels", "标签"},
		"label_manager.edit":    {"Edit %s", "编辑 %s"},
		"label_manager.delete":  {"Delete %s", "删除 %s"},
		"label_manager.empty":   {"No labels yet.", "暂无标签。"},
		"label_manager.new":     {"New label", "新建标签"},
		"label_manager.name":    {"Label name", "标签名称"},
		"label_manager.color":   {"Color", "颜色"},
		"label_manager.taken":   {"A label named %q exists.", "已有名为 %q 的标签。"},
		"label_manager.accent":  {"Crimson", "酒红"},
		"label_manager.info":    {"Blue", "蓝"},
		"label_manager.success": {"Green", "绿"},
		"label_manager.warning": {"Amber", "琥珀"},
		"label_manager.danger":  {"Red", "红"},
		"label_manager.neutral": {"Gray", "灰"},

		"task_list.label":  {"Tasks", "任务"},
		"task_list.filter": {"Show", "显示"},
		"task_list.all":    {"All %d", "全部 %d"},
		"task_list.open":   {"Open %d", "未完成 %d"},
		"task_list.done":   {"Done %d", "已完成 %d"},
		"task_list.empty":  {"No tasks here.", "这里没有任务。"},

		"task_detail_panel.close":    {"Close details", "关闭详情"},
		"task_detail_panel.title":    {"Title", "标题"},
		"task_detail_panel.status":   {"Status", "状态"},
		"task_detail_panel.priority": {"Priority", "优先级"},
		"task_detail_panel.activity": {"Activity", "活动记录"},

		"kanban_board.label": {"Board", "看板"},
		"kanban_board.full":  {"%s is at its limit of %d cards.", "%s 已达 %d 张卡片的上限。"},
		"kanban_board.wip":   {"%d of %d cards", "%d / %d 张卡片"},
		"kanban_board.over":  {"%d over the limit", "超出上限 %d 张"},
		"kanban_board.empty": {"No cards", "暂无卡片"},

		"sprint_board.time":          {"Time", "时间"},
		"sprint_board.done":          {"%d of %d done", "已完成 %d / %d"},
		"sprint_board.progress":      {"Sprint progress", "迭代进度"},
		"sprint_board.starts_in":     {"Starts in %d days", "%d 天后开始"},
		"sprint_board.starts_in_one": {"Starts in %d day", "%d 天后开始"},
		"sprint_board.ended":         {"Ended", "已结束"},
		"sprint_board.left":          {"%d days left", "剩余 %d 天"},
		"sprint_board.left_one":      {"%d day left", "剩余 %d 天"},

		"milestone_progress.count":   {"%d of %d tasks done", "已完成 %d / %d 个任务"},
		"milestone_progress.reached": {"Reached", "已达成"},

		"workload_view.label":          {"Workload", "工作负载"},
		"workload_view.hours":          {"h", "小时"},
		"workload_view.over":           {"over by %s %s", "超出 %s %s"},
		"workload_view.overloaded":     {"%d people are over capacity", "%d 人超出负荷"},
		"workload_view.overloaded_one": {"%d person is over capacity", "%d 人超出负荷"},

		"time_tracker.label":   {"Time tracked", "已记录时间"},
		"time_tracker.start":   {"Start", "开始"},
		"time_tracker.stop":    {"Stop", "停止"},
		"time_tracker.running": {"Tracking…", "计时中…"},
		"time_tracker.correct": {"Correct the total", "修正总时长"},
		"time_tracker.apply":   {"Apply", "应用"},

		"pomodoro_timer.label":      {"Pomodoro timer", "番茄钟"},
		"pomodoro_timer.focus":      {"Focus", "专注"},
		"pomodoro_timer.break":      {"Short break", "短休息"},
		"pomodoro_timer.long_break": {"Long break", "长休息"},
		"pomodoro_timer.done":       {"%s finished.", "%s结束。"},
		"pomodoro_timer.left":       {"%s, %s left", "%s，剩余 %s"},
		"pomodoro_timer.sessions":   {"%d of %d sessions before a long break", "距长休息：%d / %d 轮"},
		"pomodoro_timer.start":      {"Start", "开始"},
		"pomodoro_timer.pause":      {"Pause", "暂停"},
		"pomodoro_timer.skip":       {"Skip", "跳过"},
		"pomodoro_timer.reset":      {"Reset", "重置"},

		"burndown_chart.actual": {"Remaining", "剩余"},
		"burndown_chart.ideal":  {"Ideal", "理想"},
	})
}
