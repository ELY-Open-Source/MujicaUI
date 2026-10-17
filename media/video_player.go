package media

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// VideoState is what the caller's playback reports.
type VideoState struct {
	// Frame is the current picture; nil shows "No video".
	Frame              ui.ImageSource
	Position, Duration time.Duration
	Buffered           []MediaRange
	Playing            bool
	// Rate is the speed; zero means 1.
	Rate float64
	// Subtitle is the cue to show now; empty shows none.
	Subtitle                     string
	SubtitlesOn, Fullscreen, PiP bool
}

// VideoPlayerOptions are callbacks; nil optional ones hide their buttons.
type VideoPlayerOptions struct {
	OnPlay, OnPause func()
	OnSeek          func(time.Duration)
	OnRate          func(float64)
	OnSubtitles     func(on bool)
	OnFullscreen    func(on bool)
	OnPiP           func(on bool)
	Disabled        bool
	// Label names the video for assistive technology.
	Label string
}

// VideoBlack is the stage behind letterboxed frames in both themes.
var VideoBlack = ui.Hex("#0E0C0D")

// VideoPlayer draws the caller's frame and reports controls through callbacks.
func VideoPlayer(c *ui.Context, st VideoState, opts VideoPlayerOptions) ui.Element {
	if opts.OnPlay == nil || opts.OnPause == nil || opts.OnSeek == nil {
		panic("mujica: VideoPlayer needs OnPlay, OnPause and OnSeek")
	}
	if st.Duration < 0 || st.Rate < 0 {
		panic("mujica: VideoPlayer needs Duration and Rate >= 0")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	col := ui.Column(c).Radius(theme.CardRadius).Border(theme.BorderWidth, k.Border).Background(k.Surface).
		Clip().Role(ui.RoleGroup).Label(videoLabel(c, opts.Label)).Disabled(opts.Disabled)
	col.Children(func() {
		videoStage(c, k, st, opts)
		ui.Column(c).Padding(u, u*2, u*2, u*2).Children(func() {
			if st.Duration > 0 {
				pos := min(max(0, st.Position), st.Duration)
				if VideoScrubber(c, &pos, VideoScrubberOptions{Duration: st.Duration, Buffered: st.Buffered}).Changed() {
					opts.OnSeek(pos)
				}
			}
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
				mediaPlayToggle(c, st.Playing, opts.OnPlay, opts.OnPause)
				ui.Text(c, core.MediaClock(st.Position)+" / "+core.MediaClock(st.Duration)).FontSize(core.FontSize(c, theme.CaptionSize)).
					TextColor(k.TextMuted).FontFeatures("tnum").SingleLine().Grow(1)
				videoButtons(c, st, opts)
			})
		})
	})
	return col
}

func videoLabel(c *ui.Context, label string) string {
	if label != "" {
		return label
	}
	return core.Msg(c, "video_player.label")
}

// videoStage is the focusable picture with playback keys.
func videoStage(c *ui.Context, k theme.Tokens, st VideoState, opts VideoPlayerOptions) {
	stage := ui.Box(c).FillWidth().AspectRatio(16.0 / 9).Background(VideoBlack).Center().Focusable().
		Role(ui.RoleImage).Label(videoLabel(c, opts.Label))
	internal.Ring(stage, k.Focus, 0)
	if !stage.IsDisabled() {
		videoKeys(stage, st, opts)
		if stage.Clicked() {
			mediaToggle(st.Playing, opts.OnPlay, opts.OnPause)
		}
	}
	stage.Children(func() {
		if st.Frame == nil {
			ui.Text(c, core.Msg(c, "video_player.empty")).TextColor(ui.Hex("#BCAEAA"))
		} else {
			ui.Image(c, st.Frame).Fill()
		}
		if st.SubtitlesOn && st.Subtitle != "" {
			ui.Box(c).Absolute().Left(0).Right(0).Bottom(16).Center().PassThrough().Children(func() {
				ui.Text(c, st.Subtitle).TextColor(ui.Hex("#F2EBDD")).Background(ui.RGBA(0, 0, 0, 0.72)).
					TextAlign(ui.Center).MaxLines(3)
			})
		}
	})
}

func videoKeys(stage ui.Element, st VideoState, opts VideoPlayerOptions) {
	seek := func(d time.Duration) {
		to := min(max(0, st.Position+d), st.Duration)
		slog.Debug("mujica: video_player seek", "pos", to)
		opts.OnSeek(to)
	}
	switch {
	case stage.Shortcut(0, ui.KeySpace), stage.Shortcut(0, ui.KeyK):
		mediaToggle(st.Playing, opts.OnPlay, opts.OnPause)
	case stage.Shortcut(0, ui.KeyLeft):
		seek(-5 * time.Second)
	case stage.Shortcut(0, ui.KeyRight):
		seek(5 * time.Second)
	case opts.OnFullscreen != nil && stage.Shortcut(0, ui.KeyF):
		opts.OnFullscreen(!st.Fullscreen)
	case opts.OnSubtitles != nil && stage.Shortcut(0, ui.KeyC):
		opts.OnSubtitles(!st.SubtitlesOn)
	}
}

// videoButtons are rate, subtitles, picture-in-picture and full screen.
func videoButtons(c *ui.Context, st VideoState, opts VideoPlayerOptions) {
	if opts.OnRate != nil {
		rate := st.Rate
		if rate == 0 {
			rate = 1
		}
		if PlaybackSpeedControl(c, &rate, PlaybackSpeedControlOptions{}).Changed() {
			opts.OnRate(rate)
		}
	}
	toggle := func(icon *ui.SVG, key string, on bool, fn func(bool)) {
		if fn == nil {
			return
		}
		b := MediaToggleAction(c, icon, core.Msg(c, key), on)
		if b.Clicked() && !b.IsDisabled() {
			slog.Debug("mujica: video_player toggle", "control", key, "on", !on)
			fn(!on)
		}
	}
	toggle(core.MediaCaptions, "video_player.subtitles", st.SubtitlesOn, opts.OnSubtitles)
	toggle(core.MediaPiP, "video_player.pip", st.PiP, opts.OnPiP)
	toggle(core.MediaFullscreen, "video_player.fullscreen", st.Fullscreen, opts.OnFullscreen)
}

// MediaToggleAction is an iconAction that shows when it is on.
func MediaToggleAction(c *ui.Context, icon *ui.SVG, label string, on bool) ui.Element {
	k := core.Tokens(c)
	b := core.IconAction(c, icon, label).Role(ui.RoleToggleButton).Selected(on)
	if on {
		b.Background(k.Selection).TextColor(k.AccentText).BorderWidth(0, 0, 2, 0).BorderColor(k.AccentText)
	}
	return b
}
