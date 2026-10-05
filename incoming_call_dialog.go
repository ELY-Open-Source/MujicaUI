package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// IncomingCall is who calls; Caller is required.
type IncomingCall struct {
	Caller string
	Image  *ui.Bitmap
	// Video offers answering with the camera on.
	Video bool
}

// IncomingCallOptions configure an IncomingCallDialog.
type IncomingCallOptions struct {
	// FocusAccept focuses Accept instead of the default Decline.
	FocusAccept bool
	Disabled    bool
	// OnAccept answers; OnDecline refuses or handles Escape.
	OnAccept  func(video bool)
	OnDecline func()
}

// IncomingCallDialog presents call actions with configurable initial focus.
func IncomingCallDialog(c *ui.Context, open *bool, call IncomingCall, opts IncomingCallOptions) {
	if call.Caller == "" || opts.OnAccept == nil || opts.OnDecline == nil {
		panic("mujica: IncomingCallDialog needs Caller, OnAccept and OnDecline")
	}
	site, focus := overlaySite(c)
	if opts.Disabled || site.IsDisabled() {
		*open = false
	}
	if !*open {
		focus.closed(c, site)
		return
	}
	overlayLayer(c, open)
	k := Tokens(c)
	u := Density(c).Unit()
	title := fmt.Sprintf(msg(c, "incoming_call_dialog.title"), call.Caller)
	shown, accepted, video, declined, backClicked := true, false, false, false, false
	panel := ui.DialogBase(c, &shown, func(back, panel *ui.Element) {
		backClicked = back.Clicked() // a click outside neither answers nor refuses
		p := overlayEnter(c, panel)
		back.Background(overlayScrim(k).Alpha(p)).Padding(16)
		dialogPanel(c, panel, 420, false, p)
		panel.Role(ui.RoleAlertDialog).Label(title).AlignItems(ui.Center)
		ui.Scroll(c).Shrink(1).FillWidth().Gap(u * 3).AlignItems(ui.Center).Children(func() {
			Avatar(c, call.Caller, AvatarOptions{Image: call.Image, Size: 72})
			ui.Text(c, call.Caller).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).FontWeight(600).TextAlign(ui.Center)
			kind := msg(c, "incoming_call_dialog.voice")
			if call.Video {
				kind = msg(c, "incoming_call_dialog.video")
			}
			ui.Text(c, kind).TextColor(k.TextMuted)
		})
		ui.Row(c).Shrink(0).Gap(u*3).Wrap().Justify(ui.Center).Margin(u*2, 0, 0, 0).Children(func() {
			decline := Button(c, msg(c, "incoming_call_dialog.decline"), ButtonOptions{Variant: Danger, Icon: msgHangUp, Label: msg(c, "incoming_call_dialog.decline")})
			accept := Button(c, msg(c, "incoming_call_dialog.accept"), ButtonOptions{Icon: msgPhone, Label: msg(c, "incoming_call_dialog.accept")})
			if opts.FocusAccept {
				accept.AutoFocus()
			} else {
				decline.AutoFocus()
			}
			declined, accepted = decline.Clicked(), accept.Clicked()
			if call.Video && Button(c, msg(c, "incoming_call_dialog.acceptVideo"), ButtonOptions{Variant: Secondary, Icon: mediaCamera, Label: msg(c, "incoming_call_dialog.acceptVideo")}).Clicked() {
				accepted, video = true, true
			}
		})
	})
	focus.shown(panel)
	if !shown && !accepted && !backClicked {
		declined = true // Escape
	}
	switch {
	case accepted:
		*open = false
		slog.Debug("mujica: incoming_call accept", "video", video)
		opts.OnAccept(video)
	case declined:
		*open = false
		slog.Debug("mujica: incoming_call decline")
		opts.OnDecline()
	}
}
