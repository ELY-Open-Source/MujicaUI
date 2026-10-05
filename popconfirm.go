package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// PopconfirmOptions configure a Popconfirm. Title is required.
type PopconfirmOptions struct {
	Title       string
	Description string
	// ConfirmLabel is "Confirm" and CancelLabel "Cancel" when empty.
	ConfirmLabel string
	CancelLabel  string
	// Danger paints the confirm button as destructive.
	Danger bool
}

// PopconfirmResult reports the choice made in this frame.
type PopconfirmResult struct {
	Confirmed bool
	Cancelled bool
}

// Popconfirm asks for a light confirmation beside anchor while *open is
// true. Confirm takes the focus; Cancel, Escape and an outside click cancel.
// Use AlertDialog for high-risk actions.
func Popconfirm(c *ui.Context, anchor *ui.Element, open *bool, opts PopconfirmOptions) PopconfirmResult {
	if anchor == nil {
		panic("mujica: Popconfirm needs an anchor")
	}
	if opts.Title == "" {
		panic("mujica: Popconfirm needs a Title")
	}
	var res PopconfirmResult
	if !*open {
		return res
	}
	k := Tokens(c)
	u := Density(c).Unit()
	confirm, cancel := opts.ConfirmLabel, opts.CancelLabel
	if confirm == "" {
		confirm = msg(c, "popconfirm.confirm")
	}
	if cancel == "" {
		cancel = msg(c, "popconfirm.cancel")
	}
	variant := Primary
	if opts.Danger {
		variant = Danger
	}
	b := anchor.Bounds()
	shown := true
	ui.PopoverBase(c, anchor, &shown, func(panel *ui.Element) {
		overlayStyle(c, panel, overlayAnchored(c, panel, b))
		panel.Margin(6, 0, 0, 0).Padding(u * 3).Gap(u * 3).MaxWidth(320).Role(ui.RoleAlertDialog).
			Label(opts.Title).Description(opts.Description)
		ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Children(func() {
			color := k.Warning
			if opts.Danger {
				color = k.Danger
			}
			ui.Icon(c, icons.Must("circle-alert")).Size(18, 18).TextColor(color).Shrink(0).Margin(1, 0, 0, 0)
			ui.Column(c).Gap(u).Shrink(1).Children(func() {
				ui.Text(c, opts.Title).FontWeight(600)
				if opts.Description != "" {
					ui.Text(c, opts.Description).TextColor(k.TextMuted)
				}
			})
		})
		ui.Row(c).Gap(u * 2).Justify(ui.End).Children(func() {
			if Button(c, cancel, ButtonOptions{Variant: Secondary}).Clicked() {
				res.Cancelled = true
			}
			if Button(c, confirm, ButtonOptions{Variant: variant}).AutoFocus().Clicked() {
				res.Confirmed = true
			}
		})
	})
	if !shown {
		res.Cancelled = true
	}
	if res.Confirmed || res.Cancelled {
		slog.Debug("mujica: popconfirm close", "confirmed", res.Confirmed)
		*open = false
	}
	return res
}
