package chart

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
)

// IndicatorKind names a technical indicator.
type IndicatorKind int

const (
	IndicatorMA IndicatorKind = iota
	IndicatorEMA
	IndicatorBOLL
	IndicatorVWAP
	IndicatorMACD
	IndicatorRSI
	IndicatorKDJ
	IndicatorOBV
)

// indicatorInfo describes names, default parameters, and plotting locations.
var indicatorInfo = []struct {
	name    string
	params  []float64
	overlay bool
}{
	{"MA", []float64{20}, true}, {"EMA", []float64{20}, true}, {"BOLL", []float64{20, 2}, true}, {"VWAP", nil, true},
	{"MACD", []float64{12, 26, 9}, false}, {"RSI", []float64{14}, false}, {"KDJ", []float64{9, 3, 3}, false}, {"OBV", nil, false},
}

func (k IndicatorKind) String() string {
	if k < IndicatorMA || k > IndicatorOBV {
		panic(fmt.Sprintf("mujica: unknown IndicatorKind %d", int(k)))
	}
	return indicatorInfo[k].name
}

// DefaultParams returns a fresh copy of the kind's usual parameters.
func (k IndicatorKind) DefaultParams() []float64 {
	_ = k.String()
	return append([]float64(nil), indicatorInfo[k].params...)
}

// IndicatorOverlay configures MA, EMA, BOLL, or VWAP.
type IndicatorOverlay struct {
	Kind IndicatorKind
	// Params nil uses the indicator's defaults.
	Params []float64
	// Values overrides calculation; supply one value per candle.
	Values []float64
}

// indicatorParams resolves and checks params for kind.
func indicatorParams(k IndicatorKind, params []float64) []float64 {
	_ = k.String()
	want := indicatorInfo[k].params
	if params == nil {
		return want
	}
	if len(params) != len(want) {
		panic(fmt.Sprintf("mujica: %v takes %d parameters, not %d", k, len(want), len(params)))
	}
	for i, p := range params {
		if math.IsInf(p, 0) || !(p > 0) || (k != IndicatorBOLL || i == 0) && p != float64(int(p)) {
			panic(fmt.Sprintf("mujica: %v parameter %v must be a positive whole number", k, p))
		}
	}
	if k == IndicatorMACD && params[0] >= params[1] {
		panic("mujica: MACD fast period must be below slow period")
	}
	return params
}

// indicatorName is the kind with its parameters, as "MACD(12,26,9)".
func indicatorName(k IndicatorKind, params []float64) string {
	if len(params) == 0 {
		return k.String()
	}
	parts := make([]string, len(params))
	for i, p := range params {
		parts[i] = strconv.FormatFloat(p, 'f', -1, 64)
	}
	return k.String() + "(" + strings.Join(parts, ",") + ")"
}

// indicatorSeries computes named indicator lines from candles.
func indicatorSeries(k IndicatorKind, params []float64, candles []Candle) []ChartSeries {
	p := indicatorParams(k, params)
	closes := make([]float64, len(candles))
	for i, c := range candles {
		closes[i] = c.Close
	}
	name := indicatorName(k, p)
	switch k {
	case IndicatorMA:
		return []ChartSeries{{name, SMA(closes, int(p[0]))}}
	case IndicatorEMA:
		return []ChartSeries{{name, EMA(closes, int(p[0]))}}
	case IndicatorBOLL:
		mid, up, lo := Bollinger(closes, int(p[0]), p[1])
		return []ChartSeries{{name, mid}, {name + " ↑", up}, {name + " ↓", lo}}
	case IndicatorVWAP:
		return []ChartSeries{{name, VWAP(candles)}}
	case IndicatorMACD:
		dif, dea, hist := MACD(closes, int(p[0]), int(p[1]), int(p[2]))
		return []ChartSeries{{"DIF", dif}, {"DEA", dea}, {"MACD", hist}}
	case IndicatorRSI:
		return []ChartSeries{{name, RSI(closes, int(p[0]))}}
	case IndicatorKDJ:
		kk, d, j := KDJ(candles, int(p[0]), int(p[1]), int(p[2]))
		return []ChartSeries{{"K", kk}, {"D", d}, {"J", j}}
	}
	return []ChartSeries{{name, OBV(candles)}}
}

// mktOverlayLines caches indicators until data or parameters change.
func mktOverlayLines(ch *chart, m *mktCore, candles []Candle) []mktLine {
	pal := internal.Palette(m.k)
	var out []mktLine
	for i, o := range m.opts.Overlays {
		if _ = o.Kind.String(); !indicatorInfo[o.Kind].overlay {
			panic(fmt.Sprintf("mujica: %v is not an overlay indicator", o.Kind))
		}
		col := pal[[]int{3, 1, 5, 2, 4}[i%5]]
		if o.Values != nil {
			if len(o.Values) != len(candles) {
				panic("mujica: IndicatorOverlay Values need one value per candle")
			}
			out = append(out, mktLine{name: indicatorName(o.Kind, indicatorParams(o.Kind, o.Params)), values: o.Values, color: col, shape: len(out)})
			continue
		}
		key := fmt.Sprintf("%s|%d|%v", ChartID(candles), o.Kind, o.Params)
		series := ChartMemo(ch, fmt.Sprint("overlay.", i), key, func() []ChartSeries { return indicatorSeries(o.Kind, o.Params, candles) })
		for _, s := range series {
			out = append(out, mktLine{name: s.Name, values: s.Values, color: col, shape: len(out)})
		}
	}
	return out
}

// IndicatorParamsField edits validated indicator parameters in place.
func IndicatorParamsField(c *ui.Context, kind IndicatorKind, params *[]float64) core.ChoiceResult {
	if params == nil {
		panic("mujica: IndicatorParamsField needs params")
	}
	row := ui.Row(c).Gap(core.Density(c).Unit() * 2).AlignItems(ui.Center).Role(ui.RoleGroup).Label(kind.String())
	if row.IsDisabled() {
		copy := append([]float64(nil), (*params)...)
		params = &copy
	}
	if *params == nil {
		*params = kind.DefaultParams()
	}
	changed := false
	indicatorParams(kind, *params)
	row.Children(func() {
		for i := range *params {
			v := (*params)[i]
			step, prec := 1.0, 0
			if kind == IndicatorBOLL && i == 1 {
				step, prec = 0.1, 1
			}
			if input.NumberInput(c, &v, input.NumberInputOptions{Min: step, Max: 500, Step: step, Precision: prec, Label: fmt.Sprintf("%v %d", kind, i+1)}).Element.Width(110).Key(i).Changed() || v != (*params)[i] {
				if kind == IndicatorMACD && (i == 0 && v >= (*params)[1] || i == 1 && v <= (*params)[0]) {
					continue
				}
				(*params)[i], changed = v, true
			}
		}
	})
	if changed {
		slog.Debug("mujica: indicator_overlay params", "kind", kind)
	}
	return core.NewChoiceResult(row, changed)
}
