package core

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
)

// Trend picks the rise color: green, or red (China).
type Trend int

type trendKey struct{}

// UseTrend sets the window's rise and fall colors.
func UseTrend(c *ui.Context, t Trend) {
	if t != GreenUp && t != RedUp {
		panic(fmt.Sprintf("mujica: unknown Trend %d", t))
	}
	*ui.Local(c.Root(), trendKey{}, func() Trend { return GreenUp }) = t
}

// TrendColor colors a rise (sign > 0), fall or flat.
func TrendColor(c *ui.Context, sign int) ui.Color {
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

// TrendIcon is the arrow that says the direction without color.
func TrendIcon(sign int) *ui.SVG {
	switch {
	case sign > 0:
		return icons.Must("arrow-up")
	case sign < 0:
		return icons.Must("arrow-down")
	}
	return icons.Must("minus")
}

// TrendWord names the direction for assistive technology.
func TrendWord(c *ui.Context, sign int) string {
	switch {
	case sign > 0:
		return Msg(c, "price_text.up")
	case sign < 0:
		return Msg(c, "price_text.down")
	}
	return Msg(c, "price_text.flat")
}

func TrendSign(v float64) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

// PriceDecimals resolves a Decimals option: zero means 2, negative panics.
func PriceDecimals(name string, d int) int {
	if d < 0 {
		panic("mujica: " + name + " with negative Decimals")
	}
	if d == 0 {
		return 2
	}
	return d
}

// PriceUngrouped writes v with dec decimals and no separators.
func PriceUngrouped(v float64, dec int) string {
	s := fmt.Sprintf("%.*f", dec, math.Abs(v))
	if v < 0 && s != fmt.Sprintf("%.*f", dec, 0.0) {
		return "−" + s
	}
	return s
}

const (
	GreenUp Trend = iota
	RedUp
)
