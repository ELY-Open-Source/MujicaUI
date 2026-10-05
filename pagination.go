package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// PaginationOptions configure a Pagination.
type PaginationOptions struct {
	// Total is the number of records; it must not be negative.
	Total int
	// PageSizes lists the page sizes offered; nil hides the choice.
	PageSizes []int
	// Jump shows a field going to a typed page.
	Jump bool
	// Siblings is how many pages show each side of the current; zero means 1.
	Siblings int
	Disabled bool
}

// Pagination pages *page (from 1) through Total records, clamping it when Total shrinks.
func Pagination(c *ui.Context, page, pageSize *int, opts PaginationOptions) ChoiceResult {
	if opts.Total < 0 {
		panic("mujica: Pagination Total is negative")
	}
	if *pageSize <= 0 {
		panic("mujica: Pagination pageSize must be positive")
	}
	if opts.PageSizes != nil && !slices.Contains(opts.PageSizes, *pageSize) {
		panic("mujica: Pagination pageSize is not one of PageSizes")
	}
	sib := opts.Siblings
	if sib == 0 {
		sib = 1
	}
	if sib < 0 {
		panic("mujica: Pagination Siblings is negative")
	}
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	row := ui.Row(c).AlignItems(ui.Center).Gap(u * 3).Wrap().Role(ui.RoleGroup).Label(msg(c, "pagination.label")).Disabled(opts.Disabled)
	res := ChoiceResult{Element: row}
	pages := paginationPages(opts.Total, *pageSize)
	if p := max(1, min(*page, pages)); p != *page {
		slog.Debug("mujica: pagination correct", "from", *page, "to", p, "pages", pages)
		*page, res.changed = p, true
	}
	set := func(p int) {
		if p != *page {
			slog.Debug("mujica: pagination page", "from", *page, "to", p)
			*page, res.changed = p, true
		}
	}
	row.Children(func() {
		ui.Text(c, fmt.Sprintf(msg(c, "pagination.total"), opts.Total)).TextColor(k.TextMuted).SingleLine()
		ui.Row(c).AlignItems(ui.Center).Gap(u).FocusGroup(ui.Horizontal).Children(func() {
			if paginationButton(c, "", icons.Must("chevron-left"), msg(c, "pagination.prev"), false, *page <= 1).Clicked() {
				set(*page - 1)
			}
			for _, p := range paginationItems(*page, pages, sib) {
				if p == 0 {
					ui.Text(c, "…").TextColor(k.TextMuted).Padding(0, u).Role(ui.RoleNone)
					continue
				}
				n := strconv.Itoa(p)
				name := fmt.Sprintf(msg(c, "pagination.page"), p)
				if p == *page {
					name += ", " + msg(c, "pagination.current")
				}
				if paginationButton(c, n, nil, name, p == *page, false).Clicked() {
					set(p)
				}
			}
			if paginationButton(c, "", icons.Must("chevron-right"), msg(c, "pagination.next"), false, *page >= pages).Clicked() {
				set(*page + 1)
			}
		})
		if opts.PageSizes != nil {
			paginationSizes(c, page, pageSize, opts.PageSizes, &res)
		}
		if opts.Jump {
			paginationJump(c, pages, set)
		}
	})
	return res
}

func paginationPages(total, size int) int { return max(1, (total+size-1)/size) }

// paginationItems lists the page numbers to show, 0 standing for an ellipsis.
func paginationItems(cur, pages, sib int) []int {
	var out []int
	if pages <= 2*sib+5 {
		for p := 1; p <= pages; p++ {
			out = append(out, p)
		}
		return out
	}
	start := max(min(cur-sib, pages-2*sib-2), 3)
	end := min(max(cur+sib, 2*sib+3), pages-2)
	out = append(out, 1)
	if start > 3 {
		out = append(out, 0)
	} else {
		out = append(out, 2)
	}
	for p := start; p <= end; p++ {
		out = append(out, p)
	}
	if end < pages-2 {
		out = append(out, 0)
	} else {
		out = append(out, pages-1)
	}
	return append(out, pages)
}

