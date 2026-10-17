package finance

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// priceTextFlash is the tint's fade, three state steps long.
const priceTextFlash = 3 * theme.StateDuration

// PriceTextOptions configure a PriceText.
type PriceTextOptions struct {
	// Decimals shown; zero means 2.
	Decimals int
	// NoGrouping drops the thousands separators.
	NoGrouping bool
	// Compare marks the price's move from Reference.
	Compare   bool
	Reference float64
	// NoFlash keeps the price still when it changes.
	NoFlash bool
	// Size is the font size; zero is the body size.
	Size float32
}

type priceTextState struct {
	last float64
	seen bool
	at   time.Time
	dir  int
}

// PriceText shows a price; a change tints it once.
func PriceText(c *ui.Context, value float64, opts PriceTextOptions) ui.Element {
	dec := core.PriceDecimals("PriceText", opts.Decimals)
	if opts.Size < 0 {
		panic("mujica: PriceText with negative Size")
	}
	size := opts.Size
	if size == 0 {
		size = theme.BodySize
	}
	text := data.StatisticFormat(value, dec)
	if opts.NoGrouping {
		text = core.PriceUngrouped(value, dec)
	}
	e := ui.Row(c).Gap(core.Density(c).Unit()).AlignItems(ui.Center).Padding(0, 2).Radius(theme.ControlRadius).Shrink(0)
	st := ui.Local(e, "price", func() priceTextState { return priceTextState{} })
	if st.seen && value != st.last {
		st.at, st.dir = c.Now(), core.TrendSign(value-st.last)
		slog.Debug("mujica: price_text changed", "dir", st.dir)
	}
	st.last, st.seen = value, true
	sign := 0
	if opts.Compare {
		sign = core.TrendSign(value - opts.Reference)
	}
	col := core.Tokens(c).Text
	spoken := text
	if opts.Compare {
		col = core.TrendColor(c, sign)
		spoken += ", " + core.TrendWord(c, sign)
	}
	if d := core.Motion(c, priceTextFlash); d > 0 && !opts.NoFlash && !st.at.IsZero() {
		if t := c.Now().Sub(st.at); t < d {
			e.Background(core.TrendColor(c, st.dir).Alpha(0.2 * float32(1-float64(t)/float64(d))))
			c.AnimationFrame()
		}
	}
	e.Role(ui.RoleText).Label(spoken).TextColor(col)
	e.Children(func() {
		if opts.Compare {
			ui.Icon(c, core.TrendIcon(sign)).Size(size*0.85, size*0.85)
		}
		ui.Text(c, text).FontSize(core.FontSize(c, size)).FontFeatures("tnum").FontWeight(500).SingleLine()
	})
	return e
}
