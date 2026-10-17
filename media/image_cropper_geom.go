package media

import (
	"image"
	"math"

	"github.com/egoist/mygo/ui"
)

// cropperMin is the smallest crop side in pixels.
const cropperMin = 8

// cropperDragRect moves (mode 1) or corner-resizes (2-5) start.
func cropperDragRect(start image.Rectangle, mode, dx, dy int, size image.Point, aspect float64) image.Rectangle {
	if mode == 1 {
		dx = max(-start.Min.X, min(size.X-start.Max.X, dx))
		dy = max(-start.Min.Y, min(size.Y-start.Max.Y, dy))
		return start.Add(image.Pt(dx, dy))
	}
	corners := [4]image.Point{start.Min, {start.Max.X, start.Min.Y}, {start.Min.X, start.Max.Y}, start.Max}
	moving, anchor := corners[mode-2], corners[5-mode]
	x := max(0, min(size.X, moving.X+dx))
	y := max(0, min(size.Y, moving.Y+dy))
	sx, sy := 1, 1
	room := image.Pt(size.X-anchor.X, size.Y-anchor.Y)
	if x < anchor.X {
		sx, room.X = -1, anchor.X
	}
	if y < anchor.Y {
		sy, room.Y = -1, anchor.Y
	}
	w, h := cropperClampSize(max(cropperMin, (x-anchor.X)*sx), max(cropperMin, (y-anchor.Y)*sy), room, aspect)
	return image.Rectangle{Min: anchor, Max: anchor.Add(image.Pt(w*sx, h*sy))}.Canon()
}

// cropperClampSize fits w×h inside room, keeping aspect when locked.
func cropperClampSize(w, h int, room image.Point, aspect float64) (int, int) {
	w, h = min(w, room.X), min(h, room.Y)
	if aspect <= 0 {
		return w, h
	}
	h = int(math.Round(float64(w) / aspect))
	if h > room.Y {
		h = room.Y
		w = int(math.Round(float64(h) * aspect))
	}
	return w, h
}

// cropperKeys: arrows move, Shift+arrows resize, 1% per press.
func cropperKeys(stage ui.Element, cr *ImageCrop, size image.Point, aspect float64) {
	step := max(1, size.X/100)
	r := cr.Rect
	move := func(dx, dy int) { cr.Rect = cropperDragRect(r, 1, dx, dy, size, aspect) }
	grow := func(dw, dh int) {
		w, h := max(cropperMin, r.Dx()+dw), max(cropperMin, r.Dy()+dh)
		if aspect > 0 && dw == 0 {
			w = int(math.Round(float64(h) * aspect))
		}
		w, h = cropperClampSize(w, h, image.Pt(size.X-r.Min.X, size.Y-r.Min.Y), aspect)
		cr.Rect = image.Rectangle{Min: r.Min, Max: r.Min.Add(image.Pt(w, h))}
	}
	switch {
	case stage.Shortcut(0, ui.KeyLeft):
		move(-step, 0)
	case stage.Shortcut(0, ui.KeyRight):
		move(step, 0)
	case stage.Shortcut(0, ui.KeyUp):
		move(0, -step)
	case stage.Shortcut(0, ui.KeyDown):
		move(0, step)
	case stage.Shortcut(ui.Shift, ui.KeyRight):
		grow(step, 0)
	case stage.Shortcut(ui.Shift, ui.KeyLeft):
		grow(-step, 0)
	case stage.Shortcut(ui.Shift, ui.KeyDown):
		grow(0, step)
	case stage.Shortcut(ui.Shift, ui.KeyUp):
		grow(0, -step)
	}
}
