package messaging

import (
	"log/slog"
	"runtime"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/media"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CallState is what the caller's call reports.
type CallState struct {
	Muted, CameraOff, Sharing bool
}

// CallControlsOptions are the callbacks; every one is required.
type CallControlsOptions struct {
	OnMute, OnCamera, OnShare func(on bool)
	OnHangUp                  func()
	Disabled                  bool
}

type callKey struct {
	key    ui.Key
	letter string
	shift  bool
}

var (
	callMute   = callKey{ui.KeyD, "D", false}
	callCamera = callKey{ui.KeyE, "E", false}
	callShare  = callKey{ui.KeyS, "S", true}
	callHangUp = callKey{ui.KeyH, "H", true}
)

// pressed reports the shortcut anywhere in the window.
func (k callKey) pressed(c *ui.Context) bool {
	mods := ui.Cmd
	if k.shift {
		mods |= ui.Shift
	}
	return c.Shortcut(mods, k.key)
}

// String formats the platform’s shortcut label.
func (k callKey) String() string {
	mac := runtime.GOOS == "darwin"
	keys := []string{"Mod"}
	if k.shift {
		keys = append(keys, "Shift")
	}
	out := ""
	for i, name := range append(keys, k.letter) {
		face, _ := display.KbdKey(name, mac)
		if i > 0 && !mac {
			out += "+"
		}
		out += face
	}
	return out
}

// CallControls provides call actions and window shortcuts.
func CallControls(c *ui.Context, st CallState, opts CallControlsOptions) *ui.Element {
	if opts.OnMute == nil || opts.OnCamera == nil || opts.OnShare == nil || opts.OnHangUp == nil {
		panic("mujica: CallControls needs every callback")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	row := ui.Row(c).Gap(u*2).Padding(u*2, u*3).AlignItems(ui.Center).Radius(theme.CardRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.Border).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).
		Label(core.Msg(c, "call_controls.label")).Disabled(opts.Disabled)
	live := !row.IsDisabled()
	report := func(name string, on bool, fn func(bool)) {
		slog.Debug("mujica: call_controls "+name, "on", on)
		fn(on)
	}
	row.Children(func() {
		toggle := func(off, on *ui.SVG, offKey, onKey string, state bool, key callKey, name string, fn func(bool)) {
			icon, label := off, core.Msg(c, offKey)
			if state {
				icon, label = on, core.Msg(c, onKey)
			}
			b := media.MediaToggleAction(c, icon, label, state).Size(40, 40).Tooltip(label + " (" + key.String() + ")")
			if live && b.Clicked() {
				report(name, !state, fn)
			}
		}
		toggle(core.MediaMic, msgMicOff, "call_controls.mute", "call_controls.unmute", st.Muted, callMute, "mute", opts.OnMute)
		toggle(core.MediaCamera, msgVideoOff, "call_controls.cameraOff", "call_controls.cameraOn", st.CameraOff, callCamera, "camera", opts.OnCamera)
		toggle(msgShare, msgShare, "call_controls.share", "call_controls.stopShare", st.Sharing, callShare, "share", opts.OnShare)
		hang := input.Button(c, core.Msg(c, "call_controls.hangUp"), input.ButtonOptions{Variant: input.Danger, Icon: msgHangUp}).
			Tooltip(core.Msg(c, "call_controls.hangUp") + " (" + callHangUp.String() + ")")
		if live && hang.Clicked() {
			slog.Debug("mujica: call_controls hang up")
			opts.OnHangUp()
		}
	})
	if live {
		switch {
		case callMute.pressed(c):
			report("mute", !st.Muted, opts.OnMute)
		case callCamera.pressed(c):
			report("camera", !st.CameraOff, opts.OnCamera)
		case callShare.pressed(c):
			report("share", !st.Sharing, opts.OnShare)
		case callHangUp.pressed(c):
			slog.Debug("mujica: call_controls hang up")
			opts.OnHangUp()
		}
	}
	return row
}
