package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

type recorderCalls struct{ start, pause, resume, stop int }

func (r *recorderCalls) opts() ScreenRecorderControlsOptions {
	return ScreenRecorderControlsOptions{OnStart: func() { r.start++ }, OnPause: func() { r.pause++ },
		OnResume: func() { r.resume++ }, OnStop: func() { r.stop++ }}
}

func TestScreenRecorderControlsFlow(t *testing.T) {
	st := RecorderState{}
	var r recorderCalls
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() { ScreenRecorderControls(c, st, r.opts()) })
	}, 500, 120)
	if !tt.HasText("Ready") || tt.HasText("Stop") {
		t.Fatalf("idle shows %v", tt.Texts())
	}
	if err := tt.Click("Record"); err != nil {
		t.Fatal(err)
	}
	st = RecorderState{Status: RecorderRecording, Elapsed: 65 * time.Second}
	tt.Frame()
	if !tt.HasText("1:05") || !tt.HasText("Recording") {
		t.Fatalf("recording shows %v", tt.Texts())
	}
	if err := tt.Click("Pause"); err != nil {
		t.Fatal(err)
	}
	st.Status = RecorderPaused
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if r != (recorderCalls{1, 1, 1, 1}) {
		t.Fatalf("calls %+v", r)
	}
}

func TestScreenRecorderControlsDisabledAndBad(t *testing.T) {
	var r recorderCalls
	o := r.opts()
	o.Disabled = true
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ScreenRecorderControls(c, RecorderState{}, o)
	}, 500, 120)
	_ = tt.Click("Record")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if r.start != 0 {
		t.Fatal("disabled recorder started")
	}
	mediaPanics(t, "no callbacks", func(c *ui.Context) { ScreenRecorderControls(c, RecorderState{}, ScreenRecorderControlsOptions{}) })
	mediaPanics(t, "bad status", func(c *ui.Context) { ScreenRecorderControls(c, RecorderState{Status: 9}, r.opts()) })
}

func TestScreenRecorderControlsSnapshot(t *testing.T) {
	var r recorderCalls
	snapshot(t, "screen_recorder_controls", 420, 140, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).AlignItems(ui.Start).Children(func() {
			ScreenRecorderControls(c, RecorderState{}, r.opts())
			ScreenRecorderControls(c, RecorderState{Status: RecorderRecording, Elapsed: 3*time.Minute + 7*time.Second}, r.opts())
		})
	})
}
