package mujica

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// StatCardOptions configure a StatCard. Label is required.
type StatCardOptions struct {
	Label string
	// Change is the relative change, 0.12 for +12%; HasChange shows it.
	Change    float64
	HasChange bool
	// LowerIsBetter colors a fall as good, as for latency.
	LowerIsBetter bool
	// Trend draws a sparkline; nil hides it.
	Trend  []float64
	Status DataStatus
	Error  string
}

// StatCardResult reports this frame's events.
type StatCardResult struct {
	Element *ui.Element
	Retry   bool
}

// StatCard shows a statistic, its change and trend, or loading and error states at the same size.
func StatCard(c *ui.Context, value string, opts StatCardOptions) StatCardResult {
	if opts.Label == "" {
		panic("mujica: StatCard needs a Label")
	}
	if math.IsNaN(opts.Change) || math.IsInf(opts.Change, 0) {
		panic("mujica: StatCard Change must be finite")
	}
	if opts.Status != DataReady && opts.Status != DataLoading && opts.Status != DataFailed {
		panic(fmt.Sprintf("mujica: StatCard Status %d is not ready, loading or failed", int(opts.Status)))
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r StatCardResult
	r.Element = agentCard(c).Height(116).MinWidth(180).Label(opts.Label)
	r.Element.Children(func() {
		agentCaption(c, opts.Label)
		switch opts.Status {
		case DataLoading:
			ui.Column(c).Grow(1).Justify(ui.Center).Children(func() {
				Spinner(c, SpinnerOptions{Label: msg(c, "list.loading"), ShowLabel: true})
			})
			return
		case DataFailed:
			ui.Column(c).Grow(1).Justify(ui.Center).Gap(u).Children(func() {
				ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).TextColor(k.Danger).Children(func() {
					ui.Icon(c, icons.Must("circle-alert")).Size(16, 16)
					ui.Text(c, dataError(c, opts.Error)).Role(ui.RoleStatus).FontSize(FontSize(c, theme.CaptionSize+1))
				})
				ui.Row(c).Children(func() { r.Retry = agentToggle(c, msg(c, "list.retry"), "refresh-cw", false) })
			})
			return
		}
		ui.Text(c, value).Font(theme.SerifFont).FontSize(FontSize(c, theme.H2Size)).FontWeight(600).FontFeatures("tnum").SingleLine()
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			if opts.HasChange {
				statChange(c, opts.Change, opts.LowerIsBetter)
			}
			ui.Box(c).Grow(1)
			if len(opts.Trend) > 0 {
				Sparkline(c, opts.Trend, SparklineOptions{Frame: ChartFrame{Label: opts.Label}, Width: 88})
			}
		})
	})
	if r.Retry {
		slog.Debug("mujica: stat_card retry", "label", opts.Label)
	}
	return r
}

// statChange shows a signed percentage with an arrow, toned good or bad.
func statChange(c *ui.Context, change float64, lowerIsBetter bool) {
	k := Tokens(c)
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
	ui.Row(c).Gap(Density(c).Unit()).AlignItems(ui.Center).Label(fmt.Sprintf(msg(c, "stat_card.change"), text)).Children(func() {
		ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(tone)
		ui.Text(c, text).TextColor(tone).FontSize(FontSize(c, theme.CaptionSize+1)).FontFeatures("tnum")
	})
}
