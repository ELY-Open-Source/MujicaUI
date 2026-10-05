package git_test

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/agent"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/git"
)

func ExampleDiffViewer() {
	var state git.DiffViewerState
	from := "func main() {\n\tprintln(\"hi\")\n}\n"
	to := "func main() {\n\tprintln(\"hello\")\n}\n"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		git.DiffViewer(c, &state, from, to, git.DiffViewerOptions{Language: "go", SideBySide: true}).Height(240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleInlineDiff() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		git.InlineDiff(c, "The court meets at noon", "The court meets at dusk", git.InlineDiffOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleThreeWayMerge() {
	var state git.ThreeWayMergeState
	base, local, remote := "a\nb\nc", "a\nB1\nc", "a\nB2\nc"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := git.ThreeWayMerge(c, &state, base, local, remote, git.ThreeWayMergeOptions{})
		r.Element.Height(280)
		if done, total := state.Resolved(); r.Changed() && done == total {
			_ = state.Result() // write the merged file
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleConflictResolver() {
	var state git.ConflictResolverState
	text := "a\n<<<<<<< HEAD\nmine\n=======\ntheirs\n>>>>>>> feature\nb"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if git.ConflictResolver(c, &state, text, git.ConflictResolverOptions{}).Changed() {
			_ = state.Result() // save the file once every block is resolved
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleChangesList() {
	var state git.ChangesListState
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
	ui.Render(view, 400, 300, 1)
}

func ExampleCommitInput() {
	msg := git.CommitMessage{}
	busy := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := git.CommitInput(c, &msg, git.CommitInputOptions{Busy: busy})
		if r.Committed() || r.Amended() {
			busy = true // run git commit, then clear busy
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommitGraph() {
	var state git.CommitGraphState
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
	ui.Render(view, 400, 300, 1)
}

func ExampleCommitList() {
	var state data.ListState[string]
	commits := []git.Commit{{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Subject: "Add the commit list"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := git.CommitList(c, &state, commits, git.CommitListOptions{})
		r.Element.Height(240)
		if hash, ok := state.Selected(); ok && r.Submitted() {
			_ = hash // show the commit
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBranchSelector() {
	branch := "main"
	branches := []git.Branch{{Name: "main", Current: true}, {Name: "origin/main", Remote: true}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := git.BranchSelector(c, &branch, branches, git.BranchSelectorOptions{AllowCreate: true})
		if name, ok := r.Created(); ok {
			_ = name // run git switch -c name
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBranchList() {
	var state data.ListState[string]
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
	ui.Render(view, 400, 300, 1)
}

func ExampleTagList() {
	var state data.ListState[string]
	tags := []git.GitTag{{Name: "v1.1.0", Target: "a1b2c3d4", Annotation: "Release"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := git.TagList(c, &state, tags, git.TagListOptions{})
		r.Element.Height(200)
		if _, hash, ok := r.Action(); ok {
			_ = hash // show the tagged commit
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStashList() {
	var state data.ListState[string]
	stashes := []git.Stash{{Ref: "stash@{0}", Message: "WIP on main", Branch: "main"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := git.StashList(c, &state, stashes, git.StashListOptions{})
		r.Element.Height(200)
		if action, ref, ok := r.Action(); ok {
			_, _ = action, ref // run git stash apply, pop or drop
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBlameView() {
	var state git.BlameViewState
	lines := []git.BlameLine{
		{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Summary: "Add main", Text: "package main"},
		{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Summary: "Add main", Text: "func main() {}"},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		git.BlameView(c, &state, lines, git.BlameViewOptions{Language: "go"}).Height(240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFileHistory() {
	var state data.ListState[string]
	commits := []git.Commit{{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Subject: "Edit main.go"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := git.FileHistory(c, &state, "main.go", commits, git.FileHistoryOptions{})
		r.Element.Height(240)
		if _, hash, ok := r.Action(); ok {
			_ = hash // open the diff of main.go at hash
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePullRequestCard() {
	pr := git.PullRequest{Number: 42, Title: "Add git components", Author: "Sakiko", Head: "feature", Base: "main",
		ChecksPassed: 3, Reviewers: []git.Reviewer{{Name: "Uika", State: git.ReviewApproved}}, Added: 120, Deleted: 8, Files: 6}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if git.PullRequestCard(c, pr, git.PullRequestCardOptions{}).Changed() {
			// open the pull request
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleReviewComment() {
	thread := git.ReviewThread{File: "main.go", Line: 12, Messages: []git.ReviewMessage{
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
	ui.Render(view, 400, 300, 1)
}

func ExampleDiffStat() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		git.DiffStat(c, 120, 8, git.DiffStatOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleGitStatusBadge() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		git.GitStatusBadge(c, git.GitModified, git.GitStatusBadgeOptions{ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommandBlock() {
	open := true
	run := agent.CommandRun{Command: "go test ./...", Output: "\x1b[32mok\x1b[0m  example  0.4s", Duration: 400 * time.Millisecond}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := git.CommandBlock(c, run, &open, git.CommandBlockOptions{})
		if what, ok := r.Copied(); ok {
			_ = what // "command" or "output" is on the clipboard
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAnsiText() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		git.AnsiText(c, "\x1b[1;32mPASS\x1b[0m 12 tests\n\x1b[31mFAIL\x1b[0m 1 test", git.AnsiTextOptions{})
	}
	ui.Render(view, 400, 300, 1)
}
