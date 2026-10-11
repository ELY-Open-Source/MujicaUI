package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/devtools"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type devDemoSettings struct {
	disabled bool
	status   data.DataStatus
	empty    bool
}

func devControls(c *ui.Context, id string, log func(string), states bool) *devDemoSettings {
	s := ui.Local(c.Root(), id+"-controls", func() devDemoSettings { return devDemoSettings{} })
	dtControls(c, func() {
		if input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{}).Changed() {
			log("disabled changed")
		}
		switch id {
		case "451", "454", "456", "458":
		default:
			if input.Checkbox(c, &s.empty, "Empty data", input.CheckboxOptions{}).Changed() {
				log("data state changed")
			}
		}
		if states {
			options := []input.SelectOption[data.DataStatus]{{Value: data.DataReady, Label: "Ready"}, {Value: data.DataLoading, Label: "Loading"}, {Value: data.DataFailed, Label: "Failed"}}
			if input.Select(c, &s.status, options, input.SelectOptions{Label: "Status"}).Changed() {
				log("status changed")
			}
		}
	})
	return s
}

func devPreview(c *ui.Context, s *devDemoSettings, build func()) {
	ui.Column(c).Disabled(s.disabled).Gap(8).Children(build)
}

func devStatCardDemo(c *ui.Context, log func(string)) {
	dtSection(c, "StatCard", func() {
		s := devControls(c, "442", log, true)
		lower := ui.Local(c.Root(), "442-lower", func() bool { return false })
		if input.Checkbox(c, lower, "Lower is better", input.CheckboxOptions{}).Changed() {
			log("change direction updated")
		}
		devPreview(c, s, func() {
			value := "12,480"
			trend := []float64{3, 5, 4, 8, 7, 9}
			if s.empty {
				value, trend = "—", nil
			}
			if devtools.StatCard(c, value, devtools.StatCardOptions{Label: "Requests", Change: 0.12, HasChange: true, LowerIsBetter: *lower, Trend: trend, Status: s.status, Error: "Source unreachable"}).Changed() {
				log("retry requested")
			}
		})
	})
}

func devAlertListDemo(c *ui.Context, log func(string)) {
	dtSection(c, "AlertList", func() {
		s := devControls(c, "443", log, false)
		alerts := ui.Local(c.Root(), "443", func() []devtools.AlertItem {
			return []devtools.AlertItem{{ID: "cpu", Title: "CPU usage high", Level: devtools.AlertWarning, Source: "worker", At: time.Now()}, {ID: "api", Title: "API unavailable", Level: devtools.AlertCritical, Source: "gateway", At: time.Now()}}
		})
		devPreview(c, s, func() {
			shown := *alerts
			if s.empty {
				shown = nil
			}
			r := devtools.AlertList(c, shown, devtools.AlertListOptions{})
			for i := range *alerts {
				if r.Ack == (*alerts)[i].ID {
					(*alerts)[i].Acked = true
					log("acknowledged " + r.Ack)
				}
				if r.Silence == (*alerts)[i].ID {
					(*alerts)[i].Silence = c.Now().Add(time.Hour)
					log("silenced " + r.Silence)
				}
			}
		})
	})
}

func devIncidentCardDemo(c *ui.Context, log func(string)) {
	dtSection(c, "IncidentCard", func() {
		s := devControls(c, "444", log, false)
		status := ui.Local(c.Root(), "444", func() devtools.IncidentStatus { return devtools.IncidentMonitoring })
		choices := []input.SelectOption[devtools.IncidentStatus]{{Value: devtools.IncidentInvestigating, Label: "Investigating"}, {Value: devtools.IncidentIdentified, Label: "Identified"}, {Value: devtools.IncidentMonitoring, Label: "Monitoring"}, {Value: devtools.IncidentResolved, Label: "Resolved"}}
		if input.Select(c, status, choices, input.SelectOptions{Label: "Incident status"}).Changed() {
			log("incident status changed")
		}
		devPreview(c, s, func() {
			updates := []devtools.IncidentUpdate{{Status: devtools.IncidentIdentified, Text: "Slow query located"}, {Status: devtools.IncidentMonitoring, Text: "Recovery observed"}}
			if s.empty {
				updates = nil
			}
			devtools.IncidentCard(c, devtools.Incident{Title: "API latency", Status: *status, Impact: devtools.ImpactMinor, Updates: updates})
		})
	})
}

func devUptimeBarDemo(c *ui.Context, log func(string)) {
	dtSection(c, "UptimeBar", func() {
		s := devControls(c, "445", log, true)
		days := ui.Local(c.Root(), "445", func() []devtools.UptimeDay {
			out := make([]devtools.UptimeDay, 30)
			start := core.Date{Year: 2026, Month: 10, Day: 1}
			for i := range out {
				out[i] = devtools.UptimeDay{Date: start.AddDays(i), Uptime: 1}
			}
			out[7].Uptime, out[7].Note = 0.91, "2 h outage"
			out[12].Uptime = 0.98
			return out
		})
		devPreview(c, s, func() {
			data := *days
			if s.empty {
				data = nil
			}
			r := devtools.UptimeBar(c, data, devtools.UptimeBarOptions{Frame: chart.ChartFrame{Label: "API uptime", Status: s.status, Error: "Source unreachable"}})
			if r.Changed() {
				log(fmt.Sprintf("day %d", r.Point))
			}
			if r.Retried() {
				log("retry requested")
			}
		})
	})
}

func devServiceStatusDemo(c *ui.Context, log func(string)) {
	dtSection(c, "ServiceStatus", func() {
		s := devControls(c, "446", log, false)
		health := ui.Local(c.Root(), "446", func() devtools.ServiceHealth { return devtools.ServiceDegraded })
		choices := []input.SelectOption[devtools.ServiceHealth]{{Value: devtools.ServiceOperational, Label: "Operational"}, {Value: devtools.ServiceDegraded, Label: "Degraded"}, {Value: devtools.ServiceOutage, Label: "Outage"}}
		if input.Select(c, health, choices, input.SelectOptions{Label: "Search health"}).Changed() {
			log("service health changed")
		}
		devPreview(c, s, func() {
			data := []devtools.ServiceItem{{Name: "API", Health: devtools.ServiceOperational}, {Name: "Search", Health: *health, Note: "Search cluster"}}
			if s.empty {
				data = nil
			}
			devtools.ServiceStatus(c, data, devtools.ServiceStatusOptions{})
		})
	})
}

func devResourceGaugeDemo(c *ui.Context, log func(string)) {
	dtSection(c, "ResourceGauge", func() {
		s := devControls(c, "449", log, true)
		used := ui.Local(c.Root(), "449", func() float64 { return 45 })
		if input.NumberInput(c, used, input.NumberInputOptions{Min: 0, Max: 100, Step: 5, Label: "CPU usage"}).Changed() {
			log("resource reading changed")
		}
		devPreview(c, s, func() {
			data := []devtools.ResourceReading{{Name: "CPU", Used: *used, Total: 100}, {Name: "Memory", Used: 6.4, Total: 8, Unit: "GB"}, {Name: "Disk", Used: 470, Total: 500, Unit: "GB"}}
			if s.empty {
				data = nil
			}
			if devtools.ResourceGauge(c, data, devtools.ResourceGaugeOptions{Status: s.status, Error: "Source unreachable"}).Changed() {
				log("retry requested")
			}
		})
	})
}
