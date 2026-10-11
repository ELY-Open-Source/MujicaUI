package agent

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// AgentStep is one step of an agent run.
type AgentStep struct {
	ID       string
	Title    string
	Detail   string
	State    AgentState
	Duration time.Duration
}

// AgentStepListOptions configure an AgentStepList.
type AgentStepListOptions struct {
	Disabled bool
	Still    bool
}

// AgentStepListResult reports this frame's events.
type AgentStepListResult struct {
	Element *ui.Element
	// Retry is the failed step to retry, else "".
	Retry string
	// Toggled is the step whose detail toggled, else "".
	Toggled string
}

// AgentStepList is a step timeline with details and retries.
func AgentStepList(c *ui.Context, steps []AgentStep, opts AgentStepListOptions) AgentStepListResult {
	current := -1
	seen := make(map[string]bool, len(steps))
	for i, s := range steps {
		if s.ID == "" || s.Title == "" || seen[s.ID] {
			panic(fmt.Sprintf("mujica: AgentStepList step %d needs a unique ID and a Title", i))
		}
		seen[s.ID] = true
		if current < 0 && (s.State == AgentRunning || s.State == AgentWaiting || s.State == AgentPlanning) {
			current = i
		}
	}
	var r AgentStepListResult
	r.Element = ui.Column(c).Shrink(0).FocusGroup(ui.Vertical).Disabled(opts.Disabled)
	open := ui.Local(r.Element, "mujica.steps.open", func() map[string]bool { return map[string]bool{} })
	r.Element.Children(func() {
		for i, s := range steps {
			ui.Column(c).Key(s.ID).Children(func() {
				if agentStepRow(c, s, i, len(steps), i == current, (*open)[s.ID], opts.Still, &r) {
					(*open)[s.ID] = !(*open)[s.ID]
					r.Toggled = s.ID
				}
			})
		}
	})
	if r.Element.IsDisabled() {
		r.Retry, r.Toggled = "", ""
	}
	for id := range *open {
		if !seen[id] {
			delete(*open, id)
		}
	}
	if r.Retry != "" || r.Toggled != "" {
		slog.Debug("mujica: agent_step_list event", "retry", r.Retry, "toggled", r.Toggled)
	}
	return r
}

// agentStepRow builds one step, reporting a header click.
func agentStepRow(c *ui.Context, s AgentStep, i, n int, current, open, still bool, r *AgentStepListResult) bool {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	clicked := false
	top := (d.ControlHeight() - 18) / 2 // centers the mark on the step header
	ui.Row(c).Gap(u * 3).AlignItems(ui.Stretch).Children(func() {
		// The rail joins the marks; it is drawn, not focusable.
		ui.Column(c).Width(20).Shrink(0).AlignItems(ui.Center).Draw(func(p *ui.Painter, rc ui.Rect) {
			x := internal.Snap(p, rc.X+rc.W/2)
			if i > 0 {
				p.Line(x, rc.Y, x, rc.Y+top-3, 1, k.ControlBorder)
			}
			if i < n-1 {
				p.Line(x, rc.Y+top+21, x, rc.Y+rc.H, 1, k.ControlBorder)
			}
		}).Children(func() {
			ui.Box(c).Height(top)
			agentMark(c, s.State, 18, still)
		})
		ui.Column(c).Grow(1).Shrink(1).Gap(u).Padding(0, 0, u*2, 0).Children(func() {
			head := ui.ButtonBase(c).Justify(ui.Start).Gap(u*2).MinHeight(d.ControlHeight()).Padding(0, u*2).
				Radius(theme.ControlRadius).Label(s.Title + ", " + agentWord(c, s.State))
			internal.Ring(head, k.Focus, theme.ControlRadius)
			if current {
				head.Background(k.Selection).Border(theme.BorderWidth, k.AccentText)
			} else if head.Hovered() {
				head.Background(k.SurfaceHover)
			}
			head.Children(func() {
				ui.Text(c, s.Title).FontWeight(600).Grow(1).Shrink(1)
				if current {
					AgentCaption(c, core.Msg(c, "agent_step_list.current"))
				}
				AgentCaption(c, agentWord(c, s.State))
				if s.State == AgentDone || s.State == AgentFailed {
					AgentCaption(c, AgentDuration(s.Duration))
				}
				if s.Detail != "" {
					chev := "chevron-right"
					if open {
						chev = "chevron-down"
					}
					ui.Icon(c, icons.Must(chev)).Size(14, 14).TextColor(k.TextMuted)
				}
			})
			clicked = s.Detail != "" && !head.IsDisabled() && head.Clicked()
			if open && s.Detail != "" {
				ui.Text(c, s.Detail).TextColor(k.TextMuted).Selectable().Padding(0, u*2)
			}
			if s.State == AgentFailed {
				ui.Row(c).Children(func() {
					if AgentToggle(c, core.Msg(c, "agent_step_list.retry"), "refresh-cw", head.IsDisabled()) {
						r.Retry = s.ID
					}
				})
			}
		})
	})
	return clicked
}
