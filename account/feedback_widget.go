package account

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/navigation"
)

// FeedbackKind is what the feedback is about.
type FeedbackKind int

const (
	FeedbackBug FeedbackKind = iota
	FeedbackIdea
	FeedbackPraise
	FeedbackOther
)

// Feedback is what a FeedbackWidget sends.
type Feedback struct {
	Kind FeedbackKind
	Text string
}

// FeedbackWidgetState holds the draft and its form state.
type FeedbackWidgetState struct {
	Draft Feedback
	Form  input.FormState
}

// EndSubmit clears successful drafts and preserves failed ones.
func (s *FeedbackWidgetState) EndSubmit(err string) {
	s.Form.EndSubmit(err)
	if err == "" {
		s.Draft.Text = ""
	}
}

// FeedbackWidgetOptions configure a FeedbackWidget.
type FeedbackWidgetOptions struct {
	// Attachments builds the caller-provided screenshot or file slot.
	Attachments func()
	// MinLength is the fewest characters accepted; zero means 10.
	MinLength int
	Disabled  bool
}

// FeedbackWidgetResult carries the submitted feedback.
type FeedbackWidgetResult struct {
	Element   ui.Element
	feedback  Feedback
	submitted bool
}

// Submitted returns feedback; finish the request with EndSubmit.
func (r FeedbackWidgetResult) Submitted() (Feedback, bool) { return r.feedback, r.submitted }

// FeedbackWidget validates feedback and submits it once.
func FeedbackWidget(c *ui.Context, s *FeedbackWidgetState, opts FeedbackWidgetOptions) FeedbackWidgetResult {
	if opts.MinLength < 0 {
		panic("mujica: FeedbackWidget MinLength cannot be negative")
	}
	if s.Draft.Kind < FeedbackBug || s.Draft.Kind > FeedbackOther || opts.MinLength > 2000 {
		panic("mujica: FeedbackWidget invalid kind or length")
	}
	minLen := opts.MinLength
	if minLen == 0 {
		minLen = 10
	}
	d := core.Density(c)
	var r FeedbackWidgetResult
	r.Element = ui.Column(c).Gap(d.Unit() * 4).Width(420).Disabled(opts.Disabled)
	disabled := r.Element.IsDisabled()
	if disabled {
		frozen := *s
		frozen.Form.Order = append([]string(nil), s.Form.Order...)
		frozen.Form.Fields = make(map[string]*input.FormFieldState, len(s.Form.Fields))
		for name, field := range s.Form.Fields {
			copy := *field
			frozen.Form.Fields[name] = &copy
		}
		frozen.Form.Requested, frozen.Form.Awaiting = false, false
		s = &frozen
	}
	locked := disabled || s.Form.Submitting()
	r.Element.Children(func() {
		kind := int(s.Draft.Kind)
		var segs []navigation.Segment
		for f := FeedbackBug; f <= FeedbackOther; f++ {
			segs = append(segs, navigation.Segment{Label: core.Msg(c, fmt.Sprintf("feedback_widget.kind%d", f))})
		}
		navigation.SegmentedControl(c, &kind, segs, navigation.SegmentedControlOptions{Equal: true, Label: core.Msg(c, "feedback_widget.kind"), Disabled: locked})
		if !locked {
			s.Draft.Kind = FeedbackKind(kind)
		}
		res := input.Form(c, &s.Form, input.FormOptions{SubmitLabel: core.Msg(c, "feedback_widget.send"), ResetLabel: core.Msg(c, "feedback_widget.clear")}, func(f *input.FormScope) {
			input.FormItem(c, f, "text", &s.Draft.Text, input.FormItemOptions[string]{Label: core.Msg(c, "feedback_widget.message"), Required: true,
				Validate: func(v string) string {
					if n := len([]rune(strings.TrimSpace(v))); n < minLen {
						return fmt.Sprintf(core.Msg(c, "feedback_widget.short"), minLen)
					}
					return ""
				}}, func() ui.Element {
				return input.TextArea(c, &s.Draft.Text, input.TextAreaOptions{MinHeight: 96, MaxLength: 2000, Disabled: locked}).Input
			})
			if opts.Attachments != nil {
				ui.Column(c).Gap(d.Unit()).Disabled(locked).Children(opts.Attachments)
			}
		})
		if res.Submitted() && !disabled {
			r.feedback, r.submitted = s.Draft, true
			slog.Debug("mujica: feedback_widget submitted", "kind", s.Draft.Kind)
		}
	})
	return r
}
