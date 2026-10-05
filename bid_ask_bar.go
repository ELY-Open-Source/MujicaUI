package mujica

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// BidAskBarOptions configure a BidAskBar.
type BidAskBarOptions struct {
	// Decimals of the volumes shown; zero shows whole numbers.
	Decimals int
	// Label names the bar; empty uses the built-in name.
	Label string
}

// BidAskBar compares buy and sell volume in a split bar.
func BidAskBar(c *ui.Context, bid, ask float64, opts BidAskBarOptions) *ui.Element {
	if bid < 0 || ask < 0 || math.IsNaN(bid+ask) {
		panic("mujica: BidAskBar volumes must be zero or more")
	}
	if opts.Decimals < 0 {
		panic("mujica: BidAskBar with negative Decimals")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	total := bid + ask
	share := 0.5
	if total > 0 {
		share = bid / total
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "bid_ask_bar.label")
	}
	buy := fmt.Sprintf(msg(c, "bid_ask_bar.buy"), statisticFormat(share*100, 1), statisticFormat(bid, opts.Decimals))
	sell := fmt.Sprintf(msg(c, "bid_ask_bar.sell"), statisticFormat(ask, opts.Decimals), statisticFormat((1-share)*100, 1))
	up, down := trendColor(c, 1), trendColor(c, -1)
	upText, downText := up, down
	if total == 0 {
		up, down, upText, downText = k.Border, k.Border, k.TextMuted, k.TextMuted
	}
	e := ui.Column(c).Gap(u).MinWidth(160).Role(ui.RoleMeter).Label(label).Description(buy + ", " + sell)
	e.Children(func() {
		ui.Row(c).Justify(ui.SpaceBetween).Gap(u * 2).Children(func() {
			ui.Text(c, buy).TextColor(upText).FontSize(FontSize(c, theme.CaptionSize)).FontFeatures("tnum").FontWeight(600).SingleLine()
			ui.Text(c, sell).TextColor(downText).FontSize(FontSize(c, theme.CaptionSize)).FontFeatures("tnum").FontWeight(600).SingleLine()
		})
		ui.Box(c).Height(8).FillWidth().Draw(func(p *ui.Painter, r ui.Rect) {
			split := r.X + r.W*float32(share)
			const gap = 2
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(0, split-r.X-gap/2), H: r.H}, up, theme.ControlRadius)
			p.Fill(ui.Rect{X: split + gap/2, Y: r.Y, W: max(0, r.X+r.W-split-gap/2), H: r.H}, down, theme.ControlRadius)
		})
	})
	return e
}
