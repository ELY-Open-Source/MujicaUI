package finance

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MarginIndicatorOptions configure a MarginIndicator. Zero lines mean 80% and 100%.
type MarginIndicatorOptions struct {
	// Warning and Liquidation are margin ratios, Warning below Liquidation.
	Warning, Liquidation float64
	// Label names the indicator; empty uses the built-in "Margin ratio".
	Label string
}

// MarginIndicator shows maintenance over equity against its two lines.
func MarginIndicator(c *ui.Context, ratio float64, opts MarginIndicatorOptions) *ui.Element {
	warn, liq := opts.Warning, opts.Liquidation
	if warn == 0 && liq == 0 {
		warn, liq = 0.8, 1
	}
	if !(warn > 0 && warn < liq) || ratio < 0 || math.IsNaN(ratio) {
		panic(fmt.Sprintf("mujica: MarginIndicator ratio %v, lines %v and %v: need 0 < Warning < Liquidation", ratio, warn, liq))
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "margin_indicator.label")
	}
	status, col, icon := core.Msg(c, "margin_indicator.safe"), k.Success, "circle-check"
	switch {
	case ratio >= liq:
		status, col, icon = core.Msg(c, "margin_indicator.liquidation"), k.Danger, "circle-x"
	case ratio >= warn:
		status, col, icon = core.Msg(c, "margin_indicator.warning"), k.Warning, "triangle-alert"
	}
	pct := func(v float64) string { return data.StatisticFormat(v*100, 1) + "%" }
	top := math.Max(liq*1.25, ratio)
	e := ui.Column(c).Gap(u * 2).MinWidth(220).Role(ui.RoleMeter).Label(label).
		Description(fmt.Sprintf(core.Msg(c, "margin_indicator.spoken"), pct(ratio), status, pct(warn), pct(liq)))
	e.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, pct(ratio)).FontSize(core.FontSize(c, theme.H3Size)).FontFeatures("tnum").FontWeight(600).SingleLine()
			ui.Row(c).Gap(u / 2).AlignItems(ui.Center).TextColor(col).Children(func() {
				ui.Icon(c, icons.Must(icon)).Size(14, 14)
				ui.Text(c, status).FontSize(core.FontSize(c, theme.CaptionSize)).FontWeight(600).SingleLine()
			})
		})
		ui.Box(c).Height(14).FillWidth().Draw(func(p *ui.Painter, r ui.Rect) {
			bar := ui.Rect{X: r.X, Y: r.Y + 3, W: r.W, H: 8}
			p.Fill(bar, k.SurfaceHover, theme.ControlRadius)
			p.Fill(ui.Rect{X: bar.X, Y: bar.Y, W: bar.W * float32(ratio/top), H: bar.H}, col, theme.ControlRadius)
			// The liquidation line is drawn heavier than the warning line.
			for i, line := range []float64{warn, liq} {
				x := bar.X + bar.W*float32(line/top)
				p.Line(x, r.Y, x, r.Y+r.H, 1.5+float32(i)*1.5, k.Text)
			}
		})
		ui.Row(c).Gap(u * 4).Wrap().Children(func() {
			ui.Text(c, core.Msg(c, "margin_indicator.warning_line")+" "+pct(warn)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			ui.Text(c, core.Msg(c, "margin_indicator.liquidation_line")+" "+pct(liq)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		})
	})
	return e
}
