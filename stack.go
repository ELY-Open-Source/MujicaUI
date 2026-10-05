package mujica

import "github.com/egoist/mygo/ui"

// StackAlign places children across a Stack's direction. Zero stretches a vertical stack and centers a horizontal one.
type StackAlign int

const (
	StackAlignAuto StackAlign = iota
	StackAlignStart
	StackAlignCenter
	StackAlignEnd
	StackAlignStretch
)

// StackOptions configure a Stack.
type StackOptions struct {
	// Horizontal lays children left to right; otherwise top to bottom.
	Horizontal bool
	// Gap defaults to SpaceS.
	Gap   Space
	Align StackAlign
	// Justify places children along the direction; zero is ui.Start.
	Justify ui.Align
	// Wrap moves children of a horizontal stack onto more lines.
	Wrap bool
}

// Stack arranges children in a row or column. Build children with StackItem to hide them without leaving a gap.
func Stack(c *ui.Context, opts StackOptions, children func()) *ui.Element {
	if children == nil {
		panic("mujica: Stack needs children")
	}
	if opts.Wrap && !opts.Horizontal {
		panic("mujica: Stack Wrap needs Horizontal")
	}
	e := ui.Column(c)
	if opts.Horizontal {
		e.Row()
	}
	e.Gap(containerSpace(Density(c), opts.Gap, SpaceS)).Justify(opts.Justify)
	switch opts.Align {
	case StackAlignAuto:
	case StackAlignStart:
		e.AlignItems(ui.Start)
	case StackAlignCenter:
		e.AlignItems(ui.Center)
	case StackAlignEnd:
		e.AlignItems(ui.End)
	case StackAlignStretch:
		e.AlignItems(ui.Stretch)
	default:
		panic("mujica: invalid StackAlign")
	}
	if opts.Wrap {
		e.Wrap()
	}
	return e.Children(children)
}

// StackItem builds one child of a Stack unless hidden; a hidden item takes no room and no gap.
func StackItem(hidden bool, build func()) {
	if build == nil {
		panic("mujica: StackItem needs build")
	}
	if !hidden {
		build()
	}
}
