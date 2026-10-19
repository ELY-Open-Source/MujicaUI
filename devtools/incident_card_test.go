package devtools

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func incidentSample() Incident {
	t0 := time.Date(2026, 10, 5, 21, 4, 0, 0, time.UTC)
	return Incident{Title: "Checkout errors in eu-west", Status: IncidentMonitoring, Impact: ImpactMajor, Started: t0,
		Updates: []IncidentUpdate{
			{At: t0, Status: IncidentInvestigating, Text: "Error rate rose to 8%."},
			{At: t0.Add(12 * time.Minute), Status: IncidentIdentified, Text: "A bad config push to the payment proxy."},
			{At: t0.Add(30 * time.Minute), Status: IncidentMonitoring, Text: "Rolled back; errors are falling."},
		}}
}

func TestIncidentCardShowsStatusImpactTimeline(t *testing.T) {
	inc := incidentSample()
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		IncidentCard(c, inc)
	}, 560, 500)
	for _, s := range []string{"Checkout errors in eu-west", "Monitoring", "Impact: major", "Started 2026-10-05 21:04",
		"Investigating", "21:16", "Rolled back; errors are falling."} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	inc.Status = IncidentResolved
	tt.Frame()
	if !tt.HasText("Resolved") {
		t.Fatal("status did not update")
	}
}

func TestIncidentCardRejectsUnknownStatus(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		IncidentCard(c, Incident{Title: "x", Status: 9})
	}, 200, 200)
}

func TestIncidentCardSnapshot(t *testing.T) {
	snapshot(t, "incident_card", 540, 360, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { IncidentCard(c, incidentSample()) })
	})
}

func TestIncidentCardManyUpdates(t *testing.T) {
	inc := Incident{Title: "Recovery", Updates: make([]IncidentUpdate, 100000)}
	for i := range inc.Updates {
		inc.Updates[i].Text = "Update received"
	}
	tt := ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); IncidentCard(c, inc) }, 560, 500)
	if n := len(tt.Texts()); n > 500 {
		t.Fatalf("built %d texts for 100000 updates", n)
	}
}
