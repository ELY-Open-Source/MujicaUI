package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/git"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

var gitDemoTime = time.Date(2026, 10, 5, 9, 30, 0, 0, time.UTC)

const gitDemoOld = "package court\n\nimport \"fmt\"\n\n// Greet welcomes a guest.\nfunc Greet(name string) {\n\tfmt.Println(\"Welcome,\", name)\n}\n\nfunc Leave() {}\n\nfunc Dance() {}\n\nfunc Bow() {}\n"
const gitDemoNew = "package court\n\nimport \"fmt\"\n\n// Greet welcomes a guest by title.\nfunc Greet(title, name string) {\n\tfmt.Println(\"Welcome,\", title, name)\n}\n\nfunc Leave() {}\n\nfunc Dance() {}\n\nfunc Bow() {}\n\nfunc Curtsy() {}\n"

func gitDemoCommits(n int) []git.Commit {
	out := make([]git.Commit, n)
	authors := []string{"Sakiko", "Mutsumi", "Uika", "Umiri", "Nyamu"}
	for i := range out {
		out[i] = git.Commit{Hash: fmt.Sprintf("%016x%024x", uint64(n-i)*0x9E3779B97F4A7C15, n-i), Author: authors[i%5],
			When: gitDemoTime.Add(-time.Duration(i) * 37 * time.Minute), Subject: fmt.Sprintf("Chapter %d of the chronicle", n-i)}
		if i+1 < n {
			out[i].Parents = []string{fmt.Sprintf("%016x%024x", uint64(n-i-1)*0x9E3779B97F4A7C15, n-i-1)}
		}
	}
	if n > 0 {
		out[0].Refs = []string{"main", "origin/main"}
	}
	return out
}

func gitDemoGraph() []git.Commit {
	cs := []git.Commit{
		{Hash: "e7", Parents: []string{"e5", "f3"}, Subject: "Merge branch 'feature/lanes'", Refs: []string{"main"}},
		{Hash: "f3", Parents: []string{"f2"}, Subject: "Color lanes by branch"},
		{Hash: "e5", Parents: []string{"e4", "h1"}, Subject: "Merge hotfix"},
		{Hash: "h1", Parents: []string{"e3"}, Subject: "Fix the overture timing", Refs: []string{"hotfix"}},
		{Hash: "f2", Parents: []string{"f1"}, Subject: "Draw merge lines"},
		{Hash: "e4", Parents: []string{"e3"}, Subject: "Tune the ballroom palette"},
		{Hash: "f1", Parents: []string{"e3"}, Subject: "Start lane layout"},
		{Hash: "e3", Parents: []string{"e2"}, Subject: "Add the commit list"},
		{Hash: "e2", Parents: []string{"e1"}, Subject: "Add the blame view"},
		{Hash: "e1", Subject: "Initial commit", Refs: []string{"v1.0.0"}},
	}
	for i := range cs {
		cs[i].Author = []string{"Sakiko", "Mutsumi", "Uika"}[i%3]
		cs[i].When = gitDemoTime.Add(-time.Duration(i) * 2 * time.Hour)
	}
	return cs
}

func gitDemoBranches() []git.Branch {
	return []git.Branch{
		{Name: "main", Current: true, Upstream: "origin/main", Ahead: 2, Behind: 1},
		{Name: "feature/lanes", Upstream: "origin/feature/lanes", Ahead: 5},
		{Name: "fix/blame"},
		{Name: "origin/main", Remote: true},
		{Name: "origin/feature/lanes", Remote: true},
		{Name: "origin/release/1.1", Remote: true},
	}
}

func gitDiffViewerDemo(c *ui.Context, log func(string)) {
	type state struct {
		viewer       git.DiffViewerState
		split, plain bool
		context      float64
	}
	s := ui.Local(c.Root(), "301", func() state { return state{context: 2} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.split, "Side by side", log)
		navigationCheck(c, &s.plain, "Plain text", log)
		navigationSlider(c, &s.context, "Context", 1, 6, log)
	}, func() {
		lang := "go"
		if s.plain {
			lang = ""
		}
		git.DiffViewer(c, &s.viewer, gitDemoOld, gitDemoNew, git.DiffViewerOptions{SideBySide: s.split, Context: int(s.context), Language: lang}).Height(380)
	})
}

func gitInlineDiffDemo(c *ui.Context, log func(string)) {
	type state struct {
		from, to string
		mono     bool
	}
	s := ui.Local(c.Root(), "302", func() state {
		return state{from: "The court meets at noon in the east hall.", to: "The council meets at dusk in the east hall."}
	})
	navigationDemo(c, func() {
		navigationCheck(c, &s.mono, "Monospace", log)
	}, func() {
		if input.TextInput(c, &s.from, input.TextInputOptions{Label: "Old line"}).Changed() {
			log("Old: " + s.from)
		}
		if input.TextInput(c, &s.to, input.TextInputOptions{Label: "New line"}).Changed() {
			log("New: " + s.to)
		}
		git.InlineDiff(c, s.from, s.to, git.InlineDiffOptions{Mono: s.mono})
	})
}