func paginationButton(c *ui.Context, text string, icon *ui.SVG, name string, on, disabled bool) *ui.Element {
	k := Tokens(c)
	h := Density(c).ControlHeight()
	b := ui.ButtonBase(c).MinHeight(h).MinWidth(h).Padding(0, Density(c).Unit()*2).Radius(theme.ControlRadius).
		Label(name).Disabled(disabled)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	fg := k.Text
	switch {
	case on:
		b.Background(k.Accent).Border(theme.BorderWidth, k.Accent).FontWeight(600)
		fg = k.OnAccent
	case b.Pressed():
		b.Background(k.Selection).Border(theme.BorderWidth, k.ControlBorder)
	case b.Hovered():
		b.Background(k.SurfaceHover).Border(theme.BorderWidth, k.ControlBorder)
	default:
		b.Border(theme.BorderWidth, k.Border)
	}
	internal.Ease(b, Motion(c, theme.StateDuration), fg)
	if icon != nil {
		b.Tooltip(name)
	}
	b.Children(func() {
		if icon != nil {
			ui.Icon(c, icon).Size(16, 16)
		} else {
			ui.Text(c, text).SingleLine().FontFeatures("tnum")
		}
	})
	return b
}

func paginationSizes(c *ui.Context, page, pageSize *int, sizes []int, res *ChoiceResult) {
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	old := *pageSize
	sel := ui.SelectBase(c, pageSize)
	sel.Trigger.MinHeight(d.ControlHeight()).Padding(0, u*3).Gap(u*2).Radius(theme.ControlRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).Label(msg(c, "pagination.size"))
	internal.Ring(sel.Trigger, k.Focus, theme.ControlRadius)
	if sel.Trigger.Hovered() {
		sel.Trigger.Background(k.SurfaceHover)
	}
	internal.Ease(sel.Trigger, Motion(c, theme.StateDuration), k.Text)
	sel.Trigger.Children(func() {
		ui.Text(c, fmt.Sprintf(msg(c, "pagination.per_page"), *pageSize)).SingleLine()
		ui.Icon(c, icons.Must("chevron-down")).Size(14, 14).TextColor(k.TextMuted)
	})
	sel.Popup(func(panel *ui.Element) {
		internal.PopIn(c, panel, Motion(c, theme.PopupDuration))
		panel.Margin(u, 0, 0, 0).Padding(u).Radius(theme.PopupRadius).Background(k.Surface).
			Border(theme.BorderWidth, k.Border).Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, 0.16))
		for _, s := range sizes {
			it := sel.Item(s).Padding(u*1.5, u*3).Radius(theme.ControlRadius).Gap(u * 2)
			if it.Highlighted() {
				it.Background(k.Selection)
			}
			internal.Ease(it, Motion(c, theme.StateDuration), k.Text)
			it.Children(func() {
				ui.Text(c, fmt.Sprintf(msg(c, "pagination.per_page"), s)).SingleLine().Grow(1)
				if s == *pageSize {
					ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(k.AccentText)
				}
			})
		}
	})
	if *pageSize != old {
		// Keep the first record of the old page in view.
		first := (*page-1)*old + 1
		*page = (first-1) / *pageSize + 1
		res.changed = true
		slog.Debug("mujica: pagination size", "from", old, "to", *pageSize, "page", *page)
	}
}

func paginationJump(c *ui.Context, pages int, set func(int)) {
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	ui.Row(c).AlignItems(ui.Center).Gap(u * 2).Children(func() {
		ui.Text(c, msg(c, "pagination.jump")).TextColor(k.TextMuted).SingleLine()
		box := ui.Row(c).AlignItems(ui.Center).MinHeight(d.ControlHeight()).Width(64).Padding(0, u*2).
			Radius(theme.ControlRadius).Background(k.Surface)
		text := ui.Local(box, "text", func() string { return "" })
		bad := ui.Local(box, "bad", func() bool { return false })
		var in *ui.Element
		box.Children(func() {
			in = ui.TextInputBase(c, text).FillWidth().Label(msg(c, "pagination.field"))
		})
		if in.Changed() {
			*bad = false
		}
		if in.Submitted() {
			p, err := strconv.Atoi(strings.TrimSpace(*text))
			if err != nil || p < 1 || p > pages {
				*bad = true
				slog.Debug("mujica: pagination jump rejected", "text", *text, "pages", pages)
			} else {
				*text, *bad = "", false
				set(p)
			}
		}
		border := k.ControlBorder
		if *bad {
			border = k.Danger
			in.Error(fmt.Sprintf(msg(c, "pagination.range"), pages))
		} else if in.Focused() {
			border = k.Focus
		}
		internal.Ease(box.Border(theme.BorderWidth, border), Motion(c, theme.StateDuration), k.Text)
		if *bad {
			ui.Text(c, fmt.Sprintf(msg(c, "pagination.range"), pages)).TextColor(k.Danger).SingleLine()
		}
	})
}
