package chart

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/overlay"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// intervalGroups sorts Intervals under their headings.
var intervalGroups = []struct {
	key   string
	items []Interval
}{
	{"interval_selector.minutes", []Interval{"1m", "3m", "5m", "15m", "30m"}},
	{"interval_selector.hours", []Interval{"1h", "2h", "4h", "6h", "12h"}},
	{"interval_selector.days", []Interval{"1D", "3D", "1W", "1M"}},
}

// IntervalSelectorOptions configure an IntervalSelector.
type IntervalSelectorOptions struct {
	Disabled bool
}

// IntervalSelector groups candle periods and preserves caller-owned favorites.
func IntervalSelector(c *ui.Context, value *Interval, favorites *[]Interval, opts IntervalSelectorOptions) core.ChoiceResult {
	if value == nil || favorites == nil {
		panic("mujica: IntervalSelector needs a value and favorites")
	}
	for _, iv := range append([]Interval{*value}, *favorites...) {
		if !slices.Contains(Intervals, iv) {
			panic(fmt.Sprintf("mujica: unknown Interval %q", string(iv)))
		}
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	changed := false
	row := ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Disabled(opts.Disabled)
	disabled := row.IsDisabled()
	row.Children(func() {
		quick := slices.Clone(*favorites)
		if !slices.Contains(quick, *value) {
			quick = append(quick, *value)
		}
		slices.SortFunc(quick, func(a, b Interval) int { return slices.Index(Intervals, a) - slices.Index(Intervals, b) })
		items := make([]input.ToggleItem[Interval], len(quick))
		for i, iv := range quick {
			items[i] = input.ToggleItem[Interval]{Value: iv, Label: string(iv)}
		}
		sel := []Interval{*value}
		if input.ToggleGroup(c, &sel, items, input.ToggleGroupOptions{Label: core.Msg(c, "interval_selector.label")}).Changed() && !disabled {
			*value, changed = sel[0], true
		}
		trigger := input.Button(c, "", input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("chevron-down"), Label: core.Msg(c, "interval_selector.all")})
		open := ui.Local(row, "interval_selector.open", func() bool { return false })
		if disabled {
			*open = false
		}
		if trigger.Clicked() && !disabled {
			*open = !*open
		}
		overlay.Popover(c, trigger, open, overlay.PopoverOptions{Label: core.Msg(c, "interval_selector.all"), Width: 220}, func() {
			for _, g := range intervalGroups {
				ui.Text(c, core.Msg(c, g.key)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Margin(u, 0, 0, u)
				for _, iv := range g.items {
					if intervalSelectorRow(c, k, u, iv, iv == *value, slices.Contains(*favorites, iv), favorites) {
						*value, changed, *open = iv, true, false
					}
				}
			}
		})
	})
	if changed {
		slog.Debug("mujica: interval_selector changed", "interval", *value)
	}
	return core.NewChoiceResult(row, changed)
}

// intervalSelectorRow shows a selectable period and its favorite toggle.
func intervalSelectorRow(c *ui.Context, k theme.Tokens, u float32, iv Interval, on, fav bool, favorites *[]Interval) bool {
	chosen := false
	ui.Row(c.Key(string(iv))).AlignItems(ui.Center).Gap(u).Children(func() {
		b := ui.ButtonBase(c).Grow(1).MinHeight(core.Density(c).ControlHeight()).Padding(0, u*2).Justify(ui.Start).Radius(theme.ControlRadius).Gap(u * 2)
		internal.Ring(b, k.Focus, theme.ControlRadius)
		if on {
			b.Background(k.Selection).TextColor(k.AccentText)
		} else if b.Hovered() {
			b.Background(k.SurfaceHover)
		}
		b.Children(func() {
			ui.Box(c).Size(14, 14).Draw(func(p *ui.Painter, r ui.Rect) {
				if on {
					internal.Mark(p, r, "check", k.AccentText)
				}
			})
			ui.Text(c, string(iv)).SingleLine()
		})
		chosen = b.Clicked()
		name := fmt.Sprintf(core.Msg(c, "interval_selector.favorite"), iv)
		if fav {
			name = fmt.Sprintf(core.Msg(c, "interval_selector.unfavorite"), iv)
		}
		star := core.IconAction(c, icons.Must("star"), name)
		if fav {
			star.TextColor(k.Warning)
		}
		if star.Clicked() && !star.IsDisabled() {
			if fav {
				*favorites = slices.DeleteFunc(slices.Clone(*favorites), func(x Interval) bool { return x == iv })
			} else {
				*favorites = append(slices.Clone(*favorites), iv)
			}
			slog.Debug("mujica: interval_selector favorite", "interval", iv, "on", !fav)
		}
	})
	return chosen
}