func gitThreeWayMergeDemo(c *ui.Context, log func(string)) {
	type state struct {
		merge    git.ThreeWayMergeState
		disabled bool
	}
	s := ui.Local(c.Root(), "303", func() state { return state{} })
	base := "title: Overture\ntempo: 96\nkey: D minor\nlength: 4:10\nhall: east"
	ours := "title: Overture\ntempo: 104\nkey: D minor\nlength: 4:10\nhall: east\nencore: yes"
	theirs := "title: Overture\ntempo: 92\nkey: D minor\nlength: 4:30\nhall: east"
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		r := git.ThreeWayMerge(c, &s.merge, base, ours, theirs, git.ThreeWayMergeOptions{Disabled: s.disabled})
		r.Element.Height(360)
		if r.Changed() {
			done, total := s.merge.Resolved()
			log(fmt.Sprintf("Changed: %d of %d resolved", done, total))
		}
		ui.Text(c, s.merge.Result()).Font("monospace").TextColor(core.Tokens(c).TextMuted)
	})
}

func gitConflictResolverDemo(c *ui.Context, log func(string)) {
	type state struct {
		res      git.ConflictResolverState
		disabled bool
	}
	s := ui.Local(c.Root(), "304", func() state { return state{} })
	text := "# Program\n<<<<<<< HEAD\nOpening: Overture in D minor\n=======\nOpening: Prelude in G major\n>>>>>>> feature/program\nInterval: 15 minutes\n<<<<<<< HEAD\nFinale: Masquerade\n=======\nFinale: Ave Mujica\n>>>>>>> feature/program"
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		if git.ConflictResolver(c, &s.res, text, git.ConflictResolverOptions{Disabled: s.disabled}).Changed() {
			done, total := s.res.Resolved()
			log(fmt.Sprintf("Changed: %d of %d resolved", done, total))
		}
	})
}

func gitChangesListDemo(c *ui.Context, log func(string)) {
	type state struct {
		list     git.ChangesListState
		files    []git.ChangedFile
		disabled bool
	}
	s := ui.Local(c.Root(), "305", func() state {
		return state{files: []git.ChangedFile{
			{Path: "commit_graph.go", Status: git.GitAdded, Staged: true},
			{Path: "internal/git.go", Status: git.GitModified},
			{Path: "docs/old.md", Status: git.GitDeleted},
			{Path: "notes.txt", Status: git.GitUntracked},
			{Path: "icons/svg/tag.svg", Status: git.GitRenamed},
		}}
	})
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		r := git.ChangesList(c, &s.list, s.files, git.ChangesListOptions{Disabled: s.disabled})
		r.Element.Height(320)
		if a, path, ok := r.Action(); ok {
			log("Action: " + a + " " + path)
			s.files = gitStage(s.files, a, path)
		}
		if r.Changed() && s.list.Selected() >= 0 {
			log("Selected: " + s.files[s.list.Selected()].Path)
		}
	})
}

// gitStage applies a stage action to a copy of files.
func gitStage(files []git.ChangedFile, action, path string) []git.ChangedFile {
	out := append([]git.ChangedFile(nil), files...)
	for i := range out {
		switch {
		case action == "stage_all":
			out[i].Staged = true
		case action == "unstage_all":
			out[i].Staged = false
		case out[i].Path == path:
			out[i].Staged = action == "stage"
		}
	}
	return out
}

func gitCommitInputDemo(c *ui.Context, log func(string)) {
	type state struct {
		msg            git.CommitMessage
		busy, disabled bool
		until          time.Time
	}
	s := ui.Local(c.Root(), "306", func() state { return state{} })
	if s.busy && c.Now().After(s.until) {
		s.busy = false
		s.msg = git.CommitMessage{}
		log("Committed")
	}
	if s.busy {
		c.After(s.until.Sub(c.Now()))
	}
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		r := git.CommitInput(c, &s.msg, git.CommitInputOptions{Busy: s.busy, Disabled: s.disabled})
		if r.Committed() || r.Amended() {
			log(fmt.Sprintf("Submit: amend=%v title=%q", r.Amended(), s.msg.Title))
			s.busy, s.until = true, c.Now().Add(1200*time.Millisecond)
		}
	})
}

func gitCommitGraphDemo(c *ui.Context, log func(string)) {
	type state struct {
		graph git.CommitGraphState
		big   bool
	}
	s := ui.Local(c.Root(), "307", func() state { return state{} })
	small, big := gitDemoGraph(), gitDemoCommits(100000)
	navigationDemo(c, func() {
		navigationCheck(c, &s.big, "100,000 commits", log)
	}, func() {
		commits := small
		if s.big {
			commits = big
		}
		r := git.CommitGraph(c, &s.graph, commits, git.CommitGraphOptions{})
		r.Element.Height(380)
		if hash, ok := s.graph.List.Selected(); ok && (r.Changed() || r.Submitted()) {
			log(fmt.Sprintf("Selected %s, open=%v", hash[:min(7, len(hash))], r.Submitted()))
		}
	})
}

