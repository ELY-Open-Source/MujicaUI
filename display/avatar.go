package display

import (
	"unicode"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// AvatarOptions configure an Avatar.
type AvatarOptions struct {
	// Image is the picture; nil falls back to initials, icon.
	Image *ui.Bitmap
	// Size is the diameter in DIPs; zero means 32.
	Size float32
}

// Avatar shows a person as picture, initials or icon.
func Avatar(c *ui.Context, name string, opts AvatarOptions) *ui.Element {
	size := opts.Size
	if size < 0 {
		panic("mujica: Avatar Size must not be negative")
	}
	if size == 0 {
		size = 32
	}
	k := core.Tokens(c)
	label := name
	if label == "" {
		label = core.Msg(c, "avatar.unnamed")
	}
	e := ui.Box(c).Size(size, size).Radius(size / 2).Clip().Shrink(0).Center().Role(ui.RoleImage).Label(label)
	initials := avatarInitials(name)
	e.Children(func() {
		switch {
		case opts.Image != nil:
			ui.Image(c, opts.Image).Fit(ui.Cover).Size(size, size)
		case initials != "":
			e.Background(k.Selection).Border(theme.BorderWidth, k.Border)
			ui.Text(c, initials).Font(theme.SerifFont).FontSize(size * avatarInitialsScale(initials)).FontWeight(600).
				TextColor(k.AccentText).SingleLine()
		default:
			e.Background(k.SurfaceHover).Border(theme.BorderWidth, k.Border)
			Icon(c, "user", IconOptions{Size: size * 0.5, Tone: IconMuted})
		}
	})
	return e
}

// avatarInitials returns first letters of first and last words.
func avatarInitials(name string) string {
	var first, last rune
	inWord := false
	for _, r := range name {
		letter := unicode.IsLetter(r) || unicode.IsDigit(r)
		if letter && !inWord {
			if first == 0 {
				first = unicode.ToUpper(r)
			} else {
				last = unicode.ToUpper(r)
			}
		}
		inWord = letter
	}
	switch {
	case first == 0:
		return ""
	case last == 0:
		return string(first)
	}
	return string([]rune{first, last})
}

// avatarInitialsScale fits two letters inside the circle.
func avatarInitialsScale(initials string) float32 {
	if len([]rune(initials)) > 1 {
		return 0.34
	}
	return 0.42
}
