package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/agent"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/git"
	"github.com/ELY-Open-Source/MujicaUI/navigation"
)

func gitPullRequestCardDemo(c *ui.Context, log func(string)) {
	type state struct {
		st               int
		failed, disabled bool
	}
	s := ui.Local(c.Root(), "315", func() state { return state{} })
	navigationDemo(c, func() {
		navigation.SegmentedControl(c, &s.st, []navigation.Segment{{Label: "Open"}, {Label: "Draft"}, {Label: "Merged"}, {Label: "Closed"}}, navigation.SegmentedControlOptions{Label: "State"})
		navigationCheck(c, &s.failed, "Failing check", log)
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		pr := git.PullRequest{Number: 318, Title: "Git components 301–320", Author: "Sakiko", State: git.PullRequestState(s.st),
			Head: "feature/git", Base: "main", ChecksPassed: 12, ChecksPending: 1,
			Reviewers: []git.Reviewer{{Name: "Uika", State: git.ReviewApproved}, {Name: "Mutsumi", State: git.ReviewChangesRequested}, {Name: "Umiri"}},
			Added:     2412, Deleted: 310, Files: 46}
		if s.failed {
			pr.ChecksFailed = 1
		}
		if git.PullRequestCard(c, pr, git.PullRequestCardOptions{Disabled: s.disabled}).Changed() {
			log(fmt.Sprintf("Open #%d", pr.Number))
		}
	})
}

func gitReviewCommentDemo(c *ui.Context, log func(string)) {
	type state struct {
		thread   git.ReviewThread
		draft    string
		disabled bool
	}
	s := ui.Local(c.Root(), "316", func() state {
		return state{thread: git.ReviewThread{File: "commit_graph.go", Line: 88, Messages: []git.ReviewMessage{
			{Author: "Uika", When: gitDemoTime, Body: "Should merge lanes reuse the parent's color?"},
			{Author: "Sakiko", When: gitDemoTime.Add(time.Hour), Body: "Yes, the first parent keeps it."},
		}}}
	})
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		r := git.ReviewComment(c, s.thread, &s.draft, git.ReviewCommentOptions{Disabled: s.disabled})
		if text, ok := r.Replied(); ok {
			log("Reply: " + text)
			s.thread.Messages = append(s.thread.Messages, git.ReviewMessage{Author: "You", When: c.Now(), Body: text})
		}
		if r.Resolved() {
			log("Resolved")
			s.thread.Resolved = true
		}
		if r.Reopened() {
			log("Reopened")
			s.thread.Resolved = false
		}
	})
}

func gitDiffStatDemo(c *ui.Context, log func(string)) {
	type state struct{ added, deleted, blocks float64 }
	s := ui.Local(c.Root(), "317", func() state { return state{added: 120, deleted: 30, blocks: 5} })
	navigationDemo(c, func() {
		navigationSlider(c, &s.added, "Added", 0, 500, log)
		navigationSlider(c, &s.deleted, "Deleted", 0, 500, log)
		navigationSlider(c, &s.blocks, "Blocks", 1, 10, log)
	}, func() {
		git.DiffStat(c, int(s.added), int(s.deleted), git.DiffStatOptions{Blocks: int(s.blocks)})
	})
}

func gitStatusBadgeDemo(c *ui.Context, log func(string)) {
	type state struct{ labels bool }
	s := ui.Local(c.Root(), "318", func() state { return state{} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.labels, "Show labels", log)
	}, func() {
		ui.Row(c).Gap(16).Wrap().Children(func() {
			for st := git.GitModified; st <= git.GitCopied; st++ {
				git.GitStatusBadge(c, st, git.GitStatusBadgeOptions{ShowLabel: s.labels})
			}
		})
	})
}

func gitCommandBlockDemo(c *ui.Context, log func(string)) {
	type state struct {
		open, failed, running bool
	}
	s := ui.Local(c.Root(), "319", func() state { return state{open: true} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.failed, "Failed", log)
		navigationCheck(c, &s.running, "Running", log)
	}, func() {
		run := agent.CommandRun{Command: "go test -run Git ./...", Duration: 3214 * time.Millisecond,
			Output: "\x1b[32mok\x1b[0m  \tmujicaui\t3.214s\n\x1b[33mwarn\x1b[0m slow test: TestCommitGraphVirtualized"}
		if s.failed {
			run.ExitCode = 1
			run.Output = "\x1b[31m--- FAIL\x1b[0m: TestDiffViewer (0.02s)\n\x1b[1mFAIL\x1b[0m\tmujicaui\t0.401s"
		}
		run.Running = s.running
		r := git.CommandBlock(c, run, &s.open, git.CommandBlockOptions{})
		if r.Toggled() {
			log(fmt.Sprintf("Toggled: open=%v", s.open))
		}
		if what, ok := r.Copied(); ok {
			log("Copied " + what)
		}
	})
}

func gitAnsiTextDemo(c *ui.Context, log func(string)) {
	type state struct{ wrap bool }
	s := ui.Local(c.Root(), "320", func() state { return state{} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.wrap, "Wrap", log)
	}, func() {
		ui.Column(c).Padding(12).Background(core.Tokens(c).Surface).Radius(4).Children(func() {
			git.AnsiText(c, "\x1b[1mBuilding\x1b[0m mujica \x1b[2m(v1.1)\x1b[0m\n"+
				"\x1b[32m✓\x1b[0m 300 passed  \x1b[31m✗\x1b[0m 2 failed  \x1b[33m!\x1b[0m 1 skipped\n"+
				"\x1b[30m0\x1b[31m1\x1b[32m2\x1b[33m3\x1b[34m4\x1b[35m5\x1b[36m6\x1b[37m7\x1b[90m8\x1b[91m9\x1b[92mA\x1b[93mB\x1b[94mC\x1b[95mD\x1b[96mE\x1b[97mF\x1b[0m\n"+
				"\x1b[38;5;208m256-color\x1b[0m \x1b[38;2;125;32;52mtrue color\x1b[0m \x1b[4munderline\x1b[0m \x1b[3mitalic\x1b[0m \x1b[9mstruck\x1b[0m \x1b[7m inverse \x1b[0m\n"+
				"A long line that only breaks when wrapping is on, so the terminal layout of columns stays intact as it was printed.\n"+
				"Unknown escapes stay as text: \x1b[2K(cleared line)", git.AnsiTextOptions{Wrap: s.wrap})
		})
	})
}
