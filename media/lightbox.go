package media

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Lightbox fills the window with images[*index] while *open.
func Lightbox(c *ui.Context, open *bool, index *int, images []data.ViewerImage) core.ChoiceResult {
	var v core.ChoiceResult
	if !*open {
		return v
	}
	if len(images) == 0 || *index < 0 || *index >= len(images) {
		panic("mujica: Lightbox index out of range")
	}
	for _, im := range images {
		data.ViewerSize(im.Source)
	}
	u := core.Density(c).Unit()
	w, h := c.Size()
	light := ui.Hex("#F2EBDD")
	v.Element = ui.DialogBase(c, open, func(back, panel ui.Element) {
		back.Background(ui.RGBA(0, 0, 0, 0.92))
		panel.Size(w, h).Padding(u * 4).Gap(u * 3).TextColor(light).Label(core.Msg(c, "lightbox.title")).
			Opacity(core.OverlayEnter(c, panel))
		goTo := func(i int) {
			if i < 0 || i >= len(images) || i == *index {
				return
			}
			*index, v = i, core.NewChoiceResult(v.Element, true)
			c.Announce(fmt.Sprintf(core.Msg(c, "lightbox.position"), i+1, len(images)))
			slog.Debug("mujica: lightbox image", "index", i)
		}
		switch {
		case panel.Shortcut(0, ui.KeyLeft):
			goTo(*index - 1)
		case panel.Shortcut(0, ui.KeyRight):
			goTo(*index + 1)
		case panel.Shortcut(0, ui.KeyHome):
			goTo(0)
		case panel.Shortcut(0, ui.KeyEnd):
			goTo(len(images) - 1)
		}
		img := images[*index]
		ui.Row(c).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "lightbox.position"), *index+1, len(images))).FontFeatures("tnum").Grow(1)
			if lightboxButton(c, "x", "lightbox.close", false).Clicked() {
				*open = false
			}
		})
		ui.Row(c).Grow(1).MinHeight(0).Gap(u * 3).AlignItems(ui.Center).Children(func() {
			if b := lightboxButton(c, "chevron-left", "lightbox.previous", *index == 0); b.Clicked() && !b.IsDisabled() {
				goTo(*index - 1)
			}
			ui.Box(c).Grow(1).FillHeight().MinWidth(0).Center().Children(func() {
				ui.Image(c, img.Source).Fit(ui.Contain).Fill().Label(img.Alt)
			})
			if b := lightboxButton(c, "chevron-right", "lightbox.next", *index == len(images)-1); b.Clicked() && !b.IsDisabled() {
				goTo(*index + 1)
			}
		})
		if img.Alt != "" {
			ui.Text(c, img.Alt).TextAlign(ui.Center).FillWidth().MaxLines(3)
		}
	})
	if !*open {
		slog.Debug("mujica: lightbox close")
	}
	return v
}

// lightboxButton is a 40 DIP icon button on the backdrop.
func lightboxButton(c *ui.Context, icon, key string, disabled bool) ui.Element {
	light := ui.Hex("#F2EBDD")
	b := core.IconAction(c, icons.Must(icon), core.Msg(c, key)).Size(40, 40).TextColor(light).Disabled(disabled)
	// The theme's focus color is too dark on the backdrop.
	internal.Ring(b, light, theme.ControlRadius)
	if b.Hovered() {
		b.Background(ui.RGBA(255, 255, 255, 0.12)).TextColor(ui.Hex("#FFFFFF"))
	}
	return b
}
