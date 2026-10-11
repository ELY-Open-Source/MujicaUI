package messaging

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type callCalls struct {
	mute, camera, share []bool
	hang                int
}

func (k *callCalls) opts() CallControlsOptions {
	return CallControlsOptions{OnMute: func(on bool) { k.mute = append(k.mute, on) }, OnCamera: func(on bool) { k.camera = append(k.camera, on) },
		OnShare: func(on bool) { k.share = append(k.share, on) }, OnHangUp: func() { k.hang++ }}
}

func incomingTester(open *bool, opts IncomingCallOptions, call IncomingCall) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Ring", input.ButtonOptions{}).Clicked() {
			*open = true
		}
		IncomingCallDialog(c, open, call, opts)
	}, 600, 500)
}
