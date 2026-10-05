package media

import (
	"log/slog"
	"slices"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

// playbackRates are the default choices, 0.5× to 2×.
var playbackRates = []float64{0.5, 0.75, 1, 1.25, 1.5, 2}

// PlaybackSpeedControlOptions configure a PlaybackSpeedControl.
type PlaybackSpeedControlOptions struct {
	// Rates are the ascending choices; nil means 0.5× to 2×.
	Rates    []float64
	Disabled bool
}

// PlaybackSpeedControl chooses *rate from Rates in a drop-down.
func PlaybackSpeedControl(c *ui.Context, rate *float64, opts PlaybackSpeedControlOptions) core.ChoiceResult {
	rates := opts.Rates
	if rates == nil {
		rates = playbackRates
	}
	if len(rates) == 0 || !slices.IsSorted(rates) || rates[0] <= 0 {
		panic("mujica: PlaybackSpeedControl needs positive Rates in ascending order")
	}
	if !slices.Contains(rates, *rate) {
		panic("mujica: PlaybackSpeedControl rate is not one of Rates")
	}
	options := make([]input.SelectOption[float64], len(rates))
	for i, r := range rates {
		options[i] = input.SelectOption[float64]{Value: r, Label: playbackLabel(r)}
	}
	res := input.Select(c, rate, options, input.SelectOptions{Disabled: opts.Disabled, Label: core.Msg(c, "playback_speed_control.label")})
	res.Element.MinWidth(88).Width(96)
	if res.Changed() {
		slog.Debug("mujica: playback_speed changed", "rate", *rate)
	}
	return res
}

// playbackLabel writes a rate as "1.25×".
func playbackLabel(r float64) string { return strconv.FormatFloat(r, 'f', -1, 64) + "×" }
