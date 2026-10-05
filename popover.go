package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// PopoverOptions configure a Popover.
type PopoverOptions struct {
	// Title shows in serif at the top and names the panel.
	Title string
	// Width fixes the panel width; zero fits the content, at least the anchor's width.
	Width float32
	// Label names an untitled panel for assistive technology.
	Label string
}

// Popover shows content in a panel below anchor while *open is true. The
// panel flips above where there is no room, stays in the window and scrolls
// when taller than the room; an outside click or Escape closes it.
func Popover(c *ui.Context, anchor *ui.Element, open *bool, opts PopoverOptions, content func()) *ui.Element {
	if anchor == nil {
		panic("mujica: Popover needs an anchor")
	}
	if !*open {
		return nil
	}
	u := Density(c).Unit()
	b := anchor.Bounds()
	label := opts.Label
	if opts.Title != "" {
		label = opts.Title
	}
	panel := ui.PopoverBase(c, anchor, open, func(panel *ui.Element) {
		overlayStyle(c, panel, overlayAnchored(c, panel, b))
		panel.Margin(6, 0, 0, 0).Padding(u * 3).Gap(u * 2).MinWidth(b.W).MaxHeight(popoverRoom(c, b)).Label(label)
		if opts.Width > 0 {
			panel.Width(opts.Width)
		} else {
			panel.MaxWidth(360)
		}
		if opts.Title != "" {
			ui.Text(c, opts.Title).Font(theme.SerifFont).FontSize(theme.BodySize + 2).FontWeight(600).Shrink(0)
		}
		ui.Scroll(c).Shrink(1).Children(content)
	})
	if !*open {
		slog.Debug("mujica: popover close")
	}
	return panel
}

// popoverRoom is the tallest a panel below or above anchor b may be.
func popoverRoom(c *ui.Context, b ui.Rect) float32 {
	_, h := c.Size()
	return max(b.Y, h-b.Y-b.H) - 16
}
