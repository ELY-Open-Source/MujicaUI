package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// InputGroupButton is a button joined to an InputGroup's border.
type InputGroupButton struct {
	// Label is shown, or names an icon-only button for assistive technology.
	Label string
	Icon  *ui.SVG
	// IconOnly hides Label on screen; Icon is then required.
	IconOnly bool
	// Leading puts the button before the input.
	Leading  bool
	Disabled bool
}

// InputGroupOptions configure an InputGroup.
type InputGroupOptions struct {
	// Prefix and Suffix are fixed text beside the input, such as "https://".
	Prefix, Suffix string
	Buttons        []InputGroupButton
	Placeholder    string
	// Label names the input when it is not inside a FormField.
	Label    string
	Disabled bool
	// ReadOnly shows the value as selectable text that cannot change.
	ReadOnly bool
	// Invalid draws the border in the danger color.
	Invalid bool
}

// InputGroupParts are the parts of an InputGroup.
type InputGroupParts struct {
	Group *ui.Element
	// Input takes the focus; it is selectable text when ReadOnly.
	Input *ui.Element
	// Clicked is the index in Buttons of the button clicked, or -1.
	Clicked int
	locked  bool
}

// Changed reports whether the user edited the value this frame.
func (p InputGroupParts) Changed() bool { return !p.locked && p.Input.Changed() }

// Submitted reports whether Enter was pressed in the input.
func (p InputGroupParts) Submitted() bool { return !p.locked && p.Input.Submitted() }

// InputGroup joins an input, fixed text and buttons, each its own Tab stop, in one border.
func InputGroup(c *ui.Context, value *string, opts InputGroupOptions) InputGroupParts {
	for _, b := range opts.Buttons {
		if b.IconOnly && b.Icon == nil {
			panic("mujica: InputGroupButton IconOnly needs an Icon")
		}
		if b.Label == "" {
			panic("mujica: InputGroupButton needs a Label")
		}
	}
	k := Tokens(c)
	d := Density(c)
	parts := InputGroupParts{Clicked: -1}
	g := ui.Row(c).MinHeight(d.ControlHeight()).AlignItems(ui.Stretch).Radius(theme.ControlRadius).
		Padding(theme.BorderWidth).Background(k.Surface).Disabled(opts.Disabled)
	parts.Group = g
	r := theme.ControlRadius - theme.BorderWidth
	// Index of each segment, for its corners and separators.
	n := len(opts.Buttons) + 1
	if opts.Prefix != "" {
		n++
	}
	if opts.Suffix != "" {
		n++
	}
	seg := 0
	segment := func(e *ui.Element) {
		var left, right float32
		if seg == 0 {
			left = r
		}
		if seg == n-1 {
			right = r
		}
		e.Radius(left, right, right, left)
		if seg > 0 {
			e.BorderWidth(0, 0, 0, theme.BorderWidth).BorderColor(k.ControlBorder)
		}
		seg++
	}
	button := func(i int) {
		b := opts.Buttons[i]
		bb := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Padding(0, d.Unit()*3).Gap(d.Unit() * 2).MinWidth(d.ControlHeight()).
			Disabled(b.Disabled).TextColor(k.AccentText).Label(b.Label)
		segment(bb)
		internal.Ring(bb, k.Focus, theme.ControlRadius)
		if bb.Pressed() {
			bb.Background(k.Selection)
		} else if bb.Hovered() {
			bb.Background(k.SurfaceHover)
		}
		bb.Children(func() {
			if b.Icon != nil {
				ui.Icon(c, b.Icon).Size(16, 16)
			}
			if !b.IconOnly {
				ui.Text(c, b.Label).SingleLine()
			}
		})
		if bb.Clicked() {
			parts.Clicked = i
		}
	}
	addon := func(s string) {
		a := ui.Row(c).Padding(0, d.Unit()*3).Background(k.Background).Shrink(0)
		segment(a)
		a.Children(func() { ui.Text(c, s).TextColor(k.TextMuted).SingleLine() })
	}
	g.Children(func() {
		for i, b := range opts.Buttons {
			if b.Leading {
				button(i)
			}
		}
		if opts.Prefix != "" {
			addon(opts.Prefix)
		}
		box := ui.Row(c).Grow(1).MinWidth(d.Unit()*12).Padding(0, d.Unit()*3)
		segment(box)
		box.Children(func() {
			if opts.ReadOnly {
				// As wide as MyGo makes a text input.
				parts.Input = ui.Text(c, *value).SingleLine().Selectable().Focusable().FocusRing(false).Grow(1).Width(200).MinWidth(0)
			} else {
				parts.locked = box.IsDisabled()
				parts.Input = inputGroupEditor(c, box, value, parts.locked).FocusRing(false).Grow(1).MinWidth(0).Placeholder(opts.Placeholder)
			}
			if opts.Label != "" {
				parts.Input.Label(opts.Label)
			}
		})
		if box.Clicked() {
			parts.Input.Focus()
		}
		if opts.Suffix != "" {
			addon(opts.Suffix)
		}
		for i, b := range opts.Buttons {
			if !b.Leading {
				button(i)
			}
		}
	})
	in := parts.Input
	border := k.ControlBorder
	switch {
	case opts.Invalid:
		border = k.Danger
	case in.Focused():
		border = k.Focus
	case g.Hovered() && !opts.Disabled:
		border = k.ControlBorder.Mix(k.Text, 0.25)
	}
	g.Border(theme.BorderWidth, border)
	g.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if in.Focused() {
			internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
		}
	})
	return parts
}

// inputGroupEditor builds the input; a disabled one edits a copy, so queued input cannot reach value.
func inputGroupEditor(c *ui.Context, box *ui.Element, value *string, locked bool) *ui.Element {
	if !locked {
		return ui.TextInputBase(c, value)
	}
	shown := ui.Local(box, "locked", func() string { return *value })
	*shown = *value
	in := ui.TextInputBase(c, shown)
	*shown = *value
	return in
}
