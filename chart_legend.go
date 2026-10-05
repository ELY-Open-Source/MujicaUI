package mujica

import (
	"fmt"
	"log/slog"
	"maps"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChartLegendItem is a series key: name, color and marker shape.
type ChartLegendItem struct {
	Name  string
	Color ui.Color
	// Shape picks the marker: 0 circle, 1 square, 2 diamond, 3 triangle, 4 inverted.
	Shape int
}

// ChartLegendOptions configure a ChartLegend.
type ChartLegendOptions struct {
	Vertical bool
	// MaxItems folds the rest behind a "+N more" button; 0 shows all.
	MaxItems int
}

// ChartLegendResult is a built legend and the series switched this frame.
type ChartLegendResult struct {
	Element *ui.Element
	toggled string
}

// Toggled returns the series switched in this frame.
func (r ChartLegendResult) Toggled() (string, bool) { return r.toggled, r.toggled != "" }

// ChartLegend lists series keys; with hidden set, each is a check box hiding its series.
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
	k := Tokens(c)
	u := Density(c).Unit()
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
			ui.Row(c).Key(it.Name).Children(func() {
				if name, ok := chartLegendEntry(c, hidden, it, res.Element.IsDisabled()); ok {
					res.toggled = name
				}
			})
		}
		if opts.MaxItems > 0 && len(items) > opts.MaxItems {
			label := fmt.Sprintf(msg(c, "chart_legend.more"), len(items)-opts.MaxItems)
			if *open {
				label = msg(c, "chart_legend.fewer")
			}
			b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Padding(0, u*2).Radius(theme.ControlRadius).Disabled(res.Element.IsDisabled()).TextColor(k.AccentText)
			internal.Ring(b, k.Focus, theme.ControlRadius)
			if b.Hovered() {
				b.Background(k.SurfaceHover)
			}
			b.Children(func() { ui.Text(c, label).FontSize(FontSize(c, theme.CaptionSize)).SingleLine() })
			if b.Clicked() && !b.IsDisabled() {
				*open = !*open
				slog.Debug("mujica: chart_legend fold", "open", *open)
			}
		}
	})
	return res
}

// chartLegendEntry builds one key; a check box when hidden is set. It reports a switch.
func chartLegendEntry(c *ui.Context, hidden *map[string]bool, it ChartLegendItem, disabled bool) (string, bool) {
	k := Tokens(c)
	u := Density(c).Unit()
	on := hidden == nil || !(*hidden)[it.Name]
	var e *ui.Element
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
	e.MinHeight(theme.MinHitSize).Gap(u*2).AlignItems(ui.Center).Padding(0, u).Radius(theme.ControlRadius)
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
		t := ui.Text(c, it.Name).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Text)
		if !on {
			t.Strikethrough().TextColor(k.TextMuted)
		}
	})
	return it.Name, switched
}
