package mujica

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// QuickTradeButtonsOptions configure QuickTradeButtons. Symbol is required.
type QuickTradeButtonsOptions struct {
	Symbol string
	// Presets are the quantities offered; nil means 1, 10, 100.
	Presets []float64
	// Decimals of the quantities shown.
	Decimals int
	// Confirm asks before sending, by the button pressed.
	Confirm  bool
	Disabled bool
}

// QuickTradeResult is built QuickTradeButtons and their trade.
type QuickTradeResult struct {
	Element  *ui.Element
	side     OrderSide
	quantity float64
}

// Traded returns the side and quantity the user traded.
func (r QuickTradeResult) Traded() (OrderSide, float64, bool) { return r.side, r.quantity, r.side != 0 }

type quickTradeState struct{ asking OrderSide }

// QuickTradeButtons trade *quantity at market with one press.
func QuickTradeButtons(c *ui.Context, quantity *float64, opts QuickTradeButtonsOptions) QuickTradeResult {
	if opts.Symbol == "" {
		panic("mujica: QuickTradeButtons needs a Symbol")
	}
	if opts.Decimals < 0 {
		panic("mujica: QuickTradeButtons with negative Decimals")
	}
	presets := opts.Presets
	if presets == nil {
		presets = []float64{1, 10, 100}
	}
	if *quantity == 0 {
		*quantity = presets[0]
	}
	at := slices.Index(presets, *quantity)
	if at < 0 {
		panic(fmt.Sprintf("mujica: QuickTradeButtons quantity %v is not one of Presets", *quantity))
	}
	u := Density(c).Unit()
	var r QuickTradeResult
	r.Element = ui.Column(c).Gap(u * 2).Disabled(opts.Disabled).Role(ui.RoleGroup).Label(msg(c, "quick_trade_buttons.label"))
	st := ui.Local(r.Element, "quick", func() quickTradeState { return quickTradeState{} })
	disabled := r.Element.IsDisabled()
	if disabled {
		st.asking = 0
	}
	qty := statisticFormat(*quantity, opts.Decimals)
	r.Element.Children(func() {
		var segs []Segment
		for _, p := range presets {
			segs = append(segs, Segment{Label: statisticFormat(p, opts.Decimals)})
		}
		if SegmentedControl(c, &at, segs, SegmentedControlOptions{Equal: true, Label: msg(c, "quick_trade_buttons.quantity"), Disabled: disabled}).Changed() {
			*quantity = presets[at]
			qty = statisticFormat(*quantity, opts.Decimals)
		}
		ui.Row(c).Gap(u * 2).Children(func() {
			for _, side := range []OrderSide{SideBuy, SideSell} {
				label := fmt.Sprintf(msg(c, fmt.Sprintf("quick_trade_buttons.side%d", side)), qty)
				b := quickTradeButton(c, label, side)
				if b.Clicked() && !disabled {
					if opts.Confirm {
						st.asking = side
					} else {
						r.side, r.quantity = side, *quantity
					}
				}
				open := st.asking == side
				res := Popconfirm(c, b, &open, PopconfirmOptions{
					Title:        fmt.Sprintf(msg(c, fmt.Sprintf("quick_trade_buttons.confirm%d", side)), qty, opts.Symbol),
					Description:  msg(c, "quick_trade_buttons.market"),
					ConfirmLabel: msg(c, fmt.Sprintf("order_entry.side%d", side)),
					Danger:       true,
				})
				if res.Confirmed && !disabled {
					r.side, r.quantity = side, *quantity
				}
				if !open && st.asking == side {
					st.asking = 0
				}
			}
		})
	})
	if r.side != 0 {
		slog.Debug("mujica: quick_trade_buttons traded", "side", r.side)
	}
	return r
}

// quickTradeButton is a large button tinted by its side.
func quickTradeButton(c *ui.Context, label string, side OrderSide) *ui.Element {
	k := Tokens(c)
	d := Density(c)
	col := trendColor(c, side.sign())
	b := ui.ButtonBase(c).Grow(1).Basis(0).MinHeight(d.ControlHeight()+d.Unit()*2).Radius(theme.ControlRadius).
		Border(theme.BorderWidth*1.5, col).TextColor(col).Background(col.Alpha(0.1)).Label(label)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case b.Pressed():
		b.Background(col.Alpha(0.28))
	case b.Hovered():
		b.Background(col.Alpha(0.18))
	}
	b.Children(func() {
		ui.Icon(c, trendIcon(side.sign())).Size(16, 16)
		ui.Text(c, label).FontWeight(700).SingleLine()
	})
	return b
}
