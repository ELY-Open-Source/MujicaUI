package mujica

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// SentimentGaugeOptions configure a SentimentGauge.
type SentimentGaugeOptions struct {
	// Label names the index, such as "Fear & Greed".
	Label string
	// Size is the dial width; zero is 240 DIPs.
	Size float32
}

var sentimentBands = [5]string{"sentiment_gauge.extreme_fear", "sentiment_gauge.fear", "sentiment_gauge.neutral",
	"sentiment_gauge.greed", "sentiment_gauge.extreme_greed"}

// sentimentBand is the band index of a 0 to 100 value.
func sentimentBand(v float64) int { return min(4, int(v/20)) }

// SentimentGauge is a half dial from extreme fear to extreme greed, with the value and its word.
func SentimentGauge(c *ui.Context, value float64, opts SentimentGaugeOptions) *ui.Element {
	if math.IsNaN(value) || value < 0 || value > 100 {
		panic("mujica: SentimentGauge value must be in 0..100")
	}
	if opts.Size < 0 {
		panic("mujica: SentimentGauge Size is negative")
	}
	w := opts.Size
	if w == 0 {
		w = 240
	}
	k := Tokens(c)
	u := Density(c).Unit()
	band := sentimentBand(value)
	word := msg(c, sentimentBands[band])
	label := fmt.Sprintf("%.0f, %s", value, word)
	if opts.Label != "" {
		label = opts.Label + ": " + label
	}
	// Fear to greed runs from the fall color to the rise color.
	fear, greed := trendColor(c, -1), trendColor(c, 1)
	cols := [5]ui.Color{fear, fear.Mix(k.Surface, 0.45), k.TextMuted, greed.Mix(k.Surface, 0.45), greed}
	size := FontSize(c, theme.CaptionSize)
	col := ui.Column(c).Gap(u).AlignItems(ui.Center).Role(ui.RoleMeter).Label(label)
	col.Children(func() {
		if opts.Label != "" {
			ui.Text(c, opts.Label).Font(theme.SerifFont).TextColor(k.TextMuted)
		}
		ui.Box(c).Size(w, w/2+8).Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.W/2
			outer, inner := r.W/2-2, r.W/2-18
			for i := range 5 {
				a0 := -math.Pi/2 + float64(i)*math.Pi/5
				p.FillPath(internal.Sector(cx, cy, inner, outer, a0+0.012, a0+math.Pi/5-0.012), cols[i])
			}
			a := -math.Pi/2 + value/100*math.Pi
			tx, ty := internal.Polar(cx, cy, inner-6, a)
			var needle ui.Path
			needle.MoveTo(internal.Snap(p, cx), internal.Snap(p, cy)).LineTo(internal.Snap(p, tx), internal.Snap(p, ty))
			p.StrokePath(&needle, 2.5, k.Text)
			p.FillPath(internal.Diamond(internal.Snap(p, cx), internal.Snap(p, cy), 10), k.Text)
			chartText(p, "0", r.X+2, cy+4, size, k.TextMuted, 0, 0)
			chartText(p, "100", r.X+r.W-2, cy+4, size, k.TextMuted, 1, 0)
		})
		ui.Text(c, fmt.Sprintf("%.0f", value)).Font(theme.SerifFont).FontSize(FontSize(c, theme.H1Size)).TextColor(k.Text).FontFeatures("tnum")
		ui.Text(c, word).FontWeight(600).TextColor(cols[band].Mix(k.Text, 0.25))
	})
	return col
}
