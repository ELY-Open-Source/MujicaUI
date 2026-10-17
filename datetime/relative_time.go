package datetime

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

// RelativeTimeOptions configure a RelativeTime.
type RelativeTimeOptions struct {
	// Location is the tooltip's zone; required.
	Location *time.Location
	// Now fixes the clock and stops updates.
	Now    time.Time
	Hour12 bool
}

// RelativeTime shows t as "3 minutes ago", refreshing while shown.
func RelativeTime(c *ui.Context, t time.Time, opts RelativeTimeOptions) ui.Element {
	loc := eventLocation("RelativeTime", opts.Location)
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
		c.After(relativeTimeNext(now.Sub(t)))
	}
	text := RelativeTimeText(c, now.Sub(t))
	w := t.In(loc)
	exact := core.DateOf(w).String() + " " + timeText(c, TimeOfDay{w.Hour(), w.Minute(), 0}, timeConfig{hour12: opts.Hour12}) + " " + w.Format("MST")
	return ui.Text(c, text).Label(text + ", " + exact).Tooltip(exact).FontFeatures("tnum")
}

// RelativeTimeText words d ago (positive) or ahead (negative).
func RelativeTimeText(c *ui.Context, d time.Duration) string {
	past := d >= 0
	if !past {
		d = -d
	}
	var key string
	var n int
	switch {
	case d < 45*time.Second:
		return core.Msg(c, "relative_time.now")
	case d < 45*time.Minute:
		key, n = "minutes", max(1, int(d/time.Minute))
	case d < 22*time.Hour:
		key, n = "hours", max(1, int(d/time.Hour))
	case d < 26*24*time.Hour:
		key, n = "days", max(1, int(d/(24*time.Hour)))
	case d < 320*24*time.Hour:
		key, n = "months", max(1, int(d/(30*24*time.Hour)))
	default:
		key, n = "years", max(1, int(d/(365*24*time.Hour)))
	}
	unit := core.CalendarCount(c, "relative_time."+key, n)
	if past {
		return fmt.Sprintf(core.Msg(c, "relative_time.ago"), unit)
	}
	return fmt.Sprintf(core.Msg(c, "relative_time.in"), unit)
}

// relativeTimeNext is how long until the wording may change.
func relativeTimeNext(d time.Duration) time.Duration {
	if d < 0 {
		d = -d
	}
	step := time.Hour
	switch {
	case d < 2*time.Minute:
		step = time.Second
	case d < 2*time.Hour:
		step = time.Minute
	}
	return step - d%step
}
