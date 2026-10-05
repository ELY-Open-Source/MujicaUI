package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func voiceLevels(n int, v float32) []float32 {
	out := make([]float32, n)
	for i := range out {
		out[i] = v
	}
	return out
}

func TestVoiceWaveformSilenceHint(t *testing.T) {
	levels := append(voiceLevels(30, 0.6), voiceLevels(10, 0.01)...)
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		VoiceWaveform(c, levels, VoiceWaveformOptions{SilentFor: 12})
	}, 400, 100)
	if tt.HasText("No sound detected") {
		t.Fatal("hint before 12 silent samples")
	}
	levels = append(levels, voiceLevels(2, 0.02)...)
	tt.Frame()
	if !tt.HasText("No sound detected") {
		t.Fatal("no hint after 12 silent samples")
	}
	if _, ok := tt.Find("Input level"); !ok {
		t.Fatal("meter not named")
	}
}

func TestVoiceWaveformBadLevels(t *testing.T) {
	for _, l := range [][]float32{{-0.1}, {1.5}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%v: no panic", l)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				Use(c, Settings{})
				VoiceWaveform(c, l, VoiceWaveformOptions{})
			}, 100, 100)
		}()
	}
}

func TestVoiceWaveformSnapshot(t *testing.T) {
	levels := make([]float32, 60)
	for i := range levels {
		levels[i] = float32((i*7)%10) / 10
	}
	snapshot(t, "voice_waveform", 420, 110, func(c *ui.Context) {
		ui.Column(c).Padding(12).Gap(12).Children(func() {
			VoiceWaveform(c, levels, VoiceWaveformOptions{})
			VoiceWaveform(c, voiceLevels(40, 0), VoiceWaveformOptions{Width: 160})
		})
	})
}
