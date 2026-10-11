package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
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
	Canceled  bool
}

// Popconfirm asks for a light confirmation by anchor; use AlertDialog for high-risk actions.
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
	core.OverlayLayer(c, open)
	u := core.Density(c).Unit()
	confirm, cancel := opts.ConfirmLabel, opts.CancelLabel
	if confirm == "" {
		confirm = core.Msg(c, "popconfirm.confirm")
	}
	if cancel == "" {
		cancel = core.Msg(c, "popconfirm.cancel")
	}
	variant := input.Primary
	if opts.Danger {
		variant = input.Danger
	}
	b := anchor.Bounds()
	shown := true
	ui.PopoverBase(c, anchor, &shown, func(panel *ui.Element) {
		core.OverlayStyle(c, panel, OverlayAnchored(c, panel, b))
		panel.Margin(6, 0, 0, 0).Padding(u * 3).Gap(u * 3).MaxWidth(PopoverWide(c, 320)).MaxHeight(PopoverRoom(c, b)).Role(ui.RoleAlertDialog).
			Label(opts.Title).Description(opts.Description)
		// The text scrolls so the actions stay in the window.
		ui.Scroll(c).Shrink(1).Children(func() { popconfirmText(c, opts) })
		ui.Row(c).Gap(u * 2).Justify(ui.End).Shrink(0).Children(func() {
			// Danger focuses Cancel, so a second Enter never runs the destructive action.
			cb := input.Button(c, cancel, input.ButtonOptions{Variant: input.Secondary})
			if opts.Danger {
				cb.AutoFocus()
			}
			if cb.Clicked() {
				res.Canceled = true
			}
			ok := input.Button(c, confirm, input.ButtonOptions{Variant: variant})
			if !opts.Danger {
				ok.AutoFocus()
			}
			if ok.Clicked() {
				res.Confirmed = true
			}
		})
	})
	if !shown {
		res.Canceled = true
	}
	if res.Confirmed || res.Canceled {
		slog.Debug("mujica: popconfirm close", "confirmed", res.Confirmed)
		*open = false
	}
	return res
}

// popconfirmText is the icon, title and description of a Popconfirm.
func popconfirmText(c *ui.Context, opts PopconfirmOptions) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
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
}
