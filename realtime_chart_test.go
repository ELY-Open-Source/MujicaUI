package mujica

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var rtStart = time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)

func TestRealtimeBufferTrims(t *testing.T) {
	b := NewRealtimeBuffer(100, "a")
	for i := range 100_000 {
		b.Append(rtStart.Add(time.Duration(i)*time.Second), float64(i))
	}
	if b.Len() != 100 || len(b.times) >= 200 || b.vals[0][len(b.vals[0])-1] != 99_999 {
		t.Fatalf("kept %d samples, %d stored", b.Len(), len(b.times))
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a sample back in time")
		}
	}()
	b.Append(rtStart, 1)
}

func TestRealtimeChartWindowAndPause(t *testing.T) {
	buf := NewRealtimeBuffer(1000, "Bell", "Organ")
	var st RealtimeChartState
	values := 0
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		r := RealtimeChart(c, &st, buf, RealtimeChartOptions{Frame: ChartFrame{Label: "Sound"}, Window: 10 * time.Second})
		if r.ValueChanged() {
			values++
		}
		return r
	}}
	tt := ui.NewTester(h.view, 520, 440)
	if !tt.HasText("No data") {
		t.Fatalf("empty buffer: %q", tt.Texts())
	}
	for i := range 30 {
		buf.Append(rtStart.Add(time.Duration(i)*time.Second), float64(i), float64(2*i))
	}
	tt.Frame()
	tabTo(t, tt, "Sound")
	tt.Key(0, ui.KeyHome)
	if got := chartReadout(tt); got != "2026-10-05 09:00:19 · Bell 19 · Organ 38" {
		t.Fatalf("window start %q", got)
	}
	if err := tt.Click("Pause"); err != nil || !st.Paused || values != 1 {
		t.Fatalf("pause: %v %v %d", err, st.Paused, values)
	}
	for i := 30; i < 40; i++ {
		buf.Append(rtStart.Add(time.Duration(i)*time.Second), float64(i), 0)
	}
	tt.Frame()
	tabTo(t, tt, "Sound")
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); !strings.HasPrefix(got, "2026-10-05 09:00:29 ·") {
		t.Fatalf("paused view moved: %q", got)
	}
	_ = tt.Click("Resume")
	tabTo(t, tt, "Sound")
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); !strings.HasPrefix(got, "2026-10-05 09:00:39 ·") || st.Paused {
		t.Fatalf("resumed view %q", got)
	}
}

func TestRealtimeChartSnapshot(t *testing.T) {
	buf := NewRealtimeBuffer(200, "Bell", "Organ")
	for i := range 60 {
		x := float64(i)
		buf.Append(rtStart.Add(time.Duration(i)*time.Second), 5+3*math.Sin(x/5), 4+2*math.Sin(x/3+1))
	}
	snapshot(t, "realtime_chart", 520, 380, func(c *ui.Context) {
		var st RealtimeChartState
		ui.Column(c).Padding(16).Children(func() {
			RealtimeChart(c, &st, buf, RealtimeChartOptions{Frame: ChartFrame{Label: "Sound"}, Window: 30 * time.Second})
		})
	})
}
