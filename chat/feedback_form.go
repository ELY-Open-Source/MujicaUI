package chat

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// FeedbackRating is the thumb given; zero is FeedbackUp.
type FeedbackRating int

const (
	FeedbackUp FeedbackRating = iota
	FeedbackDown
)

// FeedbackValue is the chosen reasons and a comment.
type FeedbackValue struct {
	Reasons []string
	Comment string
}

// FeedbackFormOptions configure a FeedbackForm.
type FeedbackFormOptions struct {
	Rating FeedbackRating
	// Reasons replaces the built-in reasons for the rating.
	Reasons []string
	// Submitting shows the caller sending; it blocks the form.
	Submitting bool
	Disabled   bool
}

// FeedbackFormView is a built FeedbackForm and this frame's outcome.
type FeedbackFormView struct {
	Element              ui.Element
	Submitted, Cancelled bool
}

var feedbackReasons = [2][]string{
	{"feedback_form.up.accurate", "feedback_form.up.helpful", "feedback_form.up.clear"},
	{"feedback_form.down.wrong", "feedback_form.down.unhelpful", "feedback_form.down.unsafe", "feedback_form.down.style"},
}

// FeedbackForm asks for reasons and a comment after a rating.
func FeedbackForm(c *ui.Context, value *FeedbackValue, opts FeedbackFormOptions) FeedbackFormView {
	if value == nil {
		panic("mujica: FeedbackForm needs a value")
	}
	if opts.Rating != FeedbackUp && opts.Rating != FeedbackDown {
		panic("mujica: FeedbackForm Rating out of range")
	}
	reasons := opts.Reasons
	if reasons == nil {
		for _, key := range feedbackReasons[opts.Rating] {
			reasons = append(reasons, core.Msg(c, key))
		}
	}
	k := core.Tokens(c)
	d := core.Density(c)
	blocked := opts.Disabled || opts.Submitting
	col := ui.Column(c).Gap(d.Unit()*3).Padding(d.Unit()*4).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Radius(theme.CardRadius).Role(ui.RoleGroup).Disabled(blocked)
	v := FeedbackFormView{Element: col}
	title := map[FeedbackRating]string{FeedbackUp: "feedback_form.up.title", FeedbackDown: "feedback_form.down.title"}[opts.Rating]
	col.Label(core.Msg(c, title))
	col.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, core.Msg(c, title)).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).TextColor(k.Text)
		})
		ui.Row(c).Wrap().Gap(d.Unit() * 2).FocusGroup(ui.Horizontal).Children(func() {
			for _, r := range reasons {
				on := slices.Contains(value.Reasons, r)
				was := on
				input.Toggle(c, &on, r, input.ToggleOptions{Disabled: blocked})
				if on != was && !col.IsDisabled() {
					if on {
						value.Reasons = append(slices.Clone(value.Reasons), r)
					} else {
						value.Reasons = slices.DeleteFunc(slices.Clone(value.Reasons), func(x string) bool { return x == r })
					}
				}
			}
		})
		input.TextArea(c, &value.Comment, input.TextAreaOptions{Placeholder: core.Msg(c, "feedback_form.comment"), Label: core.Msg(c, "feedback_form.comment"),
			MinHeight: 72, Disabled: blocked})
		ui.Row(c).Gap(d.Unit() * 2).Justify(ui.End).Children(func() {
			v.Cancelled = input.Button(c, core.Msg(c, "feedback_form.cancel"), input.ButtonOptions{Variant: input.Ghost, Disabled: blocked}).Clicked()
			v.Submitted = input.Button(c, core.Msg(c, "feedback_form.submit"), input.ButtonOptions{Icon: icons.Must("check"), Loading: opts.Submitting, Disabled: opts.Disabled}).Clicked()
		})
	})
	if !col.IsDisabled() && col.Shortcut(0, ui.KeyEscape) {
		v.Cancelled = true
	}
	if col.IsDisabled() {
		v.Submitted, v.Cancelled = false, false
	}
	switch {
	case v.Submitted:
		slog.Debug("mujica: feedback_form submit", "rating", opts.Rating, "reasons", len(value.Reasons))
	case v.Cancelled:
		slog.Debug("mujica: feedback_form cancel")
	}
	return v
}
