package mujica

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
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

// ServiceStatus lists services as operational, degraded or down, under an overall summary.
func ServiceStatus(c *ui.Context, services []ServiceItem, opts ServiceStatusOptions) *ui.Element {
	if opts.MaxHeight < 0 {
		panic("mujica: ServiceStatus MaxHeight is negative")
	}
	worst, counts := ServiceOperational, [3]int{}
	for i, s := range services {
		if s.Name == "" || s.Health < ServiceOperational || s.Health > ServiceOutage {
			panic(fmt.Sprintf("mujica: ServiceStatus service %d needs a Name and a known Health", i))
		}
		worst = max(worst, s.Health)
		counts[s.Health]++
	}
	maxH := opts.MaxHeight
	if maxH == 0 {
		maxH = 360
	}
	k := Tokens(c)
	u := Density(c).Unit()
	summary := msg(c, [...]string{"service_status.all_ok", "service_status.some_degraded", "service_status.some_down"}[worst])
	col := ui.Column(c).Shrink(0).Gap(u * 2).Label(summary)
	col.Children(func() {
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Padding(u*3).Radius(theme.CardRadius).Border(theme.BorderWidth, k.Border).
			Background(k.Surface).Role(ui.RoleStatus).Label(summary).Children(func() {
			devtoolsSeverity(c, serviceSeverity(worst), summary)
			ui.Box(c).Grow(1)
			agentCaption(c, fmt.Sprintf(msg(c, "service_status.counts"), counts[0], counts[1], counts[2]))
		})
		if len(services) == 0 {
			agentCaption(c, msg(c, "service_status.empty"))
			return
		}
		agentRows(c, len(services), func(i int) {
			s := services[i]
			word := serviceWord(c, s.Health)
			ui.Row(c).Key(i).Gap(u * 3).AlignItems(ui.Center).MinHeight(Density(c).ControlHeight()).Padding(u, u*3).
				Label(s.Name + ", " + word).Children(func() {
				ui.Column(c).Grow(1).Shrink(1).Children(func() {
					ui.Text(c, s.Name).SingleLine()
					if s.Note != "" {
						agentCaption(c, s.Note)
					}
				})
				devtoolsSeverity(c, serviceSeverity(s.Health), word)
			})
		}).Shrink(1).MaxHeight(maxH).Background(k.Surface).Border(theme.BorderWidth, k.Border).
			Radius(theme.CardRadius).Dividers(theme.BorderWidth, k.Border)
	})
	return col
}

func serviceSeverity(h ServiceHealth) Severity {
	return [...]Severity{SeveritySuccess, SeverityWarning, SeverityError}[h]
}

func serviceWord(c *ui.Context, h ServiceHealth) string {
	return msg(c, [...]string{"service_status.operational", "service_status.degraded", "service_status.outage"}[h])
}
