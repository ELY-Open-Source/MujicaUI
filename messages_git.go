package mujica

func init() {
	addMessages(map[string][2]string{
		"git_status_badge.modified":  {"Modified", "已修改"},
		"git_status_badge.added":     {"Added", "已添加"},
		"git_status_badge.deleted":   {"Deleted", "已删除"},
		"git_status_badge.untracked": {"Untracked", "未跟踪"},
		"git_status_badge.renamed":   {"Renamed", "已重命名"},
		"git_status_badge.copied":    {"Copied", "已复制"},

		"diff_stat.label": {"%d additions, %d deletions", "新增 %d 行，删除 %d 行"},

		"command_block.copy_command": {"Copy command", "复制命令"},
		"command_block.copy_output":  {"Copy output", "复制输出"},
		"command_block.exit":         {"exit %d", "退出码 %d"},
		"command_block.no_output":    {"No output", "无输出"},
		"command_block.output":       {"Command output", "命令输出"},
		"command_block.running":      {"Running", "运行中"},

		"commit_list.empty": {"No commits", "没有提交"},

		"branch_list.ahead":   {"%d ahead of upstream", "领先上游 %d 个提交"},
		"branch_list.behind":  {"%d behind upstream", "落后上游 %d 个提交"},
		"branch_list.current": {"current", "当前"},
		"branch_list.empty":   {"No branches", "没有分支"},

		"branch_selector.create": {"Create branch", "新建分支"},
		"branch_selector.empty":  {"No matching branches", "没有匹配的分支"},
		"branch_selector.label":  {"Branch", "分支"},
		"branch_selector.search": {"Search branches", "搜索分支"},
		"branch_selector.local":  {"Local", "本地"},
		"branch_selector.remote": {"Remote", "远端"},

		"tag_list.empty": {"No tags", "没有标签"},
		"tag_list.jump":  {"Go to commit", "跳转到提交"},

		"stash_list.apply": {"Apply", "应用"},
		"stash_list.pop":   {"Pop", "弹出"},
		"stash_list.drop":  {"Drop", "删除"},
		"stash_list.empty": {"No stashes", "没有贮藏"},

		"blame_view.line": {"Line %d, %s, %s", "第 %d 行，%s，%s"},

		"file_history.title": {"History", "历史"},
		"file_history.diff":  {"View changes", "查看差异"},
		"file_history.empty": {"No history for this file", "此文件没有历史"},

		"pull_request_card.open":           {"Open", "开放"},
		"pull_request_card.draft":          {"Draft", "草稿"},
		"pull_request_card.merged":         {"Merged", "已合并"},
		"pull_request_card.closed":         {"Closed", "已关闭"},
		"pull_request_card.files":          {"%d files", "%d 个文件"},
		"pull_request_card.passed":         {"%d checks passed", "%d 项检查通过"},
		"pull_request_card.failed":         {"%d checks failed", "%d 项检查失败"},
		"pull_request_card.pending":        {"%d checks pending", "%d 项检查进行中"},
		"pull_request_card.approved":       {"approved", "已批准"},
		"pull_request_card.changes":        {"changes requested", "要求修改"},
		"pull_request_card.review_pending": {"review pending", "待评审"},

		"review_comment.placeholder": {"Write a reply", "写下回复"},
		"review_comment.reply":       {"Reply", "回复"},
		"review_comment.field":       {"Reply text", "回复内容"},
		"review_comment.resolve":     {"Resolve", "解决"},
		"review_comment.reopen":      {"Reopen", "重新打开"},
		"review_comment.resolved":    {"Resolved", "已解决"},

		"changes_list.staged":      {"Staged changes (%d)", "已暂存的更改（%d）"},
		"changes_list.unstaged":    {"Changes (%d)", "更改（%d）"},
		"changes_list.stage":       {"Stage", "暂存"},
		"changes_list.unstage":     {"Unstage", "取消暂存"},
		"changes_list.stage_all":   {"Stage all", "全部暂存"},
		"changes_list.unstage_all": {"Unstage all", "全部取消暂存"},

		"diff_viewer.added":   {"Added", "新增"},
		"diff_viewer.deleted": {"Deleted", "删除"},
		"diff_viewer.fold":    {"%d unchanged lines", "%d 行未更改"},

		"inline_diff.removed": {"Old", "原文"},
		"inline_diff.added":   {"New", "新文"},

		"three_way_merge.resolved":    {"%d of %d conflicts resolved", "已解决 %d / %d 个冲突"},
		"three_way_merge.ours":        {"Local", "本地"},
		"three_way_merge.base":        {"Base", "基准"},
		"three_way_merge.theirs":      {"Remote", "远端"},
		"three_way_merge.result":      {"Result", "结果"},
		"three_way_merge.unresolved":  {"Unresolved", "未解决"},
		"three_way_merge.empty":       {"No lines", "无内容"},
		"three_way_merge.take_ours":   {"Take local", "采纳本地"},
		"three_way_merge.take_theirs": {"Take remote", "采纳远端"},
		"three_way_merge.take_both":   {"Take both", "两者都要"},
		"three_way_merge.take_base":   {"Keep base", "保留基准"},

		"conflict_resolver.resolved":      {"%d of %d conflicts resolved", "已解决 %d / %d 个冲突"},
		"conflict_resolver.broken":        {"Conflict markers are not closed after line %d", "第 %d 行起的冲突标记未闭合"},
		"conflict_resolver.none":          {"No conflicts", "没有冲突"},
		"conflict_resolver.title":         {"Conflict %d · line %d", "冲突 %d · 第 %d 行"},
		"conflict_resolver.open":          {"unresolved", "未解决"},
		"conflict_resolver.done":          {"resolved", "已解决"},
		"conflict_resolver.current":       {"Current", "当前"},
		"conflict_resolver.incoming":      {"Incoming", "传入"},
		"conflict_resolver.take_current":  {"Accept current", "采纳当前"},
		"conflict_resolver.take_incoming": {"Accept incoming", "采纳传入"},
		"conflict_resolver.take_both":     {"Accept both", "两者都采纳"},
		"conflict_resolver.reset":         {"Undo", "撤销"},

		"commit_input.title":  {"Summary", "提交标题"},
		"commit_input.body":   {"Description", "提交说明"},
		"commit_input.count":  {"%d of %d characters", "%d / %d 字符"},
		"commit_input.long":   {"Summary is %d characters over the suggested length", "标题超出建议长度 %d 个字符"},
		"commit_input.commit": {"Commit", "提交"},
		"commit_input.amend":  {"Amend", "修订"},
	})
}
