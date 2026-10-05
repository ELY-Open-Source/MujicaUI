package mujica_test

import (
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleDiffViewer() {
	var state mujica.DiffViewerState
	from := "func main() {\n\tprintln(\"hi\")\n}\n"
	to := "func main() {\n\tprintln(\"hello\")\n}\n"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DiffViewer(c, &state, from, to, mujica.DiffViewerOptions{Language: "go", SideBySide: true}).Height(240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleInlineDiff() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.InlineDiff(c, "The court meets at noon", "The court meets at dusk", mujica.InlineDiffOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleThreeWayMerge() {
	var state mujica.ThreeWayMergeState
	base, local, remote := "a\nb\nc", "a\nB1\nc", "a\nB2\nc"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.ThreeWayMerge(c, &state, base, local, remote, mujica.ThreeWayMergeOptions{})
		r.Element.Height(280)
		if done, total := state.Resolved(); r.Changed() && done == total {
			_ = state.Result() // write the merged file
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleConflictResolver() {
	var state mujica.ConflictResolverState
	text := "a\n<<<<<<< HEAD\nmine\n=======\ntheirs\n>>>>>>> feature\nb"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ConflictResolver(c, &state, text, mujica.ConflictResolverOptions{}).Changed() {
			_ = state.Result() // save the file once every block is resolved
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleChangesList() {
	var state mujica.ChangesListState
	files := []mujica.ChangedFile{
		{Path: "main.go", Status: mujica.GitModified, Staged: true},
		{Path: "notes.txt", Status: mujica.GitUntracked},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.ChangesList(c, &state, files, mujica.ChangesListOptions{})
		r.Element.Height(240)
		if action, path, ok := r.Action(); ok {
			_, _ = action, path // run git add or git restore --staged
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommitInput() {
	msg := mujica.CommitMessage{}
	busy := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.CommitInput(c, &msg, mujica.CommitInputOptions{Busy: busy})
		if r.Committed() || r.Amended() {
			busy = true // run git commit, then clear busy
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommitGraph() {
	var state mujica.CommitGraphState
	commits := []mujica.Commit{
		{Hash: "c3", Parents: []string{"c2", "b1"}, Subject: "Merge branch", Refs: []string{"main"}},
		{Hash: "b1", Parents: []string{"c1"}, Subject: "Feature"},
		{Hash: "c2", Parents: []string{"c1"}, Subject: "Fix"},
		{Hash: "c1", Subject: "Initial commit"},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CommitGraph(c, &state, commits, mujica.CommitGraphOptions{}).Element.Height(240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommitList() {
	var state mujica.ListState[string]
	commits := []mujica.Commit{{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Subject: "Add the commit list"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.CommitList(c, &state, commits, mujica.CommitListOptions{})
		r.Element.Height(240)
		if hash, ok := state.Selected(); ok && r.Submitted() {
			_ = hash // show the commit
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBranchSelector() {
	branch := "main"
	branches := []mujica.Branch{{Name: "main", Current: true}, {Name: "origin/main", Remote: true}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.BranchSelector(c, &branch, branches, mujica.BranchSelectorOptions{AllowCreate: true})
		if name, ok := r.Created(); ok {
			_ = name // run git switch -c name
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBranchList() {
	var state mujica.ListState[string]
	branches := []mujica.Branch{{Name: "main", Current: true, Upstream: "origin/main", Ahead: 2}}
	actions := []mujica.GitAction{{ID: "delete", Label: "Delete", Icon: "trash-2"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.BranchList(c, &state, branches, mujica.BranchListOptions{Actions: actions})
		r.Element.Height(200)
		if action, name, ok := r.Action(); ok {
			_, _ = action, name // run the action on the branch
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTagList() {
	var state mujica.ListState[string]
	tags := []mujica.GitTag{{Name: "v1.1.0", Target: "a1b2c3d4", Annotation: "Release"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.TagList(c, &state, tags, mujica.TagListOptions{})
		r.Element.Height(200)
		if _, hash, ok := r.Action(); ok {
			_ = hash // show the tagged commit
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStashList() {
	var state mujica.ListState[string]
	stashes := []mujica.Stash{{Ref: "stash@{0}", Message: "WIP on main", Branch: "main"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.StashList(c, &state, stashes, mujica.StashListOptions{})
		r.Element.Height(200)
		if action, ref, ok := r.Action(); ok {
			_, _ = action, ref // run git stash apply, pop or drop
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBlameView() {
	var state mujica.BlameViewState
	lines := []mujica.BlameLine{
		{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Summary: "Add main", Text: "package main"},
		{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Summary: "Add main", Text: "func main() {}"},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.BlameView(c, &state, lines, mujica.BlameViewOptions{Language: "go"}).Height(240)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFileHistory() {
	var state mujica.ListState[string]
	commits := []mujica.Commit{{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Subject: "Edit main.go"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.FileHistory(c, &state, "main.go", commits, mujica.FileHistoryOptions{})
		r.Element.Height(240)
		if _, hash, ok := r.Action(); ok {
			_ = hash // open the diff of main.go at hash
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePullRequestCard() {
	pr := mujica.PullRequest{Number: 42, Title: "Add git components", Author: "Sakiko", Head: "feature", Base: "main",
		ChecksPassed: 3, Reviewers: []mujica.Reviewer{{Name: "Uika", State: mujica.ReviewApproved}}, Added: 120, Deleted: 8, Files: 6}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.PullRequestCard(c, pr, mujica.PullRequestCardOptions{}).Changed() {
			// open the pull request
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleReviewComment() {
	thread := mujica.ReviewThread{File: "main.go", Line: 12, Messages: []mujica.ReviewMessage{
		{Author: "Uika", When: time.Now(), Body: "Can this return early?"},
	}}
	draft := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.ReviewComment(c, thread, &draft, mujica.ReviewCommentOptions{})
		if text, ok := r.Replied(); ok {
			thread.Messages = append(thread.Messages, mujica.ReviewMessage{Author: "Sakiko", When: time.Now(), Body: text})
		}
		if r.Resolved() || r.Reopened() {
			thread.Resolved = !thread.Resolved
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDiffStat() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DiffStat(c, 120, 8, mujica.DiffStatOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleGitStatusBadge() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.GitStatusBadge(c, mujica.GitModified, mujica.GitStatusBadgeOptions{ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommandBlock() {
	open := true
	run := mujica.CommandRun{Command: "go test ./...", Output: "\x1b[32mok\x1b[0m  example  0.4s", Duration: 400 * time.Millisecond}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.CommandBlock(c, run, &open, mujica.CommandBlockOptions{})
		if what, ok := r.Copied(); ok {
			_ = what // "command" or "output" is on the clipboard
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAnsiText() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.AnsiText(c, "\x1b[1;32mPASS\x1b[0m 12 tests\n\x1b[31mFAIL\x1b[0m 1 test", mujica.AnsiTextOptions{})
	}
	ui.Render(view, 400, 300, 1)
}
