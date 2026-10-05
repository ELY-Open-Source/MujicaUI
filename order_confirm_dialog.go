package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// OrderSummaryLine is one labeled value of an order summary.
type OrderSummaryLine struct {
	Label, Value string
}

// OrderConfirmDialogOptions configure an OrderConfirmDialog. Lines are required.
type OrderConfirmDialogOptions struct {
	// Title is "Confirm order" when empty.
	Title string
	Lines []OrderSummaryLine
	// Risk is a warning shown above the actions.
	Risk string
	// Acknowledge, when set, must be ticked before confirming.
	Acknowledge string
	// ConfirmLabel is "Place order" when empty.
	ConfirmLabel string
}

type orderConfirmState struct{ ack bool }

// OrderConfirmDialog summarizes an order; Cancel has the focus.
func OrderConfirmDialog(c *ui.Context, open *bool, opts OrderConfirmDialogOptions) AlertDialogResult {
	if len(opts.Lines) == 0 {
		panic("mujica: OrderConfirmDialog needs summary Lines")
	}
	var res AlertDialogResult
	k := Tokens(c)
	u := Density(c).Unit()
	title, confirm := opts.Title, opts.ConfirmLabel
	if title == "" {
		title = msg(c, "order_confirm_dialog.title")
	}
	if confirm == "" {
		confirm = msg(c, "order_confirm_dialog.confirm")
	}
	was := *open
	var st *orderConfirmState
	Dialog(c, open, DialogOptions{Title: title, KeepOnBackdrop: true, HideClose: true, Actions: func() {
		if Button(c, msg(c, "order_confirm_dialog.cancel"), ButtonOptions{Variant: Secondary}).AutoFocus().Clicked() {
			res.Canceled = true
		}
		locked := opts.Acknowledge != "" && !st.ack
		if Button(c, confirm, ButtonOptions{Disabled: locked}).Clicked() && !locked {
			res.Confirmed = true
		}
	}}, func() {
		list := ui.Column(c)
		// The tick lives with the open dialog only.
		st = ui.Local(list, "ack", func() orderConfirmState { return orderConfirmState{} })
		list.Children(func() {
			for _, l := range opts.Lines {
				ui.Row(c).Gap(u*4).MinHeight(Density(c).ControlHeight()).AlignItems(ui.Center).
					BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).Children(func() {
					ui.Text(c, l.Label).TextColor(k.TextMuted).Grow(1)
					ui.Text(c, l.Value).FontFeatures("tnum").FontWeight(600)
				})
			}
		})
		if opts.Risk != "" {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Padding(u * 3).Radius(theme.ControlRadius).Background(k.Warning.Alpha(0.12)).Children(func() {
				ui.Icon(c, icons.Must("triangle-alert")).Size(16, 16).TextColor(k.Warning).Shrink(0)
				ui.Text(c, opts.Risk).Grow(1)
			})
		}
		if opts.Acknowledge != "" {
			Checkbox(c, &st.ack, opts.Acknowledge, CheckboxOptions{})
		}
	})
	if was && !*open && !res.Confirmed {
		res.Canceled = true
	}
	if res.Confirmed || res.Canceled {
		*open = false
		slog.Debug("mujica: order_confirm_dialog close", "confirmed", res.Confirmed)
	}
	return res
}
