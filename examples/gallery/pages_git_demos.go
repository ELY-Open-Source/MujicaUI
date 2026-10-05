package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

var gitDemoTime = time.Date(2026, 10, 5, 9, 30, 0, 0, time.UTC)

const gitDemoOld = "package court\n\nimport \"fmt\"\n\n// Greet welcomes a guest.\nfunc Greet(name string) {\n\tfmt.Println(\"Welcome,\", name)\n}\n\nfunc Leave() {}\n\nfunc Dance() {}\n\nfunc Bow() {}\n"
const gitDemoNew = "package court\n\nimport \"fmt\"\n\n// Greet welcomes a guest by title.\nfunc Greet(title, name string) {\n\tfmt.Println(\"Welcome,\", title, name)\n}\n\nfunc Leave() {}\n\nfunc Dance() {}\n\nfunc Bow() {}\n\nfunc Curtsy() {}\n"

func gitDemoCommits(n int) []mujica.Commit {
	out := make([]mujica.Commit, n)
	authors := []string{"Sakiko", "Mutsumi", "Uika", "Umiri", "Nyamu"}
	for i := range out {
		out[i] = mujica.Commit{Hash: fmt.Sprintf("%016x%024x", uint64(n-i)*0x9E3779B97F4A7C15, n-i), Author: authors[i%5],
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

func gitDemoGraph() []mujica.Commit {
	cs := []mujica.Commit{
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

func gitDemoBranches() []mujica.Branch {
	return []mujica.Branch{
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
		viewer       mujica.DiffViewerState
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
		mujica.DiffViewer(c, &s.viewer, gitDemoOld, gitDemoNew, mujica.DiffViewerOptions{SideBySide: s.split, Context: int(s.context), Language: lang}).Height(380)
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
		if mujica.TextInput(c, &s.from, mujica.TextInputOptions{Label: "Old line"}).Changed() {
			log("Old: " + s.from)
		}
		if mujica.TextInput(c, &s.to, mujica.TextInputOptions{Label: "New line"}).Changed() {
			log("New: " + s.to)
		}
		mujica.InlineDiff(c, s.from, s.to, mujica.InlineDiffOptions{Mono: s.mono})
	})
}

func gitThreeWayMergeDemo(c *ui.Context, log func(string)) {
	type state struct {
		merge    mujica.ThreeWayMergeState
		disabled bool
	}
	s := ui.Local(c.Root(), "303", func() state { return state{} })
	base := "title: Overture\ntempo: 96\nkey: D minor\nlength: 4:10\nhall: east"
	ours := "title: Overture\ntempo: 104\nkey: D minor\nlength: 4:10\nhall: east\nencore: yes"
	theirs := "title: Overture\ntempo: 92\nkey: D minor\nlength: 4:30\nhall: east"
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		r := mujica.ThreeWayMerge(c, &s.merge, base, ours, theirs, mujica.ThreeWayMergeOptions{Disabled: s.disabled})
		r.Element.Height(360)
		if r.Changed() {
			done, total := s.merge.Resolved()
			log(fmt.Sprintf("Changed: %d of %d resolved", done, total))
		}
		ui.Text(c, s.merge.Result()).Font("monospace").TextColor(mujica.Tokens(c).TextMuted)
	})
}

func gitConflictResolverDemo(c *ui.Context, log func(string)) {
	type state struct {
		res      mujica.ConflictResolverState
		disabled bool
	}
	s := ui.Local(c.Root(), "304", func() state { return state{} })
	text := "# Program\n<<<<<<< HEAD\nOpening: Overture in D minor\n=======\nOpening: Prelude in G major\n>>>>>>> feature/program\nInterval: 15 minutes\n<<<<<<< HEAD\nFinale: Masquerade\n=======\nFinale: Ave Mujica\n>>>>>>> feature/program"
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		if mujica.ConflictResolver(c, &s.res, text, mujica.ConflictResolverOptions{Disabled: s.disabled}).Changed() {
			done, total := s.res.Resolved()
			log(fmt.Sprintf("Changed: %d of %d resolved", done, total))
		}
	})
}

func gitChangesListDemo(c *ui.Context, log func(string)) {
	type state struct {
		list     mujica.ChangesListState
		files    []mujica.ChangedFile
		disabled bool
	}
	s := ui.Local(c.Root(), "305", func() state {
		return state{files: []mujica.ChangedFile{
			{Path: "commit_graph.go", Status: mujica.GitAdded, Staged: true},
			{Path: "internal/git.go", Status: mujica.GitModified},
			{Path: "docs/old.md", Status: mujica.GitDeleted},
			{Path: "notes.txt", Status: mujica.GitUntracked},
			{Path: "icons/svg/tag.svg", Status: mujica.GitRenamed},
		}}
	})
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		r := mujica.ChangesList(c, &s.list, s.files, mujica.ChangesListOptions{Disabled: s.disabled})
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
func gitStage(files []mujica.ChangedFile, action, path string) []mujica.ChangedFile {
	out := append([]mujica.ChangedFile(nil), files...)
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
		msg            mujica.CommitMessage
		busy, disabled bool
		until          time.Time
	}
	s := ui.Local(c.Root(), "306", func() state { return state{} })
	if s.busy && c.Now().After(s.until) {
		s.busy = false
		s.msg = mujica.CommitMessage{}
		log("Committed")
	}
	if s.busy {
		c.After(s.until.Sub(c.Now()))
	}
	navigationDemo(c, func() {
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		r := mujica.CommitInput(c, &s.msg, mujica.CommitInputOptions{Busy: s.busy, Disabled: s.disabled})
		if r.Committed() || r.Amended() {
			log(fmt.Sprintf("Submit: amend=%v title=%q", r.Amended(), s.msg.Title))
			s.busy, s.until = true, c.Now().Add(1200*time.Millisecond)
		}
	})
}

