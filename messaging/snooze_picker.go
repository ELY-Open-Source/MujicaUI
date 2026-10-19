package messaging

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/datetime"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// SnoozePickerOptions configure a SnoozePicker.
type SnoozePickerOptions struct {
	// Location is the user's time zone; nil means local.
	Location *time.Location
	// MorningHour defaults to 9; valid values are 0–23.
	MorningHour int
	// Now defaults to the frame clock.
	Now      time.Time
	Disabled bool
}

// SnoozeChoice pairs a message key with a reminder time.
type SnoozeChoice struct {
	Key string
	At  time.Time
}

// SnoozeChoices returns localized-day shortcuts using morning hours.
func SnoozeChoices(now time.Time, morning int) []SnoozeChoice {
	if morning < 0 || morning > 23 {
		panic("mujica: SnoozeChoices morning outside 0..23")
	}
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	at := func(days int) time.Time {
		return time.Date(day.Year(), day.Month(), day.Day()+days, morning, 0, 0, 0, now.Location())
	}
	var out []SnoozeChoice
	if later := now.Add(3 * time.Hour).Truncate(time.Hour); later.Day() == now.Day() {
		out = append(out, SnoozeChoice{"snooze_picker.later", later})
	}
	out = append(out, SnoozeChoice{"snooze_picker.tomorrow", at(1)})
	if wd := now.Weekday(); wd != time.Saturday && wd != time.Sunday {
		out = append(out, SnoozeChoice{"snooze_picker.weekend", at(int(time.Saturday - wd))})
	}
	monday := (8 - int(now.Weekday())) % 7
	if monday == 0 {
		monday = 7
	}
	return append(out, SnoozeChoice{"snooze_picker.nextWeek", at(monday)})
}

type snoozeState struct {
	custom bool
	pick   time.Time
}

type snoozeKey struct{}

// SnoozePicker offers shortcuts and custom reminder times.
func SnoozePicker(c *ui.Context, until *time.Time, opts SnoozePickerOptions) core.ChoiceResult {
	if opts.MorningHour < 0 || opts.MorningHour > 23 {
		panic("mujica: SnoozePicker MorningHour outside 0..23")
	}
	loc := opts.Location
	if loc == nil {
		loc = time.Local
	}
	morning := opts.MorningHour
	if morning == 0 {
		morning = 9
	}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
	}
	now = now.In(loc)
	old := *until
	col := ui.Column(c).Width(300).Padding(u).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(core.Msg(c, "snooze_picker.label")).
		Disabled(opts.Disabled)
	st := ui.Local(col, snoozeKey{}, func() *snoozeState { return &snoozeState{} })
	s := *st
	col.Children(func() {
		row := func(name, when string, on bool) bool {
			b := ui.ButtonBase(c).MinHeight(d.ControlHeight()).Justify(ui.Start).Gap(u*2).Padding(0, u*3).Radius(theme.ControlRadius).
				Role(ui.RoleListItem).Label(name + ", " + when).Selected(on).Background(ui.Transparent).Transition(core.ColorFade)
			internal.Ring(b, k.Focus, theme.ControlRadius)
			switch {
			case on:
				b.Background(k.Selection).BorderWidth(0, 0, 0, 3).BorderColor(k.AccentText)
			case b.Hovered():
				b.Background(k.SurfaceHover)
			}
			b.Children(func() {
				ui.Text(c, name).Grow(1).SingleLine()
				ui.Text(c, when).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).FontFeatures("tnum").SingleLine()
			})
			return b.Clicked() && !b.IsDisabled()
		}
		for _, ch := range SnoozeChoices(now, morning) {
			ui.Box(c.Key(ch.Key)).Children(func() {
				if row(core.Msg(c, ch.Key), snoozeWhen(c, now, ch.At), until.Equal(ch.At)) {
					*until, s.custom = ch.At, false
				}
			})
		}
		ui.Box(c).Height(1).Margin(u, 0).Background(k.Border)
		if row(core.Msg(c, "snooze_picker.custom"), "", s.custom) {
			s.custom = !s.custom
			if s.pick.IsZero() {
				s.pick = SnoozeChoices(now, morning)[0].At
			}
		}
		if s.custom {
			ui.Column(c).Gap(u * 2).Padding(u * 2).Children(func() {
				datetime.DateTimePicker(c, &s.pick, datetime.DateTimePickerOptions{Location: loc, Min: now, Label: core.Msg(c, "snooze_picker.custom"), Disabled: opts.Disabled})
				if input.Button(c, core.Msg(c, "snooze_picker.set"), input.ButtonOptions{Disabled: !s.pick.After(now)}).Clicked() {
					*until = s.pick
				}
			})
		}
	})
	changed := !until.Equal(old)
	if changed {
		slog.Debug("mujica: snooze_picker set", "until", *until)
	}
	return core.NewChoiceResult(col, changed)
}

func snoozeWhen(c *ui.Context, now, t time.Time) string {
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("15:04")
	}
	return core.MessagingDate(c, t, "snooze_picker.when")
}
