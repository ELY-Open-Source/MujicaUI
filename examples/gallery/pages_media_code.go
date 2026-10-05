package main

// Media page snippets, each its compiled Example's body.
const (
	imageCropperCode = `crop := media.ImageCrop{} // empty: starts fitted to the image
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if media.ImageCropper(c, examplePhoto, &crop, media.ImageCropperOptions{Aspect: 1}).Changed() {
		_ = crop.Rect // in pixels of the photo turned by crop.Rotation
	}
}
ui.Render(view, 400, 300, 1)`
	imageAnnotatorCode = `photo := ui.NewBitmap(examplePhoto)
var notes []media.Annotation
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if media.ImageAnnotator(c, photo, &notes, media.ImageAnnotatorOptions{}).Changed() {
		data, _ := media.ExportAnnotations(notes)
		_ = data // save the JSON
	}
}
ui.Render(view, 400, 300, 1)`
	imageThumbnailCode = `photo := ui.NewBitmap(examplePhoto)
selected := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if media.ImageThumbnail(c, photo, media.ImageThumbnailOptions{Alt: "Grand hall", Aspect: 4.0 / 3, Selected: selected}).Changed() {
		selected = !selected
	}
}
ui.Render(view, 400, 300, 1)`
	lightboxCode = `images := []data.ViewerImage{{Source: ui.NewBitmap(examplePhoto), Alt: "The grand hall"}}
open, index := false, 0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.Button(c, "View photos", input.ButtonOptions{}).Clicked() {
		open = true
	}
	media.Lightbox(c, &open, &index, images)
}
ui.Render(view, 400, 300, 1)`
	videoPlayerCode = `st := media.VideoState{Duration: 4 * time.Minute, Rate: 1}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.VideoPlayer(c, st, media.VideoPlayerOptions{
		OnPlay:  func() { st.Playing = true },  // start your decoder
		OnPause: func() { st.Playing = false }, // pause it
		OnSeek:  func(d time.Duration) { st.Position = d },
		OnRate:  func(r float64) { st.Rate = r },
	})
}
ui.Render(view, 400, 300, 1)`
	videoScrubberCode = `pos := 30 * time.Second
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.VideoScrubber(c, &pos, media.VideoScrubberOptions{
		Duration: 2 * time.Minute,
		Buffered: []media.MediaRange{{Start: 0, End: time.Minute}},
	})
}
ui.Render(view, 400, 300, 1)`
	videoStripCode = `frames := []media.StripFrame{{At: 0}, {At: 10 * time.Second}, {At: 20 * time.Second}}
current := 0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if media.VideoThumbnailStrip(c, &current, frames, media.VideoThumbnailStripOptions{}).Changed() {
		_ = frames[current].At // seek the player there
	}
}
ui.Render(view, 400, 300, 1)`
	audioPlayerCode = `st := media.AudioState{Title: "Overture", Duration: 3 * time.Minute}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.AudioPlayer(c, st, media.AudioPlayerOptions{
		OnPlay:  func() { st.Playing = true },
		OnPause: func() { st.Playing = false },
		OnSeek:  func(d time.Duration) { st.Position = d },
	})
}
ui.Render(view, 400, 300, 1)`
	audioWaveformCode = `samples := make([]float32, 48000) // from your decoder, -1 to 1
pos, sel := time.Duration(0), media.MediaRange{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.AudioWaveform(c, samples, media.AudioWaveformOptions{
		Duration: time.Second, Position: pos, Selection: &sel,
		OnSeek: func(d time.Duration) { pos = d },
	})
}
ui.Render(view, 400, 300, 1)`
	audioSpectrumCode = `bands := []float32{0.2, 0.6, 0.9, 0.4, 0.1} // set by your analyser, then c.Invalidate()
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.AudioSpectrum(c, bands, media.AudioSpectrumOptions{})
}
ui.Render(view, 400, 300, 1)`
	volumeCode = `volume, muted := 0.8, false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.VolumeControl(c, &volume, &muted, media.VolumeControlOptions{})
}
ui.Render(view, 400, 300, 1)`
	speedCode = `rate := 1.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.PlaybackSpeedControl(c, &rate, media.PlaybackSpeedControlOptions{})
}
ui.Render(view, 400, 300, 1)`
	mediaControlsCode = `st := media.MediaControlsState{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.MediaControls(c, st, media.MediaControlsOptions{
		OnPrevious: func() {}, OnNext: func() {},
		OnPlay: func() { st.Playing = true }, OnPause: func() { st.Playing = false },
		OnLoop: func(m media.LoopMode) { st.Loop = m }, OnShuffle: func(on bool) { st.Shuffle = on },
	})
}
ui.Render(view, 400, 300, 1)`
	playlistCode = `items := []media.PlaylistItem{{ID: "a", Title: "Overture"}, {ID: "b", Title: "Masquerade"}}
current := "a"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if r := media.Playlist(c, &items, media.PlaylistOptions{Current: current}); r.Played != "" {
		current = r.Played
	}
}
ui.Render(view, 400, 300, 1)`
	subtitleCode = `cues := []media.SubtitleCue{{Start: 0, End: 2 * time.Second, Text: "The curtain rises."}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.SubtitleEditor(c, &cues, media.SubtitleEditorOptions{})
	if len(media.CheckSubtitles(cues)) > 0 {
		_ = cues // block saving until the issues are fixed
	}
}
ui.Render(view, 400, 300, 1)`
	recorderCode = `st := media.RecorderState{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.ScreenRecorderControls(c, st, media.ScreenRecorderControlsOptions{
		OnStart:  func() { st.Status = media.RecorderRecording },
		OnPause:  func() { st.Status = media.RecorderPaused },
		OnResume: func() { st.Status = media.RecorderRecording },
		OnStop:   func() { st = media.RecorderState{} },
	})
}
ui.Render(view, 400, 300, 1)`
	cameraCode = `frame := media.CameraFrame{Image: examplePhoto, Seq: 1} // a new Seq for each captured frame
mirror := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.CameraPreview(c, frame, &mirror, media.CameraPreviewOptions{})
}
ui.Render(view, 400, 300, 1)`
	micCode = `level := float32(0.4) // the input's peak, 0 to 1 of full scale
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.MicLevelMeter(c, level, media.MicLevelMeterOptions{})
}
ui.Render(view, 400, 300, 1)`
	deviceCode = `devices := media.MediaDevices{
	Microphones: []media.MediaDevice{{ID: "mic-1", Name: "Built-in Microphone"}},
	Speakers:    []media.MediaDevice{{ID: "out-1", Name: "Headphones"}},
}
sel := media.DeviceSelection{Microphone: "mic-1"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.DeviceSelector(c, &sel, devices, media.DeviceSelectorOptions{})
}
ui.Render(view, 400, 300, 1)`
	compareCode = `before, after := ui.NewBitmap(examplePhoto), ui.NewBitmap(examplePhoto)
split := 0.5
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.ImageCompare(c, before, after, &split, media.ImageCompareOptions{})
}
ui.Render(view, 400, 300, 1)`
)
