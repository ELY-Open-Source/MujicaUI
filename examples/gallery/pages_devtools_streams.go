package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func devLogStreamDemo(c *ui.Context, log func(string)) {
	dtSection(c, "LogStream", func() {
		s := devControls(c, "447", log, false)
		type state struct {
			lines  []mujica.LogLine
			stream mujica.LogStreamState
		}
		st := ui.Local(c.Root(), "447", func() state {
			return state{lines: []mujica.LogLine{{Level: mujica.LogInfo, Text: "server ready"}, {Level: mujica.LogWarn, Text: "slow query"}}}
		})
		if mujica.Button(c, "Append line", mujica.ButtonOptions{}).Clicked() {
			st.lines = append(st.lines, mujica.LogLine{Time: c.Now(), Level: mujica.LogInfo, Text: fmt.Sprintf("request %d completed", len(st.lines))})
			log("line appended")
		}
		devPreview(c, s, func() {
			data := st.lines
			if s.empty {
				data = nil
			}
			if mujica.LogStream(c, data, &st.stream, mujica.LogStreamOptions{}).Changed() {
				log("pause or level changed")
			}
		})
	})
}

func devTraceWaterfallDemo(c *ui.Context, log func(string)) {
	dtSection(c, "TraceWaterfall", func() {
		s := devControls(c, "448", log, false)
		spans := ui.Local(c.Root(), "448", func() []mujica.WaterfallSpan {
			return []mujica.WaterfallSpan{{ID: "root", Name: "GET /api", Service: "gateway", Duration: 320 * time.Millisecond}, {ID: "db", Parent: "root", Name: "Query", Service: "database", Start: 40 * time.Millisecond, Duration: 90 * time.Millisecond}, {ID: "pay", Parent: "root", Name: "Charge", Service: "payments", Start: 230 * time.Millisecond, Duration: 80 * time.Millisecond, Error: true}}
		})
		devPreview(c, s, func() {
			data := *spans
			if s.empty {
				data = nil
			}
			if mujica.TraceWaterfall(c, data, mujica.TraceWaterfallOptions{}).Changed() {
				log("span selected")
			}
		})
	})
}

func devSystemMonitorDemo(c *ui.Context, log func(string)) {
	dtSection(c, "SystemMonitor", func() {
		s := devControls(c, "450", log, true)
		type state struct {
			times   []time.Time
			metrics []mujica.MonitorMetric
		}
		st := ui.Local(c.Root(), "450", func() state {
			out := state{metrics: []mujica.MonitorMetric{{Name: "CPU", Unit: "%"}, {Name: "Memory", Unit: "GB"}}}
			for i := range 20 {
				out.times = append(out.times, time.Date(2026, 10, 5, 9, 0, i*5, 0, time.UTC))
				out.metrics[0].Samples = append(out.metrics[0].Samples, 40+float64(i%7)*3)
				out.metrics[1].Samples = append(out.metrics[1].Samples, 3.2+float64(i)*0.02)
			}
			return out
		})
		if mujica.Button(c, "Append sample", mujica.ButtonOptions{}).Clicked() {
			st.times = append(append([]time.Time(nil), st.times...), st.times[len(st.times)-1].Add(5*time.Second))
			for i := range st.metrics {
				last := st.metrics[i].Samples[len(st.metrics[i].Samples)-1]
				st.metrics[i].Samples = append(append([]float64(nil), st.metrics[i].Samples...), last+0.1)
			}
			log("sample appended")
		}
		devPreview(c, s, func() {
			times, metrics := st.times, st.metrics
			if s.empty {
				times, metrics = nil, nil
			}
			if mujica.SystemMonitor(c, times, metrics, mujica.SystemMonitorOptions{Status: s.status, Error: "Agent offline"}).Changed() {
				log("retry requested")
			}
		})
	})
}

func devRefreshIntervalSelectorDemo(c *ui.Context, log func(string)) {
	dtSection(c, "RefreshIntervalSelector", func() {
		s := devControls(c, "451", log, false)
		interval := ui.Local(c.Root(), "451", func() time.Duration { return 0 })
		refreshing := ui.Local(c.Root(), "451-refreshing", func() bool { return false })
		if mujica.Checkbox(c, refreshing, "Refreshing", mujica.CheckboxOptions{}).Changed() {
			log("refresh state changed")
		}
		devPreview(c, s, func() {
			r := mujica.RefreshIntervalSelector(c, interval, mujica.RefreshIntervalSelectorOptions{Refreshing: *refreshing})
			if r.Refresh {
				log("refresh requested")
			}
			if r.Changed {
				log("interval changed")
			}
		})
	})
}

func devDashboardFilterBarDemo(c *ui.Context, log func(string)) {
	dtSection(c, "DashboardFilterBar", func() {
		s := devControls(c, "452", log, false)
		filters := ui.Local(c.Root(), "452", func() mujica.DashboardFilters { return mujica.DashboardFilters{} })
		devPreview(c, s, func() {
			opts := mujica.DashboardFilterBarOptions{Envs: []string{"production", "staging"}, Tags: []string{"api", "db", "web"}}
			if s.empty {
				opts.Envs, opts.Tags = nil, nil
				*filters = mujica.DashboardFilters{}
			}
			if mujica.DashboardFilterBar(c, filters, opts).Changed() {
				log("filters changed")
			}
		})
	})
}
