package mujica

import (
	"github.com/egoist/mygo/ui"
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
		Use(c, Settings{})
		if Button(c, "Ring", ButtonOptions{}).Clicked() {
			*open = true
		}
		IncomingCallDialog(c, open, call, opts)
	}, 600, 500)
}
