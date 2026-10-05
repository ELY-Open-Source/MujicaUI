package mujica

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
)

// CurrencyInputOptions configure a CurrencyInput.
type CurrencyInputOptions struct {
	// Symbol goes before the amount, as "$" or "¥"; empty shows none.
	Symbol string
	// Decimals is the minor unit digits; zero is 2. NoFraction makes it 0, as for JPY.
	Decimals   int
	NoFraction bool
	Label      string
	Disabled   bool
	// Error marks the value invalid and shows the message below.
	Error string
}

// currencyInputDigits resolves the decimals of opts.
func currencyInputDigits(opts CurrencyInputOptions) int {
	switch {
	case opts.Decimals < 0 || opts.Decimals > 6:
		panic(fmt.Sprintf("mujica: CurrencyInput Decimals %d outside 0..6", opts.Decimals))
	case opts.NoFraction && opts.Decimals != 0:
		panic("mujica: CurrencyInput NoFraction with Decimals")
	case opts.NoFraction:
		return 0
	case opts.Decimals == 0:
		return 2
	}
	return opts.Decimals
}

// CurrencyFormat writes minor units as an amount with thousands separators.
func CurrencyFormat(minor int64, decimals int, group bool) string {
	neg := minor < 0
	u := uint64(minor)
	if neg {
		u = -u
	}
	p := uint64(math.Pow10(decimals))
	whole := strconv.FormatUint(u/p, 10)
	if group {
		var b strings.Builder
		for i, r := range whole {
			if i > 0 && (len(whole)-i)%3 == 0 {
				b.WriteByte(',')
			}
			b.WriteRune(r)
		}
		whole = b.String()
	}
	s := whole
	if decimals > 0 {
		s += "." + fmt.Sprintf("%0*d", decimals, u%p)
	}
	if neg {
		s = "-" + s
	}
	return s
}

// currencyInputParse reads an amount into minor units; reason is a message key when it fails.
func currencyInputParse(text string, decimals int) (minor int64, reason string) {
	t := strings.NewReplacer(",", "", " ", "", " ", "").Replace(strings.TrimSpace(text))
	neg := strings.HasPrefix(t, "-")
	t = strings.TrimPrefix(t, "-")
	whole, frac, dot := strings.Cut(t, ".")
	if whole == "" && frac == "" || strings.Trim(whole+frac, "0123456789") != "" || dot && strings.Contains(frac, ".") {
		return 0, "currency_input.invalid"
	}
	if len(frac) > decimals {
		return 0, "currency_input.precision"
	}
	frac += strings.Repeat("0", decimals-len(frac))
	v, err := strconv.ParseInt(whole+frac, 10, 64)
	if err != nil {
		return 0, "currency_input.invalid"
	}
	if neg {
		v = -v
	}
	return v, ""
}

// CurrencyInput edits *value, an amount in minor units (cents); it groups thousands while not focused.
func CurrencyInput(c *ui.Context, value *int64, opts CurrencyInputOptions) InputResult {
	if value == nil {
		panic("mujica: CurrencyInput needs a value")
	}
	dec := currencyInputDigits(opts)
	k := Tokens(c)
	var r InputResult
	look := &textInputLook{disabled: opts.Disabled, err: opts.Error}
	r.Element, r.Input = textInputFrame(c, look, func(frame *ui.Element) *ui.Element {
		if opts.Symbol != "" {
			ui.Text(c, opts.Symbol).TextColor(k.TextMuted).SingleLine()
		}
		text := ui.Local(frame, "text", func() string { return CurrencyFormat(*value, dec, false) })
		bad := ui.Local(frame, "bad", func() string { return "" })
		seen := ui.Local(frame, "seen", func() int64 { return *value })
		was := ui.Local(frame, "was", func() bool { return false })
		if *seen != *value {
			*bad = ""
		}
		// Grouped while resting; plain digits once focused, for editing.
		focused := frame.FocusWithin()
		switch {
		case *bad != "":
		case !focused:
			*text = CurrencyFormat(*value, dec, true)
		case !*was || *seen != *value:
			*text = CurrencyFormat(*value, dec, false)
		}
		*was = focused
		in := textInputEditor(c, frame, text, look.disabled, ui.TextInputBase).Grow(1).MinWidth(Density(c).Unit() * 16).Label(opts.Label)
		if in.Changed() && !look.disabled {
			v, why := currencyInputParse(*text, dec)
			*bad = ""
			switch {
			case why == "currency_input.precision":
				*bad = fmt.Sprintf(msg(c, why), dec)
			case why != "":
				*bad = msg(c, why)
			case v != *value:
				*value, r.changed = v, true
			}
		}
		*seen = *value
		if look.err == "" {
			look.err = *bad
		}
		return in
	}, nil)
	if r.changed {
		slog.Debug("mujica: currency_input change", "minor", *value)
	}
	return r
}
