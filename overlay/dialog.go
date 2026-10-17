package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DialogOptions configure a Dialog. Title is required.
type DialogOptions struct {
	Title       string
	Description string
	// Width is 440 DIPs when zero; the window always bounds it.
	Width float32
	// Decorated draws the double border, corner ornaments and title rule.
	Decorated bool
	// Actions builds the footer, usually buttons, aligned to the end.
	Actions func()
	// KeepOnBackdrop ignores clicks on the backdrop.
	KeepOnBackdrop bool
	// KeepOnEscape ignores Escape.
	KeepOnEscape bool
	// HideClose removes the close button from the header.
	HideClose bool
}

// Dialog shows a modal dialog while *open is true, trapping and then restoring the focus.
func Dialog(c *ui.Context, open *bool, opts DialogOptions, content func()) ui.Element {
	if opts.Title == "" {
		panic("mujica: Dialog needs a Title")
	}
	site, focus := core.OverlaySite(c)
	if !*open {
		focus.Closed(c, site)
		return ui.Element{}
	}
	core.OverlayLayer(c, open)
	k := core.Tokens(c)
	shown, backClicked, closed := true, false, false
	panel := ui.DialogBase(c, &shown, func(back, panel ui.Element) {
		backClicked = back.Clicked()
		p := core.OverlayEnter(c, panel)
		back.Background(OverlayScrim(k).Alpha(p)).Padding(16)
		DialogPanel(c, panel, opts.Width, opts.Decorated, p)
		panel.Label(opts.Title)
		if opts.Description != "" {
			panel.Description(opts.Description)
		}
		closed = dialogHeader(c, opts.Title, opts.Decorated, !opts.HideClose, "dialog.close")
		// Description and content scroll together so the actions always keep their room.
		ui.Scroll(c).Shrink(1).Gap(core.Density(c).Unit() * 3).Children(func() {
			if opts.Description != "" {
				ui.Text(c, opts.Description).TextColor(k.TextMuted)
			}
			if content != nil {
				content()
			}
		})
		if opts.Actions != nil {
			dialogActions(c, opts.Actions)
		}
	})
	focus.Shown(panel)
	switch {
	case closed:
		slog.Debug("mujica: dialog close", "by", "button")
		*open = false
	case !shown && backClicked && !opts.KeepOnBackdrop:
		slog.Debug("mujica: dialog close", "by", "backdrop")
		*open = false
	case !shown && !backClicked && !opts.KeepOnEscape:
		slog.Debug("mujica: dialog close", "by", "escape")
		*open = false
	}
	return panel
}

// DialogPanel styles a centered dialog panel at entrance progress p.
func DialogPanel(c *ui.Context, panel ui.Element, width float32, decorated bool, p float32) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	if width == 0 {
		width = 440
	}
	w, h := c.Size()
	core.OverlayStyle(c, panel, p)
	panel.Width(width).MaxWidth(w-32).MaxHeight(h-32).Padding(u*6).Gap(u*3).Margin(8*(1-p), 0, 0, 0)
	if decorated {
		panel.Border(theme.BorderWidth, k.Ornament).Padding(u*6 + 6).DrawOver(func(pt *ui.Painter, r ui.Rect) {
			internal.Frame(pt, r, 5, 7, k.Ornament)
		})
	}
}

// dialogHeader builds the title and close button; it reports a close click.
func dialogHeader(c *ui.Context, title string, decorated, closable bool, closeKey string) bool {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	clicked := false
	ui.Column(c).Gap(u).Shrink(0).Children(func() {
		ui.Row(c).Gap(u * 3).AlignItems(ui.Start).Children(func() {
			ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600).Grow(1)
			if closable {
				clicked = core.IconAction(c, icons.Must("x"), core.Msg(c, closeKey)).Margin(-4, -4, 0, 0).Clicked()
			}
		})
		if decorated {
			dialogRule(c, k.Ornament)
		}
	})
	return clicked
}

// dialogRule is the ornament line under a decorated title.
func dialogRule(c *ui.Context, ink ui.Color) {
	ui.Box(c).FillWidth().Height(9).Margin(2, 0, 2, 0).Draw(func(p *ui.Painter, r ui.Rect) {
		internal.RuleDiamond(p, r, 10, 7, ink)
	})
}

// dialogActions lays out a dialog footer.
func dialogActions(c *ui.Context, actions func()) {
	u := core.Density(c).Unit()
	ui.Row(c).Gap(u*2).Justify(ui.End).Wrap().Shrink(0).Margin(u*2, 0, 0, 0).Children(actions)
}
