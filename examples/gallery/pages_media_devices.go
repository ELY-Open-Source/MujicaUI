package main

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/media"
)

func volumeDemo(c *ui.Context, log func(string)) {
	type state struct {
		volume          float64
		muted, disabled bool
	}
	s := ui.Local(c.Root(), "351", func() state { return state{volume: 0.7} })
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.disabled, "Disabled") })
	mediaSection(c, "Live", func() {
		if media.VolumeControl(c, &s.volume, &s.muted, media.VolumeControlOptions{Disabled: s.disabled}).Changed() {
			log(fmt.Sprintf("volume %.0f%%, muted %v", s.volume*100, s.muted))
		}
	})
}

func speedDemo(c *ui.Context, log func(string)) {
	type state struct {
		rate     float64
		disabled bool
	}
	s := ui.Local(c.Root(), "352", func() state { return state{rate: 1} })
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.disabled, "Disabled") })
	mediaSection(c, "Live", func() {
		if media.PlaybackSpeedControl(c, &s.rate, media.PlaybackSpeedControlOptions{Disabled: s.disabled}).Changed() {
			log(fmt.Sprintf("rate %v×", s.rate))
		}
	})
}

func mediaControlsDemo(c *ui.Context, log func(string)) {
	type state struct {
		st       media.MediaControlsState
		track    int
		disabled bool
	}
	s := ui.Local(c.Root(), "353", func() state { return state{} })
	tracks := []string{"Ave Mujica", "Imprisoned XII", "KiLLKiSS", "Symbol I: △"}
	s.st.NoPrevious, s.st.NoNext = s.track == 0, s.track == len(tracks)-1
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.disabled, "Disabled") })
	mediaSection(c, "Live", func() {
		ui.Text(c, fmt.Sprintf("Now: %s (%d of %d)", tracks[s.track], s.track+1, len(tracks)))
		media.MediaControls(c, s.st, media.MediaControlsOptions{
			OnPrevious: func() { s.track--; log("previous") },
			OnNext:     func() { s.track++; log("next") },
			OnPlay:     func() { s.st.Playing = true; log("play") },
			OnPause:    func() { s.st.Playing = false; log("pause") },
			OnLoop:     func(m media.LoopMode) { s.st.Loop = m; log(fmt.Sprintf("loop mode %d", m)) },
			OnShuffle:  func(on bool) { s.st.Shuffle = on; log(fmt.Sprintf("shuffle %v", on)) },
			Disabled:   s.disabled,
		})
	})
}

func playlistDemo(c *ui.Context, log func(string)) {
	type state struct {
		items    []media.PlaylistItem
		current  string
		disabled bool
	}
	s := ui.Local(c.Root(), "354", func() state {
		titles := []string{"Ave Mujica", "Imprisoned XII", "KiLLKiSS", "Symbol I: △", "Choir 'S' Choir", "Georgette Me, Georgette You"}
		items := make([]media.PlaylistItem, len(titles))
		for i, t := range titles {
			items[i] = media.PlaylistItem{ID: fmt.Sprint("t", i), Title: t, Subtitle: fmt.Sprintf("Act %d", i+1), Duration: time.Duration(200+i*17) * time.Second}
		}
		return state{items: items, current: "t1"}
	})
	mediaSection(c, "Parameters", func() {
		ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
			mediaCheck(c, log, &s.disabled, "Disabled")
			if input.Button(c, "Add track", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				id := fmt.Sprint("t", len(s.items)+int(time.Now().UnixNano()%1000000))
				s.items = append(s.items, media.PlaylistItem{ID: id, Title: "Encore " + id, Duration: 3 * time.Minute})
				log("added " + id)
			}
		})
	})
	mediaSection(c, "Live", func() {
		ui.Column(c).MaxWidth(520).Children(func() {
			r := media.Playlist(c, &s.items, media.PlaylistOptions{Current: s.current, Disabled: s.disabled})
			switch {
			case r.Played != "":
				s.current = r.Played
				log("play " + r.Played)
			case r.Removed != "":
				log("removed " + r.Removed)
			case r.Moved:
				log("reordered")
			}
		})
	})
}

func subtitleDemo(c *ui.Context, log func(string)) {
	type state struct {
		cues     []media.SubtitleCue
		disabled bool
	}
	s := ui.Local(c.Root(), "355", func() state {
		return state{cues: []media.SubtitleCue{
			{Start: 0, End: 2500 * time.Millisecond, Text: "The curtain rises."},
			{Start: 2 * time.Second, End: 5 * time.Second, Text: "Welcome to the masquerade."},
			{Start: 6 * time.Second, End: 9 * time.Second, Text: "Ave Mujica."},
		}}
	})
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.disabled, "Disabled") })
	mediaSection(c, "Live", func() {
		ui.Column(c).MaxWidth(720).Children(func() {
			if media.SubtitleEditor(c, &s.cues, media.SubtitleEditorOptions{Height: 260, Disabled: s.disabled}).Changed() {
				log(fmt.Sprintf("%d cues, %d issues", len(s.cues), len(media.CheckSubtitles(s.cues))))
			}
		})
	})
}

