package main

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/media"
	"github.com/ELY-Open-Source/MujicaUI/navigation"
)

// mediaClockState simulates the caller's playback: position advances while playing.
type mediaClockState struct {
	pos, dur time.Duration
	playing  bool
	rate     float64
	last     time.Time
}

// tick advances the position; frames are asked only while playing.
func (m *mediaClockState) tick(c *ui.Context) {
	now := c.Now()
	if m.playing {
		if !m.last.IsZero() {
			m.pos += time.Duration(float64(now.Sub(m.last)) * m.rate)
		}
		if m.pos >= m.dur {
			m.pos, m.playing = m.dur, false
		} else {
			c.AnimationFrame()
		}
	}
	m.last = now
}

func videoPlayerDemo(c *ui.Context, log func(string)) {
	type state struct {
		clock                    mediaClockState
		subs, full, pip, noFrame bool
		disabled                 bool
	}
	s := ui.Local(c.Root(), "345", func() state {
		return state{clock: mediaClockState{dur: 3 * time.Minute, rate: 1}, subs: true}
	})
	s.clock.tick(c)
	mediaSection(c, "Parameters", func() {
		ui.Row(c).Gap(16).Children(func() {
			mediaCheck(c, log, &s.noFrame, "No frame")
			mediaCheck(c, log, &s.disabled, "Disabled")
		})
	})
	cues := []string{"The curtain rises.", "Welcome to the masquerade.", "Ave Mujica."}
	var frame ui.ImageSource = mediaBitmaps[int(s.clock.pos/(20*time.Second))%3]
	if s.noFrame {
		frame = nil
	}
	st := media.VideoState{Frame: frame, Position: s.clock.pos, Duration: s.clock.dur, Playing: s.clock.playing,
		Rate: s.clock.rate, Buffered: []media.MediaRange{{Start: 0, End: s.clock.pos + 40*time.Second}},
		Subtitle: cues[int(s.clock.pos/(10*time.Second))%3], SubtitlesOn: s.subs, Fullscreen: s.full, PiP: s.pip}
	mediaSection(c, "Live", func() {
		ui.Column(c).MaxWidth(640).Children(func() {
			media.VideoPlayer(c, st, media.VideoPlayerOptions{
				OnPlay:       func() { s.clock.playing = true; log("play") },
				OnPause:      func() { s.clock.playing = false; log("pause") },
				OnSeek:       func(d time.Duration) { s.clock.pos = d; log("seek " + mediaClockLabel(d)) },
				OnRate:       func(r float64) { s.clock.rate = r; log(fmt.Sprintf("rate %v×", r)) },
				OnSubtitles:  func(on bool) { s.subs = on; log(fmt.Sprintf("subtitles %v", on)) },
				OnFullscreen: func(on bool) { s.full = on; log(fmt.Sprintf("full screen requested %v", on)) },
				OnPiP:        func(on bool) { s.pip = on; log(fmt.Sprintf("picture in picture requested %v", on)) },
				Disabled:     s.disabled,
			})
		})
	})
}

func videoScrubberDemo(c *ui.Context, log func(string)) {
	type state struct {
		pos               time.Duration
		preview, disabled bool
	}
	s := ui.Local(c.Root(), "346", func() state { return state{pos: 40 * time.Second, preview: true} })
	mediaSection(c, "Parameters", func() {
		ui.Row(c).Gap(16).Children(func() {
			mediaCheck(c, log, &s.preview, "Preview slot")
			mediaCheck(c, log, &s.disabled, "Disabled")
		})
	})
	opts := media.VideoScrubberOptions{Duration: 3 * time.Minute, Disabled: s.disabled,
		Buffered: []media.MediaRange{{Start: 0, End: 90 * time.Second}, {Start: 120 * time.Second, End: 140 * time.Second}}}
	if s.preview {
		opts.Preview = func(at time.Duration) {
			ui.Image(c, mediaBitmaps[int(at/(time.Minute))%3]).Fit(ui.Cover).Fill()
		}
	}
	mediaSection(c, "Live", func() {
		ui.Column(c).MaxWidth(560).Padding(110, 0, 0, 0).Children(func() {
			if media.VideoScrubber(c, &s.pos, opts).Changed() {
				log("seek " + mediaClockLabel(s.pos))
			}
		})
	})
}

