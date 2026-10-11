package media

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

// LoopMode is how playback repeats. The zero value repeats nothing.
type LoopMode int

const (
	LoopOff LoopMode = iota
	LoopAll
	LoopOne
)

// next cycles off, all, one.
func (m LoopMode) next() LoopMode { return (m + 1) % 3 }

// MediaControlsState is what the caller's playback reports.
type MediaControlsState struct {
	Playing bool
	Loop    LoopMode
	Shuffle bool
	// NoPrevious and NoNext disable those buttons at the ends.
	NoPrevious, NoNext bool
}

// MediaControlsOptions are the callbacks; every one is required.
type MediaControlsOptions struct {
	OnPrevious, OnPlay, OnPause, OnNext func()
	OnLoop                              func(LoopMode)
	OnShuffle                           func(bool)
	Disabled                            bool
}

// MediaControls are previous, play, next, loop and shuffle buttons.
func MediaControls(c *ui.Context, st MediaControlsState, opts MediaControlsOptions) *ui.Element {
	if opts.OnPrevious == nil || opts.OnPlay == nil || opts.OnPause == nil || opts.OnNext == nil || opts.OnLoop == nil || opts.OnShuffle == nil {
		panic("mujica: MediaControls needs every callback")
	}
	if st.Loop < LoopOff || st.Loop > LoopOne {
		panic(fmt.Sprintf("mujica: unknown LoopMode %d", st.Loop))
	}
	row := ui.Row(c).Gap(core.Density(c).Unit()).AlignItems(ui.Center).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).
		Label(core.Msg(c, "media_controls.label")).Disabled(opts.Disabled)
	row.Children(func() {
		shuffle := MediaToggleAction(c, core.MediaShuffle, core.Msg(c, "media_controls.shuffle"), st.Shuffle)
		if shuffle.Clicked() && !shuffle.IsDisabled() {
			slog.Debug("mujica: media_controls shuffle", "on", !st.Shuffle)
			opts.OnShuffle(!st.Shuffle)
		}
		prev := core.IconAction(c, core.MediaPrev, core.Msg(c, "media_controls.previous")).Disabled(st.NoPrevious)
		if prev.Clicked() && !prev.IsDisabled() {
			slog.Debug("mujica: media_controls previous")
			opts.OnPrevious()
		}
		mediaPlayToggle(c, st.Playing, opts.OnPlay, opts.OnPause)
		next := core.IconAction(c, core.MediaNext, core.Msg(c, "media_controls.next")).Disabled(st.NoNext)
		if next.Clicked() && !next.IsDisabled() {
			slog.Debug("mujica: media_controls next")
			opts.OnNext()
		}
		icon := core.MediaRepeat
		if st.Loop == LoopOne {
			icon = core.MediaRepeatOne
		}
		loop := MediaToggleAction(c, icon, core.Msg(c, []string{"media_controls.loopOff", "media_controls.loopAll", "media_controls.loopOne"}[st.Loop]), st.Loop != LoopOff)
		if loop.Clicked() && !loop.IsDisabled() {
			slog.Debug("mujica: media_controls loop", "mode", st.Loop.next())
			opts.OnLoop(st.Loop.next())
		}
	})
	return row
}
