package mujica_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleStatCard() {
	value := "12,480"
	opts := mujica.StatCardOptions{Label: "Requests", Change: 0.12, HasChange: true, Trend: []float64{3, 5, 4, 8}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.StatCard(c, value, opts).Changed() {
			fmt.Println("retry")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleAlertList() {
	alerts := []mujica.AlertItem{{ID: "cpu", Title: "CPU usage high", Level: mujica.AlertWarning, Source: "worker", At: time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.AlertList(c, alerts, mujica.AlertListOptions{})
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
	incident := mujica.Incident{Title: "API latency", Status: mujica.IncidentMonitoring, Impact: mujica.ImpactMinor, Updates: []mujica.IncidentUpdate{{Status: mujica.IncidentIdentified, Text: "Slow query located"}, {Status: mujica.IncidentMonitoring, Text: "Recovery observed"}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.IncidentCard(c, incident)
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleUptimeBar() {
	days := []mujica.UptimeDay{{Date: mujica.Date{Year: 2026, Month: 10, Day: 1}, Uptime: 1}, {Date: mujica.Date{Year: 2026, Month: 10, Day: 2}, Uptime: 0.91, Note: "Recovery"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.UptimeBar(c, days, mujica.UptimeBarOptions{Frame: mujica.ChartFrame{Label: "API uptime"}})
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
	services := []mujica.ServiceItem{{Name: "API", Health: mujica.ServiceOperational}, {Name: "Search", Health: mujica.ServiceDegraded, Note: "Slow queries"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ServiceStatus(c, services, mujica.ServiceStatusOptions{})
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleLogStream() {
	lines := []mujica.LogLine{{Level: mujica.LogInfo, Text: "server ready"}, {Level: mujica.LogWarn, Text: "slow query"}}
	var state mujica.LogStreamState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.LogStream(c, lines, &state, mujica.LogStreamOptions{}).Changed() {
			fmt.Println("filter changed")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleTraceWaterfall() {
	spans := []mujica.WaterfallSpan{{ID: "root", Name: "GET /api", Service: "gateway", Duration: 200 * time.Millisecond}, {ID: "db", Parent: "root", Name: "Query", Service: "database", Start: 20 * time.Millisecond, Duration: 80 * time.Millisecond}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.TraceWaterfall(c, spans, mujica.TraceWaterfallOptions{}).Changed() {
			fmt.Println("span selected")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleResourceGauge() {
	readings := []mujica.ResourceReading{{Name: "CPU", Used: 45, Total: 100}, {Name: "Memory", Used: 6.4, Total: 8, Unit: "GB"}, {Name: "Disk", Used: 470, Total: 500, Unit: "GB"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ResourceGauge(c, readings, mujica.ResourceGaugeOptions{}).Changed() {
			fmt.Println("retry")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleSystemMonitor() {
	times := []time.Time{time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC), time.Date(2026, 10, 5, 9, 0, 5, 0, time.UTC)}
	metrics := []mujica.MonitorMetric{{Name: "CPU", Unit: "%", Samples: []float64{42, 51}}, {Name: "Memory", Unit: "GB", Samples: []float64{3.2, 3.4}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.SystemMonitor(c, times, metrics, mujica.SystemMonitorOptions{}).Changed() {
			fmt.Println("retry")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleRefreshIntervalSelector() {
	interval := 30 * time.Second
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.RefreshIntervalSelector(c, &interval, mujica.RefreshIntervalSelectorOptions{})
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
	var filters mujica.DashboardFilters
	opts := mujica.DashboardFilterBarOptions{Envs: []string{"production", "staging"}, Tags: []string{"api", "db"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.DashboardFilterBar(c, &filters, opts).Changed() {
			fmt.Println("filters changed")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleJsonViewer() {
	data := `{"server":{"port":8080},"ready":true}`
	var state mujica.JsonViewerState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.JsonViewer(c, data, &state, mujica.JsonViewerOptions{})
		if r.Copied != "" {
			fmt.Println("path copied")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleJsonEditor() {
	text := `{"port":8080,"ready":true}`
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.JsonEditor(c, &text, mujica.JsonEditorOptions{})
		if r.Changed {
			fmt.Println("valid", r.Valid)
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleSchemaTree() {
	roots := []mujica.SchemaNode{{Name: "app", Children: []mujica.SchemaNode{{Name: "users", Kind: mujica.SchemaTable, Children: []mujica.SchemaNode{{Name: "id", Kind: mujica.SchemaColumn, Detail: "bigint"}, {Name: "users_pk", Kind: mujica.SchemaIndex, Detail: "(id)"}}}}}}
	var state mujica.SchemaTreeState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.SchemaTree(c, &state, roots, mujica.SchemaTreeOptions{}).Changed() {
			path, _ := state.Selected()
			fmt.Println("selected", path)
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleApiRequestBuilder() {
	request := mujica.ApiRequest{Method: "POST", URL: "https://example.com/api", Headers: []mujica.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.ApiRequestBuilder(c, &request, mujica.ApiRequestBuilderOptions{})
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
	response := &mujica.ApiResponse{Status: 200, Duration: 182 * time.Millisecond, Size: 32, Body: `{"ready":true}`, Headers: []mujica.KeyValue{{Key: "Content-Type", Value: "application/json"}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ResponseViewer(c, response, mujica.ResponseViewerOptions{}).Changed() {
			fmt.Println("retry")
		}
	}
	ui.Render(view, 720, 480, 1)
}

func ExampleKeyValueInput() {
	rows := []mujica.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.KeyValueInput(c, &rows, mujica.KeyValueInputOptions{}).Changed() {
			fmt.Println("rows", len(rows))
		}
	}
	ui.Render(view, 720, 480, 1)
}
