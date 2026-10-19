package finance

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

type Currency struct {
	Code, Name string
	Decimals   int
}

// CurrencyConverterState is the amount and the pair being converted.
type CurrencyConverterState struct {
	Amount           float64
	From, To         string
	rate             float64
	rateFrom, rateTo string
	rateBound        bool
}

// CurrencyConverterOptions configure a CurrencyConverter.
type CurrencyConverterOptions struct {
	Rate float64
	// RateFrom and RateTo identify refreshed rates, including equal rates.
	RateFrom, RateTo string
	RateTime         time.Time
	Disabled         bool
}

type CurrencyConverterView struct {
	Element ui.Element
	Result  float64
	// PairChanged asks the caller for the new pair's rate.
	PairChanged bool
}

func CurrencyConverter(c *ui.Context, s *CurrencyConverterState, currencies []Currency, opts CurrencyConverterOptions) CurrencyConverterView {
	if s == nil || len(currencies) < 2 {
		panic("mujica: CurrencyConverter needs a state and two currencies")
	}
	if opts.Rate < 0 || math.IsNaN(opts.Rate) || math.IsInf(opts.Rate, 0) {
		panic("mujica: CurrencyConverter Rate must be finite and not negative")
	}
	seen := map[string]bool{}
	for _, cur := range currencies {
		if cur.Code == "" || seen[cur.Code] || cur.Decimals < 0 || cur.Decimals > 12 {
			panic("mujica: currencies need unique codes and decimals in 0..12")
		}
		seen[cur.Code] = true
	}
	if s.Amount < 0 || math.IsNaN(s.Amount) || math.IsInf(s.Amount, 0) {
		panic("mujica: currency amount must be finite and nonnegative")
	}
	from, to := currencyFind(currencies, s.From), currencyFind(currencies, s.To)
	rate := opts.Rate
	if !s.rateBound || s.rate != rate || opts.RateFrom != "" || opts.RateTo != "" {
		s.rate, s.rateFrom, s.rateTo, s.rateBound = rate, s.From, s.To, true
		if opts.RateFrom != "" || opts.RateTo != "" {
			s.rateFrom, s.rateTo = opts.RateFrom, opts.RateTo
		}
	}
	if s.From != s.rateFrom || s.To != s.rateTo {
		rate = 0
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	options := make([]input.SelectOption[string], len(currencies))
	for i, cur := range currencies {
		options[i] = input.SelectOption[string]{Value: cur.Code, Label: cur.Code + " · " + cur.Name}
	}
	v := CurrencyConverterView{Result: math.NaN()}
	if rate > 0 {
		v.Result = input.NumberInputSpec{Decimals: to.Decimals}.Round(s.Amount * rate)
	}
	v.Element = ui.Column(c).Gap(u*3).Padding(u*4).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Radius(theme.CardRadius).Disabled(opts.Disabled)
	opts.Disabled = v.Element.IsDisabled()
	if opts.Disabled {
		copy := *s
		s = &copy
	}
	v.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.End).Wrap().Children(func() {
			input.FormField(c, core.Msg(c, "currency_converter.amount"), input.FormFieldOptions{}, func() ui.Element {
				r := input.NumberInput(c, &s.Amount, input.NumberInputOptions{Min: 0, Max: 1e15, Step: math.Pow(10, -float64(from.Decimals)),
					Precision: from.Decimals, Label: core.Msg(c, "currency_converter.amount"), Disabled: opts.Disabled})
				return r.Element.Width(170)
			})
			input.FormField(c, core.Msg(c, "currency_converter.from"), input.FormFieldOptions{}, func() ui.Element {
				r := input.Select(c, &s.From, options, input.SelectOptions{Label: core.Msg(c, "currency_converter.from"), Disabled: opts.Disabled})
				v.PairChanged = v.PairChanged || r.Changed()
				return r.Element
			})
			swap := core.IconAction(c, icons.Must("arrow-up-down"), core.Msg(c, "currency_converter.swap")).Rotate(90)
			if swap.Clicked() && !swap.IsDisabled() {
				s.From, s.To, v.PairChanged = s.To, s.From, true
			}
			input.FormField(c, core.Msg(c, "currency_converter.to"), input.FormFieldOptions{}, func() ui.Element {
				r := input.Select(c, &s.To, options, input.SelectOptions{Label: core.Msg(c, "currency_converter.to"), Disabled: opts.Disabled})
				v.PairChanged = v.PairChanged || r.Changed()
				return r.Element
			})
		})
		from, to = currencyFind(currencies, s.From), currencyFind(currencies, s.To)
		if v.PairChanged {
			v.Result = math.NaN()
		} else if rate > 0 {
			v.Result = input.NumberInputSpec{Decimals: to.Decimals}.Round(s.Amount * rate)
		}
		if rate == 0 || v.PairChanged {
			ui.Text(c, core.Msg(c, "currency_converter.no_rate")).TextColor(k.TextMuted).Role(ui.RoleStatus)
			return
		}
		ui.Row(c).Gap(u * 2).AlignItems(ui.End).Role(ui.RoleStatus).Label(fmt.Sprintf("%s %s = %s %s",
			data.StatisticFormat(s.Amount, from.Decimals), from.Code, data.StatisticFormat(v.Result, to.Decimals), to.Code)).Children(func() {
			ui.Text(c, data.StatisticFormat(s.Amount, from.Decimals)+" "+from.Code+" =").TextColor(k.TextMuted).FontFeatures("tnum")
			ui.Text(c, data.StatisticFormat(v.Result, to.Decimals)+" "+to.Code).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H2Size)).
				TextColor(k.Text).FontFeatures("tnum")
		})
		note := fmt.Sprintf(core.Msg(c, "currency_converter.rate"), from.Code, rate, to.Code)
		if !opts.RateTime.IsZero() {
			note += " · " + opts.RateTime.Format("2006-01-02 15:04")
		}
		ui.Text(c, note).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
	})
	if v.Element.IsDisabled() {
		v.PairChanged = false
	}
	if v.PairChanged {
		slog.Debug("mujica: currency_converter pair", "from", s.From, "to", s.To)
	}
	return v
}

// currencyFind returns the currency with code, panicking when absent.
func currencyFind(currencies []Currency, code string) Currency {
	for _, cur := range currencies {
		if cur.Decimals < 0 {
			panic("mujica: Currency Decimals is negative")
		}
		if cur.Code == code {
			return cur
		}
	}
	panic("mujica: CurrencyConverter currency " + code + " is not offered")
}
