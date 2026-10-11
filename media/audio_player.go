package media

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// AudioState is what the caller's playback reports.
type AudioState struct {
	Title    string
	Position time.Duration
	Duration time.Duration
	Playing  bool
}

// AudioPlayerOptions are the callbacks a player shell reports to.
type AudioPlayerOptions struct {
	OnPlay, OnPause func()
	OnSeek          func(time.Duration)
	Disabled        bool
}

// AudioPlayer is a play, pause and seek shell.
func AudioPlayer(c *ui.Context, st AudioState, opts AudioPlayerOptions) *ui.Element {
	if opts.OnPlay == nil || opts.OnPause == nil || opts.OnSeek == nil {
		panic("mujica: AudioPlayer needs OnPlay, OnPause and OnSeek")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	row := ui.Row(c).Gap(u*3).AlignItems(ui.Center).MinHeight(core.Density(c).ControlHeight()).Padding(u*2, u*3).
		Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).Role(ui.RoleGroup).
		Disabled(opts.Disabled)
	if st.Title != "" {
		row.Label(st.Title)
	}
	row.Children(func() {
		mediaPlayToggle(c, st.Playing, opts.OnPlay, opts.OnPause)
		ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
			if st.Title != "" {
				ui.Text(c, st.Title).SingleLine()
			}
			if st.Duration > 0 {
				pos := min(max(0, st.Position), st.Duration)
				if VideoScrubber(c, &pos, VideoScrubberOptions{Duration: st.Duration, Label: st.Title}).Changed() {
					opts.OnSeek(pos)
				}
			}
		})
		ui.Text(c, core.MediaClock(st.Position)+" / "+core.MediaClock(st.Duration)).FontSize(core.FontSize(c, theme.CaptionSize)).
			TextColor(k.TextMuted).FontFeatures("tnum").SingleLine()
	})
	// Space plays or pauses wherever the focus is inside.
	if !row.IsDisabled() && row.Shortcut(0, ui.KeySpace) {
		mediaToggle(st.Playing, opts.OnPlay, opts.OnPause)
	}
	return row
}

// mediaPlayToggle is the play/pause button; it calls the matching callback.
func mediaPlayToggle(c *ui.Context, playing bool, onPlay, onPause func()) *ui.Element {
	icon, label := core.MediaPlay, core.Msg(c, "media_controls.play")
	if playing {
		icon, label = core.MediaPause, core.Msg(c, "media_controls.pause")
	}
	b := input.IconButton(c, icon, label, input.IconButtonOptions{})
	if b.Clicked() && !b.IsDisabled() {
		mediaToggle(playing, onPlay, onPause)
	}
	return b
}

func mediaToggle(playing bool, onPlay, onPause func()) {
	if playing {
		slog.Debug("mujica: media pause")
		onPause()
		return
	}
	slog.Debug("mujica: media play")
	onPlay()
}
