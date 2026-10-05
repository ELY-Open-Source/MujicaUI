package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

type devDemoSettings struct {
	disabled bool
	status   mujica.DataStatus
	empty    bool
}

func devControls(c *ui.Context, id string, log func(string), states bool) *devDemoSettings {
	s := ui.Local(c.Root(), id+"-controls", func() devDemoSettings { return devDemoSettings{} })
	dtControls(c, func() {
		if mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
			log("disabled changed")
		}
		switch id {
		case "451", "454", "456", "458":
		default:
			if mujica.Checkbox(c, &s.empty, "Empty data", mujica.CheckboxOptions{}).Changed() {
				log("data state changed")
			}
		}
		if states {
			options := []mujica.SelectOption[mujica.DataStatus]{{Value: mujica.DataReady, Label: "Ready"}, {Value: mujica.DataLoading, Label: "Loading"}, {Value: mujica.DataFailed, Label: "Failed"}}
			if mujica.Select(c, &s.status, options, mujica.SelectOptions{Label: "Status"}).Changed() {
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
		if mujica.Checkbox(c, lower, "Lower is better", mujica.CheckboxOptions{}).Changed() {
			log("change direction updated")
		}
		devPreview(c, s, func() {
			value := "12,480"
			trend := []float64{3, 5, 4, 8, 7, 9}
			if s.empty {
				value, trend = "—", nil
			}
			if mujica.StatCard(c, value, mujica.StatCardOptions{Label: "Requests", Change: 0.12, HasChange: true, LowerIsBetter: *lower, Trend: trend, Status: s.status, Error: "Source unreachable"}).Changed() {
				log("retry requested")
			}
		})
	})
}

func devAlertListDemo(c *ui.Context, log func(string)) {
	dtSection(c, "AlertList", func() {
		s := devControls(c, "443", log, false)
		alerts := ui.Local(c.Root(), "443", func() []mujica.AlertItem {
			return []mujica.AlertItem{{ID: "cpu", Title: "CPU usage high", Level: mujica.AlertWarning, Source: "worker", At: time.Now()}, {ID: "api", Title: "API unavailable", Level: mujica.AlertCritical, Source: "gateway", At: time.Now()}}
		})
		devPreview(c, s, func() {
			shown := *alerts
			if s.empty {
				shown = nil
			}
			r := mujica.AlertList(c, shown, mujica.AlertListOptions{})
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
		status := ui.Local(c.Root(), "444", func() mujica.IncidentStatus { return mujica.IncidentMonitoring })
		choices := []mujica.SelectOption[mujica.IncidentStatus]{{Value: mujica.IncidentInvestigating, Label: "Investigating"}, {Value: mujica.IncidentIdentified, Label: "Identified"}, {Value: mujica.IncidentMonitoring, Label: "Monitoring"}, {Value: mujica.IncidentResolved, Label: "Resolved"}}
		if mujica.Select(c, status, choices, mujica.SelectOptions{Label: "Incident status"}).Changed() {
			log("incident status changed")
		}
		devPreview(c, s, func() {
			updates := []mujica.IncidentUpdate{{Status: mujica.IncidentIdentified, Text: "Slow query located"}, {Status: mujica.IncidentMonitoring, Text: "Recovery observed"}}
			if s.empty {
				updates = nil
			}
			mujica.IncidentCard(c, mujica.Incident{Title: "API latency", Status: *status, Impact: mujica.ImpactMinor, Updates: updates})
		})
	})
}

func devUptimeBarDemo(c *ui.Context, log func(string)) {
	dtSection(c, "UptimeBar", func() {
		s := devControls(c, "445", log, true)
		days := ui.Local(c.Root(), "445", func() []mujica.UptimeDay {
			out := make([]mujica.UptimeDay, 30)
			start := mujica.Date{Year: 2026, Month: 10, Day: 1}
			for i := range out {
				out[i] = mujica.UptimeDay{Date: start.AddDays(i), Uptime: 1}
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
			r := mujica.UptimeBar(c, data, mujica.UptimeBarOptions{Frame: mujica.ChartFrame{Label: "API uptime", Status: s.status, Error: "Source unreachable"}})
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
		health := ui.Local(c.Root(), "446", func() mujica.ServiceHealth { return mujica.ServiceDegraded })
		choices := []mujica.SelectOption[mujica.ServiceHealth]{{Value: mujica.ServiceOperational, Label: "Operational"}, {Value: mujica.ServiceDegraded, Label: "Degraded"}, {Value: mujica.ServiceOutage, Label: "Outage"}}
		if mujica.Select(c, health, choices, mujica.SelectOptions{Label: "Search health"}).Changed() {
			log("service health changed")
		}
		devPreview(c, s, func() {
			data := []mujica.ServiceItem{{Name: "API", Health: mujica.ServiceOperational}, {Name: "Search", Health: *health, Note: "Search cluster"}}
			if s.empty {
				data = nil
			}
			mujica.ServiceStatus(c, data, mujica.ServiceStatusOptions{})
		})
	})
}

func devResourceGaugeDemo(c *ui.Context, log func(string)) {
	dtSection(c, "ResourceGauge", func() {
		s := devControls(c, "449", log, true)
		used := ui.Local(c.Root(), "449", func() float64 { return 45 })
		if mujica.NumberInput(c, used, mujica.NumberInputOptions{Min: 0, Max: 100, Step: 5, Label: "CPU usage"}).Changed() {
			log("resource reading changed")
		}
		devPreview(c, s, func() {
			data := []mujica.ResourceReading{{Name: "CPU", Used: *used, Total: 100}, {Name: "Memory", Used: 6.4, Total: 8, Unit: "GB"}, {Name: "Disk", Used: 470, Total: 500, Unit: "GB"}}
			if s.empty {
				data = nil
			}
			if mujica.ResourceGauge(c, data, mujica.ResourceGaugeOptions{Status: s.status, Error: "Source unreachable"}).Changed() {
				log("retry requested")
			}
		})
	})
}
