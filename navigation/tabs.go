package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Tab is one tab of Tabs.
type Tab struct {
	// Label is required; it names the tab and its panel.
	Label    string
	Icon     *ui.SVG
	Disabled bool
	// Panel builds the tab's content while it is chosen; nil shows none.
	Panel func()
}

// TabsOptions configure Tabs.
type TabsOptions struct {
	// Vertical stacks the tabs at the start with the panel beside them.
	Vertical bool
	// Label names the tab list for assistive technology.
	Label string
}

// Tabs shows tabs choosing *selected and the chosen tab's panel. Disabled tabs are skipped by the arrows.
func Tabs(c *ui.Context, selected *int, tabs []Tab, opts TabsOptions) core.ChoiceResult {
	if len(tabs) == 0 {
		panic("mujica: Tabs needs at least one tab")
	}
	if *selected < 0 || *selected >= len(tabs) {
		panic("mujica: Tabs selected index out of range")
	}
	if tabs[*selected].Disabled {
		panic("mujica: Tabs selected tab is disabled")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	pad := theme.FocusGap + theme.FocusWidth
	// TabsBase only sees enabled tabs, so the arrows skip disabled ones.
	var enabled []int
	at := 0
	for i, t := range tabs {
		if !t.Disabled {
			if i == *selected {
				at = len(enabled)
			}
			enabled = append(enabled, i)
		}
	}
	root := ui.Column(c).AlignItems(ui.Stretch)
	if root.IsDisabled() {
		for _, key := range []ui.Key{ui.KeyRight, ui.KeyDown, ui.KeyLeft, ui.KeyUp, ui.KeyHome, ui.KeyEnd} {
			root.Shortcut(0, key)
		}
	}
	if opts.Vertical {
		root.Row()
	}
	res := core.ChoiceResult{Element: root}
	root.Children(func() {
		var strip ui.Element
		if opts.Vertical {
			strip = ui.Scroll(c).Shrink(0).Padding(pad)
		} else {
			strip = ui.ScrollHorizontal(c).Shrink(0).MinWidth(0).Padding(pad, pad, 0, pad)
		}
		strip.Draw(func(p *ui.Painter, r ui.Rect) {
			if opts.Vertical {
				p.Fill(ui.Rect{X: r.X + r.W - 1, Y: r.Y, W: 1, H: r.H}, k.Border, 0)
			} else {
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, k.Border, 0)
			}
		})
		strip.Children(func() {
			base := ui.TabsBase(c, &at, len(enabled))
			list := base.List.Gap(u)
			if opts.Vertical {
				list.Column().AlignItems(ui.Stretch)
			}
			if opts.Label != "" {
				list.Label(opts.Label)
			}
			built := make([]ui.Element, len(tabs))
			list.Children(func() {
				n := 0
				for i, t := range tabs {
					var tab ui.Element
					if t.Disabled {
						tab = ui.Row(c).Shrink(0).Role(ui.RoleTab).Disabled(true)
					} else {
						tab = base.Tab(n)
						n++
					}
					TabsStyle(c, tab, t, i == enabled[at], opts.Vertical)
					built[i] = tab
				}
			})
			if !list.IsDisabled() && list.Changed() {
				prev := *selected
				*selected = enabled[at]
				res = core.NewChoiceResult(res.Element, true)
				built[*selected].ScrollIntoView()
				slog.Debug("mujica: tabs change", "from", prev, "to", *selected)
			}
		})
		panel := ui.Column(c).MinWidth(0).Padding(u * 4).Role(ui.RoleGroup).Label(tabs[*selected].Label)
		if opts.Vertical {
			panel.Grow(1)
		}
		if p := tabs[*selected].Panel; p != nil {
			panel.Children(p)
		}
	})
	return res
}

func TabsStyle(c *ui.Context, tab ui.Element, t Tab, on, vertical bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	tab.MinHeight(d.ControlHeight()).Padding(0, u*4).Gap(u * 2).AlignItems(ui.Center).Cursor(ui.CursorPointer)
	if vertical {
		tab.Radius(theme.ControlRadius, 0, 0, theme.ControlRadius)
	} else {
		tab.Radius(theme.ControlRadius, theme.ControlRadius, 0, 0)
	}
	if t.Disabled {
		tab.Cursor(ui.CursorNotAllowed)
	}
	if t.Label == "" {
		panic("mujica: Tab needs a Label")
	}
	fg := k.TextMuted
	switch {
	case on:
		tab.Background(k.Selection).FontWeight(600)
		fg = k.Text
	case !t.Disabled && tab.Pressed():
		tab.Background(k.Selection)
		fg = k.Text
	case !t.Disabled && tab.Hovered():
		tab.Background(k.SurfaceHover)
		fg = k.Text
	}
	internal.Ease(tab, core.Motion(c, theme.StateDuration), fg)
	tab.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if on {
			if vertical {
				p.Fill(ui.Rect{X: r.X + r.W - 2, Y: r.Y, W: 2, H: r.H}, k.AccentText, 0)
			} else {
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 2, W: r.W, H: 2}, k.AccentText, 0)
			}
		}
		if tab.FocusVisible() {
			internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
		}
	})
	tab.Children(func() {
		if t.Icon != nil {
			ui.Icon(c, t.Icon).Size(16, 16)
		}
		ui.Text(c, t.Label).SingleLine()
	})
}
