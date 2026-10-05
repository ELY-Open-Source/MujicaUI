package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// IconTone is an icon's semantic color. The zero value takes the text color.
type IconTone int

const (
	IconInherit IconTone = iota
	IconMuted
	IconAccent
	IconSuccess
	IconWarning
	IconDanger
	IconInfo
)

// IconOptions configure an Icon.
type IconOptions struct {
	// Size is the icon's width and height in DIPs; zero means 16.
	Size float32
	Tone IconTone
	// Label names a meaningful icon for assistive technology; empty hides it.
	Label string
}

type iconMissingKey struct{}

// Icon shows the named icon of the icons package. A missing name shows a
// dashed fallback box in the Danger color and logs a warning.
func Icon(c *ui.Context, name string, opts IconOptions) *ui.Element {
	if name == "" {
		panic("mujica: Icon needs a name")
	}
	size := opts.Size
	if size < 0 {
		panic("mujica: Icon Size must not be negative")
	}
	if size == 0 {
		size = 16
	}
	k := Tokens(c)
	svg, ok := icons.Get(name)
	if !ok {
		return iconMissing(c, name, size, k)
	}
	e := ui.Icon(c, svg).Size(size, size).Shrink(0)
	if col, set := iconColor(k, opts.Tone); set {
		e.TextColor(col)
	}
	if opts.Label != "" {
		e.Label(opts.Label)
	}
	return e
}

func iconMissing(c *ui.Context, name string, size float32, k theme.Tokens) *ui.Element {
	text := fmt.Sprintf(msg(c, "icon.missing"), name)
	e := ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleImage).Label(text).Tooltip(text)
	logged := ui.Local(e, iconMissingKey{}, func() string { return "" })
	if *logged != name {
		*logged = name
		slog.Warn("mujica: icon missing", "name", name)
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		in := ui.Rect{X: r.X + 1, Y: r.Y + 1, W: r.W - 2, H: r.H - 2}
		p.StrokeDashed(in, k.Danger, theme.ControlRadius, 1)
		p.Line(in.X+3, in.Y+3, in.X+in.W-3, in.Y+in.H-3, 1, k.Danger)
	})
}

func iconColor(k theme.Tokens, t IconTone) (ui.Color, bool) {
	switch t {
	case IconInherit:
		return ui.Color{}, false
	case IconMuted:
		return k.TextMuted, true
	case IconAccent:
		return k.AccentText, true
	case IconSuccess:
		return k.Success, true
	case IconWarning:
		return k.Warning, true
	case IconDanger:
		return k.Danger, true
	case IconInfo:
		return k.Info, true
	}
	panic(fmt.Sprintf("mujica: unknown IconTone %d", t))
}
