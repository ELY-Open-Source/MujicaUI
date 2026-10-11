package finance

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SentimentGaugeOptions configure a SentimentGauge.
type SentimentGaugeOptions struct {
	// Label names the index, such as "Fear & Greed".
	Label string
	// Size is the dial width; zero is 240 DIPs.
	Size   float32
	Status data.DataStatus
	Error  string
	Empty  bool
	Retry  func()
}

var sentimentBands = [5]string{"sentiment_gauge.extreme_fear", "sentiment_gauge.fear", "sentiment_gauge.neutral",
	"sentiment_gauge.greed", "sentiment_gauge.extreme_greed"}

func sentimentBand(v float64) int { return min(4, int(v/20)) }

func SentimentGauge(c *ui.Context, value float64, opts SentimentGaugeOptions) *ui.Element {
	if math.IsNaN(value) || value < 0 || value > 100 {
		panic("mujica: SentimentGauge value must be in 0..100")
	}
	if math.IsNaN(float64(opts.Size)) || math.IsInf(float64(opts.Size), 0) || opts.Size != 0 && opts.Size < 48 {
		panic("mujica: SentimentGauge Size is negative")
	}
	w := opts.Size
	if w == 0 {
		w = 240
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	dialHeight := max(w/2+8, u*30)
	band := sentimentBand(value)
	word := core.Msg(c, sentimentBands[band])
	label := fmt.Sprintf("%.0f, %s", value, word)
	if opts.Status == data.DataLoading || opts.Status == data.DataLoadingMore || opts.Status == data.DataUnloaded {
		label = core.Msg(c, "list.loading")
	} else if opts.Status == data.DataFailed {
		label = data.DataError(c, opts.Error)
	} else if opts.Empty {
		label = core.Msg(c, "chart.empty")
	}
	if opts.Label != "" {
		label = opts.Label + ": " + label
	}
	pal := internal.Palette(k)
	cols := [5]ui.Color{pal[0], pal[3], pal[5], pal[1], pal[2]}
	size := core.FontSize(c, theme.CaptionSize)
	col := ui.Column(c).Gap(u).AlignItems(ui.Center).Role(ui.RoleMeter).Label(label).Focusable()
	internal.Ring(col, k.Focus, theme.ControlRadius)
	if !col.IsDisabled() {
		for _, key := range []ui.Key{ui.KeyLeft, ui.KeyRight, ui.KeyUp, ui.KeyDown, ui.KeyHome, ui.KeyEnd} {
			if col.Shortcut(0, key) {
				c.Announce(label)
			}
		}
	}
	col.Children(func() {
		if opts.Label != "" {
			ui.Text(c, opts.Label).Font(theme.SerifFont).TextColor(k.TextMuted)
		}
		if opts.Status != data.DataReady || opts.Empty {
			ui.Box(c).Size(w, dialHeight).Children(func() {
				ui.Column(c).Fill().Children(func() {
					if data.DataPlaceholder(c, opts.Status, opts.Error, "", core.Msg(c, "chart.empty")) && opts.Retry != nil {
						opts.Retry()
					}
				})
			})
			ui.Box(c).Height(core.FontSize(c, theme.H1Size)*theme.LineHeight + u + core.FontSize(c, theme.BodySize)*theme.LineHeight)
			return
		}
		ui.Box(c).Size(w, dialHeight).Draw(func(p *ui.Painter, r ui.Rect) {
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
			chart.ChartText(p, "0", r.X+2, cy+4, size, k.TextMuted, 0, 0)
			chart.ChartText(p, "100", r.X+r.W-2, cy+4, size, k.TextMuted, 1, 0)
		})
		ui.Text(c, fmt.Sprintf("%.0f", value)).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H1Size)).TextColor(k.Text).FontFeatures("tnum")
		ui.Text(c, word).FontWeight(600).TextColor(k.Text)
	})
	return col
}
