package finance

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/overlay"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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
func OrderConfirmDialog(c *ui.Context, open *bool, opts OrderConfirmDialogOptions) overlay.AlertDialogResult {
	if len(opts.Lines) == 0 {
		panic("mujica: OrderConfirmDialog needs summary Lines")
	}
	var res overlay.AlertDialogResult
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	title, confirm := opts.Title, opts.ConfirmLabel
	if title == "" {
		title = core.Msg(c, "order_confirm_dialog.title")
	}
	if confirm == "" {
		confirm = core.Msg(c, "order_confirm_dialog.confirm")
	}
	was := *open
	var st *orderConfirmState
	overlay.Dialog(c, open, overlay.DialogOptions{Title: title, KeepOnBackdrop: true, HideClose: true, Actions: func() {
		if input.Button(c, core.Msg(c, "order_confirm_dialog.cancel"), input.ButtonOptions{Variant: input.Secondary}).AutoFocus().Clicked() {
			res.Canceled = true
		}
		locked := opts.Acknowledge != "" && !st.ack
		if input.Button(c, confirm, input.ButtonOptions{Disabled: locked}).Clicked() && !locked {
			res.Confirmed = true
		}
	}}, func() {
		list := ui.Column(c)
		// The tick lives with the open dialog only.
		st = ui.Local(list, "ack", func() orderConfirmState { return orderConfirmState{} })
		list.Children(func() {
			for _, l := range opts.Lines {
				ui.Row(c).Gap(u*4).MinHeight(core.Density(c).ControlHeight()).AlignItems(ui.Center).
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
			input.Checkbox(c, &st.ack, opts.Acknowledge, input.CheckboxOptions{})
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
