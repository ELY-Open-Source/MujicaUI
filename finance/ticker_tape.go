package finance

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TickerItem is one symbol on a TickerTape.
type TickerItem struct {
	Symbol        string
	Price         float64
	ChangePercent float64
}

// TickerTapeOptions configure a TickerTape.
type TickerTapeOptions struct {
	// Decimals of the prices; zero means 2.
	Decimals int
	// Speed is in DIPs a second; zero means 40.
	Speed float32
	// PerPage is items per reduced-motion page; zero means 4.
	PerPage int
	// Label names the tape; empty uses the built-in "Market ticker".
	Label string
}

type tickerTapeState struct {
	offset float32
	last   time.Time
	page   int
	paused bool
}

// TickerTape loops quotes, pausing on hover or focus.
func TickerTape(c *ui.Context, items []TickerItem, opts TickerTapeOptions) *ui.Element {
	if opts.Speed < 0 || opts.PerPage < 0 {
		panic("mujica: TickerTape Speed and PerPage cannot be negative")
	}
	dec := core.PriceDecimals("TickerTape", opts.Decimals)
	speed, per := opts.Speed, opts.PerPage
	if speed == 0 {
		speed = 40
	}
	if per == 0 {
		per = 4
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "ticker_tape.label")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	e := ui.Row(c).MinHeight(d.ControlHeight()).AlignItems(ui.Center).Background(k.Surface).
		BorderWidth(theme.BorderWidth, 0).BorderColor(k.Border).Focusable().Role(ui.RoleGroup).Label(label)
	e.FocusRing(false)
	st := ui.Local(e, "tape", func() tickerTapeState { return tickerTapeState{} })
	if c.Preferences().ReduceMotion || len(items) == 0 {
		tickerTapePages(c, e, st, items, per, dec)
		return e
	}
	paused := e.Hovered() || e.FocusWithin()
	if paused != st.paused {
		st.paused = paused
		slog.Debug("mujica: ticker_tape paused", "paused", paused)
	}
	now := c.Now()
	if !paused && !st.last.IsZero() {
		st.offset += float32(now.Sub(st.last).Seconds()) * speed
	}
	st.last = now
	if !paused {
		c.AnimationFrame()
	}
	e.Description(core.Msg(c, "ticker_tape.pause_hint"))
	e.ClipX().Children(func() {
		strip := ui.Row(c).Absolute().Top(0).Bottom(0).AlignItems(ui.Center)
		if w := strip.Bounds().W / 2; w > 0 {
			st.offset = float32(math.Mod(float64(st.offset), float64(w)))
		}
		// Quarter steps keep the moving text on the renderer's grid.
		strip.Left(-float32(math.Round(float64(st.offset)*4)) / 4)
		strip.Children(func() {
			// The second copy fills in as the first scrolls out.
			for n := range 2 {
				r := ui.Row(c).Key(n).AlignItems(ui.Center).Shrink(0)
				if n == 1 {
					r.Role(ui.RoleNone)
				}
				r.Children(func() {
					for _, it := range items {
						tickerTapeItem(c, it, dec)
					}
				})
			}
		})
	})
	return e
}

// tickerTapePages shows one page of items, for reduced motion.
func tickerTapePages(c *ui.Context, e *ui.Element, st *tickerTapeState, items []TickerItem, per, dec int) {
	pages := max(1, (len(items)+per-1)/per)
	st.page = min(st.page, pages-1)
	turn := func(by int) {
		if !e.IsDisabled() {
			st.page = (st.page + by + pages) % pages
			slog.Debug("mujica: ticker_tape page", "page", st.page)
		}
	}
	if e.Shortcut(0, ui.KeyRight) {
		turn(1)
	}
	if e.Shortcut(0, ui.KeyLeft) {
		turn(-1)
	}
	k := core.Tokens(c)
	e.Children(func() {
		ui.Row(c).Grow(1).MinWidth(0).ClipX().AlignItems(ui.Center).Children(func() {
			for _, it := range items[st.page*per : min(len(items), (st.page+1)*per)] {
				tickerTapeItem(c, it, dec)
			}
		})
		if pages < 2 {
			return
		}
		if core.IconAction(c, icons.Must("chevron-left"), core.Msg(c, "ticker_tape.previous")).Clicked() {
			turn(-1)
		}
		ui.Text(c, fmt.Sprintf("%d / %d", st.page+1, pages)).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum")
		if core.IconAction(c, icons.Must("chevron-right"), core.Msg(c, "ticker_tape.next")).Clicked() {
			turn(1)
		}
	})
}

// tickerTapeItem is one symbol, price and change.
func tickerTapeItem(c *ui.Context, it TickerItem, dec int) {
	u := core.Density(c).Unit()
	sign := core.TrendSign(it.ChangePercent)
	ui.Row(c).Gap(u*1.5).AlignItems(ui.Center).Padding(0, u*4).Shrink(0).Children(func() {
		ui.Text(c, it.Symbol).FontWeight(600).SingleLine()
		ui.Text(c, data.StatisticFormat(it.Price, dec)).FontFeatures("tnum").SingleLine()
		ui.Row(c).Gap(u / 2).AlignItems(ui.Center).TextColor(core.TrendColor(c, sign)).Label(core.TrendWord(c, sign) + " " + data.StatisticSigned(it.ChangePercent, 2) + "%").Children(func() {
			ui.Icon(c, core.TrendIcon(sign)).Size(12, 12)
			ui.Text(c, data.StatisticSigned(it.ChangePercent, 2)+"%").FontFeatures("tnum").SingleLine()
		})
	})
}
