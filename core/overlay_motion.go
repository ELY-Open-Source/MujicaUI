package core

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// OverlayEnter returns how far, 0 to 1, the popup born with e has entered.
func OverlayEnter(c *ui.Context, e ui.Element) float32 {
	born := ui.Local(e, "mujica.overlay.born", func() time.Time { return c.Now() })
	d := Motion(c, theme.PopupDuration)
	if d == 0 {
		return 1
	}
	t := float32(c.Now().Sub(*born)) / float32(d)
	if t >= 1 {
		return 1
	}
	c.AnimationFrame()
	return ui.EaseOut(t)
}

// OverlayStyle gives a popup panel the surface, edge and shadow of the theme.
func OverlayStyle(c *ui.Context, panel ui.Element, p float32) {
	k := Tokens(c)
	panel.Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.PopupRadius).
		TextColor(k.Text).Opacity(p).Shadow(0, 6, 18, 0, ui.RGBA(0, 0, 0, 0.16*p))
}

// PlaceOverlay puts a w×h box below anchor a, above when there is no room, inside the window.
func PlaceOverlay(c *ui.Context, a ui.Rect, w, h, gap float32) (x, y float32) {
	ww, wh := c.Size()
	const edge = 4
	x, y = a.X, a.Y+a.H+gap
	if y+h > wh-edge && a.Y-gap-h >= edge {
		y = a.Y - gap - h
	}
	x = max(edge, min(x, ww-edge-w))
	y = max(edge, min(y, wh-edge-h))
	return x, y
}
