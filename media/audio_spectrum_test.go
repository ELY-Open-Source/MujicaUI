package media

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
)

func TestAudioSpectrumKeyboardReadout(t *testing.T) {
	bands := []float32{0.2, 0.9, 0.5}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(300).Children(func() {
			AudioSpectrum(c, bands, AudioSpectrumOptions{BandLabel: func(i int) string { return []string{"63 Hz", "1 kHz", "8 kHz"}[i] }})
		})
	}, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnd)
	if got := tt.Announcements(); !slices.Equal(got, []string{"63 Hz: 20%", "1 kHz: 90%", "8 kHz: 50%"}) {
		t.Fatalf("announced %v", got)
	}
}

func TestAudioSpectrumNeverAnimates(t *testing.T) {
	builds := 0
	bands := []float32{0.3, 0.6}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		builds++
		AudioSpectrum(c, bands, AudioSpectrumOptions{})
	}, 300, 200)
	start := builds
	tt.Frame()
	if builds-start != 1 {
		t.Fatalf("one frame built %d times: the spectrum asked for more", builds-start)
	}
}

func TestAudioSpectrumStates(t *testing.T) {
	var got []ui.Rect
	for _, o := range []AudioSpectrumOptions{{}, {Status: data.DataLoading}, {Status: data.DataFailed, Error: "Device lost"}} {
		var r ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			ui.Column(c).Width(300).Children(func() { r = AudioSpectrum(c, nil, o).Bounds() })
		}, 400, 200)
		tt.Frame()
		got = append(got, r)
		if o.Status == data.DataFailed && !tt.HasText("Device lost") {
			t.Fatal("error not shown")
		}
		if o.Status == data.DataReady && !tt.HasText("No signal") {
			t.Fatal("empty state not shown")
		}
	}
	if got[1] != got[0] || got[2] != got[0] {
		t.Fatalf("state sizes %v", got)
	}
	mediaPanics(t, "band over 1", func(c *ui.Context) { AudioSpectrum(c, []float32{2}, AudioSpectrumOptions{}) })
	mediaPanics(t, "failed without error", func(c *ui.Context) { AudioSpectrum(c, nil, AudioSpectrumOptions{Status: data.DataFailed}) })
}

func TestAudioSpectrumSnapshot(t *testing.T) {
	bands := make([]float32, 32)
	for i := range bands {
		bands[i] = float32(0.15 + 0.8*float64((i*37)%29)/29)
	}
	snapshot(t, "audio_spectrum", 420, 160, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(380).Children(func() { AudioSpectrum(c, bands, AudioSpectrumOptions{}) })
	})
}
