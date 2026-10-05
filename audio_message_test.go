package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var audioMessageWave = []float32{0.2, 0.5, 0.9, 0.4, 0.7, 0.3, 0.8, 0.6, 0.2, 0.5}

func TestAudioMessagePlaySeekAndKeys(t *testing.T) {
	playing, pos := false, 10*time.Second
	var seeks []time.Duration
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := AudioMessage(c, AudioMessageOptions{Duration: time.Minute, Position: pos, Playing: playing, Waveform: audioMessageWave})
		if r.Toggled() {
			playing = !playing
		}
		if d, ok := r.Seek(); ok {
			pos = d
			seeks = append(seeks, d)
		}
	}, 400, 200)
	if !tt.HasText("0:10 / 1:00") {
		t.Fatalf("times %q", tt.Texts())
	}
	_ = tt.Click("Play")
	if !playing {
		t.Fatal("play not reported")
	}
	if _, ok := tt.Find("Pause"); !ok {
		t.Fatal("button did not turn to pause")
	}
	r, _ := tt.Find("Playback position")
	tt.ClickAt(r.X+r.W*0.75, r.Y+r.H/2)
	if len(seeks) == 0 || pos < 40*time.Second || pos > 50*time.Second {
		t.Fatalf("click at 3/4 seeked to %v", pos)
	}
	tt.Key(0, ui.KeyHome)
	if pos != 0 {
		t.Fatalf("Home gave %v", pos)
	}
	tt.Key(0, ui.KeyEnd)
	if pos != time.Minute {
		t.Fatalf("End gave %v", pos)
	}
}

func TestAudioMessageDisabled(t *testing.T) {
	events := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(true).Children(func() {
			r := AudioMessage(c, AudioMessageOptions{Duration: time.Minute})
			if _, ok := r.Seek(); ok || r.Toggled() {
				events++
			}
		})
	}, 400, 200)
	_ = tt.Click("Play")
	r, _ := tt.Find("Playback position")
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if events != 0 {
		t.Fatalf("disabled audio reported %d events", events)
	}
}

func TestAudioMessageBadOptions(t *testing.T) {
	for _, o := range []AudioMessageOptions{{}, {Duration: time.Second, Position: 2 * time.Second}, {Duration: time.Second, Waveform: []float32{2}}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: no panic", o)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				Use(c, Settings{})
				AudioMessage(c, o)
			}, 100, 100)
		}()
	}
}

func TestAudioMessageSnapshot(t *testing.T) {
	snapshot(t, "audio_message", 340, 80, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			AudioMessage(c, AudioMessageOptions{Duration: 47 * time.Second, Position: 18 * time.Second, Playing: true, Waveform: audioMessageWave})
		})
	})
}
