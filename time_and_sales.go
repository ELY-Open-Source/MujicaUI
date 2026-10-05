package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// TradePrint is one executed trade on the tape.
type TradePrint struct {
	// ID identifies the trade; required and unique.
	ID    string
	Time  time.Time
	Price float64
	Size  float64
	Side  OrderSide
}

// TimeAndSalesOptions configure a TimeAndSales.
type TimeAndSalesOptions struct {
	// Decimals of prices and sizes; zero means 2 and 0.
	Decimals, SizeDecimals int
	// LargeSize marks trades at or above it; zero marks none.
	LargeSize float64
	// Height of the list in DIPs; zero means 320.
	Height float32
	Status DataStatus
	Error  string
}

// TimeAndSalesState is the tape's follow state and list place.
type TimeAndSalesState struct {
	// Paused stops new trades from showing until resumed.
	Paused bool
	frozen int
	shown  int
	list   ListState[string]
}

// TimeAndSales lists trades newest on top; Paused holds the view.
func TimeAndSales(c *ui.Context, s *TimeAndSalesState, trades []TradePrint, opts TimeAndSalesOptions) ChoiceResult {
	dec := priceDecimals("TimeAndSales", opts.Decimals)
	if opts.SizeDecimals < 0 || opts.LargeSize < 0 || opts.Height < 0 {
		panic("mujica: TimeAndSales SizeDecimals, LargeSize and Height cannot be negative")
	}
	height := opts.Height
	if height == 0 {
		height = 320
	}
	if !s.Paused {
		s.frozen = len(trades)
	}
	s.frozen = min(s.frozen, len(trades))
	n := s.frozen
	if !s.Paused && n != s.shown {
		s.list.ScrollTo(0)
	}
	s.shown = n
	k := Tokens(c)
	d := Density(c)
	var r ChoiceResult
	r.Element = ui.Column(c).Gap(d.Unit() * 2).Height(height)
	r.Element.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, msg(c, "time_and_sales.title")).Font(theme.SerifFont).FontWeight(600).Grow(1)
			label := msg(c, "time_and_sales.pause")
			if s.Paused {
				label = msg(c, "time_and_sales.resume")
				if waiting := len(trades) - n; waiting > 0 {
					label = fmt.Sprintf(msg(c, "time_and_sales.resume_count"), waiting)
				}
			}
			if Button(c, label, ButtonOptions{Variant: Secondary}).Clicked() {
				s.Paused, r.changed = !s.Paused, true
				slog.Debug("mujica: time_and_sales follow", "paused", s.Paused)
			}
		})
		ui.Row(c).Padding(0, d.Unit()*4).Gap(d.Unit() * 2).Children(func() {
			for i, key := range []string{"time_and_sales.time", "time_and_sales.price", "time_and_sales.size"} {
				t := ui.Text(c, msg(c, key)).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Grow(1).Basis(0)
				if i > 0 {
					t.TextAlign(ui.End)
				}
			}
			ui.Box(c).Width(16)
		})
		at := func(i int) TradePrint { return trades[n-1-i] }
		List(c, &s.list, n, ListOptions[string]{
			Key: func(i int) string { return at(i).ID },
			Label: func(i int) string {
				t := at(i)
				return timeAndSalesSpoken(c, t, dec, opts)
			},
			Status: opts.Status, Error: opts.Error, Empty: msg(c, "time_and_sales.empty"),
		}, func(i int) { timeAndSalesRow(c, at(i), dec, opts) }).Element.Grow(1)
	})
	return r
}

func timeAndSalesLarge(t TradePrint, opts TimeAndSalesOptions) bool {
	return opts.LargeSize > 0 && t.Size >= opts.LargeSize
}

// timeAndSalesSpoken reads a trade with its side and size.
func timeAndSalesSpoken(c *ui.Context, t TradePrint, dec int, opts TimeAndSalesOptions) string {
	s := fmt.Sprintf("%s %s %s %s", t.Time.Format("15:04:05"), msg(c, fmt.Sprintf("order_book.side%d", t.Side)),
		statisticFormat(t.Price, dec), statisticFormat(t.Size, opts.SizeDecimals))
	if timeAndSalesLarge(t, opts) {
		s += ", " + msg(c, "time_and_sales.large")
	}
	return s
}

// timeAndSalesRow is one trade with its side and large mark.
func timeAndSalesRow(c *ui.Context, t TradePrint, dec int, opts TimeAndSalesOptions) {
	k := Tokens(c)
	large := timeAndSalesLarge(t, opts)
	sign := t.Side.sign()
	ui.Text(c, t.Time.Format("15:04:05")).FontFeatures("tnum").TextColor(k.TextMuted).Grow(1).Basis(0).SingleLine()
	ui.Row(c).Gap(2).Justify(ui.End).AlignItems(ui.Center).Grow(1).Basis(0).TextColor(trendColor(c, sign)).Children(func() {
		ui.Icon(c, trendIcon(sign)).Size(12, 12)
		ui.Text(c, statisticFormat(t.Price, dec)).FontFeatures("tnum").SingleLine()
	})
	size := ui.Text(c, statisticFormat(t.Size, opts.SizeDecimals)).FontFeatures("tnum").TextAlign(ui.End).Grow(1).Basis(0).SingleLine()
	mark := ui.Box(c).Size(16, 16).Shrink(0)
	if large {
		size.FontWeight(700)
		mark.Draw(func(p *ui.Painter, r ui.Rect) {
			p.FillPath(internal.Diamond(r.X+r.W/2, r.Y+r.H/2, 10), k.AccentText)
		})
	}
}
