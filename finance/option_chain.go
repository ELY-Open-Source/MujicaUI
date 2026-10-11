package finance

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// OptionQuote is one contract's market.
type OptionQuote struct {
	Bid, Ask, Last float64
}

// OptionStrike is a strike with its call and put.
type OptionStrike struct {
	Strike    float64
	Call, Put OptionQuote
}

type OptionSide int

const (
	OptionCall OptionSide = iota
	OptionPut
)

// OptionChainState is the expiry shown and the chosen contract.
type OptionChainState struct {
	Expiry string
	Strike float64
	Side   OptionSide
	Chosen bool
	rows   ui.ListState
	cursor int
	center string
}

// OptionChainOptions configure an OptionChain.
type OptionChainOptions struct {
	Spot     float64
	Decimals int
	Disabled bool
}

// OptionChainView is a built OptionChain and this frame's events.
type OptionChainView struct {
	Element *ui.Element
	// ExpiryChanged asks the caller for the strikes of State.Expiry.
	ExpiryChanged      bool
	Changed, Submitted bool
}

func OptionChain(c *ui.Context, s *OptionChainState, expiries []string, strikes []OptionStrike, opts OptionChainOptions) OptionChainView {
	if s == nil || len(expiries) == 0 {
		panic("mujica: OptionChain needs a state and an expiry")
	}
	if opts.Spot < 0 || math.IsNaN(opts.Spot) || math.IsInf(opts.Spot, 0) || s.Side < OptionCall || s.Side > OptionPut {
		panic("mujica: OptionChain price or side is invalid")
	}
	for i, o := range strikes {
		if i > 0 && o.Strike == strikes[i-1].Strike {
			panic("mujica: OptionChain strikes must be unique")
		}
		for _, v := range []float64{o.Strike, o.Call.Bid, o.Call.Ask, o.Call.Last, o.Put.Bid, o.Put.Ask, o.Put.Last} {
			if v < 0 || math.IsNaN(v) || math.IsInf(v, 0) {
				panic("mujica: OptionChain quotes must be finite and nonnegative")
			}
		}
	}
	if !slices.IsSortedFunc(strikes, func(a, b OptionStrike) int { return cmp.Compare(a.Strike, b.Strike) }) {
		panic("mujica: OptionChain strikes must ascend")
	}

	dec := core.PriceDecimals("OptionChain", opts.Decimals)
	k := core.Tokens(c)
	d := core.Density(c)
	var v OptionChainView
	v.Element = ui.Column(c).Gap(d.Unit() * 2).Grow(1).MinHeight(0).Disabled(opts.Disabled)
	opts.Disabled = v.Element.IsDisabled()
	if opts.Disabled {
		copy := *s
		s = &copy
	}
	if s.Expiry == "" {
		s.Expiry = expiries[0]
	} else if !slices.Contains(expiries, s.Expiry) {
		panic("mujica: OptionChain expiry is not listed")
	}
	v.Element.Children(func() {
		ui.Row(c).Gap(d.Unit() * 3).AlignItems(ui.Center).Children(func() {
			items := make([]input.SelectOption[string], len(expiries))
			for i, e := range expiries {
				items[i] = input.SelectOption[string]{Value: e, Label: e}
			}
			v.ExpiryChanged = input.Select(c, &s.Expiry, items, input.SelectOptions{Label: core.Msg(c, "option_chain.expiry"), Disabled: opts.Disabled}).Changed()
			if opts.Spot > 0 {
				ui.Text(c, fmt.Sprintf(core.Msg(c, "option_chain.spot"), data.StatisticFormat(opts.Spot, dec))).TextColor(k.TextMuted)
			}
		})
		optionChainHeader(c)
		v.Changed, v.Submitted = optionChainRows(c, s, strikes, opts, dec)
	})
	if v.Element.IsDisabled() {
		v.ExpiryChanged, v.Changed, v.Submitted = false, false, false
	}
	if v.Changed {
		slog.Debug("mujica: option_chain choose", "strike", s.Strike, "side", s.Side)
	}
	return v
}

var optionChainCols = [3]string{"option_chain.bid", "option_chain.ask", "option_chain.last"}

func optionChainHeader(c *ui.Context) {
	t := core.Tokens(c)
	d := core.Density(c)
	size := core.FontSize(c, theme.CaptionSize)
	ui.Row(c).Padding(0, d.Unit()*3).Children(func() {
		ui.Text(c, core.Msg(c, "option_chain.calls")).Grow(1).Basis(0).TextAlign(ui.Center).FontSize(size).TextColor(t.TextMuted)
		ui.Box(c).Width(90)
		ui.Text(c, core.Msg(c, "option_chain.puts")).Grow(1).Basis(0).TextAlign(ui.Center).FontSize(size).TextColor(t.TextMuted)
	})
	ui.Row(c).Padding(0, d.Unit()*3).Children(func() {
		for side, title := range []string{"option_chain.calls", "option_chain.puts"} {
			if side == 1 {
				ui.Text(c, core.Msg(c, "option_chain.strike")).Width(90).TextAlign(ui.Center).FontSize(size).TextColor(t.TextMuted)
			}
			ui.Row(c).Grow(1).Basis(0).Children(func() {
				for _, col := range optionChainCols {
					ui.Text(c, core.Msg(c, col)).Grow(1).Basis(0).TextAlign(ui.End).FontSize(size).TextColor(t.TextMuted)
				}
			}).Label(core.Msg(c, title))
		}
	})
}

