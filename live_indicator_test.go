package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var liveNow = time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)

type liveProbe struct {
	status  LiveStatus
	err     string
	retried int
}

func (p *liveProbe) view(c *ui.Context) {
	Use(c, Settings{})
	ui.Column(c).Padding(12).Children(func() {
		if LiveIndicator(c, p.status, LiveIndicatorOptions{LastSynced: liveNow.Add(-3 * time.Minute), Now: liveNow, Error: p.err}).Changed() {
			p.retried++
		}
	})
}

func TestLiveIndicatorStates(t *testing.T) {
	p := &liveProbe{}
	tt := ui.NewTester(p.view, 500, 120)
	if !tt.HasText("All changes saved") || !tt.HasText("3 min ago") || tt.HasText("Retry") {
		t.Fatalf("synced: %q", tt.Texts())
	}
	p.status = LiveOffline
	tt.Frame()
	if !tt.HasText("Offline, changes kept on this device") {
		t.Fatalf("offline: %q", tt.Texts())
	}
	ann := tt.Announcements()
	if len(ann) == 0 || ann[len(ann)-1] != "Offline, changes kept on this device" {
		t.Fatalf("not announced: %q", ann)
	}
	if err := tt.Click("Retry"); err != nil {
		t.Fatal(err)
	}
	p.status, p.err = LiveError, "server refused"
	tt.Frame()
	if !tt.HasText("Sync failed: server refused") {
		t.Fatalf("error: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if p.retried != 2 {
		t.Fatalf("retried %d", p.retried)
	}
	p.status = LiveSyncing
	tt.Frame()
	if !tt.HasText("Saving…") || tt.HasText("Retry") {
		t.Fatalf("syncing: %q", tt.Texts())
	}
}

func TestLiveAgo(t *testing.T) {
	var got []string
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		got = []string{liveAgo(c, 10*time.Second), liveAgo(c, 5*time.Hour), liveAgo(c, 50*time.Hour)}
	}, 10, 10)
	if got[0] != "just now" || got[1] != "5 h ago" || got[2] != "2 days ago" {
		t.Fatalf("ago %q", got)
	}
}

func TestLiveIndicatorInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("bad status did not panic")
		}
	}()
	ui.NewTester((&liveProbe{status: 9}).view, 400, 100)
}

func TestLiveIndicatorSnapshot(t *testing.T) {
	snapshot(t, "live_indicator", 520, 150, func(c *ui.Context) {
		ui.Column(c).Gap(4).Children(func() {
			LiveIndicator(c, LiveSynced, LiveIndicatorOptions{LastSynced: liveNow.Add(-3 * time.Minute), Now: liveNow})
			LiveIndicator(c, LiveOffline, LiveIndicatorOptions{})
		})
	})
}
