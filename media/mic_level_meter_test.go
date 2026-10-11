package media

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestMicLevelMeterPeakHold(t *testing.T) {
	level := float32(1)
	var e *ui.Element
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		e = MicLevelMeter(c, level, MicLevelMeterOptions{Hold: 40 * time.Millisecond})
	}, 400, 100)
	if !tt.HasText("0 dB") {
		t.Fatalf("texts %v", tt.Texts())
	}
	level = 0.1
	tt.Frame()
	if !tt.HasText("-20 dB") {
		t.Fatalf("level readout %v", tt.Texts())
	}
	if r, ok := tt.Find("Microphone level"); !ok || r.W == 0 {
		t.Fatal("meter not named")
	}
	if !strings.Contains(micTestDescription(e), "peak 0 dB") {
		t.Fatal("peak not held")
	}
	time.Sleep(60 * time.Millisecond)
	tt.Frame()
	if !strings.Contains(micTestDescription(e), "peak -20 dB") {
		t.Fatal("peak did not fall after the hold")
	}
}

// micTestDescription rebuilds the meter's spoken value from its state.
func micTestDescription(e *ui.Element) string {
	pk := *ui.Local(e, micPeak{}, func() micPeak { return micPeak{} })
	return "peak " + micDB(pk.level)
}

func TestMicLevelMeterDB(t *testing.T) {
	for v, want := range map[float32]string{0: "-∞ dB", 1: "0 dB", 0.5: "-6 dB"} {
		if got := micDB(v); got != want {
			t.Errorf("%v: %s, want %s", v, got, want)
		}
	}
	mediaPanics(t, "over full scale", func(c *ui.Context) { MicLevelMeter(c, 1.2, MicLevelMeterOptions{}) })
}

func TestMicLevelMeterSnapshot(t *testing.T) {
	snapshot(t, "mic_level_meter", 360, 150, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(8).Children(func() {
			MicLevelMeter(c, 0.35, MicLevelMeterOptions{})
			MicLevelMeter(c, 0.8, MicLevelMeterOptions{})
			MicLevelMeter(c, 0.97, MicLevelMeterOptions{})
		})
	})
}