func optionChainRows(c *ui.Context, s *OptionChainState, strikes []OptionStrike, opts OptionChainOptions, dec int) (changed, submitted bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	if len(strikes) == 0 {
		ui.Text(c, core.Msg(c, "option_chain.empty")).TextColor(k.TextMuted).Padding(d.Unit() * 3)
		return false, false
	}
	s.cursor = -1
	if s.Chosen {
		s.cursor = slices.IndexFunc(strikes, func(o OptionStrike) bool { return o.Strike == s.Strike })
	}
	before, side := s.cursor, s.Side
	cursor := &s.cursor
	if opts.Disabled {
		cp := s.cursor
		cursor = &cp
	}
	atm, _ := slices.BinarySearchFunc(strikes, opts.Spot, func(o OptionStrike, v float64) int { return cmp.Compare(o.Strike, v) })
	if atm == len(strikes) || atm > 0 && opts.Spot-strikes[atm-1].Strike < strikes[atm].Strike-opts.Spot {
		atm--
	}
	if key := fmt.Sprint(s.Expiry, len(strikes)); key != s.center {
		s.center = key
		s.rows.ScrollTo(min(atm, len(strikes)-1), ui.Center)
	}
	s.rows.Key = func(i int) any { return strikes[i].Strike }
	s.rows.Label = func(i int) string { return data.StatisticFormat(strikes[i].Strike, dec) }
	s.rows.Selected = cursor
	clickedRow, clickedSide := -1, OptionCall
	double := false
	restore := data.DataSelectTheme(c, 0)
	list := ui.List(c, &s.rows, len(strikes), func(i int) {
		o := strikes[i]
		chosen := s.Chosen && i == s.cursor
		ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).AlignItems(ui.Stretch).Children(func() {
			for sd, q := range []OptionQuote{o.Call, o.Put} {
				if sd == 1 {
					optionChainStrike(c, o.Strike, dec, i == atm && opts.Spot > 0)
				}
				itm := opts.Spot > 0 && (sd == 0 && o.Strike < opts.Spot || sd == 1 && o.Strike > opts.Spot)
				half := ui.Row(c).Grow(1).Basis(0).AlignItems(ui.Center).Radius(theme.ControlRadius).
					Label(optionChainName(c, o.Strike, OptionSide(sd), q, dec, itm))
				if half.Clicked() && !half.IsDisabled() {
					clickedRow, clickedSide, double = i, OptionSide(sd), half.DoubleClicked()
				}
				if itm {
					half.Background(k.Selection)
				}
				half.DrawOver(func(p *ui.Painter, r ui.Rect) {
					if itm {
						p.FillPath(internal.Diamond(internal.Snap(p, r.X+5), internal.Snap(p, r.Y+r.H/2), 6), k.AccentText)
					}
					if chosen && OptionSide(sd) == s.Side {
						internal.Frame(p, r, 1, 5, k.AccentText)
					}
				})
				half.Children(func() {
					for _, v := range []float64{q.Bid, q.Ask, q.Last} {
						ui.Text(c, data.StatisticFormat(v, dec)).Grow(1).Basis(0).TextAlign(ui.End).FontFeatures("tnum").SingleLine()
					}
				})
			}
		})
	}).Grow(1).MinHeight(0).Padding(d.Unit()).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
	restore()
	internal.Ring(list, k.Focus, theme.ControlRadius)
	pointed := OptionSide(-1)
	if x, _, over := list.PointerPosition(); over {
		pointed = OptionCall
		if x > list.Bounds().W/2 {
			pointed = OptionPut
		}
	}
	if !list.IsDisabled() {
		if clickedRow >= 0 {
			s.cursor, s.Side = clickedRow, clickedSide
			list.Focus()
		}
		switch {
		case clickedRow < 0 && s.cursor != before && pointed >= 0:
			s.Side = pointed
		case list.Shortcut(0, ui.KeyLeft):
			s.Side = OptionCall
		case list.Shortcut(0, ui.KeyRight):
			s.Side = OptionPut
		}
		if s.cursor >= 0 && (s.cursor != before || s.Side != side) {
			s.Strike, s.Chosen, changed = strikes[s.cursor].Strike, true, true
		}
		submitted = (list.Submitted() || double) && s.Chosen
	}
	return changed, submitted
}

func optionChainStrike(c *ui.Context, strike float64, dec int, atm bool) {
	k := core.Tokens(c)
	ui.Row(c).Width(90).Gap(4).Justify(ui.Center).AlignItems(ui.Center).Background(k.Background).Children(func() {
		if atm {
			ui.Box(c).Size(8, 8).Draw(func(p *ui.Painter, r ui.Rect) { p.FillPath(internal.Diamond(r.X+4, r.Y+4, 8), k.AccentText) })
		}
		ui.Text(c, data.StatisticFormat(strike, dec)).FontWeight(700).FontFeatures("tnum")
	})
}

func optionChainName(c *ui.Context, strike float64, side OptionSide, q OptionQuote, dec int, itm bool) string {
	s := fmt.Sprintf(core.Msg(c, []string{"option_chain.call_name", "option_chain.put_name"}[side]), data.StatisticFormat(strike, dec),
		data.StatisticFormat(q.Bid, dec), data.StatisticFormat(q.Ask, dec), data.StatisticFormat(q.Last, dec))
	if itm {
		s += ", " + core.Msg(c, "option_chain.itm")
	}
	return s
}
