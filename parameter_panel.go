package mujica

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// ModelParam is one generation parameter, as temperature or top-p.
type ModelParam struct {
	// Label and Value are required.
	Label string
	Value *float64
	// Min, Max and Step bound the value; Default lies within them.
	Min, Max, Step, Default float64
}

// ParameterPanel adjusts generation parameters with sliders; each shows its range and resets to its default.
func ParameterPanel(c *ui.Context, params []ModelParam) ChoiceResult {
	for _, p := range params {
		if p.Label == "" || p.Value == nil || !(p.Min < p.Max) || p.Step <= 0 || p.Default < p.Min || p.Default > p.Max {
			panic("mujica: ModelParam needs Label, Value, Min < Max, Step > 0 and Default in range")
		}
	}
	k := Tokens(c)
	u := Density(c).Unit()
	small := FontSize(c, theme.CaptionSize)
	var r ChoiceResult
	r.Element = ui.Column(c).Gap(u * 4).MinWidth(240).Role(ui.RoleGroup).Label(msg(c, "parameter_panel.label"))
	disabled := r.Element.IsDisabled()
	set := func(p ModelParam, v float64) {
		if *p.Value != v && !disabled {
			*p.Value, r.changed = v, true
			slog.Debug("mujica: parameter_panel set", "param", p.Label, "value", v)
		}
	}
	r.Element.Children(func() {
		for _, p := range params {
			format := func(v float64) string {
				return strconv.FormatFloat(v, 'f', sliderDecimals(p.Min, p.Step, p.Default), 64)
			}
			ui.Column(c).Key(p.Label).Gap(u).Children(func() {
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					ui.Text(c, p.Label).FontWeight(600).Grow(1).SingleLine()
					ui.Text(c, format(*p.Value)).FontFeatures("tnum").TextColor(k.AccentText).SingleLine()
					ui.Row(c).Disabled(*p.Value == p.Default).Children(func() {
						if iconAction(c, icons.Must("rotate-ccw"), fmt.Sprintf(msg(c, "parameter_panel.reset"), p.Label)).
							Tooltip(fmt.Sprintf(msg(c, "parameter_panel.default"), format(p.Default))).Clicked() {
							set(p, p.Default)
						}
					})
				})
				v := *p.Value
				if Slider(c, &v, SliderOptions{Min: p.Min, Max: p.Max, Step: p.Step, Label: p.Label, Disabled: disabled}).Changed() {
					set(p, v)
				}
				ui.Row(c).Children(func() {
					ui.Text(c, format(p.Min)).FontSize(small).TextColor(k.TextMuted).Grow(1)
					ui.Text(c, fmt.Sprintf(msg(c, "parameter_panel.default"), format(p.Default))).FontSize(small).TextColor(k.TextMuted)
					ui.Text(c, format(p.Max)).FontSize(small).TextColor(k.TextMuted).Grow(1).TextAlign(ui.End)
				})
			})
		}
		ui.Row(c).Justify(ui.End).Children(func() {
			if Button(c, msg(c, "parameter_panel.reset_all"), ButtonOptions{Variant: Ghost, Icon: icons.Must("rotate-ccw")}).Clicked() {
				for _, p := range params {
					set(p, p.Default)
				}
			}
		})
	})
	return r
}
