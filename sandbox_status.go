package mujica

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// SandboxState is a sandbox's lifecycle. The zero value is SandboxStarting.
type SandboxState int

const (
	SandboxStarting SandboxState = iota
	SandboxRunning
	SandboxStopped
	SandboxFailed
)

// SandboxResource is one metered resource, as CPU or memory.
type SandboxResource struct {
	Name        string
	Used, Limit float64
	// Unit follows the numbers, as "GB" or "cores".
	Unit string
}

// SandboxInfo is the sandbox the caller reports.
type SandboxInfo struct {
	Name      string
	State     SandboxState
	Resources []SandboxResource
	// Expires is the shutdown time; zero means none.
	Expires time.Time
}

// SandboxStatusOptions configure a SandboxStatus.
type SandboxStatusOptions struct {
	Still bool
}

// SandboxStatus shows sandbox state, resources and time left.
func SandboxStatus(c *ui.Context, info SandboxInfo, opts SandboxStatusOptions) *ui.Element {
	agent := map[SandboxState]AgentState{SandboxStarting: AgentPlanning, SandboxRunning: AgentRunning, SandboxStopped: AgentSkipped, SandboxFailed: AgentFailed}
	state, ok := agent[info.State]
	if !ok {
		panic(fmt.Sprintf("mujica: unknown SandboxState %d", int(info.State)))
	}
	for _, r := range info.Resources {
		if r.Name == "" || !(r.Limit > 0) || r.Used < 0 || math.IsNaN(r.Used) || r.Used > r.Limit {
			panic(fmt.Sprintf("mujica: SandboxStatus resource %q needs a name and 0 <= Used <= Limit", r.Name))
		}
	}
	k := Tokens(c)
	u := Density(c).Unit()
	word := msg(c, "sandbox_status."+[...]string{"starting", "running", "stopped", "failed"}[info.State])
	card := agentCard(c).Label(info.Name + ", " + word)
	card.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, info.Name).FontWeight(600).SingleLine().Grow(1).Shrink(1)
			AgentStatus(c, state, AgentStatusOptions{Label: word, Still: opts.Still})
		})
		for _, r := range info.Resources {
			frac := r.Used / r.Limit
			ui.Column(c).Gap(u).Children(func() {
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					agentCaption(c, r.Name).Grow(1)
					if frac >= 0.9 {
						ui.Icon(c, icons.Must("triangle-alert")).Size(12, 12).TextColor(k.Warning).Label(msg(c, "sandbox_status.high"))
					}
					ui.Text(c, fmt.Sprintf("%s / %s %s", sandboxNum(r.Used), sandboxNum(r.Limit), r.Unit)).
						FontSize(FontSize(c, theme.CaptionSize)).FontFeatures("tnum").TextColor(k.Text)
				})
				agentBar(c, frac, fmt.Sprintf("%s %.0f%%", r.Name, frac*100))
			})
		}
		if !info.Expires.IsZero() && info.State != SandboxStopped && info.State != SandboxFailed {
			left := info.Expires.Sub(c.Now())
			text := msg(c, "sandbox_status.expired")
			if left > 0 {
				text = fmt.Sprintf(msg(c, "sandbox_status.left"), agentDuration(left.Truncate(time.Second)))
				wait := left % time.Second
				if wait == 0 {
					wait = time.Second
				}
				c.After(wait)
			}
			ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).Role(ui.RoleStatus).Label(text).Children(func() {
				ui.Icon(c, icons.Must("clock")).Size(14, 14).TextColor(k.TextMuted)
				agentCaption(c, text)
			})
		}
	})
	return card
}

// sandboxNum prints a resource amount with at most one decimal.
func sandboxNum(v float64) string {
	if v == math.Trunc(v) {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.1f", v)
}
