package mujica_test

import (
	"image"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

// examplePhoto stands in for a decoded picture.
var examplePhoto = image.NewRGBA(image.Rect(0, 0, 320, 180))

func ExampleImageCropper() {
	crop := mujica.ImageCrop{} // empty: starts fitted to the image
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ImageCropper(c, examplePhoto, &crop, mujica.ImageCropperOptions{Aspect: 1}).Changed() {
			_ = crop.Rect // in pixels of the photo turned by crop.Rotation
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleImageAnnotator() {
	photo := ui.NewBitmap(examplePhoto)
	var notes []mujica.Annotation
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ImageAnnotator(c, photo, &notes, mujica.ImageAnnotatorOptions{}).Changed() {
			data, _ := mujica.ExportAnnotations(notes)
			_ = data // save the JSON
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleImageThumbnail() {
	photo := ui.NewBitmap(examplePhoto)
	selected := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ImageThumbnail(c, photo, mujica.ImageThumbnailOptions{Alt: "Grand hall", Aspect: 4.0 / 3, Selected: selected}).Changed() {
			selected = !selected
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLightbox() {
	images := []mujica.ViewerImage{{Source: ui.NewBitmap(examplePhoto), Alt: "The grand hall"}}
	open, index := false, 0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Button(c, "View photos", mujica.ButtonOptions{}).Clicked() {
			open = true
		}
		mujica.Lightbox(c, &open, &index, images)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVideoPlayer() {
	st := mujica.VideoState{Duration: 4 * time.Minute, Rate: 1}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.VideoPlayer(c, st, mujica.VideoPlayerOptions{
			OnPlay:  func() { st.Playing = true },  // start your decoder
			OnPause: func() { st.Playing = false }, // pause it
			OnSeek:  func(d time.Duration) { st.Position = d },
			OnRate:  func(r float64) { st.Rate = r },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVideoScrubber() {
	pos := 30 * time.Second
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.VideoScrubber(c, &pos, mujica.VideoScrubberOptions{
			Duration: 2 * time.Minute,
			Buffered: []mujica.MediaRange{{Start: 0, End: time.Minute}},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVideoThumbnailStrip() {
	frames := []mujica.StripFrame{{At: 0}, {At: 10 * time.Second}, {At: 20 * time.Second}}
	current := 0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.VideoThumbnailStrip(c, &current, frames, mujica.VideoThumbnailStripOptions{}).Changed() {
			_ = frames[current].At // seek the player there
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAudioPlayer() {
	st := mujica.AudioState{Title: "Overture", Duration: 3 * time.Minute}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.AudioPlayer(c, st, mujica.AudioPlayerOptions{
			OnPlay:  func() { st.Playing = true },
			OnPause: func() { st.Playing = false },
			OnSeek:  func(d time.Duration) { st.Position = d },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAudioWaveform() {
	samples := make([]float32, 48000) // from your decoder, -1 to 1
	pos, sel := time.Duration(0), mujica.MediaRange{}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.AudioWaveform(c, samples, mujica.AudioWaveformOptions{
			Duration: time.Second, Position: pos, Selection: &sel,
			OnSeek: func(d time.Duration) { pos = d },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAudioSpectrum() {
	bands := []float32{0.2, 0.6, 0.9, 0.4, 0.1} // set by your analyser, then c.Invalidate()
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.AudioSpectrum(c, bands, mujica.AudioSpectrumOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleVolumeControl() {
	volume, muted := 0.8, false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.VolumeControl(c, &volume, &muted, mujica.VolumeControlOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePlaybackSpeedControl() {
	rate := 1.0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.PlaybackSpeedControl(c, &rate, mujica.PlaybackSpeedControlOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMediaControls() {
	st := mujica.MediaControlsState{}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MediaControls(c, st, mujica.MediaControlsOptions{
			OnPrevious: func() {}, OnNext: func() {},
			OnPlay: func() { st.Playing = true }, OnPause: func() { st.Playing = false },
			OnLoop: func(m mujica.LoopMode) { st.Loop = m }, OnShuffle: func(on bool) { st.Shuffle = on },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePlaylist() {
	items := []mujica.PlaylistItem{{ID: "a", Title: "Overture"}, {ID: "b", Title: "Masquerade"}}
	current := "a"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if r := mujica.Playlist(c, &items, mujica.PlaylistOptions{Current: current}); r.Played != "" {
			current = r.Played
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSubtitleEditor() {
	cues := []mujica.SubtitleCue{{Start: 0, End: 2 * time.Second, Text: "The curtain rises."}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.SubtitleEditor(c, &cues, mujica.SubtitleEditorOptions{})
		if len(mujica.CheckSubtitles(cues)) > 0 {
			_ = cues // block saving until the issues are fixed
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleScreenRecorderControls() {
	st := mujica.RecorderState{}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ScreenRecorderControls(c, st, mujica.ScreenRecorderControlsOptions{
			OnStart:  func() { st.Status = mujica.RecorderRecording },
			OnPause:  func() { st.Status = mujica.RecorderPaused },
			OnResume: func() { st.Status = mujica.RecorderRecording },
			OnStop:   func() { st = mujica.RecorderState{} },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCameraPreview() {
	frame := mujica.CameraFrame{Image: examplePhoto, Seq: 1} // a new Seq for each captured frame
	mirror := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CameraPreview(c, frame, &mirror, mujica.CameraPreviewOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMicLevelMeter() {
	level := float32(0.4) // the input's peak, 0 to 1 of full scale
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MicLevelMeter(c, level, mujica.MicLevelMeterOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDeviceSelector() {
	devices := mujica.MediaDevices{
		Microphones: []mujica.MediaDevice{{ID: "mic-1", Name: "Built-in Microphone"}},
		Speakers:    []mujica.MediaDevice{{ID: "out-1", Name: "Headphones"}},
	}
	sel := mujica.DeviceSelection{Microphone: "mic-1"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DeviceSelector(c, &sel, devices, mujica.DeviceSelectorOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleImageCompare() {
	before, after := ui.NewBitmap(examplePhoto), ui.NewBitmap(examplePhoto)
	split := 0.5
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ImageCompare(c, before, after, &split, mujica.ImageCompareOptions{})
	}
	ui.Render(view, 400, 300, 1)
}