func recorderDemo(c *ui.Context, log func(string)) {
	type state struct {
		st          media.RecorderState
		since       time.Time
		disabled    bool
		accumulated time.Duration
	}
	s := ui.Local(c.Root(), "356", func() state { return state{} })
	if s.st.Status == media.RecorderRecording {
		s.st.Elapsed = s.accumulated + c.Now().Sub(s.since)
		// The clock changes each second; ask for that frame only.
		c.After(time.Second - s.st.Elapsed%time.Second)
	}
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.disabled, "Disabled") })
	mediaSection(c, "Live", func() {
		media.ScreenRecorderControls(c, s.st, media.ScreenRecorderControlsOptions{
			OnStart: func() {
				s.st, s.since, s.accumulated = media.RecorderState{Status: media.RecorderRecording}, c.Now(), 0
				log("start")
			},
			OnPause:  func() { s.st.Status, s.accumulated = media.RecorderPaused, s.st.Elapsed; log("pause") },
			OnResume: func() { s.st.Status, s.since = media.RecorderRecording, c.Now(); log("resume") },
			OnStop:   func() { log("stop at " + mediaClockLabel(s.st.Elapsed)); s.st = media.RecorderState{} },
			Disabled: s.disabled,
		})
	})
}

func cameraDemo(c *ui.Context, log func(string)) {
	type state struct {
		mirror, signal, disabled bool
	}
	s := ui.Local(c.Root(), "357", func() state { return state{mirror: true, signal: true} })
	mediaSection(c, "Parameters", func() {
		ui.Row(c).Gap(16).Children(func() {
			mediaCheck(c, log, &s.signal, "Signal")
			mediaCheck(c, log, &s.disabled, "Disabled")
		})
	})
	frame := media.CameraFrame{}
	if s.signal {
		frame = media.CameraFrame{Image: mediaPhotos[2], Seq: 1}
	}
	mediaSection(c, "Live", func() {
		ui.Column(c).MaxWidth(480).Children(func() {
			if media.CameraPreview(c, frame, &s.mirror, media.CameraPreviewOptions{Disabled: s.disabled}).Changed() {
				log(fmt.Sprintf("mirror %v", s.mirror))
			}
		})
	})
}

func micDemo(c *ui.Context, log func(string)) {
	type state struct{ live bool }
	s := ui.Local(c.Root(), "358", func() state { return state{live: true} })
	mediaSection(c, "Parameters", func() { mediaCheck(c, log, &s.live, "Live input") })
	level := float32(0.3)
	if s.live {
		t := float64(c.Now().UnixMilli()) / 1000
		level = float32(math.Min(1, 0.15+0.8*math.Abs(math.Sin(t*1.3)*math.Sin(t*0.37))))
		c.AnimationFrame()
	}
	mediaSection(c, "Live", func() {
		ui.Column(c).MaxWidth(420).Children(func() { media.MicLevelMeter(c, level, media.MicLevelMeterOptions{}) })
	})
}

func deviceDemo(c *ui.Context, log func(string)) {
	type state struct {
		sel                media.DeviceSelection
		noCamera, disabled bool
	}
	s := ui.Local(c.Root(), "359", func() state { return state{sel: media.DeviceSelection{Microphone: "m1", Speaker: "s2"}} })
	devices := media.MediaDevices{
		Microphones: []media.MediaDevice{{ID: "m1", Name: "Built-in Microphone"}, {ID: "m2", Name: "Studio Condenser"}},
		Cameras:     []media.MediaDevice{{ID: "c1", Name: "FaceTime HD Camera"}},
		Speakers:    []media.MediaDevice{{ID: "s1", Name: "Built-in Speakers"}, {ID: "s2", Name: "Headphones"}},
	}
	mediaSection(c, "Parameters", func() {
		ui.Row(c).Gap(16).Children(func() {
			before := s.noCamera
			mediaCheck(c, log, &s.noCamera, "No cameras")
			if s.noCamera != before {
				s.sel.Camera = ""
			}
			mediaCheck(c, log, &s.disabled, "Disabled")
		})
	})
	if s.noCamera {
		devices.Cameras = []media.MediaDevice{}
	}
	mediaSection(c, "Live", func() {
		ui.Column(c).MaxWidth(480).Children(func() {
			if media.DeviceSelector(c, &s.sel, devices, media.DeviceSelectorOptions{Disabled: s.disabled}).Changed() {
				log(fmt.Sprintf("microphone %q, camera %q, speaker %q", s.sel.Microphone, s.sel.Camera, s.sel.Speaker))
			}
		})
	})
}
