package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

var servicesSample = []ServiceItem{
	{Name: "API", Health: ServiceOperational},
	{Name: "Search", Health: ServiceDegraded, Note: "Indexing is slow"},
	{Name: "Payments", Health: ServiceOutage, Note: "Provider unreachable"},
}

func TestServiceStatusSummaryAndRows(t *testing.T) {
	services := servicesSample
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ServiceStatus(c, services, ServiceStatusOptions{})
	}, 520, 400)
	for _, s := range []string{"Some services down", "1 up · 1 degraded · 1 down", "Degraded", "Outage", "Provider unreachable"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	services = servicesSample[:1]
	tt.Frame()
	if !tt.HasText("All systems operational") {
		t.Fatalf("all ok: %q", tt.Texts())
	}
	services = nil
	tt.Frame()
	if !tt.HasText("No services") {
		t.Fatal("empty")
	}
}

func TestServiceStatusManyVirtualized(t *testing.T) {
	many := make([]ServiceItem, 100000)
	for i := range many {
		many[i] = ServiceItem{Name: fmt.Sprint("svc ", i)}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ServiceStatus(c, many, ServiceStatusOptions{})
	}, 520, 400)
	if n := len(tt.Texts()); n > 300 {
		t.Fatalf("%d texts", n)
	}
}

func TestServiceStatusSnapshot(t *testing.T) {
	snapshot(t, "service_status", 520, 280, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { ServiceStatus(c, servicesSample, ServiceStatusOptions{}) })
	})
}
