package media

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// RecorderStatus is where a recording is; zero is idle.
type RecorderStatus int

const (
	RecorderIdle RecorderStatus = iota
	RecorderRecording
	RecorderPaused
)

// RecorderState is what the caller's recorder reports.
type RecorderState struct {
	Status  RecorderStatus
	Elapsed time.Duration
}

// ScreenRecorderControlsOptions are the callbacks; every one is required.
type ScreenRecorderControlsOptions struct {
	OnStart, OnPause, OnResume, OnStop func()
	Disabled                           bool
}

// ScreenRecorderControls start, pause, resume and stop a recording.
func ScreenRecorderControls(c *ui.Context, st RecorderState, opts ScreenRecorderControlsOptions) *ui.Element {
	if opts.OnStart == nil || opts.OnPause == nil || opts.OnResume == nil || opts.OnStop == nil {
		panic("mujica: ScreenRecorderControls needs every callback")
	}
	if st.Status < RecorderIdle || st.Status > RecorderPaused || st.Elapsed < 0 {
		panic(fmt.Sprintf("mujica: bad RecorderState %+v", st))
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	row := ui.Row(c).Gap(u*2).Padding(u, u*2).AlignItems(ui.Center).Radius(theme.CardRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.Border).Role(ui.RoleToolbar).
		Label(core.Msg(c, "screen_recorder_controls.label")).Disabled(opts.Disabled)
	act := func(icon *ui.SVG, key string, variant input.ButtonVariant, fn func()) {
		if input.Button(c, core.Msg(c, key), input.ButtonOptions{Variant: variant, Icon: icon}).Clicked() && !row.IsDisabled() {
			slog.Debug("mujica: recorder " + key)
			fn()
		}
	}
	row.Children(func() {
		status := []string{"screen_recorder_controls.idle", "screen_recorder_controls.recording", "screen_recorder_controls.paused"}[st.Status]
		ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).Role(ui.RoleStatus).Label(core.Msg(c, status) + " " + core.MediaClock(st.Elapsed)).Children(func() {
			dot := k.TextMuted
			if st.Status == RecorderRecording {
				dot = k.Danger
			}
			ui.Box(c).Size(10, 10).Radius(5).Background(dot)
			ui.Text(c, core.Msg(c, status)).TextColor(k.TextMuted).SingleLine()
			ui.Text(c, core.MediaClock(st.Elapsed)).FontFeatures("tnum").SingleLine().MinWidth(48)
		})
		switch st.Status {
		case RecorderIdle:
			act(core.MediaRecord, "screen_recorder_controls.start", input.Danger, opts.OnStart)
		case RecorderRecording:
			act(core.MediaPause, "screen_recorder_controls.pause", input.Secondary, opts.OnPause)
			act(core.MediaStop, "screen_recorder_controls.stop", input.Danger, opts.OnStop)
		case RecorderPaused:
			act(core.MediaRecord, "screen_recorder_controls.resume", input.Secondary, opts.OnResume)
			act(core.MediaStop, "screen_recorder_controls.stop", input.Danger, opts.OnStop)
		}
	})
	return row
}
