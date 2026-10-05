package chart

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// SparkBarOptions configure a SparkBar; Frame.Height 0 means 28 DIPs.
type SparkBarOptions struct {
	Frame ChartFrame
	// BarWidth defaults to 6 DIPs; width fits the window.
	BarWidth float32
	Format   func(v float64) string
}

// SparkBar draws signed inline bars and highlights the last.
func SparkBar(c *ui.Context, values []float64, opts SparkBarOptions) ChartResult {
	if opts.BarWidth < 0 || math.IsNaN(float64(opts.BarWidth)) || math.IsInf(float64(opts.BarWidth), 0) {
		panic("mujica: SparkBar BarWidth must be ≥ 0")
	}
	bw := opts.BarWidth
	if bw == 0 {
		bw = 6
	}
	w, _ := c.Size()
	frame := opts.Frame
	if frame.Empty == "" {
		frame.Empty = core.Msg(c, "spark_bar.empty")
	}
	return Sparkline(c, values, SparklineOptions{Frame: frame, Kind: SparklineBar, Width: min(max(28, float32(len(values))*(bw+1)+6), max(28, w-48)), Format: opts.Format})
}
