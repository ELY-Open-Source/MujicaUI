package mujica

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func waveSamples(n int) []float32 {
	s := make([]float32, n)
	for i := range s {
		s[i] = float32(math.Sin(float64(i)/40) * (0.3 + 0.6*math.Abs(math.Sin(float64(i)/3000))))
	}
	return s
}

func waveTester(samples []float32, opts *AudioWaveformOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(20).Width(400).Children(func() {
			if AudioWaveform(c, samples, *opts).Changed() {
				*n++
			}
		})
	}, 500, 200)
}

func TestAudioWaveformSeekAndSelect(t *testing.T) {
	var seeks []time.Duration
	sel := MediaRange{}
	opts := AudioWaveformOptions{Duration: 100 * time.Second, Position: 10 * time.Second, Selection: &sel,
		OnSeek: func(d time.Duration) { seeks = append(seeks, d) }}
	n := 0
	tt := waveTester(waveSamples(5000), &opts, &n)
	r, ok := tt.Find("Waveform")
	if !ok {
		t.Fatal("not named")
	}
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if len(seeks) != 1 || seeks[0] < 48*time.Second || seeks[0] > 52*time.Second {
		t.Fatalf("click seeks %v", seeks)
	}
	tt.Press(r.X+r.W*0.2, r.Y+r.H/2)
	tt.Move(r.X+r.W*0.3, r.Y+r.H/2)
	tt.Move(r.X+r.W*0.4, r.Y+r.H/2)
	tt.Release(r.X+r.W*0.4, r.Y+r.H/2)
	if sel.Start < 18*time.Second || sel.Start > 22*time.Second || sel.End < 38*time.Second || sel.End > 42*time.Second || len(seeks) != 1 {
		t.Fatalf("drag selection %+v, seeks %v", sel, seeks)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Key(ui.Shift, ui.KeyRight)
	if sel != (MediaRange{10 * time.Second, 15 * time.Second}) || n < 3 {
		t.Fatalf("Shift+Right selection %+v, changes %d", sel, n)
	}
	tt.Key(0, ui.KeyRight)
	if seeks[len(seeks)-1] != 15*time.Second {
		t.Fatalf("Right seek %v", seeks)
	}
}

func TestAudioWaveformComputesPeaksOnce(t *testing.T) {
	samples := waveSamples(100000)
	opts := AudioWaveformOptions{Duration: time.Minute, OnSeek: func(time.Duration) {}}
	n := 0
	tt := waveTester(samples, &opts, &n)
	before := waveformBuckets
	for range 5 {
		opts.Position += time.Second
		tt.Frame()
	}
	if got := waveformBuckets - before; got != 0 {
		t.Fatalf("peaks rebuilt %d times without new data", got)
	}
	samples = waveSamples(100000)
	tt = waveTester(samples, &opts, &n)
	tt.Frame()
	if waveformBuckets-before == 0 {
		t.Fatal("new data did not rebuild peaks")
	}
}

func TestAudioWaveformStatesKeepSize(t *testing.T) {
	var got []ui.Rect
	seek := func(time.Duration) {}
	for _, o := range []AudioWaveformOptions{{Duration: time.Second, OnSeek: seek}, {Status: DataLoading}, {Status: DataFailed, Error: "Decode failed"}} {
		var r ui.Rect
		data := waveSamples(100)
		if o.Status != DataReady {
			data = nil
		}
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			ui.Column(c).Width(300).Children(func() { r = AudioWaveform(c, data, o).Element.Bounds() })
		}, 400, 200)
		tt.Frame()
		got = append(got, r)
	}
	if got[1] != got[0] || got[2] != got[0] {
		t.Fatalf("state sizes %v", got)
	}
}

func TestAudioWaveformDisabled(t *testing.T) {
	var seeks []time.Duration
	opts := AudioWaveformOptions{Duration: time.Minute, OnSeek: func(d time.Duration) { seeks = append(seeks, d) }}
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(true).Width(300).Children(func() { AudioWaveform(c, waveSamples(500), opts) })
	}, 400, 200)
	r, _ := tt.Find("Waveform")
	tt.ClickAt(r.X+10, r.Y+10)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if len(seeks) != 0 || n != 0 {
		t.Fatalf("disabled waveform seeked %v", seeks)
	}
}

func TestAudioWaveformSnapshot(t *testing.T) {
	samples := waveSamples(20000)
	sel := MediaRange{50 * time.Second, 70 * time.Second}
	snapshot(t, "audio_waveform", 460, 150, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(420).Children(func() {
			AudioWaveform(c, samples, AudioWaveformOptions{Duration: 2 * time.Minute, Position: 40 * time.Second,
				Selection: &sel, OnSeek: func(time.Duration) {}})
		})
	})
}