func gitCommitListDemo(c *ui.Context, log func(string)) {
	type state struct {
		list    data.ListState[string]
		loading bool
	}
	s := ui.Local(c.Root(), "308", func() state { return state{} })
	commits := gitDemoCommits(40)
	navigationDemo(c, func() {
		navigationCheck(c, &s.loading, "Loading more", log)
	}, func() {
		st := data.DataReady
		if s.loading {
			st = data.DataLoadingMore
		}
		r := git.CommitList(c, &s.list, commits, git.CommitListOptions{Status: st})
		r.Element.Height(360)
		if hash, ok := s.list.Selected(); ok && (r.Changed() || r.Submitted()) {
			log(fmt.Sprintf("Selected %s, open=%v", hash[:7], r.Submitted()))
		}
	})
}

func gitBranchSelectorDemo(c *ui.Context, log func(string)) {
	type state struct {
		branch           string
		create, disabled bool
		branches         []git.Branch
	}
	s := ui.Local(c.Root(), "309", func() state { return state{branch: "main", create: true, branches: gitDemoBranches()} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.create, "Allow create", log)
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		r := git.BranchSelector(c, &s.branch, s.branches, git.BranchSelectorOptions{AllowCreate: s.create, Disabled: s.disabled})
		if r.Changed() {
			log("Changed: " + s.branch)
		}
		if name, ok := r.Created(); ok {
			log("Create: " + name)
			s.branches = append(append([]git.Branch(nil), s.branches...), git.Branch{Name: name})
			s.branch = name
		}
	})
}

func gitBranchListDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "310", func() data.ListState[string] { return data.ListState[string]{} })
	actions := []git.GitAction{{ID: "checkout", Label: "Check out", Icon: "git-branch"}, {ID: "delete", Label: "Delete", Icon: "trash-2"}}
	r := git.BranchList(c, s, gitDemoBranches(), git.BranchListOptions{Actions: actions})
	r.Element.Height(300)
	if a, name, ok := r.Action(); ok {
		log("Action: " + a + " " + name)
	}
	if key, ok := s.Selected(); ok && r.Changed() {
		log("Selected: " + key)
	}
}

func gitTagListDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "311", func() data.ListState[string] { return data.ListState[string]{} })
	tags := []git.GitTag{
		{Name: "v1.1.0", Target: "9e3779b97f4a7c15", Annotation: "Five hundred components"},
		{Name: "v1.0.0", Target: "3c6ef372fe94f82a", Annotation: "The first hundred"},
		{Name: "nightly", Target: "daa66d2c7ddf743f"},
	}
	r := git.TagList(c, s, tags, git.TagListOptions{})
	r.Element.Height(220)
	if _, hash, ok := r.Action(); ok {
		log("Jump to " + hash[:7])
	}
	if key, ok := s.Selected(); ok && r.Changed() {
		log("Selected: " + key)
	}
}

func gitStashListDemo(c *ui.Context, log func(string)) {
	type state struct {
		list    data.ListState[string]
		stashes []git.Stash
	}
	s := ui.Local(c.Root(), "312", func() state {
		return state{stashes: []git.Stash{
			{Ref: "stash@{0}", Message: "WIP on main: lane colors", Branch: "main"},
			{Ref: "stash@{1}", Message: "Try serif headings", Branch: "feature/type"},
			{Ref: "stash@{2}", Message: "Blame hover experiment", Branch: "fix/blame"},
		}}
	})
	r := git.StashList(c, &s.list, s.stashes, git.StashListOptions{})
	r.Element.Height(240)
	if a, ref, ok := r.Action(); ok {
		log("Action: " + a + " " + ref)
		if a != "apply" {
			var keep []git.Stash
			for _, st := range s.stashes {
				if st.Ref != ref {
					keep = append(keep, st)
				}
			}
			s.stashes = keep
		}
	}
}

func gitBlameViewDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "313", func() git.BlameViewState { return git.BlameViewState{} })
	var lines []git.BlameLine
	for i, l := range strings.Split(strings.TrimSuffix(gitDemoNew, "\n"), "\n") {
		who := []git.BlameLine{
			{Hash: "9e3779b97f4a7c15", Author: "Sakiko", When: gitDemoTime, Summary: "Create the court package"},
			{Hash: "3c6ef372fe94f82a", Author: "Uika", When: gitDemoTime.Add(-48 * time.Hour), Summary: "Greet guests by title"},
		}[min(i/4%2, 1)]
		who.Text = l
		lines = append(lines, who)
	}
	git.BlameView(c, s, lines, git.BlameViewOptions{Language: "go"}).Height(380)
	ui.Text(c, "Hover the left column for commit details.").TextColor(core.Tokens(c).TextMuted)
}

func gitFileHistoryDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "314", func() data.ListState[string] { return data.ListState[string]{} })
	r := git.FileHistory(c, s, "court/greet.go", gitDemoCommits(12), git.FileHistoryOptions{})
	r.Element.Height(360)
	if _, hash, ok := r.Action(); ok {
		log("Diff at " + hash[:7])
	}
	if key, ok := s.Selected(); ok && r.Changed() {
		log("Selected: " + key[:7])
	}
}
