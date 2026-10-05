package chart

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CandleCountdownOptions configure a CandleCountdown.
type CandleCountdownOptions struct {
	// Location defaults to UTC.
	Location *time.Location
	// Now freezes the clock; zero uses the frame time.
	Now time.Time
}

// CandleCountdown displays closing time and schedules visible updates.
func CandleCountdown(c *ui.Context, interval Interval, opts CandleCountdownOptions) *ui.Element {
	loc := opts.Location
	if loc == nil {
		loc = time.UTC
	}
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
	}
	left := interval.Close(now, loc).Sub(now)
	k := core.Tokens(c)
	text := candleCountdownText(left)
	e := ui.Row(c).Gap(core.Density(c).Unit()).AlignItems(ui.Center).Role(ui.RoleText).Label(fmt.Sprintf(core.Msg(c, "candle_countdown.label"), string(interval), text))
	e.Children(func() {
		ui.Text(c, string(interval)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		ui.Text(c, text).FontFeatures("tnum").FontWeight(500).TextColor(k.Text).SingleLine()
	})
	// The next whole second, so the text never skips one.
	if opts.Now.IsZero() {
		e.DrawOver(func(_ *ui.Painter, _ ui.Rect) { c.After(left - left.Truncate(time.Second) + time.Millisecond) })
	}
	return e
}

// candleCountdownText writes d as 59:58, 3:59:58 or 2d 03:59:58.
func candleCountdownText(d time.Duration) string {
	s := int(d.Seconds())
	days, h, m, sec := s/86400, s/3600%24, s/60%60, s%60
	switch {
	case days > 0:
		return fmt.Sprintf("%dd %02d:%02d:%02d", days, h, m, sec)
	case h > 0:
		return fmt.Sprintf("%d:%02d:%02d", h, m, sec)
	}
	return fmt.Sprintf("%02d:%02d", m, sec)
}
