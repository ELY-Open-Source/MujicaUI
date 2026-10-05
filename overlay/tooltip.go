package overlay

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TooltipOptions configure a Tooltip.
type TooltipOptions struct {
	// Delay is how long hover or focus rests before the tip shows: 500 ms when zero.
	Delay time.Duration
}

// TooltipState tracks when the anchor became hot and whether Escape hid the tip.
type TooltipState struct {
	Since     time.Time
	Dismissed bool
}

const tooltipMaxWidth = 240

// Tooltip shows a non-interactive tip by anchor after hover or keyboard focus rests on it for Delay.
func Tooltip(c *ui.Context, anchor *ui.Element, text string, opts TooltipOptions) *ui.Element {
	if anchor == nil {
		panic("mujica: Tooltip needs an anchor")
	}
	if opts.Delay < 0 {
		panic("mujica: Tooltip Delay is negative")
	}
	if text == "" {
		return anchor
	}
	anchor.Description(text)
	st := ui.Local(anchor, "mujica.tooltip", func() TooltipState { return TooltipState{} })
	hot := anchor.Hovered() || anchor.FocusVisible()
	if !hot {
		st.Since, st.Dismissed = time.Time{}, false
		return anchor
	}
	if st.Dismissed {
		return anchor
	}
	if st.Since.IsZero() {
		st.Since = c.Now()
	}
	delay := opts.Delay
	if delay == 0 {
		delay = 500 * time.Millisecond
	}
	if wait := delay - c.Now().Sub(st.Since); wait > 0 {
		c.After(wait)
		return anchor
	}
	if core.OverlayEscape(c, st) {
		slog.Debug("mujica: tooltip dismiss")
		st.Dismissed = true
		return anchor
	}
	tooltipShow(c, anchor.Bounds(), text)
	return anchor
}

// tooltipShow draws the tip below a, or above it when there is no room.
func tooltipShow(c *ui.Context, a ui.Rect, text string) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	size := core.FontSize(c, theme.CaptionSize+1)
	tw, th := c.MeasureText(tooltipMaxWidth, ui.Span{Text: text, Size: size, Font: c.Theme().Font})
	padX, padY := u*2.5, u*1.5
	w, h := tw+2*(padX+theme.BorderWidth)+1, th+2*(padY+theme.BorderWidth)
	x, y := core.PlaceOverlay(c, ui.Rect{X: a.X + a.W/2 - w/2, Y: a.Y, W: w, H: a.H}, w, h, 6)
	ui.Overlay(c, func() {
		ui.Box(c).Absolute().Left(x).Top(y).Width(w).Padding(padY, padX).Radius(theme.ControlRadius).
			Background(k.Text).TextColor(k.Background).Border(theme.BorderWidth, k.Text).
			PassThrough().Role(ui.RoleTooltip).Children(func() {
			ui.Text(c, text).FontSize(size)
		})
	})
}
