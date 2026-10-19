package devtools

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
)

func monitorSample(n int) ([]time.Time, []MonitorMetric) {
	t0 := time.Date(2026, 10, 5, 21, 0, 0, 0, time.UTC)
	times := make([]time.Time, n)
	cpu, mem := make([]float64, n), make([]float64, n)
	for i := range n {
		times[i] = t0.Add(time.Duration(i) * 5 * time.Second)
		cpu[i] = 40 + 20*math.Sin(float64(i)/4)
		mem[i] = 3.2 + float64(i)*0.02
	}
	return times, []MonitorMetric{{Name: "CPU", Unit: "%", Samples: cpu}, {Name: "Memory", Unit: "GB", Samples: mem}}
}

func TestSystemMonitorLiveAndKeyboard(t *testing.T) {
	times, metrics := monitorSample(20)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SystemMonitor(c, times, metrics, SystemMonitorOptions{})
	}, 600, 420)
	if !tt.HasText("3.6 GB") || !tt.HasText("20.0 %") {
		t.Fatalf("latest samples missing: %q", tt.Texts())
	}
	times, metrics = monitorSample(21)
	tt.Frame()
	if !tt.HasText("20.8 %") || tt.HasText("20.0 %") {
		t.Fatalf("new sample: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab) // the legend key
	tt.Key(0, ui.KeyTab) // the plot
	tt.Key(0, ui.KeyEnd)
	found := false
	for _, s := range tt.Texts() {
		found = found || strings.Contains(s, "21:01:40")
	}
	if !found {
		t.Fatalf("keyboard readout: %q", tt.Texts())
	}
}

func TestSystemMonitorStatesAndRules(t *testing.T) {
	times, metrics := monitorSample(5)
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if SystemMonitor(c, times, metrics[:1], SystemMonitorOptions{Status: data.DataFailed, Error: "agent offline"}).Changed() {
			retries++
		}
	}, 600, 300)
	if !tt.HasText("agent offline") {
		t.Fatalf("error: %q", tt.Texts())
	}
	_ = tt.Click("Retry")
	if retries != 1 {
		t.Fatal("retry not reported")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("mismatched samples did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SystemMonitor(c, times, []MonitorMetric{{Name: "x", Samples: []float64{1}}}, SystemMonitorOptions{})
	}, 200, 200)
}

func TestSystemMonitorSnapshot(t *testing.T) {
	times, metrics := monitorSample(40)
	snapshot(t, "system_monitor", 560, 460, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { SystemMonitor(c, times, metrics, SystemMonitorOptions{}) })
	})
}

func TestSystemMonitorEmptyLoadingAndDisabled(t *testing.T) {
	status := data.DataReady
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(true).Children(func() {
			if SystemMonitor(c, nil, nil, SystemMonitorOptions{Status: status, Error: "offline"}).Changed() {
				retries++
			}
		})
	}, 600, 300)
	if !tt.HasText("No resources") {
		t.Fatal("empty resources")
	}
	status = data.DataLoading
	tt.Frame()
	if !tt.HasText("Loading") {
		t.Fatal("loading resources")
	}
	status = data.DataFailed
	tt.Frame()
	_ = tt.Click("Retry")
	if retries != 0 {
		t.Fatal("disabled retry")
	}
}

func TestSystemMonitorManySamples(t *testing.T) {
	times, metrics := monitorSample(100000)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SystemMonitor(c, times, metrics[:1], SystemMonitorOptions{})
	}, 600, 420)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if len(tt.Texts()) > 200 {
		t.Fatal("unbounded sample elements")
	}
	if !tt.HasText(times[len(times)-1].Format("15:04:05")) {
		t.Fatal("last sample readout")
	}
}
