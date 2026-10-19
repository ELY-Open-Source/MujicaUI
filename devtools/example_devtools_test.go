package devtools_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/code"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/devtools"
)

func ExampleStatCard() {
	value := "12,480"
	opts := devtools.StatCardOptions{Label: "Requests", Change: 0.12, HasChange: true, Trend: []float64{3, 5, 4, 8}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if devtools.StatCard(c, value, opts).Changed() {
			fmt.Println("retry")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleAlertList() {
	alerts := []devtools.AlertItem{{ID: "cpu", Title: "CPU usage high", Level: devtools.AlertWarning, Source: "worker", At: time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := devtools.AlertList(c, alerts, devtools.AlertListOptions{})
		if r.Ack != "" {
			fmt.Println("ack", r.Ack)
		}
		if r.Silence != "" {
			fmt.Println("silence", r.Silence)
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleIncidentCard() {
	incident := devtools.Incident{Title: "API latency", Status: devtools.IncidentMonitoring, Impact: devtools.ImpactMinor, Updates: []devtools.IncidentUpdate{{Status: devtools.IncidentIdentified, Text: "Slow query located"}, {Status: devtools.IncidentMonitoring, Text: "Recovery observed"}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		devtools.IncidentCard(c, incident)
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleUptimeBar() {
	days := []devtools.UptimeDay{{Date: core.Date{Year: 2026, Month: 10, Day: 1}, Uptime: 1}, {Date: core.Date{Year: 2026, Month: 10, Day: 2}, Uptime: 0.91, Note: "Recovery"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := devtools.UptimeBar(c, days, devtools.UptimeBarOptions{Frame: chart.ChartFrame{Label: "API uptime"}})
		if r.Changed() {
			fmt.Println("day", r.Point)
		}
		if r.Retried() {
			fmt.Println("retry")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleServiceStatus() {
	services := []devtools.ServiceItem{{Name: "API", Health: devtools.ServiceOperational}, {Name: "Search", Health: devtools.ServiceDegraded, Note: "Slow queries"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		devtools.ServiceStatus(c, services, devtools.ServiceStatusOptions{})
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleLogStream() {
	lines := []code.LogLine{{Level: code.LogInfo, Text: "server ready"}, {Level: code.LogWarn, Text: "slow query"}}
	var state devtools.LogStreamState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if devtools.LogStream(c, lines, &state, devtools.LogStreamOptions{}).Changed() {
			fmt.Println("filter changed")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleTraceWaterfall() {
	spans := []devtools.WaterfallSpan{{ID: "root", Name: "GET /api", Service: "gateway", Duration: 200 * time.Millisecond}, {ID: "db", Parent: "root", Name: "Query", Service: "database", Start: 20 * time.Millisecond, Duration: 80 * time.Millisecond}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if devtools.TraceWaterfall(c, spans, devtools.TraceWaterfallOptions{}).Changed() {
			fmt.Println("span selected")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleResourceGauge() {
	readings := []devtools.ResourceReading{{Name: "CPU", Used: 45, Total: 100}, {Name: "Memory", Used: 6.4, Total: 8, Unit: "GB"}, {Name: "Disk", Used: 470, Total: 500, Unit: "GB"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if devtools.ResourceGauge(c, readings, devtools.ResourceGaugeOptions{}).Changed() {
			fmt.Println("retry")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleSystemMonitor() {
	times := []time.Time{time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC), time.Date(2026, 10, 5, 9, 0, 5, 0, time.UTC)}
	metrics := []devtools.MonitorMetric{{Name: "CPU", Unit: "%", Samples: []float64{42, 51}}, {Name: "Memory", Unit: "GB", Samples: []float64{3.2, 3.4}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if devtools.SystemMonitor(c, times, metrics, devtools.SystemMonitorOptions{}).Changed() {
			fmt.Println("retry")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleRefreshIntervalSelector() {
	interval := 30 * time.Second
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := devtools.RefreshIntervalSelector(c, &interval, devtools.RefreshIntervalSelectorOptions{})
		if r.Refresh {
			fmt.Println("refresh")
		}
		if r.Changed {
			fmt.Println("interval", interval)
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleDashboardFilterBar() {
	var filters devtools.DashboardFilters
	opts := devtools.DashboardFilterBarOptions{Envs: []string{"production", "staging"}, Tags: []string{"api", "db"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if devtools.DashboardFilterBar(c, &filters, opts).Changed() {
			fmt.Println("filters changed")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleJsonViewer() {
	data := `{"server":{"port":8080},"ready":true}`
	var state devtools.JsonViewerState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := devtools.JsonViewer(c, data, &state, devtools.JsonViewerOptions{})
		if r.Copied != "" {
			fmt.Println("path copied")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleJsonEditor() {
	text := `{"port":8080,"ready":true}`
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := devtools.JsonEditor(c, &text, devtools.JsonEditorOptions{})
		if r.Changed {
			fmt.Println("valid", r.Valid)
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleSchemaTree() {
	roots := []devtools.SchemaNode{{Name: "app", Children: []devtools.SchemaNode{{Name: "users", Kind: devtools.SchemaTable, Children: []devtools.SchemaNode{{Name: "id", Kind: devtools.SchemaColumn, Detail: "bigint"}, {Name: "users_pk", Kind: devtools.SchemaIndex, Detail: "(id)"}}}}}}
	var state devtools.SchemaTreeState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if devtools.SchemaTree(c, &state, roots, devtools.SchemaTreeOptions{}).Changed() {
			path, _ := state.Selected()
			fmt.Println("selected", path)
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleApiRequestBuilder() {
	request := devtools.ApiRequest{Method: "POST", URL: "https://example.com/api", Headers: []devtools.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := devtools.ApiRequestBuilder(c, &request, devtools.ApiRequestBuilderOptions{})
		if r.Send {
			fmt.Println("send", request.Method)
		}
		if r.Changed {
			fmt.Println("request changed")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleResponseViewer() {
	response := &devtools.ApiResponse{Status: 200, Duration: 182 * time.Millisecond, Size: 32, Body: `{"ready":true}`, Headers: []devtools.KeyValue{{Key: "Content-Type", Value: "application/json"}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if devtools.ResponseViewer(c, response, devtools.ResponseViewerOptions{}).Changed() {
			fmt.Println("retry")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleKeyValueInput() {
	rows := []devtools.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if devtools.KeyValueInput(c, &rows, devtools.KeyValueInputOptions{}).Changed() {
			fmt.Println("rows", len(rows))
		}
	}
	ui.Render(view, 720, 480, 1)
}
