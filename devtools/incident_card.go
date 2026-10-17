package devtools

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/agent"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// IncidentStatus tracks investigation; zero means IncidentInvestigating.
type IncidentStatus int

const (
	IncidentInvestigating IncidentStatus = iota
	IncidentIdentified
	IncidentMonitoring
	IncidentResolved
)

// IncidentImpact tracks severity; zero means ImpactNone.
type IncidentImpact int

const (
	ImpactNone IncidentImpact = iota
	ImpactMinor
	ImpactMajor
	ImpactCritical
)

// IncidentUpdate is one timeline entry, oldest first.
type IncidentUpdate struct {
	At     time.Time
	Status IncidentStatus
	Text   string
}

// Incident is what the caller reports about an incident.
type Incident struct {
	Title   string
	Status  IncidentStatus
	Impact  IncidentImpact
	Started time.Time
	Updates []IncidentUpdate
}

// IncidentCard shows an incident's title, status, impact and update timeline.
func IncidentCard(c *ui.Context, inc Incident) ui.Element {
	if inc.Title == "" {
		panic("mujica: IncidentCard needs a Title")
	}
	status := incidentStatusWord(c, inc.Status)
	if inc.Impact < ImpactNone || inc.Impact > ImpactCritical {
		panic(fmt.Sprintf("mujica: unknown IncidentImpact %d", int(inc.Impact)))
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	sev := [...]core.Severity{core.SeverityInfo, core.SeverityWarning, core.SeverityError, core.SeverityError}[inc.Impact]
	impact := core.Msg(c, [...]string{"incident_card.none", "incident_card.minor", "incident_card.major", "incident_card.critical"}[inc.Impact])
	card := agent.AgentCard(c).Label(inc.Title + ", " + status)
	card.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			devtoolsTitle(c, inc.Title).Grow(1).Shrink(1)
			if inc.Status == IncidentResolved {
				devtoolsSeverity(c, core.SeveritySuccess, status)
			} else {
				devtoolsSeverity(c, core.SeverityWarning, status)
			}
		})
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Children(func() {
			devtoolsSeverity(c, sev, fmt.Sprintf(core.Msg(c, "incident_card.impact"), impact))
			if !inc.Started.IsZero() {
				agent.AgentCaption(c, fmt.Sprintf(core.Msg(c, "incident_card.started"), inc.Started.Format("2006-01-02 15:04")))
			}
		})
		if len(inc.Updates) > 0 {
			agent.AgentRows(c, len(inc.Updates), func(i int) { incidentUpdate(c, k, inc.Updates[i], i == len(inc.Updates)-1) }).Shrink(1).MaxHeight(320)
		}
	})
	return card
}

// incidentUpdate draws one timeline entry on a diamond rail.
func incidentUpdate(c *ui.Context, k theme.Tokens, up IncidentUpdate, last bool) {
	u := core.Density(c).Unit()
	ui.Row(c).Gap(u * 3).AlignItems(ui.Stretch).Children(func() {
		ui.Box(c).Width(12).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			cx := internal.Snap(p, r.X+r.W/2)
			cy := internal.Snap(p, r.Y+9)
			if !last {
				p.Line(cx, cy+5, cx, r.Y+r.H, 1, k.ControlBorder)
			}
			p.FillPath(internal.Diamond(cx, cy, 9), k.Ornament)
		})
		ui.Column(c).Grow(1).Shrink(1).Gap(u/2).Padding(0, 0, u*2, 0).Children(func() {
			ui.Row(c).Gap(u * 2).Children(func() {
				ui.Text(c, incidentStatusWord(c, up.Status)).FontWeight(600)
				if !up.At.IsZero() {
					agent.AgentCaption(c, up.At.Format("15:04"))
				}
			})
			ui.Text(c, up.Text).TextColor(k.TextMuted)
		})
	})
}

func incidentStatusWord(c *ui.Context, s IncidentStatus) string {
	if s < IncidentInvestigating || s > IncidentResolved {
		panic(fmt.Sprintf("mujica: unknown IncidentStatus %d", int(s)))
	}
	return core.Msg(c, [...]string{"incident_card.investigating", "incident_card.identified", "incident_card.monitoring", "incident_card.resolved"}[s])
}
