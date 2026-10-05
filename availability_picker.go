package mujica

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// AvailabilityPerson is another person's free slots, [day][slot].
type AvailabilityPerson struct {
	Name string
	Free [][]bool
}

// AvailabilityPickerOptions configure an AvailabilityPicker.
type AvailabilityPickerOptions struct {
	// Days are the grid's columns; required.
	Days []Date
	// From and To are hours; zero means 9 to 18.
	From, To int
	// SlotMinutes must divide an hour; zero means 30.
	SlotMinutes int
	// Others are counted as free per slot.
	Others   []AvailabilityPerson
	Hour12   bool
	Disabled bool
	Label    string
}

// availabilityDrag is a paint stroke: the value being painted.
type availabilityDrag struct {
	active bool
	paint  bool
}

// AvailabilityPicker paints free slots and counts others' overlap.
func AvailabilityPicker(c *ui.Context, mine *[][]bool, opts AvailabilityPickerOptions) ChoiceResult {
	from, to := opts.From, opts.To
	if from == 0 && to == 0 {
		from, to = 9, 18
	}
	slot := timeGridSlot("AvailabilityPicker", opts.SlotMinutes)
	if len(opts.Days) == 0 || from < 0 || to > 24 || to <= from || slot > 60 {
		panic(fmt.Sprintf("mujica: AvailabilityPicker needs days and hours 0 <= %d < %d <= 24", from, to))
	}
	slots := (to - from) * 60 / slot
	availabilityCheck("value", *mine, len(opts.Days), slots)
	for _, p := range opts.Others {
		availabilityCheck(p.Name, p.Free, len(opts.Days), slots)
	}
	k := Tokens(c)
	d := Density(c)
	label := opts.Label
	if label == "" {
		label = msg(c, "availability_picker.label")
	}
	gutter := d.ControlHeight() * 1.8
	rowH := float32(theme.MinHitSize)
	headH := d.ControlHeight()
	changed := false
	set := func(day, s int, v bool) {
		if (*mine)[day][s] != v {
			(*mine)[day][s], changed = v, true
		}
	}
	grid := ui.Column(c).Focusable().Role(ui.RoleTable).Label(label).Disabled(opts.Disabled)
	internal.Ring(grid, k.Focus, theme.ControlRadius)
	drag := ui.Local(grid, "availability-drag", func() availabilityDrag { return availabilityDrag{} })
	cursor := ui.Local(grid, "availability-cursor", func() [2]int { return [2]int{} })
	if !grid.IsDisabled() {
		b := grid.Bounds()
		px, py, _ := grid.PointerPosition()
		cw := (b.W - gutter) / float32(len(opts.Days))
		day, s := int((px-gutter)/max(cw, 1)), int((py-headH)/rowH)
		inside := px >= gutter && py >= headH && day < len(opts.Days) && s < slots
		if _, _, pressing := grid.Dragged(); pressing && inside {
			if !drag.active {
				*drag = availabilityDrag{true, !(*mine)[day][s]}
				*cursor = [2]int{day, s}
			}
			set(day, s, drag.paint)
		} else if !pressing {
			drag.active = false
		}
		for key, dd := range map[ui.Key][2]int{ui.KeyLeft: {-1, 0}, ui.KeyRight: {1, 0}, ui.KeyUp: {0, -1}, ui.KeyDown: {0, 1}} {
			if grid.Shortcut(0, key) {
				*cursor = [2]int{max(0, min(len(opts.Days)-1, cursor[0]+dd[0])), max(0, min(slots-1, cursor[1]+dd[1]))}
			}
		}
		if grid.Shortcut(0, ui.KeySpace) || grid.Shortcut(0, ui.KeyEnter) {
			set(cursor[0], cursor[1], !(*mine)[cursor[0]][cursor[1]])
		}
	}
	focused := grid.FocusVisible() || grid.Focused() && !drag.active
	grid.Children(func() {
		ui.Row(c).Height(headH).Children(func() {
			ui.Box(c).Width(gutter).Shrink(0)
			for _, day := range opts.Days {
				ui.Text(c, msg(c, "calendar.weekday_"+strconv.Itoa(int(day.Weekday())))+" "+strconv.Itoa(day.Day)).Grow(1).Basis(0).
					TextAlign(ui.Center).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Role(ui.RoleColumnHeader)
			}
		})
		for s := range slots {
			m := from*60 + s*slot
			ui.Row(c).Height(rowH).Role(ui.RoleRow).Children(func() {
				hour := ""
				if m%60 == 0 {
					hour = timeText(c, TimeOfDay{m / 60, 0, 0}, timeConfig{hour12: opts.Hour12})
				}
				ui.Text(c, hour).Width(gutter).Shrink(0).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).FontFeatures("tnum")
				for di, day := range opts.Days {
					free := 0
					for _, p := range opts.Others {
						if p.Free[di][s] {
							free++
						}
					}
					availabilityCell(c, day, m, (*mine)[di][s], free, len(opts.Others), focused && *cursor == [2]int{di, s}, opts.Hour12)
				}
			})
		}
	})
	if changed {
		slog.Debug("mujica: availability_picker change")
	}
	return ChoiceResult{Element: grid, changed: changed}
}

func availabilityCheck(name string, v [][]bool, days, slots int) {
	if len(v) != days {
		panic(fmt.Sprintf("mujica: AvailabilityPicker %s has %d days, want %d", name, len(v), days))
	}
	for _, row := range v {
		if len(row) != slots {
			panic(fmt.Sprintf("mujica: AvailabilityPicker %s has %d slots, want %d", name, len(row), slots))
		}
	}
}

// availabilityCell shows one slot's choice and free count.
func availabilityCell(c *ui.Context, day Date, m int, mine bool, free, others int, cursor, hour12 bool) {
	k := Tokens(c)
	text := datetimeDayLabel(c, day) + " " + timeText(c, TimeOfDay{m / 60, m % 60, 0}, timeConfig{hour12: hour12})
	if mine {
		text += ", " + msg(c, "availability_picker.mine")
	}
	if others > 0 {
		text += ", " + fmt.Sprintf(msg(c, "availability_picker.others"), free, others)
	}
	everyone := mine && others > 0 && free == others
	cell := ui.Box(c).Grow(1).Basis(0).MinWidth(0).FillHeight().Center().Role(ui.RoleCell).Label(text).BorderWidth(0, 0, theme.BorderWidth, theme.BorderWidth).BorderColor(k.Border)
	if mine {
		cell.Background(k.Selection)
	}
	cell.Draw(func(p *ui.Painter, r ui.Rect) {
		if others > 0 && free > 0 {
			h := (r.H - 4) * float32(free) / float32(others)
			p.Fill(ui.Rect{X: r.X + 2, Y: r.Y + r.H - 2 - h, W: 3, H: h}, k.Info, 1)
		}
		if mine {
			internal.Mark(p, ui.Rect{X: r.X + r.W - 16, Y: r.Y + 2, W: 12, H: 12}, "check", k.AccentText)
		}
		if cursor {
			datetimeSelectedMark(p, r, k.Focus)
		}
	})
	cell.Children(func() {
		if others > 0 && free > 0 {
			t := ui.Text(c, strconv.Itoa(free)).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).FontFeatures("tnum")
			if everyone {
				t.Bold().TextColor(k.AccentText)
			}
		}
	})
}
