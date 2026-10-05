package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// BreadcrumbItem is one level of a Breadcrumb path.
type BreadcrumbItem struct {
	Label string
	Icon  *ui.SVG
}

// BreadcrumbOptions configure a Breadcrumb.
type BreadcrumbOptions struct {
	// MaxItems shown before middle items fold into a menu; zero means 4, minimum 3.
	MaxItems int
	// MaxItemWidth cuts longer labels with an ellipsis; zero means 160 DIPs.
	MaxItemWidth float32
	// Label names the path for assistive technology; empty uses the built-in name.
	Label string
}

// BreadcrumbResult is the built Breadcrumb.
type BreadcrumbResult struct {
	*ui.Element
	chosen int
}

// Chosen returns the index of the item the user activated in this frame.
func (r BreadcrumbResult) Chosen() (int, bool) { return r.chosen, r.chosen >= 0 }

// Breadcrumb shows a path whose last item is the current one; earlier items are links.
func Breadcrumb(c *ui.Context, items []BreadcrumbItem, opts BreadcrumbOptions) BreadcrumbResult {
	if len(items) == 0 {
		panic("mujica: Breadcrumb needs at least one item")
	}
	limit := opts.MaxItems
	if limit == 0 {
		limit = 4
	}
	if limit < 3 {
		panic("mujica: Breadcrumb MaxItems must be at least 3")
	}
	width := opts.MaxItemWidth
	if width == 0 {
		width = 160
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "breadcrumb.label")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	row := ui.Row(c).AlignItems(ui.Center).Gap(u).MinWidth(0).Role(ui.RoleGroup).Label(label)
	res := BreadcrumbResult{Element: row, chosen: -1}
	choose := func(i int) {
		res.chosen = i
		slog.Debug("mujica: breadcrumb choose", "index", i, "label", items[i].Label)
	}
	// Folded items are those between the first and the last limit-2.
	fold := len(items) > limit
	hideFrom, hideTo := 1, len(items)-(limit-2)
	row.Children(func() {
		for i, it := range items {
			if it.Label == "" {
				panic("mujica: BreadcrumbItem needs a Label")
			}
			if fold && i >= hideFrom && i < hideTo {
				if i == hideFrom {
					breadcrumbSeparator(c)
					breadcrumbMore(c, items[hideFrom:hideTo], func(j int) { choose(hideFrom + j) })
				}
				continue
			}
			if i > 0 {
				breadcrumbSeparator(c)
			}
			if i == len(items)-1 {
				cur := ui.Row(c).AlignItems(ui.Center).Gap(u).Padding(u, u*1.5).MinWidth(0).
					FontWeight(600).TextColor(k.Text).Tooltip(it.Label).Label(it.Label + ", " + msg(c, "breadcrumb.current"))
				cur.Children(func() { breadcrumbContent(c, it, width) })
				continue
			}
			link := ui.ButtonBase(c).Role(ui.RoleLink).Gap(u).Padding(u, u*1.5).Radius(theme.ControlRadius).
				Shrink(1).MinWidth(0).Cursor(ui.CursorPointer).Tooltip(it.Label).Label(it.Label)
			internal.Ring(link, k.Focus, theme.ControlRadius)
			link.TextColor(k.TextMuted)
			if link.Hovered() {
				link.Background(k.SurfaceHover).TextColor(k.Text)
			}
			link.Children(func() { breadcrumbContent(c, it, width) })
			if link.Clicked() {
				choose(i)
			}
		}
	})
	return res
}

func breadcrumbContent(c *ui.Context, it BreadcrumbItem, width float32) {
	if it.Icon != nil {
		ui.Icon(c, it.Icon).Size(14, 14).Shrink(0)
	}
	ui.Text(c, it.Label).SingleLine().MaxWidth(width).MinWidth(0)
}

func breadcrumbSeparator(c *ui.Context) {
	ui.Icon(c, icons.Must("chevron-right")).Size(12, 12).Shrink(0).TextColor(Tokens(c).Ornament).Role(ui.RoleNone)
}

// breadcrumbMore is the ellipsis button listing the folded items.
func breadcrumbMore(c *ui.Context, hidden []BreadcrumbItem, choose func(int)) {
	k := Tokens(c)
	u := Density(c).Unit()
	b := ui.ButtonBase(c).Padding(u, u*1.5).Radius(theme.ControlRadius).MinWidth(theme.MinHitSize).
		Label(msg(c, "breadcrumb.more")).Tooltip(msg(c, "breadcrumb.more")).TextColor(k.TextMuted).Cursor(ui.CursorPointer)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	open := ui.Local(b, "open", func() bool { return false })
	if b.Hovered() || *open {
		b.Background(k.SurfaceHover).TextColor(k.Text)
	}
	if b.Clicked() {
		*open = !*open
	}
	b.Children(func() { ui.Icon(c, icons.Must("ellipsis")).Size(16, 16) })
	ui.PopoverBase(c, b, open, func(panel *ui.Element) {
		panel.Margin(u, 0, 0, 0).Padding(u).Gap(u*0.5).Radius(theme.PopupRadius).Background(k.Surface).
			Border(theme.BorderWidth, k.Border).Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, 0.16)).MinWidth(160).MaxWidth(320)
		for j, it := range hidden {
			item := ui.ButtonBase(c).Role(ui.RoleLink).Justify(ui.Start).Gap(u*2).Padding(u*1.5, u*2).
				Radius(theme.ControlRadius).MinHeight(Density(c).ControlHeight() - u*2).Label(it.Label).TextColor(k.Text)
			if j == 0 {
				item.AutoFocus()
			}
			internal.Ring(item, k.Focus, theme.ControlRadius)
			if item.Hovered() {
				item.Background(k.SurfaceHover)
			}
			item.Children(func() {
				if it.Icon != nil {
					ui.Icon(c, it.Icon).Size(14, 14).Shrink(0)
				}
				ui.Text(c, it.Label).SingleLine().MinWidth(0)
			})
			if item.Clicked() {
				*open = false
				choose(j)
			}
		}
	})
}
