package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func pullRequestSample() PullRequest {
	return PullRequest{Number: 318, Title: "Git components 301–320", Author: "Sakiko", State: PullRequestOpen,
		Head: "feature/git", Base: "main", ChecksPassed: 12, ChecksFailed: 1, ChecksPending: 2,
		Reviewers: []Reviewer{{Name: "Uika", State: ReviewApproved}, {Name: "Mutsumi", State: ReviewChangesRequested}, {Name: "Umiri"}},
		Added:     2412, Deleted: 310, Files: 46}
}

func TestPullRequestCardShowsAndOpens(t *testing.T) {
	opens := 0
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if PullRequestCard(c, pullRequestSample(), PullRequestCardOptions{Disabled: disabled}).Changed() {
			opens++
		}
	}, 700, 240)
	for _, name := range []string{"Open", "#318", "feature/git → main", "12 checks passed", "1 checks failed", "2 checks pending",
		"Uika: approved", "Mutsumi: changes requested", "Umiri: review pending", "46 files", "2412 additions, 310 deletions"} {
		if _, ok := tt.Find(name); !ok {
			t.Errorf("missing %q", name)
		}
	}
	if err := tt.Click("Git components 301–320"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnter)
	if opens != 2 {
		t.Fatalf("opens %d", opens)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Git components 301–320")
	if opens != 2 {
		t.Fatal("disabled card opened")
	}
}

func TestPullRequestCardStatesAndInvalid(t *testing.T) {
	for st, word := range map[PullRequestState]string{PullRequestDraft: "Draft", PullRequestMerged: "Merged", PullRequestClosed: "Closed"} {
		pr := pullRequestSample()
		pr.State = st
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			PullRequestCard(c, pr, PullRequestCardOptions{})
		}, 700, 240)
		if !tt.HasText(word) {
			t.Errorf("state %d: %q", st, tt.Texts())
		}
	}
	for _, pr := range []PullRequest{{Title: "x"}, {Number: 1}, {Number: 1, Title: "x", Files: -1}, {Number: 1, Title: "x", State: 9}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v did not panic", pr)
				}
			}()
			ui.Render(func(c *ui.Context) {
				Use(c, Settings{})
				PullRequestCard(c, pr, PullRequestCardOptions{})
			}, 300, 100, 1)
		}()
	}
}

func TestPullRequestCardSnapshot(t *testing.T) {
	snapshot(t, "pull_request_card", 640, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			PullRequestCard(c, pullRequestSample(), PullRequestCardOptions{})
			pr := pullRequestSample()
			pr.Number, pr.State, pr.Title, pr.Reviewers, pr.ChecksFailed, pr.ChecksPending = 301, PullRequestMerged, "Diff viewer", nil, 0, 0
			PullRequestCard(c, pr, PullRequestCardOptions{})
		})
	})
}