func gitCommitGraphDemo(c *ui.Context, log func(string)) {
	type state struct {
		graph mujica.CommitGraphState
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
		r := mujica.CommitGraph(c, &s.graph, commits, mujica.CommitGraphOptions{})
		r.Element.Height(380)
		if hash, ok := s.graph.List.Selected(); ok && (r.Changed() || r.Submitted()) {
			log(fmt.Sprintf("Selected %s, open=%v", hash[:min(7, len(hash))], r.Submitted()))
		}
	})
}

func gitCommitListDemo(c *ui.Context, log func(string)) {
	type state struct {
		list    mujica.ListState[string]
		loading bool
	}
	s := ui.Local(c.Root(), "308", func() state { return state{} })
	commits := gitDemoCommits(40)
	navigationDemo(c, func() {
		navigationCheck(c, &s.loading, "Loading more", log)
	}, func() {
		st := mujica.DataReady
		if s.loading {
			st = mujica.DataLoadingMore
		}
		r := mujica.CommitList(c, &s.list, commits, mujica.CommitListOptions{Status: st})
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
		branches         []mujica.Branch
	}
	s := ui.Local(c.Root(), "309", func() state { return state{branch: "main", create: true, branches: gitDemoBranches()} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.create, "Allow create", log)
		navigationCheck(c, &s.disabled, "Disabled", log)
	}, func() {
		r := mujica.BranchSelector(c, &s.branch, s.branches, mujica.BranchSelectorOptions{AllowCreate: s.create, Disabled: s.disabled})
		if r.Changed() {
			log("Changed: " + s.branch)
		}
		if name, ok := r.Created(); ok {
			log("Create: " + name)
			s.branches = append(append([]mujica.Branch(nil), s.branches...), mujica.Branch{Name: name})
			s.branch = name
		}
	})
}

func gitBranchListDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "310", func() mujica.ListState[string] { return mujica.ListState[string]{} })
	actions := []mujica.GitAction{{ID: "checkout", Label: "Check out", Icon: "git-branch"}, {ID: "delete", Label: "Delete", Icon: "trash-2"}}
	r := mujica.BranchList(c, s, gitDemoBranches(), mujica.BranchListOptions{Actions: actions})
	r.Element.Height(300)
	if a, name, ok := r.Action(); ok {
		log("Action: " + a + " " + name)
	}
	if key, ok := s.Selected(); ok && r.Changed() {
		log("Selected: " + key)
	}
}

func gitTagListDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "311", func() mujica.ListState[string] { return mujica.ListState[string]{} })
	tags := []mujica.GitTag{
		{Name: "v1.1.0", Target: "9e3779b97f4a7c15", Annotation: "Five hundred components"},
		{Name: "v1.0.0", Target: "3c6ef372fe94f82a", Annotation: "The first hundred"},
		{Name: "nightly", Target: "daa66d2c7ddf743f"},
	}
	r := mujica.TagList(c, s, tags, mujica.TagListOptions{})
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
		list    mujica.ListState[string]
		stashes []mujica.Stash
	}
	s := ui.Local(c.Root(), "312", func() state {
		return state{stashes: []mujica.Stash{
			{Ref: "stash@{0}", Message: "WIP on main: lane colors", Branch: "main"},
			{Ref: "stash@{1}", Message: "Try serif headings", Branch: "feature/type"},
			{Ref: "stash@{2}", Message: "Blame hover experiment", Branch: "fix/blame"},
		}}
	})
	r := mujica.StashList(c, &s.list, s.stashes, mujica.StashListOptions{})
	r.Element.Height(240)
	if a, ref, ok := r.Action(); ok {
		log("Action: " + a + " " + ref)
		if a != "apply" {
			var keep []mujica.Stash
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
	s := ui.Local(c.Root(), "313", func() mujica.BlameViewState { return mujica.BlameViewState{} })
	var lines []mujica.BlameLine
	for i, l := range strings.Split(strings.TrimSuffix(gitDemoNew, "\n"), "\n") {
		who := []mujica.BlameLine{
			{Hash: "9e3779b97f4a7c15", Author: "Sakiko", When: gitDemoTime, Summary: "Create the court package"},
			{Hash: "3c6ef372fe94f82a", Author: "Uika", When: gitDemoTime.Add(-48 * time.Hour), Summary: "Greet guests by title"},
		}[min(i/4%2, 1)]
		who.Text = l
		lines = append(lines, who)
	}
	mujica.BlameView(c, s, lines, mujica.BlameViewOptions{Language: "go"}).Height(380)
	ui.Text(c, "Hover the left column for commit details.").TextColor(mujica.Tokens(c).TextMuted)
}

func gitFileHistoryDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "314", func() mujica.ListState[string] { return mujica.ListState[string]{} })
	r := mujica.FileHistory(c, s, "court/greet.go", gitDemoCommits(12), mujica.FileHistoryOptions{})
	r.Element.Height(360)
	if _, hash, ok := r.Action(); ok {
		log("Diff at " + hash[:7])
	}
	if key, ok := s.Selected(); ok && r.Changed() {
		log("Selected: " + key[:7])
	}
}