func audioPlayerDemo(c *ui.Context, log func(string)) {
	type state struct {
		clock    mediaClockState
		disabled bool
	}
	s := ui.Local(c.Root(), "348", func() state { return state{clock: mediaClockState{dur: 4*time.Minute + 5*time.Second, rate: 1}} })
	s.clock.tick(c)
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.disabled, "Disabled") })
	mediaSection(c, "Live", func() {
		ui.Column(c).MaxWidth(560).Children(func() {
			media.AudioPlayer(c, media.AudioState{Title: "Ave Mujica — Overture", Position: s.clock.pos, Duration: s.clock.dur, Playing: s.clock.playing},
				media.AudioPlayerOptions{
					OnPlay:   func() { s.clock.playing = true; log("play") },
					OnPause:  func() { s.clock.playing = false; log("pause") },
					OnSeek:   func(d time.Duration) { s.clock.pos = d; log("seek " + mediaClockLabel(d)) },
					Disabled: s.disabled,
				})
		})
	})
}

// mediaSamples is a synthetic voice-like recording, made once.
var mediaSamples = func() []float32 {
	s := make([]float32, 120000)
	for i := range s {
		env := 0.25 + 0.7*math.Abs(math.Sin(float64(i)/9000))*math.Abs(math.Sin(float64(i)/2100))
		s[i] = float32(env * math.Sin(float64(i)/7))
	}
	return s
}()

func audioWaveformDemo(c *ui.Context, log func(string)) {
	type state struct {
		pos       time.Duration
		sel       media.MediaRange
		status    int
		selecting bool
	}
	s := ui.Local(c.Root(), "349", func() state { return state{pos: 20 * time.Second, selecting: true} })
	mediaSection(c, "Parameters", func() {
		ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
			before := s.status
			navigation.SegmentedControl(c, &s.status, []navigation.Segment{{Label: "Ready"}, {Label: "Loading"}, {Label: "Failed"}, {Label: "Empty"}}, navigation.SegmentedControlOptions{Label: "Status"})
			if s.status != before {
				log(fmt.Sprintf("status %d", s.status))
			}
			mediaCheck(c, log, &s.selecting, "Selectable")
		})
	})
	opts := media.AudioWaveformOptions{Duration: time.Minute, Position: s.pos, OnSeek: func(d time.Duration) { s.pos = d; log("seek " + mediaClockLabel(d)) }}
	if s.selecting {
		opts.Selection = &s.sel
	}
	samples := mediaSamples
	switch s.status {
	case 1:
		opts.Status, samples = data.DataLoading, nil
	case 2:
		opts.Status, opts.Error, samples = data.DataFailed, "The recording could not be decoded.", nil
	case 3:
		samples = nil
	}
	mediaSection(c, "Live", func() {
		if media.AudioWaveform(c, samples, opts).Changed() {
			log(fmt.Sprintf("selection %s – %s", mediaClockLabel(s.sel.Start), mediaClockLabel(s.sel.End)))
		}
	})
}

func audioSpectrumDemo(c *ui.Context, log func(string)) {
	type state struct{ live bool }
	s := ui.Local(c.Root(), "350", func() state { return state{live: true} })
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.live, "Live input") })
	freqs := []string{"32 Hz", "63 Hz", "125 Hz", "250 Hz", "500 Hz", "1 kHz", "2 kHz", "4 kHz", "8 kHz", "16 kHz"}
	bands := make([]float32, len(freqs))
	t := float64(c.Now().UnixMilli()) / 1000
	for i := range bands {
		bands[i] = float32(0.2 + 0.6*math.Abs(math.Sin(float64(i)*0.9+t*(1+float64(i)*0.2))))
	}
	if s.live {
		// The demo is its own analyser; the component only draws.
		c.AnimationFrame()
	}
	mediaSection(c, "Live", func() {
		media.AudioSpectrum(c, bands, media.AudioSpectrumOptions{BandLabel: func(i int) string { return freqs[i] }})
	})
}
