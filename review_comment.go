package mujica

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// ReviewMessage is one comment of a review thread.
type ReviewMessage struct {
	Author string
	When   time.Time
	Body   string
}

// ReviewThread is a line-level review discussion.
type ReviewThread struct {
	File     string
	Line     int
	Resolved bool
	Messages []ReviewMessage
}

// ReviewCommentOptions configure a ReviewComment.
type ReviewCommentOptions struct {
	// TimeFormat lays out message times; empty means "2006-01-02 15:04".
	TimeFormat string
	Disabled   bool
}

// ReviewCommentResult is a built ReviewComment and its events.
type ReviewCommentResult struct {
	Element            *ui.Element
	reply              string
	resolved, reopened bool
}

// Replied returns the reply sent this frame.
func (r ReviewCommentResult) Replied() (string, bool) { return r.reply, r.reply != "" }

// Resolved reports the Resolve button this frame.
func (r ReviewCommentResult) Resolved() bool { return r.resolved }

// Reopened reports the Reopen button this frame.
func (r ReviewCommentResult) Reopened() bool { return r.reopened }

// ReviewComment shows a line's review thread; *draft is the reply.
func ReviewComment(c *ui.Context, thread ReviewThread, draft *string, opts ReviewCommentOptions) ReviewCommentResult {
	if thread.File == "" || thread.Line <= 0 || len(thread.Messages) == 0 {
		panic("mujica: ReviewThread needs a File, a positive Line and a message")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	small := FontSize(c, theme.CaptionSize)
	layout := gitTimeLayout(opts.TimeFormat)
	card := ui.Column(c).AlignItems(ui.Stretch).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Radius(theme.CardRadius).Disabled(opts.Disabled).Clip()
	res := ReviewCommentResult{Element: card}
	card.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(u*2).Padding(u*2, u*3).Background(k.Background).FontSize(small).Children(func() {
			ui.Icon(c, icons.Must("message-square")).Size(14, 14).TextColor(k.TextMuted)
			ui.Text(c, fmt.Sprintf("%s:%d", thread.File, thread.Line)).Font(theme.MonoFont).SingleLine().Grow(1).MinWidth(0)
			if thread.Resolved {
				ui.Row(c).AlignItems(ui.Center).Gap(u).TextColor(k.Success).Children(func() {
					ui.Icon(c, icons.Must("circle-check")).Size(14, 14)
					ui.Text(c, msg(c, "review_comment.resolved")).FontWeight(600)
				})
			}
		})
		for i, m := range thread.Messages {
			ui.Column(c).Gap(u).Padding(u*3, u*4).Children(func() {
				if i > 0 {
					ui.Box(c).Height(1).Background(k.Border).Margin(0, 0, u*2, 0).Role(ui.RoleNone)
				}
				ui.Row(c).AlignItems(ui.Center).Gap(u * 2).Children(func() {
					Avatar(c, m.Author, AvatarOptions{Size: 20})
					ui.Text(c, m.Author).FontWeight(600).SingleLine()
					ui.Text(c, m.When.Format(layout)).FontSize(small).TextColor(k.TextMuted).FontFeatures("tnum")
				})
				ui.Text(c, m.Body)
			})
		}
		ui.Box(c).Height(1).Background(k.Border).Role(ui.RoleNone)
		ui.Column(c).Gap(u * 2).Padding(u * 3).AlignItems(ui.Stretch).Children(func() {
			if !thread.Resolved {
				in := TextArea(c, draft, TextAreaOptions{Placeholder: msg(c, "review_comment.placeholder"), Label: msg(c, "review_comment.field"),
					MinHeight: 56, MaxHeight: 160, Disabled: opts.Disabled})
				send := in.Submitted()
				ui.Row(c).Gap(u * 2).Justify(ui.End).Children(func() {
					if Button(c, msg(c, "review_comment.resolve"), ButtonOptions{Variant: Secondary, Icon: icons.Must("check"), Disabled: opts.Disabled}).Clicked() {
						res.resolved = true
					}
					if Button(c, msg(c, "review_comment.reply"), ButtonOptions{Icon: icons.Must("reply"), Disabled: opts.Disabled || strings.TrimSpace(*draft) == ""}).Clicked() {
						send = true
					}
				})
				if send && !card.IsDisabled() && strings.TrimSpace(*draft) != "" {
					res.reply, *draft = strings.TrimSpace(*draft), ""
				}
				return
			}
			ui.Row(c).Justify(ui.End).Children(func() {
				if Button(c, msg(c, "review_comment.reopen"), ButtonOptions{Variant: Secondary, Icon: icons.Must("rotate-ccw"), Disabled: opts.Disabled}).Clicked() {
					res.reopened = true
				}
			})
		})
	})
	if card.IsDisabled() {
		res.resolved, res.reopened = false, false
	}
	switch {
	case res.reply != "":
		slog.Debug("mujica: review comment reply", "file", thread.File, "line", thread.Line)
	case res.resolved:
		slog.Debug("mujica: review comment resolve", "file", thread.File, "line", thread.Line)
	case res.reopened:
		slog.Debug("mujica: review comment reopen", "file", thread.File, "line", thread.Line)
	}
	return res
}
