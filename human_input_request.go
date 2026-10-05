package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// HumanInputRequestOptions configure a HumanInputRequest. Question is required.
type HumanInputRequestOptions struct {
	Question string
	// Options are suggested answers; AllowFree adds a free one.
	Options   []string
	AllowFree bool
	// Deadline is when the agent stops waiting; zero waits forever.
	Deadline time.Time
	Disabled bool
}

// HumanInputResult reports this frame's events.
type HumanInputResult struct {
	Element   *ui.Element
	Submitted bool
	// TimedOut reports the frame the deadline passed.
	TimedOut bool
}

type humanInputState struct {
	draft   string
	expired bool
}

// HumanInputRequest asks the user a question before a deadline.
func HumanInputRequest(c *ui.Context, answer *string, opts HumanInputRequestOptions) HumanInputResult {
	if answer == nil || opts.Question == "" {
		panic("mujica: HumanInputRequest needs an answer and a Question")
	}
	if len(opts.Options) == 0 && !opts.AllowFree {
		panic("mujica: HumanInputRequest needs Options or AllowFree")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r HumanInputResult
	r.Element = agentCard(c).Border(theme.BorderWidth, k.Info).Role(ui.RoleGroup).Label(opts.Question)
	st := ui.Local(r.Element, "mujica.humaninput", func() humanInputState { return humanInputState{} })
	left := time.Duration(0)
	if !opts.Deadline.IsZero() {
		left = opts.Deadline.Sub(c.Now())
		if left <= 0 && !st.expired {
			st.expired, r.TimedOut = true, true
		}
	}
	r.Element.Disabled(opts.Disabled || st.expired)
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Children(func() {
			ui.Icon(c, icons.Must("user")).Size(18, 18).TextColor(k.Info).Shrink(0).Margin(2, 0, 0, 0)
			ui.Text(c, opts.Question).FontWeight(600).Grow(1).Shrink(1)
		})
		if len(opts.Options) > 0 {
			items := make([]RadioItem[string], len(opts.Options))
			for i, o := range opts.Options {
				items[i] = RadioItem[string]{Value: o, Label: o}
			}
			chosen := *answer
			if RadioGroup(c, &chosen, items, RadioGroupOptions{Label: opts.Question, Disabled: locked}).Changed() && !locked {
				*answer, st.draft = chosen, ""
			}
		}
		submit := false
		if opts.AllowFree {
			in := TextInput(c, &st.draft, TextInputOptions{Placeholder: msg(c, "human_input_request.free"),
				Label: msg(c, "human_input_request.free"), Disabled: locked})
			if in.Changed() && !locked {
				*answer = st.draft
			}
			submit = in.Submitted()
		}
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			humanInputClock(c, opts.Deadline, left, st.expired)
			blank := strings.TrimSpace(*answer) == ""
			if Button(c, msg(c, "human_input_request.send"), ButtonOptions{Disabled: locked || blank}).Clicked() {
				submit = true
			}
			r.Submitted = submit && !locked && !blank
		})
	})
	if len(opts.Options) > 0 && !opts.AllowFree && *answer != "" && !slices.Contains(opts.Options, *answer) {
		panic(fmt.Sprintf("mujica: HumanInputRequest answer %q is not an option", *answer))
	}
	if r.Submitted || r.TimedOut {
		slog.Debug("mujica: human_input_request event", "submitted", r.Submitted, "timed_out", r.TimedOut)
	}
	return r
}

// humanInputClock shows the time left, or that it ran out.
func humanInputClock(c *ui.Context, deadline time.Time, left time.Duration, expired bool) {
	k := Tokens(c)
	text, tone := "", k.TextMuted
	switch {
	case expired:
		text, tone = msg(c, "human_input_request.timeout"), k.Warning
	case !deadline.IsZero():
		text = fmt.Sprintf(msg(c, "human_input_request.left"), agentDuration(left.Truncate(time.Second)))
		wait := left % time.Second
		if wait == 0 {
			wait = time.Second
		}
		c.After(wait)
	}
	ui.Row(c).Grow(1).Gap(Density(c).Unit()).AlignItems(ui.Center).Role(ui.RoleStatus).Label(text).Children(func() {
		if text == "" {
			return
		}
		ui.Icon(c, icons.Must("clock")).Size(14, 14).TextColor(tone)
		ui.Text(c, text).TextColor(tone).FontSize(FontSize(c, theme.CaptionSize+1))
	})
}
