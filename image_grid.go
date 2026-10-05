package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// GridImage is one picture of an ImageGrid; the caller loads it.
type GridImage struct {
	Source ui.ImageSource
	Alt    string
	// Loading and Err show the caller's load in progress or failure.
	Loading bool
	Err     error
}

// ImageGridResult is an ImageGrid's element and a retry the user asked for.
type ImageGridResult struct {
	Element *ui.Element
	retry   int
}

// Retried returns the index of the failed picture whose retry was pressed this frame.
func (r ImageGridResult) Retried() (int, bool) { return r.retry, r.retry >= 0 }

// imageGridState is the viewer opened from the grid.
type imageGridState struct {
	open  bool
	index int
}

// ImageGrid lays out 1 to 9 pictures; a click on a loaded one opens the ImageViewer.
func ImageGrid(c *ui.Context, images []GridImage) ImageGridResult {
	n := len(images)
	if n < 1 || n > 9 {
		panic("mujica: ImageGrid shows 1 to 9 images")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	cols := 3
	switch n {
	case 1:
		cols = 1
	case 2, 4:
		cols = 2
	}
	cell := float32(120)
	if n == 1 {
		cell = 280
	}
	r := ImageGridResult{retry: -1}
	r.Element = ui.Grid(c).Columns(cols).Gap(u).AlignSelf(ui.Start)
	st := ui.Local(r.Element, "viewer", func() imageGridState { return imageGridState{} })
	var viewable []ViewerImage
	var at []int
	for i, im := range images {
		if im.Source != nil && !im.Loading && im.Err == nil {
			at = append(at, i)
			viewable = append(viewable, ViewerImage{Source: im.Source, Alt: im.Alt})
		}
	}
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		for i, im := range images {
			loaded := im.Source != nil && !im.Loading && im.Err == nil
			// Only a loaded picture is a button; a failed one keeps its retry.
			var b *ui.Element
			if loaded {
				b = ui.ButtonBase(c).Label(im.Alt)
				internal.Ring(b, k.Focus, theme.ControlRadius)
			} else {
				b = ui.Box(c)
			}
			b.Size(cell, cell).Radius(theme.ControlRadius).Clip().Children(func() {
				v := Image(c, im.Source, ImageOptions{Width: cell, Height: cell, Crop: true, Alt: im.Alt, Loading: im.Loading, Err: im.Err})
				if v.Retry && !disabled {
					r.retry = i
					slog.Debug("mujica: image_grid retry", "index", i)
				}
			})
			if loaded && b.Clicked() && !disabled {
				for vi, ii := range at {
					if ii == i {
						st.open, st.index = true, vi
					}
				}
				slog.Debug("mujica: image_grid open", "index", i)
			}
		}
	})
	if disabled {
		st.open = false
	}
	if st.open && len(viewable) > 0 {
		st.index = min(st.index, len(viewable)-1)
		ImageViewer(c, &st.open, &st.index, viewable)
	}
	return r
}
