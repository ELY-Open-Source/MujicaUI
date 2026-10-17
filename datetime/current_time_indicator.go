package datetime

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
)

// CurrentTimeIndicatorOptions configure a CurrentTimeIndicator.
type CurrentTimeIndicatorOptions struct {
	// Location is the zone of the day shown; required.
	Location *time.Location
	// Day is the day the parent spans; zero means today.
	Day core.Date
	// Now overrides the clock; zero uses the system clock.
	Now    time.Time
	Hour12 bool
}

// CurrentTimeIndicator rules the current time across its parent.
func CurrentTimeIndicator(c *ui.Context, opts CurrentTimeIndicatorOptions) ui.Element {
	loc := eventLocation("CurrentTimeIndicator", opts.Location)
	k := core.Tokens(c)
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
		c.After(now.Truncate(time.Minute).Add(time.Minute).Sub(now))
	}
	wall := now.In(loc)
	day := opts.Day
	if day.IsZero() {
		day = core.DateOf(wall)
	}
	core.DatetimeCheckDate("CurrentTimeIndicator Day", day)
	e := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).PassThrough()
	if core.DateOf(wall) != day {
		return e.Role(ui.RoleNone)
	}
	frac := float32(wall.Hour()*60+wall.Minute()) / (24 * 60)
	when := timeText(c, TimeOfDay{wall.Hour(), wall.Minute(), 0}, timeConfig{hour12: opts.Hour12})
	e.Role(ui.RoleStatus).Label(fmt.Sprintf(core.Msg(c, "current_time_indicator.label"), when))
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		y := internal.Snap(p, r.Y+r.H*frac)
		p.Line(r.X, y, r.X+r.W, y, 2, k.Danger)
		p.FillPath(internal.Diamond(internal.Snap(p, r.X+4), y, 8), k.Danger)
	})
	return e
}
