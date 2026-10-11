package git

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PullRequestState is where a pull request stands.
type PullRequestState int

const (
	PullRequestOpen PullRequestState = iota
	PullRequestDraft
	PullRequestMerged
	PullRequestClosed
)

// ReviewState is a reviewer's verdict.
type ReviewState int

const (
	ReviewPending ReviewState = iota
	ReviewApproved
	ReviewChangesRequested
)

// Reviewer is a person asked to review.
type Reviewer struct {
	Name  string
	State ReviewState
}

// PullRequest is what a PullRequestCard shows; the caller fetches it.
type PullRequest struct {
	Number     int
	Title      string
	Author     string
	State      PullRequestState
	Head, Base string
	// Checks counts CI results.
	ChecksPassed, ChecksFailed, ChecksPending int
	Reviewers                                 []Reviewer
	Added, Deleted, Files                     int
}

// PullRequestCardOptions configure a PullRequestCard.
type PullRequestCardOptions struct {
	Disabled bool
}

// PullRequestCard summarizes a pull request; Changed reports opening.
func PullRequestCard(c *ui.Context, pr PullRequest, opts PullRequestCardOptions) core.ChoiceResult {
	if pr.Number <= 0 || pr.Title == "" {
		panic("mujica: PullRequest needs a positive Number and a Title")
	}
	if pr.ChecksPassed < 0 || pr.ChecksFailed < 0 || pr.ChecksPending < 0 || pr.Files < 0 {
		panic("mujica: PullRequest counts must not be negative")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	card := ui.Column(c).AlignItems(ui.Stretch).Gap(u*2).Padding(u*4).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).Disabled(opts.Disabled)
	res := core.ChoiceResult{Element: card}
	card.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(u * 2).Children(func() {
			pullRequestStateBadge(c, pr.State)
			title := ui.ButtonBase(c).Justify(ui.Start).Gap(u*2).Padding(u, u).Radius(theme.ControlRadius).
				Role(ui.RoleLink).Label(fmt.Sprintf("#%d %s", pr.Number, pr.Title)).Shrink(1).MinWidth(0).Cursor(ui.CursorPointer)
			internal.Ring(title, k.Focus, theme.ControlRadius)
			if title.Hovered() {
				title.Background(k.SurfaceHover)
			}
			title.Transition(core.ColorFade)
			title.Children(func() {
				ui.Text(c, pr.Title).FontWeight(600).SingleLine().MinWidth(0)
				ui.Text(c, fmt.Sprintf("#%d", pr.Number)).TextColor(k.TextMuted).FontFeatures("tnum").Shrink(0)
			})
			if title.Clicked() && !title.IsDisabled() {
				res = core.NewChoiceResult(res.Element, true)
				slog.Debug("mujica: pull request card open", "number", pr.Number)
			}
		})
		ui.Row(c).AlignItems(ui.Center).Gap(u * 2).FontSize(small).TextColor(k.TextMuted).Wrap().Children(func() {
			ui.Text(c, pr.Author).FontWeight(600).TextColor(k.Text)
			if pr.Head != "" && pr.Base != "" {
				ui.Icon(c, icons.Must("git-pull-request")).Size(12, 12)
				ui.Text(c, pr.Head+" → "+pr.Base).Font(theme.MonoFont).SingleLine().MaxWidth(320)
			}
		})
		ui.Row(c).AlignItems(ui.Center).Gap(u * 4).Wrap().FontSize(small).Children(func() {
			pullRequestChecks(c, pr)
			if len(pr.Reviewers) > 0 {
				ui.Row(c).AlignItems(ui.Center).Gap(u * 2).Children(func() {
					for _, r := range pr.Reviewers {
						pullRequestReviewer(c, r)
					}
				})
			}
			ui.Spacer(c)
			ui.Text(c, fmt.Sprintf(core.Msg(c, "pull_request_card.files"), pr.Files)).TextColor(k.TextMuted).FontFeatures("tnum")
			DiffStat(c, pr.Added, pr.Deleted, DiffStatOptions{})
		})
	})
	return res
}

func pullRequestStateBadge(c *ui.Context, st PullRequestState) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	icon, col, key := "git-pull-request", k.Success, "open"
	switch st {
	case PullRequestOpen:
	case PullRequestDraft:
		icon, col, key = "git-pull-request", k.TextMuted, "draft"
	case PullRequestMerged:
		icon, col, key = "git-merge", k.AccentText, "merged"
	case PullRequestClosed:
		icon, col, key = "circle-x", k.Danger, "closed"
	default:
		panic("mujica: unknown PullRequestState")
	}
	ui.Row(c).AlignItems(ui.Center).Gap(u).Padding(2, u*2).Radius(theme.ControlRadius).Border(theme.BorderWidth, col).
		TextColor(col).Shrink(0).Children(func() {
		ui.Icon(c, icons.Must(icon)).Size(14, 14)
		ui.Text(c, core.Msg(c, "pull_request_card."+key)).FontSize(core.FontSize(c, theme.CaptionSize)).FontWeight(600).SingleLine()
	})
}

func pullRequestChecks(c *ui.Context, pr PullRequest) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	ui.Row(c).AlignItems(ui.Center).Gap(u * 3).Children(func() {
		for _, ch := range []struct {
			n         int
			icon, key string
			col       ui.Color
		}{
			{pr.ChecksPassed, "circle-check", "pull_request_card.passed", k.Success},
			{pr.ChecksFailed, "circle-x", "pull_request_card.failed", k.Danger},
			{pr.ChecksPending, "clock", "pull_request_card.pending", k.Warning},
		} {
			if ch.n == 0 {
				continue
			}
			name := fmt.Sprintf(core.Msg(c, ch.key), ch.n)
			ui.Row(c).AlignItems(ui.Center).Gap(u).TextColor(ch.col).Label(name).Tooltip(name).Children(func() {
				ui.Icon(c, icons.Must(ch.icon)).Size(14, 14)
				ui.Text(c, fmt.Sprint(ch.n)).FontFeatures("tnum")
			})
		}
	})
}

func pullRequestReviewer(c *ui.Context, r Reviewer) {
	k := core.Tokens(c)
	key, icon, col := "pull_request_card.review_pending", "clock", k.TextMuted
	switch r.State {
	case ReviewPending:
	case ReviewApproved:
		key, icon, col = "pull_request_card.approved", "check", k.Success
	case ReviewChangesRequested:
		key, icon, col = "pull_request_card.changes", "circle-alert", k.Danger
	default:
		panic("mujica: unknown ReviewState")
	}
	name := r.Name + ": " + core.Msg(c, key)
	ui.Box(c).Size(28, 28).Shrink(0).Label(name).Tooltip(name).Children(func() {
		display.Avatar(c, r.Name, display.AvatarOptions{Size: 24}).Absolute().Left(0).Top(0)
		ui.Box(c).Size(14, 14).Radius(7).Background(k.Surface).Center().Absolute().Right(0).Bottom(0).TextColor(col).
			Children(func() { ui.Icon(c, icons.Must(icon)).Size(10, 10) })
	})
}
