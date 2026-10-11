package datetime

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TimezoneSelectOptions configure a TimezoneSelect.
type TimezoneSelectOptions struct {
	// Zones are the IANA names offered; empty uses built-ins.
	Zones []string
	// Common are pinned first; empty uses a built-in short list.
	Common []string
	// Now dates the offsets; zero uses the system clock.
	Now      time.Time
	Disabled bool
	Label    string
}

var timezoneCommon = []string{"UTC", "Asia/Shanghai", "Asia/Tokyo", "Europe/London", "Europe/Paris", "America/New_York", "America/Los_Angeles"}

var timezoneWorld = []string{"Pacific/Honolulu", "America/Anchorage", "America/Los_Angeles", "America/Denver", "America/Chicago",
	"America/Mexico_City", "America/New_York", "America/Toronto", "America/Bogota", "America/Lima", "America/Caracas",
	"America/Santiago", "America/Halifax", "America/St_Johns", "America/Sao_Paulo", "America/Argentina/Buenos_Aires",
	"Atlantic/Azores", "UTC", "Europe/London", "Europe/Dublin", "Europe/Lisbon", "Africa/Lagos", "Europe/Paris",
	"Europe/Berlin", "Europe/Madrid", "Europe/Rome", "Europe/Amsterdam", "Europe/Stockholm", "Europe/Warsaw",
	"Africa/Cairo", "Africa/Johannesburg", "Europe/Athens", "Europe/Helsinki", "Europe/Kyiv", "Europe/Istanbul",
	"Asia/Jerusalem", "Europe/Moscow", "Asia/Riyadh", "Africa/Nairobi", "Asia/Tehran", "Asia/Dubai", "Asia/Kabul",
	"Asia/Karachi", "Asia/Tashkent", "Asia/Kolkata", "Asia/Kathmandu", "Asia/Dhaka", "Asia/Yangon", "Asia/Bangkok",
	"Asia/Jakarta", "Asia/Ho_Chi_Minh", "Asia/Shanghai", "Asia/Hong_Kong", "Asia/Taipei", "Asia/Singapore",
	"Asia/Manila", "Australia/Perth", "Asia/Seoul", "Asia/Tokyo", "Australia/Darwin", "Australia/Adelaide",
	"Australia/Brisbane", "Australia/Sydney", "Pacific/Guam", "Pacific/Noumea", "Pacific/Auckland", "Pacific/Fiji",
	"Pacific/Chatham", "Pacific/Tongatapu", "Pacific/Kiritimati"}

var timezoneCache sync.Map

// timezoneLoad resolves an IANA name once or panics.
func timezoneLoad(name string) *time.Location {
	if v, ok := timezoneCache.Load(name); ok {
		return v.(*time.Location)
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("mujica: TimezoneSelect zone %q: %v", name, err))
	}
	timezoneCache.Store(name, loc)
	return loc
}

// TimezoneOffset formats a zone's offset at t, as "UTC+09:00".
func TimezoneOffset(loc *time.Location, t time.Time) string {
	_, off := t.In(loc).Zone()
	sign := "+"
	if off < 0 {
		sign, off = "−", -off
	}
	return fmt.Sprintf("UTC%s%02d:%02d", sign, off/3600, off%3600/60)
}

