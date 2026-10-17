package chat

import (
	"fmt"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CostEstimatorOptions configure a CostEstimator; prices are per million tokens.
type CostEstimatorOptions struct {
	InputTokens, OutputTokens int
	InputPrice, OutputPrice   float64
	// Currency prefixes amounts; empty is "$".
	Currency string
	// Digits is the decimals shown for amounts; zero is 4.
	Digits int
}

// CostEstimate is the input, output and total cost.
type CostEstimate struct{ Input, Output, Total float64 }

// EstimateCost prices token counts at per-million rates.
func EstimateCost(inputTokens, outputTokens int, inputPrice, outputPrice float64) CostEstimate {
	in := float64(inputTokens) * inputPrice / 1e6
	out := float64(outputTokens) * outputPrice / 1e6
	return CostEstimate{in, out, in + out}
}

// CostEstimator shows token counts and cost at caller prices.
func CostEstimator(c *ui.Context, opts CostEstimatorOptions) ui.Element {
	if opts.InputTokens < 0 || opts.OutputTokens < 0 {
		panic("mujica: CostEstimator token counts must not be negative")
	}
	for _, p := range []float64{opts.InputPrice, opts.OutputPrice} {
		if p < 0 || math.IsNaN(p) || math.IsInf(p, 0) {
			panic("mujica: CostEstimator prices must be finite and not negative")
		}
	}
	if opts.Digits < 0 {
		panic("mujica: CostEstimator Digits is negative")
	}
	digits := opts.Digits
	if digits == 0 {
		digits = 4
	}
	cur := opts.Currency
	if cur == "" {
		cur = "$"
	}
	money := func(v float64) string { return cur + strconv.FormatFloat(v, 'f', digits, 64) }
	est := EstimateCost(opts.InputTokens, opts.OutputTokens, opts.InputPrice, opts.OutputPrice)
	k := core.Tokens(c)
	d := core.Density(c)
	total := fmt.Sprintf(core.Msg(c, "cost_estimator.total_label"), money(est.Total))
	g := ui.Grid(c).ColumnTracks(ui.Fr(1), ui.FitContent(), ui.FitContent()).GapX(d.Unit()*4).GapY(d.Unit()).
		Padding(d.Unit()*3).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).
		Role(ui.RoleGroup).Label(total)
	cell := func(s string, ink ui.Color, right bool) {
		t := ui.Text(c, s).TextColor(ink).SingleLine().FontFeatures("tnum")
		if right {
			t.JustifySelf(ui.End)
		}
	}
	g.Children(func() {
		cell(core.Msg(c, "cost_estimator.kind"), k.TextMuted, false)
		cell(core.Msg(c, "cost_estimator.tokens"), k.TextMuted, true)
		cell(core.Msg(c, "cost_estimator.cost"), k.TextMuted, true)
		cell(core.Msg(c, "cost_estimator.input"), k.Text, false)
		cell(strconv.Itoa(opts.InputTokens), k.Text, true)
		cell(money(est.Input), k.Text, true)
		cell(core.Msg(c, "cost_estimator.output"), k.Text, false)
		cell(strconv.Itoa(opts.OutputTokens), k.Text, true)
		cell(money(est.Output), k.Text, true)
		ui.Box(c).ColumnSpan(3).Height(theme.BorderWidth).MarginY(d.Unit()).Background(k.Ornament)
		ui.Text(c, core.Msg(c, "cost_estimator.total")).Font(theme.SerifFont).TextColor(k.Text)
		cell(strconv.Itoa(opts.InputTokens+opts.OutputTokens), k.Text, true)
		ui.Text(c, money(est.Total)).Bold().TextColor(k.AccentText).JustifySelf(ui.End).FontFeatures("tnum")
		ui.Text(c, fmt.Sprintf(core.Msg(c, "cost_estimator.rates"), cur+strconv.FormatFloat(opts.InputPrice, 'f', -1, 64),
			cur+strconv.FormatFloat(opts.OutputPrice, 'f', -1, 64))).
			ColumnSpan(3).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
	})
	return g
}
