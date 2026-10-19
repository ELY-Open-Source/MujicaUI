package devtools

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/agent"
	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// StatCardOptions configure a StatCard. Label is required.
type StatCardOptions struct {
	Label string
	// Change is relative; 0.12 means +12%, shown with HasChange.
	Change    float64
	HasChange bool
	// LowerIsBetter colors a fall as good, as for latency.
	LowerIsBetter bool
	// Trend draws a sparkline; nil hides it.
	Trend  []float64
	Status data.DataStatus
	Error  string
}

// StatCard shows values, changes, trends, loading and errors.
func StatCard(c *ui.Context, value string, opts StatCardOptions) core.ChoiceResult {
	if opts.Label == "" {
		panic("mujica: StatCard needs a Label")
	}
	if math.IsNaN(opts.Change) || math.IsInf(opts.Change, 0) {
		panic("mujica: StatCard Change must be finite")
	}
	if opts.Status != data.DataReady && opts.Status != data.DataLoading && opts.Status != data.DataFailed {
		panic(fmt.Sprintf("mujica: StatCard Status %d is not ready, loading or failed", int(opts.Status)))
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = agent.AgentCard(c).MinHeight(core.FontSize(c, 116)).MinWidth(180).Label(opts.Label)
	r.Element.Children(func() {
		agent.AgentCaption(c, opts.Label)
		switch opts.Status {
		case data.DataLoading:
			ui.Column(c).Grow(1).Justify(ui.Center).Children(func() {
				core.Spinner(c, core.SpinnerOptions{Label: core.Msg(c, "list.loading"), ShowLabel: true})
			})
			return
		case data.DataFailed:
			ui.Column(c).Grow(1).Justify(ui.Center).Gap(u).Children(func() {
				ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).TextColor(k.Danger).Children(func() {
					ui.Icon(c, icons.Must("circle-alert")).Size(16, 16)
					ui.Text(c, data.DataError(c, opts.Error)).Role(ui.RoleStatus).FontSize(core.FontSize(c, theme.CaptionSize+1))
				})
				ui.Row(c).Children(func() {
					r = core.NewChoiceResult(r.Element, agent.AgentToggle(c, core.Msg(c, "list.retry"), "refresh-cw", false))
				})
			})
			return
		}
		ui.Text(c, value).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H2Size)).FontWeight(600).FontFeatures("tnum").SingleLine()
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			if opts.HasChange {
				statChange(c, opts.Change, opts.LowerIsBetter)
			}
			ui.Box(c).Grow(1)
			if len(opts.Trend) > 0 {
				chart.Sparkline(c, opts.Trend, chart.SparklineOptions{Frame: chart.ChartFrame{Label: opts.Label}, Width: 88})
			}
		})
	})
	if r.Changed() {
		slog.Debug("mujica: stat_card retry", "label", opts.Label)
	}
	return r
}

// statChange shows a signed percentage and direction arrow.
func statChange(c *ui.Context, change float64, lowerIsBetter bool) {
	k := core.Tokens(c)
	icon, tone, sign := "minus", k.TextMuted, ""
	switch {
	case change > 0:
		icon, tone, sign = "arrow-up", k.Success, "+"
	case change < 0:
		icon, tone = "arrow-down", k.Danger
	}
	if lowerIsBetter && change != 0 {
		if tone == k.Success {
			tone = k.Danger
		} else {
			tone = k.Success
		}
	}
	text := fmt.Sprintf("%s%.1f%%", sign, change*100)
	ui.Row(c).Gap(core.Density(c).Unit()).AlignItems(ui.Center).Label(fmt.Sprintf(core.Msg(c, "stat_card.change"), text)).Children(func() {
		ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(tone)
		ui.Text(c, text).TextColor(tone).FontSize(core.FontSize(c, theme.CaptionSize+1)).FontFeatures("tnum")
	})
}
