package media

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

type mediaCalls struct {
	play, pause int
	seeks       []time.Duration
}

func (m *mediaCalls) opts() AudioPlayerOptions {
	return AudioPlayerOptions{OnPlay: func() { m.play++ }, OnPause: func() { m.pause++ }, OnSeek: func(d time.Duration) { m.seeks = append(m.seeks, d) }}
}

func audioTester(st *AudioState, m *mediaCalls, disabled bool) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		o := m.opts()
		o.Disabled = disabled
		ui.Column(c).Padding(20).Children(func() { AudioPlayer(c, *st, o) })
	}, 600, 200)
}

func TestAudioPlayerPlayPause(t *testing.T) {
	st := AudioState{Title: "Ave Mujica", Position: 30 * time.Second, Duration: 3 * time.Minute}
	var m mediaCalls
	tt := audioTester(&st, &m, false)
	if err := tt.Click("Play"); err != nil {
		t.Fatal(err)
	}
	st.Playing = true
	tt.Frame()
	if err := tt.Click("Pause"); err != nil {
		t.Fatal(err)
	}
	if m.play != 1 || m.pause != 1 {
		t.Fatalf("play %d pause %d", m.play, m.pause)
	}
	if !tt.HasText("0:30 / 3:00") {
		t.Fatalf("texts %v", tt.Texts())
	}
}

func TestAudioPlayerKeyboard(t *testing.T) {
	st := AudioState{Title: "Ave Mujica", Position: 30 * time.Second, Duration: 3 * time.Minute}
	var m mediaCalls
	tt := audioTester(&st, &m, false)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if len(m.seeks) != 1 || m.seeks[0] != 35*time.Second {
		t.Fatalf("seeks %v", m.seeks)
	}
	tt.Key(0, ui.KeySpace)
	if m.play != 1 {
		t.Fatalf("Space on the timeline: play %d", m.play)
	}
}

func TestAudioPlayerDisabled(t *testing.T) {
	st := AudioState{Title: "Ave Mujica", Duration: time.Minute}
	var m mediaCalls
	tt := audioTester(&st, &m, true)
	_ = tt.Click("Play")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyRight)
	if m.play != 0 || len(m.seeks) != 0 {
		t.Fatalf("disabled player reported %+v", m)
	}
}

func TestAudioPlayerNeedsCallbacks(t *testing.T) {
	mediaPanics(t, "no callbacks", func(c *ui.Context) { AudioPlayer(c, AudioState{}, AudioPlayerOptions{}) })
}

func TestAudioPlayerSnapshot(t *testing.T) {
	var m mediaCalls
	snapshot(t, "audio_player", 560, 200, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).Children(func() {
			AudioPlayer(c, AudioState{Title: "Ave Mujica — Overture", Position: 72 * time.Second, Duration: 4 * time.Minute, Playing: true}, m.opts())
			AudioPlayer(c, AudioState{Title: "Imprisoned XII", Duration: 5 * time.Minute}, m.opts())
		})
	})
}
