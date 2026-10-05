package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// IndicatorSelectorOptions configure an IndicatorSelector.
type IndicatorSelectorOptions struct {
	Disabled bool
}

// IndicatorSelector searches, adds, removes, and configures chart indicators.
func IndicatorSelector(c *ui.Context, overlays *[]IndicatorOverlay, panes *IndicatorPanes, opts IndicatorSelectorOptions) ChoiceResult {
	if overlays == nil || panes == nil {
		panic("mujica: IndicatorSelector needs overlays and panes")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	changed := false
	box := ui.Row(c).Disabled(opts.Disabled)
	box.Children(func() {
		open := ui.Local(box, "indicator_selector.open", func() bool { return false })
		query := ui.Local(box, "indicator_selector.query", func() string { return "" })
		if box.IsDisabled() {
			*open = false
		}
		label := fmt.Sprintf(msg(c, "indicator_selector.button"), len(*overlays)+len(panes.Panes))
		trigger := Button(c, label, ButtonOptions{Variant: Secondary, Icon: icons.Must("plus")})
		if trigger.Clicked() {
			*open = !*open
		}
		Popover(c, trigger, open, PopoverOptions{Title: msg(c, "indicator_selector.title"), Width: 340}, func() {
			TextInput(c, query, TextInputOptions{Placeholder: msg(c, "indicator_selector.search"), Label: msg(c, "indicator_selector.search"), Clearable: true})
			ui.Text(c, msg(c, "indicator_selector.available")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Margin(u, 0, 0, 0)
			found := 0
			for kind := IndicatorMA; kind <= IndicatorOBV; kind++ {
				name := kind.String()
				if q := strings.TrimSpace(*query); q != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(q)) {
					continue
				}
				found++
				ui.Row(c).Key(name).AlignItems(ui.Center).Gap(u * 2).MinHeight(Density(c).ControlHeight()).Children(func() {
					ui.Text(c, name).FontWeight(600).Width(56)
					where := "indicator_selector.pane"
					if indicatorInfo[kind].overlay {
						where = "indicator_selector.overlay"
					}
					ui.Text(c, msg(c, where)).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Grow(1)
					if Button(c, msg(c, "indicator_selector.add"), ButtonOptions{Variant: Ghost, Label: fmt.Sprintf(msg(c, "indicator_selector.add_named"), name)}).Clicked() {
						if indicatorInfo[kind].overlay {
							*overlays = append(slices.Clone(*overlays), IndicatorOverlay{Kind: kind, Params: kind.DefaultParams()})
						} else {
							panes.Panes = append(slices.Clone(panes.Panes), IndicatorPaneSpec{Kind: kind, Params: kind.DefaultParams()})
						}
						changed = true
						slog.Debug("mujica: indicator_selector add", "kind", name)
					}
				})
			}
			if found == 0 {
				ui.Text(c, msg(c, "indicator_selector.none")).TextColor(k.TextMuted)
			}
			if len(*overlays)+len(panes.Panes) == 0 {
				return
			}
			ui.Text(c, msg(c, "indicator_selector.active")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Margin(u*2, 0, 0, 0)
			remove := -1
			for i := range *overlays {
				if indicatorSelectorActive(c, u, fmt.Sprint("o", i), (*overlays)[i].Kind, &(*overlays)[i].Params, &changed) {
					remove = i
				}
			}
			if remove >= 0 {
				*overlays = slices.Delete(slices.Clone(*overlays), remove, remove+1)
				changed = true
			}
			remove = -1
			for i := range panes.Panes {
				if indicatorSelectorActive(c, u, fmt.Sprint("p", i), panes.Panes[i].Kind, &panes.Panes[i].Params, &changed) {
					remove = i
				}
			}
			if remove >= 0 {
				panes.Panes = slices.Delete(slices.Clone(panes.Panes), remove, remove+1)
				changed = true
			}
		})
	})
	return ChoiceResult{box, changed}
}

// indicatorSelectorActive edits parameters and reports removal of an indicator.
func indicatorSelectorActive(c *ui.Context, u float32, key string, kind IndicatorKind, params *[]float64, changed *bool) bool {
	removed := false
	ui.Row(c).Key(key).AlignItems(ui.Center).Gap(u * 2).Wrap().Children(func() {
		ui.Text(c, kind.String()).FontWeight(600).Width(56)
		if len(kind.DefaultParams()) > 0 && IndicatorParamsField(c, kind, params).Changed() {
			*changed = true
		}
		ui.Spacer(c)
		if iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "indicator_pane.remove"), indicatorName(kind, *params))).Clicked() {
			removed = true
			slog.Debug("mujica: indicator_selector remove", "kind", kind)
		}
	})
	return removed
}
