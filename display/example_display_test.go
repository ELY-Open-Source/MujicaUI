package display_test

import (
	"errors"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
)

func ExampleText() {
	name := "Coronation.md"
	saveErr := ""
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		display.Text(c, "Order ID: 4821-AX", display.TextOptions{Selectable: true})
		display.Text(c, "A long note that stops after two lines.", display.TextOptions{MaxLines: 2})
		if display.EditableText(c, &name, display.EditableTextOptions{Error: saveErr}).Changed() {
			saveErr = "" // save name here; set saveErr on failure
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleHeading() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		display.Heading(c, "The Grand Hall", display.HeadingOptions{Subtitle: "Winter ceremonies", Divider: true})
		display.Heading(c, "Guest Ledger", display.HeadingOptions{Level: 2})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleIcon() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Gap(8).Children(func() {
			display.Icon(c, "star", display.IconOptions{})
			display.Icon(c, "circle-check", display.IconOptions{Size: 20, Tone: display.IconSuccess, Label: "Done"})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleImage() {
	var photo *ui.Bitmap // decode once with ui.DecodeBitmap, outside the view
	loading, loadErr := false, errors.New("timeout")
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		img := display.Image(c, photo, display.ImageOptions{
			Width: 160, Height: 120, Crop: true, Alt: "Hall of mirrors",
			Loading: loading, Err: loadErr,
		})
		if img.Retry {
			loading, loadErr = true, nil // start loading again here
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAvatar() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Gap(8).Children(func() {
			display.Avatar(c, "Sakiko Togawa", display.AvatarOptions{Size: 40})
			display.Avatar(c, "", display.AvatarOptions{})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAvatarGroup() {
	people := []display.AvatarItem{
		{Name: "Sakiko Togawa"}, {Name: "Uika Misumi"}, {Name: "Mutsumi Wakaba"},
		{Name: "Umiri Yahata"}, {Name: "Nyamu Yutenji"},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		display.AvatarGroup(c, people, display.AvatarGroupOptions{Max: 3})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBadge() {
	unread := 120
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		display.Badge(c, display.BadgeOptions{Text: "New", Tone: display.BadgeNeutral}, nil)
		display.Badge(c, display.BadgeOptions{Count: unread, Label: "Unread", Position: display.BadgeTopRight}, func() {
			display.Avatar(c, "Sakiko Togawa", display.AvatarOptions{})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTag() {
	tags := []string{"Baroque", "Gothic"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Gap(8).Children(func() {
			for i, t := range tags {
				var v display.TagView
				ui.Box(c.Key(t)).Children(func() {
					v = display.Tag(c, t, display.TagOptions{Icon: "star", Clickable: true, Closable: true})
				})
				if v.Closed {
					tags = append(tags[:i:i], tags[i+1:]...)
				}
			}
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLink() {
	details := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		display.Link(c, "Read the charter", display.LinkOptions{URL: "https://example.com/charter"})
		if display.Link(c, "Show details", display.LinkOptions{}).Clicked() {
			details = !details
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleKbd() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Gap(6).Children(func() {
			display.Text(c, "Command palette", display.TextOptions{})
			display.Kbd(c, "Mod", "Shift", "P")
		})
	}
	ui.Render(view, 400, 300, 1)
}
