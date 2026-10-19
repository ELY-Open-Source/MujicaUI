package finance

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

type MarketPhase int

const (
	MarketClosed MarketPhase = iota
	MarketPreMarket
	MarketOpen
	MarketAfterHours
	MarketHoliday
)

var marketPhases = [5]struct{ key, icon string }{
	{"market_status.closed", "lock"}, {"market_status.pre", "sun"}, {"market_status.open", "circle-check"},
	{"market_status.after", "moon"}, {"market_status.holiday", "calendar"},
}

// MarketStatusOptions configure a MarketStatus.
type MarketStatusOptions struct {
	// Market names the venue, such as "NYSE".
	Market    string
	Next      time.Time
	NextPhase MarketPhase
	// Now counts down to Next; zero uses the window clock.
	Now time.Time
}

func MarketStatus(c *ui.Context, phase MarketPhase, opts MarketStatusOptions) ui.Element {
	marketPhaseCheck(phase)
	marketPhaseCheck(opts.NextPhase)
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
	}
	word := core.Msg(c, marketPhases[phase].key)
	ink := k.TextMuted
	if phase == MarketOpen {
		ink = k.Success
	}
	label := word
	if opts.Market != "" {
		label = opts.Market + ": " + word
	}
	next := ""
	if !opts.Next.IsZero() {
		next = fmt.Sprintf(core.Msg(c, "market_status.next"), core.Msg(c, marketPhases[opts.NextPhase].key), opts.Next.Format("15:04"),
			marketStatusIn(c, opts.Next.Sub(now)))
		label += ", " + next

	}
	row := ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Role(ui.RoleStatus).Label(label)
	row.Children(func() {
		ui.Row(c).Gap(u).AlignItems(ui.Center).Padding(u/2, u*2).Radius(theme.ControlRadius).Border(theme.BorderWidth, ink).Children(func() {
			ui.Icon(c, icons.Must(marketPhases[phase].icon)).Size(14, 14).TextColor(ink)
			if opts.Market != "" {
				ui.Text(c, opts.Market).Bold().FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Text).SingleLine()
			}
			ui.Text(c, word).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(ink).FontWeight(600).SingleLine()
		})
		if next != "" {
			ui.Text(c, next).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
		}
	})
	if opts.Now.IsZero() && opts.Next.After(now) {
		row.DrawOver(func(p *ui.Painter, _ ui.Rect) { c.After(time.Minute - time.Duration(p.Now().Second())*time.Second) })
	}
	return row
}

func marketPhaseCheck(p MarketPhase) {
	if p < MarketClosed || p > MarketHoliday {
		panic("mujica: MarketPhase out of range")
	}
}

func marketStatusIn(c *ui.Context, d time.Duration) string {
	if d <= 0 {
		return core.Msg(c, "market_status.now")
	}
	m := int(d.Round(time.Minute) / time.Minute)
	if m >= 60 {
		return fmt.Sprintf(core.Msg(c, "market_status.hm"), m/60, m%60)
	}
	return fmt.Sprintf(core.Msg(c, "market_status.m"), max(m, 1))
}
