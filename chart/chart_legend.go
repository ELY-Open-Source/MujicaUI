package chart

import (
	"fmt"
	"log/slog"
	"maps"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ChartLegendItem is a series key: name, color and marker shape.
type ChartLegendItem struct {
	Name  string
	Color ui.Color
	// Shape selects circle, square, diamond, triangle, or inverted triangle.
	Shape int
}

// ChartLegendOptions configure a ChartLegend.
type ChartLegendOptions struct {
	Vertical bool
	// MaxItems folds overflow; zero shows all items.
	MaxItems int
}

// ChartLegendResult exposes the element and legend events.
type ChartLegendResult struct {
	Element          ui.Element
	toggled          string
	expanded, folded bool
}

// Toggled returns the series switched in this frame.
func (r ChartLegendResult) Toggled() (string, bool) { return r.toggled, r.toggled != "" }

// Expanded reports the new state when expansion changes.
func (r ChartLegendResult) Expanded() (bool, bool) { return r.expanded, r.folded }

// ChartLegend lists series; hidden enables visibility toggles.
func ChartLegend(c *ui.Context, hidden *map[string]bool, items []ChartLegendItem, opts ChartLegendOptions) ChartLegendResult {
	if opts.MaxItems < 0 {
		panic("mujica: ChartLegend MaxItems must be ≥ 0")
	}
	seen := map[string]bool{}
	for _, it := range items {
		if it.Name == "" || seen[it.Name] {
			panic("mujica: ChartLegend items need unique names")
		}
		seen[it.Name] = true
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var res ChartLegendResult
	res.Element = ui.Row(c).Gap(u * 3).Wrap()
	orient := ui.Horizontal
	if opts.Vertical {
		res.Element.Column().Gap(u).AlignItems(ui.Start)
		orient = ui.Vertical
	}
	open := ui.Local(res.Element, "mujica.legend.open", func() bool { return false })
	res.Element.FocusGroup(orient).Children(func() {
		shown := items
		if opts.MaxItems > 0 && len(items) > opts.MaxItems && !*open {
			shown = items[:opts.MaxItems]
		}
		for _, it := range shown {
			ui.Row(c.Key(it.Name)).Children(func() {
				if name, ok := chartLegendEntry(c, hidden, it, res.Element.IsDisabled()); ok {
					res.toggled = name
				}
			})
		}
		if opts.MaxItems > 0 && len(items) > opts.MaxItems {
			label := fmt.Sprintf(core.Msg(c, "chart_legend.more"), len(items)-opts.MaxItems)
			if *open {
				label = core.Msg(c, "chart_legend.fewer")
			}
			b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Padding(0, u*2).Radius(theme.ControlRadius).Disabled(res.Element.IsDisabled()).TextColor(k.AccentText)
			internal.Ring(b, k.Focus, theme.ControlRadius)
			if b.Hovered() {
				b.Background(k.SurfaceHover)
			}
			b.Children(func() { ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine() })
			if b.Clicked() && !b.IsDisabled() {
				*open = !*open
				res.expanded, res.folded = *open, true
				slog.Debug("mujica: chart_legend fold", "open", *open)
			}
		}
	})
	return res
}

// chartLegendEntry builds a series visibility control.
func chartLegendEntry(c *ui.Context, hidden *map[string]bool, it ChartLegendItem, disabled bool) (string, bool) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	on := hidden == nil || !(*hidden)[it.Name]
	var e ui.Element
	switched := false
	if hidden != nil {
		e = ui.CheckboxBase(c, &on).Disabled(disabled).Label(it.Name)
		internal.Ring(e, k.Focus, theme.ControlRadius)
		if e.Hovered() && !disabled {
			e.Background(k.SurfaceHover)
		}
		if e.Changed() && !e.IsDisabled() {
			m := maps.Clone(*hidden)
			if m == nil {
				m = map[string]bool{}
			}
			m[it.Name] = !on
			*hidden, switched = m, true
			slog.Debug("mujica: chart_legend toggle", "series", it.Name, "visible", on)
		}
	} else {
		e = ui.Row(c)
	}
	e.MinWidth(theme.MinHitSize).MinHeight(core.Density(c).ControlHeight()).Gap(u*2).AlignItems(ui.Center).Padding(0, u).Radius(theme.ControlRadius)
	col := it.Color
	if !on {
		col = k.TextMuted
	}
	e.Children(func() {
		ui.Box(c).Size(12, 12).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			if on {
				internal.Marker(p, it.Shape, r.X+6, r.Y+6, 10, col)
			} else {
				p.Stroke(r, col, 2, 1)
			}
		})
		t := ui.Text(c, it.Name).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Text)
		if !on {
			t.Strikethrough().TextColor(k.TextMuted)
		}
	})
	return it.Name, switched
}
