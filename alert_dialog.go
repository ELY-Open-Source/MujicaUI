package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
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
	Cancelled bool
}

// alertDialogState is the typed confirmation, kept while the alert shows.
type alertDialogState struct{ typed string }

// AlertDialog asks to confirm a destructive action while *open is true.
// Cancel takes the focus, the backdrop ignores clicks, Escape cancels,
// and Enter never confirms before ConfirmText is typed.
func AlertDialog(c *ui.Context, open *bool, opts AlertDialogOptions) AlertDialogResult {
	if opts.Title == "" {
		panic("mujica: AlertDialog needs a Title")
	}
	var res AlertDialogResult
	if !*open {
		return res
	}
	k := Tokens(c)
	u := Density(c).Unit()
	confirmLabel, cancelLabel := opts.ConfirmLabel, opts.CancelLabel
	if confirmLabel == "" {
		confirmLabel = msg(c, "alert_dialog.confirm")
	}
	if cancelLabel == "" {
		cancelLabel = msg(c, "alert_dialog.cancel")
	}
	shown, backClicked := true, false
	ui.DialogBase(c, &shown, func(back, panel *ui.Element) {
		backClicked = back.Clicked()
		st := ui.Local(panel, "mujica.alert", func() alertDialogState { return alertDialogState{} })
		p := overlayEnter(c, panel)
		back.Background(overlayScrim(k).Alpha(p)).Padding(16)
		dialogPanel(c, panel, 420, false, p)
		panel.Role(ui.RoleAlertDialog).Label(opts.Title).Description(opts.Description)
		ui.Row(c).Gap(u * 3).AlignItems(ui.Start).Shrink(0).Children(func() {
			ui.Icon(c, icons.Must("triangle-alert")).Size(22, 22).TextColor(k.Danger).Shrink(0).Margin(1, 0, 0, 0)
			ui.Column(c).Gap(u * 2).Grow(1).Children(func() {
				ui.Text(c, opts.Title).Font(theme.SerifFont).FontSize(theme.H3Size).FontWeight(600)
				if opts.Description != "" {
					ui.Text(c, opts.Description).TextColor(k.TextMuted)
				}
			})
		})
		submitted := opts.ConfirmText != "" && alertDialogTyped(c, st, opts.ConfirmText)
		unlocked := opts.ConfirmText == "" || st.typed == opts.ConfirmText
		res.Confirmed = submitted && unlocked
		dialogActions(c, func() {
			if Button(c, cancelLabel, ButtonOptions{Variant: Secondary}).AutoFocus().Clicked() {
				res.Cancelled = true
			}
			if Button(c, confirmLabel, ButtonOptions{Variant: Danger, Disabled: !unlocked}).Clicked() {
				res.Confirmed = true
			}
		})
	})
	if !shown && !backClicked {
		res.Cancelled = true // Escape
	}
	if res.Confirmed {
		res.Cancelled = false
	}
	if res.Confirmed || res.Cancelled {
		slog.Debug("mujica: alert_dialog close", "confirmed", res.Confirmed)
		*open = false
	}
	return res
}

// alertDialogTyped builds the confirmation field and reports Enter in it.
func alertDialogTyped(c *ui.Context, st *alertDialogState, want string) bool {
	k := Tokens(c)
	d := Density(c)
	prompt := fmt.Sprintf(msg(c, "alert_dialog.type_to_unlock"), want)
	submitted := false
	ui.Column(c).Gap(d.Unit()).Shrink(0).Children(func() {
		ui.Text(c, prompt).FontSize(theme.CaptionSize).TextColor(k.TextMuted)
		in := ui.TextInputBase(c, &st.typed).Label(prompt).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).
			Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder)
		internal.Ring(in, k.Focus, theme.ControlRadius)
		submitted = in.Submitted()
	})
	return submitted
}
