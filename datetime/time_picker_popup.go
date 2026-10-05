package datetime

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// timeColumns lists choosable values per segment and reports whether a choice edited *draft.
func timeColumns(c *ui.Context, draft *TimeOfDay, cfg timeConfig, fresh bool) (picked timePart, edited bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	h := d.ControlHeight() - 4
	within := func(lo, hi TimeOfDay) bool { return hi.Compare(cfg.min) >= 0 && lo.Compare(cfg.max) <= 0 }
	ui.Row(c).Gap(d.Unit()).Children(func() {
		for _, part := range timeParts(cfg) {
			t := *draft
			var options []TimeOfDay
			switch part {
			case timeHourPart:
				for v := range 24 {
					if !cfg.hour12 || v/12 == t.Hour/12 {
						options = append(options, TimeOfDay{v, t.Minute, t.Second})
					}
				}
			case timeMinutePart:
				for v := 0; v < 60; v += cfg.minuteStep {
					options = append(options, TimeOfDay{t.Hour, v, t.Second})
				}
			case timeSecondPart:
				for v := 0; v < 60; v += cfg.secondStep {
					options = append(options, TimeOfDay{t.Hour, t.Minute, v})
				}
			case timePeriodPart:
				options = []TimeOfDay{{t.Hour % 12, t.Minute, t.Second}, {t.Hour%12 + 12, t.Minute, t.Second}}
			}
			name := core.Msg(c, timePartNames[part])
			ui.Scroll(c).Height(h * 6).Width(d.ControlHeight() * 1.5).Children(func() {
				ui.Column(c).FocusGroup(ui.Vertical).Role(ui.RoleRadioGroup).Label(name).Children(func() {
					for _, o := range options {
						var lo, hi TimeOfDay
						text := ""
						switch part {
						case timeHourPart:
							lo, hi = TimeOfDay{o.Hour, 0, 0}, TimeOfDay{o.Hour, 59, 59}
							text = fmt.Sprintf("%02d", o.Hour)
							if cfg.hour12 {
								text = fmt.Sprintf("%d", timeHour12(o.Hour))
							}
						case timeMinutePart:
							lo, hi = TimeOfDay{o.Hour, o.Minute, 0}, TimeOfDay{o.Hour, o.Minute, 59}
							text = fmt.Sprintf("%02d", o.Minute)
						case timeSecondPart:
							lo, hi = o, o
							text = fmt.Sprintf("%02d", o.Second)
						case timePeriodPart:
							lo, hi = TimeOfDay{o.Hour / 12 * 12, 0, 0}, TimeOfDay{o.Hour/12*12 + 11, 59, 59}
							text = timePeriod(c, o.Hour)
						}
						sel := t
						item := ui.RadioBase(c, &sel, o).Height(h).FillWidth().Center().Radius(theme.ControlRadius).
							Disabled(!within(lo, hi)).Label(text)
						fg := k.Text
						switch {
						case o == t:
							item.Background(k.Accent)
							fg = k.OnAccent
							if fresh {
								item.ScrollIntoView()
							}
						case item.Pressed():
							item.Background(k.Selection)
						case item.Hovered():
							item.Background(k.SurfaceHover)
						}
						internal.Ring(item, k.Focus, theme.ControlRadius)
						if o == t {
							item.Draw(func(p *ui.Painter, r ui.Rect) { core.DatetimeSelectedMark(p, r, k.OnAccent) })
						}
						item.Children(func() { ui.Text(c, text).TextColor(fg).FontFeatures("tnum") })
						if item.Clicked() {
							*draft, picked, edited = o, part, true
						}
					}
				})
			})
		}
	})
	return picked, edited
}
