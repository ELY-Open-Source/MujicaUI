package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/messaging"
)

func callControlsDemo(c *ui.Context, log func(string)) {
	type state struct {
		call     messaging.CallState
		disabled bool
		ended    bool
	}
	s := ui.Local(c.Root(), "409", func() state { return state{} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if input.Button(c, "Reconnect", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
		s.ended = false
		log("call joined")
	}
	messaging.CallControls(c, s.call, messaging.CallControlsOptions{
		Disabled: s.disabled || s.ended,
		OnMute:   func(on bool) { s.call.Muted = on; log(fmt.Sprintf("mute %v", on)) },
		OnCamera: func(off bool) { s.call.CameraOff = off; log(fmt.Sprintf("camera off %v", off)) },
		OnShare:  func(on bool) { s.call.Sharing = on; log(fmt.Sprintf("sharing %v", on)) },
		OnHangUp: func() { s.ended = true; log("call left") },
	})
	ui.Text(c, fmt.Sprintf("Muted: %v; camera off: %v; sharing: %v; ended: %v", s.call.Muted, s.call.CameraOff, s.call.Sharing, s.ended))
}

func videoCallGridDemo(c *ui.Context, log func(string)) {
	type state struct {
		speaking, muted, cameraOff, frame bool
		columns                           float64
	}
	s := ui.Local(c.Root(), "410", func() state { return state{speaking: true, frame: true, columns: 2} })
	mediaCheck(c, log, &s.speaking, "Speaking")
	mediaCheck(c, log, &s.muted, "Muted")
	mediaCheck(c, log, &s.cameraOff, "Camera off")
	mediaCheck(c, log, &s.frame, "Caller-provided frame")
	if input.NumberInput(c, &s.columns, input.NumberInputOptions{Min: 0, Max: 3, Step: 1, Label: "Columns (0 = auto)"}).Changed() {
		log(fmt.Sprintf("columns %.0f", s.columns))
	}
	people := []messaging.CallParticipant{
		{ID: "1", Name: "Sakiko", Speaking: s.speaking, Muted: s.muted, CameraOff: s.cameraOff},
		{ID: "2", Name: "Uika", Muted: true}, {ID: "3", Name: "Mutsumi", CameraOff: true}, {ID: "4", Name: "Umiri"},
	}
	if s.frame {
		people[0].Frame = mediaBitmaps[0]
		people[3].Frame = mediaBitmaps[1]
	}
	messaging.VideoCallGrid(c, people, messaging.VideoCallGridOptions{Columns: int(s.columns)})
}

func incomingCallDialogDemo(c *ui.Context, log func(string)) {
	type state struct {
		open, video, focusAccept, disabled bool
		result                             string
	}
	s := ui.Local(c.Root(), "411", func() state { return state{video: true} })
	mediaCheck(c, log, &s.video, "Video call")
	mediaCheck(c, log, &s.focusAccept, "Focus Accept")
	mediaCheck(c, log, &s.disabled, "Disabled")
	if input.Button(c, "Incoming call", input.ButtonOptions{Disabled: s.disabled}).Clicked() {
		s.open = true
		log("incoming call")
	}
	messaging.IncomingCallDialog(c, &s.open, messaging.IncomingCall{Caller: "Sakiko Togawa", Video: s.video}, messaging.IncomingCallOptions{
		Disabled: s.disabled, FocusAccept: s.focusAccept,
		OnAccept:  func(video bool) { s.result = fmt.Sprintf("Accepted, video: %v", video); log(s.result) },
		OnDecline: func() { s.result = "Declined"; log(s.result) },
	})
	ui.Text(c, s.result)
}
