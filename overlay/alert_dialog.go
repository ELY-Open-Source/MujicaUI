package overlay

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// AlertDialogOptions configure an AlertDialog. Title is required.
type AlertDialogOptions struct {
	Title       string
	Description string
	// ConfirmLabel is "Delete" when empty; CancelLabel is "Cancel".
	ConfirmLabel string
	CancelLabel  string
	// ConfirmText, when set, must be typed before Confirm unlocks.
	ConfirmText string
}

// AlertDialogResult reports the choice made in this frame.
type AlertDialogResult struct {
	Confirmed bool
	Canceled  bool
}

// alertDialogState is the typed confirmation, kept while the alert shows.
type alertDialogState struct{ typed string }

// AlertDialog confirms a destructive action; Cancel has the focus.
func AlertDialog(c *ui.Context, open *bool, opts AlertDialogOptions) AlertDialogResult {
	if opts.Title == "" {
		panic("mujica: AlertDialog needs a Title")
	}
	var res AlertDialogResult
	site, focus := core.OverlaySite(c)
	if !*open {
		focus.Closed(c, site)
		return res
	}
	core.OverlayLayer(c, open)
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	confirmLabel, cancelLabel := opts.ConfirmLabel, opts.CancelLabel
	if confirmLabel == "" {
		confirmLabel = core.Msg(c, "alert_dialog.confirm")
	}
	if cancelLabel == "" {
		cancelLabel = core.Msg(c, "alert_dialog.cancel")
	}
	shown, backClicked := true, false
	panel := ui.DialogBase(c, &shown, func(back, panel ui.Element) {
		backClicked = back.Clicked()
		st := ui.Local(panel, "mujica.alert", func() alertDialogState { return alertDialogState{} })
		p := core.OverlayEnter(c, panel)
		back.Background(OverlayScrim(k).Alpha(p)).Padding(16)
		DialogPanel(c, panel, 420, false, p)
		panel.Role(ui.RoleAlertDialog).Label(opts.Title).Description(opts.Description)
		ui.Row(c).Gap(u * 3).AlignItems(ui.Start).Shrink(0).Children(func() {
			ui.Icon(c, icons.Must("triangle-alert")).Size(22, 22).TextColor(k.Danger).Shrink(0).Margin(1, 0, 0, 0)
			ui.Text(c, opts.Title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600).Grow(1)
		})
		if opts.Description != "" {
			// A long description scrolls so the field and the actions stay in the window.
			ui.Scroll(c).Shrink(1).Padding(0, 0, 0, 22+u*3).Children(func() {
				ui.Text(c, opts.Description).TextColor(k.TextMuted)
			})
		}
		if opts.ConfirmText != "" {
			alertDialogTyped(c, st, opts.ConfirmText)
		}
		unlocked := opts.ConfirmText == "" || st.typed == opts.ConfirmText
		dialogActions(c, func() {
			if core.AutoFocus(c, input.Button(c, cancelLabel, input.ButtonOptions{Variant: input.Secondary})).Clicked() {
				res.Canceled = true
			}
			if input.Button(c, confirmLabel, input.ButtonOptions{Variant: input.Danger, Disabled: !unlocked}).Clicked() {
				res.Confirmed = true
			}
		})
	})
	focus.Shown(panel)
	if !shown && !backClicked {
		res.Canceled = true // Escape
	}
	if res.Confirmed {
		res.Canceled = false
	}
	if res.Confirmed || res.Canceled {
		slog.Debug("mujica: alert_dialog close", "confirmed", res.Confirmed)
		*open = false
	}
	return res
}

// alertDialogTyped builds the confirmation field; Enter in it never confirms, even mid-composition.
func alertDialogTyped(c *ui.Context, st *alertDialogState, want string) {
	k := core.Tokens(c)
	d := core.Density(c)
	prompt := fmt.Sprintf(core.Msg(c, "alert_dialog.type_to_unlock"), want)
	ui.Column(c).Gap(d.Unit()).Shrink(0).Children(func() {
		ui.Text(c, prompt).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		in := ui.TextInputBase(c, &st.typed).Label(prompt).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).
			Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder)
		internal.Ring(in, k.Focus, theme.ControlRadius)
	})
}
