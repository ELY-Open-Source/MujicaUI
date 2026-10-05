package layout

import (
	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Space is a spacing step scaled by the window's density. Zero picks the component's default.
type Space int

const (
	SpaceDefault Space = iota
	SpaceNone
	SpaceXS  // 1 unit, 4 DIPs at default density
	SpaceS   // 2 units
	SpaceM   // 3 units
	SpaceL   // 4 units
	SpaceXL  // 6 units
	SpaceXXL // 8 units
)

// containerSpace returns s in DIPs for density d, using def when s is SpaceDefault.
func containerSpace(d theme.Density, s, def Space) float32 {
	if s == SpaceDefault {
		s = def
	}
	var u float32
	switch s {
	case SpaceNone:
	case SpaceXS:
		u = 1
	case SpaceS:
		u = 2
	case SpaceM:
		u = 3
	case SpaceL:
		u = 4
	case SpaceXL:
		u = 6
	case SpaceXXL:
		u = 8
	default:
		panic("mujica: invalid Space value")
	}
	return u * d.Unit()
}

// ContainerOptions configure a Container.
type ContainerOptions struct {
	// Padding defaults to SpaceL.
	Padding Space
	// MaxWidth caps the width in DIPs; zero leaves it unbounded.
	MaxWidth float32
	// Center places the container in the middle of its parent's width.
	Center bool
	// Fill takes the parent's whole content box.
	Fill bool
	// Clip cuts content that overflows the padding box.
	Clip bool
}

// Container pads its children, caps their width and can center, fill and clip.
func Container(c *ui.Context, opts ContainerOptions, children func()) *ui.Element {
	if opts.MaxWidth < 0 {
		panic("mujica: Container MaxWidth must not be negative")
	}
	if children == nil {
		panic("mujica: Container needs children")
	}
	e := ui.Column(c).Padding(containerSpace(core.Density(c), opts.Padding, SpaceL))
	if opts.MaxWidth > 0 {
		e.MaxWidth(opts.MaxWidth)
	}
	if opts.Fill {
		e.Fill()
	}
	if opts.Center {
		e.FillWidth().Margin(0, ui.Auto)
	}
	if opts.Clip {
		e.Clip()
	}
	return e.Children(children)
}
