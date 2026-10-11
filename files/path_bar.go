package files

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PathBarOptions configure a PathBar.
type PathBarOptions struct {
	// Separator splits the path; empty is "/".
	Separator string
	// RootLabel defaults to the absolute path separator.
	RootLabel string
	Disabled  bool
}

// PathBarResult exposes the toolbar and requested path.
type PathBarResult struct {
	Element *ui.Element
	// Navigated is the chosen ancestor path, or empty.
	Navigated string
}

// pathSegment is a crumb and the path up to it.
type pathSegment struct{ label, path string }

// pathSegments splits p into crumbs with cumulative paths.
func pathSegments(p, sep, root string) []pathSegment {
	var out []pathSegment
	acc := ""
	if strings.HasPrefix(p, sep) {
		label := root
		if label == "" {
			label = sep
		}
		out = append(out, pathSegment{label, sep})
		acc = sep
	}
	for _, part := range strings.Split(strings.Trim(p, sep), sep) {
		if part == "" {
			continue
		}
		if acc != "" && !strings.HasSuffix(acc, sep) {
			acc += sep
		}
		acc += part
		out = append(out, pathSegment{part, acc})
	}
	return out
}

// PathBar displays clickable segments and collapses overflowing middle segments.
func PathBar(c *ui.Context, p string, opts PathBarOptions) PathBarResult {
	sep := opts.Separator
	if sep == "" {
		sep = "/"
	}
	segs := pathSegments(p, sep, opts.RootLabel)
	if len(segs) == 0 {
		panic("mujica: PathBar needs a path")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var res PathBarResult
	e := ui.Row(c).Grow(1).MinWidth(0).Gap(u).AlignItems(ui.Center).Role(ui.RoleToolbar).Label(core.Msg(c, "path_bar.label")).
		FocusGroup(ui.Horizontal).Disabled(opts.Disabled)
	res.Element = e
	was := e.Bounds()
	internal.Remeasure(c, e, was.W, was.H)
	hidden := pathBarFold(c, segs, was.W)
	disabled := e.IsDisabled()
	open := ui.Local(e, "path_bar.open", func() bool { return false })
	*open = *open && hidden > 0 && !disabled
	goTo := func(s pathSegment) {
		if !disabled {
			res.Navigated, *open = s.path, false
			slog.Debug("mujica: path_bar navigate", "depth", strings.Count(s.path, sep))
		}
	}
	e.Children(func() {
		for i, s := range segs {
			if i > 0 && (i <= 1 || i > hidden) {
				ui.Icon(c, icons.Must("chevron-right")).Size(12, 12).Shrink(0).TextColor(k.TextMuted)
			}
			if i >= 1 && i <= hidden {
				if i == 1 {
					more := pathBarCrumb(c, "…", false).Label(core.Msg(c, "path_bar.more"))
					if more.Clicked() && !disabled {
						*open = !*open
					}
					overlay.Popover(c, more, open, overlay.PopoverOptions{Label: core.Msg(c, "path_bar.more")}, func() {
						for _, h := range segs[1 : hidden+1] {
							if pathBarCrumb(c, h.label, false).Key(h.path).Clicked() {
								goTo(h)
							}
						}
					})
				}
				continue
			}
			if pathBarCrumb(c, s.label, i == len(segs)-1).Key(s.path).Clicked() {
				goTo(s)
			}
		}
	})
	return res
}

// pathBarFold counts middle segments to collapse into the available width.
func pathBarFold(c *ui.Context, segs []pathSegment, room float32) int {
	if room <= 0 {
		return 0
	}
	u := core.Density(c).Unit()
	width := func(s string) float32 {
		w, _ := c.MeasureText(0, ui.Span{Text: s, Size: core.FontSize(c, theme.BodySize)})
		return max(w+4*u, theme.MinHitSize) + 12 + 2*u
	}
	total := float32(0)
	for _, s := range segs {
		total += width(s.label)
	}
	hidden := 0
	for total > room && hidden < len(segs)-2 {
		hidden++
		total -= width(segs[hidden].label)
		if hidden == 1 {
			total += width("…")
		}
	}
	return hidden
}

// pathBarCrumb builds a button, emphasizing the current segment.
func pathBarCrumb(c *ui.Context, label string, current bool) *ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).MinWidth(theme.MinHitSize).Padding(0, d.Unit()*2).
		Radius(theme.ControlRadius).Shrink(0).TextColor(k.TextMuted)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case b.Pressed():
		b.Background(k.Selection).TextColor(k.Text)
	case b.Hovered():
		b.Background(k.SurfaceHover).TextColor(k.Text)
	}
	if current {
		b.TextColor(k.Text).FontWeight(600).Description(core.Msg(c, "path_bar.current"))
	}
	return b.Children(func() { ui.Text(c, label).SingleLine() })
}
