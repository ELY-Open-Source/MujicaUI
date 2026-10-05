package mujica

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// PlanItem is a plan task: pending, running, done or skipped.
type PlanItem struct {
	ID    string
	Text  string
	State AgentState
}

// AgentPlanOptions configure an AgentPlan.
type AgentPlanOptions struct {
	// Title heads the checklist; empty uses "Plan".
	Title    string
	Disabled bool
	Still    bool
}

type agentPlanState struct{ draft string }

// AgentPlan is a live checklist with skip, restore and add.
func AgentPlan(c *ui.Context, items *[]PlanItem, opts AgentPlanOptions) ChoiceResult {
	if items == nil {
		panic("mujica: AgentPlan needs items")
	}
	seen := map[string]bool{}
	done, counted := 0, 0
	for i, it := range *items {
		if it.ID == "" || seen[it.ID] {
			panic(fmt.Sprintf("mujica: AgentPlan item %d needs a unique ID", i))
		}
		switch it.State {
		case AgentPending, AgentRunning, AgentDone, AgentSkipped:
		default:
			panic(fmt.Sprintf("mujica: AgentPlan item %q has state %v", it.ID, it.State))
		}
		seen[it.ID] = true
		if it.State != AgentSkipped {
			counted++
		}
		done += agentB2I(it.State == AgentDone)
	}
	u := Density(c).Unit()
	title := opts.Title
	if title == "" {
		title = msg(c, "agent_plan.title")
	}
	var r ChoiceResult
	r.Element = agentCard(c).Disabled(opts.Disabled).Label(title)
	st := ui.Local(r.Element, "mujica.plan", func() agentPlanState { return agentPlanState{} })
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, title).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).FontWeight(600).Grow(1)
			agentCaption(c, fmt.Sprintf(msg(c, "agent_plan.progress"), done, counted))
		})
		frac := 0.0
		if counted > 0 {
			frac = float64(done) / float64(counted)
		}
		agentBar(c, frac, fmt.Sprintf(msg(c, "agent_plan.progress"), done, counted))
		list := *items
		for i := range list {
			ui.Row(c).Key(list[i].ID).Gap(u * 2).AlignItems(ui.Center).MinHeight(Density(c).ControlHeight()).Children(func() {
				if agentPlanRow(c, &list[i], locked, opts.Still) {
					r.changed = true
				}
			})
		}
		if !locked && agentPlanAdd(c, &st.draft) {
			*items = append(*items, PlanItem{ID: agentPlanID(seen), Text: strings.TrimSpace(st.draft)})
			st.draft = ""
			r.changed = true
		}
	})
	if r.changed {
		slog.Debug("mujica: agent_plan change", "items", len(*items))
	}
	return r
}

// agentPlanRow builds one item, reporting a user change.
func agentPlanRow(c *ui.Context, it *PlanItem, locked, still bool) bool {
	k := Tokens(c)
	checked := it.State == AgentDone
	box := Checkbox(c, &checked, "", CheckboxOptions{Disabled: locked || it.State == AgentSkipped}).Label(it.Text)
	changed := !locked && box.Changed()
	if changed {
		it.State = AgentPending
		if checked {
			it.State = AgentDone
		}
	}
	t := ui.Text(c, it.Text).Grow(1).Shrink(1)
	switch it.State {
	case AgentSkipped:
		t.Strikethrough().TextColor(k.TextMuted)
	case AgentRunning:
		agentMark(c, AgentRunning, 14, still)
		agentCaption(c, msg(c, "agent_plan.doing"))
	}
	if it.State == AgentDone || locked {
		return changed
	}
	label, icon, next := msg(c, "agent_plan.skip"), "ban", AgentSkipped
	if it.State == AgentSkipped {
		label, icon, next = msg(c, "agent_plan.restore"), "refresh-cw", AgentPending
	}
	if agentToggle(c, label, icon, false) {
		it.State = next
		changed = true
	}
	return changed
}

// agentPlanAdd builds the new-item field and reports a non-empty submit.
func agentPlanAdd(c *ui.Context, draft *string) bool {
	u := Density(c).Unit()
	add := false
	ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
		in := TextInput(c, draft, TextInputOptions{Placeholder: msg(c, "agent_plan.placeholder"), Label: msg(c, "agent_plan.placeholder")})
		in.Element.Grow(1)
		blank := strings.TrimSpace(*draft) == ""
		add = Button(c, msg(c, "agent_plan.add"), ButtonOptions{Variant: Secondary, Disabled: blank}).Clicked()
		add = (add || in.Submitted()) && !blank
	})
	return add
}

// agentPlanID returns an ID no item uses yet.
func agentPlanID(seen map[string]bool) string {
	for n := len(seen) + 1; ; n++ {
		if id := fmt.Sprintf("item-%d", n); !seen[id] {
			return id
		}
	}
}
