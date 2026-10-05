package main

var gitPages = []page{
	{
		ID: "301", Name: "DiffViewer", ZhName: "差异查看器 比较 diff", Category: "Git",
		Demo: gitDiffViewerDemo,
		Code: `var state git.DiffViewerState
from := "func main() {\n\tprintln(\"hi\")\n}\n"
to := "func main() {\n\tprintln(\"hello\")\n}\n"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.DiffViewer(c, &state, from, to, git.DiffViewerOptions{Language: "go", SideBySide: true}).Height(240)
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab focuses the list; arrows and Page Up/Down scroll. Fold buttons take Tab; Enter or Space shows the hidden lines.",
	},
	{
		ID: "302", Name: "InlineDiff", ZhName: "行内差异 词级差异", Category: "Git",
		Demo: gitInlineDiffDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.InlineDiff(c, "The court meets at noon", "The court meets at dusk", git.InlineDiffOptions{})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; each line is read with its full text.",
	},
	{
		ID: "303", Name: "ThreeWayMerge", ZhName: "三方合并 合并冲突", Category: "Git",
		Demo: gitThreeWayMergeDemo,
		Code: `var state git.ThreeWayMergeState
base, local, remote := "a\nb\nc", "a\nB1\nc", "a\nB2\nc"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.ThreeWayMerge(c, &state, base, local, remote, git.ThreeWayMergeOptions{})
	r.Element.Height(280)
	if done, total := state.Resolved(); r.Changed() && done == total {
		_ = state.Result() // write the merged file
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab reaches the take buttons of each conflict; Enter or Space chooses. The list scrolls with the arrows.",
	},
	{
		ID: "304", Name: "ConflictResolver", ZhName: "冲突解决 冲突标记", Category: "Git",
		Demo: gitConflictResolverDemo,
		Code: `var state git.ConflictResolverState
text := "a\n<<<<<<< HEAD\nmine\n=======\ntheirs\n>>>>>>> feature\nb"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if git.ConflictResolver(c, &state, text, git.ConflictResolverOptions{}).Changed() {
		_ = state.Result() // save the file once every block is resolved
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab moves through each block's buttons; Enter or Space accepts a side; Undo clears the choice.",
	},
	{
		ID: "305", Name: "ChangesList", ZhName: "更改列表 暂存 未暂存", Category: "Git",
		Demo: gitChangesListDemo,
		Code: `var state git.ChangesListState
files := []git.ChangedFile{
	{Path: "main.go", Status: git.GitModified, Staged: true},
	{Path: "notes.txt", Status: git.GitUntracked},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.ChangesList(c, &state, files, git.ChangesListOptions{})
	r.Element.Height(240)
	if action, path, ok := r.Action(); ok {
		_, _ = action, path // run git add or git restore --staged
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a file, skipping group titles; typing jumps by name; stage buttons take Tab, Enter or Space.",
	},
	{
		ID: "306", Name: "CommitInput", ZhName: "提交输入 提交信息", Category: "Git",
		Demo: gitCommitInputDemo,
		Code: `msg := git.CommitMessage{}
busy := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.CommitInput(c, &msg, git.CommitInputOptions{Busy: busy})
	if r.Committed() || r.Amended() {
		busy = true // run git commit, then clear busy
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Enter in the summary or Cmd/Ctrl+Enter in the description commits; Tab reaches Amend and Commit.",
	},
	{
		ID: "307", Name: "CommitGraph", ZhName: "提交图谱 泳道 分支图", Category: "Git",
		Demo: gitCommitGraphDemo,
		Code: `var state git.CommitGraphState
commits := []git.Commit{
	{Hash: "c3", Parents: []string{"c2", "b1"}, Subject: "Merge branch", Refs: []string{"main"}},
	{Hash: "b1", Parents: []string{"c1"}, Subject: "Feature"},
	{Hash: "c2", Parents: []string{"c1"}, Subject: "Fix"},
	{Hash: "c1", Subject: "Initial commit"},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.CommitGraph(c, &state, commits, git.CommitGraphOptions{}).Element.Height(240)
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down/Home/End choose a commit; typing jumps by subject; Enter opens it.",
	},
	{
		ID: "308", Name: "CommitList", ZhName: "提交列表 提交历史", Category: "Git",
		Demo: gitCommitListDemo,
		Code: `var state data.ListState[string]
commits := []git.Commit{{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Subject: "Add the commit list"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.CommitList(c, &state, commits, git.CommitListOptions{})
	r.Element.Height(240)
	if hash, ok := state.Selected(); ok && r.Submitted() {
		_ = hash // show the commit
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down/Home/End choose a commit; typing jumps by subject; Enter opens it.",
	},
	{
		ID: "309", Name: "BranchSelector", ZhName: "分支选择器 切换分支", Category: "Git",
		Demo: gitBranchSelectorDemo,
		Code: `branch := "main"
branches := []git.Branch{{Name: "main", Current: true}, {Name: "origin/main", Remote: true}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.BranchSelector(c, &branch, branches, git.BranchSelectorOptions{AllowCreate: true})
	if name, ok := r.Created(); ok {
		_ = name // run git switch -c name
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Enter or Space opens; type to filter; Up/Down move; Enter chooses or creates; Escape closes.",
	},
	{
		ID: "310", Name: "BranchList", ZhName: "分支列表 上游 领先落后", Category: "Git",
		Demo: gitBranchListDemo,
		Code: `var state data.ListState[string]
branches := []git.Branch{{Name: "main", Current: true, Upstream: "origin/main", Ahead: 2}}
actions := []git.GitAction{{ID: "delete", Label: "Delete", Icon: "trash-2"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.BranchList(c, &state, branches, git.BranchListOptions{Actions: actions})
	r.Element.Height(200)
	if action, name, ok := r.Action(); ok {
		_, _ = action, name // run the action on the branch
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a branch; row actions take Tab, Enter or Space.",
	},
	{
		ID: "311", Name: "TagList", ZhName: "标签列表 tag", Category: "Git",
		Demo: gitTagListDemo,
		Code: `var state data.ListState[string]
tags := []git.GitTag{{Name: "v1.1.0", Target: "a1b2c3d4", Annotation: "Release"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.TagList(c, &state, tags, git.TagListOptions{})
	r.Element.Height(200)
	if _, hash, ok := r.Action(); ok {
		_ = hash // show the tagged commit
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a tag; the jump button takes Tab, Enter or Space.",
	},
	{
		ID: "312", Name: "StashList", ZhName: "贮藏列表 stash", Category: "Git",
		Demo: gitStashListDemo,
		Code: `var state data.ListState[string]
stashes := []git.Stash{{Ref: "stash@{0}", Message: "WIP on main", Branch: "main"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.StashList(c, &state, stashes, git.StashListOptions{})
	r.Element.Height(200)
	if action, ref, ok := r.Action(); ok {
		_, _ = action, ref // run git stash apply, pop or drop
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a stash; Apply, Pop and Drop take Tab, Enter or Space.",
	},
	{
		ID: "313", Name: "BlameView", ZhName: "逐行追溯 blame 作者", Category: "Git",
		Demo: gitBlameViewDemo,
		Code: `var state git.BlameViewState
lines := []git.BlameLine{
	{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Summary: "Add main", Text: "package main"},
	{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Summary: "Add main", Text: "func main() {}"},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.BlameView(c, &state, lines, git.BlameViewOptions{Language: "go"}).Height(240)
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab focuses the code; arrows scroll; typing jumps to a line starting with the text. Hover the left column for details.",
	},
	{
		ID: "314", Name: "FileHistory", ZhName: "文件历史 单文件提交", Category: "Git",
		Demo: gitFileHistoryDemo,
		Code: `var state data.ListState[string]
commits := []git.Commit{{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Subject: "Edit main.go"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.FileHistory(c, &state, "main.go", commits, git.FileHistoryOptions{})
	r.Element.Height(240)
	if _, hash, ok := r.Action(); ok {
		_ = hash // open the diff of main.go at hash
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Up/Down choose a commit; View changes takes Tab, Enter or Space.",
	},
	{
		ID: "315", Name: "PullRequestCard", ZhName: "拉取请求卡片 合并请求 PR", Category: "Git",
		Demo: gitPullRequestCardDemo,
		Code: `pr := git.PullRequest{Number: 42, Title: "Add git components", Author: "Sakiko", Head: "feature", Base: "main",
	ChecksPassed: 3, Reviewers: []git.Reviewer{{Name: "Uika", State: git.ReviewApproved}}, Added: 120, Deleted: 8, Files: 6}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if git.PullRequestCard(c, pr, git.PullRequestCardOptions{}).Changed() {
		// open the pull request
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab reaches the title link; Enter or Space opens the pull request.",
	},
	{
		ID: "316", Name: "ReviewComment", ZhName: "评审评论 代码评审", Category: "Git",
		Demo: gitReviewCommentDemo,
		Code: `thread := git.ReviewThread{File: "main.go", Line: 12, Messages: []git.ReviewMessage{
	{Author: "Uika", When: time.Now(), Body: "Can this return early?"},
}}
draft := ""
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.ReviewComment(c, thread, &draft, git.ReviewCommentOptions{})
	if text, ok := r.Replied(); ok {
		thread.Messages = append(thread.Messages, git.ReviewMessage{Author: "Sakiko", When: time.Now(), Body: text})
	}
	if r.Resolved() || r.Reopened() {
		thread.Resolved = !thread.Resolved
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab moves through the reply field and buttons; Cmd/Ctrl+Enter sends the reply.",
	},
	{
		ID: "317", Name: "DiffStat", ZhName: "差异统计 增删行数", Category: "Git",
		Demo: gitDiffStatDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.DiffStat(c, 120, 8, git.DiffStatOptions{})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; read as additions and deletions.",
	},
	{
		ID: "318", Name: "GitStatusBadge", ZhName: "Git 状态标记 文件状态", Category: "Git",
		Demo: gitStatusBadgeDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.GitStatusBadge(c, git.GitModified, git.GitStatusBadgeOptions{ShowLabel: true})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; each badge is read by its status name.",
	},
	{
		ID: "319", Name: "CommandBlock", ZhName: "命令块 终端输出 退出码", Category: "Git",
		Demo: gitCommandBlockDemo,
		Code: `open := true
run := agent.CommandRun{Command: "go test ./...", Output: "\x1b[32mok\x1b[0m  example  0.4s", Duration: 400 * time.Millisecond}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.CommandBlock(c, run, &open, git.CommandBlockOptions{})
	if what, ok := r.Copied(); ok {
		_ = what // "command" or "output" is on the clipboard
	}
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Tab reaches the header; Enter or Space folds the output; the copy buttons take Tab.",
	},
	{
		ID: "320", Name: "AnsiText", ZhName: "ANSI 文本 终端颜色", Category: "Git",
		Demo: gitAnsiTextDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.AnsiText(c, "\x1b[1;32mPASS\x1b[0m 12 tests\n\x1b[31mFAIL\x1b[0m 1 test", git.AnsiTextOptions{})
}
ui.Render(view, 400, 300, 1)`,
		Keys: "Not focusable; shows the text with its styles.",
	},
}
