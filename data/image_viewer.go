package data

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ViewerImage is one picture of an ImageViewer.
type ViewerImage struct {
	// Source is a *ui.Bitmap or *ui.SVG; required.
	Source ui.ImageSource
	// Alt describes the picture for assistive technology and the caption.
	Alt string
}

// Zoom limits, as multiples of the fitted size.
const (
	viewerMinZoom  = 0.25
	viewerMaxZoom  = 8
	viewerZoomStep = 1.25
	viewerPanStep  = 48
)

type viewerState struct {
	image      int
	zoom, x, y float32
}

// ImageViewer shows images[*index] over the window while *open, with zoom, pan and switching; Element is nil while closed.
func ImageViewer(c *ui.Context, open *bool, index *int, images []ViewerImage) core.ChoiceResult {
	var v core.ChoiceResult
	if !*open {
		return v
	}
	if len(images) == 0 || *index < 0 || *index >= len(images) {
		panic("mujica: ImageViewer index out of range")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	w, h := c.Size()
	v.Element = ui.DialogBase(c, open, func(back, panel ui.Element) {
		back.Background(ui.RGBA(0, 0, 0, 0.55))
		panel.Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.PopupRadius).
			Size(w-u*12, h-u*12).Padding(u * 3).Gap(u * 3).Label(core.Msg(c, "image_viewer.title"))
		st := ui.Local(panel, "view", func() viewerState { return viewerState{image: -1} })
		if st.image != *index {
			*st = viewerState{image: *index, zoom: 1}
		}
		goTo := func(i int) {
			if i < 0 || i >= len(images) || i == *index {
				return
			}
			*index, v = i, core.NewChoiceResult(v.Element, true)
			*st = viewerState{image: i, zoom: 1}
			c.Announce(fmt.Sprintf(core.Msg(c, "image_viewer.position"), i+1, len(images)))
			slog.Debug("mujica: image_viewer image", "index", i)
		}
		zoom := func(f float32) {
			st.zoom = float32(math.Max(viewerMinZoom, math.Min(viewerMaxZoom, float64(st.zoom*f))))
		}
		viewerKeys(panel, st, zoom, goTo, *index)
		img := images[*index]
		tool := func(icon, label string, disabled bool) bool {
			return input.Button(c, "", input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must(icon), Label: core.Msg(c, label), Disabled: disabled}).Clicked()
		}
		ui.Row(c).Gap(u).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "image_viewer.position"), *index+1, len(images))).TextColor(k.TextMuted).Grow(1)
			if tool("zoom-out", "image_viewer.zoomOut", st.zoom <= viewerMinZoom) {
				zoom(1 / viewerZoomStep)
			}
			ui.Text(c, fmt.Sprintf("%d%%", int(math.Round(float64(st.zoom*100))))).Width(48).TextAlign(ui.Center).SingleLine()
			if tool("zoom-in", "image_viewer.zoomIn", st.zoom >= viewerMaxZoom) {
				zoom(viewerZoomStep)
			}
			if tool("maximize", "image_viewer.fit", st.zoom == 1 && st.x == 0 && st.y == 0) {
				*st = viewerState{image: *index, zoom: 1}
			}
			if tool("x", "image_viewer.close", false) {
				*open = false
			}
		})
		ui.Row(c).Grow(1).MinHeight(0).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			if tool("chevron-left", "image_viewer.prev", *index == 0) {
				goTo(*index - 1)
			}
			stage := core.AutoFocus(c, ui.Box(c).Grow(1).FillHeight().MinWidth(0).Focusable()).Clip().
				Radius(theme.ControlRadius).Background(k.Background).Role(ui.RoleImage).Label(img.Alt).Cursor(ui.CursorGrab)
			internal.Ring(stage, k.Focus, theme.ControlRadius)
			if dx, dy, ok := stage.Dragged(); ok {
				st.x, st.y = st.x+dx, st.y+dy
			}
			stage.HandleInput(func(ev ui.InputEvent) bool {
				if ev.Kind != ui.InputScroll || ev.DY == 0 {
					return false
				}
				if ev.DY < 0 {
					zoom(viewerZoomStep)
				} else {
					zoom(1 / viewerZoomStep)
				}
				return true
			})
			iw, ih := ViewerSize(img.Source)
			place := func(r ui.Rect) (dw, dh, mx, my float32) {
				fit := float32(math.Min(float64(r.W/iw), float64(r.H/ih)))
				dw, dh = iw*fit*st.zoom, ih*fit*st.zoom
				// The picture cannot be dragged past its edges.
				return dw, dh, max(0, (dw-r.W)/2), max(0, (dh-r.H)/2)
			}
			if b := stage.Bounds(); b.W > 0 {
				_, _, mx, my := place(b)
				st.x, st.y = min(mx, max(-mx, st.x)), min(my, max(-my, st.y))
			}
			x, y := st.x, st.y
			stage.Draw(func(p *ui.Painter, r ui.Rect) {
				dw, dh, mx, my := place(r)
				ox, oy := min(mx, max(-mx, x)), min(my, max(-my, y))
				p.Clip(r, theme.ControlRadius, func() {
					p.Image(img.Source, ui.Rect{X: r.X + (r.W-dw)/2 + ox, Y: r.Y + (r.H-dh)/2 + oy, W: dw, H: dh}, ui.FillBox)
				})
			})
			if tool("chevron-right", "image_viewer.next", *index == len(images)-1) {
				goTo(*index + 1)
			}
		})
		if img.Alt != "" {
			ui.Text(c, img.Alt).TextColor(k.TextMuted).TextAlign(ui.Center).FillWidth()
		}
	})
	if !*open {
		slog.Debug("mujica: image_viewer close")
	}
	return v
}

// viewerKeys, anywhere in the panel, zooms with + - 0, switches with Left and Right, and pans with Shift+arrows.
func viewerKeys(panel ui.Element, st *viewerState, zoom func(float32), goTo func(int), at int) {
	switch {
	case panel.Shortcut(0, ui.KeyEqual), panel.Shortcut(ui.Shift, ui.KeyEqual):
		zoom(viewerZoomStep)
	case panel.Shortcut(0, ui.KeyMinus):
		zoom(1 / viewerZoomStep)
	case panel.Shortcut(0, ui.Key0):
		*st = viewerState{image: st.image, zoom: 1}
	case panel.Shortcut(0, ui.KeyLeft):
		goTo(at - 1)
	case panel.Shortcut(0, ui.KeyRight):
		goTo(at + 1)
	case panel.Shortcut(ui.Shift, ui.KeyLeft):
		st.x += viewerPanStep
	case panel.Shortcut(ui.Shift, ui.KeyRight):
		st.x -= viewerPanStep
	case panel.Shortcut(ui.Shift, ui.KeyUp):
		st.y += viewerPanStep
	case panel.Shortcut(ui.Shift, ui.KeyDown):
		st.y -= viewerPanStep
	}
}

// ViewerSize returns a picture's own size in DIPs.
func ViewerSize(src ui.ImageSource) (float32, float32) {
	var w, h float32
	switch s := src.(type) {
	case *ui.Bitmap:
		iw, ih := s.Size()
		w, h = float32(iw), float32(ih)
	case *ui.SVG:
		w, h = s.Size()
	default:
		panic(fmt.Sprintf("mujica: ImageViewer source %T is not a *ui.Bitmap or *ui.SVG", src))
	}
	if w <= 0 || h <= 0 {
		panic("mujica: ImageViewer source without a size")
	}
	return w, h
}
