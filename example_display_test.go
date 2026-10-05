package mujica_test

import (
	"errors"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleText() {
	name := "Coronation.md"
	saveErr := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Text(c, "Order ID: 4821-AX", mujica.TextOptions{Selectable: true})
		mujica.Text(c, "A long note that stops after two lines.", mujica.TextOptions{MaxLines: 2})
		if mujica.EditableText(c, &name, mujica.EditableTextOptions{Error: saveErr}).Changed() {
			saveErr = "" // save name here; set saveErr on failure
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleHeading() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Heading(c, "The Grand Hall", mujica.HeadingOptions{Subtitle: "Winter ceremonies", Divider: true})
		mujica.Heading(c, "Guest Ledger", mujica.HeadingOptions{Level: 2})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleIcon() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Row(c).Gap(8).Children(func() {
			mujica.Icon(c, "star", mujica.IconOptions{})
			mujica.Icon(c, "circle-check", mujica.IconOptions{Size: 20, Tone: mujica.IconSuccess, Label: "Done"})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleImage() {
	var photo *ui.Bitmap // decode once with ui.DecodeBitmap, outside the view
	loading, loadErr := false, errors.New("timeout")
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		img := mujica.Image(c, photo, mujica.ImageOptions{
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
		mujica.Use(c, mujica.Settings{})
		ui.Row(c).Gap(8).Children(func() {
			mujica.Avatar(c, "Sakiko Togawa", mujica.AvatarOptions{Size: 40})
			mujica.Avatar(c, "", mujica.AvatarOptions{})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAvatarGroup() {
	people := []mujica.AvatarItem{
		{Name: "Sakiko Togawa"}, {Name: "Uika Misumi"}, {Name: "Mutsumi Wakaba"},
		{Name: "Umiri Yahata"}, {Name: "Nyamu Yutenji"},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.AvatarGroup(c, people, mujica.AvatarGroupOptions{Max: 3})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleBadge() {
	unread := 120
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Badge(c, mujica.BadgeOptions{Text: "New", Tone: mujica.BadgeNeutral}, nil)
		mujica.Badge(c, mujica.BadgeOptions{Count: unread, Label: "Unread", Position: mujica.BadgeTopRight}, func() {
			mujica.Avatar(c, "Sakiko Togawa", mujica.AvatarOptions{})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTag() {
	tags := []string{"Baroque", "Gothic"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Row(c).Gap(8).Children(func() {
			for i, t := range tags {
				var v mujica.TagView
				ui.Box(c).Key(t).Children(func() {
					v = mujica.Tag(c, t, mujica.TagOptions{Icon: "star", Clickable: true, Closable: true})
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
		mujica.Use(c, mujica.Settings{})
		mujica.Link(c, "Read the charter", mujica.LinkOptions{URL: "https://example.com/charter"})
		if mujica.Link(c, "Show details", mujica.LinkOptions{}).Clicked() {
			details = !details
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleKbd() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Row(c).Gap(6).Children(func() {
			mujica.Text(c, "Command palette", mujica.TextOptions{})
			mujica.Kbd(c, "Mod", "Shift", "P")
		})
	}
	ui.Render(view, 400, 300, 1)
}
