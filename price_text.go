package mujica

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// Trend picks the rise color: green, or red (China).
type Trend int

const (
	GreenUp Trend = iota
	RedUp
)

type trendKey struct{}

// UseTrend sets the window's rise and fall colors.
func UseTrend(c *ui.Context, t Trend) {
	if t != GreenUp && t != RedUp {
		panic(fmt.Sprintf("mujica: unknown Trend %d", t))
	}
	*ui.Local(c.Root(), trendKey{}, func() Trend { return GreenUp }) = t
}

// trendColor colors a rise (sign > 0), fall or flat.
func trendColor(c *ui.Context, sign int) ui.Color {
	k := Tokens(c)
	up, down := k.Success, k.Danger
	if *ui.Local(c.Root(), trendKey{}, func() Trend { return GreenUp }) == RedUp {
		up, down = k.Danger, k.Success
	}
	switch {
	case sign > 0:
		return up
	case sign < 0:
		return down
	}
	return k.TextMuted
}

// trendIcon is the arrow that says the direction without color.
func trendIcon(sign int) *ui.SVG {
	switch {
	case sign > 0:
		return icons.Must("arrow-up")
	case sign < 0:
		return icons.Must("arrow-down")
	}
	return icons.Must("minus")
}

// trendWord names the direction for assistive technology.
func trendWord(c *ui.Context, sign int) string {
	switch {
	case sign > 0:
		return msg(c, "price_text.up")
	case sign < 0:
		return msg(c, "price_text.down")
	}
	return msg(c, "price_text.flat")
}

func trendSign(v float64) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

// priceDecimals resolves a Decimals option: zero means 2, negative panics.
func priceDecimals(name string, d int) int {
	if d < 0 {
		panic("mujica: " + name + " with negative Decimals")
	}
	if d == 0 {
		return 2
	}
	return d
}

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
func PriceText(c *ui.Context, value float64, opts PriceTextOptions) *ui.Element {
	dec := priceDecimals("PriceText", opts.Decimals)
	if opts.Size < 0 {
		panic("mujica: PriceText with negative Size")
	}
	size := opts.Size
	if size == 0 {
		size = theme.BodySize
	}
	text := statisticFormat(value, dec)
	if opts.NoGrouping {
		text = priceUngrouped(value, dec)
	}
	e := ui.Row(c).Gap(Density(c).Unit()).AlignItems(ui.Center).Padding(0, 2).Radius(theme.ControlRadius).Shrink(0)
	st := ui.Local(e, "price", func() priceTextState { return priceTextState{} })
	if st.seen && value != st.last {
		st.at, st.dir = c.Now(), trendSign(value-st.last)
		slog.Debug("mujica: price_text changed", "dir", st.dir)
	}
	st.last, st.seen = value, true
	sign := 0
	if opts.Compare {
		sign = trendSign(value - opts.Reference)
	}
	col := Tokens(c).Text
	spoken := text
	if opts.Compare {
		col = trendColor(c, sign)
		spoken += ", " + trendWord(c, sign)
	}
	if d := Motion(c, priceTextFlash); d > 0 && !opts.NoFlash && !st.at.IsZero() {
		if t := c.Now().Sub(st.at); t < d {
			e.Background(trendColor(c, st.dir).Alpha(0.2 * float32(1-float64(t)/float64(d))))
			c.AnimationFrame()
		}
	}
	e.Role(ui.RoleText).Label(spoken).TextColor(col)
	e.Children(func() {
		if opts.Compare {
			ui.Icon(c, trendIcon(sign)).Size(size*0.85, size*0.85)
		}
		ui.Text(c, text).FontSize(FontSize(c, size)).FontFeatures("tnum").FontWeight(500).SingleLine()
	})
	return e
}

// priceUngrouped writes v with dec decimals and no separators.
func priceUngrouped(v float64, dec int) string {
	s := fmt.Sprintf("%.*f", dec, math.Abs(v))
	if v < 0 && s != fmt.Sprintf("%.*f", dec, 0.0) {
		return "−" + s
	}
	return s
}
