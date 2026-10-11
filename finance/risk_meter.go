package finance

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// RiskLevel is one step of a RiskMeter's scale.
type RiskLevel struct {
	Label       string
	Description string
}

// RiskMeterOptions configure a RiskMeter.
type RiskMeterOptions struct {
	// Levels rise in risk; nil uses five built-in levels.
	Levels []RiskLevel
	// Label names the meter; empty uses the built-in "Risk level".
	Label string
}

// RiskMeter shows level, from 1, on a named graded scale.
func RiskMeter(c *ui.Context, level int, opts RiskMeterOptions) *ui.Element {
	levels := opts.Levels
	if levels == nil {
		for i := range 5 {
			levels = append(levels, RiskLevel{Label: core.Msg(c, fmt.Sprintf("risk_meter.level%d", i+1))})
		}
	}
	if len(levels) < 2 || level < 1 || level > len(levels) {
		panic(fmt.Sprintf("mujica: RiskMeter level %d outside 1..%d", level, len(levels)))
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "risk_meter.label")
	}
	cur := levels[level-1]
	spoken := fmt.Sprintf(core.Msg(c, "risk_meter.value"), level, len(levels), cur.Label)
	col := riskMeterColor(k, level, len(levels))
	e := ui.Column(c).Gap(u * 2).MinWidth(200).Role(ui.RoleMeter).Label(label).Description(spoken + " " + cur.Description)
	e.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, cur.Label).TextColor(col).FontWeight(600).SingleLine()
			ui.Text(c, fmt.Sprintf("%d / %d", level, len(levels))).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum")
		})
		ui.Box(c).Height(18).FillWidth().Draw(func(p *ui.Painter, r ui.Rect) {
			n := len(levels)
			const gap = 3
			w := (r.W - gap*float32(n-1)) / float32(n)
			for i := range n {
				seg := ui.Rect{X: r.X + float32(i)*(w+gap), Y: r.Y + 8, W: w, H: 8}
				if i < level {
					p.Fill(seg, riskMeterColor(k, i+1, n), theme.ControlRadius)
				} else {
					p.Stroke(seg, k.ControlBorder, theme.ControlRadius, theme.BorderWidth)
				}
				if i == level-1 {
					p.FillPath(internal.Diamond(seg.X+seg.W/2, r.Y+3, 7), k.Text)
				}
			}
		})
		if cur.Description != "" {
			ui.Text(c, cur.Description).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize))
		}
	})
	return e
}

// riskMeterColor grades from success through warning to danger.
func riskMeterColor(k theme.Tokens, level, n int) ui.Color {
	switch f := float64(level-1) / float64(n-1); {
	case f < 0.34:
		return k.Success
	case f < 0.67:
		return k.Warning
	}
	return k.Danger
}
