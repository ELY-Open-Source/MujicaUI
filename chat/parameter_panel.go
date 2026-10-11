package chat

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ModelParam is one generation parameter, as temperature or top-p.
type ModelParam struct {
	// Label and Value are required.
	Label string
	Value *float64
	// Min, Max and Step define the grid; Default selects reset.
	Min, Max, Step, Default float64
}

// ParameterPanel displays bounded generation sliders and reset actions.
func ParameterPanel(c *ui.Context, params []ModelParam) core.ChoiceResult {
	for _, p := range params {
		if p.Label == "" || p.Value == nil || !(p.Min < p.Max) || !(p.Step > 0) || !(p.Default >= p.Min && p.Default <= p.Max) {
			panic("mujica: ModelParam needs Label, Value, Min < Max, Step > 0 and Default in range")
		}
		for _, v := range []float64{p.Min, p.Max, p.Step, p.Default, *p.Value} {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				panic("mujica: ModelParam values must be finite")
			}
		}
		if *p.Value < p.Min || *p.Value > p.Max {
			panic("mujica: ModelParam value must be in range")
		}
		spec := input.NumberInputResolve(input.NumberInputOptions{Min: p.Min, Max: p.Max, Step: p.Step, Precision: input.SliderDecimals(p.Min, p.Step)})
		if spec.StepBy(p.Default, 0) != p.Default {
			panic("mujica: ModelParam default must lie on the step grid")
		}
		if spec.StepBy(*p.Value, 0) != *p.Value {
			panic("mujica: ModelParam value must lie on the step grid")
		}
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	var r core.ChoiceResult
	r.Element = ui.Column(c).Gap(u * 4).MinWidth(240).Role(ui.RoleGroup).Label(core.Msg(c, "parameter_panel.label"))
	disabled := r.Element.IsDisabled()
	set := func(p ModelParam, v float64) {
		if *p.Value != v && !disabled {
			*p.Value, r = v, core.NewChoiceResult(r.Element, true)
			slog.Debug("mujica: parameter_panel set", "param", p.Label, "value", v)
		}
	}
	r.Element.Children(func() {
		for _, p := range params {
			format := func(v float64) string {
				return strconv.FormatFloat(v, 'f', input.SliderDecimals(p.Min, p.Step, p.Default), 64)
			}
			ui.Column(c).Key(p.Label).Gap(u).Children(func() {
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					ui.Text(c, p.Label).FontWeight(600).Grow(1).SingleLine()
					ui.Text(c, format(*p.Value)).FontFeatures("tnum").TextColor(k.AccentText).SingleLine()
					ui.Row(c).Disabled(*p.Value == p.Default).Children(func() {
						if core.IconAction(c, icons.Must("rotate-ccw"), fmt.Sprintf(core.Msg(c, "parameter_panel.reset"), p.Label)).
							Tooltip(fmt.Sprintf(core.Msg(c, "parameter_panel.default"), format(p.Default))).Clicked() {
							set(p, p.Default)
						}
					})
				})
				v := *p.Value
				if input.Slider(c, &v, input.SliderOptions{Min: p.Min, Max: p.Max, Step: p.Step, Label: p.Label, Disabled: disabled}).Changed() {
					set(p, v)
				}
				ui.Row(c).Children(func() {
					ui.Text(c, format(p.Min)).FontSize(small).TextColor(k.TextMuted).Grow(1)
					ui.Text(c, fmt.Sprintf(core.Msg(c, "parameter_panel.default"), format(p.Default))).FontSize(small).TextColor(k.TextMuted)
					ui.Text(c, format(p.Max)).FontSize(small).TextColor(k.TextMuted).Grow(1).TextAlign(ui.End)
				})
			})
		}
		ui.Row(c).Justify(ui.End).Children(func() {
			if input.Button(c, core.Msg(c, "parameter_panel.reset_all"), input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("rotate-ccw")}).Clicked() {
				for _, p := range params {
					set(p, p.Default)
				}
			}
		})
	})
	return r
}
