package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

type videoCalls struct {
	mediaCalls
	rates           []float64
	subs, full, pip []bool
}

func (v *videoCalls) opts() VideoPlayerOptions {
	a := v.mediaCalls.opts()
	return VideoPlayerOptions{OnPlay: a.OnPlay, OnPause: a.OnPause, OnSeek: a.OnSeek,
		OnRate:       func(r float64) { v.rates = append(v.rates, r) },
		OnSubtitles:  func(on bool) { v.subs = append(v.subs, on) },
		OnFullscreen: func(on bool) { v.full = append(v.full, on) },
		OnPiP:        func(on bool) { v.pip = append(v.pip, on) }}
}

func videoTester(st *VideoState, o *VideoPlayerOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(640).Children(func() { VideoPlayer(c, *st, *o) })
	}, 700, 520)
}

func TestVideoPlayerStageKeys(t *testing.T) {
	st := VideoState{Frame: mediaTestBitmap(160, 90), Position: 20 * time.Second, Duration: time.Minute}
	var v videoCalls
	o := v.opts()
	tt := videoTester(&st, &o)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Video") {
		t.Fatal("stage not first in Tab order")
	}
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyF)
	tt.Key(0, ui.KeyC)
	if v.play != 1 || len(v.seeks) != 1 || v.seeks[0] != 25*time.Second || len(v.full) != 1 || !v.full[0] || len(v.subs) != 1 || !v.subs[0] {
		t.Fatalf("calls %+v", v)
	}
}

func TestVideoPlayerButtons(t *testing.T) {
	st := VideoState{Position: 0, Duration: time.Minute, Fullscreen: true, Rate: 1}
	var v videoCalls
	o := v.opts()
	tt := videoTester(&st, &o)
	if !tt.HasText("No video") {
		t.Fatal("empty stage not shown")
	}
	for _, name := range []string{"Picture in picture", "Full screen", "Subtitles"} {
		if err := tt.Click(name); err != nil {
			t.Fatal(err)
		}
	}
	if len(v.pip) != 1 || !v.pip[0] || len(v.full) != 1 || v.full[0] || len(v.subs) != 1 || !v.subs[0] {
		t.Fatalf("toggles %+v", v)
	}
	if err := tt.Click("Playback speed"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("1.5×"); err != nil {
		t.Fatal(err)
	}
	if len(v.rates) != 1 || v.rates[0] != 1.5 {
		t.Fatalf("rates %v", v.rates)
	}
}

func TestVideoPlayerSubtitleAndOptionalButtons(t *testing.T) {
	st := VideoState{Duration: time.Minute, Subtitle: "Welcome to the masquerade"}
	var v videoCalls
	o := v.opts()
	o.OnPiP, o.OnRate = nil, nil
	tt := videoTester(&st, &o)
	if tt.HasText("Welcome to the masquerade") {
		t.Fatal("subtitle shown while off")
	}
	if _, ok := tt.Find("Picture in picture"); ok {
		t.Fatal("PiP button without OnPiP")
	}
	st.SubtitlesOn = true
	tt.Frame()
	if !tt.HasText("Welcome to the masquerade") {
		t.Fatal("subtitle not shown")
	}
}

func TestVideoPlayerDisabled(t *testing.T) {
	st := VideoState{Duration: time.Minute}
	var v videoCalls
	o := v.opts()
	o.Disabled = true
	tt := videoTester(&st, &o)
	_ = tt.Click("Play")
	_ = tt.Click("Full screen")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if v.play != 0 || len(v.full) != 0 {
		t.Fatalf("disabled player reported %+v", v)
	}
}

func TestVideoPlayerSnapshot(t *testing.T) {
	var v videoCalls
	snapshot(t, "video_player", 560, 400, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(520).Children(func() {
			VideoPlayer(c, VideoState{Frame: mediaTestBitmap(320, 180), Position: 95 * time.Second, Duration: 4 * time.Minute,
				Buffered: []MediaRange{{0, 150 * time.Second}}, Playing: true, Rate: 1, SubtitlesOn: true,
				Subtitle: "The curtain rises."}, v.opts())
		})
	})
}
