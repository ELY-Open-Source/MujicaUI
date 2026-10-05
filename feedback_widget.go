package mujica

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"
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
	Form  FormState
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
	Element   *ui.Element
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
	d := Density(c)
	var r FeedbackWidgetResult
	r.Element = ui.Column(c).Gap(d.Unit() * 4).Width(420).Disabled(opts.Disabled)
	disabled := r.Element.IsDisabled()
	if disabled {
		frozen := *s
		frozen.Form.order = append([]string(nil), s.Form.order...)
		frozen.Form.fields = make(map[string]*formFieldState, len(s.Form.fields))
		for name, field := range s.Form.fields {
			copy := *field
			frozen.Form.fields[name] = &copy
		}
		frozen.Form.requested, frozen.Form.awaiting = false, false
		s = &frozen
	}
	locked := disabled || s.Form.Submitting()
	r.Element.Children(func() {
		kind := int(s.Draft.Kind)
		var segs []Segment
		for f := FeedbackBug; f <= FeedbackOther; f++ {
			segs = append(segs, Segment{Label: msg(c, fmt.Sprintf("feedback_widget.kind%d", f))})
		}
		SegmentedControl(c, &kind, segs, SegmentedControlOptions{Equal: true, Label: msg(c, "feedback_widget.kind"), Disabled: locked})
		if !locked {
			s.Draft.Kind = FeedbackKind(kind)
		}
		res := Form(c, &s.Form, FormOptions{SubmitLabel: msg(c, "feedback_widget.send"), ResetLabel: msg(c, "feedback_widget.clear")}, func(f *FormScope) {
			FormItem(c, f, "text", &s.Draft.Text, FormItemOptions[string]{Label: msg(c, "feedback_widget.message"), Required: true,
				Validate: func(v string) string {
					if n := len([]rune(strings.TrimSpace(v))); n < minLen {
						return fmt.Sprintf(msg(c, "feedback_widget.short"), minLen)
					}
					return ""
				}}, func() *ui.Element {
				return TextArea(c, &s.Draft.Text, TextAreaOptions{MinHeight: 96, MaxLength: 2000, Disabled: locked}).Input
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