// TimezoneSelect searches zones by name or offset into *zone.
func TimezoneSelect(c *ui.Context, zone *string, opts TimezoneSelectOptions) core.ChoiceResult {
	zones, common := opts.Zones, opts.Common
	if len(zones) == 0 {
		zones = timezoneWorld
	}
	if len(common) == 0 {
		common = timezoneCommon
	}
	for _, z := range append(append([]string{}, zones...), common...) {
		timezoneLoad(z)
	}
	if *zone != "" {
		timezoneLoad(*zone)
	}
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
	}
	k := core.Tokens(c)
	d := core.Density(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "timezone_select.label")
	}
	old := *zone
	col := ui.Column(c).Width(d.ControlHeight() * 9).Gap(d.Unit()).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	open := ui.Local(col, "timezone-open", func() bool { return false })
	query := ui.Local(col, "timezone-query", func() string { return "" })
	hi := ui.Local(col, "timezone-highlight", func() int { return 0 })
	shown := ui.Local(col, "timezone-shown", func() int { return -1 })
	text := ""
	if *zone != "" {
		text = TimezoneOffset(timezoneLoad(*zone), now) + " · " + timezoneName(*zone)
	}
	col.Children(func() {
		field := core.DatetimeField(c, ui.Row(c), false, false)
		field.Children(func() {
			t := core.DatetimeTrigger(c, text, core.Msg(c, "timezone_select.placeholder"), "clock").Label(label).Description(text)
			if !t.IsDisabled() && (t.Clicked() || !*open && t.Shortcut(0, ui.KeyDown)) {
				*open, *query, *hi = !*open, "", 0
			}
		})
		core.DatetimePopup(c, field, open, false, func() {
			list := timezoneFilter(zones, common, *query, now)
			in := input.TextInput(c, query, input.TextInputOptions{Placeholder: core.Msg(c, "timezone_select.search"), Label: core.Msg(c, "timezone_select.search")})
			in.Input.AutoFocus()
			if in.Changed() {
				*hi = 0
			}
			if !in.Input.IsDisabled() {
				switch {
				case in.Input.Shortcut(0, ui.KeyDown):
					*hi = min(*hi+1, len(list)-1)
				case in.Input.Shortcut(0, ui.KeyUp):
					*hi = max(*hi-1, 0)
				case in.Submitted() && len(list) > 0:
					*zone, *open = list[*hi].name, false
				}
			}
			defer func() { *shown = *hi }()
			_, wh := c.Size()
			ui.Scroll(c).MaxHeight(min(d.ControlHeight()*8, wh/2)).Children(func() {
				if len(list) == 0 {
					ui.Text(c, core.Msg(c, "timezone_select.empty")).TextColor(k.TextMuted).Padding(d.Unit() * 2)
				}
				for i, z := range list {
					if i == 0 || z.common != list[i-1].common {
						key := "timezone_select.all"
						if z.common {
							key = "timezone_select.common"
						}
						ui.Text(c, core.Msg(c, key)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Padding(d.Unit(), d.Unit()*2)
					}
					if timezoneRow(c, z, i == *hi, i == *hi && *hi != *shown, z.name == *zone) {
						*zone, *open = z.name, false
					}
				}
			})
		})
	})
	changed := *zone != old
	if changed {
		slog.Debug("mujica: timezone_select change", "from", old, "to", *zone)
	}
	return core.NewChoiceResult(col, changed)
}

type timezoneEntry struct {
	name, offset string
	common       bool
}

// timezoneFilter matches names and offsets, common zones first.
func timezoneFilter(zones, common []string, query string, now time.Time) []timezoneEntry {
	q := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(query, " ", "_")))
	var out []timezoneEntry
	add := func(names []string, isCommon bool) {
		for _, z := range names {
			off := TimezoneOffset(timezoneLoad(z), now)
			if q == "" || strings.Contains(strings.ToLower(z), q) || strings.Contains(strings.ToLower(off), q) {
				out = append(out, timezoneEntry{z, off, isCommon})
			}
		}
	}
	add(common, true)
	add(zones, false)
	return out
}

// timezoneName is a zone's city part, readable.
func timezoneName(z string) string {
	return strings.ReplaceAll(z, "_", " ")
}

// timezoneRow draws one zone; it reports a click.
func timezoneRow(c *ui.Context, z timezoneEntry, highlighted, reveal, chosen bool) bool {
	k := core.Tokens(c)
	d := core.Density(c)
	b := ui.ButtonBase(c).MinHeight(d.ControlHeight()).Justify(ui.Start).Gap(d.Unit() * 3).PaddingX(d.Unit() * 2).Radius(theme.ControlRadius).
		Label(z.offset + " " + timezoneName(z.name))
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case chosen || highlighted:
		b.Background(k.Selection)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	if reveal {
		b.ScrollIntoView()
	}
	b.Children(func() {
		ui.Text(c, z.offset).Width(d.ControlHeight() * 2.4).Shrink(0).FontFeatures("tnum").TextColor(k.TextMuted)
		ui.Text(c, timezoneName(z.name)).TextColor(k.Text).Grow(1).Shrink(1).SingleLine().Ellipsis("…")
		if chosen {
			ui.Box(c).Size(16, 16).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) { internal.Mark(p, r, "check", k.AccentText) })
		}
	})
	return b.Clicked() && !b.IsDisabled()
}
