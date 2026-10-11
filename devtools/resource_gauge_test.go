package devtools

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

var resourceSample = []ResourceReading{
	{Name: "CPU", Used: 45, Total: 100},
	{Name: "Memory", Used: 6.4, Total: 8, Unit: "GB"},
	{Name: "Disk", Used: 470, Total: 500, Unit: "GB"},
}

func TestResourceGaugeThresholdsAndKeyboard(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ResourceGauge(c, resourceSample, ResourceGaugeOptions{})
	}, 700, 300)
	if !tt.HasText("6.4 / 8 GB") || !tt.HasText("470 / 500 GB") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab) // CPU dial
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("45 % · Normal") {
		t.Fatalf("CPU readout: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("80 % · High") {
		t.Fatalf("memory readout: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("94 % · Critical") {
		t.Fatalf("disk readout: %q", tt.Texts())
	}
}

func TestResourceGaugeRejectsBadReading(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ResourceGauge(c, []ResourceReading{{Name: "CPU", Used: 5, Total: 4}}, ResourceGaugeOptions{})
	}, 200, 200)
}

func TestResourceGaugeSnapshot(t *testing.T) {
	snapshot(t, "resource_gauge", 680, 250, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { ResourceGauge(c, resourceSample, ResourceGaugeOptions{}) })
	})
}

func TestResourceGaugeRejectsInfiniteTotal(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("infinite capacity accepted")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ResourceGauge(c, []ResourceReading{{Name: "CPU", Total: math.Inf(1)}}, ResourceGaugeOptions{})
	}, 400, 300)
}

func TestResourceGaugeEmptyState(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); ResourceGauge(c, nil, ResourceGaugeOptions{}) }, 600, 300)
	if !tt.HasText("No resources") {
		t.Fatal("empty resource state")
	}
}

func TestResourceGaugeLoadingErrorAndDisabledRetry(t *testing.T) {
	status, disabled := data.DataLoading, false
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			if ResourceGauge(c, resourceSample[:1], ResourceGaugeOptions{Status: status, Error: "offline"}).Changed() {
				retries++
			}
		})
	}, 400, 300)
	if !tt.HasText("Loading") {
		t.Fatal("loading resource")
	}
	status = data.DataFailed
	tt.Frame()
	_ = tt.Click("Retry")
	if retries != 1 || !tt.HasText("offline") {
		t.Fatal("error retry")
	}
	disabled = true
	tt.Key(0, ui.KeyEnter)
	if retries != 1 {
		t.Fatal("disabled resource retry")
	}
}
