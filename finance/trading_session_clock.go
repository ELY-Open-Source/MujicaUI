package finance

import (
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

type TradingSession struct {
	Phase      MarketPhase
	Start, End int
}

type SessionMarket struct {
	Name     string
	Location *time.Location
	Sessions []TradingSession
	// Holiday closes the market all day.
	Holiday bool
}

// TradingSessionClockOptions configure a TradingSessionClock.
type TradingSessionClockOptions struct {
	// Now is the instant shown; zero uses the window clock.
	Now time.Time
}

func TradingSessionClock(c *ui.Context, markets []SessionMarket, opts TradingSessionClockOptions) ui.Element {
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	size := core.FontSize(c, theme.CaptionSize)
	g := ui.Grid(c).ColumnTracks(ui.FitContent(), ui.FitContent(), ui.Fr(1)).GapX(u * 4).GapY(u * 3).AlignItems(ui.Center)
	g.Children(func() {
		for _, m := range markets {
			tradingSessionCheck(m)
			local := now.In(m.Location)
			phase, next, nextPhase := tradingSessionPhase(m, local)
			ui.Column(c.Key(m.Name)).Gap(2).Children(func() {
				ui.Text(c, m.Name).Font(theme.SerifFont).FontWeight(600)
				ui.Text(c, local.Format("15:04 MST")).FontSize(size).TextColor(k.TextMuted).FontFeatures("tnum")
			})
			MarketStatus(c, phase, MarketStatusOptions{Next: next, NextPhase: nextPhase, Now: now})
			minute := float32(local.Hour()*60 + local.Minute())
			ui.Box(c).MinWidth(160).Height(22).Draw(func(p *ui.Painter, r ui.Rect) {
				bar := ui.Rect{X: r.X, Y: r.Y + 4, W: r.W, H: 10}
				p.Fill(bar, k.Background, 2)
				p.Stroke(bar, k.Border, 2, 1)
				for _, s := range m.Sessions {
					col := k.Info.Alpha(0.45)
					if s.Phase == MarketOpen {
						col = k.Success
					}
					x0, x1 := bar.X+bar.W*float32(s.Start)/1440, bar.X+bar.W*float32(s.End)/1440
					p.Fill(ui.Rect{X: x0, Y: bar.Y + 1, W: x1 - x0, H: bar.H - 2}, col, 1)
				}
				x := internal.Snap(p, bar.X+bar.W*minute/1440)
				p.FillPath(internal.Diamond(x, bar.Y+bar.H+3, 8), k.Text)
				p.Line(x, bar.Y-2, x, bar.Y+bar.H, 1.5, k.Text)
			}).Role(ui.RoleImage).Label(m.Name + " " + core.Msg(c, marketPhases[phase].key))
		}
	})
	// Times and phases move on the minute.
	if opts.Now.IsZero() && len(markets) > 0 {
		g.DrawOver(func(p *ui.Painter, _ ui.Rect) { c.After(time.Minute - time.Duration(p.Now().Second())*time.Second) })
	}
	return g
}

func tradingSessionCheck(m SessionMarket) {
	if m.Name == "" || m.Location == nil {
		panic("mujica: SessionMarket needs a Name and a Location")
	}
	for i, s := range m.Sessions {
		if s.Start < 0 || s.End > 1440 || s.Start >= s.End || i > 0 && s.Start < m.Sessions[i-1].End {
			panic("mujica: TradingSession spans must be ordered within 0..1440 minutes")
		}
		if s.Phase != MarketPreMarket && s.Phase != MarketOpen && s.Phase != MarketAfterHours {
			panic("mujica: TradingSession Phase must be pre-market, open or after-hours")
		}
	}
}

func tradingSessionPhase(m SessionMarket, local time.Time) (MarketPhase, time.Time, MarketPhase) {
	if m.Holiday || len(m.Sessions) == 0 {
		return MarketHoliday, time.Time{}, MarketClosed
	}
	mid := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, m.Location)
	at := func(min int) time.Time {
		return time.Date(mid.Year(), mid.Month(), mid.Day(), min/60, min%60, 0, 0, m.Location)
	}
	cur := local.Hour()*60 + local.Minute()
	for i, s := range m.Sessions {
		switch {
		case cur < s.Start:
			return MarketClosed, at(s.Start), s.Phase
		case cur < s.End:
			next := MarketClosed
			if i+1 < len(m.Sessions) && m.Sessions[i+1].Start == s.End {
				next = m.Sessions[i+1].Phase
			}
			return s.Phase, at(s.End), next
		}
	}
	return MarketClosed, at(1440 + m.Sessions[0].Start), m.Sessions[0].Phase
}
