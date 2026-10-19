package devtools

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/agent"
	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ServiceHealth is a service's state. The zero value is ServiceOperational.
type ServiceHealth int

const (
	ServiceOperational ServiceHealth = iota
	ServiceDegraded
	ServiceOutage
)

// ServiceItem is one service the caller reports.
type ServiceItem struct {
	Name   string
	Health ServiceHealth
	// Note explains a degraded or down service.
	Note string
}

// ServiceStatusOptions configure a ServiceStatus.
type ServiceStatusOptions struct {
	// MaxHeight caps the scrolling list; zero means 360.
	MaxHeight float32
}

// ServiceStatus lists service health and an overall summary.
func ServiceStatus(c *ui.Context, services []ServiceItem, opts ServiceStatusOptions) ui.Element {
	if opts.MaxHeight < 0 {
		panic("mujica: ServiceStatus MaxHeight is negative")
	}
	col := ui.Column(c).Shrink(0).Gap(core.Density(c).Unit() * 2)
	summaryState := ui.Local(col, "mujica.services.summary", func() serviceSummary { return serviceSummary{key: "unset"} })
	if key := chart.ChartID(services); key != summaryState.key {
		worst, counts := ServiceOperational, [3]int{}
		for i, s := range services {
			if s.Name == "" || s.Health < ServiceOperational || s.Health > ServiceOutage {
				panic(fmt.Sprintf("mujica: ServiceStatus service %d needs a Name and a known Health", i))
			}
			worst = max(worst, s.Health)
			counts[s.Health]++
		}
		summaryState.key, summaryState.worst, summaryState.counts = key, worst, counts
	}
	worst, counts := summaryState.worst, summaryState.counts
	maxH := opts.MaxHeight
	if maxH == 0 {
		maxH = 360
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	summary := core.Msg(c, [...]string{"service_status.all_ok", "service_status.some_degraded", "service_status.some_down"}[worst])
	col.Label(summary)
	col.Children(func() {
		ui.Row(c).Gap(u*3).AlignItems(ui.Center).Padding(u*3).Radius(theme.CardRadius).Border(theme.BorderWidth, k.Border).
			Background(k.Surface).Role(ui.RoleStatus).Label(summary).Children(func() {
			devtoolsSeverity(c, serviceSeverity(worst), summary)
			ui.Box(c).Grow(1)
			agent.AgentCaption(c, fmt.Sprintf(core.Msg(c, "service_status.counts"), counts[0], counts[1], counts[2]))
		})
		if len(services) == 0 {
			agent.AgentCaption(c, core.Msg(c, "service_status.empty"))
			return
		}
		agent.AgentRows(c, len(services), func(i int) {
			s := services[i]
			word := serviceWord(c, s.Health)
			ui.Row(c.Key(i)).Gap(u*3).AlignItems(ui.Center).MinHeight(core.Density(c).ControlHeight()).Padding(u, u*3).
				Label(s.Name + ", " + word).Children(func() {
				ui.Column(c).Grow(1).Shrink(1).Children(func() {
					ui.Text(c, s.Name).SingleLine()
					if s.Note != "" {
						agent.AgentCaption(c, s.Note)
					}
				})
				devtoolsSeverity(c, serviceSeverity(s.Health), word)
			})
		}).Shrink(1).MaxHeight(maxH).Background(k.Surface).Border(theme.BorderWidth, k.Border).
			Radius(theme.CardRadius).Dividers(theme.BorderWidth, k.Border)
	})
	return col
}

func serviceSeverity(h ServiceHealth) core.Severity {
	return [...]core.Severity{core.SeveritySuccess, core.SeverityWarning, core.SeverityError}[h]
}

func serviceWord(c *ui.Context, h ServiceHealth) string {
	return core.Msg(c, [...]string{"service_status.operational", "service_status.degraded", "service_status.outage"}[h])
}

type serviceSummary struct {
	key    string
	worst  ServiceHealth
	counts [3]int
}
