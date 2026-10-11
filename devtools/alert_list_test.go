package devtools

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func alertSample() []AlertItem {
	now := time.Now()
	return []AlertItem{
		{ID: "a", Title: "Disk 92% full", Level: AlertCritical, Source: "db-1", At: now.Add(-5 * time.Minute)},
		{ID: "b", Title: "Latency above 500 ms", Level: AlertWarning, Source: "api", At: now.Add(-time.Hour), Acked: true},
		{ID: "c", Title: "Deploy finished", Level: AlertInfo, Source: "ci", Silence: now.Add(time.Hour)},
	}
}

func alertTester(alerts []AlertItem, opts AlertListOptions, got *[]AlertListResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := AlertList(c, alerts, opts)
		if r.Ack != "" || r.Silence != "" {
			*got = append(*got, r)
		}
	}, 640, 400)
}

func TestAlertListLevelsAndActions(t *testing.T) {
	var got []AlertListResult
	tt := alertTester(alertSample(), AlertListOptions{}, &got)
	for _, s := range []string{"Critical", "Warning", "Info", "db-1 · 5 m 00 s ago", "Acknowledged"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	_ = tt.Click("Acknowledge")
	silences, _ := tt.Find("Silence 1 h")
	tt.ClickAt(silences.X+4, silences.Y+4)
	if len(got) != 2 || got[0].Ack != "a" || got[1].Silence != "a" {
		t.Fatalf("actions %+v", got)
	}
	tt = alertTester(alertSample(), AlertListOptions{}, &got)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(got) != 3 || got[2].Ack != "a" {
		t.Fatalf("keyboard ack: %+v", got)
	}
}

func TestAlertListDisabledEmptyMany(t *testing.T) {
	var got []AlertListResult
	tt := alertTester(alertSample(), AlertListOptions{Disabled: true}, &got)
	_ = tt.Click("Acknowledge")
	_ = tt.Click("Silence 1 h")
	if len(got) != 0 {
		t.Fatalf("disabled list acted: %+v", got)
	}
	tt = alertTester(nil, AlertListOptions{}, &got)
	if !tt.HasText("No alerts") {
		t.Fatal("empty text missing")
	}
	many := make([]AlertItem, 100000)
	for i := range many {
		many[i] = AlertItem{ID: fmt.Sprint(i), Title: fmt.Sprintf("alert %d", i)}
	}
	tt = alertTester(many, AlertListOptions{}, &got)
	if n := len(tt.Texts()); n > 300 {
		t.Fatalf("%d texts for 100000 alerts", n)
	}
}

func TestAlertListSnapshot(t *testing.T) {
	snapshot(t, "alert_list", 640, 260, func(c *ui.Context) {
		alerts := alertSample()
		for i := range alerts {
			alerts[i].At = time.Time{}
		}
		alerts[2].Silence = time.Time{}
		ui.Column(c).Padding(16).Children(func() { AlertList(c, alerts, AlertListOptions{}) })
	})
}

func TestAlertListRejectsDuplicateIDs(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate alerts accepted")
		}
	}()
	var results []AlertListResult
	alertTester([]AlertItem{{ID: "a", Title: "One"}, {ID: "a", Title: "Two"}}, AlertListOptions{}, &results)
}
