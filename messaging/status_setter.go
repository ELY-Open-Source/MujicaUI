package messaging

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CustomStatus stores content and an optional expiry.
type CustomStatus struct {
	Emoji, Text string
	Until       time.Time
}

// StatusSetterOptions configure a StatusSetter.
type StatusSetterOptions struct {
	// Location defines day boundaries; nil uses local time.
	Location *time.Location
	Disabled bool
}

var statusClear = []struct {
	key   string
	after time.Duration
}{
	{"status_setter.never", 0}, {"status_setter.30m", 30 * time.Minute}, {"status_setter.1h", time.Hour},
	{"status_setter.4h", 4 * time.Hour}, {"status_setter.today", -1}, {"status_setter.week", -2},
}

type statusSetterState struct {
	picking bool
	clear   int
	until   time.Time
}

type statusSetterKey struct{}

// StatusSetter edits *st's emoji and text, and when it clears.
func StatusSetter(c *ui.Context, st *CustomStatus, opts StatusSetterOptions) core.ChoiceResult {
	loc := opts.Location
	if loc == nil {
		loc = time.Local
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	old := *st
	col := ui.Column(c).Gap(u * 2).Role(ui.RoleGroup).Label(core.Msg(c, "status_setter.label")).Disabled(opts.Disabled)
	ss := ui.Local(col, statusSetterKey{}, func() *statusSetterState { return &statusSetterState{} })
	s := *ss
	live := !col.IsDisabled()
	if !live {
		copy := *st
		st = &copy
	}
	if !st.Until.Equal(s.until) {
		s.until = st.Until
		s.clear = 0
		if !st.Until.IsZero() {
			s.clear = -1
		}
	}
	col.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			face := st.Emoji
			if face == "" {
				face = "🙂"
			}
			b := ui.ButtonBase(c).Size(core.Density(c).ControlHeight(), core.Density(c).ControlHeight()).Radius(theme.ControlRadius).
				Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).Label(core.Msg(c, "status_setter.emoji"))
			internal.Ring(b, k.Focus, theme.ControlRadius)
			if b.Hovered() {
				b.Background(k.SurfaceHover)
			}
			b.Children(func() { ui.Text(c, face).FontSize(core.FontSize(c, 18)) })
			if b.Clicked() && !b.IsDisabled() {
				s.picking = !s.picking
			}
			if b.IsDisabled() {
				s.picking = false
			}
			overlay.Popover(c, b, &s.picking, overlay.PopoverOptions{Label: core.Msg(c, "emoji_picker.label")}, func() {
				if e := EmojiPicker(c, EmojiPickerOptions{}); e.Picked != "" {
					st.Emoji, s.picking = e.Picked, false
				}
			})
			input.TextInput(c, &st.Text, input.TextInputOptions{Placeholder: core.Msg(c, "status_setter.placeholder"), Label: core.Msg(c, "status_setter.text"),
				Clearable: true, Disabled: opts.Disabled}).Element.Grow(1)
		})
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, core.Msg(c, "status_setter.clearAfter")).TextColor(k.TextMuted).SingleLine().Shrink(0)
			options := make([]input.SelectOption[int], len(statusClear))
			for i, o := range statusClear {
				options[i] = input.SelectOption[int]{Value: i, Label: core.Msg(c, o.key)}
			}
			ui.Box(c.Key(live)).Children(func() {
				if input.Select(c, &s.clear, options, input.SelectOptions{Label: core.Msg(c, "status_setter.clearAfter"), Placeholder: core.MessagingDate(c, st.Until.In(loc), "status_setter.until"), Disabled: !live}).Changed() {
					st.Until = statusUntil(c.Now().In(loc), statusClear[s.clear].after)
				}
			})
			if !st.Until.IsZero() {
				ui.Text(c, core.MessagingDate(c, st.Until.In(loc), "status_setter.until")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			}
			ui.Box(c).Grow(1)
			if input.Button(c, core.Msg(c, "status_setter.clear"), input.ButtonOptions{Variant: input.Ghost, Disabled: *st == CustomStatus{}}).Clicked() {
				*st, s.clear = CustomStatus{}, 0
			}
		})
	})
	s.until = st.Until
	changed := live && *st != old
	if changed {
		slog.Debug("mujica: status_setter changed", "emoji", st.Emoji, "until", st.Until)
	}
	return core.NewChoiceResult(col, changed)
}

func statusUntil(now time.Time, d time.Duration) time.Time {
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	switch {
	case d == 0:
		return time.Time{}
	case d == -1:
		return day.AddDate(0, 0, 1)
	case d == -2:
		// Weeks end on Sunday night: the next Monday at midnight.
		add := (8 - int(day.Weekday())) % 7
		if add == 0 {
			add = 7
		}
		return day.AddDate(0, 0, add)
	}
	return now.Add(d)
}
