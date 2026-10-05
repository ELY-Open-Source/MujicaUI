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

// Currency is a currency the converter offers; Decimals is its minor-unit precision.
type Currency struct {
	Code, Name string
	Decimals   int
}

// CurrencyConverterState is the amount and the pair being converted.
type CurrencyConverterState struct {
	Amount   float64
	From, To string
}

// CurrencyConverterOptions configure a CurrencyConverter.
type CurrencyConverterOptions struct {
	// Rate converts one From into To; zero shows it is unavailable.
	Rate float64
	// RateTime is when the caller's rate was quoted; zero hides it.
	RateTime time.Time
	Disabled bool
}

// CurrencyConverterView is a built converter, the result, and whether the pair changed.
type CurrencyConverterView struct {
	Element *ui.Element
	// Result is Amount in To, rounded to its decimals; NaN without a rate.
	Result float64
	// PairChanged asks the caller for the new pair's rate.
	PairChanged bool
}

// CurrencyConverter picks and swaps two currencies and converts at the caller's rate.
func CurrencyConverter(c *ui.Context, s *CurrencyConverterState, currencies []Currency, opts CurrencyConverterOptions) CurrencyConverterView {
	if s == nil || len(currencies) < 2 {
		panic("mujica: CurrencyConverter needs a state and two currencies")
	}
	if opts.Rate < 0 || math.IsNaN(opts.Rate) || math.IsInf(opts.Rate, 0) {
		panic("mujica: CurrencyConverter Rate must be finite and not negative")
	}
	from, to := currencyFind(currencies, s.From), currencyFind(currencies, s.To)
	k := Tokens(c)
	u := Density(c).Unit()
	options := make([]SelectOption[string], len(currencies))
	for i, cur := range currencies {
		options[i] = SelectOption[string]{Value: cur.Code, Label: cur.Code + " · " + cur.Name}
	}
	v := CurrencyConverterView{Result: math.NaN()}
	if opts.Rate > 0 {
		v.Result = numberInputSpec{decimals: to.Decimals}.round(s.Amount * opts.Rate)
	}
	v.Element = ui.Column(c).Gap(u*3).Padding(u*4).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Radius(theme.CardRadius).Disabled(opts.Disabled)
	v.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.End).Wrap().Children(func() {
			FormField(c, msg(c, "currency_converter.amount"), FormFieldOptions{}, func() *ui.Element {
				r := NumberInput(c, &s.Amount, NumberInputOptions{Min: 0, Max: 1e15, Step: math.Pow(10, -float64(from.Decimals)),
					Precision: from.Decimals, Label: msg(c, "currency_converter.amount"), Disabled: opts.Disabled})
				return r.Element.Width(170)
			})
			FormField(c, msg(c, "currency_converter.from"), FormFieldOptions{}, func() *ui.Element {
				r := Select(c, &s.From, options, SelectOptions{Label: msg(c, "currency_converter.from"), Disabled: opts.Disabled})
				v.PairChanged = v.PairChanged || r.Changed()
				return r.Element
			})
			swap := iconAction(c, icons.Must("arrow-up-down"), msg(c, "currency_converter.swap")).Rotate(90)
			if swap.Clicked() && !swap.IsDisabled() {
				s.From, s.To, v.PairChanged = s.To, s.From, true
			}
			FormField(c, msg(c, "currency_converter.to"), FormFieldOptions{}, func() *ui.Element {
				r := Select(c, &s.To, options, SelectOptions{Label: msg(c, "currency_converter.to"), Disabled: opts.Disabled})
				v.PairChanged = v.PairChanged || r.Changed()
				return r.Element
			})
		})
		if opts.Rate == 0 {
			ui.Text(c, msg(c, "currency_converter.no_rate")).TextColor(k.TextMuted).Role(ui.RoleStatus)
			return
		}
		ui.Row(c).Gap(u * 2).AlignItems(ui.End).Role(ui.RoleStatus).Label(fmt.Sprintf("%s %s = %s %s",
			statisticFormat(s.Amount, from.Decimals), from.Code, statisticFormat(v.Result, to.Decimals), to.Code)).Children(func() {
			ui.Text(c, statisticFormat(s.Amount, from.Decimals)+" "+from.Code+" =").TextColor(k.TextMuted).FontFeatures("tnum")
			ui.Text(c, statisticFormat(v.Result, to.Decimals)+" "+to.Code).Font(theme.SerifFont).FontSize(FontSize(c, theme.H2Size)).
				TextColor(k.Text).FontFeatures("tnum")
		})
		note := fmt.Sprintf(msg(c, "currency_converter.rate"), from.Code, opts.Rate, to.Code)
		if !opts.RateTime.IsZero() {
			note += " · " + opts.RateTime.Format("2006-01-02 15:04")
		}
		ui.Text(c, note).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
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
