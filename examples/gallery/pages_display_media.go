package main

import (
	"errors"
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/display"
)

type imageState struct {
	crop, loading, fail, empty bool
	err                        error
}

var displayPicture = func() *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="240" height="160" viewBox="0 0 240 160">
<rect width="240" height="160" fill="#2B2126"/><rect x="20" y="20" width="200" height="120" fill="none" stroke="#BDA172" stroke-width="2"/>
<path d="M120 40 L150 80 L120 120 L90 80 Z" fill="#7D2034"/><circle cx="120" cy="80" r="8" fill="#F2EBDD"/></svg>`))
}()

func imageDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "014", func() imageState { return imageState{crop: true} })
	display.Heading(c, "Image", display.HeadingOptions{Subtitle: "位图与 SVG；包含或裁切；占位、加载、错误；替代说明", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			displayLog(log, ui.Checkbox(c, &s.crop, "Crop"), "crop", &s.crop)
			displayLog(log, ui.Checkbox(c, &s.empty, "No source"), "empty", &s.empty)
			displayLog(log, ui.Checkbox(c, &s.loading, "Loading"), "loading", &s.loading)
			if ui.Checkbox(c, &s.fail, "Error").Changed() {
				log(fmt.Sprintf("fail = %v", s.fail))
				s.err = nil
				if s.fail {
					s.err = errors.New("connection reset")
				}
			}
		})
	})
	displaySection(c, "Preview", func() {
		var src ui.ImageSource = displayPicture
		if s.empty {
			src = nil
		}
		v := display.Image(c, src, display.ImageOptions{Width: 200, Height: 200, Crop: s.crop, Alt: "Gothic emblem", Loading: s.loading, Err: s.err})
		if v.Retry {
			s.fail, s.err = false, nil
			log("retry")
		}
	})
	displaySection(c, "States", func() {
		ui.Row(c).Gap(12).Wrap().Children(func() {
			display.Image(c, displayPicture, display.ImageOptions{Width: 120, Height: 120})
			display.Image(c, displayPicture, display.ImageOptions{Width: 120, Height: 120, Crop: true})
			display.Image(c, nil, display.ImageOptions{Width: 120, Height: 120})
			display.Image(c, nil, display.ImageOptions{Width: 120, Height: 120, Loading: true})
			if display.Image(c, nil, display.ImageOptions{Width: 160, Height: 140, Err: errors.New("404")}).Retry {
				log("retry (static example)")
			}
		})
	})
}

const imageCode = `img := display.Image(c, photo, display.ImageOptions{
	Width: 160, Height: 120, Crop: true, Alt: "Hall of mirrors",
	Loading: loading, Err: loadErr,
})
if img.Retry {
	loading, loadErr = true, nil // start loading again here
}`

type avatarState struct {
	name string
	size float64
	img  bool
}

var displayBitmap = func() *ui.Bitmap {
	return ui.NewBitmap(ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Dark})
		ui.Box(c).Fill().Background(ui.Hex("#7D2034"))
	}, 64, 64, 1))
}()

func avatarDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "015", func() avatarState { return avatarState{name: "Sakiko Togawa", size: 48} })
	display.Heading(c, "Avatar", display.HeadingOptions{Subtitle: "图片、姓名首字母与默认图标回退", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			if ui.TextInput(c, &s.name).Width(200).Changed() {
				log("name " + s.name)
			}
			ui.Text(c, fmt.Sprintf("Size %d", int(s.size)))
			displayLog(log, ui.StepSlider(c, &s.size, 24, 80, 4).Width(160), "size", &s.size)
			displayLog(log, ui.Checkbox(c, &s.img, "Image"), "img", &s.img)
		})
	})
	displaySection(c, "Preview", func() {
		var img *ui.Bitmap
		if s.img {
			img = displayBitmap
		}
		display.Avatar(c, s.name, display.AvatarOptions{Size: float32(s.size), Image: img})
	})
	displaySection(c, "Fallbacks", func() {
		ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
			display.Avatar(c, "Photo", display.AvatarOptions{Image: displayBitmap, Size: 40})
			display.Avatar(c, "Uika Misumi", display.AvatarOptions{Size: 40})
			display.Avatar(c, "若叶睦", display.AvatarOptions{Size: 40})
			display.Avatar(c, "", display.AvatarOptions{Size: 40})
		})
	})
}

const avatarCode = `ui.Row(c).Gap(8).Children(func() {
	display.Avatar(c, "Sakiko Togawa", display.AvatarOptions{Size: 40})
	display.Avatar(c, "", display.AvatarOptions{})
})`
