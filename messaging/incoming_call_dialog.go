package messaging

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
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
	site, focus := core.OverlaySite(c)
	if opts.Disabled || site.IsDisabled() {
		*open = false
	}
	if !*open {
		focus.Closed(c, site)
		return
	}
	core.OverlayLayer(c, open)
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	title := fmt.Sprintf(core.Msg(c, "incoming_call_dialog.title"), call.Caller)
	shown, accepted, video, declined, backClicked := true, false, false, false, false
	panel := ui.DialogBase(c, &shown, func(back, panel *ui.Element) {
		backClicked = back.Clicked() // a click outside neither answers nor refuses
		p := core.OverlayEnter(c, panel)
		back.Background(overlay.OverlayScrim(k).Alpha(p)).Padding(16)
		overlay.DialogPanel(c, panel, 420, false, p)
		panel.Role(ui.RoleAlertDialog).Label(title).AlignItems(ui.Center)
		ui.Scroll(c).Shrink(1).FillWidth().Gap(u * 3).AlignItems(ui.Center).Children(func() {
			display.Avatar(c, call.Caller, display.AvatarOptions{Image: call.Image, Size: 72})
			ui.Text(c, call.Caller).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600).TextAlign(ui.Center)
			kind := core.Msg(c, "incoming_call_dialog.voice")
			if call.Video {
				kind = core.Msg(c, "incoming_call_dialog.video")
			}
			ui.Text(c, kind).TextColor(k.TextMuted)
		})
		ui.Row(c).Shrink(0).Gap(u*3).Wrap().Justify(ui.Center).Margin(u*2, 0, 0, 0).Children(func() {
			decline := input.Button(c, core.Msg(c, "incoming_call_dialog.decline"), input.ButtonOptions{Variant: input.Danger, Icon: msgHangUp, Label: core.Msg(c, "incoming_call_dialog.decline")})
			accept := input.Button(c, core.Msg(c, "incoming_call_dialog.accept"), input.ButtonOptions{Icon: msgPhone, Label: core.Msg(c, "incoming_call_dialog.accept")})
			if opts.FocusAccept {
				accept.AutoFocus()
			} else {
				decline.AutoFocus()
			}
			declined, accepted = decline.Clicked(), accept.Clicked()
			if call.Video && input.Button(c, core.Msg(c, "incoming_call_dialog.acceptVideo"), input.ButtonOptions{Variant: input.Secondary, Icon: core.MediaCamera, Label: core.Msg(c, "incoming_call_dialog.acceptVideo")}).Clicked() {
				accepted, video = true, true
			}
		})
	})
	focus.Shown(panel)
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
